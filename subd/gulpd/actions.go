/*
** Copyright [2013-2015] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */
package gulpd

import (
	"github.com/megamsys/libgo/action"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/gulp/carton"
)


var publishStatus = action.Action{
	Name: "publish-status",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		args := ctx.Params[0].(*Service)

		err := args.pubStatus(provision.StatusBootstrapped)				
		
		return args, err
	},
	Backward: func(ctx action.BWContext) {
		s := ctx.Params[0].(*Service)
		
		asm, _ := carton.NewAssembly(s.Gulpd.CatID)	
		
		asm.SetStatus(provision.StatusError)
	},
}

var updateStatusInRiak = action.Action{
	Name: "updateStatusInRiak",
	Forward: func(ctx action.FWContext) (action.Result, error) {
		s := ctx.Params[0].(*Service)

		asm, _ := carton.NewAssembly(s.Gulpd.CatID)	
		
		asm.SetStatus(provision.StatusBootstrapped)
		
		return asm, nil
	},
	Backward: func(ctx action.BWContext) {
		s := ctx.Params[0].(*Service)
		
		asm, _ := carton.NewAssembly(s.Gulpd.CatID)	
		
		asm.SetStatus(provision.StatusError)
	},
}


