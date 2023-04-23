// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package generator

import (
	"fmt"
	"strings"

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
		processRows := svc.buildLandingPage(rs, gg)
		if processRows == 0 {
			break
		}
	}

	ctx.Log("Generate landing pages done")

	return nil
}

func (svc *Generator) buildLandingPage(
	rs *SheetsResultSet,
	gg gcloud.IGCloud) int /*process rows*/ {

	ctx := svc.ctx

	i := -1

	for rs.Next() {
		i++

		item := map[string]interface{}{}
		err := rs.MapScan(item)
		if err != nil {
			ctx.WrapError(err, err)
			return i
		}

		for k, v := range item {
			ctx.Log(fmt.Sprintf("k=%s", k))
			ctx.Log(fmt.Sprintf("v=%s", v))
			if strings.HasPrefix(k, "Your Presentations") && v != nil {
				err := svc.buildPresentationThumbnail(v.(string), gg)
				if err != nil {
					ctx.WrapError(err, err)
				}
			}
			ctx.Log("---")
		}

	}

	return i
}

func (svc *Generator) buildPresentationThumbnail(
	slidesURL string,
	gg gcloud.IGCloud) error {

	ctx := svc.ctx
	cfg := svc.cfg

	token := cfg.GoogleToken()
	slideID := gcloud.GetSlideID(slidesURL)
	filePaths, err := gg.ReadSlidesThumbnails(token, slideID, 10)
	if err != nil {
		return ctx.WrapError(err, err)
	}

	ctx.LogObj("TEST", "filePaths", filePaths)
	return nil
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
