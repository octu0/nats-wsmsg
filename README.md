# nats-wsmsg

[nats](https://nats.io/) based websocket message pubsub/queue server.

nats-wsmsg embeds [gnatsd](https://github.com/nats-io/gnatsd) and provides high performance, portable portability, messaging capabilities.

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

Build requires Go version 1.11+ installed.

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
   1.2.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -i value, --ip value         server bind-ip (default: "0.0.0.0") [$WSMSG_BIND_IP]
   -p value, --port value       server bind-port (default: 8080) [$WSMSG_BIND_PORT]
   --max-payload value          msg max payload size (default: 1048576) [$WSMSG_MAX_PAYLOAD]
   --log-dir value              /path/to/log directory (default: "/tmp") [$WSMSG_LOG_DIR]
   --ws-max-message-size value  websocket max message size(byte) (default: 1048576) [$WSMSG_WS_MAX_MSG_SIZE]
   --procs value, -P value      attach cpu(s) (default: 8) [$WSMSG_PROCS]
   --debug, -d                  debug mode [$WSMSG_DEBUG]
   --verbose, -V                verbose. more message [$WSMSG_VERBOSE]
   --stdout-http-log            http-log outputs to standard out [$WSMSG_STDOUT_HTTP_LOG]
   --stdout-nats-log            nats-log outputs to standard out [$WSMSG_STDOUT_NATS_LOG]
   --help, -h                   show help
   --version, -v                print the version
```
