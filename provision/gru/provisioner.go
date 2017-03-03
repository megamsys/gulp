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
package gru

import (
	"bytes"
//	"encoding/json"
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
	GRU_GIT      = "gru"
	GRU_TARBALL  = "gru_tarball"
	GRUCTL_TAR  = "gructl_tar"
)

var mainGruProvisioner *gruProvisioner

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

// Provisioner is a provisioner based on Gructl.
type gruProvisioner struct {
//	RunList    []string
	Attributes string
	Format     string
	LogLevel   string
	RootPath   string
	Sudo       bool
}

func init() {
	mainGruProvisioner = &gruProvisioner{}
	provision.Register(provision.GRU, mainGruProvisioner)
}

//initialize the provisioner and setup the requirements for provisioner
func (p *gruProvisioner) Initialize(m map[string]string) error {
	var outBuffer bytes.Buffer
	start := time.Now()

	logWriter := carton.NewLogWriter(&provision.Box{CartonName: m[NAME]})
	writer := io.MultiWriter(&outBuffer, &logWriter)
	defer logWriter.Close()

	cr := NewGruRepo(m, writer)
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

func (p *gruProvisioner) StartupMessage() (string, error) {
	w := new(tabwriter.Writer)
	var b bytes.Buffer
	w.Init(&b, 0, 8, 0, '\t', 0)
	b.Write([]byte(cmd.Colorfy("  > gructl ", "white", "", "bold") + "\t" +
		cmd.Colorfy(p.String(), "cyan", "", "")))
	fmt.Fprintln(w)
	w.Flush()
	return strings.TrimSpace(b.String()), nil
}

func (p *gruProvisioner) String() string {
	return "ready"
}

func (p *gruProvisioner) Bootstrap(box *provision.Box, w io.Writer) error {
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

func (p *gruProvisioner) Stateup(b *provision.Box, w io.Writer) error {
	fmt.Fprintf(w, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("\n--- stateup box (%s)\n", b.GetFullName())))
	var repo, src string
	if b.Repo != nil {
		repo = b.Repo.Gitr()
		src = b.Repo.RepoProvider()
	}

	atr := &Attributes{
		ToscaType:  b.GetShortTosca(),
		RepoURL:    repo,
		RepoSource: src,
		Version:    b.ImageVersion,
	}


DefaultAttributes := fmt.Sprintf( "tosca_type = \"%s\"\n  scm =  \"%s\"\n ", atr.ToscaType , atr.RepoURL )
DefaultAttributes +=  fmt.Sprintf("provider = \"%s\"\n version = \"%s\"\n",  atr.RepoSource , atr.Version )
	p.Attributes = DefaultAttributes
	p.Format = DefaultFormat
	p.LogLevel = DefaultLogLevel
	p.RootPath = meta.MC.Home
	p.Sudo = DefaultSudo
	return p.kickOffSolo(b, w)
}

func (p *gruProvisioner) StateupBitnami(b *provision.Box, w io.Writer) error {
	fmt.Fprintf(w, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("\n--- stateup box (%s)\n", b.GetFullName())))
	p.Attributes = p.setBitnamiAttributes(b)
	p.Format = DefaultFormat
	p.LogLevel = DefaultLogLevel
	p.RootPath = meta.MC.Home
	p.Sudo = DefaultSudo
	return p.kickOffSolo(b, w)
}

