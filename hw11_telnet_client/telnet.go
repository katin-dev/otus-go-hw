package main

import (
	"bufio"
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
		address, timeout, in, out, nil, nil, nil,
	}
}

// P.S. Author's solution takes no more than 50 lines.

type SimpleTelnet struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	con     net.Conn
	inScan  *bufio.Scanner
	conScan *bufio.Scanner
}

func (c *SimpleTelnet) Connect() error {
	con, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}

	c.con = con

	c.inScan = bufio.NewScanner(c.in)
	c.conScan = bufio.NewScanner(c.con)

	return nil
}

func (c *SimpleTelnet) Send() error {
	if c.con == nil {
		return fmt.Errorf("connection is not established")
	}

	var err error
	if c.inScan.Scan() {
		_, err = c.con.Write([]byte(c.inScan.Text() + "\n"))
	} else {
		c.Close()
	}

	return err
}

func (c *SimpleTelnet) Receive() error {
	if c.con == nil {
		return fmt.Errorf("connection is not established")
	}

	if c.conScan.Scan() {
		c.out.Write([]byte(c.conScan.Text() + "\n"))
		return nil
	}

	return fmt.Errorf("failed to read from remote server")
}

func (c *SimpleTelnet) Close() error {
	if c.con != nil {
		return c.con.Close()
	}

	return nil
}
