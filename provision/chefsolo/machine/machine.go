package machine

import (
	"encoding/json"
	"net"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	nsqp "github.com/crackcomm/nsqueue/producer"
	"github.com/megamsys/gulp/carton"
	"github.com/megamsys/gulp/db"
	"github.com/megamsys/gulp/meta"
	"github.com/megamsys/gulp/provision"
)

const (
	TOPIC          = "vms"
	SSHFILESBUCKET = "sshfiles"
)

type Machine struct {
	Name      string
	Id        string
	CartonId  string
	CartonsId string
	Level     provision.BoxLevel
	PublicIp  string
	Status    provision.Status
}

func (m *Machine) SetStatus(status provision.Status) error {
	log.Debugf("  set status[%s] of machine (%s, %s)", m.Id, m.Name, status.String())

	if asm, err := carton.NewAmbly(m.CartonId); err != nil {
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

// FindAndSetIps returns the non loopback local IP4 (can be public or private)
// we also have to add it in for ipv6
func (m *Machine) FindAndSetIps() error {
	ips := make(map[string][]string)
	ips["ip"] = m.findIps()

	log.Debugf("  find and setips of machine (%s, %s)", m.Id, m.Name)

	if asm, err := carton.NewAmbly(m.CartonId); err != nil {
		return err
	} else if err = asm.NukeAndSetOutputs(ips); err != nil {
		return err
	}
	return nil
}

// FindAIps returns the non loopback local IP4 (can be public or private)
func (m *Machine) findIps() []string {
	var ips = []string{}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips
}

// append user sshkey into authorized_keys file
func (m *Machine) AppendAuthorizedKeys() error {
	sshkey, err := db.FetchObject(SSHFILESBUCKET, "m.Assembly.Sshkey()"+"_pub")
	if err != nil {
		return err
	}

	f, err := os.OpenFile("/root/.ssh/authorized_keys", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(sshkey); err != nil {
		return err
	}
	return nil
}

func (m *Machine) ChangeState(status provision.Status) error {
	log.Debugf("  change state of machine (%s, %s)", m.Name, status.String())

	pons := nsqp.New()
	if err := pons.Connect(meta.MC.NSQd[0]); err != nil {
		return err
	}

	bytes, err := json.Marshal(
		carton.Requests{
			CatId:     m.CartonsId,
			Action:    status.String(),
			Category:  carton.STATE,
			CreatedAt: time.Now().Local().Format(time.RFC822),
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
