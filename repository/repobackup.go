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

package repository

import (
	"os"
	"path/filepath"
)

type RepoBackup struct {
	f string
	t string
}

func (s *RepoBackup) fromPath(suffix string) string {
	return filepath.Join(s.f, suffix)
}

func (s *RepoBackup) toPath(suffix string) string {
	return filepath.Join(s.t, suffix+"_old")
}

func NewRepoBackup(from string, to string) *RepoBackup {
	return &RepoBackup{
		f: from,
		t: to,
	}
}

func (s *RepoBackup) Backup(suffix string) error {
	if _, err := os.Stat(s.fromPath(suffix)); err == nil {
		if err = os.Rename(s.fromPath(suffix), s.toPath(suffix)); err != nil {
			return err
		}
		if err = os.RemoveAll(s.fromPath(suffix)); err != nil {
			return err
		}
	}
	return nil
}

func (s *RepoBackup) Revert(suffix string) error {
	if err := s.removeall(s.fromPath(suffix)); err != nil {
		return err
	}

	return s.rename(s.toPath(suffix), s.fromPath(suffix))
}

func (s *RepoBackup) Cleanup(suffix string) error {
	return s.removeall(s.toPath(suffix))
}

func (s *RepoBackup) removeall(path string) error {
	if _, err := os.Stat(path); err == nil {
		if err = os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}

func (s *RepoBackup) rename(from string, to string) error {
	if _, err := os.Stat(from); err == nil {
		if err = os.Rename(from, to); err != nil {
			return err
		}
	}
	return nil
}
