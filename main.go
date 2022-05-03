package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/CheHsuan/markdown-tree/pkg/mdtree"
)

var (
	inputDir   = flag.String("input-dir", "", "input directory path")
	outputFile = flag.String("output-file", "", "output file path")
)

func main() {
	flag.Parse()

	if *inputDir == "" || *outputFile == "" {
		log.Panic("missing input directory path or output file path")
	}

	tree := mdtree.NewMarkdownTree(filepath.Base(*inputDir))

	wd, err := os.Getwd()
	if err != nil {
		log.Panic(err.Error())
	}

	// change working dir to target dir
	if err := os.Chdir(*inputDir); err != nil {
		log.Panic(err.Error())
	}

	// walk through the directory
	if err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && info.Name() == ".git" {
				return filepath.SkipDir
			}

			tree.InsertNode(path, info.IsDir())

			return nil
		}); err != nil {
		log.Panic(err.Error())
	}

	// change working dir back
	if err := os.Chdir(wd); err != nil {
		log.Panic(err.Error())
	}

	if err := tree.GenerateMarkdown(*outputFile); err != nil {
		log.Panic(err.Error())
	}

	log.Printf("done")
}
