package docker


import (
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/gulp/global"
)


func StreamLogs(logs *global.DockerLogsInfo) error {
	actions := []*action.Action{&streamLogs}

	pipeline := action.NewPipeline(actions...)
	err := pipeline.Execute(logs)
	if err != nil {
		return &AppLifecycleError{app: logs.Name, Err: err}
	}
	return nil
}
