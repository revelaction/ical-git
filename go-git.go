package main

import (
	"fmt"
	"io"
	"log"
	"os"

	//"io"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"golang.org/x/crypto/ssh"
)

func main() {

	//// Create an authentication method
	//authMethod, err := ssh.DefaultAuthBuilder("github-revelaction")
	//if err != nil {
	//	log.Fatalf("default auth builder: %v", err)
	//}

	url := "git@github.com:revelaction/privage.git"
	//url := "git@github.com:revelaction/ical-git.git"

	sshPassphrase := os.Getenv("SSH_PASSPHRASE")
	if sshPassphrase == "" {
		log.Fatalf("SSH_PASSPHRASE environment variable is not set")
	}

	s := fmt.Sprintf("%s/.ssh/id_revelaction", os.Getenv("HOME"))
	sshKey, err := os.ReadFile(s)
	if err != nil {
		log.Fatalf("plain clone: %v", err)
	}
	signer, err := ssh.ParsePrivateKeyWithPassphrase([]byte(sshKey), []byte(sshPassphrase))
	if err != nil {
		log.Fatalf("plain clone: %v", err)
	}
	auth := &gitssh.PublicKeys{User: "git", Signer: signer}

	//var publicKey *ssh.PublicKeys
	//sshPath := os.Getenv("HOME") + "/.ssh/id_revelaction"
	//sshKey, _ := os.ReadFile(sshPath)
	//publicKey, keyError := ssh.NewPublicKeys("git", []byte(sshKey), "")
	//if keyError != nil {
	//	fmt.Println(keyError)
	//}
	fs := memfs.New()
	repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL:  url,
		Auth: auth,
		//Progress: io.Discard, // Suppress output
	})

	if err != nil {
		log.Fatalf("plain clone: %v", err)
	}

	// 1) WORKS

	ref, err := repo.Reference(plumbing.HEAD, true)
	if err != nil {
		log.Fatalf("reference: %v", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		log.Fatalf("commit object: %v", err)
	}

	tree, err := commit.Tree()
	if err != nil {
		log.Fatalf("tree: %v", err)
	}

	// List the files in the root tree
	err = tree.Files().ForEach(func(file *object.File) error {
		fmt.Println(file.Name)

		f, err := fs.Open(file.Name)
		if err != nil {
			log.Fatalf("Error when opening file: %s", err)
		}
		defer f.Close()

		// Read the first line of the file
		fileContent, err := io.ReadAll(f)
		if err != nil {
			log.Fatalf("Error while reading file: %s", err)
		}

		fmt.Println(string(fileContent))
		return nil

	})
	if err != nil {
		log.Fatalf("for each: %v", err)
	}

	// 2) DOES NOT WORK
	//ref, err := repo.Reference(plumbing.HEAD, true)
	//if err != nil {
	//    log.Fatalf("reference: %v", err)
	//}
	//fmt.Printf("HEAD: %v hash: %s \n", ref, ref.Hash())
	//tree, err := repo.TreeObject(ref.Hash())
	//if err != nil {
	//    log.Fatalf("tree object: %v", err)
	//}

	//// List the files in the root tree
	//err = tree.Files().ForEach(func(file *object.File) error {
	//    fmt.Println(file.Name)
	//    return nil
	//})
	//if err != nil {
	//    log.Fatalf("for each: %v", err)
	//}
}
