// Copyright 2020 arugal, zhangwei24@apache.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github/arugal/frp-notify/pkg/cli/interceptor"
	"github/arugal/frp-notify/pkg/cmd"
	"github/arugal/frp-notify/pkg/config"
	"github/arugal/frp-notify/pkg/controller"
	"github/arugal/frp-notify/pkg/controller/handler"
	"github/arugal/frp-notify/pkg/ip"
	"github/arugal/frp-notify/pkg/logger"
	_ "github/arugal/frp-notify/pkg/notify/dingtalk"
	_ "github/arugal/frp-notify/pkg/notify/gotify"
	_ "github/arugal/frp-notify/pkg/notify/log"
	"github/arugal/frp-notify/pkg/version"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
)

var (
	log      *logrus.Logger
	cmdStart = cli.Command{
		Name:  "start",
		Usage: "start frp notify",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "config, c",
				Usage:    "load notify plugin configuration from `FILE`",
				Required: false,
				EnvVar:   "FRP_NOTIFY_PLUGIN_CONF",
				Value:    "frp-notify.json",
			},
			cli.StringFlag{
				Name:     "bind-address, b",
				Usage:    "manager server listen `ADDRESS`",
				Required: false,
				EnvVar:   "FRP_NOTIFY_MANAGER_ADDRESS",
				Value:    ":80",
			},
			cli.Int64Flag{
				Name:     "window-interval",
				Usage:    "user conn cache time window interval (unit: MINUTE)",
				Required: false,
				EnvVar:   "FRP_NOTIFY_WINDOW_INTERVAL",
				Value:    60,
			},
			cli.BoolTFlag{
				Name:     "ip-query",
				Usage:    "enable the query for IP address ownership",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			configPath := ctx.String("config")
			bindAddress := ctx.String("bind-address")
			windowInterval := ctx.Int64("window-interval")
			enable := ctx.Bool("ip-query")

			// Create the stop channel for all of the servers.
			stop := make(chan struct{})

			mux := http.NewServeMux()

			ms := controller.NewManagerController(controller.WithHandlerChains(
				handler.NewWhitelistHandler(),
				handler.NewBlacklistHandler(),
				handler.NewNotifyHandler(handler.WithAddressService(enable, func() ip.AddressService {
					return ip.NewDefaultAddressService()
				}))))

			ms.Register(mux)

			// config
			configController := config.NewConfigController(bindAddress, windowInterval, configPath)
			go configController.Start(stop)

			err := ms.Start(stop)
			if err != nil {
				return err
			}

			serve := http.Server{
				Handler: mux,
				Addr:    bindAddress,
			}

			go func() {
				err := serve.ListenAndServe()
				if err != nil {
					panic(err)
				}
			}()

			cmd.WaitSignal(stop)
			_ = serve.Close()
			return nil
		},
	}
)

func init() {
	log = logger.Log
}

func main() {
	app := cli.NewApp()
	app.Version = version.Version
	app.Name = "frp-notify"
	app.Usage = "https://github.com/arugal/frp-notify"
	app.Compiled = time.Now()
	app.Copyright = "(c) " + strconv.Itoa(time.Now().Year()) + " arugal"
	app.Description = "frp server manager plugin implement, focus on notify."

	flags := []cli.Flag{
		altsrc.NewStringFlag(cli.StringFlag{
			Name:     "log-level",
			Required: false,
			Usage:    "set log level, support: panic，fatal, error, warn, info, debug, trace",
			EnvVar:   "FRP_NOTIFY_LOG_LEVEL",
			Value:    "info",
		}),
	}

	app.Flags = flags

	app.Commands = []cli.Command{
		cmdStart,
	}

	app.Before = interceptor.BeforeChain([]cli.BeforeFunc{
		setUpCommandLineContext,
	})

	app.Action = func(c *cli.Context) error {
		err := cli.ShowAppHelp(c)
		if err != nil {
			return err
		}
		c.App.Setup()
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func setUpCommandLineContext(c *cli.Context) error {
	level := c.GlobalString("log-level")
	logger.SetLogLevel(level)
	return nil
}
