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

// Package chefsolo implements a provisioner using Chef Solo.
package chefsolo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
	"path"
	"strings"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/carton"
	"github.com/megamsys/gulp/meta"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/cmd"
)

const (
	// DefaultFormat is the default output format of Chef.
	DefaultFormat = "doc"

	// DefaultLogLevel is the set log level (default: info)
	DefaultLogLevel = "info"

	//Do not run commands with sudo (enabled by default)
	DefaultSudo       = true
	NAME              = "name"
	CHEFREPO_GIT      = "chefrepo"
	CHEFREPO_TARBALL  = "chefrepo_tarball"
	CHEFREPO_COOKBOOK = "cookbook"
)

var mainChefSoloProvisioner *chefsoloProvisioner

type Attributes struct {
	RunList   []string `json:"run_list"`
	ToscaType string   `json:"tosca_type"`
	Scm       string   `json:"scm"`
}

// Provisioner is a provisioner based on Chef Solo.
type chefsoloProvisioner struct {
	RunList    []string
	Attributes string
	Format     string
	LogLevel   string
	Cookbook   string
	RootPath   string
	Sudo       bool
}

func init() {
	mainChefSoloProvisioner = &chefsoloProvisioner{}
	provision.Register(provision.CHEFSOLO, mainChefSoloProvisioner)
}

//initialize the provisioner and setup the requirements for provisioner
func (p *chefsoloProvisioner) Initialize(m map[string]string) error {
	var outBuffer bytes.Buffer
	start := time.Now()

	p.Cookbook = m[CHEFREPO_COOKBOOK]
	logWriter := carton.NewLogWriter(&provision.Box{CartonName: m[NAME]})
	writer := io.MultiWriter(&outBuffer, &logWriter)
	defer logWriter.Close()

	cr := NewChefRepo(m, writer)
	if err := cr.Download(true); err != nil {
		return err
	}
	if err := cr.Torr(); err != nil {
		return err
	}
	elapsed := time.Since(start)

	log.Debugf("%s in (%s)\n%s",
		cmd.Colorfy(m[NAME], "cyan", "", "bold"),
		cmd.Colorfy(elapsed.String(), "green", "", "bold"),
		cmd.Colorfy(outBuffer.String(), "yellow", "", ""))
	return nil
}

func (p *chefsoloProvisioner) StartupMessage() (string, error) {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("  > chefsolo ", "white", "", "bold") + "\t" +
		cmd.Colorfy(p.String(), "cyan", "", "")))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String()), nil
}

func (p *chefsoloProvisioner) String() string {
	return "ready"
}

func (p *chefsoloProvisioner) Bootstrap(box *provision.Box, w io.Writer) error {
	fmt.Fprintf(w, "--- bootstrap box (%s)\n", box.GetFullName())
	actions := []*action.Action{
		&createMachine,
		&updateStatusInRiak,
		&updateIpsInRiak,
		&appendAuthKeys,
		&updateStatusInRiak,
		&changeStateofMachine,
	}

	pipeline := action.NewPipeline(actions...)

	args := runMachineActionsArgs{
		box:           box,
		writer:        w,
		machineStatus: provision.StatusBootstrapping,
		provisioner:   p,
	}

	if err := pipeline.Execute(args); err != nil {
		return err
	}
	fmt.Fprintf(w, "--- bootstrap box (%s) OK\n", box.GetFullName())
	return nil
}

func (p *chefsoloProvisioner) Stateup(b *provision.Box, w io.Writer) error {
	fmt.Fprintf(w, "--- stateup box (%s)\n", b.GetFullName())
	var repo string
	if b.Repo != nil {
		repo = b.Repo.Gitr()
	}

	DefaultAttributes, _ := json.Marshal(&Attributes{
		RunList:   []string{"recipe[" + p.Cookbook + "]"},
		ToscaType: b.GetShortTosca(),
		Scm:       repo,
	})

	p.Attributes = string(DefaultAttributes)
	p.Format = DefaultFormat
	p.LogLevel = DefaultLogLevel
	p.RootPath = meta.MC.Dir
	p.Sudo = DefaultSudo
	if b.Level != provision.BoxNone {
			return p.kickOffSolo(b, w)
	}
	fmt.Fprintf(w, "--- stateup box (%s) OK\n", b.GetFullName())
  return nil
}

