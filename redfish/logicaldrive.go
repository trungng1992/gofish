package redfish

import (
	"encoding/json"

	"github.com/trungng1992/gofish/common"
)

type LogicalDrive struct {
	common.Entity

	// ODataContext is the odata context
	ODataContext string `json:"@odata.context"`
	// ODataType is the odata type
	ODataType   string `json:"@odata.type"`
	VolumeCount int    `json:"Members@odata.count"`
	volumes     []string
}

func (logicaldrive *LogicalDrive) UnmarshalJSON(b []byte) error {
	type temp LogicalDrive
	var t struct {
		temp
		Volumes common.Links `json:"members"`
	}
	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}

	*logicaldrive = LogicalDrive(t.temp)
	// Extract the links to other entities for lates
	logicaldrive.volumes = t.Volumes.ToStrings()
	return nil
}

func GetLogicalDrive(c common.Client, uri string) (*LogicalDrive, error) {
	resp, err := c.Get(uri)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var logicalDrive LogicalDrive
	err = json.NewDecoder(resp.Body).Decode(&logicalDrive)
	if err != nil {
		return nil, err
	}

	logicalDrive.SetClient(c)
	return &logicalDrive, nil
}

func ListReferencedLogicalDrive(c common.Client, link string) (*LogicalDrive, error) {
	var result *LogicalDrive
	if link == "" {
		return result, nil
	}

	logicaldrive, err := GetLogicalDrive(c, link)
	if err != nil {
		return logicaldrive, err
	}

	return logicaldrive, nil
}

func (logicaldrive *LogicalDrive) Volumes() ([]*Logical, error) {
	var result []*Logical

	collectionError := common.NewCollectionError()
	for _, volumeLink := range logicaldrive.volumes {
		logical, err := GetLogical(logicaldrive.Client, volumeLink)
		if err != nil {
			collectionError.Failures[volumeLink] = err
		} else {
			result = append(result, logical)
		}
	}

	if collectionError.Empty() {
		return result, nil
	}

	return result, collectionError
}
