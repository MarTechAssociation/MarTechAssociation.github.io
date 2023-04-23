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
			sb.Append("\n")
			sb.Append(fmt.Sprintf(`## %s

|Platform Name|Description||
|---|---|---|`, page.Category))
		}

		currentCategory = string(page.Category)

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

[back](./)

## %s

> %s

%s

`, page.Name,
		page.Category,
		page.Description))

	// Write Contacts
	// 	### Contacts
	// | Channels        | Contact
	// |:----------------|:------------------------------------|
	// | Website         | [https://pams.ai](https://pams.ai)  |
	// | Email           | chananya@pams.ai                    |
	// | Mobile          | +66909832659                        |
	// | Facebook        | [https://www.facebook.com/PAMmarketingAutomation](https://www.facebook.com/PAMmarketingAutomation)|
	// | LINE            | line-id                             |
	sb.Append("\n")
	sb.Append(`### Contacts

| Channels        | Contact |
|:----------------|:------------------------------------|
`)
	if len(page.Website) > 0 {
		sb.Append(`| Website |`)
		sb.Append(fmt.Sprintf("[%s](%s)", page.Website, page.Website))
		sb.Append("|\n")
	}
	if len(page.Email) > 0 {
		sb.Append(`| Email |`)
		sb.Append(page.Email)
		sb.Append("|\n")
	}
	if len(page.Mobile) > 0 {
		sb.Append(`| Mobile |`)
		sb.Append(page.Mobile)
		sb.Append("|\n")
	}
	if len(page.Facebook) > 0 {
		sb.Append(`| Facebook |`)
		sb.Append(fmt.Sprintf("[%s](%s)", page.Facebook, page.Facebook))
		sb.Append("|\n")
	}
	if len(page.LINE) > 0 {
		sb.Append(`| LINE |`)
		sb.Append(page.LINE)
		sb.Append("|\n")
	}

	// Presentations images table
	sb.Append("\n")
	sb.Append(`### Presentations`)
	sb.Append("\n")
	err := svc.buildImagesTable(sb, page.PresentSlides, page.GetLandingPageFileName())
	if err != nil {
		return ctx.WrapError(err, err)
	}
	sb.Append("\n")

	// Screenshots images table
	sb.Append(`### Screenshots`)
	sb.Append("\n")
	err = svc.buildImagesTable(sb, page.ScreenSlides, page.GetLandingPageFileName())
	if err != nil {
		return ctx.WrapError(err, err)
	}

	// Write all content to file
	fileName := fmt.Sprintf("%s.md", page.GetLandingPageFileName())
	err = svc.writeFileContent(sb, fileName)
	if err != nil {
		return ctx.WrapError(err, err)
	}

	// Copy image to asset dir / delete the old one
	err = svc.copyImagesToAssetsDir(page)
	if err != nil {
		return ctx.WrapError(err, err)
	}
	return nil
}

func (svc *Generator) buildImagesTable(
	sb *StringBuilder,
	images []string,
	landingPageName string) error {

	sb.Append("\n")

	if len(images) == 0 {
		return nil
	}

	// <table>
	// 	<tr>
	// 	  <td><img src="assets/img/pam/presents-slides-00-27995.png"></td>
	// 	  <td><img src="assets/img/pam/presents-slides-01-16836.png"></td>
	//   </tr>
	//   <tr>
	// 	  <td><img src="assets/img/pam/presents-slides-02-14868.png"></td>
	// 	  <td><img src="assets/img/pam/presents-slides-03-30686.png"></td>
	// 	</tr>
	//   </table>

	sb.Append("<table>\n")
	// Draw 2 images per row
	for _, image := range images {
		sb.Append("<tr>\n")
		destPath := svc.imageAssetPath(image, landingPageName)
		sb.Append(fmt.Sprintf(`<td><img src="%s"></td>`, destPath))
		sb.Append("</tr>\n")
	}
	sb.Append("</table>\n")

	return nil
}

func (svc *Generator) copyImagesToAssetsDir(page *LandingPage) error {
	ctx := svc.ctx

	err := svc.createImageDir(page.GetLandingPageFileName())
	if err != nil {
		return ctx.WrapError(err, err)
	}

	// Copy presentations
	for _, filePath := range page.PresentSlides {
		destPath := svc.imageAssetPath(filePath, page.GetLandingPageFileName())
		err := svc.copyImage(filePath, destPath)
		if err != nil {
			ctx.WrapError(err, err)
			continue
		}
	}

	// Copy screenshots
	for _, filePath := range page.ScreenSlides {
		destPath := svc.imageAssetPath(filePath, page.GetLandingPageFileName())
		err := svc.copyImage(filePath, destPath)
		if err != nil {
			ctx.WrapError(err, err)
			continue
		}
	}
	return nil
}

func (svc *Generator) imageAssetPath(filePath string, landingPageName string) string {
	destFileName := filepath.Base(filePath)
	destPath := fmt.Sprintf("%s%s/%s",
		ImageAssetsDirPrefix,
		landingPageName,
		destFileName)
	return destPath
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
