// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package gcloud

import (
	"context"

	"github.com/martechassociation/martechassociation.github.io/f"
	"google.golang.org/api/option"
	sheets "google.golang.org/api/sheets/v4"
)

func (g *GCloud) ReadSheetColumns(
	token string,
	fileID string,
	sheetName string) ([]string, error) {
	ctx := g.ctx

	client, err := g.getSheetsClient(token)
	if err != nil {
		return nil, ctx.WrapError(err, err)
	}

	maxColsLen := GetSheetMaxColumnsLen()
	readRange := GetSheetRangeLabel(sheetName, 0, 0, 0, maxColsLen-1)

	resp, err := client.Spreadsheets.Values.Get(fileID, readRange).
		ValueRenderOption(string(GSheetValueTypeUnformattedValue)).
		Do()
	if err != nil {
		return nil, ctx.WrapError(err, err)
	}

	resItems := make([]string, 0)
	for _, row := range resp.Values {
		for _, cell := range row {
			resItems = append(resItems, f.InterfaceToString(cell))
		}
	}

	return resItems, nil
}

func (g *GCloud) ReadSheetMetaData(token string, fileID string) ([]*SheetMetaData, error) {
	ctx := g.ctx

	client, err := g.getSheetsClient(token)
	if err != nil {
		return nil, ctx.WrapError(err, err)
	}

	sheet, err := client.Spreadsheets.Get(fileID).IncludeGridData(false).Do()
	if err != nil {
		return nil, ctx.WrapError(err, err)
	}

	resItems := make([]*SheetMetaData, 0)
	for _, sheetItem := range sheet.Sheets {
		sheetProp := sheetItem.Properties
		// SheetType Possible values:
		//   "SHEET_TYPE_UNSPECIFIED" - Default value, do not use.
		//   "GRID" - The sheet is a grid.
		//   "OBJECT" - The sheet has no grid and instead has an object like a
		// chart or image.
		//   "DATA_SOURCE" - The sheet connects with an external DataSource and
		// shows the preview of data.
		sheetType := sheetProp.SheetType
		// We use only sheet type = GRID
		if sheetType != "GRID" {
			continue
		}

		title := sheetProp.Title
		columnCount := sheetProp.GridProperties.ColumnCount
		rowCount := sheetProp.GridProperties.RowCount

		resItems = append(resItems, &SheetMetaData{
			TabName:      title,
			ColumnsCount: int(columnCount),
			RowsCount:    int(rowCount),
		})
	}

	return resItems, nil
}

// OpenSheet open the spreadsheet
func (g *GCloud) OpenSheet(
	token string, fileID string, readRange string) (*GSheet, error) {
	return g.OpenSheetR(token, fileID, readRange, GSheetValueTypeFormattedValue)
}

// OpenSheetR open the spreadsheet with read option
func (g *GCloud) OpenSheetR(
	token string,
	fileID string,
	readRange string,
	renderOption GSheetValueType) (*GSheet, error) {
	ctx := g.ctx

	client, err := g.getSheetsClient(token)
	if err != nil {
		return nil, ctx.WrapError(err, err)
	}

	resp, err := client.Spreadsheets.Values.Get(fileID, readRange).
		ValueRenderOption(string(renderOption)).
		Do()
	if err != nil {
		return nil, ctx.WrapError(err, err)
	}

	retFile := &GSheet{
		Cells: resp.Values,
	}

	return retFile, nil
}

func (g *GCloud) getSheetsClient(token string) (*sheets.Service, error) {
	ctx := g.ctx

	svc, err := sheets.NewService(
		context.TODO(),
		option.WithCredentialsJSON([]byte(token)))
	if err != nil {
		return nil, ctx.WrapError(err, err)
	}

	return svc, nil
}
