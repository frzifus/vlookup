package arp

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/mdlayher/arp"
	"github.com/mdlayher/ethernet"
)

// NewDiscovery creates a new arp Discovery service for the given interface
func NewDiscovery(iface *net.Interface) (*Discovery, error) {
	c, err := arp.Dial(iface)
	if err != nil {
		return nil, err
	}
	addresses, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	targets := make([]net.IP, 0)
	for _, a := range addresses {
		t, err := hosts(a.String())
		if err != nil {
			return nil, err
		}
		targets = append(targets, t...)
	}

	ips := make(map[string]net.IP, 0)
	for _, a := range addresses {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if x := ipnet.IP.To4(); x != nil {
				ips[ipnet.String()] = ipnet.IP
			}
		}
	}

	return &Discovery{
		client:       c,
		sendTimeout:  10 * time.Millisecond,
		wTimeout:     2 * time.Second,
		sendInterval: 5 * time.Second,
		rTimeout:     10 * time.Second,
		myAddresses:  ips,
		targets:      targets,
		hosts:        make(chan net.IP),
		iface:        iface,
	}, nil
}

type arpClient interface {
	Request(net.IP) error
	Read() (*arp.Packet, *ethernet.Frame, error)
	SetWriteDeadline(time.Time) error
	Close() error
}

// Discovery is used to locate devices on the network using the Address
// Resolution Protocol (ARP). Send and read timeouts can be set individually.
// The results are returned in real time via the host channel. A scan can be
// started using the "Find" method. This method blocks and ends only when
// the context has done.
// NOTE: to receive arp replies over the network interface cap_net_raw is
// required because a raw socket is used.
type Discovery struct {
	client       arpClient
	myAddresses  map[string]net.IP
	targets      []net.IP
	wTimeout     time.Duration
	rTimeout     time.Duration
	sendTimeout  time.Duration
	sendInterval time.Duration
	hosts        chan net.IP
	discovered   discoveryTable
	iface        *net.Interface
}

// Close the unix raw socket used for sending and receiving
func (a *Discovery) Close() error {
	return a.client.Close()
}

// Find device entries in the network where the initialized interface is
// located. This method blocks and returns the results via the passed entry
// channel. The process can be terminated by canceling the passed context.
func (a *Discovery) Find(ctx context.Context, response chan<- Entry) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go a.scan(ctx)
	return a.receive(ctx, response)
}

func (a *Discovery) scan(ctx context.Context) {
	for _, ip := range a.targets {
		select {
		case <-ctx.Done():
			return
		default:
		}
		// Set request deadline from flag
		if err := a.client.SetWriteDeadline(time.Now().Add(a.wTimeout)); err != nil {
			log.Printf("error: %w\n", err)
			continue
		}

		if err := a.client.Request(ip); err != nil {
			log.Printf("error: %w\n", err)
		}
		time.Sleep(a.sendTimeout)
	}
}

func (a *Discovery) receive(ctx context.Context, response chan<- Entry) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		resp, _, err := a.client.Read()
		if err != nil {
			return err
		}

		if resp.Operation != arp.OperationReply {
			log.Println("warn: invalid operation")
			continue
		}

		if _, ok := a.myAddresses[resp.SenderIP.String()]; ok {
			response <- Entry{
				Address: resp.SenderIP,
				Type:    byte(resp.HardwareType),
				Flags:   byte(resp.ProtocolType),
				Mac:     resp.SenderHardwareAddr,
				Device:  a.iface,
			}
		}
	}
}
