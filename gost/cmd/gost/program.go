package main

import (
	"net/http"
	"os"

	"github.com/judwhite/go-svc"
	"github.com/wenewzhang/core/logger"
	"github.com/wenewzhang/x/config"
	"github.com/wenewzhang/x/config/parsing"
	"github.com/wenewzhang/x/registry"
)

type program struct {
}

func (p *program) Init(env svc.Environment) error {
	// fmt.Fprint(os.Stdout, "lite-gost2 gost program Init\n")
	cfg := &config.Config{}
	if cfgFile != "" {
		if err := cfg.ReadFile(cfgFile); err != nil {
			return err
		}
	}

	cmdCfg, err := buildConfigFromCmd(services, nodes)
	if err != nil {
		return err
	}
	cfg = p.mergeConfig(cfg, cmdCfg)

	if len(cfg.Services) == 0 {
		if err := cfg.Load(); err != nil {
			return err
		}
	}

	if v := os.Getenv("GOST_LOGGER_LEVEL"); v != "" {
		cfg.Log = &config.LogConfig{
			Level: v,
		}
	}
	if v := os.Getenv("GOST_PROFILING"); v != "" {
		cfg.Profiling = &config.ProfilingConfig{
			Addr: v,
		}
	}

	if debug {
		if cfg.Log == nil {
			cfg.Log = &config.LogConfig{}
		}
		cfg.Log.Level = string(logger.DebugLevel)
	}

	logger.SetDefault(logFromConfig(cfg.Log))

	if outputFormat != "" {
		if err := cfg.Write(os.Stdout, outputFormat); err != nil {
			return err
		}
		os.Exit(0)
	}

	parsing.BuildDefaultTLSConfig(cfg.TLS)

	config.Set(cfg)

	return nil
}

func (p *program) Start() error {
	// fmt.Fprint(os.Stdout, "lite-gost2 gost program Start\n")

	log := logger.Default()
	cfg := config.Global()

	if cfg.Profiling != nil {
		go func() {
			addr := cfg.Profiling.Addr
			if addr == "" {
				addr = ":6060"
			}
			log.Info("profiling server on ", addr)
			log.Fatal(http.ListenAndServe(addr, nil))
		}()
	}

	for _, svc := range buildService(cfg) {
		svc := svc
		go func() {
			svc.Serve()
		}()
	}

	return nil
}

func (p *program) Stop() error {
	// fmt.Fprint(os.Stdout, "lite-gost2 gost program Stop\n")
	for name, srv := range registry.ServiceRegistry().GetAll() {
		srv.Close()
		logger.Default().Debugf("service %s shutdown", name)
	}
	return nil
}

func (p *program) mergeConfig(cfg1, cfg2 *config.Config) *config.Config {
	if cfg1 == nil {
		return cfg2
	}
	if cfg2 == nil {
		return cfg1
	}

	cfg := &config.Config{
		Services:  append(cfg1.Services, cfg2.Services...),
		Chains:    append(cfg1.Chains, cfg2.Chains...),
		TLS:       cfg1.TLS,
		Log:       cfg1.Log,
		Profiling: cfg1.Profiling,
	}
	if cfg2.TLS != nil {
		cfg.TLS = cfg2.TLS
	}
	if cfg2.Log != nil {
		cfg.Log = cfg2.Log
	}
	if cfg2.Profiling != nil {
		cfg.Profiling = cfg2.Profiling
	}

	return cfg
}
