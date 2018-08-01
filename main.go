package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/libgit2/git2go"
)

func main() {
	// prompt storage
	buf := new(bytes.Buffer)

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	fmt.Fprintf(buf, exPath)

	found, err := git.Discover(exPath, false, []string{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(found)
	fmt.Println(buf.String())
}
