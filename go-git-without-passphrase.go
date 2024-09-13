package main

import (
    "fmt"

    "github.com/go-git/go-git/v5"
    "github.com/go-git/go-git/v5/plumbing"
    "github.com/go-git/go-git/v5/plumbing/object"
    "github.com/go-git/go-git/v5/plumbing/transport/ssh"
    "github.com/go-git/go-git/v5/storage/memory"
    "github.com/go-git/go-billy/v5/memfs"
)

func main() {
    // Repository URL
    repoURL := "git@github.com:revelaction....."

    // Path to your SSH private key
    privateKeyPath := "/home/revelac...."

    // Create SSH auth method
    auth, err := ssh.NewPublicKeysFromFile("git", privateKeyPath, "")
    if err != nil {
        fmt.Printf("Failed to create SSH auth method: %v\n", err)
        return
    }

    // Create in-memory filesystem
    fs := memfs.New()

    // Clone the repository
    fmt.Println("Cloning repository...")
    repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
        URL:           repoURL,
        Auth:          auth,
        ReferenceName: plumbing.Master,
        SingleBranch:  true,
        Depth:         1,
    })
    if err != nil {
        fmt.Printf("Failed to clone repository: %v\n", err)
        return
    }

    // Get the latest commit on master
    ref, err := repo.Head()
    if err != nil {
        fmt.Printf("Failed to get HEAD: %v\n", err)
        return
    }

    commit, err := repo.CommitObject(ref.Hash())
    if err != nil {
        fmt.Printf("Failed to get commit: %v\n", err)
        return
    }

    // Get the tree for the commit
    tree, err := commit.Tree()
    if err != nil {
        fmt.Printf("Failed to get tree: %v\n", err)
        return
    }

    // Walk through the tree and print file contents
    err = tree.Files().ForEach(func(f *object.File) error {
        fmt.Printf("\n--- File: %s ---\n", f.Name)
        contents, err := f.Contents()
        if err != nil {
            fmt.Printf("Error reading file: %v\n", err)
        } else {
            fmt.Println(contents)
        }
        return nil
    })

    if err != nil {
        fmt.Printf("Error walking through repository: %v\n", err)
    }
}
