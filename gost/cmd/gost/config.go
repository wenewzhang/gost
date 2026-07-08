package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/wenewzhang/core/logger"
	"github.com/wenewzhang/core/service"
	"github.com/wenewzhang/x/config"
	"github.com/wenewzhang/x/config/parsing"
	xlogger "github.com/wenewzhang/x/logger"
	"github.com/wenewzhang/x/registry"
	"gopkg.in/natefinch/lumberjack.v2"
)

func buildService(cfg *config.Config) (services []service.Service) {
	// fmt.Fprint(os.Stdout, "lite-gost2 gost config buildService\n")
	if cfg == nil {
		return
	}

	log := logger.Default()
	// str, _ := json.MarshalIndent(cfg, "", "\t")
	// fmt.Fprint(os.Stdout, "lite-gost2 gost config cfg: %s\n", string(str))

	for _, chainCfg := range cfg.Chains {
		c, err := parsing.ParseChain(chainCfg)
		if err != nil {
			log.Fatal(err)
		}
		if c != nil {
			// log.Info("lite-gost2 registry.ChainRegistry")
			if err := registry.ChainRegistry().Register(chainCfg.Name, c); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, svcCfg := range cfg.Services {
		svc, err := parsing.ParseService(svcCfg)
		if err != nil {
			log.Fatal(err)
		}
		if svc != nil {
			// log.Info("lite-gost2 registry.ServiceRegistry")

			if err := registry.ServiceRegistry().Register(svcCfg.Name, svc); err != nil {
				log.Fatal(err)
			}
		}
		services = append(services, svc)
	}

	return
}

func logFromConfig(cfg *config.LogConfig) logger.Logger {
	// fmt.Fprint(os.Stdout, "lite-gost2 gost config logFromConfig\n")
	if cfg == nil {
		cfg = &config.LogConfig{}
	}
	opts := []xlogger.LoggerOption{
		xlogger.FormatLoggerOption(logger.LogFormat(cfg.Format)),
		xlogger.LevelLoggerOption(logger.LogLevel(cfg.Level)),
	}
	// str, _ := json.MarshalIndent(cfg, "", "\t")
	// fmt.Fprint(os.Stdout, "lite-gost2 gost config logFromConfig: %s\n", string(str))

	var out io.Writer = os.Stderr
	switch cfg.Output {
	case "none", "null":
		return xlogger.Nop()
	case "stdout":
		out = os.Stdout
	case "stderr", "":
		out = os.Stderr
	default:
		if cfg.Rotation != nil {
			out = &lumberjack.Logger{
				Filename:   cfg.Output,
				MaxSize:    cfg.Rotation.MaxSize,
				MaxAge:     cfg.Rotation.MaxAge,
				MaxBackups: cfg.Rotation.MaxBackups,
				LocalTime:  cfg.Rotation.LocalTime,
				Compress:   cfg.Rotation.Compress,
			}
		} else {
			os.MkdirAll(filepath.Dir(cfg.Output), 0755)
			f, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				logger.Default().Warn(err)
			} else {
				out = f
			}
		}
	}
	opts = append(opts, xlogger.OutputLoggerOption(out))
	// fmt.Fprint(os.Stdout, "lite-gost2 gost config logFromConfig end!\n")
	return xlogger.NewLogger(opts...)
}

// func buildAPIService(cfg *config.APIConfig) (service.Service, error) {
// 	fmt.Fprint(os.Stdout, "lite-gost2 gost config buildAPIService!\n")
// 	auther := parsing.ParseAutherFromAuth(cfg.Auth)
// 	if cfg.Auther != "" {
// 		auther = registry.AutherRegistry().Get(cfg.Auther)
// 	}
// 	return api.NewService(
// 		cfg.Addr,
// 		api.PathPrefixOption(cfg.PathPrefix),
// 		api.AccessLogOption(cfg.AccessLog),
// 		api.AutherOption(auther),
// 	)
// }
