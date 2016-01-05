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

package github

import (
	"github.com/megamsys/gulp/repository"
	skia "go.skia.org/infra/go/gitinfo"
)

func init() {
	repository.Register("github", gitManager{})
}

type gitHubManager struct{}

func (m gitHubManager) Clone(r repository.Repository) error {
	sn, err := r.GetShortName()

	if err !=nil {
		return err
	}

	_, err = skia.Clone(r.Gitr(), sn, false)
	return err

}
