package http2

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/wenewzhang/core/listener"
	"github.com/wenewzhang/core/logger"
	md "github.com/wenewzhang/core/metadata"
	xnet "github.com/wenewzhang/x/internal/net"
	"github.com/wenewzhang/x/internal/net/proxyproto"
	mdx "github.com/wenewzhang/x/metadata"
	"github.com/wenewzhang/x/registry"
	"golang.org/x/net/http2"
)

func init() {
	registry.ListenerRegistry().Register("http2", NewListener)
}

type http2Listener struct {
	server  *http.Server
	addr    net.Addr
	cqueue  chan net.Conn
	errChan chan error
	logger  logger.Logger
	md      metadata
	options listener.Options
}

func NewListener(opts ...listener.Option) listener.Listener {
	options := listener.Options{}
	for _, opt := range opts {
		opt(&options)
	}
	return &http2Listener{
		logger:  options.Logger,
		options: options,
	}
}

func (l *http2Listener) Init(md md.Metadata) (err error) {
	if err = l.parseMetadata(md); err != nil {
		return
	}

	l.server = &http.Server{
		Addr:      l.options.Addr,
		Handler:   http.HandlerFunc(l.handleFunc),
		TLSConfig: l.options.TLSConfig,
	}
	if err := http2.ConfigureServer(l.server, nil); err != nil {
		return err
	}

	network := "tcp"
	if xnet.IsIPv4(l.options.Addr) {
		network = "tcp4"
	}
	ln, err := net.Listen(network, l.options.Addr)
	if err != nil {
		return err
	}
	l.addr = ln.Addr()
	ln = proxyproto.WrapListener(l.options.ProxyProtocol, ln, 10*time.Second)
	ln = tls.NewListener(
		ln,
		l.options.TLSConfig,
	)

	l.cqueue = make(chan net.Conn, l.md.backlog)
	l.errChan = make(chan error, 1)

	go func() {
		if err := l.server.Serve(ln); err != nil {
			l.logger.Error(err)
		}
	}()

	return
}

func (l *http2Listener) Accept() (conn net.Conn, err error) {
	var ok bool
	select {
	case conn = <-l.cqueue:
	case err, ok = <-l.errChan:
		if !ok {
			err = listener.ErrClosed
		}
	}
	return
}

func (l *http2Listener) Addr() net.Addr {
	return l.addr
}

func (l *http2Listener) Close() (err error) {
	select {
	case <-l.errChan:
	default:
		err = l.server.Close()
		l.errChan <- err
		close(l.errChan)
	}
	return nil
}

func (l *http2Listener) handleFunc(w http.ResponseWriter, r *http.Request) {
	raddr, _ := net.ResolveTCPAddr("tcp", r.RemoteAddr)
	conn := &conn{
		laddr:  l.addr,
		raddr:  raddr,
		closed: make(chan struct{}),
		md: mdx.NewMetadata(map[string]any{
			"r": r,
			"w": w,
		}),
	}
	select {
	case l.cqueue <- conn:
	default:
		l.logger.Warnf("connection queue is full, client %s discarded", r.RemoteAddr)
		return
	}

	<-conn.Done()
}
