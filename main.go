package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/divpro/puml/output"
	"github.com/divpro/puml/parser"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("usage: puml <path to dir with go files> <path to out.puml>")
	}
	path := os.Args[1]
	files, err := readDir(path)
	if err != nil {
		log.Fatal(err)
	}

	var (
		wg            sync.WaitGroup
		mx            sync.Mutex
		parsedStructs []parser.PStruct
	)
	wg.Add(len(files))
	for _, file := range files {
		go func(file string) {
			defer wg.Done()
			mx.Lock()
			s, err := parser.ParseFile(file)
			if err != nil {
				log.Fatal(err.Error())
			}
			parsedStructs = append(parsedStructs, s...)
			mx.Unlock()
		}(file)
	}
	wg.Wait()

	out := output.Out(parsedStructs)
	err = os.WriteFile(os.Args[2], []byte(out), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func readDir(dir string) (files []string, err error) {
	if _, err := os.Stat(dir); err != nil {
		return nil, err
	}

	const ext = ".go"
	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ext {
			files = append(files, path)
		}
		return nil
	})

	return
}
