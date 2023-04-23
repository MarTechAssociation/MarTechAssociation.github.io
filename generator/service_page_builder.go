// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (svc *Generator) buildIndexMarkdown(pages []*LandingPage) error {
	ctx := svc.ctx

	sb := NewStringBuilder()
	sb.Append(`---
layout: default
---

The list of **MarTech platforms** below is categorized by types, for each platform in each type is order by _alphabetical_ order.

`)

	currentCategory := ""
	for _, page := range pages {
		// When category has changed, we write header for category
		if currentCategory != string(page.Category) {
			sb.Append("\n")
			sb.Append(fmt.Sprintf(`## %s

|Platform Name|Description||
|---|---|---|`, page.Category))
		}

		// Write each platform items
		// sample
		// |AIYA|A.I. For your business|[Get started](./aiya.html)|
		sb.Append("\n")
		sb.Append(fmt.Sprintf(`|%s|%s|[Get started](./%s.html)|`,
			page.Name,
			page.GetShortDescription(),
			page.GetLandingPageFileName()))
	}
	// Write index to the file index.md
	fileName := "index.md"

	err := svc.writeFileContent(sb, fileName)
	if err != nil {
		return ctx.WrapError(err, err)
	}

	return nil
}

func (svc *Generator) buildDetailMarkdown(page *LandingPage) error {
	ctx := svc.ctx

	sb := NewStringBuilder()

	sb.Append(fmt.Sprintf(`---
layout: default
---

## %s

[back](./)
`, page.Name))

	fileName := fmt.Sprintf("%s.md", page.GetLandingPageFileName())
	err := svc.writeFileContent(sb, fileName)
	if err != nil {
		return ctx.WrapError(err, err)
	}

	err = svc.copyImagesToAssetsDir(page)
	if err != nil {
		return ctx.WrapError(err, err)
	}
	return nil
}

func (svc *Generator) copyImagesToAssetsDir(page *LandingPage) error {
	ctx := svc.ctx

	err := svc.createImageDir(page.GetLandingPageFileName())
	if err != nil {
		return ctx.WrapError(err, err)
	}

	for _, filePath := range page.PresentSlides {
		destFileName := filepath.Base(filePath)
		destPath := fmt.Sprintf("%s%s/%s",
			ImageAssetsDirPrefix,
			page.GetLandingPageFileName(),
			destFileName)
		err := svc.copyImage(filePath, destPath)
		if err != nil {
			ctx.WrapError(err, err)
			continue
		}
	}
	return nil
}

func (svc *Generator) copyImage(sourceFile string, destFile string) error {
	ctx := svc.ctx

	ctx.Log(fmt.Sprintf("copy files fron %s to %s", sourceFile, destFile))

	source, err := os.Open(sourceFile)
	if err != nil {
		return ctx.WrapError(err, err)
	}
	defer source.Close()

	destination, err := os.Create(destFile)
	if err != nil {
		return ctx.WrapError(err, err)
	}
	defer destination.Close()

	ctx.Log(fmt.Sprintf("copying files fron %s to %s", sourceFile, destFile))

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destination, source)
	if err != nil {
		return ctx.WrapError(err, err)
	}

	return nil
}

func (svc *Generator) createImageDir(dirName string) error {
	ctx := svc.ctx

	if !strings.HasPrefix(dirName, ImageAssetsDirPrefix) {
		dirName = ImageAssetsDirPrefix + dirName
	}

	fi, err := os.Stat(dirName)
	if !os.IsNotExist(err) {
		if !fi.IsDir() {
			if err := os.Remove(dirName); err != nil {
				return ctx.WrapError(err, err)
			}
		} else {
			if err := os.RemoveAll(dirName); err != nil {
				return ctx.WrapError(err, err)
			}
		}
	}

	err = os.Mkdir(dirName, 0755)
	if err != nil {
		return ctx.WrapError(err, err)
	}

	return nil
}

func (svc *Generator) writeFileContent(sb *StringBuilder, fileName string) error {
	ctx := svc.ctx

	err := svc.deleteFileIfExists(fileName)
	if err != nil {
		ctx.WrapError(err, err)
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

func (svc *Generator) deleteFileIfExists(fileName string) error {
	ctx := svc.ctx
	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		err := os.Remove(fileName)
		if err != nil {
			return ctx.WrapError(err, err)
		}
	}
	return nil
}
