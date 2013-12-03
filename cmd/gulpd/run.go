package main

import (
	"github.com/indykish/gulp/cmd"
	"launchpad.net/gnuflag"
	"log"
)

type GulpStart struct {
	manager *cmd.Manager
	fs      *gnuflag.FlagSet
	dry     bool
}

func (g *GulpStart) Info() *cmd.Info {
	desc := `starts the gulpd daemon, and connects to queue.

If you use the '--dry' flag gulpd will do a dry run(parse conf/jsons) and exit.

`
	return &cmd.Info{
		Name:    "start",
		Usage:   `start [--dry] [--config]`,
		Desc:    desc,
		MinArgs: 0,
	}
}

func (c *GulpStart) Run(context *cmd.Context) error {
	log.Printf("arg 0    = %s", context.Args[0])
	log.Printf("manager  = %s", c.manager)
	// The struc will also have the c.manager
	// c.manager
	// Now using this value start the queue.
	RunServer(c.dry)
	return nil
}

func (c *GulpStart) Flags() *gnuflag.FlagSet {
	if c.fs == nil {
		c.fs = gnuflag.NewFlagSet("gulpd", gnuflag.ExitOnError)
		c.fs.BoolVar(&c.dry, "config", false, "config: the configuration file to use")
		c.fs.BoolVar(&c.dry, "c", false, "dry-run: does not start the gulpd (for testing purpose)")
		c.fs.BoolVar(&c.dry, "dry", false, "dry-run: does not start the gulpd (for testing purpose)")
		c.fs.BoolVar(&c.dry, "d", false, "dry-run: does not start the gulpd (for testing purpose)")
	}
	return c.fs
}

type GulpStop struct {
	fs   *gnuflag.FlagSet
	bark bool
}

func (g *GulpStop) Info() *cmd.Info {
	desc := `stops the gulpd daemon, and shutsdown the queue.

If you use the '--bark' flag gulpd will notify daemon status.

`
	return &cmd.Info{
		Name:    "stop",
		Usage:   `stop [--bark]`,
		Desc:    desc,
		MinArgs: 0,
	}
}

func (c *GulpStop) Run(context *cmd.Context) error {
	//api.RunServer(c.bark)
	// The struc will also have the started Handler to the Queue
	// c.handler
	// Now using the handler call
	// c.handler.stop()
	//on successful stop, propagate the status if there was no err and bark is turned on..

	/*appName, err := c.dry()
	if err != nil {
		return err
	}

	url, err := cmd.GetURL(fmt.Sprintf("/apps/%s/run?once=%t", appName, c.once))
	if err != nil {
		return err
	}
	b := strings.NewReader(strings.Join(context.Args, " "))
	request, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}
	r, err := client.Do(request)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	_, err = io.Copy(context.Stdout, r.Body)

	return err
	*/
	return nil
}

func (c *GulpStop) Flags() *gnuflag.FlagSet {
	if c.fs == nil {
		c.fs = gnuflag.NewFlagSet("gulpd", gnuflag.ExitOnError)
		c.fs.BoolVar(&c.bark, "bark", false, "bark: does a notify of the daemon status (to zk)")
		c.fs.BoolVar(&c.bark, "b", false, "bark: does a notify of the daemon status (to zk)")
	}
	return c.fs
}
