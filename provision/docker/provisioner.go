package docker

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/carton"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/libgo/action"
	"io"
)

type DockerProvisioner struct {
	Id          string
	ContainerId string
	Name        string
	IpAddr      string
	Gateway     string
	Bridge      string
	HomeDir     string
}

func (p *DockerProvisioner) Initialize(m string) error {
	return nil
}

func (p *DockerProvisioner) LogExec() {
	var outBuffer bytes.Buffer
	var closeChan chan bool

	box := &provision.Box{Id: p.ContainerId, Name: p.Name}
	logWriter := carton.LogWriter{Box: box}
	logWriter.Async()

	writer := io.MultiWriter(&outBuffer, &logWriter)
	p.createLogPipeline(writer, closeChan, &logWriter)

	go func(closeChan chan bool, logWriter carton.LogWriter) {
		select {
		case <-closeChan:
			logWriter.Close()
		default:
		}
	}(closeChan, logWriter)
}

func (p *DockerProvisioner) NetworkExec() {
	p.createNetworkPipeline()
}

func (p *DockerProvisioner) createNetworkPipeline() error {
	actions := []*action.Action{
		&setNetwork,
	}
	pipeline := action.NewPipeline(actions...)
	args := runNetworkActionsArgs{
		Id:      p.ContainerId,
		IpAddr:  p.IpAddr,
		Bridge:  p.Bridge,
		Gateway: p.Gateway,
		HomeDir: p.HomeDir,
	}

	err := pipeline.Execute(args)
	if err != nil {
		log.Errorf("Error on executing Network setup")
		return err
	}
	return nil
}

func (p *DockerProvisioner) createLogPipeline(writer io.Writer, closeChan chan bool, logWriter *carton.LogWriter) error {
	actions := []*action.Action{
		&setLogs,
	}
	pipeline := action.NewPipeline(actions...)
	args := runLogsActionsArgs{
		Id:        p.ContainerId,
		Name:      p.Name,
		HomeDir:   p.HomeDir,
		Writer:    writer,
		CloseChan: closeChan,
		LogWriter: logWriter,
	}

	err := pipeline.Execute(args)
	if err != nil {
		log.Errorf("Error on executing Log setup")
		return err
	}
	return nil
}
