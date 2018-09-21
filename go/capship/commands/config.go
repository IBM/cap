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
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/IBM/cap/go/caperrors"
	"github.com/IBM/cap/go/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	// DefaultRootDir - location used by capship to store data
	DefaultRootDir = "/var/lib/capship/"
	// DefaultConfigPath - default path to the config file
	DefaultConfigPath = "/etc/capship/config.toml"
	// DefaultLogLevel - default loging level
	DefaultLogLevel = "info"
	// DefaultMaxUploadSize - for all uploads
	DefaultMaxUploadSize = 2 * 1048576 // 2mb
)

// Config provides capship configuration data
type Config struct {
	// Root is the path to a directory where capship will store persistent data
	Root     string `toml:"root"`
	LogLevel string `toml:"log_level"`
	// MaxUploadSize is the maximum file size permitted for uploads
	MaxUploadSize int64 `toml:"max_upload_size"`

	md toml.MetaData
}

// LoadConfig loads the capship config from the provided path
func LoadConfig(path string, v *Config) error {
	if v == nil {
		return errors.Wrapf(caperrors.ErrInvalidArgument, "argument v must not be nil")
	}
	md, err := toml.DecodeFile(path, v)
	if err != nil {
		return err
	}
	v.md = md
	return nil
}

// WriteTo marshals the config to the provided writer
func (c *Config) WriteTo(w io.Writer) (int64, error) {
	return 0, toml.NewEncoder(w).Encode(c)
}

var configCommand = cli.Command{
	Name:  "config",
	Usage: "information on the capship config",
	Subcommands: []cli.Command{
		{
			Name:  "default",
			Usage: "see the output of the default config",
			Action: func(context *cli.Context) error {
				config := defaultConfig()
				_, err := config.WriteTo(os.Stdout)
				return err
			},
		},
	},
}

func defaultConfig() *Config {
	return &Config{
		Root:          DefaultRootDir,
		LogLevel:      DefaultLogLevel,
		MaxUploadSize: DefaultMaxUploadSize,
	}
}

func applyFlags(c *cli.Context, config *Config) error {
	// flags override config values
	if err := setLevel(c, config); err != nil {
		return err
	}
	root := c.GlobalString("root")
	if root != "" {
		config.Root = root
	}
	uls := c.GlobalInt64("max_upload_size")
	if uls != 0 {
		config.MaxUploadSize = uls
	}
	return nil
}

func setLevel(context *cli.Context, config *Config) error {
	l := context.GlobalString("log-level")
	if l == "" {
		l = config.LogLevel
	}
	if l != "" {
		lvl, err := log.ParseLevel(l)
		if err != nil {
			return err
		}
		logrus.SetLevel(lvl)
	}
	return nil
}
