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

const format = "\x1b[%dm"
const blue = 34
const red = 31
const yellow = 33 // TODO: replace with magenta
const reset = "\x1b[0m"

func colored(s string, color int) string {
	// ghetto AF. Refactor later.  the %42{...%} dance around \e[3m escape
	// codes is telling the string length to the terminal. It unfucks the
	// prompt getting out of control when autocompleting.
	return "%" + fmt.Sprintf("%d", len(s)) + "{" + fmt.Sprintf(format+s+reset, color) + "%}"
}

func main() {
	// prompt storage
	buf := new(bytes.Buffer)

	wd, err := os.Getwd()
	fatalIfError(err)

	usr, err := user.Current()
	var wdout string
	if err != nil {
		wdout = wd
	} else {
		r, err := filepath.Rel(usr.HomeDir, wd)
		fatalIfError(err)
		wdout = "~/" + r
	}
	buf.WriteString(colored(wdout, blue))

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

	var refout string
	if detached {
		name, err := head.Peel(git.ObjectAny)
		fatalIfError(err)
		headName, err := name.ShortId()
		fatalIfError(err)
		refout = headName
	} else {
		// Not detached, it should be in a branch
		branch, err := head.Branch().Name()
		fatalIfError(err)
		refout = "⌥ " + branch
	}

	out, err := exec.Command("git", "status", "--porcelain").Output()
	fatalIfError(err)
	if len(out) > 0 {
		// dirty
		refout = colored(refout, red)
	} else if detached {
		refout = colored(refout, yellow)
	}
	buf.WriteString(" " + refout)

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
	// buf.WriteTo(os.Stdout)
	fmt.Print(buf.String())
}

func fatalIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
