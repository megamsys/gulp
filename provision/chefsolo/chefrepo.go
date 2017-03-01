package chefsolo

import (
	"fmt"
	lb "github.com/megamsys/gulp/logbox"
	"github.com/megamsys/gulp/meta"
	"github.com/megamsys/gulp/provision"
	"github.com/megamsys/gulp/repository"
	_ "github.com/megamsys/gulp/repository/github"
	constants "github.com/megamsys/libgo/utils"
	"github.com/mitchellh/ioprogress"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ChefRepo struct {
	git     string
	tar     string
	dir     string
	version string
	writer  io.Writer
}

func NewChefRepo(m map[string]string, w io.Writer) *ChefRepo {
	return &ChefRepo{
		git:    m[CHEFREPO_GIT],
		tar:    m[CHEFREPO_TARBALL],
		dir:    meta.MC.Dir,
		writer: w,
	}
}

//try downloading tar first, if not, do a clone of the chef-repo
func (ch *ChefRepo) Download(force bool) error {
	_ = provision.EventNotify(constants.StatusCookbookDownloading)
	fmt.Fprintf(ch.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("--- download (%s)\n", ch.repodir())))
	if !ch.exists() || !ch.isUptodate() {
		if err := ch.download(force); err != nil {
			return scm().Clone(repository.Repo{URL: ch.git})
		}
	}
	fmt.Fprintf(ch.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("--- download (%s)OK\n", ch.repodir())))
	return nil
}

func (ch *ChefRepo) Torr() error {
	fmt.Fprintf(ch.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("--- torr (%s)\n", ch.tarfile())))
	if !ch.exists() {
		tr := NewTorr(ch.tarfile())
		tr.Base = ch.repodir()
		tr.writer = ch.writer
		if err := tr.untar(); err != nil {
			return err
		}
		return tr.cleanup()
	}
	fmt.Fprintf(ch.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("--- torr (%s) OK\n", ch.tarfile())))
	return nil
}

func (ch *ChefRepo) filename() (string, error) {
	return (repository.Repo{URL: ch.git}).GetShortName()
}

func (ch *ChefRepo) repodir() string {
	f, err := ch.filename()
	if err != nil {
		return ""
	}
	return filepath.Join(ch.dir, f)
}

func (ch *ChefRepo) tarfile() string {
	tokens := strings.Split(ch.tar, "/")
	return filepath.Join(ch.dir, tokens[len(tokens)-1])
}

//bit screwy we are doing it twice inside here and in provisioner
func scm() repository.RepositoryManager {
	return repository.Manager("github")
}

func (ch *ChefRepo) exists() bool {
	var exists = false
	if f := ch.repodir(); len(strings.TrimSpace(f)) > 0 {
		if _, err := os.Stat(ch.repodir()); err == nil {
			exists = true
		}
	}
	return exists
}

//for now its always uptodate
func (ch *ChefRepo) isUptodate() bool {
	return true
}

func (ch *ChefRepo) download(force bool) error {
	if force {
		_ = os.RemoveAll(ch.tarfile())
	}
	fmt.Fprintf(ch.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  create tar (%s)\n", ch.tarfile())))
	output, err := os.Create(ch.tarfile())
	if err != nil {
		return err
	}
	defer output.Close()
	response, err := http.Get(ch.tar)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	fmt.Fprintf(ch.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  http GET tar (%s) \n", ch.tar)))
	// Create the progress reader
	progressR := &ioprogress.Reader{
		Reader: response.Body,
		Size:   response.ContentLength,
	}

	_, err = io.Copy(output, progressR)
	if err != nil {
		return err
	}
	fmt.Fprintf(ch.writer, lb.W(lb.VM_DEPLOY, lb.INFO, fmt.Sprintf("  http GET, write tar (%s) OK\n", ch.tar)))
	return nil
}
