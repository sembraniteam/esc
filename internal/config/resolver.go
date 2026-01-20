package config

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func ResolveName(cfg *Schema, name string) (string, error) {
	for _, a := range cfg.Aliases {
		if a.Name == name {
			return a.Target, nil
		}
	}

	return name, nil
}

func FindConnection(cfg *Schema, name string) (*Connection, error) {
	for _, c := range cfg.Connections {
		if c.Name == name {
			return &c, nil
		}
	}

	return nil, fmt.Errorf("connection '%s' not found", name)
}

func FindConnectionBlock(f *hclwrite.File, name string) *hclwrite.Block {
	for _, b := range f.Body().Blocks() {
		if b.Type() == "connection" && b.Labels()[0] == name {
			return b
		}
	}

	return nil
}

func StringValTrim(s string) cty.Value {
	return cty.StringVal(strings.TrimSpace(s))
}

func ResolveTargets(cfg *Schema, pattern string) []string {
	var result []string
	for _, c := range cfg.Connections {
		if Match(pattern, c.Name) {
			result = append(result, c.Name)
		}
	}
	return result
}
