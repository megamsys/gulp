package scm

import (
	"github.com/tsuru/config"
	"log"
)

func GetRemotePath() (string, error) {
	return config.GetString("scm:remote_repo")
}

// GetPath returns the path to the repository where the app code is in its
// units.
func GetPath() (string, error) {
	return config.GetString("scm:local_repo")
}

func Builder() (string, error) {
	return config.GetString("scm:builder")
}

func Project() (string, error) {
	return config.GetString("scm:project")
}

func ServerURL() string {
	server, err := config.GetString("scm:api_server")
	if err != nil {
		log.Print("scm:api-server config not found")
		panic(err)
	}
	return server
}
