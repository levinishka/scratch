package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/causelovem/scratch/internal/generator"
)

func main() {
	projectPathPtr := flag.String("project", "", `path to new project directory
(last element in a path - project name)`)
	repoPtr := flag.String("repo", "", `git repository path for new project
(e.g. github.com/causelovem)`)
	helpPtr := flag.Bool("help", false, "prints this message")
	flag.Parse()

	// print help if needed
	if *helpPtr || flag.NFlag() < 2 {
		flag.Usage()
		return
	}

	projectPath := *projectPathPtr
	projectName := filepath.Base(*projectPathPtr)

	// create project directory
	if err := os.MkdirAll(projectPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// generate project code
	if err := generator.Generate(projectPath, projectName, *repoPtr); err != nil {
		log.Fatal(err)
	}

	log.Printf("Project '%s' successfully created at %s", projectName, projectPath)
}