func (p *gruProvisioner) setBitnamiAttributes(b *provision.Box) string {
	var repo, src, ip string
	if b.Repo != nil {
		repo = b.Repo.Gitr()
		src = b.Repo.RepoProvider()
	}

	bitAtr := &ReposBitnami{
		ToscaType:  b.GetShortTosca(),
		BitnamiURL: repo,
		RepoSource: src,
	}

	DefaultAttributes := fmt.Sprintf( "tosca_type =  \"%s\"\n  bitnami_url = \"%s\"\n  provider = \"%s\"\n", bitAtr.ToscaType, bitAtr.BitnamiURL, bitAtr.RepoSource)

	if b.Outputs[constants.PUBLICIPV4] != "" {
    ip = b.Outputs[constants.PUBLICIPV4]
	} else if b.Outputs[constants.PRIVATEIPV4] != "" {
		ip = b.Outputs[constants.PRIVATEIPV4]
	}
	 for _,v := range provision.BitnamiAttributes {
		 switch true {
		 case v == provision.BITUSERNAME && b.Inputs[provision.BITUSERNAME] != "":
					bitAtr.BitnamiUserName = b.Inputs[provision.BITUSERNAME]
				bitAtr.BitnamiEmail = b.Inputs[provision.BITUSERNAME]
       	DefaultAttributes += fmt.Sprintf("bitnami_username =  \"%s\"\n  bitnami_email = \"%s\"\n", bitAtr.BitnamiUserName, bitAtr.BitnamiEmail)

		 case v == provision.BITPASSWORD && b.Inputs[provision.BITPASSWORD] != "":
					bitAtr.BitnamiPassword = b.Inputs[provision.BITPASSWORD]
				DefaultAttributes += fmt.Sprintf("bitnami_password = \"%s\"\n", bitAtr.BitnamiPassword)
	   case v == provision.BITNAMI_DB_PASSWORD && b.Environments[provision.BITNAMI_DB_PASSWORD] != "":
			  	bitAtr.BitnamiDBPassword = b.Inputs[provision.BITPASSWORD]
				DefaultAttributes +=  fmt.Sprintf("bitnami_database_password = \"%s\"\n ", bitAtr.BitnamiDBPassword)
		 case v == provision.BITNAMI_PROSTASHOP_IP && b.Environments[provision.BITNAMI_PROSTASHOP_IP] != "":
		    	bitAtr.PrestashopSite = ip
				DefaultAttributes += fmt.Sprintf("bitnami_prestashop_site = \"%s\"\n", bitAtr.PrestashopSite)
	  case v == provision.BITNAMI_OWNCLOUD_IP && b.Environments[provision.BITNAMI_OWNCLOUD_IP] != "":
	    		bitAtr.OwncloudSite = ip
				DefaultAttributes +=  fmt.Sprintf("bitnami_owncloud_site = \"%s\"\n", bitAtr.OwncloudSite)
		 }
	 }

	//res, _ := json.Marshal(bitAtr)
  return DefaultAttributes
}
//1. &prepareJSON in generate the json file for chefsolo
//2. &prepareConfig in generate the config file for gru.
//3. &updateStatus in Riak - Creating..
func (p *gruProvisioner) kickOffSolo(b *provision.Box, w io.Writer) error {
	fmt.Fprintf(w, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("\n--- kickofff gru box (%s)\n", b.GetFullName())))
	gruAction := make([]*action.Action, 0, 4)
	gruAction = append(gruAction, &updateStatusInScylla,  &generateGruParam, &updateStatusInScylla, &cloneBox, &updateStatusInScylla)
	if b.Level != provision.BoxNone {
		gruAction = append(gruAction, &setGruStatus, &updateStatusInScylla, &gructlRun, &updateStatusInScylla)
	}
	gruAction = append(gruAction, &setFinalState, &changeDoneNotify, &mileStoneUpdate, &updateStatusInScylla)
	actions := gruAction
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
		log.Errorf("error on execute gru pipeline for box %s - %s", b.GetFullName(), err)
		return err
	}
	fmt.Fprintf(w, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("--- kickofff gru box (%s) OK\n", b.GetFullName())))
	return nil
}

func (p *gruProvisioner) Start(b *provision.Box, w io.Writer) error {
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

func (p *gruProvisioner) Stop(b *provision.Box, w io.Writer) error {
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

func (p *gruProvisioner) Restart(b *provision.Box, w io.Writer) error {
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
func (p gruProvisioner) Command() []string {

	cmd := []string{
		path.Join(p.RootPath,"gru/gulp/gructl"),
		"apply", path.Join(p.RootPath, "gru/site/route/route.lua"),
	//	"--format", format,
	//	"--log_level", logLevel,
	}
	return cmd
}


func (p *gruProvisioner) ResetPassword(b *provision.Box, w io.Writer) error {
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
