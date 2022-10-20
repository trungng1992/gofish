//
// SPDX-License-Identifier: BSD-3-Clause
//

package redfish

import (
	"encoding/json"
	"reflect"

	"github.com/trungng1992/gofish/common"
)

// DiskDrive is used to represent a disk drive or other physical storage
// medium for a Redfish implementation.
type DiskDrive struct {
	common.Entity

	// ODataContext is the odata context.
	ODataContext string `json:"@odata.context"`
	// ODataType is the odata type.
	ODataType string `json:"@odata.type"`
	// BlockSizeBytes shall contain size of the smallest addressable unit of the
	// associated drive.
	BlockSizeBytes int
	// CapableSpeedGbs shall contain fastest capable bus speed of the associated
	// drive.
	CapableSpeedGbs float32
	// CapacityBytes shall contain the raw size in bytes of the associated drive.
	CapacityBytes int64
	CapacityGB    int
	CapacityMiB   int
	// Name
	Name string
	// CapacityLogicalBlocks
	CapacityLogicalBlocks int64
	// Description provides a description of this resource.
	Description string
	// FailurePredicted shall contain failure information as defined by the
	// manufacturer for the associated drive.
	FailurePredicted bool
	// Location shall contain location information of the associated drive.
	Location string
	// Manufacturer shall be the name of the organization responsible for
	// producing the drive. This organization might be the entity from whom the
	// drive is purchased, but this is not necessarily true.
	Manufacturer string
	// MediaType shall contain the type of media contained in the associated
	// drive.
	MediaType MediaType
	// Model shall be the name by which the manufacturer generally refers to the
	// drive.
	Model string
	// Multipath shall indicate whether the drive is
	// accessible by an initiator from multiple paths allowing for failover
	// capabilities upon a path failure.
	Multipath bool
	// NegotiatedSpeedGbs shall contain current bus speed of the associated
	// drive.
	NegotiatedSpeedGbs float32

	// PartNumber shall be a part number assigned by the organization that is
	// responsible for producing or manufacturing the drive.
	PartNumber string
	// PredictedMediaLifeLeftPercent shall contain an indicator of the
	// percentage of life remaining in the Drive's media.
	PredictedMediaLifeLeftPercent float32
	//InterfaceSpeedMbps
	InterfaceSpeedMbps float32
	// Revision
	Revision string
	// RotationSpeedRPM shall contain rotation speed of the associated drive.
	RotationSpeedRPM float32
	// SerialNumber is used to identify the drive.
	SerialNumber string
	// Status shall contain any status or health properties of the resource.
	Status common.Status
	// WriteCacheEnabled shall indicate whether the drive
	// write cache is enabled.
	WriteCacheEnabled bool

	PowerOnHours float32

	Protocol               string
	UncorrectedReadErrors  int
	UncorrectedWriteErrors int
	// Temperature
	CurrentTemperatureCelsius int
	MaximumTemperatureCelsius int
	// rawData holds the original serialized JSON so we can compare updates.
	rawData []byte
}

// UnmarshalJSON unmarshals a Drive object from the raw JSON.
func (drive *DiskDrive) UnmarshalJSON(b []byte) error {
	type temp DiskDrive

	var t struct {
		temp
		SSDEnduranceUtilizationPercentage float32
		RotationalSpeedRpm                float32
		InterfaceType                     string
		FirmwareVersion                   common.FirmwareVersion
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}

	// Extract the links to other entities for later
	*drive = DiskDrive(t.temp)
	drive.CapacityBytes = int64(t.CapacityMiB * 1e6)
	drive.PredictedMediaLifeLeftPercent = t.SSDEnduranceUtilizationPercentage
	drive.Protocol = t.InterfaceType
	if t.FirmwareVersion.Current.VersionString != "" {
		drive.Revision = t.FirmwareVersion.Current.VersionString
	}
	if t.RotationalSpeedRpm > 0 {
		drive.RotationSpeedRPM = t.RotationalSpeedRpm
	}
	// This is a read/write object, so we need to save the raw object data for later
	drive.rawData = b

	return nil
}

// Update commits updates to this object's properties to the running system.
func (drive *DiskDrive) Update() error {
	// Get a representation of the object's original state so we can find what
	// to update.
	original := new(DiskDrive)
	err := original.UnmarshalJSON(drive.rawData)
	if err != nil {
		return err
	}

	readWriteFields := []string{
		"AssetTag",
		"HotspareReplacementMode",
		"IndicatorLED",
		"StatusIndicator",
		"WriteCacheEnabled",
	}

	originalElement := reflect.ValueOf(original).Elem()
	currentElement := reflect.ValueOf(drive).Elem()

	return drive.Entity.Update(originalElement, currentElement, readWriteFields)
}

// GetDrive will get a Drive instance from the service.
func GetDiskDrive(c common.Client, uri string) (*DiskDrive, error) {
	resp, err := c.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var drive DiskDrive
	err = json.NewDecoder(resp.Body).Decode(&drive)
	if err != nil {
		return nil, err
	}

	drive.SetClient(c)
	return &drive, nil
}

// ListReferencedDrives gets the collection of Drives from a provided reference.
func ListReferencedDiskDrives(c common.Client, link string) ([]*DiskDrive, error) { //nolint:dupl
	var result []*DiskDrive
	if link == "" {
		return result, nil
	}

	links, err := common.GetCollection(c, link)
	if err != nil {
		return result, err
	}

	collectionError := common.NewCollectionError()
	for _, driveLink := range links.ItemLinks {
		drive, err := GetDiskDrive(c, driveLink)
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
