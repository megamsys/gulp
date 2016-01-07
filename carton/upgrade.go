package carton

import (
	"fmt"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/gulp/upgrade"
	"github.com/megamsys/libgo/action"
)

type Upgradeable struct {
	B             *provision.Box
	w             io.Writer
	ShouldRestart bool
}

func NewUpgradeable(box *provision.Box) *Upgradeable {
	u := &Upgradeable{
		B:             box,
		ShouldRestart: true,
	}
	u.register()
	return u
}

func (u *Upgradeable) register() {
	err := upgrade.Register("ci", u.opsBuild)
	if err != nil {
		log.Fatalf("unable to register ops ci: %s", err)
	}
	err = upgrade.Register("bind", u.opsBind)
	if err != nil {
		log.Fatalf("unable to register ops bind: %s", err)
	}
}

// Boot runs the boot of the vm.
func (u *Upgradeable) Upgrade() error {
	logWriter := NewLogWriter(u.B)
	defer logWriter.Close()
	writer := io.MultiWriter(&logWriter)
	err := u.operateBox(writer)
	if err != nil {
		return err
	}
	return nil
}

func (u *Upgradeable) operateBox(writer io.Writer) error {
	u.w = writer
	start := time.Now()
	opsRan, err := upgrade.Run(upgrade.RunArgs{
		Name:   u.B.GetFullName(),
		O:      u.B.Operations,
		Writer: writer,
		Force:  false,
	})
	if err != nil {
		return err
	}

	if !opsRan.Successful() {
		return nil
	}
	elapsed := time.Since(start)

	if err := u.saveData(opsRan, elapsed); err != nil {
		log.Errorf("WARNING: couldn't save ops data, ops opts: %#v", u)
		return err
	}

	if !u.ShouldRestart {
		return nil
	}

	return Provisioner.Restart(u.B, u.w)
}

func (u *Upgradeable) opsBuild() error {
	fmt.Fprintf(u.w, "---- ops ci (%s) is kicking ----\n", u.B.GetFullName())

	actions := []*action.Action{
		&cloneBox,
		&buildBox, //buildpack does everthing
	}
	pipeline := action.NewPipeline(actions...)
	args := runOpsPipelineArgs{
		box:    u.B,
		writer: u.w,
	}
	if err := pipeline.Execute(args); err != nil {
		return err
	}

	return nil
}

func (u *Upgradeable) opsBind() error {
	fmt.Fprintf(u.w, "---- ops bind (%s) is kicking ----\n", u.B.GetFullName())
	actions := []*action.Action{
		&setEnvsAction,
	}
	pipeline := action.NewPipeline(actions...)
	args := runOpsPipelineArgs{
		box:    u.B,
		writer: u.w,
	}
	if err := pipeline.Execute(args); err != nil {
		return err
	}
	return nil
}

func (u *Upgradeable) saveData(opsRan upgrade.OperationsRan, elapsed time.Duration) error {
	if u.B.Level == provision.BoxSome {
		log.Debugf("  update operation run for box (%s, %s)", u.B.Id, u.B.GetFullName())
		if comp, err := NewComponent(u.B.Id); err != nil {
			return err
		} else if err = comp.UpdateOpsRun(opsRan); err != nil {
			return err
		}
	}
	return nil
}
