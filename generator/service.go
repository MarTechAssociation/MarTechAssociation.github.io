// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package generator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/martechassociation/martechassociation.github.io/f"
	"github.com/martechassociation/martechassociation.github.io/gcloud"
	"github.com/martechassociation/martechassociation.github.io/microservices"
)

// Generator is service to genereate martech page files
type Generator struct {
	ctx microservices.IContext
	cfg microservices.IConfig
}

// NewGenerator return new generator service
func NewGenerator(ctx microservices.IContext, cfg microservices.IConfig) *Generator {
	return &Generator{
		ctx: ctx,
		cfg: cfg,
	}
}

func (svc *Generator) GenerateLandingPages() error {
	ctx := svc.ctx
	cfg := svc.cfg

	gg := gcloud.NewGCloud(ctx, cfg)

	sheetName := cfg.MarTechSheetName()
	token := cfg.GoogleToken()

	maxColsLen := gcloud.GetSheetMaxColumnsLen()
	rangeLabel := gcloud.GetSheetRangeLabel(sheetName, 0, 0, 0, maxColsLen-1)
	sheetID := gcloud.GetSheetID(cfg.MarTechSheetURL())
	firstRow, err := gg.OpenSheetR(cfg.GoogleToken(), sheetID, rangeLabel, gcloud.GSheetValueTypeUnformattedValue)
	if err != nil {
		return ctx.WrapError(err, err)
	}

	// fromRow=1, fromCol=0, toRow=1000, toCol=colLen-1
	columns := buildColumnsFromSheets(firstRow)
	fromCol := 0
	toCol := len(columns) - 1
	fromRow := 1

	allPages := []*LandingPage{}
	for {

		// Move read last row pointer to the next 1000 records
		toRow := fromRow + SheetsReadBatchSize
		dataRange := gcloud.GetSheetRangeLabel(sheetName, fromRow, fromCol, toRow, toCol)
		sheets, err := gg.OpenSheetR(token, sheetID, dataRange, gcloud.GSheetValueTypeUnformattedValue)
		if err != nil {
			isExceedMaxRow := strings.Contains(strings.ToLower(err.Error()), "exceeds grid limits")
			if !isExceedMaxRow {
				return ctx.WrapError(err, err)
			}
		}

		rs := NewSheetsResultSet(firstRow, sheets)
		pages := svc.buildLandingPage(rs, gg)
		if len(pages) == 0 {
			break
		}

		allPages = append(allPages, pages...)
		fromRow += len(allPages)
	}

	// Sort by category and name
	sort.Slice(allPages, func(i, j int) bool {
		return allPages[i].Category < allPages[j].Category &&
			allPages[i].GetName() < allPages[j].GetName()
	})

	ctx.LogObj("RESULT", "allPages", allPages)

	if len(allPages) > 0 {
		err := svc.buildIndexMarkdown(allPages)
		if err != nil {
			ctx.WrapError(err, err)
		}

		for _, page := range allPages {
			err := svc.buildDetailMarkdown(page)
			if err != nil {
				ctx.WrapError(err, err)
			}
		}
	}

	return nil
}

func (svc *Generator) buildLandingPage(
	rs *SheetsResultSet,
	gg gcloud.IGCloud) []*LandingPage {

	ctx := svc.ctx

	pages := []*LandingPage{}
	i := -1

	for rs.Next() {
		i++

		item := map[string]interface{}{}
		err := rs.MapScan(item)
		if err != nil {
			ctx.WrapError(err, err)
			return pages
		}

		page := &LandingPage{}
		for k, v := range item {
			ctx.Log(fmt.Sprintf("i=%d", i))
			ctx.Log(fmt.Sprintf("k=%s", k))
			ctx.Log(fmt.Sprintf("v=%s", v))
			ctx.Log("---")

			if strings.HasPrefix(k, "Your MarTech name") {
				page.Name = f.InterfaceToString(v)
			} else if strings.HasPrefix(k, "What is your MarTech categories?") {
				page.Category = MarTechCategory(f.InterfaceToString(v))
			} else if strings.HasPrefix(k, "Short explanation of what is your product or service") {
				page.Description = f.InterfaceToString(v)
			} else if strings.HasPrefix(k, "Email for customer to contact") {
				page.Email = f.InterfaceToString(v)
			} else if strings.HasPrefix(k, "Mobile for customer to contact") {
				page.Mobile = f.InterfaceToString(v)
			} else if strings.HasPrefix(k, "Your Facebook Page for customer to contact") {
				page.Facebook = f.InterfaceToString(v)
			} else if strings.HasPrefix(k, "Your Website for customer to contact") {
				page.Website = f.InterfaceToString(v)
			} else if strings.HasPrefix(k, "Your LINE ID or LINE OA") {
				page.LINE = f.InterfaceToString(v)
			} else if strings.HasPrefix(k, "Sample user interfaces of your software") {
				url := f.InterfaceToString(v)
				filePaths, err := svc.buildSlidesThumbnails(url, "screens", gg)
				if err != nil {
					ctx.WrapError(err, err)
					continue
				}
				page.ScreenSlides = filePaths
			} else if strings.HasPrefix(k, "Your Presentations") && v != nil {
				url := f.InterfaceToString(v)
				filePaths, err := svc.buildSlidesThumbnails(url, "presents", gg)
				if err != nil {
					ctx.WrapError(err, err)
					continue
				}
				page.PresentSlides = filePaths
			}
		}
		pages = append(pages, page)
	}

	return pages
}

func (svc *Generator) buildSlidesThumbnails(
	slidesURL string,
	imagePrefix string,
	gg gcloud.IGCloud) ([]string, error) {

	ctx := svc.ctx
	cfg := svc.cfg

	token := cfg.GoogleToken()
	slideID := gcloud.GetSlideID(slidesURL)
	filePaths, err := gg.ReadSlidesThumbnails(token, slideID, imagePrefix, MaxSlidesToRead)
	if err != nil {
		return nil, ctx.WrapError(err, err)
	}

	ctx.LogObj("TEST", "filePaths", filePaths)
	return filePaths, nil
}

func buildColumnsFromSheets(sheet *gcloud.GSheet) []string {
	columns := []string{}
	if len(sheet.Cells) > 0 {
		row := sheet.Cells[0]
		for _, cell := range row {
			colName := fmt.Sprintf("%v", cell)
			if len(colName) == 0 {
				break
			}
			columns = append(columns, colName)
		}
	}
	return columns
}
