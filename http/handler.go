package http

import (
	"log"

	"github.com/fasthttp/router"
	"github.com/nats-io/nats.go"
	"github.com/valyala/fasthttp"

	"github.com/octu0/nats-wsmsg"
	"github.com/octu0/nats-wsmsg/logger"
)

func Handler(natsUrl string, lg *logger.HttpLogger) fasthttp.RequestHandler {
	natsPool := NewNatsConnPool(natsUrl,
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Printf(
				"error: conn:%s sub:%s.%s err:%+v",
				nc.ConnectedAddr(),
				sub.Subject,
				sub.Queue,
				err,
			)
		}),
	)

	r := router.New()
	ws := r.Group("/ws")
	ws.GET("/sub/{topic}", WebsocketHandleSubscribe(natsPool))
	ws.GET("/sub/{topic}/{group}", WebsocketHandleSubscribeWithQueue(natsPool))
	r.POST("/pub/{topic}", HandlePublish(natsPool))

	r.GET("/_version", HandleVersion())
	r.GET("/_chk", HandleHealthcheck())
	r.GET("/_status", HandleStatus())
	r.GET("/", HandleStatus())

	return logging(lg, r.Handler)
}

func WebsocketHandleSubscribe(natsPool *NatsConnPool) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		nc, err := natsPool.Get()
		if err != nil {
			log.Printf("error: %+v", err)
			na(ctx)
			return
		}
		defer natsPool.Put(nc)

		topic := ctx.UserValue("topic").(string)
		if err := WebsocketSubscribe(ctx, nc, topic); err != nil {
			log.Printf("error: %+v", err)
			na(ctx)
			return
		}
	}
}

func WebsocketHandleSubscribeWithQueue(natsPool *NatsConnPool) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		nc, err := natsPool.Get()
		if err != nil {
			log.Printf("error: %+v", err)
			na(ctx)
			return
		}
		defer natsPool.Put(nc)

		topic := ctx.UserValue("topic").(string)
		group := ctx.UserValue("group").(string)
		if err := WebsocketSubscribeWithQueue(ctx, nc, topic, group); err != nil {
			log.Printf("error: %+v", err)
			na(ctx)
			return
		}
	}
}

func HandlePublish(natsPool *NatsConnPool) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		nc, err := natsPool.Get()
		if err != nil {
			log.Printf("error: %+v", err)
			na(ctx)
			return
		}
		defer natsPool.Put(nc)

		topic := ctx.UserValue("topic").(string)
		nc.Publish(topic, ctx.PostBody())
	}
}

func HandleVersion() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ok(ctx, wsmsg.Version)
	}
}

func HandleHealthcheck() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ok(ctx, `OK`)
	}
}

func HandleStatus() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ok(ctx, `OK`)
	}
}

func writeln(ctx *fasthttp.RequestCtx, text string) {
	ctx.Write([]byte(text))
	ctx.Write([]byte("\n"))
}

func na(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusNotAcceptable)

	writeln(ctx, `{"status":"not acceptable"}`)
}

func ok(ctx *fasthttp.RequestCtx, text string) {
	ctx.SetContentType("text/plain")
	ctx.SetStatusCode(fasthttp.StatusOK)

	writeln(ctx, text)
}

func logging(lg *logger.HttpLogger, handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		handler(ctx)

		lg.Accesslog(
			ctx.URI().Host(),
			ctx.URI().RequestURI(),
			ctx.Request.Header.Method(),
			ctx.Request.Header.UserAgent(),
			ctx.Response.Header.StatusCode(),
		)
	}
}
