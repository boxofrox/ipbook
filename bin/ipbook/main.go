/*
 *  ipbook          two-purpose client for ipbookd.  a) command line client, b)
 *                  update daemon.
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
 *      ipbook [-c CONFIG] [-h HOST] [-p PORT] [-t SECONDS] get NAME
 *      ipbook [-c CONFIG] [-h HOST] [-p PORT] [-t SECONDS] set NAME [IP]
 *      ipbook -? | --help
 *      ipbook -V | --version
 *
 * Options:
 *   -c, --config CONFIG    - set config file to load. (default location ~/.config/ipbook/ipbook.conf,/etc/ipbook.conf)
 *   -h, --host HOST        - set host to connect to. (default pulled from config)
 *   -p, --port PORT        - set port to connect to. (default port 7000)
 *   -t, --timeout SECONDS  - set timeout in seconds.
 *   -?, --help             - print this help message.
 *   -V, --version          - print version info.
 */
package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"

	"github.com/boxofrox/ipbook/bin/ipbook/config"
	"github.com/boxofrox/ipbook/lib/client"
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
	TimeoutError
	OtherError
)

type ProgError int

func main() {
	var (
		arguments map[string]interface{}
		c         *client.Client
		conf      config.Config
		err       error
		ip        string
	)

	usage := `IP Book Client

Usage:
    ipbook [-c CONFIG] [-h HOST] [-p PORT] [-t SECONDS] get NAME
    ipbook [-c CONFIG] [-h HOST] [-p PORT] [-t SECONDS] set NAME [IP]
    ipbook -? | --help
    ipbook -V | --version

Options:
    -c, --config CONFIG    - set config file to load.
                             (default location ~/.config/ipbook/ipbook.conf, /etc/ipbook.conf)
    -h, --host HOST        - set host to connect to. (default pulled from config)
    -p, --port PORT        - set port to connect to. (default port 7000)
    -t, --timeout SECONDS  - set timeout in seconds.
    -?, --help             - print this help message.
    -V, --version          - print version info.

    If IP is not specified, then client will upload its publicly assigned IP address.
`

	arguments, err = docopt.Parse(usage, nil, true, version(), false)

	if err != nil {
		panic(err)
	}

	conf = config.Load(arguments)

	if c, err = client.New(conf.Host, conf.Port, conf.Timeout); nil != err {
		panic(err)
	}

	switch {
	case arguments["get"].(bool):
		if ip, err = c.RequestIp(arguments["NAME"].(string)); nil != err {
			check(err)
		}

		fmt.Printf(ip)

	case arguments["set"].(bool):

		if nil == arguments["IP"] {
			if err = c.RegisterPublicIp(arguments["NAME"].(string)); nil != err {
				check(err)
			}
		} else {
			if err = c.RegisterIp(arguments["NAME"].(string), arguments["IP"].(string)); nil != err {
				check(err)
			}
		}
	}

	os.Exit(int(Ok))
}

func check(err error) {
	if err.(*client.Error).Timeout() {
		fmt.Printf("Timeout: no response.\n")
		os.Exit(int(TimeoutError))
	} else {
		fmt.Println(err)
		os.Exit(int(OtherError))
	}
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
