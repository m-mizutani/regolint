package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/zlog"
	"github.com/urfave/cli/v2"
)

var logger = zlog.New()

const outputStdout = "-"

type config struct {
	PolicyFile string
	OutputFile string

	LogLevel string
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
				Usage:       "lint policy file/dir. If no policy file, output only parsed rego files",
				EnvVars:     []string{"REGOLINT_POLICY"},
				Destination: &cfg.PolicyFile,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "specify output file. `-` means stdout",
				EnvVars:     []string{"REGOLINT_OUTPUT"},
				Destination: &cfg.OutputFile,
				Value:       outputStdout,
			},

			&cli.StringFlag{
				Name:        "log-level",
				Aliases:     []string{"l"},
				Usage:       "Log level [trace|debug|info|warn|error]",
				Destination: &cfg.LogLevel,
				Value:       "info",
			},
		},
		Before: func(c *cli.Context) error {
			created, err := zlog.NewWithError(zlog.WithLogLevel(cfg.LogLevel))
			if err != nil {
				return err
			}
			logger = created
			logger.With("config", cfg).Debug("starting regolint...")
			return nil
		},
		Action: func(c *cli.Context) error {
			targets, err := loadDirs(c.Args().Slice()...)
			if err != nil {
				return err
			}

			var output io.Writer = os.Stdout
			if cfg.OutputFile != outputStdout {
				f, err := os.Create(cfg.OutputFile)
				if err != nil {
					return goerr.Wrap(err)
				}
				output = f
			}

			if cfg.PolicyFile != "" {
				if err := evalWithFile(cfg.PolicyFile, targets, output); err != nil {
					return err
				}
			} else {
				raw, err := json.MarshalIndent(input{Files: targets}, "", "  ")
				if err != nil {
					return goerr.Wrap(err)
				}
				fmt.Fprint(output, string(raw))
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
