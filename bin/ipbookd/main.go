/*
 *  ipbookd         network server used to record and provide ip address.
 *
 *  Copyright (c) 2015 Justin Charette <charetjc@gmail.com> (@boxofrox)
 *                All Rights Reserved
 *
 *  This program is free software. It comes without any warranty, to
 *  the extent permitted by applicable law. You can redistribute it
 *  and/or modify it under the terms of the Do What the Fuck You Want
 *  to Public License, Version 2, as published by Sam Hocevar. See
 *  http://www.wtfpl.net/ for more details.
 *
 *
 * Usage:
 *      ipbookd [-p port]
 *      ipbookd [-h]
 *      ipbookd [-V]
 *
 * Options:
 *   -p, --port PORT - port server will listen on.
 *   -h, --help      - print this help message.
 *   -V, --version   - print version info.
 */
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/docopt/docopt-go"

	"github.com/boxofrox/ipbook/lib/pool"
	"github.com/boxofrox/ipbook/lib/protocol"
	"github.com/boxofrox/ipbook/lib/registry"
)

var (
	VERSION    string
	GIT_COMMIT string
	BUILD_DATE string
)

const (
	Ok ProgError = iota
	ArgParseError
	ConfigParseError
)

const (
	MAX_UDP_PACKET_SIZE = 65535
)

type ProgError int

func main() {
	var (
		err       error
		arguments map[string]interface{}
		port      int
	)

	usage := `UDP Echo Server

Usage:
    ipbookd [-p PORT]
    ipbookd -h | --help
    ipbookd -V | --version

Options:
    -p, --port PORT          - set port the server listens on.
    -h, --help               - print this help message.
    -V, --version            - print version info.
`

	arguments, err = docopt.Parse(usage, nil, true, version(), false)

	if err != nil {
		log.Fatal(err)
	}

	if arguments["--port"] != nil {
		port, err = strconv.Atoi(arguments["--port"].(string))
	} else {
		port = 3000
	}

	log.Printf("Listening on port %d\n", port)

	run(port)

	os.Exit(int(Ok))
}

func createBuffer() []byte {
	return make([]byte, MAX_UDP_PACKET_SIZE)
}

func createListener(port int) (*net.UDPConn, error) {
	addr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: port,
	}
	return net.ListenUDP("udp", &addr)
}

func envOr(name string, def string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}
	return def
}

func listen(conn *net.UDPConn) {
	// create a buffer pool
	pool := pool.New(5, createBuffer)
	registry := registry.New()

	for {
		buffer := pool.GetFreeBuffer()
		n, addr, err := conn.ReadFromUDP(buffer)

		if err != nil {
			log.Printf("Error: reading udp packet from %s. %s", addr.String(), err)

			if nil != addr {
				protocol.SendErrorResponse(conn, addr, protocol.BAD_REQUEST, "unable to read request")
			}

			continue
		}

		go func(conn *net.UDPConn, addr *net.UDPAddr, n int, buffer []byte) {
			// handle request

			defer pool.Recycle(buffer)

			var err error

			object, err := protocol.Decode(buffer[0:n])
			if nil != err {
				log.Printf("Error: unable to decode request. %s", err)
				protocol.SendErrorResponse(conn, addr, protocol.BAD_REQUEST, "unable to decode request")
				return
			}

			switch object.GetType() {
			case protocol.TYPE_GET_IP_REQUEST:
				request := object.(*protocol.GetIpRequest)

				if false == registry.Contains(request.Name) {
					log.Printf("Host (%s): requested name (%s) not found in registry.",
						addr.String(), request.Name)
					protocol.SendErrorResponse(conn, addr, protocol.NAME_NOT_FOUND,
						"name not found in registry")
					return
				}

				ip, _ := registry.Get(request.Name)
				_, err = protocol.SendGetIpResponse(conn, addr, request.Name, ip)
				if nil != err {
					log.Printf("Error: unable to send get-ip response to Host (%s). %s", addr.String(), err)
				} else {
					log.Printf("Host (%s): requested IP of (%s).", addr.String(), request.Name)
				}

			case protocol.TYPE_SET_IP_REQUEST:
				request := object.(*protocol.SetIpRequest)

				if registry.Contains(request.Name) {
					log.Printf("Host (%s): changed IP of (%s) to (%s)", addr.String(), request.Name, request.Ip)
				} else {
					log.Printf("Host (%s): set IP of (%s) to (%s)", addr.String(), request.Name, request.Ip)
				}

				err = registry.Put(request.Name, request.Ip)
				if nil != err {
					log.Printf("Host (%s): unable to change IP of (%s) to (%s). %s",
						addr.String(), request.Name, request.Ip, err)

					_, err := protocol.SendErrorResponse(conn, addr, protocol.INVALID_NAME, err.Error())
					if nil != err {
						log.Printf("Error: failed to send error response to Host (%s). %s",
							addr.String(), err)
					}

					return
				}

				_, err = protocol.SendSetIpResponse(conn, addr, "ok", "")

				if nil != err {
					log.Printf("Error: failed to send set-ip response to Host (%s). %s",
						addr.String(), err)
				}
			}

		}(conn, addr, n, buffer)
	}
}

func run(port int) error {
	var (
		err  error
		conn *net.UDPConn
		sigs chan os.Signal
	)

	if conn, err = createListener(port); err != nil {
		return err
	}
	defer conn.Close()

	go listen(conn)

	sigs = make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	return nil
}

func version() string {
	return fmt.Sprintf("%s: version %s, commit %s, build %s\n\n", os.Args[0], VERSION, GIT_COMMIT, BUILD_DATE)
}
