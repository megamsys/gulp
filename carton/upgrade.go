package carton

import (
	"bytes"
	"fmt"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/gulp/upgrade"
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/libgo/cmd"
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

func (u *Upgradeable) canCycle() bool {
	return u.B.Status == provision.StatusRunning ||
		u.B.Status == provision.StatusStarted ||
		u.B.Status == provision.StatusStopped ||
		u.B.Status == provision.StatusUpgraded
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
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := NewLogWriter(u.B)
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	if !u.canCycle() {
		fmt.Fprintf(writer, "  skip upgrade for box (%s)\n", u.B.GetFullName())
		return nil
	}
	err := u.operateBox(writer)
	elapsed := time.Since(start)
	saveErr := saveUpgradeData(u, outBuffer.String(), elapsed)

	if saveErr != nil {
		log.Errorf("WARNING: couldn't save upgrade data, deploy opts: %#v", u)
	}
	if err != nil {
		return err
	}
	return nil
}

func (u *Upgradeable) operateBox(writer io.Writer) error {
	u.w = writer
	fmt.Fprintf(u.w, "---- operate box (%s)\n", u.B.GetFullName())

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
	fmt.Fprintf(u.w, "---- operate box (%s) OK\n", u.B.GetFullName())
	return Provisioner.Restart(u.B, u.w)
}

func (u *Upgradeable) opsBuild() error {
	fmt.Fprintf(u.w, "  ops ci (%s) is kicking\n", u.B.GetFullName())

	actions := []*action.Action{
		&cloneBox,
		&buildBox, //buildpack does everthing
	}
	pipeline := action.NewPipeline(actions...)
	args := runOpsPipelineArgs{
		box:    u.B,
		writer: u.w,
	}
	if err := pipeline.Execute(&args); err != nil {
		return err
	}
	fmt.Fprintf(u.w, "  ops ci (%s) OK\n", u.B.GetFullName())
	return nil
}

func (u *Upgradeable) opsBind() error {
	fmt.Fprintf(u.w, "  ops bind (%s) is kicking\n", u.B.GetFullName())

	actions := []*action.Action{
		&setEnvsAction,
	}
	pipeline := action.NewPipeline(actions...)
	args := runOpsPipelineArgs{
		box:    u.B,
		writer: u.w,
	}
	if err := pipeline.Execute(&args); err != nil {
		return err
	}
	fmt.Fprintf(u.w, "  ops bind (%s) OK\n", u.B.GetFullName())
	return nil
}

func (u *Upgradeable) saveData(opsRan upgrade.OperationsRan, elapsed time.Duration) error {
	if u.B.Level == provision.BoxSome {
		fmt.Fprintf(u.w, "  operate box saving.. (%s)\n", u.B.GetFullName())
		if comp, err := NewComponent(u.B.Id); err != nil {
			return err
		} else if err = comp.UpdateOpsRun(opsRan); err != nil {
			return err
		}
	}
	fmt.Fprintf(u.w, "  operate box saving (%s) OK\n", u.B.GetFullName())
	return nil
}

func saveUpgradeData(opts *Upgradeable, ulog string, duration time.Duration) error {
	log.Debugf("%s in (%s)\n%s",
		cmd.Colorfy(opts.B.GetFullName(), "cyan", "", "bold"),
		cmd.Colorfy(duration.String(), "green", "", "bold"),
		cmd.Colorfy(ulog, "yellow", "", ""))
	return nil
}
