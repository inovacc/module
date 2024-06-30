package git

import (
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"io"
	"net/url"
)

type Git struct {
	*git.Repository
}

type Repo struct {
	g *git.Repository
}

func NewGitNewStorage(repoUrl string) (*Git, error) {
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: repoUrl,
	})
	if err != nil {
		return nil, err
	}
	return &Git{Repository: r}, nil
}

func NewGitExistingStorage(repoUrl string) (*Git, error) {
	r, err := git.PlainOpenWithOptions(repoUrl, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, err
	}
	return &Git{Repository: r}, nil
}

func NewGitBareStorage() (*Git, error) {
	var wt billy.Filesystem
	r, err := git.InitWithOptions(memory.NewStorage(), wt, git.InitOptions{})
	if err != nil {
		return nil, err
	}
	return &Git{Repository: r}, nil
}

//dir = C:/Users/ddaniels/AppData/Local/Temp/modload-test-1906163017/pkg/mod/cache/vcs/{7d9b3b49b55db5b40e68a94007f21a05905d3fda866f685220de88f9c9bad98a} <- hash url
//hashId = cd70d50baa6daa949efa12e295e10829f3a7bd46
//repoUrl = https://go.googlesource.com/tools
//git init --bare
//git remote add origin -- {repoUrl}
//git config core.longpaths true
//git ls-remote -q origin > list hashes
//git tag -l
//git -c log.showsignature=false log --no-decorate -n1 --format=format:%H %ct %D {hashId} --
//git -c protocol.version=2 fetch -f --depth=1 origin refs/tags/{repoName}/{tags}:refs/tags/{repoName}/{tags}
//git -c log.showsignature=false log --no-decorate -n1 --format=format:%H %ct %D refs/tags/{repoName}/{tags} --
//git cat-file blob {hashId}:{repoName}/{fileName}
//
//
//
//
//
//git -c log.showsignature=false log --no-decorate -n1 --format=format:%H %ct %D cd70d50baa6daa949efa12e295e10829f3a7bd46 --

func (g *Git) Repo() *Repo {
	return &Repo{g: g.Repository}
}

func (r *Repo) SetRemote(remoteUrl string) error {
	u, err := url.Parse(remoteUrl)
	if err != nil {
		return err
	}
	u.Scheme = "https"
	_, err = r.g.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{u.String()},
	})
	return err
}

func (r *Repo) Fetch() error {
	return r.g.Fetch(&git.FetchOptions{
		RemoteName: "origin",
	})
}

func (r *Repo) GetFileContent(hashId, fileName string) (io.ReadCloser, error) {
	file, err := r.g.BlobObject(plumbing.NewHash(fmt.Sprintf("%s:%s", hashId, fileName)))
	if err != nil {
		return nil, err
	}
	return file.Reader()
}

func (r *Repo) GetCommitInfo(hashId string) (string, error) {
	commit, err := r.g.CommitObject(plumbing.NewHash(hashId))
	if err != nil {
		return "", err
	}
	return commit.Message, nil
}

func (r *Repo) GetTags() ([]string, error) {
	tags := make([]string, 0)
	iter, err := r.g.Tags()
	if err != nil {
		return nil, err
	}
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		tags = append(tags, ref.Name().Short())
		return nil
	})
	return tags, err
}

func (r *Repo) GetHashId() (string, error) {
	ref, err := r.g.Head()
	if err != nil {
		return "", err
	}
	return ref.Hash().String(), nil
}
