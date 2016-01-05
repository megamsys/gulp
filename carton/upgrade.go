package carton

import (
	"fmt"
	"log"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/gulp/upgrade"
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
	err := upgrade.Register("snaps-ci", u.snapCI)
	if err != nil {
		log.Fatalf("unable to register snap-ci: %s", err)
	}
	err = upgrade.Register("snap-bindser", u.snapBind)
	if err != nil {
		log.Fatalf("unable to register snap-bindser: %s", err)
	}
}

// Boot runs the boot of the vm.
func (u *Upgradeable) Upgrade() error {
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := NewLogWriter(u.B)
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	u.w = writer

	err := u.kickOffSerially(writer)
	elapsed := time.Since(start)
	saveErr := saveData(outBuffer.String(), elapsed)
	if saveErr != nil {
		log.Errorf("WARNING: couldn't save ops data, ops opts: %#v", u)
	}

	if err != nil {
		return err
	}
	return nil
}

func (u *Upgradeable) kickOffSerially(writer io.Writer) error {
	return upgrade.Run(upgrade.RunArgs{
		Writer: writer,
		Box:    u.B,
		Force:  true,
	})
}

func (u *Upgradeable) snapBuild(ops *Operation) error {
	fmt.Fprintf(u.w, "---- ci snap (%s) is kicking ----\n", u.B.GetFullName())

	actions := []*action.Action{
		&cloneBox,
		&buildBox,
	}
	pipeline := action.NewPipeline(actions...)
	args := runSnapsPipelineArgs{
		box:    box,
		writer: w,
	}
	if err := pipeline.Execute(args); err != nil {
		return err
	}

	if err := saveSnaps(ps, provision.StatusSnapped); err != nil {
		return err
	}

	if !u.ShouldRestart {
		return nil
	}
	return Provisioner.Restart(u.B, u.w)
}

func (u *Upgradeable) snapBind(ops *Operation) error {
	fmt.Fprintf(u.w, "---- Bind snap (%s) is kicking ----\n", u.B.GetFullName())
	actions := []*action.Action{
		&setEnvsAction,
		&setUpgraded,
	}
	pipeline := action.NewPipeline(actions...)
	args := runSnapsPipelineArgs{
		box:    box,
		writer: w,
	}
	if err := pipeline.Execute(args); err != nil {
		return err
	}

	if err := saveSnaps(ops, provision.StatusSnapped); err != nil {
		return err
	}

	if !u.ShouldRestart {
		return nil
	}
	return Provisioner.Restart(u.B, u.w)
}

func (u *Upgradeable) saveSnaps(ops *Operation, status provison.Status) error {
	if u.B.Level == provision.BoxSome {
		log.Debugf("  save snap status[%s] of box (%s, %s)", m.Id, m.Name, status.String())

		if comp, err := carton.NewComponent(m.Id); err != nil {
			return err
		} else if err = comp.SetDoneOperation(ops, status); err != nil {
			return err
		}
	}
	return nil
}

func (u *Upgradeable) saveData(out string, elapsed time.Duration) error {
	return nil
}
