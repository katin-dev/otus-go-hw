package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &SimpleTelnet{
		address, timeout, in, out, nil,
	}
}

// P.S. Author's solution takes no more than 50 lines.

type SimpleTelnet struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	con     net.Conn
}

func (c *SimpleTelnet) Connect() error {
	con, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}

	c.con = con

	return nil
}

func (c *SimpleTelnet) Send() error {
	if c.con == nil {
		return fmt.Errorf("connection is not established")
	}

	_, err := io.Copy(c.con, c.in)

	return err
}

func (c *SimpleTelnet) Receive() error {
	if c.con == nil {
		return fmt.Errorf("connection is not established")
	}

	_, err := io.Copy(c.out, c.con)

	return err
}

func (c *SimpleTelnet) Close() error {
	if c.con != nil {
		return c.con.Close()
	}

	return nil
}
