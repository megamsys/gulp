package machine

import (
	"encoding/json"
	"fmt"
	"time"
	"os"
	"github.com/megamsys/gulp/carton"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/libgo/amqp"
	"github.com/megamsys/gulp/meta"
	"net"
	"github.com/megamsys/gulp/db"
	log "github.com/Sirupsen/logrus"

)

const ( 
	QUEUE = "cloudstandup"
	
	SSHFILESBUCKET = "sshfiles"
	)

type Machine struct {
	CatID     string
	CatsID	  string
	Assembly  *carton.Ambly
	IP        string
}


//just publish a message bootstrapped to the service.
func (m *Machine) PubStatus(status provision.Status) error {
	log.Debugf("publish bootstrapped service %s to %s", "Gulpd", status.String())

	p, err := amqp.NewRabbitMQ(meta.MC.AMQP, QUEUE)
	if err != nil {
		return err
	}
	
    //before publish the queue, we need to verify assembly status
    jsonMsg, err := json.Marshal(
		carton.Requests{
			CatId: 		m.CatsID, 
			Action:     status.String(),
			Category:   carton.STATE,
			CreatedAt:  time.Now().String(),
		})

	if err != nil {
		return err
	}	

	if err := p.Pub(jsonMsg); err != nil {
		return err
	}
	return nil
}

// GetLocalIP returns the non loopback local IP of the host
func (m *Machine) GetLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}

// append user sshkey into authorized_keys file
func (m *Machine) UpdateSshkey() error {
	sshkey, err := db.FetchObject(SSHFILESBUCKET, m.Assembly.Sshkey()+"_pub")
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

