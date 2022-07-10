package redfish

import (
	"encoding/json"

	"github.com/trungng1992/gofish/common"
)

type ArrayController struct {
	common.Entity

	// ODataContext is the odata context
	ODataContext string `json:"@odata.context"`
	// ODataType is the odata type.
	ODataType string `json:"@odata.type"`
	// Location shall contain location information of the
	// associated array controller.
	Location string
	// Manufacturer shall contain the name of the
	// organization responsible for producing the chassis. This organization
	// might be the entity from whom the chassis is purchased, but this is
	// not necessarily true.
	Manufacturer string
	// Model shall contain the name by which the
	// manufacturer generally refers to the array controller.
	Model string
	// SerialNumber shall contain a manufacturer-allocated
	// number that identifies the array controller.
	SerialNumber string

	// ReadCachePercent
	ReadCachePercent int

	// Status shall contain any status or health properties
	// of the resource.
	Status common.Status

	// PhysicalDrives
	physicaldrives     string
	storageenclosures  string
	logicaldrives      string
	unconfigureddrives string

	// rawData holds the original serialized JSON so we can compare updates.
	rawData []byte
}

// UnmarshalJSON unmarshals a Chassis object from the raw JSON.
func (arraycontroller *ArrayController) UnmarshalJSON(b []byte) error {
	type temp ArrayController
	type linkReference struct {
		LogicalDrives      common.Link
		StorageEnclosures  common.Link
		PhysicalDrives     common.Link
		UnconfiguredDrives common.Link
	}

	var t struct {
		temp
		Links linkReference
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}

	*arraycontroller = ArrayController(t.temp)

	// Extract the links to other entities for later
	arraycontroller.physicaldrives = string(t.Links.PhysicalDrives)
	arraycontroller.logicaldrives = string(t.Links.LogicalDrives)
	arraycontroller.storageenclosures = string(t.Links.StorageEnclosures)
	arraycontroller.unconfigureddrives = string(t.Links.UnconfiguredDrives)
	// This is a read/write object, so we need to save the raw object data for later
	arraycontroller.rawData = b

	return nil
}

func GetArrayController(c common.Client, uri string) (*ArrayController, error) {
	resp, err := c.Get(uri)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var arrayController ArrayController
	err = json.NewDecoder(resp.Body).Decode(&arrayController)
	if err != nil {
		return nil, err
	}

	arrayController.SetClient(c)
	return &arrayController, nil
}

func ListReferencedArrayControllers(c common.Client, link string) ([]*ArrayController, error) {
	var result []*ArrayController
	if link == "" {
		return result, nil
	}

	links, err := common.GetCollection(c, link)
	if err != nil {
		return result, err
	}

	collectionError := common.NewCollectionError()
	for _, arrayControllerLink := range links.ItemLinks {
		arrayController, err := GetArrayController(c, arrayControllerLink)
		if err != nil {
			collectionError.Failures[arrayControllerLink] = err
		} else {
			result = append(result, arrayController)
		}
	}

	if collectionError.Empty() {
		return result, nil
	}

	return result, collectionError
}

func (arrayController *ArrayController) PhysicalDrive() (*PhysicalDrive, error) {
	return ListReferencedPhysicalDrive(arrayController.Client, arrayController.physicaldrives)
}

func (arrayController *ArrayController) LogicalDrive() (*LogicalDrive, error) {
	return ListReferencedLogicalDrive(arrayController.Client, arrayController.logicaldrives)
}
