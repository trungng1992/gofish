package redfish

import (
	"encoding/json"

	"github.com/trungng1992/gofish/common"
)

type TelemetryService struct {
	common.Entity

	// ODataContext is the odata context.
	ODataContext string `json:"@odata.context"`
	// ODataType is the odata type.
	ODataType string `json:"@odata.type"`

	// Description provides a description of this resource.
	Description string
	// ID uniquely identifies the resource.
	ID string `json:"Id"`

	// ServiceEnabled is used for the service. The value shall be true if
	// enabled and false if disabled.
	ServiceEnabled bool

	// Status shall contain any status or health properties
	// of the resource.
	Status common.Status

	logService              string
	metricDefinitions       string
	netricReportDefinitions string
	metricReports           string
	// rawData holds the original serialized JSON so we can compare updates.
	rawData []byte
}

// UnmarshallJSON unmarshals a TelemetryService object from raw JSON.
func (telemetryService *TelemetryService) UnmarshalJSON(b []byte) error {
	type temp TelemetryService
	type actions struct {
		TestMetric struct {
			ActionInfo string `json:"@Redfish.ActionInfo"`
			Target     string
		} `json:"#TelemetryService.SubmitTestMetricReport"`

		Oem json.RawMessage // OEM actions will be stored here
	}

	var t struct {
		temp
		LogService              common.Link
		MetricDefinitions       common.Link
		MetricReportDefinitions common.Link
		MetricReports           common.Link
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}
	*telemetryService = TelemetryService(t.temp)
	telemetryService.logService = string(t.LogService)
	telemetryService.metricDefinitions = string(t.MetricDefinitions)
	telemetryService.metricReports = string(t.MetricReports)
	telemetryService.rawData = b

	return nil
}

// ListReferencedTelemetryService gets the collection of TelemetryServices
func ListReferencedTelemetryService(c common.Client, link string) ([]*TelemetryService, error) {
	var result []*TelemetryService
	links, err := common.GetCollection(c, link)
	if err != nil {
		return result, err
	}

	collectionError := common.NewCollectionError()
	for _, telemetryServiceLink := range links.ItemLinks {

		telemetryService, err := GetTelemetryService(c, telemetryServiceLink)
		if err != nil {
			collectionError.Failures[telemetryServiceLink] = err
		} else {
			result = append(result, telemetryService)
		}
	}

	if collectionError.Empty() {
		return result, nil
	}

	return result, collectionError
}

// GetTelemetryService will get a TelemetryService instance from the Redfish service.
func GetTelemetryService(c common.Client, uri string) (*TelemetryService, error) {
	resp, err := c.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var telemetryService TelemetryService
	err = json.NewDecoder(resp.Body).Decode(&telemetryService)
	if err != nil {
		return nil, err
	}

	telemetryService.SetClient(c)
	return &telemetryService, nil
}

func (telemetryService *TelemetryService) MetricReports() ([]*MetricReport, error) {
	return ListReferencedMetricReports(telemetryService.Client, telemetryService.metricReports)
}
