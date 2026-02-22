// Copyright 2025 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License").
// This small server is used by the kernel stack AF_PACKET functional test
// to generate traffic scenarios (conntrack fill, listen overflow, TCP rcvbuf).
//
//go:build linux

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/sys/unix"
)

var (
	port    = flag.Int("port", 9999, "Listen port")
	backlog = flag.Int("backlog", 128, "Listen backlog (use 1 for listen-overflow scenario)")
	rcvbuf  = flag.Int("rcvbuf", 0, "SO_RCVBUF size (use small value for TCPRcvQDrop scenario)")
	hold    = flag.Int("hold", 0, "Max connections to accept and hold (0 = accept and close; use with conntrack fill)")
	sleep   = flag.Duration("read-delay", 0, "Delay between reads (use for slow drain / TCPRcvQDrop)")
)

func main() {
	flag.Parse()
	var ln net.Listener
	var err error
	if *backlog == 1 {
		ln, err = listenBacklog1(*port)
	} else {
		ln, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	}
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	var held sync.WaitGroup
	acceptLimit := *hold
	if acceptLimit <= 0 {
		acceptLimit = 1 << 30
	}
	accepted := 0
	var mu sync.Mutex
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		mu.Lock()
		accepted++
		if accepted > acceptLimit {
			mu.Unlock()
			conn.Close()
			continue
		}
		mu.Unlock()
		if *rcvbuf > 0 {
			if tcp, ok := conn.(*net.TCPConn); ok {
				_ = tcp.SetReadBuffer(*rcvbuf)
			}
		}
		if *hold > 0 {
			held.Add(1)
			go func(c net.Conn) {
				defer held.Done()
				defer c.Close()
				buf := make([]byte, 1)
				for {
					if *sleep > 0 {
						time.Sleep(*sleep)
					}
					_, err := c.Read(buf)
					if err != nil {
						return
					}
				}
			}(conn)
		} else {
			conn.Close()
		}
	}
}

// listenBacklog1 creates a TCP listener with backlog 1 (for listen-overflow scenario).
func listenBacklog1(port int) (net.Listener, error) {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}
	if err := unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
		unix.Close(fd)
		return nil, err
	}
	addr := unix.SockaddrInet4{Port: port}
	if err := unix.Bind(fd, &addr); err != nil {
		unix.Close(fd)
		return nil, err
	}
	if err := unix.Listen(fd, 1); err != nil {
		unix.Close(fd)
		return nil, err
	}
	f := os.NewFile(uintptr(fd), "listener")
	// FileListener takes ownership; do not close f here.
	return net.FileListener(f)
}
