package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/libgit2/git2go"
)

func main() {
	// prompt storage
	buf := new(bytes.Buffer)

	wd, err := os.Getwd()
	fatalIfError(err)

	usr, err := user.Current()
	if err != nil {
		buf.WriteString(wd)
	} else {
		r, err := filepath.Rel(usr.HomeDir, wd)
		fatalIfError(err)
		buf.WriteString("~/" + r)
	}

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
	_, err = exec.Command("git", "status", "--porcelain").Output()
	if err == nil {
		buf.WriteString(" M")
	}

	// stashes
	var count int
	repo.Stashes.Foreach(func(index int, msg string, id *git.Oid) error {
		count++
		return nil
	})

	if count > 0 {
		fmt.Fprintf(buf, " S%v", count)
	}
	// final Exit
	buf.WriteString(" $ ")
	buf.WriteTo(os.Stdout)
}

func fatalIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
