package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

const (
	lineStart = 1
	colStart  = 1
	modeRW    = 0o600
	modeDir   = 0o700
)

func Init(path string) error {
	f := hclwrite.NewEmptyFile()
	body := f.Body()

	settings := body.AppendNewBlock("settings", nil).Body()
	settings.SetAttributeValue(
		"timeout_seconds",
		cty.NumberIntVal(DefaultTimeoutSeconds),
	)
	// settings.SetAttributeValue("prefer_ssh_agent", cty.BoolVal(true))
	body.AppendNewline()

	return WriteAtomic(path, f.Bytes())
}

func WriteAtomic(path string, data []byte) error {
	tmp := path + ".tmp"

	if err := os.WriteFile(tmp, data, modeRW); err != nil {
		return err
	}

	if err := os.Rename(tmp, path); err != nil {
		return err
	}

	if err := os.Chmod(path, modeRW); err != nil {
		return err
	}

	return nil
}

func Writable(path string) (*hclwrite.File, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = Init(path); err != nil {
			return nil, err
		}
	}

	// #nosec G304 -- path is user config location, not arbitrary input.
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	file, diags := hclwrite.ParseConfig(
		src,
		path,
		hcl.Pos{Line: lineStart, Column: colStart},
	)

	if diags.HasErrors() {
		return nil, errors.New(diags.Error())
	}

	return file, nil
}

func HasConnection(file *hclwrite.File, name string) bool {
	for _, block := range file.Body().Blocks() {
		if block.Type() == "connection" && block.Labels()[0] == name {
			return true
		}
	}
	return false
}

func AppendConnection(file *hclwrite.File, c Connection) error {
	body := file.Body()

	block := body.AppendNewBlock("connection", []string{c.Name})
	b := block.Body()

	b.SetAttributeValue("host", cty.StringVal(c.Host))
	b.SetAttributeValue("user", cty.StringVal(c.User))
	b.SetAttributeValue("port", cty.NumberIntVal(c.Port))

	auth := b.AppendNewBlock("auth", nil).Body()
	auth.SetAttributeValue("type", cty.StringVal(c.Auth.Type.String()))
	body.AppendNewline()

	return nil
}

func RemoveConnection(file *hclwrite.File, name string) bool {
	body := file.Body()

	for _, block := range body.Blocks() {
		if block.Type() == "connection" && block.Labels()[0] == name {
			body.RemoveBlock(block)
			return true
		}
	}

	return false
}

func EnsureConfigExists(path string) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, modeDir); err != nil {
		return err
	}

	if err := os.Chmod(dir, modeDir); err != nil {
		return err
	}

	if _, err := os.Stat(path); err == nil {
		return nil
	}

	return Init(path)
}
