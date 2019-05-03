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

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/IBM/cap/go/atom"
	"github.com/urfave/cli"
)

const (
	// DefaultAtomFeed - location used by captn to get and push atom data
	DefaultAtomFeed = atom.NwsNationalAtomFeedURL
)

func main() {
	app := cli.NewApp()
	app.Name = "captn"
	app.Usage = "Client to work with National Weather Service ATOM Feed and CAP Alerts"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "atom-feed,a",
			Usage: "url to atom feed host url",
			Value: DefaultAtomFeed,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:        "pull",
			Aliases:     []string{"p"},
			Usage:       "pull nws atom feed",
			Description: "Get the national weather service atom feed and dumps it as output",
			Action: func(c *cli.Context) error {
				_, raw, err := atom.GetFeedFrom(c.GlobalString("atom-feed"))
				if err != nil {
					return err
				}
				fmt.Printf("%s", raw)
				return nil
			},
		},
		{
			Name:      "alert",
			Aliases:   []string{"a"},
			Usage:     "get CAP alerts of a certain type",
			ArgsUsage: "TYPE",
			Description: `loads all CAP alert(s) of TYPE and dumps them as output, default TYPE is any/all.

   Examples: captn alert fire
             captn alert flood`,
			Action: func(c *cli.Context) error {
				alertType := strings.ToLower(c.Args().Get(0))
				feed, _, err := atom.GetFeedFrom(c.GlobalString("atom-feed"))
				if err != nil {
					return err
				}
				for _, entry := range feed.Entries {
					if strings.Contains(strings.ToLower(entry.Event), alertType) {
						_, raw, err := entry.Link[0].GetAlert()
						if err != nil {
							return err
						}
						fmt.Printf("%s", raw)
					}
				}
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
