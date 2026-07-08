package service

import (
	"context"
	"io"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/wenewzhang/core/handler"
	"github.com/wenewzhang/core/listener"
	"github.com/wenewzhang/core/logger"
	"github.com/wenewzhang/core/recorder"
	"github.com/wenewzhang/core/service"
	sx "github.com/wenewzhang/x/internal/util/selector"
)

type options struct {
	recorders []recorder.RecorderObject
	preUp     []string
	postUp    []string
	preDown   []string
	postDown  []string
	logger    logger.Logger
}

type Option func(opts *options)

func RecordersOption(recorders ...recorder.RecorderObject) Option {
	return func(opts *options) {
		opts.recorders = recorders
	}
}

func PreUpOption(cmds []string) Option {
	return func(opts *options) {
		opts.preUp = cmds
	}
}

func PreDownOption(cmds []string) Option {
	return func(opts *options) {
		opts.preDown = cmds
	}
}

func PostUpOption(cmds []string) Option {
	return func(opts *options) {
		opts.postUp = cmds
	}
}

func PostDownOption(cmds []string) Option {
	return func(opts *options) {
		opts.postDown = cmds
	}
}

func LoggerOption(logger logger.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

type defaultService struct {
	name     string
	listener listener.Listener
	handler  handler.Handler
	options  options
}

func NewService(name string, ln listener.Listener, h handler.Handler, opts ...Option) service.Service {
	var options options
	for _, opt := range opts {
		opt(&options)
	}
	s := &defaultService{
		name:     name,
		listener: ln,
		handler:  h,
		options:  options,
	}

	s.execCmds("pre-up", s.options.preUp)

	return s
}

func (s *defaultService) Addr() net.Addr {
	return s.listener.Addr()
}

func (s *defaultService) Close() error {
	s.execCmds("pre-down", s.options.preDown)
	defer s.execCmds("post-down", s.options.postDown)

	if closer, ok := s.handler.(io.Closer); ok {
		closer.Close()
	}
	return s.listener.Close()
}

func (s *defaultService) Serve() error {
	// fmt.Fprintf(os.Stdout, "lite-gost x service Serve!\n")
	s.execCmds("post-up", s.options.postUp)

	var tempDelay time.Duration
	for {
		conn, e := s.listener.Accept()
		if e != nil {
			// TODO: remove Temporary checking
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 1 * time.Second
				} else {
					tempDelay *= 2
				}
				if max := 5 * time.Second; tempDelay > max {
					tempDelay = max
				}
				s.options.logger.Warnf("accept: %v, retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			s.options.logger.Errorf("accept: %v", e)
			return e
		}
		tempDelay = 0

		host := conn.RemoteAddr().String()
		if h, _, _ := net.SplitHostPort(host); h != "" {
			host = h
		}
		for _, rec := range s.options.recorders {
			if rec.Record == recorder.RecorderServiceClientAddress {
				if err := rec.Recorder.Record(context.Background(), []byte(host)); err != nil {
					s.options.logger.Errorf("record %s: %v", rec.Record, err)
				}
				break
			}
		}

		go func() {
			// start := time.Now()
			ctx := sx.ContextWithHash(context.Background(), &sx.Hash{Source: host})
			ctx = ContextWithSid(ctx, xid.New().String())

			if err := s.handler.Handle(ctx, conn); err != nil {
				s.options.logger.Error(err)
			}
		}()
	}
}

func (s *defaultService) execCmds(phase string, cmds []string) {
	for _, cmd := range cmds {
		cmd := strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}
		s.options.logger.Info(cmd)

		if err := exec.Command("/bin/sh", "-c", cmd).Run(); err != nil {
			s.options.logger.Warnf("[%s] %s: %v", phase, cmd, err)
		}
	}
}

type sidKey struct{}

var (
	ssid sidKey
)

func ContextWithSid(ctx context.Context, sid string) context.Context {
	return context.WithValue(ctx, ssid, sid)
}

func SidFromContext(ctx context.Context) string {
	v, _ := ctx.Value(ssid).(string)
	return v
}
