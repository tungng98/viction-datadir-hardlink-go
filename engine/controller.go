package engine

import (
	"os"
	"viction-datadir-clone-go/config"

	"github.com/rs/zerolog"
)

type Controller struct {
	Logger  zerolog.Logger
	Root    *config.RootConfig
	logFile *os.File
}

func NewController(useFS bool) *Controller {
	cfg := &config.RootConfig{}
	logger, logFile, err2 := config.InitZerolog(cfg.ConfigDir, useFS)
	if err2 != nil {
		logger.Err(err2).Msg("error initializing log file")
	}
	return &Controller{
		Logger:  logger,
		Root:    cfg,
		logFile: logFile,
	}
}

func (c *Controller) Close() {
	if c.logFile != nil {
		c.logFile.Close()
		c.logFile = nil
	}
}

func (c *Controller) CommandLogger(module, command string) zerolog.Logger {
	return c.Logger.With().Str("module", module).Str("command", command).Logger()
}

func (c *Controller) ModuleLogger(module string) zerolog.Logger {
	return c.Logger.With().Str("module", module).Logger()
}
