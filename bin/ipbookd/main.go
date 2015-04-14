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

	"github.com/boxofrox/ipbook/lib/server"
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

func run(port int) {
	server, err := server.New(port)
	if nil != err {
		log.Printf("Error: %s", err)
		return
	}

	go server.Listen()

	// terminate gracefully.  ie let server finish responding to requests.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	server.Stop()
}

func version() string {
	return fmt.Sprintf("%s: version %s, commit %s, build %s\n\n", os.Args[0], VERSION, GIT_COMMIT, BUILD_DATE)
}
