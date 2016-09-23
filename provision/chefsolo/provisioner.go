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
	"path"
	"strings"
	"text/tabwriter"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/carton"
	lb "github.com/megamsys/gulp/logbox"
	"github.com/megamsys/gulp/meta"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/cmd"
	constants "github.com/megamsys/libgo/utils"
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
	RunList    []string `json:"run_list"`
	ToscaType  string   `json:"tosca_type"`
	RepoURL    string   `json:"scm"`
	RepoSource string   `json:"provider"`
	Version    string   `json:"version"`
}

// Repos for Bitnami
type ReposBitnami struct {
	RunList         []string `json:"run_list"`
	ToscaType       string   `json:"tosca_type"`
	BitnamiURL      string   `json:"bitnami_url"`
	BitnamiUserName string   `json:"bitnami_username"`
	BitnamiPassword string   `json:"bitnami_password"`
	RepoSource      string   `json:"provider"`
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
		err = provision.EventNotify(constants.StatusCookbookFailure)
		return err
	}
	if err := cr.Torr(); err != nil {
		err = provision.EventNotify(constants.StatusCookbookFailure)
		return err
	}
	elapsed := time.Since(start)

	log.Debugf("%s in (%s)\n%s",
		cmd.Colorfy(m[NAME], "cyan", "", "bold"),
		cmd.Colorfy(elapsed.String(), "green", "", "bold"),
		cmd.Colorfy(outBuffer.String(), "yellow", "", ""))
	_ = provision.EventNotify(constants.StatusCookbookDownloaded)
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
	fmt.Fprintf(w, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("--- bootstrap box (%s)\n", box.GetFullName())))
	actions := []*action.Action{
		&updateStatusInScylla,
		&createMachine,
		&updateStatusInScylla,
		&updateIpsInSyclla,
		&updateStatusInScylla,
		&appendAuthKeys,
		&updateStatusInScylla,
		&changeStateofMachine,
		&MileStoneUpdate,
		&updateStatusInScylla,
	}

	pipeline := action.NewPipeline(actions...)

	args := runMachineActionsArgs{
		box:           box,
		writer:        w,
		machineStatus: constants.StatusBootstrapping,
		machineState:  constants.StateBootstrapped,
		provisioner:   p,
	}

	if err := pipeline.Execute(args); err != nil {
		return err
	}
   switch box.GetShortTosca() {
	 case "bitnami":
		 p.StateupBitnami(box,w)
	 default:
		 p.Stateup(box, w)
   }

	fmt.Fprintf(w, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("--- bootstrap box (%s) OK\n", box.GetFullName())))
	return nil
}

func (p *chefsoloProvisioner) Stateup(b *provision.Box, w io.Writer) error {
	fmt.Fprintf(w, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("\n--- stateup box (%s)\n", b.GetFullName())))
	var repo, src string
	if b.Repo != nil {
		repo = b.Repo.Gitr()
		src = b.Repo.RepoProvider()
	}

	DefaultAttributes, _ := json.Marshal(&Attributes{
		RunList:    []string{"recipe[" + p.Cookbook + "]"},
		ToscaType:  b.GetShortTosca(),
		RepoURL:    repo,
		RepoSource: src,
		Version:    b.ImageVersion,
	})

	p.Attributes = string(DefaultAttributes)
	p.Format = DefaultFormat
	p.LogLevel = DefaultLogLevel
	p.RootPath = meta.MC.Dir
	p.Sudo = DefaultSudo
	return p.kickOffSolo(b, w)
}

func (p *chefsoloProvisioner) StateupBitnami(b *provision.Box, w io.Writer) error {
	fmt.Fprintf(w, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("\n--- stateup box (%s)\n", b.GetFullName())))
	var repo, src, username,pswd string
	if b.Repo != nil {
		repo = b.Repo.Gitr()
		src = b.Repo.RepoProvider()
	}

	if len(b.Inputs) > 0  {
		username = b.Inputs[provision.BITUSERNAME]
		pswd = b.Inputs[provision.BITPASSWORD]
	}

	DefaultAttributes, _ := json.Marshal(&ReposBitnami{
		RunList:    []string{"recipe[" + p.Cookbook + "]"},
		ToscaType:  b.GetShortTosca(),
		BitnamiURL:    repo,
		BitnamiUserName: username,
		BitnamiPassword: pswd,
		RepoSource: src,
	})

	p.Attributes = string(DefaultAttributes)
	p.Format = DefaultFormat
	p.LogLevel = DefaultLogLevel
	p.RootPath = meta.MC.Dir
	p.Sudo = DefaultSudo
	return p.kickOffSolo(b, w)
}

