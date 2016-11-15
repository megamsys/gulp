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

package carton

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	lb "github.com/megamsys/gulp/logbox"
	"github.com/megamsys/gulp/meta"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/gulp/repository"
)

const BuildFile = "build"

var (
	ErrBuildableNotFound = errors.New("no build scripts found. (searched for platform/custom scripts)")
)

type RepoBuilder struct {
	R      repository.Repository
	BP     BashBuild
	writer io.Writer
}

type BashBuild struct {
	platform string
	custom   string
}

func (rb *RepoBuilder) custom() (string, error) {
	repoName, err := rb.R.GetShortName()
	if err != nil {
		return "", err
	}
	custompath := filepath.Join(meta.MC.Dir, repoName, BuildFile)
	if _, err = os.Stat(custompath); err == nil {
		return custompath, nil
	}
	return "", nil
}

func NewRepoBuilder(r repository.Repository, w io.Writer) *RepoBuilder {
	return &RepoBuilder{
		R:      r,
		writer: w,
		BP: BashBuild{
			platform: filepath.Join(meta.MC.Dir, BuildFile),
		},
	}
}

func (rb *RepoBuilder) Build(force bool) error {
	fmt.Fprintf(rb.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  %s (%s)\n", BuildFile, rb.R.Gitr())))
	custom, err := rb.custom()

	if err != nil {
		return err
	}
	rb.BP.custom = custom

	if _, err = os.Stat(rb.BP.custom); err == nil {
		return provision.ExecuteCommandOnce(strings.Fields(rb.BP.custom), rb.writer)
	}
	if _, err = os.Stat(rb.BP.platform); err == nil {
		return provision.ExecuteCommandOnce(strings.Fields(rb.BP.platform), rb.writer)
	}
	fmt.Fprintf(rb.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  %s (%s) failed\n", BuildFile, rb.R.Gitr())))
	return ErrBuildableNotFound
}

func (s *RepoBuilder) Cleanup(suffix string) error {
	return nil
}
