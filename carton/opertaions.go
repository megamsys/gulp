package carton

import (
	"bytes"
	"io"
	"time"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/libgo/cmd"
)

type ResetOpts struct {
	B     *provision.Box
}

// Deploy runs a deployment of an application.
func ResetPassword(opts *ResetOpts) error {
	var outBuffer bytes.Buffer
	start := time.Now()

	logWriter := NewLogWriter(opts.B)
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	err := updateNewPassword(opts, writer)
	elapsed := time.Since(start)
  if err != nil {
		return err
	}
  log.Debugf("%s in (%s)\n%s",
		cmd.Colorfy(opts.B.GetFullName(), "cyan", "", "bold"),
		cmd.Colorfy(elapsed.String(), "green", "", "bold"),
		cmd.Colorfy(outBuffer.String(), "yellow", "", ""))
	return nil
}

func updateNewPassword(opts *ResetOpts, writer io.Writer) error {
   	if deployer, ok := Provisioner.(provision.Deployer); ok {
   		deployer.ResetPassword(opts.B, writer)
   	}
  return nil
}
