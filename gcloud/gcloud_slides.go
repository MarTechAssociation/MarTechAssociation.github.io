// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package gcloud

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/martechassociation/martechassociation.github.io/f"
	"google.golang.org/api/option"
	slides "google.golang.org/api/slides/v1"
)

func (g *GCloud) ReadSlidesThumbnails(
	token string,
	fileID string,
	filePrefix string,
	maxSlides int) ([]string /*file paths*/, error) {

	ctx := g.ctx

	client, err := g.getSlidesClient(token)
	if err != nil {
		return nil, ctx.WrapError(err, err)
	}

	present, err := client.Presentations.Get(fileID).Do()
	if err != nil {
		return nil, ctx.WrapError(err, err)
	}

	retFiles := []string{}

	for i, page := range present.Slides {
		if i >= maxSlides {
			break
		}

		// Retrieve the thumbnail for the page
		// resp = https://url-to-download-thumbnail
		resp, err := client.Presentations.Pages.GetThumbnail(fileID, page.ObjectId).
			ThumbnailPropertiesMimeType("PNG").
			ThumbnailPropertiesThumbnailSize("MEDIUM").
			Do()
		if err != nil {
			return nil, ctx.WrapError(err, err)
		}

		// Download the image from the URL
		respBody, err := http.Get(resp.ContentUrl)
		if err != nil {
			return nil, ctx.WrapError(err, err)
		}
		defer respBody.Body.Close()

		// Write the image data to a file
		fileSuffix := f.RandomMinMax(10001, 99999)
		filePath := fmt.Sprintf("%s%s-slides-%02d-%d.png", os.TempDir(), filePrefix, i, fileSuffix)

		// Save the image to a file
		outFile, err := os.Create(filePath)
		if err != nil {
			return nil, ctx.WrapError(err, err)
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, respBody.Body)
		if err != nil {
			return nil, ctx.WrapError(err, err)
		}

		retFiles = append(retFiles, filePath)

	}

	return retFiles, nil
}

func (g *GCloud) getSlidesClient(token string) (*slides.Service, error) {
	ctx := g.ctx

	svc, err := slides.NewService(
		context.TODO(),
		option.WithCredentialsJSON([]byte(token)))
	if err != nil {
		return nil, ctx.WrapError(err, err)
	}

	return svc, nil
}
