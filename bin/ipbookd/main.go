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
	"os"

	"github.com/docopt/docopt-go"

	"github.com/boxofrox/ipbook/bin/ipbookd/config"
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
		s         *server.Server
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

	conf := config.Load(arguments)

	log.Printf("Listening on port %d\n", conf.Port)

	if s, err = server.New(conf.Host, conf.Port); nil != err {
		log.Fatalf("Error: %s", err)
	}

	s.Run()

	os.Exit(int(Ok))
}

func envOr(name string, def string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}
	return def
}

func version() string {
	return fmt.Sprintf("%s: version %s, commit %s, build %s\n\n", os.Args[0], VERSION, GIT_COMMIT, BUILD_DATE)
}
