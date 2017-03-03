package machine

import (
	"encoding/json"
	b64 "encoding/base64"
	"fmt"
	log "github.com/Sirupsen/logrus"
	nsqp "github.com/crackcomm/nsqueue/producer"
	"github.com/megamsys/gulp/carton"
	"github.com/megamsys/gulp/meta"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/libgo/utils"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	TOPIC         = "vms"
	SSHKEYSBUCKET = "sshkeys"
	cmd           = `#!/bin/sh
id -u %s &>/dev/null || useradd %s
/usr/bin/passwd %s <<EOF
%s
%s
EOF
`
resetPwd = `echo %s:%s | /usr/sbin/chpasswd
`
)

type Machine struct {
	Name      string
	Id        string
	CartonId  string
	CartonsId string
	OrgId     string
	Level     provision.BoxLevel
	SSH       provision.BoxSSH
	PublicIp  string
	Status    utils.Status
	State     utils.State
}

func (m *Machine) SetStatus(status utils.Status) error {
	log.Debugf("  set status[%s] of machine (%s, %s)", m.Id, m.Name, status.String())

 if asm, err := carton.NewAssembly(m.CartonId); err != nil {
		return err
	} else if err = asm.SetStatus(status); err != nil {

		return err
	}
	if m.Level == provision.BoxSome {
		log.Debugf("  set status[%s] of machine (%s, %s)", m.Id, m.Name, status.String())

		if comp, err := carton.NewComponent(m.Id); err != nil {
			return err
		} else if err = comp.SetStatus(status); err != nil {
			return err
		}
	}
	return nil
}

func (m *Machine) SetState(state utils.State) error {
	log.Debugf("  set state[%s] of machine (%s, %s)", m.Id, m.Name, state.String())


	if asm, err := carton.NewAssembly(m.CartonId); err != nil {
		return err
	} else if err = asm.SetState(state); err != nil {

		return err
	}

	if m.Level == provision.BoxSome {
		log.Debugf("  set state[%s] of machine (%s, %s)", m.Id, m.Name, state.String())

		if comp, err := carton.NewComponent(m.Id); err != nil {
			return err
		} else if err = comp.SetState(state); err != nil {
			return err
		}
	}
	return nil
}

// FindAndSetIps returns the non loopback local IP4 (can be public or private)
// we also have to add it in for ipv6
func (m *Machine) FindAndSetIps(b *provision.Box) error {
	ips := m.mergeSameIPtype(m.findIps())
	log.Debugf("  find and setips of machine (%s, %s)", m.Id, m.Name)
  asm, err := carton.NewAssembly(m.CartonId)
	if  err != nil {
		return err
	} else if err = asm.NukeAndSetOutputs(ips); err != nil {
		return err
	}
  b.Outputs = asm.Outputs.ToMap()
	return nil
}


func (m *Machine) mergeSameIPtype(mm map[string][]string)  map[string][]string {
  for IPtype, ips := range mm {
		var sameIp string
		for _, ip := range ips {
			sameIp = sameIp +  ip + ", "
		}
		if sameIp != "" {
			mm[IPtype] = []string{strings.TrimRight(sameIp, ", ")}
		}
	}
	return mm
}

// FindIps returns the non loopback local IP4 (can be public or private)
// if an iface contains a string "pub", then we consider it a public interface
func (m *Machine) findIps() map[string][]string {
	var ips = make(map[string][]string)
	pubipv4s := []string{}
	priipv4s := []string{}

	ifaces, err := net.Interfaces()
	if err != nil {
		return ips
	}
	for _, iface := range ifaces {
		ifaddress, err := iface.Addrs()
		if err != nil {
			return ips
		}
		for _, address := range ifaddress {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsMulticast() {
				if ip4 := ipnet.IP.To4(); ip4 != nil {
					if ip4[0] == 192 || ip4[0] == 10 || ip4[0] == 172 {
						priipv4s = append(priipv4s, ipnet.IP.String())
					} else {
						pubipv4s = append(pubipv4s, ipnet.IP.String())
					}
				}
			}
		}
	}
	ips[utils.PUBLICIPV4] = pubipv4s
	ips[utils.PRIVATEIPV4] = priipv4s
	return ips
}

// append user sshkey into authorized_keys file
func (m *Machine) AppendAuthKeys() error {
	if strings.TrimSpace(m.SSH.Password) == "" || strings.TrimSpace(m.SSH.User) == "" {
    c, err := carton.GetSSHKeys(m.SSH.Pub())
		if err != nil {
			return err
		}

		f, err := os.OpenFile(m.SSH.AuthKeysFile(), os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}

		defer f.Close()

		if _, err = f.WriteString(c.Publickey); err != nil {
			return err
		}
	} else {
		pwd, _ := b64.StdEncoding.DecodeString(m.SSH.Password)
		d1 := []byte(fmt.Sprintf(cmd, m.SSH.User,m.SSH.User,m.SSH.User, string(pwd), string(pwd)))
		err := ioutil.WriteFile("dat1", d1, 0755)
		_, err = exec.Command("./dat1").Output()
		if err != nil {
			return err
		}
		err = os.Remove("./dat1")
		if err != nil {
			return err
		}

		if asm, err := carton.NewAssembly(m.CartonId); err != nil {
			return err
		} else if err = asm.NukeKeysInputs(carton.PASSWORD); err != nil {
			return err
		}

	}

	return nil
}

func (m *Machine) ChangeState(state string) error {
	log.Debugf("  change state of machine (%s, %s)", m.Name, m.Status.String())

	pons := nsqp.New()
	if err := pons.Connect(meta.MC.NSQd[0]); err != nil {
		return err
	}
	bytes, err := json.Marshal(
		carton.Requests{
			CatId:     m.CartonsId,
			Action:    m.Status.String(),
			AccountId: meta.MC.AccountId,
			Category:  state,
			CreatedAt: time.Now().String(),
		})

	if err != nil {
		return err
	}

	log.Debugf("  pub to topic (%s, %s)", TOPIC, bytes)
	if err = pons.Publish(TOPIC, bytes); err != nil {
		return err
	}

	defer pons.Stop()
	return nil
}

func (m *Machine) ResetPassword() error {
	pwd, _ := b64.StdEncoding.DecodeString(m.SSH.Password)
	_, err := exec.Command("sh","-c",fmt.Sprintf(resetPwd, m.SSH.User, string(pwd))).Output()
	if err != nil {
		return err
	}

	if asm, err := carton.NewAssembly(m.CartonId); err != nil {
		return err
	} else if err = asm.NukeKeysInputs(carton.PASSWORD); err != nil {
		return err
	}

	log.Debugf("  password updated successfully")
	return nil
}
