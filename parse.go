package main

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/loader"
)

type RegoFile struct {
	Path []string    `json:"path"`
	Rego *ast.Module `json:"rego"`

	filePath string
}

func loadDirs(dirs ...string) ([]*RegoFile, error) {
	var resp []*RegoFile

	for _, dir := range dirs {
		files, err := loadFiles(dir)
		if err != nil {
			return nil, err
		}
		resp = append(resp, files...)
	}

	return resp, nil
}

func loadFiles(dir string) ([]*RegoFile, error) {
	var regoFiles []*RegoFile
	handler := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return goerr.Wrap(err).With("path", path)
		}

		if strings.HasSuffix(path, ".rego") {
			path := loader.CleanPath(path)
			logger.With("path", path).Debug("load .rego")
			regoFile, err := loader.RegoWithOpts(path, ast.ParserOptions{
				ProcessAnnotation: true,
			})
			if err != nil {
				return goerr.Wrap(err).With("path", path)
			}
			logger.With("path", path).With("rego", regoFile).Trace("loaded")

			regoFiles = append(regoFiles, &RegoFile{
				filePath: path,
				Rego:     regoFile.Parsed,
				Path:     strings.Split(path, string(filepath.Separator)),
			})
		}

		return nil
	}

	if err := filepath.WalkDir(dir, handler); err != nil {
		return nil, goerr.Wrap(err)
	}
	return regoFiles, nil
}
