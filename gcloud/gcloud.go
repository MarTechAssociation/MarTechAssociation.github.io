// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package gcloud

import (
	"github.com/martechassociation/martechassociation.github.io/microservices"
)

// IGCloud is interface use to connect google drive
type IGCloud interface {
	ReadSheetMetaData(token string, fileID string) ([]*SheetMetaData, error)
	ReadSheetColumns(token string, fileID string, tabName string) ([]string, error)
	OpenSheet(token string, fileID string, readRange string) (*GSheet, error)
	OpenSheetR(token string, fileID string, readRange string, renderOption GSheetValueType) (*GSheet, error)

	ReadSlidesThumbnails(
		token string,
		fileID string,
		filePrefix string,
		maxSlides int) ([]string /*file paths*/, error)
}

// GFile is struct represent file in gdrive
type GFile struct {
	ID       string
	FileName string
	Kind     string
}

// GSheet is struct represent sheet file
type GSheet struct {
	Cells [][]interface{}
}

type SheetMetaData struct {
	TabName      string `json:"tab_name"`
	ColumnsCount int    `json:"columns_count"`
	RowsCount    int    `json:"rows_count"`
}

type GSheetValueType string

//   "FORMATTED_VALUE" - Values will be calculated & formatted in the
// reply according to the cell's formatting. Formatting is based on the
// spreadsheet's locale, not the requesting user's locale. For example,
// if `A1` is `1.23` and `A2` is `=A1` and formatted as currency, then
// `A2` would return "$1.23".
//   "UNFORMATTED_VALUE" - Values will be calculated, but not formatted
// in the reply. For example, if `A1` is `1.23` and `A2` is `=A1` and
// formatted as currency, then `A2` would return the number `1.23`.
//   "FORMULA" - Values will not be calculated. The reply will include
// the formulas. For example, if `A1` is `1.23` and `A2` is `=A1` and
// formatted as currency, then A2 would return "=A1".

const (
	GSheetValueTypeFormula          GSheetValueType = "FORMULA"
	GSheetValueTypeUnformattedValue GSheetValueType = "UNFORMATTED_VALUE"
	GSheetValueTypeFormattedValue   GSheetValueType = "FORMATTED_VALUE"
)

// GCloud is struct implement IGCloud
type GCloud struct {
	ctx microservices.IContext
	cfg microservices.IConfig
}

// NewGCloud return new gdrive service
func NewGCloud(ctx microservices.IContext, cfg microservices.IConfig) *GCloud {
	return &GCloud{
		ctx: ctx,
		cfg: cfg,
	}
}
