package control

import (
	//"fmt"
	"log"

	"github.com/megamsys/gulp/app"
	"github.com/megamsys/gulp/handlers"
)

const (
	assemblyBucket = "assembly"
	comBucket      = "components"
)

func ControlHandler(request *Request) error {

	riakUrl := "127.0.0.1:8087"
	conn, cerr := handlers.RiakConnection(riakUrl, assemblyBucket)
	defer conn.Close()
	if cerr != nil {
		return cerr
	}

	AwC, aerr := handlers.Deep(request.AppId)
	if aerr != nil {
		log.Print(aerr)
	}

	switch request.Action {
	case "reboot":
		go app.RebootApp(AwC)
		break
	case "start":
		go app.StartApp(AwC)
		break
	case "stop":
		go app.StopApp(AwC)
		break
	case "restart":
		go app.RestartApp(AwC)
		break
	case "redeploy":
		go app.BuildApp(AwC)
		break
	}

	return nil
}
