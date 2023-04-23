// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package generator

import (
	"os"
)

func (svc *Generator) buildIndexMarkdown(pages []*LandingPage) error {
	ctx := svc.ctx

	sb := NewStringBuilder()
	sb.Append(`---
layout: default
---

The list of **MarTech platforms** below is categorized by types, for each platform in each type is order by _alphabetical_ order.

`)

	fileName := "index.md"

	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		err := os.Remove(fileName)
		if err != nil {
			ctx.WrapError(err, err)
		}
	}

	file, err := os.Create(fileName)
	if err != nil {
		return ctx.WrapError(err, err)
	}

	defer file.Close()

	_, err = file.WriteString(sb.String())
	if err != nil {
		return ctx.WrapError(err, err)
	}

	return nil
}

func (svc *Generator) buildDetailMarkdown(page *LandingPage) error {
	return nil
}
