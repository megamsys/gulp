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
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/megamsys/libgo/cmd"
	"github.com/tsuru/config"
)

const (
	version = "0.3.0"
	header  = "Supported-Gulp"
)

const defaultConfigPath = "conf/gulpd.conf"

//const defaultConfigPath = "/home/megam/bin/conf/gulpd.conf"

func buildManager(name string) *cmd.Manager {
	m := cmd.BuildBaseManager(name, version, header)
	m.Register(&GulpcStart{m, nil, false}) //start the gulpc server
	m.Register(&GulpcUpdate{}) //stop  the gulpc server	
	return m
}

func main() {
	p, _ := filepath.Abs(defaultConfigPath)
	log.Println(fmt.Errorf("Conf: %s", p))
	config.ReadConfigFile(defaultConfigPath)
	name := cmd.ExtractProgramName(os.Args[0])
	manager := buildManager(name)
	manager.Run(os.Args[1:])
}
