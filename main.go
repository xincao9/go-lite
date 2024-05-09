package main

import (
	"embed"
	"flag"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/*
var embeddedFS embed.FS

type projectInfo struct {
	Project string
}

func main() {
	project := flag.String("p", "demo", "Project Name")
	if flag.Parsed() == false {
		flag.Parse()
	}
	err := copy(*project)
	if err != nil {
		_ = os.RemoveAll(*project)
		log.Fatalf("internal %v", err)
	}
	err = render(*project, &projectInfo{Project: *project})
	if err != nil {
		_ = os.RemoveAll(*project)
		log.Fatalf("internal %v", err)
	}
}

func render(tmplDir string, pi *projectInfo) error {
	return filepath.WalkDir(tmplDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".tmpl") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		err = os.Remove(path)
		if err != nil {
			return err
		}
		tmplName := d.Name()
		tmpl, err := template.New(tmplName).Parse(string(content))
		if err != nil {
			return err
		}
		path, _ = strings.CutSuffix(path, ".tmpl")
		localFile, err := os.Create(path)
		if err != nil {
			return err
		}
		err = tmpl.Execute(localFile, pi)
		if err != nil {
			return err
		}
		return localFile.Close()
	})
}

func copy(localDir string) error {
	err := os.MkdirAll(localDir, 0755)
	if err != nil {
		return err
	}
	return fs.WalkDir(embeddedFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		s, _ := strings.CutPrefix(path, "templates")
		localPath := filepath.Join(localDir, s)
		if localPath == "" || localPath == "." {
			return nil
		}
		if d.IsDir() {
			return os.MkdirAll(localPath, 0755)
		}
		embeddedFile, err := embeddedFS.Open(path)
		if err != nil {
			return err
		}
		defer embeddedFile.Close()
		localFile, err := os.Create(localPath)
		if err != nil {
			return err
		}
		defer localFile.Close()
		_, err = io.Copy(localFile, embeddedFile)
		return err
	})
}
