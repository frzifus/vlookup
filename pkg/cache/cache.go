package cache

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

const defaultCacheDir = "/tmp/vlookup"

type Option func(*Cache) error

func WithUserDir() Option {
	return func(c *Cache) error {
		usr, err := user.Current()
		if err != nil {
			return err
		}
		if usr.HomeDir == "" {
			return fmt.Errorf("missing home directory for user %s", usr.Name)
		}
		path := fmt.Sprintf("%s/.cache/vlookup", usr.HomeDir)
		c.path = path
		return nil
	}
}

func WithPath(p string) Option {
	return func(c *Cache) error {
		f, err := os.Stat(p)
		if err != nil {
			return err
		}
		if !f.IsDir() {
			return fmt.Errorf("%s is not a directory", p)
		}
		c.path = p
		return nil
	}
}

func New(opts ...Option) (*Cache, error) {
	c := &Cache{intern: make(map[string][]byte)}
	for _, o := range opts {
		if err := o(c); err != nil {
			return nil, err
		}
	}
	if err := os.MkdirAll(c.path, os.ModePerm); err != nil {
		return nil, err
	}
	return c, c.init()
}

type Cache struct {
	path   string
	intern map[string][]byte
}

func (c *Cache) init() error {
	files, err := ioutil.ReadDir(c.path)
	if err != nil {
		return err
	}
	// NOTE: file example ML3_hajfh3j43k4j32.csv [TAG]_[HASH].csv
	// sha256([]bytes("file"))
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		b, err := ioutil.ReadFile(c.path + "/" + f.Name())
		if err != nil {
			return err
		}
		h := sha256.Sum256(b)
		s := string(h[:])
		a := strings.Split(f.Name(), "_")
		if len(a) < 1 {
			return fmt.Errorf("error: TODO")
		}
		tag, hash := a[0], strings.Join(a[1:], "_")
		if len(hash) != sha256.Size {
			return fmt.Errorf("error: TODO")
		}
		if hash != s {
			return fmt.Errorf("error: TODO")
		}
		c.intern[tag] = b
	}
	return nil
}

func (c *Cache) Set(tag string, r io.Reader) error {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return err
	}
	c.intern[tag] = buf.Bytes()
	return nil
}

func (c *Cache) Get(k string) (io.Reader, error) {
	buf, ok := c.intern[k]
	if !ok {
		return nil, fmt.Errorf("error: TODO")
	}
	return bytes.NewBuffer(buf), nil
}

func (c *Cache) Flush() error {
	for tag, payload := range c.intern {
		hash := sha256.Sum256(payload)
		name := fmt.Sprintf("%s_%s", tag, hash[:])
		if _, err := os.Stat(name); err == nil {
			continue
		}
		f, err := os.Create(fmt.Sprintf("%s/%s.csv", c.path, name))
		if err != nil {
			return err
		}
		defer f.Close()
		n, err := f.Write(payload)
		if err != nil {
			return err
		}
		if n != len(payload) {
			return fmt.Errorf("error: TODO")
		}
	}
	return nil
}
