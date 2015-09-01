package handler



import (
//  "github.com/megamsys/libgo/db"
"log"

"github.com/megamsys/gulp/control"
"github.com/megamsys/gulp/handlers"
//"github.com/tsuru/config"
//"strings"
)


/*
 * service-deployd ---calls----> this handler
 *
 */


func Handler(requestId string) (error) {

  log.Print("[x] Crunching the request data, hold on! ")
 /*riakUrl, err := config.GetString("riak:url")
	if err != nil {
		log.Fatal(err)
	}

  reqBucket, berr := config.GetString("riak:bucket")
	if berr != nil {
		log.Fatal(berr)
	}
	*/

riakUrl := "127.0.0.1:8087"
reqBucket := "catreqs"


  conn, cerr := handlers.RiakConnection(riakUrl, reqBucket)
  defer conn.Close()
  if cerr != nil {
    log.Print(cerr)
  }
  req := &handlers.Request{}
  log.Print("testing---->>>")
  ferr := conn.FetchStruct(requestId, req)
  if ferr != nil {
     log.Print(ferr)
     log.Print("------>>>")
   }

  switch req.Category {

  case "control":
         control.ControlHandler(req)
  //case "state":
  //       control.StateHandler(req)
  //case "policy":
  //       control.PolicyHandler(req)
   }

   return nil
}
