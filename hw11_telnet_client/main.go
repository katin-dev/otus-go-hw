package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"

	flag "github.com/spf13/pflag"
)

func main() {
	timeout := flag.Duration("timeout", 0, "connection timeout")
	flag.Parse()

	host := strings.TrimSpace(flag.Arg(0))
	if len(host) == 0 {
		usage()
	}

	port := strings.TrimSpace(flag.Arg(1))
	if len(port) == 0 {
		usage()
	}

	t := NewTelnetClient(host+":"+port, *timeout, os.Stdin, os.Stdout)

	if err := t.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to connect to server: ", err)
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			t.Close()
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		var err error
		for {
			err = t.Send()
			if err != nil {
				t.Close()
				log.Printf("stdin closed: %s\n", err)
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		for {
			err = t.Receive()
			if err != nil {
				os.Stdin.Close()
				log.Printf("connection colsed: %s\n", err)
				return
			}
		}
	}()

	wg.Wait()

	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}

func usage() {
	fmt.Println("Usage:\n\tgo-telnet [--timeout=10] host port")
	os.Exit(0)
}
