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

package chefsolo_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mlafeldt/chef-runner/provisioner"
	. "github.com/mlafeldt/chef-runner/provisioner/chefsolo"
	"github.com/mlafeldt/chef-runner/util"
	"github.com/stretchr/testify/assert"
)


func TestProvisionerInterface(t *testing.T) {
	assert.Implements(t, (*provisioner.Provisioner)(nil), new(Provisioner))
}

func TestPrepareFiles(t *testing.T) {
	util.InTestDir(func() {
		os.MkdirAll(".chef-runner/sandbox", 0755)

		p := Provisioner{
			SandboxPath: ".chef-runner/sandbox",
			RootPath:    "/tmp/chef-runner",
		}
		assert.NoError(t, p.PrepareFiles())

		json, err := ioutil.ReadFile(".chef-runner/sandbox/dna.json")
		assert.NoError(t, err)
		assert.Equal(t, "{}\n", string(json))

		expect := "cookbook_path \"/tmp/chef-runner/cookbooks\"\n"
		expect += "ssl_verify_mode :verify_peer\n"

		config, err := ioutil.ReadFile(".chef-runner/sandbox/solo.rb")
		assert.NoError(t, err)
		assert.Equal(t, expect, string(config))
	})
}

func TestPrepareFiles_CustomJSON(t *testing.T) {
	util.InTestDir(func() {
		os.MkdirAll(".chef-runner/sandbox", 0755)

		p := Provisioner{
			Attributes:  `{"foo": "bar"}`,
			SandboxPath: ".chef-runner/sandbox",
			RootPath:    "/tmp/chef-runner",
		}
		assert.NoError(t, p.PrepareFiles())

		json, err := ioutil.ReadFile(".chef-runner/sandbox/dna.json")
		assert.NoError(t, err)
		assert.Equal(t, `{"foo": "bar"}`, string(json))

		expect := "cookbook_path \"/tmp/chef-runner/cookbooks\"\n"
		expect += "ssl_verify_mode :verify_peer\n"

		config, err := ioutil.ReadFile(".chef-runner/sandbox/solo.rb")
		assert.NoError(t, err)
		assert.Equal(t, expect, string(config))
	})
}

var commandTests = []struct {
	provisioner Provisioner
	cmd         []string
}{
	{
		Provisioner{
			SandboxPath: ".chef-runner/sandbox",
			RootPath:    "/tmp/chef-runner",
		},
		[]string{
			"chef-solo", "--config", "/tmp/chef-runner/solo.rb",
			"--json-attributes", "/tmp/chef-runner/dna.json",
			"--format", "doc", "--log_level", "info",
		},
	},
	{
		Provisioner{
			RunList:     []string{"cats::foo"},
			SandboxPath: ".chef-runner/sandbox",
			RootPath:    "/tmp/chef-runner",
		},
		[]string{
			"chef-solo", "--config", "/tmp/chef-runner/solo.rb",
			"--json-attributes", "/tmp/chef-runner/dna.json",
			"--format", "doc", "--log_level", "info",
			"--override-runlist", "cats::foo",
		},
	},
	{
		Provisioner{
			RunList:     []string{"cats::foo"},
			Format:      "null",
			SandboxPath: ".chef-runner/sandbox",
			RootPath:    "/tmp/chef-runner",
		},
		[]string{
			"chef-solo", "--config", "/tmp/chef-runner/solo.rb",
			"--json-attributes", "/tmp/chef-runner/dna.json",
			"--format", "null", "--log_level", "info",
			"--override-runlist", "cats::foo",
		},
	},
	{
		Provisioner{
			RunList:     []string{"cats::foo"},
			LogLevel:    "error",
			SandboxPath: ".chef-runner/sandbox",
			RootPath:    "/tmp/chef-runner",
		},
		[]string{
			"chef-solo", "--config", "/tmp/chef-runner/solo.rb",
			"--json-attributes", "/tmp/chef-runner/dna.json",
			"--format", "doc", "--log_level", "error",
			"--override-runlist", "cats::foo",
		},
	},
	{
		Provisioner{
			RunList:     []string{"cats::foo", "dogs::bar"},
			Format:      "min",
			LogLevel:    "warn",
			SandboxPath: ".chef-runner/sandbox",
			RootPath:    "/tmp/chef-runner",
		},
		[]string{
			"chef-solo", "--config", "/tmp/chef-runner/solo.rb",
			"--json-attributes", "/tmp/chef-runner/dna.json",
			"--format", "min", "--log_level", "warn",
			"--override-runlist", "cats::foo,dogs::bar",
		},
	},
	{
		Provisioner{
			RunList:     []string{"cats::foo"},
			Sudo:        true,
			SandboxPath: ".chef-runner/sandbox",
			RootPath:    "/tmp/chef-runner",
		},
		[]string{
			"sudo", "chef-solo", "--config", "/tmp/chef-runner/solo.rb",
			"--json-attributes", "/tmp/chef-runner/dna.json",
			"--format", "doc", "--log_level", "info",
			"--override-runlist", "cats::foo",
		},
	},
}


