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
package provision

import (
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/megamsys/gulp/carton/bind"
	"github.com/megamsys/gulp/repository"
	_ "github.com/megamsys/gulp/repository/github"
	"github.com/megamsys/gulp/upgrade"
	"github.com/megamsys/libgo/utils"
	constants "github.com/megamsys/libgo/utils"
	"gopkg.in/yaml.v2"
)

const (
	CPU = "cpu"
	RAM = "ram"
	HDD = "hdd"

	BITUSERNAME = "bitnami_username"
	BITPASSWORD = "bitnami_password"
	BITNAMI_DB_PASSWORD = "bitnami_database_password"
	BITNAMI_PROSTASHOP_IP = "bitnami_prestashop_site"
	BITNAMI_OWNCLOUD_IP = "bitnami_owncloud_site"


	// BoxSome indicates that there is atleast one box to deploy or delete.
	BoxSome BoxLevel = iota

	// BoxNone indicates that there are no boxes to deploy or delete but its parent can be.
	BoxNone
)

var BitnamiAttributes = []string{BITUSERNAME, BITPASSWORD,BITNAMI_DB_PASSWORD, BITNAMI_PROSTASHOP_IP, BITNAMI_OWNCLOUD_IP}
var cnameRegexp = regexp.MustCompile(`^(\*\.)?[a-zA-Z0-9][\w-.]+$`)

// Boxlevel represents the deployment level.
type BoxLevel int

// Boxlog represents a log entry.
type Boxlog struct {
	Timestamp string
	Message   string
	Source    string
	Name      string
	Unit      string
}

type BoxSSH struct {
	User   string
	Prefix string
	Password string
}

func (bs *BoxSSH) Pub() string {
	//return bs.Prefix + "_pub"
	return bs.Prefix
}

//authorized_keys path is same in all linux i think
func (bs *BoxSSH) AuthKeysFile() string {
	dotssh_dir := ""
	dotssh := ""
	switch runtime.GOOS {
	case "linux":
		dotssh_dir = filepath.Join(home(bs.User), ".ssh")
		dotssh = filepath.Join(dotssh_dir, "authorized_keys")
	default:
		dotssh_dir = filepath.Join(home(bs.User), ".ssh")
		dotssh = filepath.Join(dotssh_dir, "authorized_keys")
	}

	if _, err := os.Stat(dotssh_dir); err != nil { //create  authorized_keys file, if it aint there
		os.Mkdir(dotssh_dir, 755)
	}

	if _, err := os.Stat(dotssh); err != nil { //create  authorized_keys file, if it aint there
		w, _ := os.Create(dotssh)
		defer w.Close()
	}
	return dotssh
}

func home(name string) string {
	if auth_user, err := user.Lookup(name); err == nil {
		return auth_user.HomeDir
	}
	curr_user, _ := user.Current()
	return curr_user.HomeDir // hmm no error trap ?
}

type BoxCompute struct {
	Cpushare string
	Memory   string
	Swap     string
	HDD      string
}

func (bc *BoxCompute) numCpushare() int64 {
	if cs, err := strconv.ParseInt(bc.Cpushare, 10, 64); err != nil {
		return 0
	} else {
		return cs
	}
}

func (bc *BoxCompute) numMemory() int64 {
	if cp, err := strconv.ParseInt(bc.Memory, 10, 64); err != nil {
		return 0
	} else {
		return cp
	}
}

func (bc *BoxCompute) numSwap() int64 {
	if cs, err := strconv.ParseInt(bc.Swap, 10, 64); err != nil {
		return 0
	} else {
		return cs
	}
}

func (bc *BoxCompute) numHDD() int64 {
	if cp, err := strconv.ParseInt(bc.HDD, 10, 64); err != nil {
		return 10
	} else {
		return cp
	}
}

func (bc *BoxCompute) String() string {
	return "(" + strings.Join([]string{
		CPU + ":" + bc.Cpushare,
		RAM + ":" + bc.Memory,
		HDD + ":" + bc.HDD},
		",") + " )"
}

// BoxDeploy represents a log entry.
type BoxDeploy struct {
	Date    time.Time
	HookId  string
	ImageId string
	Name    string
	Unit    string
}

// Box represents a provision unit. Can be a machine, container or anything
// IP-addressable.
type Box struct {
	Id           string
	CartonsId    string
	CartonId     string
	CartonName   string
	Name         string
	Level        BoxLevel
	DomainName   string
	Inputs       map[string]string
	Outputs      map[string]string
	Tosca        string
	ImageVersion string
	Compute      BoxCompute
	SSH          BoxSSH
	PublicIp     string
	Repo         *repository.Repo
	Status       utils.Status
	State        utils.State
	Provider     string
	Commit       string
	Environments map[string]string
	Envs         bind.EnvVars
	Address      *url.URL
	Operations   []*upgrade.Operation //MEGAMD
}

func (b *Box) String() string {
	if d, err := yaml.Marshal(b); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

func (b *Box) GetMemory() int64 {
	return b.Compute.numMemory()
}

func (b *Box) GetSwap() int64 {
	return b.Compute.numSwap()
}

func (b *Box) GetCpushare() int64 {
	return b.Compute.numCpushare()
}

// GetName returns the assemblyname.domain(assembly001YeahBoy.megambox.com) of the box.
func (b *Box) GetFullName() string {
	if len(strings.TrimSpace(b.DomainName)) > 0 {
		return strings.Join([]string{b.CartonName, b.DomainName}, ".")
	}
	return b.CartonName
}

func (b *Box) GetShortTosca() string {
	a := strings.Split(b.Tosca, ".")
	if a[0] == "bitnami" {
		return a[0]
	} else {
		return a[2]
	}
}

// GetIp returns the Unit.IP.
func (b *Box) GetPublicIp() string {
	return b.PublicIp
}

func (box *Box) GetRouter() (string, error) {
	return "route53", nil //dns.LoadConfig()
}

func (b *Box) Clone() error {
	if b.Repo != nil && b.Repo.Type != repository.IMAGE && !b.Repo.OneClick {
		scm := repository.Manager(b.Repo.Source)
		if scm == nil {
			return fmt.Errorf("couldn't locate the repository manager (%s)", b.Repo.Source)
		}
		if err := scm.Clone(b.Repo); err != nil {
			return err
		}
	}
	return nil
}

// Available returns true if the unit is available. It will return true
// whenever the unit itself is available, even when the application process is
// not.
func (b *Box) Available() bool {
	return b.Status == constants.StatusBootstrapping ||
		b.Status == constants.StatusRunning ||
		b.Status == constants.StatusBootstrapped ||
		b.Status == constants.StatusStateupped ||
		b.Status == constants.StatusError ||
		b.Status == constants.StatusStarted ||
		b.Status == constants.StatusStopped ||
		b.Status == constants.StatusRestarted
}

// Log adds a log message to the app. Specifying a good source is good so the
// user can filter where the message come from.
func (box *Box) Log(message, source, unit string) error {
	messages := strings.Split(message, "\n")
	logs := make([]interface{}, 0, len(messages))
	for _, msg := range messages {
		if len(strings.TrimSpace(msg)) > 0 {
			bl := Boxlog{
				Timestamp: time.Now().Local().Format(time.RFC822),
				Message:   msg,
				Source:    source,
				Name:      box.Name,
				Unit:      box.Id,
			}
			logs = append(logs, bl)
		}
	}
	if len(logs) > 0 {
		if box.Tosca == "docker" {
			_ = notify(box.Name, logs)
		} else {
			_ = notify(box.GetFullName(), logs)
		}
	}
	return nil
}
