# vlookup

[![Build Status](https://github.com/frzifus/vlookup/workflows/Go/badge.svg)](https://github.com/frzifus/vlookup/actions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/frzifus/vlookup)](https://github.com/frzifus/vlookup/blob/master/go.mod)
[![License](https://img.shields.io/github/license/frzifus/vlookup)](https://github.com/frzifus/vlookup/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/frzifus/vlookup.svg)](https://github.com/frzifus/vlookup/releases/)

**Network device discovery and MAC address vendor lookup via ARP**

> vlookup is a lightweight, Linux-based command-line tool that discovers devices on your local network using ARP scanning and identifies their manufacturers from IEEE vendor databases — with over 2.7 million entries embedded directly in the binary for offline use.

## Table of Contents

- [Why vlookup?](#why-vlookup)
- [Features](#features)
- [How It Works](#how-it-works)
- [Quick Start](#quick-start)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
  - [Basic Scanning](#basic-scanning)
  - [Passive Mode](#passive-mode)
  - [Interface Filtering](#interface-filtering)
  - [Custom Timeout](#custom-timeout)
  - [Save Results](#save-results)
- [Examples](#examples)
- [Command Reference](#command-reference)
  - [vlookup](#vlookup-1)
  - [crawler](#crawler)
- [Database Management](#database-management)
- [Comparison with Similar Tools](#comparison-with-similar-tools)
- [Troubleshooting](#troubleshooting)
- [Performance](#performance)
- [Security](#security)
- [FAQ](#faq)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgments](#acknowledgments)

## Why vlookup?

If you need to answer any of these questions, vlookup is for you:

- **"What devices are on my network?"** — Discover every host on your local network segment in seconds
- **"Who made this device?"** — Instantly identify the manufacturer of any device by its MAC address
- **"Is an unauthorized device on my network?"** — Audit your network and spot unknown devices
- **"Can I do this without installing a huge framework?"** — Single static binary, no runtime dependencies, works offline

vlookup stands out because it embeds the full IEEE vendor database directly in the binary. No network access is needed for vendor lookups, and the tool works entirely offline (in passive mode).

## Features

- **Active ARP Scanning**: Discover all devices on your network by sending ARP probes
- **Passive Mode**: Query system ARP cache without active scanning (no root required)
- **MAC Vendor Lookup**: Match MAC addresses to manufacturers using IEEE databases
- **Multi-Interface Support**: Scan all network interfaces or filter to specific ones
- **Embedded Databases**: Includes small, medium, and large IEEE vendor databases — no external files needed
- **Database Updates**: Download latest IEEE databases with included crawler utility
- **Custom Databases**: Support for organization-specific MAC address mappings
- **Flexible Output**: Display results in terminal and/or save to files
- **Configurable Timeouts**: Adjust scanning duration for different network sizes
- **Cross-Architecture**: Builds for AMD64 and ARM Linux systems

## How It Works

```
┌─────────────────────────────────────────────────────────────┐
│                        vlookup flow                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐    ┌───────────────┐    ┌──────────────┐  │
│  │  ARP Scan    │───▶│  ARP Cache    │───▶│  MAC Lookup  │  │
│  │  (active)    │    │  (passive)    │    │  (IEEE DB)   │  │
│  └──────────────┘    └───────────────┘    └──────────────┘  │
│         │                    │                    │          │
│         │                    │                    │          │
│         ▼                    ▼                    ▼          │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                  Merged Results                      │   │
│  │  idx | interface | IP | MAC | Vendor | Address      │   │
│  └──────────────────────────────────────────────────────┘   │
│                           │                                 │
│                           ▼                                 │
│                    ┌─────────────┐                          │
│                    │   stdout /  │                           │
│                    │   file (-o) │                           │
│                    └─────────────┘                           │
└─────────────────────────────────────────────────────────────┘
```

1. **ARP Scan** (active mode): Sends ARP request packets to every host on each network interface, collecting responses
2. **ARP Cache** (passive mode): Reads the kernel's existing ARP table — no network traffic generated
3. **MAC Lookup**: Matches each discovered MAC address against the IEEE OUI database to identify the vendor
4. **Output**: Results are printed as a formatted table and optionally saved to a file

## Quick Start

Get up and running in under a minute:

```bash
# Install (choose one method)
## Option A: Download binary from https://github.com/frzifus/vlookup/releases
## Option B: Build from source (requires Go 1.18+)
git clone https://github.com/frzifus/vlookup.git && cd vlookup && make amd64
sudo cp build/bin/vlookup-linux-amd64 /usr/local/bin/vlookup

# Scan all network interfaces (requires root)
sudo vlookup

# Or, just check your ARP cache (no root needed!)
vlookup --arp.scan=false

# Scan a specific interface
sudo vlookup -i eth0

# Save results to a file
sudo vlookup -o scan-results.txt
```

See [Installation](#installation) for more details, or [Examples](#examples) for advanced usage.

## Requirements

**System Requirements**:
- Linux kernel 2.6 or later
- IPv4 network connectivity
- AMD64 or ARM architecture

**Privileges**:
- **Active scanning**: Root privileges or `CAP_NET_RAW` capability
- **Passive mode**: No special privileges required

**Grant capability to avoid sudo**:
```bash
sudo setcap cap_net_raw+ep /usr/local/bin/vlookup
vlookup  # Now works without sudo
```

**Limitations**:
- Linux only (no Windows or macOS support)
- IPv4 only (IPv6 not supported)
- Ethernet networks only

## Installation

### Option 1: Download Binary (Recommended)

Download pre-built binaries from the [Releases page](https://github.com/frzifus/vlookup/releases):

```bash
# Download latest release (replace vX.Y.Z with actual version)
wget https://github.com/frzifus/vlookup/releases/download/vX.Y.Z/vlookup-linux-amd64
chmod +x vlookup-linux-amd64
sudo mv vlookup-linux-amd64 /usr/local/bin/vlookup

# Verify installation
vlookup --version
```

**Platform-specific notes**:
- **AMD64**: Use `vlookup-linux-amd64` (most common for servers and desktops)
- **ARM**: Use `vlookup-linux-arm` (for Raspberry Pi, embedded devices, etc.)
- **ARM64/AArch64**: Build from source with `GOARCH=arm64 make` (see Option 2)

### Option 2: Build from Source

**Prerequisites**: Go 1.18 or later

```bash
# Clone repository
git clone https://github.com/frzifus/vlookup.git
cd vlookup

# Build for your architecture
make amd64  # For AMD64
# or
make arm    # For ARM

# Install binaries
sudo cp build/bin/vlookup-linux-amd64 /usr/local/bin/vlookup
sudo cp build/bin/crawler-linux-amd64 /usr/local/bin/crawler

# Verify installation
vlookup --version
crawler --version
```

**Cross-compiling for ARM64**:
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o build/bin/vlookup-linux-arm64 ./cmd/vlookup/
```

### Option 3: Using Go Install (Developers)

```bash
# Install main tool
go install github.com/frzifus/vlookup/cmd/vlookup@latest

# Install database crawler
go install github.com/frzifus/vlookup/cmd/crawler@latest

# Ensure $GOPATH/bin is in your PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

> **Note**: Binaries installed via `go install` do not include version/commit information in `--version` output. Build from source with `make` to include build metadata.

### Verification

```bash
# Check vlookup is installed
which vlookup

# Check version (shows commit hash and build timestamp when built with make)
vlookup --version

# Test passive mode (no root required)
vlookup --arp.scan=false
```

## Usage

### Basic Scanning

Scan all active network interfaces (requires root):

```bash
sudo vlookup
```

Output format:
```
idx   interface  IP                   MAC                  Name                 Address
---   ---------  --                   ---                  ----                 -------
0     eth0       192.168.1.1          aa:bb:cc:dd:ee:ff    Cisco Systems        San Jose, CA US
1     eth0       192.168.1.100        11:22:33:44:55:66    Apple Inc            Cupertino, CA US
```

### Passive Mode

Query ARP cache without active scanning (no root required):

```bash
vlookup --arp.scan=false
```

Use this mode for:
- Forensic analysis
- Environments where active scanning is prohibited
- Quick MAC address lookups
- Running as non-root user

### Interface Filtering

Scan only a specific network interface:

```bash
sudo vlookup -i eth0
```

List available interfaces:
```bash
ip addr show
```

### Custom Timeout

Adjust scanning timeout for large or slow networks:

```bash
# Wait 30 seconds for responses (default: 10s)
sudo vlookup --arp.timeout=30s
```

### Save Results

Write results to a file:

```bash
# Output to both terminal and file
sudo vlookup -o network-scan.txt

# Suppress terminal output
sudo vlookup -o network-scan.txt > /dev/null
```

## Examples

### For Network Administrators

#### Example 1: Full Network Inventory

Discover all devices on all network interfaces:

```bash
sudo vlookup
```

**Expected Output**:
```
idx   interface  IP                   MAC                  Name                 Address
---   ---------  --                   ---                  ----                 -------
0     eth0       192.168.1.1          aa:bb:cc:dd:ee:ff    Cisco Systems        San Jose, CA US
1     eth0       192.168.1.100        11:22:33:44:55:66    Apple Inc            Cupertino, CA US
2     eth0       192.168.1.150        ff:ee:dd:cc:bb:aa    Dell Inc             Round Rock, TX US
```

#### Example 2: Periodic Audit

Save scan results for later comparison:

```bash
# Initial baseline scan
sudo vlookup -o baseline-$(date +%Y%m%d).txt

# Follow-up scan
sudo vlookup -o audit-$(date +%Y%m%d).txt

# Compare results to detect changes
diff baseline-20251201.txt audit-20251208.txt
```

**Use Case**: Track network changes, detect unauthorized devices, compliance auditing

#### Example 3: Scan Specific Interface

On multi-homed systems, target a single interface:

```bash
# Scan only wireless interface
sudo vlookup -i wlan0

# Scan only wired interface
sudo vlookup -i eth0
```

**Use Case**: Focus troubleshooting on a specific network segment without disrupting others

### For Security Professionals

#### Example 4: Passive Reconnaissance

Look up vendors for devices already in ARP cache without generating network traffic:

```bash
vlookup --arp.scan=false
```

**Use Case**: Forensic analysis, stealthy reconnaissance, read-only environments

**Note**: Results depend on existing ARP cache entries; use `ip neigh show` to view cache contents

#### Example 5: Rapid Assessment with Combined Filters

Combine options for scripted security audits:

```bash
# Quick scan, short timeout, save results
sudo vlookup -i eth0 --arp.timeout=5s -o quick-scan.txt
```

**Use Case**: Automated security scans, cron jobs, CI/CD pipeline integration

### For Developers and Home Lab Users

#### Example 6: Quick Check Without Root

Check what devices your machine has recently communicated with:

```bash
vlookup --arp.scan=false
```

**Use Case**: Quick device identification on a home network, no privileges needed

#### Example 7: Custom Vendor Database

Use your own organization-specific MAC address database:

```bash
# Create custom database (IEEE CSV format)
cat > my-vendors.csv << EOF
Registry,Assignment,Organization Name,Organization Address
MA-L,001122,Internal Device Pool,"Corporate IT, Building 1"
MA-L,AABBCC,Lab Equipment,"R&D Lab, Building 2"
EOF

# Use custom database
vlookup --arp.scan=false --src.local-file=my-vendors.csv
```

**Use Case**: Corporate environments with private MAC ranges, custom device naming

### Advanced Usage

#### Example 8: Large Network Scan

Scan a large or slow network with extended timeout:

```bash
sudo vlookup --arp.timeout=30s
```

**Performance**: /24 network typically completes in 10-15 seconds

#### Example 9: Fast Scan with Small Database

Use smaller database for faster startup:

```bash
vlookup -src embd-s
```

**Trade-off**: Faster load time (~0.5s vs ~2s) but less vendor coverage

#### Example 10: Database Update and Scan

Download the latest vendor databases and use them:

```bash
# Download all databases
crawler --src.fetch-all -o updated

# Use the updated large database
vlookup --arp.scan=false --src.local-file=2025-12-03_0_updated.csv
```

**Recommendation**: Update databases quarterly or when encountering "not found" vendors

## Command Reference

### vlookup

Network scanning and MAC address lookup tool.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-src` | string | `"embd-l"` | Database source: `ieee-s`, `ieee-m`, `ieee-l`, `embd-s`, `embd-m`, `embd-l` |
| `--src.local-file` | string | `""` | Path to local CSV database file |
| `--trim.address` | int | `40` | Maximum length of vendor address field in output |
| `--arp.scan` | bool | `true` | Enable active ARP scanning (requires root) |
| `--arp.timeout` | duration | `10s` | Maximum time to wait for ARP responses |
| `-i` | string | `""` | Filter to specific network interface (empty = all) |
| `-o` | string | `""` | Output file path (empty = stdout only) |
| `--version` | bool | `false` | Print version information and exit |

**Database Source Options**:
- `ieee-l`, `ieee-m`, `ieee-s`: Download from IEEE (large, medium, small)
- `embd-l`, `embd-m`, `embd-s`: Use embedded databases (large, medium, small)

### crawler

Database download utility for updating IEEE vendor databases.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--src.fetch-l` | bool | `false` | Download large (OUI) database (~2.7MB) |
| `--src.fetch-m` | bool | `false` | Download medium (MAM) database (~380KB) |
| `--src.fetch-s` | bool | `false` | Download small (OUI36) database (~390KB) |
| `--src.fetch-all` | bool | `false` | Download all three database sizes |
| `--src.fetch-custom` | string | `""` | Download from custom URL |
| `--timeout` | duration | `30s` | HTTP request timeout |
| `-o` | string | `""` | Output filename prefix for downloaded files |
| `--version` | bool | `false` | Print version information and exit |

**Database Sizes**:
- **Large (OUI)**: ~2.7M entries, most comprehensive
- **Medium (MAM)**: ~380K entries, moderate coverage
- **Small (OUI36)**: ~390K entries, minimal size

## Database Management

### Embedded Databases

vlookup includes three pre-loaded IEEE vendor databases:

| Size | Entries | File Size | Load Time | Use Case |
|------|---------|-----------|-----------|----------|
| Large (OUI) | ~2.7M | ~2.7MB | ~2s | Maximum coverage |
| Medium (MAM) | ~380K | ~380KB | ~0.8s | Balance coverage/speed |
| Small (OUI36) | ~390K | ~390KB | ~0.5s | Minimal size, major vendors |

**Default**: Large database (`embd-l`)

### Updating Databases

IEEE allocates new MAC address ranges monthly. Update databases quarterly:

```bash
# Download latest databases
crawler --src.fetch-all -o updated

# Verify download
ls -lh *_updated.csv

# Use updated database
vlookup --src.local-file=2025-12-03_0_updated.csv --arp.scan=false
```

### Database Sources

**Embedded** (`embd-*`):
- Compiled into binary
- No network access required
- Updated with each vlookup release

**Remote** (`ieee-*`):
- Downloaded from IEEE servers
- Always current
- Requires internet connectivity
- May be slow on first run

**Local** (`--src.local-file`):
- Custom or downloaded databases
- Organization-specific mappings
- Complete control over content

### CSV Format

Custom databases must follow IEEE format:

```csv
Registry,Assignment,Organization Name,Organization Address
MA-L,AABBCC,Example Corporation,"123 Main Street, City, State, Country ZIP"
MA-M,AABBCCD,Medium Block Corp,"456 Oak Avenue, City, Country"
MA-S,AABBCCDEF,Small Block Inc,"789 Elm Street, City, ZIP"
```

**Requirements**:
- Header row required
- Exactly 4 columns
- Registry: MA-L, MA-M, or MA-S
- Assignment: Hex MAC prefix (6, 7, or 9 characters, no colons)
- Organization Name: UTF-8 string
- Organization Address: Quoted if contains commas

## Comparison with Similar Tools

| Feature | vlookup | arp-scan | nmap |
|---------|---------|----------|------|
| ARP scanning | Yes | Yes | Yes (`-sn`) |
| MAC vendor lookup | Yes (embedded) | Yes (external file) | Yes (external DB) |
| Offline operation | Yes (passive mode) | No | Partial |
| Root required | Active mode only | Yes | Varies |
| Embedded vendor DB | Yes (~2.7M entries) | No (needs oui.txt) | No (needs nmap-mac-prefixes) |
| Custom databases | Yes | Limited | Limited |
| Binary size | Small (~10MB) | Small (~300KB) | Large (~5MB+) |
| Network impact | Minimal (ARP only) | Low (ARP only) | Varies |
| Passive mode | Yes | No | No |

**When to choose vlookup**:
- You need offline MAC vendor lookup
- You want a single self-contained binary
- You prefer not to manage external database files
- You need passive mode for read-only environments

**When to choose other tools**:
- **arp-scan**: You need fine-grained control over ARP timing and packet crafting
- **nmap**: You need port scanning, service detection, or OS fingerprinting beyond ARP

## Troubleshooting

### Error: "user has insufficient permissions"

**Cause**: Active ARP scanning requires root or CAP_NET_RAW capability

**Solutions**:

1. **Run with sudo** (quick fix):
   ```bash
   sudo vlookup
   ```

2. **Grant capability** (persistent):
   ```bash
   sudo setcap cap_net_raw+ep /usr/local/bin/vlookup
   vlookup  # Now works without sudo
   ```

3. **Use passive mode** (no privileges):
   ```bash
   vlookup --arp.scan=false
   ```

---

### No Devices Found

**Symptoms**: Scan completes but results table is empty

**Causes**:
- Network isolation (no other devices)
- Wrong interface selected
- Timeout too short
- Firewall blocking ARP

**Diagnostics**:
```bash
# Check network connectivity
ping 192.168.1.1

# List interfaces
ip addr show

# Check ARP cache
cat /proc/net/arp
```

**Solutions**:
```bash
# Scan specific interface
sudo vlookup -i eth0

# Increase timeout
sudo vlookup --arp.timeout=30s

# Try passive mode
vlookup --arp.scan=false
```

---

### Vendor Shows "not found"

**Cause**: MAC address not in database

**Reasons**:
- Newly allocated MAC range (update database)
- Locally administered MAC (won't have IEEE vendor)
- Private/custom MAC range

**Solutions**:
```bash
# Update database
crawler --src.fetch-all -o fresh

# Use updated database
vlookup --src.local-file=*_fresh.csv --arp.scan=false

# Create custom database for private MACs
# (see Database Management section)
```

**Note**: Locally administered MACs have second hex digit of 2, 6, A, or E (e.g., `02:`, `06:`, `0A:`, `0E:`)

---

### Tool Hangs or Very Slow

**Cause**: Large network or network issues

**Solutions**:
```bash
# Reduce timeout
sudo vlookup --arp.timeout=5s

# Scan specific interface
sudo vlookup -i eth0

# Use passive mode
vlookup --arp.scan=false
```

**Expected Performance**:
- /24 network (254 hosts): 10-15 seconds
- /16 network (65K hosts): Several minutes (not recommended)

---

### CSV Parsing Error

**Error**: `failed to parse CSV at line X`

**Cause**: Corrupted or invalid database file

**Solutions**:
```bash
# Re-download database
crawler --src.fetch-all -o fresh

# Verify CSV format
head -5 *_fresh.csv

# Expected format:
# Registry,Assignment,Organization Name,Organization Address
# MA-L,AABBCC,Example Corp,"123 Main St, City, Country"
```

---

### Interface Not Found

**Symptoms**: Specified interface produces no results

**Diagnostics**:
```bash
# List all interfaces (case-sensitive)
ip link show

# Bring interface up if needed
sudo ip link set eth0 up
```

**Solution**:
```bash
# Use exact interface name
sudo vlookup -i eth0  # Not ETH0 or Eth0
```

---

### Still Having Issues?

1. Check [GitHub Issues](https://github.com/frzifus/vlookup/issues) for known problems
2. Search for your specific error message
3. Report new issues with:
   - Output of `vlookup --version`
   - Full error text
   - Network environment details (Linux distro, kernel version)
   - Steps to reproduce

## Performance

### Expected Scan Times

| Network Size | Active Scan | Passive Mode |
|--------------|-------------|--------------|
| /24 (254 hosts) | 10-15 seconds | < 1 second |
| /23 (510 hosts) | 15-25 seconds | < 1 second |
| Single device | 10 seconds (timeout) | Instant |

**Note**: Times assume default 10-second timeout and responsive devices

### Memory Usage

- Embedded databases: ~10MB total
- Scan results: ~100 bytes per device
- 1000 devices: ~10.1MB total memory

### Optimization Tips

**For Large Networks**:
```bash
# Reduce timeout
sudo vlookup --arp.timeout=5s

# Scan interfaces individually
sudo vlookup -i eth0
sudo vlookup -i eth1
```

**For Speed**:
```bash
# Use smaller database
vlookup -src embd-s

# Skip active scanning
vlookup --arp.scan=false
```

### Limitations

- **Network Size**: Practical limit around 1000-2000 active devices
- **Timeout**: Minimum ~10 seconds for network round-trip
- **Concurrency**: Single-threaded; one scan at a time per interface
- **Protocols**: ARP only; no IPv6, no ICMP, no TCP/UDP

## Security

### Privilege Requirements

**Active Scanning** requires:
- Root privileges (`sudo vlookup`), OR
- `CAP_NET_RAW` capability

**Why**: Raw socket access needed to send/receive ARP packets

**Grant Capability**:
```bash
sudo setcap cap_net_raw+ep /usr/local/bin/vlookup
```

**Risks**:
- Grants raw packet access to the binary
- Binary can craft arbitrary network packets
- Only grant to trusted binaries

---

### Network Impact

**Active Scanning**:
- Sends one ARP request per IP in subnet
- /24 network = 254 ARP packets
- May trigger IDS/IPS alerts
- Generates visible network traffic

**Passive Mode**:
- Read-only access to ARP cache
- No network traffic generated
- Safe for production environments

---

### Data Privacy

**Collected Data**:
- IP addresses
- MAC addresses
- Network interface names
- Vendor information (public IEEE data)

**Storage**:
- In-memory only (unless `-o` used)
- No telemetry or external transmission
- No persistent storage by default

**Saved Files** (`-o` option):
- Contain network topology information
- Should be protected as sensitive data
- Recommend secure file permissions:
  ```bash
  sudo vlookup -o scan.txt
  chmod 600 scan.txt  # Owner read/write only
  ```

---

### Best Practices

1. **Authorization**: Only scan networks you own or have permission to scan
2. **Notification**: Inform security teams before scanning production networks
3. **Scheduling**: Avoid scanning during peak hours
4. **Storage**: Protect scan results as sensitive data
5. **Updates**: Keep vlookup and databases current
6. **Auditing**: Log scan activities for security compliance
7. **Least Privilege**: Use passive mode when active scanning isn't needed

---

### Compliance

**Legal Considerations**:
- Scanning networks you don't own may be illegal
- Corporate policies may restrict network scanning
- Compliance frameworks (PCI DSS, HIPAA) may have requirements

**Your Responsibilities**:
- Obtain proper authorization
- Follow organizational policies
- Secure scan results appropriately
- Use tool ethically and legally

## FAQ

### General

**Q: Does vlookup work on macOS or Windows?**
A: No. vlookup currently only supports Linux due to its reliance on Linux-specific ARP socket APIs.

**Q: Can I scan remote networks (not locally connected)?**
A: No. ARP operates at Layer 2 and only works on directly connected network segments. Use nmap for remote network scanning.

**Q: Does vlookup support IPv6?**
A: Not currently. vlookup only supports IPv4 ARP scanning.

**Q: What's the difference between active and passive mode?**
A: Active mode (`--arp.scan=true`, the default) sends ARP probes to discover devices. Passive mode (`--arp.scan=false`) only reads the kernel's existing ARP cache without generating traffic. Passive mode requires no root privileges.

### Database

**Q: How often should I update the vendor databases?**
A: IEEE allocates new MAC address ranges monthly. We recommend updating quarterly, or whenever you encounter "not found" vendor results.

**Q: Can I use a custom vendor database?**
A: Yes. Use `--src.local-file=path/to/your.csv` with a CSV file in IEEE format. See [Database Management](#database-management) for format details.

**Q: Which embedded database size should I use?**
A: The default `embd-l` (large) is recommended for most use cases. Use `embd-s` (small) if you need faster startup and don't need comprehensive vendor coverage.

### Troubleshooting

**Q: vlookup says "user has insufficient permissions" — what do I do?**
A: Either run with `sudo`, grant the binary `CAP_NET_RAW` capability, or use passive mode (`--arp.scan=false`). See [Troubleshooting](#troubleshooting) for details.

**Q: Why are some MAC vendors showing "not found"?**
A: The vendor may have been recently allocated by IEEE (update your database), or the MAC may be locally administered (second hex digit is 2, 6, A, or E), which won't have an IEEE vendor entry.

**Q: The scan is taking too long. How can I speed it up?**
A: Reduce the timeout with `--arp.timeout=5s`, scan a specific interface with `-i`, or use passive mode with `--arp.scan=false`.

## Contributing

We welcome contributions! Whether you're fixing a bug, adding a feature, or improving documentation, your help is appreciated.

### Development Environment Setup

**Prerequisites**:
- Go 1.18 or later
- Make (GNU or compatible)
- Git

**Clone and build**:
```bash
# Fork and clone the repository
git clone https://github.com/YOUR_USERNAME/vlookup.git
cd vlookup

# Install linting tools
make build_deps  # Installs golint

# Build the project
make amd64

# Verify everything works
./build/bin/vlookup-linux-amd64 --version
```

### Running Tests

```bash
# Run all tests with verbose output
make test

# Run tests directly with Go
go test -v ./...

# Run tests for a specific package
go test -v ./pkg/arp/
go test -v ./pkg/macpack/
```

### Linting

```bash
# Run linter (required before submitting PRs)
make lint

# Or run golint directly
golint -set_exit_status ./pkg/... ./cmd/...
```

### Development Workflow

1. **Fork** the repository on GitHub
2. **Clone** your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/vlookup.git
   cd vlookup
   git remote add upstream https://github.com/frzifus/vlookup.git
   ```
3. **Create** a feature branch:
   ```bash
   git checkout -b feature/my-feature
   ```
4. **Make** your changes
5. **Test** your changes:
   ```bash
   make test
   make lint
   ```
6. **Format** your code:
   ```bash
   gofmt -w .
   ```
7. **Commit** with clear, descriptive messages:
   ```bash
   git commit -m "Add support for XYZ feature"
   ```
8. **Push** to your fork:
   ```bash
   git push origin feature/my-feature
   ```
9. **Open** a Pull Request against the `master` branch

### Code Style

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Run `gofmt` before committing — unformatted code will not be merged
- Add tests for new functionality
- Keep functions small and focused
- Add comments for exported types and functions
- Update documentation (README, command help) as needed

### Pull Request Guidelines

- **One feature per PR**: Keep changes focused and reviewable
- **Include tests**: All new functionality should have corresponding tests
- **Pass CI**: All tests and linting must pass (`make test && make lint`)
- **Update docs**: If you change behavior, update README and command reference
- **Describe changes**: Fill out the PR template with motivation, approach, and testing steps

### Reporting Issues

Found a bug? [Open an issue](https://github.com/frzifus/vlookup/issues) with:

- **vlookup version**: Output of `vlookup --version`
- **System info**: Linux distro, kernel version (`uname -a`)
- **Network details**: Interface type (eth, wlan, etc.), network size
- **Error message**: Full error text
- **Steps to reproduce**: Exact commands run
- **Expected vs actual**: What should happen vs what does happen

### Feature Requests

Want a new feature? [Open an issue](https://github.com/frzifus/vlookup/issues) describing:

- **Use case**: What problem are you trying to solve?
- **Proposed solution**: How should it work?
- **Alternatives considered**: Other ways to achieve the goal
- **Impact**: Who benefits? How common is this need?

### Database Contributions

Custom vendor databases or IEEE URL updates:
- Verify CSV format (see [Database Management](#database-management))
- Test with sample data
- Document any special requirements

## License

vlookup is licensed under the **GNU General Public License v3.0**.

See [LICENSE](LICENSE) for the full license text.

**In Summary**:
- Use for any purpose
- Study and modify source code
- Distribute copies
- Distribute modified versions
- Must disclose source
- Must license under GPL-3.0
- Must state changes
- No warranty provided

**Commercial Use**: Permitted under GPL-3.0 terms

## Acknowledgments

**Dependencies**:
- [mdlayher/arp](https://github.com/mdlayher/arp) - ARP protocol implementation
- [mdlayher/ethernet](https://github.com/mdlayher/ethernet) - Ethernet frame handling
- [google/go-cmp](https://github.com/google/go-cmp) - Testing comparisons

**Data Sources**:
- [IEEE Registration Authority](https://standards.ieee.org/products-programs/regauth/) - MAC address vendor databases

**Inspiration**:
- [arp-scan](https://github.com/royhills/arp-scan) - Original ARP scanning utility
- [nmap](https://nmap.org/) - Network exploration and security auditing

**Contributors**:
- See [GitHub Contributors](https://github.com/frzifus/vlookup/graphs/contributors)

**Maintained By**: [frzifus](https://github.com/frzifus)