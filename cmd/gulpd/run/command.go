/*
** Copyright [2013-2016] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */
package run

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"encoding/base64"
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"launchpad.net/gnuflag"
)

type configFile struct {
	value string
}

func (v *configFile) String() string {
	return v.value
}

func (v *configFile) Set(value string) error {
	v.value = value
	return nil
}

type Start struct {
	fs   *gnuflag.FlagSet
	dry  bool
	file configFile
}

func (g *Start) Info() *cmd.Info {
	desc := `starts gulpd.
`
	return &cmd.Info{
		Name:    "start",
		Usage:   `start [--config]`,
		Desc:    desc,
		MinArgs: 0,
	}
}

func (c *Start) Run(context *cmd.Context) error {
	log.Infof("gulpd.")

	config, err := c.ParseConfig(c.file.String())
	if err != nil {
		return fmt.Errorf("parse config: %s", err)
	}
	cmd := NewCommand()
	cmd.Version = "1.5-1"

	if err := cmd.Gpd(config, cmd.Version); err != nil {
		return fmt.Errorf("run: %s", err)
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	fmt.Println("Listening for signals")

	// Block until one of the signals above is received
	select {
	case <-signalCh:
		fmt.Println("Signal received, initializing clean shutdown...")
		go func() {
			cmd.Close()
		}()
	}

	// Block again until another signal is received, a shutdown timeout elapses,
	// or the Command is gracefully closed
	fmt.Println("Waiting for clean shutdown...")
	select {
	case <-signalCh:
		fmt.Println("second signal received, initializing hard shutdown")
	case <-time.After(time.Second * 30):
		fmt.Println("time limit reached, initializing hard shutdown")
	case <-cmd.Closed:
		fmt.Println("server shutdown completed")
	}
	// goodbye.
	return nil
}

func (c *Start) Flags() *gnuflag.FlagSet {
	if c.fs == nil {
		c.fs = gnuflag.NewFlagSet("gulpd", gnuflag.ExitOnError)
		c.fs.Var(&c.file, "config", "Path to configuration file (default to /conf/gulpd.conf)")
		c.fs.Var(&c.file, "c", "Path to configuration file (default to /conf/gulpd.conf)")

	}
	return c.fs
}

// ParseConfig parses the config at path.
func (c *Start) ParseConfig(path string) (*Config, error) {
	config := NewConfig()
	if path == "" {
		path = config.Meta.Dir + "/gulpd.conf"
	}
	log.Warnf("Using configuration at: %s", path)
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	user, _ := base64.StdEncoding.DecodeString(config.Meta.ScyllaUsername)
	username := string(user)
	pass, _ := base64.StdEncoding.DecodeString(config.Meta.ScyllaPassword)
	password := string(pass)
	config.Meta.ScyllaUsername = username
	config.Meta.ScyllaPassword = password
	log.Debug(config)
	return config, nil
}
