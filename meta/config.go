/*
** Copyright [2013-2015] [Megam Systems]
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

package meta

import (
	"os"
	"time"
	//"fmt"
	"path/filepath"

	"github.com/megamsys/gulp/toml"
	"github.com/megamsys/libgo/cmd"
)

const (
	// DefaultHostname is the default hostname if one is not provided.
	DefaultHostname = "localhost"

	// DefaultBindAddress is the default address to bind to.
	DefaultBindAddress = ":9999"

	// DefaultRiak is the default riak if one is not provided.
	DefaultRiak = "localhost:8087"

	// DefaultApi is the default megam gateway if one is not provided.
	DefaultApi = "https://api.megam.io/v2"
	
	// DefaultProvider is the default provisioner used by our engine.
	DefaultProvider = "chefsolo"

	// DefaultAMQP is the default rabbitmq if one is not provided.
	DefaultAMQP = "amqp://guest:guest@localhost:5672/"

	// DefaultHeartbeatTimeout is the default heartbeat timeout for the store.
	DefaultHeartbeatTimeout = 1000 * time.Millisecond

	// DefaultElectionTimeout is the default election timeout for the store.
	DefaultElectionTimeout = 1000 * time.Millisecond

	// DefaultLeaderLeaseTimeout is the default leader lease for the store.
	DefaultLeaderLeaseTimeout = 500 * time.Millisecond

	//DefaultDockerPath is the detault docker path
	DefaultDockerPath = "/var/lib/docker/containers/"
)

var MC *Config

// Config represents the meta configuration.
type Config struct {
	Home               string        `toml:"home"`
	Dir                string        `toml:"dir"`
	Hostname           string        `toml:"hostname"`
	BindAddress        string        `toml:"bind_address"`
	Riak               string        `toml:"riak"`
	Api                string        `toml:"api"`
	AMQP               string        `toml:"amqp"`
	Peers              []string      `toml:"-"`
	Provider 		   string 	     `toml:"provider"`
	ElectionTimeout    toml.Duration `toml:"election-timeout"`
	HeartbeatTimeout   toml.Duration `toml:"heartbeat-timeout"`
	LeaderLeaseTimeout toml.Duration `toml:"leader-lease-timeout"`
	DockerPath         string        `toml:"docker_path"`
}

func (c Config) String() string {
	table := cmd.NewTable()
	table.AddRow(cmd.Row{cmd.Colorfy("Config:", "white", "", "bold"), cmd.Colorfy("Meta", "green", "", "")})
	table.AddRow(cmd.Row{"Home", c.Home})
	table.AddRow(cmd.Row{"Dir", c.Dir})
	table.AddRow(cmd.Row{"Riak", c.Riak})
	table.AddRow(cmd.Row{"API", c.Api})
	table.AddRow(cmd.Row{"AMQP", c.AMQP})
	table.AddRow(cmd.Row{"Hostname", c.Hostname})
	table.AddRow(cmd.Row{"", ""})
	return table.String()
}

func NewConfig() *Config {
	var homeDir string
	// By default, store logs, meta and load conf files in current users home directory
	if os.Getenv("MEGAM_HOME") != "" {
		homeDir = os.Getenv("MEGAM_HOME")
		//} else if u, err := user.Current(); err == nil {
		//	homeDir = u.HomeDir
		//} else {
		//	return nil, fmt.Errorf("failed to determine home directory")
	}

	defaultDir := filepath.Join(homeDir, "gulp/meta")

	// Config represents the configuration format for the gulpd.
	return &Config{
		Home:               homeDir, //Need to remove
		Dir:                defaultDir,
		Hostname:           DefaultHostname,
		BindAddress:        DefaultBindAddress,
		Riak:               DefaultRiak,
		Api:                DefaultApi,
		AMQP:               DefaultAMQP,
		Provider: 			DefaultProvider,
		DockerPath:         DefaultDockerPath,
		ElectionTimeout:    toml.Duration(DefaultElectionTimeout),
		HeartbeatTimeout:   toml.Duration(DefaultHeartbeatTimeout),
		LeaderLeaseTimeout: toml.Duration(DefaultLeaderLeaseTimeout),
	}
}

func (c *Config) MC() {
	MC = c
}
