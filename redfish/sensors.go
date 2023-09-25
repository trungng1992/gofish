package redfish

import (
	"encoding/json"

	"github.com/trungng1992/gofish/common"
)

type Sensors struct {
	common.Entity

	// ODataContext is the odata context.
	ODataContext string `json:"@odata.context"`
	// ODataType is the odata type.
	ODataType string `json:"@odata.type"`

	// ArrayController
	Reading         float64
	ReadingRangeMax float64
	ReadingRangeMin float64
	ReadingType     string
	ReadingUnits    string

	ThresholdUpperFatal    float64
	ThresholdLowerCaution  float64
	ThresholdLowerCritical float64
	ThresholdLowerFatal    float64
	ThresholdUpperCaution  float64
	ThresholdUpperCritical float64
	Status                 common.Status
	// rawData holds the original serialized JSON so we can compare updates.
	rawData []byte
}

func (sensors *Sensors) UnmarshalJSON(b []byte) error {
	type temp Sensors

	type threshold struct {
		UpperFatal struct {
			Reading float64
		}
		LowerCaution struct {
			Reading float64
		}
		LowerCritical struct {
			Reading float64
		}
		LowerFatal struct {
			Reading float64
		}
		UpperCaution struct {
			Reading float64
		}
		UpperCritical struct {
			Reading float64
		}
	}
	var t struct {
		temp
		Thresholds threshold
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}
	*sensors = Sensors(t.temp)

	// Lower Bound
	sensors.ThresholdLowerCaution = t.Thresholds.LowerCaution.Reading
	sensors.ThresholdLowerCritical = t.Thresholds.LowerCritical.Reading
	sensors.ThresholdLowerFatal = t.Thresholds.LowerFatal.Reading

	// Upper Power
	sensors.ThresholdUpperFatal = t.Thresholds.UpperFatal.Reading
	sensors.ThresholdUpperCaution = t.Thresholds.UpperCaution.Reading
	sensors.ThresholdUpperCritical = t.Thresholds.UpperCritical.Reading
	sensors.rawData = b

	return nil
}

// GetMetricReport will get a metric report instance from the service.
func GetSensors(c common.Client, uri string) (*Sensors, error) {
	resp, err := c.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sensors Sensors
	err = json.NewDecoder(resp.Body).Decode(&sensors)
	if err != nil {
		return nil, err
	}

	sensors.SetClient(c)
	return &sensors, nil
}

func ListReferencedSensors(c common.Client, link string) ([]*Sensors, error) {
	var result []*Sensors
	if link == "" {
		return result, nil
	}

	links, err := common.GetCollection(c, link)
	if err != nil {
		return result, err
	}

	collectionError := common.NewCollectionError()
	for _, sensorsLink := range links.ItemLinks {
		sensors, err := GetSensors(c, sensorsLink)
		if err != nil {
			collectionError.Failures[sensorsLink] = err
		} else {
			result = append(result, sensors)
		}
	}

	if collectionError.Empty() {
		return result, nil
	}

	return result, collectionError
}
