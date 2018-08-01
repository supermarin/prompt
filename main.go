package main

import (
	"bytes"
	"log"
	"os"
	"path"

	"github.com/libgit2/git2go"
)

func main() {
	// prompt storage
	buf := new(bytes.Buffer)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	buf.WriteString(path.Clean(wd))

	var repo *git.Repository
	if found, err := git.Discover(wd, false, nil); err == nil {
		wat, err := git.OpenRepository(found)
		if err != nil {
			log.Fatal(err)
		}
		repo = wat
	} else {
		buf.WriteString(" $ ")
		buf.WriteTo(os.Stdout)
		os.Exit(0)
		return
	}

	// in git repo
	head, err := repo.Head()
	if err != nil {
		log.Fatal(err)
	}

	// Add branch if found
	detached, err := repo.IsHeadDetached()
	if err != nil {
		log.Fatal(err)
	}

	if detached {
		name, err := head.Peel(git.ObjectAny)
		if err != nil {
			log.Fatalf("Can't resolve! %v", err)
		}
		str, err := name.ShortId()
		buf.WriteString(" " + str)
	} else {
		branch, err := head.Branch().Name()
		if err != nil {
			log.Fatal(err)
		}
		buf.WriteString(" ⌥ " + branch)
	}

	// dirty

	// stashes

	// final Exit
	buf.WriteString(" $ ")
	buf.WriteTo(os.Stdout)
}
