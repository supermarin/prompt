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
	fatalIfError(err)

	buf.WriteString(path.Clean(wd))

	var repo *git.Repository
	if found, err := git.Discover(wd, false, nil); err == nil {
		wat, err := git.OpenRepository(found)
		fatalIfError(err)
		repo = wat
	} else {
		buf.WriteString(" $ ")
		buf.WriteTo(os.Stdout)
		os.Exit(0)
		return
	}

	// in git repo
	head, err := repo.Head()
	fatalIfError(err)

	// Add branch if found
	detached, err := repo.IsHeadDetached()
	fatalIfError(err)

	if detached {
		name, err := head.Peel(git.ObjectAny)
		fatalIfError(err)
		headName, err := name.ShortId()
		fatalIfError(err)
		buf.WriteString(" " + headName)
	} else {
		// Not detached, it should be in a branch
		branch, err := head.Branch().Name()
		fatalIfError(err)
		buf.WriteString(" ⌥ " + branch)
	}

	// dirty

	// stashes

	// final Exit
	buf.WriteString(" $ ")
	buf.WriteTo(os.Stdout)
}

func fatalIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
