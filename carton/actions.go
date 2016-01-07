package carton

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/megamsys/gulp/meta"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/libgo/action"
)

type runOpsPipelineArgs struct {
	box    *provision.Box
	writer io.Writer
}

var setEnvsAction = action.Action{
	Name: "set-envs",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args, _ := ctx.Params[0].(*runOpsPipelineArgs)
		if args == nil {
			return nil, errors.New("invalid arguments for pipeline, expected *runOpsPipelineArgs")
		}

		fmt.Fprintf(args.writer, "  set envs for box (%s)\n", args.box.GetFullName())

		if len(args.box.Envs) > 0 {
			envFile := filepath.Join(meta.MC.Home, "env.sh")
			if _, err := os.Stat(envFile); err == nil {
				fmt.Fprintf(args.writer, "  set envs for box (%s) appending ...\n", args.box.GetFullName())
				aenvFile, err := os.OpenFile(envFile, os.O_APPEND|os.O_WRONLY, 0755)
				if err != nil {
					return err, nil
				}
				io.WriteString(aenvFile, args.box.Envs.WrapForInitds())
				aenvFile.Close()
			}
		}
		fmt.Fprintf(args.writer, "  set envs for box (%s) OK\n", args.box.GetFullName())
		return nil, nil
	},
	Backward: func(ctx action.BWContext) {
		_, _ = ctx.Params[0].(*runOpsPipelineArgs)
	},
	MinParams: 1,
}

var cloneBox = action.Action{
	Name: "clone-box",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args, _ := ctx.Params[0].(*runOpsPipelineArgs)
		if args == nil {
			return nil, errors.New("invalid arguments for pipeline, expected *runOpsPipelineArgs")
		}
		fmt.Fprintf(args.writer, "  clone repository for box (%s)\n", args.box.GetFullName())
		if err := args.box.Clone(); err != nil {
			fmt.Fprintf(args.writer, "  clone repository for box (%s) failed\n%s", args.box.GetFullName(), err.Error())
			return nil, err
		}
		return nil, nil
	},
	Backward: func(ctx action.BWContext) {
		//nothing to backup as cleanup is handled in Clone()
	},
}

var buildBox = action.Action{
	Name: "build-box",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args, _ := ctx.Params[0].(*runOpsPipelineArgs)
		if args == nil {
			return nil, errors.New("invalid arguments for pipeline, expected *runOpsPipelineArgs")
		}
		fmt.Fprintf(args.writer, "  build repository for box (%s)\n", args.box.GetFullName())

		if err := NewRepoBuilder(args.box.Repo, args.writer).Build(false); err != nil {
			return nil, err
		}

		fmt.Fprintf(args.writer, "  build repository for box (%s) OK\n", args.box.GetFullName())
		return nil, nil
	},
	Backward: func(ctx action.BWContext) {
		_, _ = ctx.Params[0].(*runOpsPipelineArgs)
	},
	MinParams: 1,
}
