package redirect

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	dissector "github.com/go-gost/tls-dissector"
	"github.com/wenewzhang/core/chain"
	"github.com/wenewzhang/core/handler"
	"github.com/wenewzhang/core/logger"
	md "github.com/wenewzhang/core/metadata"
	xio "github.com/wenewzhang/x/internal/io"
	netpkg "github.com/wenewzhang/x/internal/net"
	"github.com/wenewzhang/x/registry"
)

func init() {
	registry.HandlerRegistry().Register("red", NewHandler)
	registry.HandlerRegistry().Register("redir", NewHandler)
	registry.HandlerRegistry().Register("redirect", NewHandler)
}

type redirectHandler struct {
	router  *chain.Router
	md      metadata
	options handler.Options
}

func NewHandler(opts ...handler.Option) handler.Handler {
	options := handler.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &redirectHandler{
		options: options,
	}
}

func (h *redirectHandler) Init(md md.Metadata) (err error) {
	if err = h.parseMetadata(md); err != nil {
		return
	}

	h.router = h.options.Router
	if h.router == nil {
		h.router = chain.NewRouter(chain.LoggerRouterOption(h.options.Logger))
	}

	return
}

func (h *redirectHandler) Handle(ctx context.Context, conn net.Conn, opts ...handler.HandleOption) (err error) {
	// fmt.Fprint(os.Stdout, "lite-gost x handle redirect Handle\n")
	defer conn.Close()

	var dstAddr net.Addr

	if h.md.tproxy {
		dstAddr = conn.LocalAddr()
	} else {
		dstAddr, err = h.getOriginalDstAddr(conn)
		if err != nil {
			return
		}
	}

	var rw io.ReadWriter = conn

	cc, err := h.router.Dial(ctx, dstAddr.Network(), dstAddr.String())
	if err != nil {
		return err
	}
	defer cc.Close()

	netpkg.Transport(rw, cc)
	return nil
}

func (h *redirectHandler) handleHTTP(ctx context.Context, rw io.ReadWriter, raddr net.Addr, log logger.Logger) error {
	// fmt.Fprint(os.Stdout, "lite-gost x handle redirect handleHTTP\n")
	req, err := http.ReadRequest(bufio.NewReader(rw))
	if err != nil {
		return err
	}

	if log.IsLevelEnabled(logger.TraceLevel) {
		dump, _ := httputil.DumpRequest(req, false)
		log.Trace(string(dump))
	}

	host := req.Host
	if _, _, err := net.SplitHostPort(host); err != nil {
		host = net.JoinHostPort(host, "80")
	}
	log = log.WithFields(map[string]any{
		"host": host,
	})

	cc, err := h.router.Dial(ctx, "tcp", host)
	if err != nil {
		log.Error(err)
		return err
	}
	defer cc.Close()

	t := time.Now()
	log.Infof("%s <-> %s", raddr, host)
	defer func() {
		log.WithFields(map[string]any{
			"duration": time.Since(t),
		}).Infof("%s >-< %s", raddr, host)
	}()

	if err := req.Write(cc); err != nil {
		log.Error(err)
		return err
	}

	var rw2 io.ReadWriter = cc
	if log.IsLevelEnabled(logger.TraceLevel) {
		var buf bytes.Buffer
		resp, err := http.ReadResponse(bufio.NewReader(io.TeeReader(cc, &buf)), req)
		if err != nil {
			log.Error(err)
			return err
		}
		defer resp.Body.Close()

		dump, _ := httputil.DumpResponse(resp, false)
		log.Trace(string(dump))

		rw2 = xio.NewReadWriter(io.MultiReader(&buf, cc), cc)
	}

	netpkg.Transport(rw, rw2)

	return nil
}

func (h *redirectHandler) handleHTTPS(ctx context.Context, rw io.ReadWriter, raddr, dstAddr net.Addr, log logger.Logger) error {
	// fmt.Fprint(os.Stdout, "lite-gost x handle redirect handleHTTPS\n")
	buf := new(bytes.Buffer)
	host, err := h.getServerName(ctx, io.TeeReader(rw, buf))
	if err != nil {
		log.Error(err)
		return err
	}
	if host == "" {
		host = dstAddr.String()
	} else {
		if _, _, err := net.SplitHostPort(host); err != nil {
			_, port, _ := net.SplitHostPort(dstAddr.String())
			if port == "" {
				port = "443"
			}
			host = net.JoinHostPort(host, port)
		}
	}

	log = log.WithFields(map[string]any{
		"host": host,
	})

	cc, err := h.router.Dial(ctx, "tcp", host)
	if err != nil {
		log.Error(err)
		return err
	}
	defer cc.Close()

	t := time.Now()
	log.Infof("%s <-> %s", raddr, host)
	netpkg.Transport(xio.NewReadWriter(io.MultiReader(buf, rw), rw), cc)
	log.WithFields(map[string]any{
		"duration": time.Since(t),
	}).Infof("%s >-< %s", raddr, host)

	return nil
}

func (h *redirectHandler) getServerName(ctx context.Context, r io.Reader) (host string, err error) {
	record, err := dissector.ReadRecord(r)
	if err != nil {
		return
	}

	clientHello := dissector.ClientHelloMsg{}
	if err = clientHello.Decode(record.Opaque); err != nil {
		return
	}

	for _, ext := range clientHello.Extensions {
		if ext.Type() == dissector.ExtServerName {
			snExtension := ext.(*dissector.ServerNameExtension)
			host = snExtension.Name
			break
		}
	}

	return
}

func isHTTP(s string) bool {
	return strings.HasPrefix(http.MethodGet, s[:3]) ||
		strings.HasPrefix(http.MethodPost, s[:4]) ||
		strings.HasPrefix(http.MethodPut, s[:3]) ||
		strings.HasPrefix(http.MethodDelete, s) ||
		strings.HasPrefix(http.MethodOptions, s) ||
		strings.HasPrefix(http.MethodPatch, s) ||
		strings.HasPrefix(http.MethodHead, s[:4]) ||
		strings.HasPrefix(http.MethodConnect, s) ||
		strings.HasPrefix(http.MethodTrace, s)
}
