package config

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func RenameConnection(file *hclwrite.File, oldName, newName string) error {
	body := file.Body()

	var connBlock *hclwrite.Block
	for _, b := range body.Blocks() {
		if b.Type() == "connection" && b.Labels()[0] == oldName {
			connBlock = b
			break
		}
	}

	if connBlock == nil {
		return fmt.Errorf("connection '%s' not found", oldName)
	}

	for _, b := range body.Blocks() {
		if b.Type() == "connection" && b.Labels()[0] == newName {
			return fmt.Errorf("connection '%s' already exists", newName)
		}
	}

	connBlock.SetLabels([]string{newName})

	for _, b := range body.Blocks() {
		if b.Type() == "alias" {
			target := b.Body().GetAttribute("target")
			if target == nil {
				continue
			}

			val := target.Expr().BuildTokens(nil).Bytes()
			if string(val) == `"`+oldName+`"` {
				b.Body().SetAttributeValue("target", cty.StringVal(newName))
			}
		}
	}

	return nil
}
