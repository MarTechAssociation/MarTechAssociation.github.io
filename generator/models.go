// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package generator

import (
	"github.com/martechassociation/martechassociation.github.io/f"
	"github.com/martechassociation/martechassociation.github.io/gcloud"
)

type SheetsResultSet struct {
	result  *gcloud.GSheet
	columns []string
	rowIdx  int
}

func NewSheetsResultSet(firstRow *gcloud.GSheet, contentRows *gcloud.GSheet) *SheetsResultSet {
	columns := buildColumnsFromSheets(firstRow)
	return &SheetsResultSet{
		result:  contentRows,
		columns: columns,
		rowIdx:  -1,
	}
}

func (rs *SheetsResultSet) Columns() []string {
	return rs.columns
}

func (rs *SheetsResultSet) Close() error {
	return nil
}

func (rs *SheetsResultSet) Err() error {
	return nil
}

func (rs *SheetsResultSet) Next() bool {
	nextRow := rs.rowIdx + 1
	if nextRow < len(rs.result.Cells) {
		rs.rowIdx = nextRow
		return true
	}
	return false
}

func (rs *SheetsResultSet) MapScan(item map[string]interface{}) error {
	colsLen := len(rs.columns)
	if colsLen == 0 {
		return nil
	}

	row := rs.result.Cells[rs.rowIdx]
	for colIdx, cell := range row {
		if colIdx >= colsLen {
			break
		}

		colVal := f.InterfaceToString(cell)
		colName := rs.columns[colIdx]
		item[colName] = colVal
	}

	return nil
}

type MarTechCategory string

const (
	MarketingAnalyticsPerformanceTrackingAndAttribution MarTechCategory = "Marketing Analytics, Performance Tracking & Attribution"
	CloudDataIntegrationPlatform                        MarTechCategory = "Cloud/Data Integration Platform"
	BusinessCustomerDataVisualizationTechnologies       MarTechCategory = "Business/Customer Data Visualization Technologies"
	ConversionRateOptimizationAndPersonalization        MarTechCategory = "Conversion Rate Optimization / Personalization"
	AdvertisingTechnology                               MarTechCategory = "Advertising Technology"
	VisitorIdentificationSoftware                       MarTechCategory = "Visitor Identification Software"
	AffiliateMarketing                                  MarTechCategory = "Affiliate Marketing"
	ContentMarketingTools                               MarTechCategory = "Content Marketing Tools"
	SEOtool                                             MarTechCategory = "SEO tool"
	SocialMediaMarketing                                MarTechCategory = "Social Media Marketing"
	CRMCDPMarketingAutomation                           MarTechCategory = "CRM, CDP, Marketing Automation"
	MarketingCloudSuites                                MarTechCategory = "Marketing Cloud Suites"
)

type LandingPage struct {
	Name          string          `json:"name"`
	Category      MarTechCategory `json:"category"`
	Description   string          `json:"description"`
	PresentSlides []string        `json:"present_slides"`
}

func (m *LandingPage) GetShortDescription() string {
	if len(m.Description) < MaxDescriptionInHomePage {
		return m.Description
	}
	return m.Description[:MaxDescriptionInHomePage] + "..."
}

func (m *LandingPage) GetLandingPageFileName() string {
	return f.EscapeName(m.Name)
}
