package app



import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	log "code.google.com/p/log4go"
	"github.com/megamsys/gulp/global"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/exec"
	"github.com/tsuru/config"
)



func DockerLogExecutor(logs *global.DockerLogsInfo) (action.Result, error) {

  fmt.Println("Entered dockerllogexecutor")
  fmt.Println(logs)
}



var streamLogs = action.Action{
	Name: "streamLogs",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		var logs *global.DockerLogsInfo
//		cib.Command = cobbler
fmt.Println("--------------------------------------------------")
fmt.Println(logs.ContainerId)
		exec, err1 := DockerLogExecutor(&logs)
		if err1 != nil {
			fmt.Println("server insert error")
			return &cib, err1
		}
		return exec, err1

	},
	Backward: func(ctx action.BWContext) {
		//app := ctx.FWResult.(*App)
		db := orm.OpenDB()
		dbmap := orm.GetDBMap(db)
		err := orm.DeleteRowFromServerName(dbmap, "COBBLER")
		if err != nil {
			log.Printf("Server delete error")
			///return &cib, err
		}
		log.Printf(" Nothing to recover")
	},
	MinParams: 1,
}
