package pindb

import (
	"github.com/google/uuid"
)

type Client struct {
	stores *stores
}

func (c *Client) Stores() []*Store {
	return c.stores.list()
}

func (c *Client) Store(uuid uuid.UUID) (*Store, error) {
	return c.stores.get(uuid)
}

func (c *Client) Has(uuid uuid.UUID) bool {
	return c.stores.has(uuid)
}

func (c *Client) Add(token, name string) (*Store, error) {
	s, err := newStore(c, token, name)
	if err != nil {
		return nil, err
	}

	c.stores.set(s)
	return s, nil
}

func (c *Client) Read(path string) (*Store, error) {
	s, err := c.stores.read(path)
	s.client = c
	return s, err
}

func (c *Client) ReadEncrypted(path string, passphrase string) (*Store, error) {
	s, err := c.stores.readEncrypted(path, passphrase)
	s.client = c
	return s, err
}

func (c *Client) ReadBytes(data []byte) (*Store, error) {
	s, err := c.stores.readBytes(data)
	s.client = c
	return s, err
}

func (c *Client) ReadBase64(data string) (*Store, error) {
	s, err := c.stores.readBase64(data)
	s.client = c
	return s, err
}

func (c *Client) JSON() StoresJSON {
	return c.stores.json()
}

func New() *Client {
	return &Client{
		stores: newStores(),
	}
}
