package config

import (
	"sort"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

func SortBlocks(file *hclwrite.File) {
	oldBody := file.Body()
	newFile := hclwrite.NewEmptyFile()
	newBody := newFile.Body()

	var (
		settings    []*hclwrite.Block
		connections []*hclwrite.Block
		aliases     []*hclwrite.Block
		others      []*hclwrite.Block
	)

	for _, b := range oldBody.Blocks() {
		switch b.Type() {
		case "settings":
			settings = append(settings, b)
		case "connection":
			connections = append(connections, b)
		case "alias":
			aliases = append(aliases, b)
		default:
			others = append(others, b)
		}
	}

	sort.Slice(connections, func(i, j int) bool {
		return connections[i].Labels()[0] < connections[j].Labels()[0]
	})

	sort.Slice(aliases, func(i, j int) bool {
		return aliases[i].Labels()[0] < aliases[j].Labels()[0]
	})

	appendAll := func(blocks []*hclwrite.Block) {
		for _, b := range blocks {
			newBody.AppendBlock(b)
			newBody.AppendNewline()
		}
	}

	appendAll(settings)
	appendAll(connections)
	appendAll(aliases)
	appendAll(others)

	*file = *newFile
}