//1. &prepareJSON in generate the json file for chefsolo
//2. &prepareConfig in generate the config file for chefsolo.
//3. &updateStatus in Riak - Creating..
func (p *chefsoloProvisioner) kickOffSolo(b *provision.Box, w io.Writer) error {
	fmt.Fprintf(w, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("\n--- kickofff chefsolo box (%s)\n", b.GetFullName())))
	soloAction := make([]*action.Action, 0, 4)
	soloAction = append(soloAction, &updateStatusInScylla, &generateSoloJson, &generateSoloConfig, &updateStatusInScylla, &cloneBox, &updateStatusInScylla)
	if b.Level != provision.BoxNone {
		soloAction = append(soloAction, &setChefsoloStatus, &updateStatusInScylla, &chefSoloRun, &updateStatusInScylla)
	}
	soloAction = append(soloAction, &setFinalState, &MileStoneUpdate, &updateStatusInScylla)
	actions := soloAction
	pipeline := action.NewPipeline(actions...)
	args := runMachineActionsArgs{
		box:           b,
		writer:        w,
		machineStatus: constants.StatusChefConfigSetupping,
		machineState:  constants.StateRunning,
		provisioner:   p,
	}

	if err := pipeline.Execute(args); err != nil {
		log.Errorf("error on execute chefsolo pipeline for box %s - %s", b.GetFullName(), err)
		return err
	}
	fmt.Fprintf(w, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("--- kickofff chefsolo box (%s) OK\n", b.GetFullName())))
	return nil
}

func (p *chefsoloProvisioner) Start(b *provision.Box, w io.Writer) error {
	fmt.Fprintf(w, lb.W(lb.VM_STARTING, lb.INFO, fmt.Sprintf("\n--- start box (%s)\n", b.GetFullName())))
	actions := []*action.Action{
		&updateStatusInScylla,
		&startBox,
		&MileStoneUpdate,
		&updateStatusInScylla,
	}
	pipeline := action.NewPipeline(actions...)
	args := runMachineActionsArgs{
		box:           b,
		writer:        w,
		machineStatus: constants.StatusStarting,
		machineState:  constants.StateRunning,
		provisioner:   p,
	}

	if err := pipeline.Execute(args); err != nil {
		log.Errorf("error on execute start pipeline for box %s - %s", b.GetFullName(), err)
		return err
	}
	fmt.Fprintf(w, lb.W(lb.VM_STARTING, lb.INFO, fmt.Sprintf("--- start box (%s) OK\n", b.GetFullName())))
	return nil
}

func (p *chefsoloProvisioner) Stop(b *provision.Box, w io.Writer) error {
	fmt.Fprintf(w, lb.W(lb.VM_STOPPING, lb.INFO, fmt.Sprintf("\n--- stop box (%s)\n", b.GetFullName())))
	actions := []*action.Action{
		&updateStatusInScylla,
		&stopBox,
		&updateStatusInScylla,
	}
	pipeline := action.NewPipeline(actions...)
	args := runMachineActionsArgs{
		box:           b,
		writer:        w,
		machineStatus: constants.StatusStopping,
		provisioner:   p,
	}

	if err := pipeline.Execute(args); err != nil {
		log.Errorf("error on execute stop pipeline for box %s - %s", b.GetFullName(), err)
		return err
	}
	fmt.Fprintf(w, lb.W(lb.VM_STOPPING, lb.INFO, fmt.Sprintf("--- stop box (%s) OK\n", b.GetFullName())))
	return nil
}

func (p *chefsoloProvisioner) Restart(b *provision.Box, w io.Writer) error {
	fmt.Fprintf(w, lb.W(lb.VM_RESTARTING, lb.INFO, fmt.Sprintf("\n--- restart box (%s)\n", b.GetFullName())))
	actions := []*action.Action{
		&updateStatusInScylla,
		&stopBox,
		&MileStoneUpdate,
		&startBox,
		&MileStoneUpdate,
		&updateStatusInScylla,
	}
	pipeline := action.NewPipeline(actions...)
	args := runMachineActionsArgs{
		box:           b,
		writer:        w,
		machineStatus: constants.StatusRestarting,
		machineState:  constants.StateStopped,
		provisioner:   p,
	}

	if err := pipeline.Execute(args); err != nil {
		log.Errorf("error on execute restart pipeline for box %s - %s", b.GetFullName(), err)
		return err
	}
	fmt.Fprintf(w, lb.W(lb.VM_RESTARTING, lb.INFO, fmt.Sprintf("--- restart box (%s) OK\n", b.GetFullName())))
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
