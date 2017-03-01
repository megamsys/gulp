/*
** Copyright [2013-2016] [Megam Systems]
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
package bind

import (
	"bufio"
	"fmt"
	lb "github.com/megamsys/gulp/logbox"
	"github.com/megamsys/gulp/meta"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// BindFile represents a file with  environment variable
type BindFile struct {
	Name      string
	BoxName   string
	LogWriter io.Writer
}

func (bi *BindFile) exists() bool {
	return false
}

func (bi *BindFile) envPath(name string) string {
	return filepath.Join(meta.MC.Home, name)
}

func (bi *BindFile) envSH() string {
	return bi.envPath(bi.Name)
}

func (bi *BindFile) envSHWriter() (*os.File, error) {
	envFile, err := os.Create(bi.envSH())
	if err != nil {
		return nil, err
	}
	return envFile, nil
}

func (bi *BindFile) envSHString() (string, error) {
	if _, err := os.Stat(bi.envSH()); err == nil {
		rawBytes, err := ioutil.ReadFile(bi.envSH())
		if err != nil {
			return "", err
		}
		return string(rawBytes), nil
	}
	return "", nil
}

func (bi *BindFile) envSHdot() string {
	return bi.envPath(bi.Name + ".save")
}

func (bi *BindFile) Mutate(newEnvReader io.Reader) (err error) {
	ei, err := bi.envSHString()
	if err != nil {
		return err
	}

	backup(bi)

	ew, err := bi.envSHWriter()
	if err != nil {
		return err
	}

	reader := bufio.NewReader(strings.NewReader(ei))
	writer := bufio.NewWriter(ew)

	defer func() {
		if err == nil {
			err = writer.Flush()
		}
	}()
	fmt.Fprintf(bi.LogWriter, lb.W(lb.VM_UPGRADING, lb.INFO, fmt.Sprintf("  set envs replacing (%s)\n", bi.envSH())))

	var replacer func(string) string
	if replacer, err = mkReplaceFunction(newEnvReader); err != nil {
		return err
	}

	envRx := regexp.MustCompile(`([a-zA-Z0-9_//-]+)=([a-zA-Z0-9_//-]+)`)

	eof := false
	for !eof {
		var line string
		line, err = reader.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			return err
		}
		line = envRx.ReplaceAllStringFunc(line, replacer)
		if _, err = writer.WriteString(line); err != nil {
			return err
		}
	}

	cleanup(bi)
	return nil
}

func mkReplaceFunction(newEnvsReader io.Reader) (func(string) string, error) {
	reader := bufio.NewReader(newEnvsReader)

	eof := false
	envsNewInput := make(map[string]string)
	for !eof {
		var line string
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			return nil, err
		}
		fields := strings.Fields(line)
		sfields := strings.Split(fields[len(fields)-1], "=")

		if len(sfields) == 2 {
			envsNewInput[sfields[0]] = sfields[1]
		}
	}

	return func(word string) string {
		rew := strings.Split(word, "=")
		if len(rew) >= 2 {
			if envWord, found := envsNewInput[rew[0]]; found {
				return rew[0] + "=" + envWord
			}
		}
		return word
	}, nil

}

//remove env.sh.save, and rename env.sh to env.sh.save
func backup(bi *BindFile) {
	fmt.Fprintf(bi.LogWriter, lb.W(lb.VM_UPGRADING, lb.INFO, fmt.Sprintf("  set envs backing (%s)\n", bi.envSHdot())))

	if _, err := os.Stat(bi.envSHdot()); err == nil {
		if err = os.Remove(bi.envSHdot()); err != nil {
			return
		}
	}

	if _, err := os.Stat(bi.envSH()); err == nil {
		if err = os.Rename(bi.envSH(), bi.envSHdot()); err != nil {
			return
		}
	}
	return
}

//remove env.sh.save file
func cleanup(bi *BindFile) {
	fmt.Fprintf(bi.LogWriter, lb.W(lb.VM_UPGRADING, lb.INFO, fmt.Sprintf("  set envs cleaning (%s)\n", bi.envSHdot())))
	if _, err := os.Stat(bi.envSHdot()); err == nil {
		if err = os.Remove(bi.envSHdot()); err != nil {
			return
		}
	}
}

//remove env.sh, and rename env.sh.save to env.sh
func Revert(bi *BindFile) {
	fmt.Fprintf(bi.LogWriter, lb.W(lb.VM_UPGRADING, lb.INFO, fmt.Sprintf("  set envs reverting (%s)\n", bi.envSHdot())))
	if _, err := os.Stat(bi.envSH()); err == nil {
		if err = os.Remove(bi.envSH()); err != nil {
			return
		}
	}

	if _, err := os.Stat(bi.envSHdot()); err == nil {
		if err = os.Rename(bi.envSHdot(), bi.envSH()); err != nil {
			return
		}
	}
	return
}
