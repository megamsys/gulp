package deploy

import (
	//	"bytes"
	//	"errors"
	//	"fmt"
	//	"github.com/tsuru/config"
	//	"github.com/megamsys/gulp/scm"
	"gopkg.in/check.v1"
)

func (s *S) TestDeploy(c *check.C) {
	/*	provisioner := testing.NewFakeProvisioner()
		provisioner.PrepareOutput([]byte("cloned"))
		provisioner.PrepareOutput([]byte("updated"))
		app := testing.NewFakeApp("cribcaged", "python", 1)
		provisioner.Provision(app)
		w := &bytes.Buffer{}
		err := Git("5734f0042844fdeb5bbc1b72b18f2dc1779cade7", w)
		c.Assert(err, check.IsNil)
		c.Assert(app.Commands, check.DeepEquals, []string{"restart"})
		c.Assert(provisioner.InstalledDeps(app), check.Equals, 1)
		cloneCommand := "git clone git://github.com/indykish/aryabhata.git test/dir --depth 1"
		c.Assert(provisioner.GetCmds(cloneCommand, app), check.HasLen, 1)
		path, _ := repository.GetPath()
		checkoutCommand := fmt.Sprintf("cd %s && git checkout 5734f0042844fdeb5bbc1b72b18f2dc1779cade7", path)
		c.Assert(provisioner.GetCmds(checkoutCommand, app), check.HasLen, 1)
	*/
}

func (s *S) TestDeployLogsActions(c *check.C) {
	/*provisioner := testing.NewFakeProvisioner()
		provisioner.PrepareOutput([]byte(""))
		provisioner.PrepareOutput([]byte("updated"))
		app := testing.NewFakeApp("cribcaged", "python", 1)
		provisioner.Provision(app)
		w := &bytes.Buffer{}
		err := Git(provisioner, app, "5734f0042844fdeb5bbc1b72b18f2dc1779cade7", w)
		c.Assert(err, check.IsNil)
		logs := w.String()
		expected := `
	 ---> Megam receiving push

	 ---> Replicating the application repository across units

	 ---> Installing dependencies

	 ---> Restarting application
	Restarting app...
	 ---> Deploy done!

	`
		c.Assert(logs, check.Equals, expected)
	*/
}

func (s *S) TestCloneRepository(c *check.C) {
	/*p := testing.NewFakeProvisioner()
	p.PrepareOutput([]byte("something"))
	app := testing.NewFakeApp("your", "python", 1)
	out, err := clone(p, app)
	c.Assert(err, check.IsNil)
	c.Assert(string(out), check.Equals, "something")
	url := repository.ReadOnlyURL(app.GetName())
	path, _ := repository.GetPath()
	expectedCommand := fmt.Sprintf("git clone %s %s --depth 1", url, path)
	c.Assert(p.GetCmds(expectedCommand, app), check.HasLen, 1)
	*/
}

func (s *S) TestCloneRepositoryUndefinedPath(c *check.C) {
	/*	old, _ := config.Get("git:unit-repo")
		config.Unset("git:unit-repo")
		defer config.Set("git:unit-repo", old)
		_, err := clone(nil, nil)
		c.Assert(err, check.NotNil)
		c.Assert(err.Error(), check.Equals, `Megam is misconfigured: key "git:unit-repo" not found`)\
	*/
}

func (s *S) TestPullRepository(c *check.C) {
	/*p := testing.NewFakeProvisioner()
	p.PrepareOutput([]byte("pulled"))
	app := testing.NewFakeApp("your", "python", 1)
	out, err := fetch(p, app)
	c.Assert(err, check.IsNil)
	c.Assert(string(out), check.Equals, "pulled")
	path, _ := repository.GetPath()
	expectedCommand := fmt.Sprintf("cd %s && git fetch origin", path)
	c.Assert(p.GetCmds(expectedCommand, app), check.HasLen, 1)
	*/
}
