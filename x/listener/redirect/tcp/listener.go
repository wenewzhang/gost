package tcp

import (
	"context"
	"net"
	"time"

	"github.com/wenewzhang/core/listener"
	"github.com/wenewzhang/core/logger"
	md "github.com/wenewzhang/core/metadata"
	xnet "github.com/wenewzhang/x/internal/net"
	"github.com/wenewzhang/x/internal/net/proxyproto"
	"github.com/wenewzhang/x/registry"
)

func init() {
	registry.ListenerRegistry().Register("red", NewListener)
	// registry.ListenerRegistry().Register("redir", NewListener)
	// registry.ListenerRegistry().Register("redirect", NewListener)
}

type redirectListener struct {
	ln      net.Listener
	logger  logger.Logger
	md      metadata
	options listener.Options
}

func NewListener(opts ...listener.Option) listener.Listener {
	options := listener.Options{}
	for _, opt := range opts {
		opt(&options)
	}
	return &redirectListener{
		logger:  options.Logger,
		options: options,
	}
}

func (l *redirectListener) Init(md md.Metadata) (err error) {
	if err = l.parseMetadata(md); err != nil {
		return
	}

	lc := net.ListenConfig{}
	if l.md.tproxy {
		lc.Control = l.control
	}
	network := "tcp"
	if xnet.IsIPv4(l.options.Addr) {
		network = "tcp4"
	}
	ln, err := lc.Listen(context.Background(), network, l.options.Addr)
	if err != nil {
		return err
	}

	ln = proxyproto.WrapListener(l.options.ProxyProtocol, ln, 10*time.Second)
	l.ln = ln
	return
}

func (l *redirectListener) Accept() (conn net.Conn, err error) {
	return l.ln.Accept()
}

func (l *redirectListener) Addr() net.Addr {
	return l.ln.Addr()
}

func (l *redirectListener) Close() error {
	return l.ln.Close()
}
