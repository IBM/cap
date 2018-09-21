/*
Copyright 2018 The cap Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"context"
	"os"
	"os/signal"

	"github.com/IBM/cap/go/log"
	"github.com/IBM/cap/go/version"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const usage = `
                                .__     .__
  ____  _____   ______    ______|  |__  |__|______
_/ ___\ \__  \  \____ \  /  ___/|  |  \ |  |\____ \
\  \___  / __ \_|  |_> > \___ \ |   Y  \|  ||  |_> >
 \___  >(____  /|   __/ /____  >|___|  /|__||   __/
     \/      \/ |__|         \/      \/     |__|

a server for CAP Alerts and Atomfeeds with CAP Alert Summaries
`

// App returns a *cli.App instance.
func App() *cli.App {
	app := cli.NewApp()
	app.Name = "capship"
	app.Version = version.Version
	app.Usage = usage
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config,c",
			Usage: "path to the configuration file",
			Value: DefaultConfigPath,
		},
		cli.StringFlag{
			Name:  "root",
			Usage: "capship root directory",
		},
		cli.Int64Flag{
			Name:  "max-upload-size,m",
			Usage: "int64",
			Value: DefaultMaxUploadSize,
		},
		cli.StringFlag{
			Name:  "log-level,l",
			Usage: "logging level [trace, debug, info, warn, error, fatal, panic]",
		},
	}
	app.Commands = []cli.Command{
		configCommand,
	}
	app.Action = func(c *cli.Context) error {
		var (
			signalsCH = make(chan os.Signal, 2048)
			serverCH  = make(chan *Server, 1)
			bctx      = context.Background()
			config    = defaultConfig()
		)

		done := handleSignals(bctx, signalsCH, serverCH)
		signal.Notify(signalsCH, handledSignals...)

		if err := LoadConfig(c.GlobalString("config"), config); err != nil && !os.IsNotExist(err) {
			return err
		}
		// apply flags to the config
		if err := applyFlags(c, config); err != nil {
			return err
		}

		log.G(bctx).WithFields(logrus.Fields{
			"version":  version.Version,
			"revision": version.Revision,
		}).Info("starting capship")

		server, err := New(bctx, config)
		if err != nil {
			return err
		}
		log.G(bctx).WithFields(logrus.Fields{
			"log-level": log.PrintLevel(log.G(bctx).Logger),
		}).Info("logger")

		serverCH <- server
		go server.serve()

		log.G(bctx).Infof("capship ready:")
		log.G(bctx).Infof("   use '/cap/' to pull the cap alert feed")
		log.G(bctx).Infof("   use '/cap/{reference}' to pull a cap alert file")
		log.G(bctx).Infof("   use '/upload' to upload unique alert files using curl etc. Example:")
		log.G(bctx).Infof("      $ curl -F 'uploadFile=@KAR0-0306112239-SW.xml' http://localhost:8080/upload")
		log.G(bctx).Infof("   use '/feeds/{fileName}' to download feed files")
		log.G(bctx).Infof("   use '/alerts/{fileName}' to download alert files")

		<-done
		return nil
	}
	return app
}
