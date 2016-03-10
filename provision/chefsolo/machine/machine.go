package machine

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	nsqp "github.com/crackcomm/nsqueue/producer"
	"github.com/megamsys/gulp/carton"
	"github.com/megamsys/gulp/meta"
	"github.com/megamsys/gulp/provision"
	ldb "github.com/megamsys/libgo/db"
	"net"
	"os"
	"time"
)

const (
	TOPIC         = "vms"
	SSHKEYSBUCKET = "sshkeys"
)

type SshKeys struct {
	OrgId      string `json:"org_id" cql:"org_id"`
	Name       string `json:"name" cql:"name"`
	CreatedAt  string `json:"created_at" cql:"created_at"`
	Id         string `json:"id" cql:"id"`
	JsonClaz   string `json:"json_claz" cql:"json_claz"`
	Privatekey string `json:"privatekey" cql:"privatekey"`
	Publickey  string `json:"publickey" cql:"publickey"`
}

type Machine struct {
	Name      string
	Id        string
	CartonId  string
	CartonsId string
	Level     provision.BoxLevel
	SSH       provision.BoxSSH
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
	ips := m.findIps()

	log.Debugf("  find and setips of machine (%s, %s)", m.Id, m.Name)

	if asm, err := carton.NewAmbly(m.CartonId); err != nil {
		return err
	} else if err = asm.NukeAndSetOutputs(ips); err != nil {
		return err
	}
	return nil
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
	ips[carton.PUBLICIPV4] = pubipv4s
	ips[carton.PRIVATEIPV4] = priipv4s
	return ips
}

// append user sshkey into authorized_keys file
func (m *Machine) AppendAuthKeys() error {
	c := &SshKeys{}
	ops := ldb.Options{
		TableName:   SSHKEYSBUCKET,
		Pks:         []string{"Name"},
		Ccms:        []string{},
		Hosts:       meta.MC.Scylla,
		Keyspace:    meta.MC.ScyllaKeyspace,
		PksClauses:  map[string]interface{}{"Name": m.SSH.Pub()},
		CcmsClauses: make(map[string]interface{}),
	}
	if err := ldb.Fetchdb(ops, c); err != nil {
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
