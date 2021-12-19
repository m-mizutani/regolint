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

	filePath string // original given file path
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
	basePath, err := filepath.Abs(dir)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	var regoFiles []*RegoFile
	handler := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return goerr.Wrap(err).With("path", path)
		}

		if strings.HasSuffix(path, ".rego") {
			relPath, err := filepath.Rel(basePath, path)
			if err != nil {
				return goerr.Wrap(err)
			}

			logger.With("path", path).With("rel", relPath).Debug("load .rego")
			regoFile, err := loader.RegoWithOpts(path, ast.ParserOptions{
				ProcessAnnotation: true,
			})
			if err != nil {
				return goerr.Wrap(err).With("path", relPath)
			}
			logger.With("path", relPath).With("rego", regoFile).Trace("loaded")

			regoFiles = append(regoFiles, &RegoFile{
				filePath: path,
				Rego:     regoFile.Parsed,
				Path:     strings.Split(relPath, string(filepath.Separator)),
			})
		}

		return nil
	}

	if err := filepath.WalkDir(basePath, handler); err != nil {
		return nil, goerr.Wrap(err)
	}
	return regoFiles, nil
}
