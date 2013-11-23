package scm 


import (
	"log"
	"github.com/globocom/config"
)

func ServerURL() string {
	server, err := config.GetString("scm:api-server")
	if err != nil {
		log.Print("scm:api-server config not found")
		panic(err)
	}
	return server
}


// GetPath returns the path to the repository where the app code is in its
// units.
func GetPath() (string, error) {
	return config.GetString("scm:repo")
}


