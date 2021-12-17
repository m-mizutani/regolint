package main

import (
	"os"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/zlog"
	"github.com/urfave/cli/v2"
)

var logger *zlog.Logger

type config struct {
	policyFile string
	opaURL     string

	logLevel string
}

func (x *config) isValid() error {
	if x.policyFile == "" && x.opaURL == "" {
		return goerr.Wrap(errInvalidConfig, "either one of --policy and --url is required")
	}

	if x.policyFile != "" && x.opaURL != "" {
		return goerr.Wrap(errInvalidConfig, "only either one of --policy and --url is allowed")
	}

	return nil
}

func Run(args []string) error {
	var cfg config
	app := &cli.App{
		Name:      "regolint",
		Usage:     "Linting Rego file with policy written by Rego",
		ArgsUsage: "<rego dir> [<rego dir> [...]]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "policy",
				Aliases:     []string{"p"},
				Usage:       "lint policy file",
				EnvVars:     []string{"REGOLINT_POLICY"},
				Destination: &cfg.policyFile,
				Required:    true,
			},

			&cli.StringFlag{
				Name:        "log-level",
				Aliases:     []string{"l"},
				Usage:       "Log level [trace|debug|info|warn|error]",
				Destination: &cfg.logLevel,
				Value:       "info",
			},
		},
		Before: func(c *cli.Context) error {
			logger = zlog.New(zlog.WithLogLevel(cfg.logLevel))
			return nil
		},
		Action: func(c *cli.Context) error {
			if err := cfg.isValid(); err != nil {
				return err
			}

			targets, err := loadDirs(c.Args().Slice()...)
			if err != nil {
				return err
			}

			if err := evalWithFile(cfg.policyFile, targets); err != nil {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Error(err.Error())
		logger.Err(err).Debug("Error detail")
		return err
	}
	return nil
}
