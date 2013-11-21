
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/indykish/gulp/cmd"
	"io/ioutil"
	"launchpad.net/gnuflag"
	"net/http"
)

type AppCreate struct{}

func (AppCreate) Run(context *cmd.Context) error {
	appName := context.Args[0]
	platform := context.Args[1]
	b := bytes.NewBufferString(fmt.Sprintf(`{"name":"%s","platform":"%s"}`, appName, platform))
	url, err := cmd.GetURL("/apps")
	if err != nil {
		return err
	}
	/*request, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	out := make(map[string]string)
	err = json.Unmarshal(result, &out)
	if err != nil {
		return err
	}
	fmt.Fprintf(context.Stdout, "App %q is being created!\n", appName)
	fmt.Fprintln(context.Stdout, "Use app-info to check the status of the app and its units.")
	fmt.Fprintf(context.Stdout, "Your repository for %q project is %q\n", appName, out["repository_url"])
	*/
	return nil
}

func (AppCreate) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "app-create",
		Usage:   "app-create <appname> <platform>",
		Desc:    "create a new app.",
		MinArgs: 2,
	}
}

type AppRemove struct {
	yes bool
	fs  *gnuflag.FlagSet
}

func (c *AppRemove) Info() *cmd.Info {
	return &cmd.Info{
		Name:  "app-remove",
		Usage: "app-remove [--app appname] [--assume-yes]",
		Desc: `removes an app.

If you don't provide the app name, megam will try to guess it.`,
		MinArgs: 0,
	}
}

func (c *AppRemove) Run(context *cmd.Context) error {
	appName, err := c.Guess()
	if err != nil {
		return err
	}
	var answer string
	if !c.yes {
		fmt.Fprintf(context.Stdout, `Are you sure you want to remove app "%s"? (y/n) `, appName)
		fmt.Fscanf(context.Stdin, "%s", &answer)
		if answer != "y" {
			fmt.Fprintln(context.Stdout, "Abort.")
			return nil
		}
	}
	/*url, err := cmd.GetURL(fmt.Sprintf("/apps/%s", appName))
	if err != nil {
		return err
	}
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	_, err = client.Do(request)
	if err != nil {
		return err
	}
	fmt.Fprintf(context.Stdout, `App "%s" successfully removed!`+"\n", appName)
	*/
	return nil
}

func (c *AppRemove) Flags() *gnuflag.FlagSet {
	if c.fs == nil {
		c.fs = c.GuessingCommand.Flags()
		c.fs.BoolVar(&c.yes, "assume-yes", false, "Don't ask for confirmation, just remove the app.")
		c.fs.BoolVar(&c.yes, "y", false, "Don't ask for confirmation, just remove the app.")
	}
	return c.fs
}