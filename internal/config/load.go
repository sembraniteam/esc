package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

func Load(path string) (*Schema, error) {
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return nil, errors.New(diags.Error())
	}

	if file == nil {
		return nil, errors.New("no config file found")
	}

	var schema Schema
	if diags = gohcl.DecodeBody(
		file.Body,
		&hcl.EvalContext{},
		&schema,
	); diags.HasErrors() {
		return nil, errors.New(diags.Error())
	}

	return &schema, nil
}

func WithLock(path string, fn func() error) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, modeDir); err != nil {
		return err
	}

	_ = os.Chmod(dir, modeDir)

	lock := flock.New(path)

	locked, err := lock.TryLock()
	if err != nil {
		return err
	}

	if !locked {
		return fmt.Errorf("ESC config is currently in use by another process")
	}

	_ = os.Chmod(path, modeRW)
	defer func(lock *flock.Flock) {
		if err = lock.Unlock(); err != nil {
			panic(err)
		}
	}(lock)

	return fn()
}
