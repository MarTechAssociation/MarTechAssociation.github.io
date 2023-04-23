// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package gcloud

import (
	"fmt"
	"strings"
)

func GetSlideID(slideEndpoint string) string {
	// From endpoint: https://docs.google.com/presentation/d/1kq9cyMIxl3GAHXwx1kW46C_Zbti6A6udfljI2k9SWW0/edit#slide=id.g1969b0baa7c_0_4005
	// To slideID: 1kq9cyMIxl3GAHXwx1kW46C_Zbti6A6udfljI2k9SWW0
	slidePrefix := "https://docs.google.com/presentation/d/"
	// If not prefix with domain, we assume the endpoint is already be a sheet ID
	if !strings.HasPrefix(slideEndpoint, slidePrefix) {
		return slideEndpoint
	}

	// Otherwise, we remove all url parts from sheet ID
	slideEndpoint = strings.ReplaceAll(slideEndpoint, slidePrefix, "")
	urlParts := strings.Split(slideEndpoint, "/")
	if len(urlParts) > 0 {
		return urlParts[0]
	}

	// For default, we return slide endpoint
	return slideEndpoint
}

func GetSheetID(sheetEndpoint string) string {
	// From endpoint: https://docs.google.com/spreadsheets/d/1LMTjWG-PPYPmMZcQiKHZkCQqtO-zLQd_tGJV1uw_zWM/edit#gid=0
	// To sheetID: 1LMTjWG-PPYPmMZcQiKHZkCQqtO-zLQd_tGJV1uw_zWM
	sheetPrefix := "https://docs.google.com/spreadsheets/d/"
	// If not prefix with domain, we assume the endpoint is already be a sheet ID
	if !strings.HasPrefix(sheetEndpoint, sheetPrefix) {
		return sheetEndpoint
	}

	// Otherwise, we remove all url parts from sheet ID
	sheetEndpoint = strings.ReplaceAll(sheetEndpoint, sheetPrefix, "")
	urlParts := strings.Split(sheetEndpoint, "/")
	if len(urlParts) > 0 {
		return urlParts[0]
	}

	// For default, we return sheet endpoint
	return sheetEndpoint
}

func GetSheetRangeLabel(sheetName string, fromRow int, fromCol int, toRow int, toCol int) string {
	if fromRow < 0 || toRow < 0 {
		return ""
	}
	fromColLabel := GetSheetColumnLabel(fromCol)
	if len(fromColLabel) == 0 {
		return ""
	}
	toColLabel := GetSheetColumnLabel(toCol)
	if len(toColLabel) == 0 {
		return ""
	}

	if len(sheetName) == 0 {
		return fmt.Sprintf(`%s%d:%s%d`, fromColLabel, fromRow+1, toColLabel, toRow+1)
	}
	return fmt.Sprintf(`'%s'!%s%d:%s%d`, sheetName, fromColLabel, fromRow+1, toColLabel, toRow+1)
}

func GetSheetRowLabel(row int, col int) string {
	if row < 0 {
		return ""
	}
	colLabel := GetSheetColumnLabel(col)
	if len(colLabel) == 0 {
		return ""
	}

	return fmt.Sprintf("%s%d", colLabel, row+1)
}

func GetSheetMaxColumnsLen() int {
	colsLen := len(SheetsColumnLetters)
	// support column name only 2 digits column (AB, BA, ZZ, ..) and single digit column (A, B, C, .., Z)
	return (colsLen + (colsLen * colsLen))
}

func GetSheetColumnLabel(col int) string {
	colsLen := len(SheetsColumnLetters)

	// support column name only 2 digits column (AB, BA, ZZ, ..) and single digit column (A, B, C, .., Z)
	maxColsLen := GetSheetMaxColumnsLen()

	if col < 0 || col >= maxColsLen {
		return ""
	}

	lastDigit := col % colsLen
	lastChar := SheetsColumnLetters[lastDigit]

	firstDigit := col / colsLen
	firstChar := ""
	if firstDigit > 0 {
		firstChar = SheetsColumnLetters[firstDigit-1]
	}

	return fmt.Sprintf("%s%s", firstChar, lastChar)
}
