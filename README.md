# nats-wsmsg

[nats](https://nats.io/) based websocket message queue frontend.

nats-wsmsg embeds [gnatsd](https://github.com/nats-io/gnatsd) and provides high performance, portable portability, messaging capabilities.

## Quick Start

Download latest [release](https://github.com/octu0/nats-wsmsg/releases) version appropriate for operating architecture.  
Run.

```
$ ./nats-wsmsg -p 8080
```

## Build

Build requires Go version 1.11+ installed.

```
$ go version
```

Run `make pkg` to Build and package for linux, darwin.

```
$ git clone https://github.com/octu0/dynomite-floridalist
$ make pkg
```

## Help

```
NAME:
   nats-wsmsg

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -i value, --ip value     server bind-ip (default: "0.0.0.0") [$WSMSG_BIND_IP]
   -p value, --port value   server bind-port (default: 8080) [$WSMSG_BIND_PORT]
   --max-payload value      msg max payload size (default: 1048576) [$WSMSG_MAX_PAYLOAD]
   --log-dir value          /path/to/log directory (default: "/tmp") [$WSMSG_LOG_DIR]
   --procs value, -P value  attach cpu(s) (default: 8) [$WSMSG_PROCS]
   --debug, -d              debug mode [$WSMSG_DEBUG]
   --verbose, -V            verbose. more message [$WSMSG_VERBOSE]
   --help, -h               show help
   --version, -v            print the version
```
