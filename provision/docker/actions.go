package docker

import (
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/megamsys/gulp/carton"
	//"github.com/megamsys/gulp/provision"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/exec"
	//constants "github.com/megamsys/libgo/utils"
	"io"
	"strings"
)

const (
	PIPEWORK       = "pipework"
	DOCKER_PATH    = "/var/lib/docker/"
	JSON_LOG       = "-json.log"
	CONTAINER_PATH = "/var/lib/docker/containers/"
)

type runNetworkActionsArgs struct {
	Id      string
	IpAddr  string
	Bridge  string
	Gateway string
	Command string
	HomeDir string
	CartonId  string
}

type runLogsActionsArgs struct {
	Id        string
	Name      string
	Command   string
	HomeDir   string
	Writer    io.Writer
	CloseChan chan bool
}

var setNetwork = action.Action{
	Name: "attach-network-docker",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runNetworkActionsArgs)
	network_command := args.HomeDir + "/" + PIPEWORK + " " + args.Bridge + " " + parseID(args.Id) + " " + args.IpAddr + "/27@" + args.Gateway
		args.Command = network_command
		return networkExecutor(&args)
	},
	Backward: func(ctx action.BWContext) {
	//	_ = provision.EventNotify(constants.StatusContainerNetworkFailure)
	},
}

var setLogs = action.Action{
	Name: "attachlogs-docker",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runLogsActionsArgs)
		return logExecutor(&args)
	},
	Backward: func(ctx action.BWContext) {

	},
}

var updateIpsInScylla = action.Action{
	Name: "update-ips-riak",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(runNetworkActionsArgs)
		err := args.Ips()
		if err != nil {
			return nil, err
		}
		return err, nil
	},
	Backward: func(ctx action.BWContext) {

	},
}
func (n *runNetworkActionsArgs) Ips() error {
	var ips = make(map[string][]string)
	pubipv4s := []string{}
	pubipv4s = []string{n.IpAddr}
 ips[carton.PUBLICIPV4] = pubipv4s
	if asm, err := carton.NewAmbly(n.CartonId); err != nil {
		return err
	} else if err = asm.NukeAndSetOutputs(ips); err != nil {
		return err
	}
	return nil
}

func networkExecutor(networks *runNetworkActionsArgs) (action.Result, error) {
	var e exec.OsExecutor
	var commandWords []string
	commandWords = strings.Fields(networks.Command)
	if len(commandWords) > 0 {
		if err := e.Execute(commandWords[0], commandWords[1:], nil, nil, nil); err != nil {
		//	_ = provision.EventNotify(constants.StatusContainerNetworkFailure)
			return nil, err
		}
	}
	//_ = provision.EventNotify(constants.StatusContainerNetworkSuccess)
	return &networks, nil
}

func logExecutor(logs *runLogsActionsArgs) (action.Result, error) {
	filePath := CONTAINER_PATH + logs.Id + "/" + logs.Id + JSON_LOG
	cs := make(chan []byte)
	go tailLog(cs, filePath, logs.Writer, logs.CloseChan)
	return &logs, nil
}

func tailLog(cs chan []byte, filePath string, w io.Writer, closechan chan bool) {
	t, _ := tail.TailFile(filePath, tail.Config{Follow: true})
	for line := range t.Lines {
		fmt.Fprintf(w, line.Text)
	}
	closechan <- true
}

func parseID(id string) string {
	if len(strings.TrimSpace(id)) > 12 {
		return string([]rune(id)[:12])
	}
	return id
}
