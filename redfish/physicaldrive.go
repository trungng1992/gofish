package redfish

import (
	"encoding/json"

	"github.com/trungng1992/gofish/common"
)

type PhysicalDrive struct {
	common.Entity

	// ODataContext is the odata context
	ODataContext string `json:"@odata.context"`
	// ODataType is the odata type.
	ODataType string `json:"@odata.type"`

	DrivesCount int `json:"Members@odata.count"`

	drives []string
}

// UnmarshallJSON unmarshalls a PhysicalDrive object from the raw JSON
func (physicaldrive *PhysicalDrive) UnmarshalJSON(b []byte) error {
	type temp PhysicalDrive
	var t struct {
		temp
		Drives common.Links `json:"members"`
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}

	*physicaldrive = PhysicalDrive(t.temp)

	// Extract the links to other entities for lates
	physicaldrive.drives = t.Drives.ToStrings()
	return nil
}

func (physicaldrive *PhysicalDrive) GetListDrives() []string {
	return physicaldrive.drives
}

func GetPhysicalDrive(c common.Client, uri string) (*PhysicalDrive, error) {
	resp, err := c.Get(uri)
	if err != nil {
		return nil, err
	}

	var physicalDrive PhysicalDrive
	err = json.NewDecoder(resp.Body).Decode(&physicalDrive)
	if err != nil {
		return nil, err
	}

	physicalDrive.SetClient(c)
	return &physicalDrive, nil
}

func ListReferencedPhysicalDrive(c common.Client, link string) (*PhysicalDrive, error) {
	var result *PhysicalDrive
	if link == "" {
		return result, nil
	}

	physicaldrive, err := GetPhysicalDrive(c, link)
	if err != nil {
		return physicaldrive, err
	}

	return physicaldrive, nil
}

func (physicaldrive *PhysicalDrive) Drives() ([]*DiskDrive, error) {
	var result []*DiskDrive

	collectionError := common.NewCollectionError()
	for _, driveLink := range physicaldrive.drives {
		drive, err := GetDiskDrive(physicaldrive.Client, driveLink)
		if err != nil {
			collectionError.Failures[driveLink] = err
		} else {
			result = append(result, drive)
		}
	}

	if collectionError.Empty() {
		return result, nil
	}

	return result, collectionError
}
