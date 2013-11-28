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
	"github.com/globocom/config"
	"github.com/indykish/gulp/cmd"
	"github.com/indykish/gulp/cmd/gulpd_jump"
	"os"
)

const (
	version = "0.1.0"
	header  = "Supported-Gulp"
)

const defaultConfigPath = "/conf/gulpd.conf"

/*
type tsrCommand struct {
	cmd.Command
	fs   *gnuflag.FlagSet
	file configFile
}
	m.Register(&tsrCommand{Command: &apiCmd{}})

*/

func buildManager(name string) *cmd.Manager {
	m := cmd.BuildBaseManager(name, version, header)
	gs:= &jump.GulpStart{}
	m.Register(gs) //start the gulpd daemon
	/*m.Register(Stop{})		   //stop  the gulpd daemon
	m.Register(&gulp.AppStart{})   //sudo service <appname> start
	m.Register(&gulp.AppStop{})    //sudo service <appname> stop
	m.Register(&gulp.AppRestart{}) //sudo service <apppname> restart
	m.Register(&gulp.AppBuild{})   //git fetch -q
	m.Register(&gulp.AppMaintain{})//sudo service nginx maintain ?
	m.Register(&gulp.SSLAdd{})     //download node_name.pub, crt from S3, mk ssl_template, cp to sites_available, ln to sites_enabled. && AppRestart
	m.Register(&gulp.SSLRemove{})  //rm node_name.pub, crt, mk regular non_ssl_template, cp to sites_available, ln to sites_enabled. && AppRestart
	m.Register(&gulp.MeterStop{})  //sudo service gmond start
	m.Register(&gulp.MeterStart{}) //sudo service gmond stop
	m.Register(&gulp.LogStart{})   //sudo service beaver start
	m.Register(&gulp.LogStop{})    //sudo service beaver stop
	m.Register(&gulp.EnvGet{})     //ENV['JMP_UP_PATH'] '~/sofware/kangaroo'
	m.Register(&gulp.EnvSet{})     //ENV['JMP_UP_PATH'] = '~/software/kangaroo'
	m.Register(&gulp.EnvUnset{})   //ENV['JMP_UP_PATH'] = blank
	m.Register(&KeyAdd{})          //add the id_rsa/pub
	m.Register(&KeyRemove{})       //remove the id_rsa/pub
	m.Register(gulp.ServiceList{}) //ps -ef
	*/
	return m
}

func main() {
	config.ReadConfigFile(defaultConfigPath)
	name := cmd.ExtractProgramName(os.Args[0])
	manager := buildManager(name)
	manager.Run(os.Args[1:])
}
