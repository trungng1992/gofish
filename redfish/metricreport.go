package redfish

import (
	"encoding/json"

	"github.com/trungng1992/gofish/common"
)

type MetricValue struct {
	MetricProperty string
	MetricValue    string
}

// MetricReports is used to represent a metric reports
// medium for a Redfish implementation.
type MetricReport struct {
	common.Entity
	// ODataContext is the odata context.
	ODataContext string `json:"@odata.context"`
	// ODataType is the odata type.
	ODataType string `json:"@odata.type"`

	// MetricReportDefinition is metric report definition of telemetry
	metricReportDefinition string
	MetricValues           []MetricValue
	Name                   string
	// rawData holds the original serialized JSON so we can compare updates.
	rawData []byte
}

func (metricReport *MetricReport) UnmarshalJSON(b []byte) error {
	type temp MetricReport
	var t struct {
		temp
		MetricReportDefinition common.Link
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}
	*metricReport = MetricReport(t.temp)
	metricReport.metricReportDefinition = string(t.MetricReportDefinition)
	metricReport.rawData = b

	return nil
}

// GetMetricReport will get a metric report instance from the service.
func GetMetricReports(c common.Client, uri string) (*MetricReport, error) {
	resp, err := c.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var metricReport MetricReport
	err = json.NewDecoder(resp.Body).Decode(&metricReport)
	if err != nil {
		return nil, err
	}

	metricReport.SetClient(c)
	return &metricReport, nil
}

func ListReferencedMetricReports(c common.Client, link string) ([]*MetricReport, error) {
	var result []*MetricReport
	if link == "" {
		return result, nil
	}

	links, err := common.GetCollection(c, link)
	if err != nil {
		return result, err
	}

	collectionError := common.NewCollectionError()
	for _, metricReportsLink := range links.ItemLinks {
		metricReports, err := GetMetricReports(c, metricReportsLink)
		if err != nil {
			collectionError.Failures[metricReportsLink] = err
		} else {
			result = append(result, metricReports)
		}
	}

	if collectionError.Empty() {
		return result, nil
	}

	return result, collectionError
}
