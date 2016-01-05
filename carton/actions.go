package carton

import (
	"errors"
	"io"
	"net/http"
	"sync"
)

type snapsPipelineArgs struct {
	box    provision.Box
	writer io.Writer
}

var setEnvsAction = action.Action{
	Name: "set-envs",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args, _ := ctx.Params[0].(*snapsPipelineArgs)
		if args == nil {
			return nil, errors.New("invalid arguments for pipeline, expected *snapsPipelineArgs")
		}

		if len(args.box.Envs) > 0 {
			envFile := meta.MC.HomeDir + "/" + "env.sh"
			if _, err := os.Stat(filename); err == nil {
				file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0755)
				if err != nil {
					return err, nil
				}
				io.WriteString(file, args.box.Envs.ForInitService())
				file.Close()
			}
		}
		return nil, nil
	},
	Backward: func(ctx action.BWContext) {
		_, _ := ctx.Params[0].(*snapsPipelineArgs)
	},
	MinParams: 1,
}

var cloneBox = action.Action{
	Name: "clone-box",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args, _ := ctx.Params[0].(*snapsPipelineArgs)
		if args == nil {
			return nil, errors.New("invalid arguments for pipeline, expected *snapsPipelineArgs")
		}
		fmt.Fprintf(args.writer, "  clone repository for box (%s)", args.box.GetFullName())
		if err := args.box.Clone(); err != nil {
			return nil, err
		}
		return nil, nil
	},
	Backward: func(ctx action.BWContext) {
		//delete the repository directory
	},
}

var buildBox = action.Action{
	Name: "build-box",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args, _ := ctx.Params[0].(*snapsPipelineArgs)
		if args == nil {
			return nil, stderrors.New("invalid arguments for pipeline, expected *snapsPipelineArgs")
		}

		return nil, nil
	},
	Backward: func(ctx action.BWContext) {
		_, _ := ctx.Params[0].(*snapsPipelineArgs)
	},
	MinParams: 1,
}
