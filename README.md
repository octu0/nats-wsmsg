# nats-wsmsg
[![MIT License](https://img.shields.io/github/license/octu0/nats-wsmsg)](https://github.com/octu0/nats-wsmsg/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/octu0/nats-wsmsg?status.svg)](https://godoc.org/github.com/octu0/nats-wsmsg)
[![Go Report Card](https://goreportcard.com/badge/github.com/octu0/nats-wsmsg)](https://goreportcard.com/report/github.com/octu0/nats-wsmsg)
[![Releases](https://img.shields.io/github/v/release/octu0/nats-wsmsg)](https://github.com/octu0/nats-wsmsg/releases)

[nats.io](https://nats.io/) based websocket message pubsub/queue server.

nats-wsmsg embeds [nats-server](https://github.com/nats-io/nats-server) and provides high performance, message (queue) server.

## Quick Start

Download latest [release](https://github.com/octu0/nats-wsmsg/releases) version appropriate for operating architecture.  
Run.

```
$ ./nats-wsmsg -p 8080
```

### example Pub/Sub

![output1](https://user-images.githubusercontent.com/42143893/50048366-70316a00-010d-11e9-8196-d84c00c0bc82.gif)

### example Pub/Queue

![output2](https://user-images.githubusercontent.com/42143893/50048371-8d663880-010d-11e9-833a-eeb3cbdcf294.gif)

see more [example](https://github.com/octu0/nats-wsmsg/tree/master/example).

## Build

Build requires Go version 1.17+ installed.

```
$ go version
```

Run `make pkg` to Build and package for linux, darwin.

```
$ git clone https://github.com/octu0/nats-wsmsg
$ make pkg
```

## Help

```
NAME:
   nats-wsmsg

USAGE:
   nats-wsmsg [global options] command [command options] [arguments...]

VERSION:
   1.4.0

COMMANDS:
     websocket  run websocket server
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-dir value  /path/to/log directory (default: "/tmp") [$WSMSG_LOG_DIR]
   --debug, -d      debug mode [$WSMSG_DEBUG]
   --verbose, -V    verbose. more message [$WSMSG_VERBOSE]
   --help, -h       show help
   --version, -v    print the version
```

### subcommand: websocket

```
NAME:
   nats-wsmsg websocket - run websocket server

USAGE:
   nats-wsmsg websocket [command options] [arguments...]

OPTIONS:
   -i value, --ip value        server bind-ip (default: "[0.0.0.0]") [$WSMSG_BIND_IP]
   -p value, --port value      server bind-port (default: "8080") [$WSMSG_BIND_PORT]
   --max-payload value         nats msg max payload size(byte) (default: 1048576) [$WSMSG_MAX_PAYLOAD]
   --http-read-timeout value   http server read timeout(seconds) (default: 10)
   --http-write-timeout value  http server write timeout(seconds) (default: 10)
   --http-idle-timeout value   http server idle timeout(seconds) (default: 15)
```
