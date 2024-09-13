package git

import (
	"fmt"
	"github.com/revelaction/ical-git/fetch"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
)

// FileSystem implements the fetch.Fetcher interface
type Git struct {
	url            string
	privateKeyPath string
	ch             chan fetch.File
}

// New creates a new instance of FileSystem
func New(url string, privateKeyPath string) fetch.Fetcher {
	return &Git{
		url:            url,
		privateKeyPath: privateKeyPath,
		ch:             make(chan fetch.File),
	}
}

// GetCh implements the fetch.Fetcher interface method
func (g *Git) GetCh() <-chan fetch.File {
	go func() {
		defer close(g.ch)

		// Create SSH auth method
		auth, err := ssh.NewPublicKeysFromFile("git", g.privateKeyPath, "")
		if err != nil {
			g.ch <- fetch.File{Error: fmt.Errorf("failed to create SSH auth method: %w", err)}
			return
		}

		// Create in-memory filesystem
		fs := memfs.New()

		// Clone the repository
		repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
			URL:           g.url,
			Auth:          auth,
			ReferenceName: plumbing.Master,
			SingleBranch:  true,
			Depth:         1,
		})
		if err != nil {
			g.ch <- fetch.File{Error: fmt.Errorf("failed to clone repository: %w", err)}
			return
		}

		// Get the latest commit on master
		ref, err := repo.Head()
		if err != nil {
			g.ch <- fetch.File{Error: fmt.Errorf("failed to get HEAD: %w", err)}
			return
		}

		commit, err := repo.CommitObject(ref.Hash())
		if err != nil {
			g.ch <- fetch.File{Error: fmt.Errorf("failed to get commit: %w", err)}
			return
		}

		// Get the tree for the commit
		tree, err := commit.Tree()
		if err != nil {
			g.ch <- fetch.File{Error: fmt.Errorf("failed to get tree: %w", err)}
			return
		}

		err = tree.Files().ForEach(func(f *object.File) error {

			contents, err := f.Contents()
			if err != nil {
				// TODO
				return nil
			}

			g.ch <- fetch.File{Path: f.Name, Content: []byte(contents)}
			return nil
		})

		// TODO
		if err != nil {
			fmt.Printf("error walking through repository: %v\n", err)
		}
	}()

	return g.ch
}