//1. &prepareJSON in generate the json file for chefsolo
//2. &prepareConfig in generate the config file for chefsolo.
//3. &updateStatus in Riak - Creating..
func (p *chefsoloProvisioner) kickOffSolo(b *provision.Box, w io.Writer) error {
	fmt.Fprintf(w, "--- kickofff chefsolo box (%s)\n", b.GetFullName())
	actions := []*action.Action{
		&generateSoloJson,
		&generateSoloConfig,
		&cloneBox,
		&chefSoloRun,
		&updateStatusInRiak,
	}
	pipeline := action.NewPipeline(actions...)
	args := runMachineActionsArgs{
		box:           b,
		writer:        w,
		machineStatus: provision.StatusRunning,
		provisioner:   p,
	}

	if err := pipeline.Execute(args); err != nil {
		log.Errorf("error on execute chefsolo pipeline for box %s - %s", b.GetFullName(), err)
		return err
	}
	fmt.Fprintf(w, "--- kickofff chefsolo box (%s) OK\n", b.GetFullName())
	return nil
}

func (p *chefsoloProvisioner) Start(b *provision.Box, w io.Writer) error {
	fmt.Fprintf(w, "--- start box (%s)\n", b.GetFullName())
	actions := []*action.Action{
		&updateStatusInRiak,
		&startBox,
		&updateStatusInRiak,
	}
	pipeline := action.NewPipeline(actions...)
	args := runMachineActionsArgs{
		box:           b,
		writer:        w,
		machineStatus: provision.StatusStarting,
		provisioner:   p,
	}

	if err := pipeline.Execute(args); err != nil {
		log.Errorf("error on execute start pipeline for box %s - %s", b.GetFullName(), err)
		return err
	}
	fmt.Fprintf(w, "--- start box (%s) OK\n", b.GetFullName())
	return nil
}

func (p *chefsoloProvisioner) Stop(b *provision.Box, w io.Writer) error {
	fmt.Fprintf(w, "--- stop box (%s)\n", b.GetFullName())
	actions := []*action.Action{
		&updateStatusInRiak,
		&stopBox,
		&updateStatusInRiak,
	}
	pipeline := action.NewPipeline(actions...)
	args := runMachineActionsArgs{
		box:           b,
		writer:        w,
		machineStatus: provision.StatusStopping,
		provisioner:   p,
	}

	if err := pipeline.Execute(args); err != nil {
		log.Errorf("error on execute stop pipeline for box %s - %s", b.GetFullName(), err)
		return err
	}
	fmt.Fprintf(w, "--- stop box (%s) OK\n", b.GetFullName())
	return nil
}

func (p *chefsoloProvisioner) Restart(b *provision.Box, w io.Writer) error {
	fmt.Fprintf(w, "--- restart box (%s)\n", b.GetFullName())
	actions := []*action.Action{
		&updateStatusInRiak,
		&stopBox,
		&startBox,
		&updateStatusInRiak,
	}
	pipeline := action.NewPipeline(actions...)
	args := runMachineActionsArgs{
		box:           b,
		writer:        w,
		machineStatus: provision.StatusRestarting,
		provisioner:   p,
	}

	if err := pipeline.Execute(args); err != nil {
		log.Errorf("error on execute restart pipeline for box %s - %s", b.GetFullName(), err)
		return err
	}
	fmt.Fprintf(w, "--- restart box (%s) OK\n", b.GetFullName())
	return nil
}

// Command returns the command string which will invoke the provisioner on the
// prepared machine.
func (p chefsoloProvisioner) Command() []string {
	format := p.Format
	if format == "" {
		format = DefaultFormat
	}

	logLevel := p.LogLevel
	if logLevel == "" {
		logLevel = DefaultLogLevel
	}

	cmd := []string{
		"chef-solo",
		"--config", path.Join(p.RootPath, "solo.rb"),
		"--json-attributes", path.Join(p.RootPath, "solo.json"),
		"--format", format,
		"--log_level", logLevel,
	}
	return cmd
}
