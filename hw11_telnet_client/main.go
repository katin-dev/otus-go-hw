package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

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

	ctx, ctxCancel := signal.NotifyContext(context.Background(), os.Interrupt)

	go func() {
		if err := t.Send(); err != nil {
			t.Close()
			ctxCancel()
			log.Printf("stdin closed: %s\n", err)
			return
		}
	}()

	go func() {
		if err := t.Receive(); err != nil {
			os.Stdin.Close()
			ctxCancel()
			log.Printf("connection colsed: %s\n", err)
			return
		}
	}()

	<-ctx.Done()

	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}

func usage() {
	fmt.Println("Usage:\n\tgo-telnet [--timeout=10] host port")
	os.Exit(0)
}
