package carton

import (
	"errors"
	"fmt"
	"github.com/megamsys/gulp/carton/bind"
	lb "github.com/megamsys/gulp/logbox"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/libgo/action"
	"io"
	"strings"
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
		fmt.Fprintf(args.writer, lb.W(lb.VM_UPGRADING, lb.INFO, fmt.Sprintf("  set envs for box (%s) \n", args.box.GetFullName())))
		bi := &bind.BindFile{}

		if len(args.box.Envs) > 0 {
			bi.Name = "env.sh"
			bi.BoxName = args.box.GetFullName()
			bi.LogWriter = args.writer

			if err := bi.Mutate(strings.NewReader(args.box.Envs.WrapForInitds())); err != nil {
				return bi, err
			}
		}
		fmt.Fprintf(args.writer, lb.W(lb.VM_UPGRADING, lb.INFO, fmt.Sprintf("  set envs for box (%s) OK\n", args.box.GetFullName())))
		return bi, nil
	},
	Backward: func(ctx action.BWContext) {
		_, _ = ctx.Params[0].(*runOpsPipelineArgs)
		c := ctx.FWResult.(bind.BindFile)
		bind.Revert(&c)
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
		fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  clone repository for box (%s)\n", args.box.GetFullName())))
		if err := args.box.Clone(); err != nil {
			fmt.Fprintf(args.writer, lb.W(lb.VM_DEPLOY, lb.ERROR, fmt.Sprintf("  clone repository for box (%s) failed \n%s", args.box.GetFullName(), err.Error())))
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
		fmt.Fprintf(args.writer, lb.W(lb.VM_UPGRADING, lb.INFO, fmt.Sprintf("  build repository for box (%s)\n", args.box.GetFullName())))

		if err := NewRepoBuilder(args.box.Repo, args.writer).Build(false); err != nil {
			return nil, err
		}
		fmt.Fprintf(args.writer, lb.W(lb.VM_UPGRADING, lb.INFO, fmt.Sprintf("  build repository for box (%s) OK\n", args.box.GetFullName())))
		return nil, nil
	},
	Backward: func(ctx action.BWContext) {
		_, _ = ctx.Params[0].(*runOpsPipelineArgs)
	},
	MinParams: 1,
}
