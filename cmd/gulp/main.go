/* 
** Copyright [2012-2013] [Megam Systems]
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
package main

import (
	"github.com/indykish/gulp/cmd"
	"os"
)

const (
	version = "0.1.0"
	header  = "Supported-Gulp"
)

func buildManager(name string) *cmd.Manager {
	m := cmd.BuildBaseManager(name, version, header)
	m.Register(&gulp.AppRun{})
	m.Register(&gulp.AppRestart{})
	m.Register(&gulp.AppStop{})
	m.Register(&gulp.AppRestart{})
	m.Register(&gulp.AppBuild{})
	m.Register(&gulp.SSLAdd{})
	m.Register(&gulp.SSLRemove{})
	m.Register(&gulp.MeterStop{})
	m.Register(&gulp.MeterStart{})
	m.Register(&gulp.LogStart{})
	m.Register(&gulp.LogStop{})
	m.Register(&gulp.EnvGet{})
	m.Register(&gulp.EnvSet{})
	m.Register(&gulp.EnvUnset{})
	m.Register(&KeyAdd{})
	m.Register(&KeyRemove{})
	m.Register(gulp.ServiceList{})
	return m
}

func main() {
	name := cmd.ExtractProgramName(os.Args[0])
	manager := buildManager(name)
	manager.Run(os.Args[1:])
}
