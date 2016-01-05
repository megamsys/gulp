package chefsolo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/megamsys/gulp/meta"
	"github.com/megamsys/gulp/repository"
	_ "github.com/megamsys/gulp/repository/github"
	"github.com/mitchellh/ioprogress"
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

func (ch *ChefRepo) Download() error {
	if !ch.exists() || !ch.isUptodate() {
		if err := ch.download(); err != nil {
			return scm().Clone(repository.Repo{URL: ch.git})
		}
	}
	return nil
}

func (ch *ChefRepo) Torr() error {
	if !ch.exists() {
		return NewTorr(ch.tarfile()).untar()
	}
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
	return ch.dir + "/" + f
}

func (ch *ChefRepo) tarfile() string {
	tokens := strings.Split(ch.tar, "/")
	return ch.dir + "/" + tokens[len(tokens)-1]
}

func scm() repository.RepositoryManager {
	return repository.Manager("github")
}

func (ch *ChefRepo) exists() bool {
	var exists = false
	if f := ch.repodir(); len(strings.TrimSpace(f)) > 0 {
		if _, err := os.Stat(ch.repodir()); err != nil {
			exists = true
		}
	}
	return exists
}

func (ch *ChefRepo) isUptodate() bool {
	return true
}

func (ch *ChefRepo) download() error {
	output, err := os.Create(ch.tarfile())
	if err != nil {
		fmt.Println("Error while creating", ch.tarfile(), "-", err)
		return err
	}
	defer output.Close()

	response, err := http.Get(ch.tar)
	if err != nil {
		fmt.Println("Error while downloading", ch.tar, "-", err)
		return err
	}
	defer response.Body.Close()

	// Create the progress reader
	progressR := &ioprogress.Reader{
		Reader: response.Body,
		Size:   response.ContentLength,
	}

	n, err := io.Copy(output, progressR)
	if err != nil {
		fmt.Println("Error while downloading", ch.tar, "-", err)
		return err
	}

	fmt.Println(n, "bytes downloaded.")
	return nil
}
