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
	ToscaType  string   `json:"tosca_type,omitempty"`
	RepoURL    string   `json:"scm,omitempty"`
	RepoSource string   `json:"provider,omitempty"`
	Version    string   `json:"version,omitempty"`
}

// Repos for Bitnami
type ReposBitnami struct {
	RunList         []string `json:"run_list,omitempty"`
	ToscaType       string   `json:"tosca_type,omitempty"`
	BitnamiURL      string   `json:"bitnami_url,omitempty"`
	BitnamiUserName string   `json:"bitnami_username,omitempty"`
	BitnamiPassword string   `json:"bitnami_password,omitempty"`
	BitnamiEmail    string   `json:"bitnami_email,omitempty"`
	BitnamiDBPassword string `json:"bitnami_database_password,omitempty"`
	OwncloudSite      string   `json:"bitnami_owncloud_site,omitempty"`
	PrestashopSite    string   `json:"bitnami_prestashop_site,omitempty"`
	RepoSource        string   `json:"provider,omitempty"`
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
		&mileStoneUpdate,
		&updateStatusInScylla,
	}

	pipeline := action.NewPipeline(actions...)

	args := runMachineActionsArgs{
		box:           box,
		writer:        w,
		machineStatus: constants.StatusBootstrapping,
		machineState:  constants.StateBootstrapped,
		provisioner:   p,
		state:         carton.STATE,
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
	p.Attributes = string(p.setBitnamiAttributes(b))
	p.Format = DefaultFormat
	p.LogLevel = DefaultLogLevel
	p.RootPath = meta.MC.Dir
	p.Sudo = DefaultSudo
	return p.kickOffSolo(b, w)
}

func (p *chefsoloProvisioner) setBitnamiAttributes(b *provision.Box) []byte {
	var repo, src, ip string
	if b.Repo != nil {
		repo = b.Repo.Gitr()
		src = b.Repo.RepoProvider()
	}
	bitAtr := &ReposBitnami{
		RunList:    []string{"recipe[" + p.Cookbook + "]"},
		ToscaType:  b.GetShortTosca(),
		BitnamiURL: repo,
		RepoSource: src,
	}

	if b.Outputs[carton.PUBLICIPV4] != "" {
    ip = b.Outputs[carton.PUBLICIPV4]
	} else if b.Outputs[carton.PRIVATEIPV4] != "" {
		ip = b.Outputs[carton.PRIVATEIPV4]
	}
	 for _,v := range provision.BitnamiAttributes {
		 switch true {
		 case v == provision.BITUSERNAME && b.Inputs[provision.BITUSERNAME] != "":
				bitAtr.BitnamiUserName = b.Inputs[provision.BITUSERNAME]
				bitAtr.BitnamiEmail = b.Inputs[provision.BITUSERNAME]
		 case v == provision.BITPASSWORD && b.Inputs[provision.BITPASSWORD] != "":
				bitAtr.BitnamiPassword = b.Inputs[provision.BITPASSWORD]
	   case v == provision.BITNAMI_DB_PASSWORD && b.Environments[provision.BITNAMI_DB_PASSWORD] != "":
			  bitAtr.BitnamiDBPassword = b.Inputs[provision.BITPASSWORD]
		 case v == provision.BITNAMI_PROSTASHOP_IP && b.Environments[provision.BITNAMI_PROSTASHOP_IP] != "":
		    bitAtr.PrestashopSite = ip
	  case v == provision.BITNAMI_OWNCLOUD_IP && b.Environments[provision.BITNAMI_OWNCLOUD_IP] != "":
	    	bitAtr.OwncloudSite = ip
		 }
	 }

	res, _ := json.Marshal(bitAtr)
  return res
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
	soloAction = append(soloAction, &setFinalState, &changeDoneNotify, &mileStoneUpdate, &updateStatusInScylla)
	actions := soloAction
	pipeline := action.NewPipeline(actions...)
	args := runMachineActionsArgs{
		box:           b,
		writer:        w,
		machineStatus: constants.StatusChefConfigSetupping,
		machineState:  constants.StateRunning,
		provisioner:   p,
		state:         carton.DONE,
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
		&mileStoneUpdate,
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
		&mileStoneUpdate,
		&updateStatusInScylla,
	}
	pipeline := action.NewPipeline(actions...)
	args := runMachineActionsArgs{
		box:           b,
		writer:        w,
		machineStatus: constants.StatusStopping,
		machineState:  constants.StateStopped,
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
		&mileStoneUpdate,
		&startBox,
		&mileStoneUpdate,
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


func (p *chefsoloProvisioner) ResetPassword(b *provision.Box, w io.Writer) error {
	fmt.Fprintf(w, lb.W(lb.VM_UPGRADING, lb.INFO, fmt.Sprintf("\n--- reset machine root password (%s)\n", b.GetFullName())))
	actions := []*action.Action{
		&updateStatusInScylla,
		&resetNewPassword,
		&updateStatusInScylla,
	}
	pipeline := action.NewPipeline(actions...)
	args := runMachineActionsArgs{
		box:           b,
		writer:        w,
		machineStatus: constants.StatusResetPassword,
		provisioner:   p,
	}

	if err := pipeline.Execute(args); err != nil {
		log.Errorf("error on execute reset password pipeline for machine %s - %s", b.GetFullName(), err)
		return err
	}
	fmt.Fprintf(w, lb.W(lb.VM_UPGRADING, lb.INFO, fmt.Sprintf("--- reset machine root password (%s) OK\n", b.GetFullName())))
	return nil
}
