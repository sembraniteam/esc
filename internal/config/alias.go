package config

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func HasAlias(file *hclwrite.File, name string) bool {
	for _, b := range file.Body().Blocks() {
		if b.Type() == "alias" && b.Labels()[0] == name {
			return true
		}
	}
	return false
}

func AppendAlias(file *hclwrite.File, name, target string) {
	body := file.Body()
	block := body.AppendNewBlock("alias", []string{name})
	b := block.Body()
	b.SetAttributeValue("target", cty.StringVal(target))
	body.AppendNewline()
}

func RemoveAlias(file *hclwrite.File, name string) bool {
	body := file.Body()
	for _, b := range body.Blocks() {
		if b.Type() == "alias" && b.Labels()[0] == name {
			body.RemoveBlock(b)
			return true
		}
	}
	return false
}

func ResolveAlias(cfg *Schema, name string) (string, bool) {
	for _, a := range cfg.Aliases {
		if a.Name == name {
			return a.Target, true
		}
	}
	return "", false
}

func EnsureTargetExists(cfg *Schema, target string) error {
	for _, c := range cfg.Connections {
		if c.Name == target {
			return nil
		}
	}
	return fmt.Errorf("target connection '%s' does not exist", target)
}
