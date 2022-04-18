package http

import (
	"sync"

	"github.com/nats-io/nats.go"

	"github.com/octu0/nats-wsmsg"
)

type NatsConnFunc func() (*nats.Conn, error)

type NatsConnPool struct {
	connFunc NatsConnFunc
	pool     *sync.Pool
}

func (p *NatsConnPool) Get() (*nats.Conn, error) {
	v := p.pool.Get()
	if v == nil {
		return p.connFunc()
	}

	conn := v.(*nats.Conn)
	if conn.IsClosed() || conn.IsDraining() {
		return p.connFunc()
	}

	return conn, nil
}

func (p *NatsConnPool) Put(nc *nats.Conn) {
	if nc.IsConnected() {
		nc.Flush()
	}
	p.pool.Put(nc)
}

func ConnFunc(url string, opts []nats.Option) NatsConnFunc {
	return func() (*nats.Conn, error) {
		return nats.Connect(url, opts...)
	}
}

func NewNatsConnPool(url string, customOptions ...nats.Option) *NatsConnPool {
	opts := make([]nats.Option, 0, len(customOptions)+3)
	opts = append(opts, nats.DontRandomize(), nats.NoEcho(), nats.Name(wsmsg.UA))
	opts = append(opts, customOptions...)

	return &NatsConnPool{
		connFunc: ConnFunc(url, opts),
		pool:     new(sync.Pool),
	}
}
