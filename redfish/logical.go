//
// SPDX-License-Identifier: BSD-3-Clause
//

package redfish

import (
	"encoding/json"
	"fmt"

	"github.com/trungng1992/gofish/common"
)

// Volume is used to represent a volume, virtual disk, logical disk, LUN,
// or other logical storage for a Redfish implementation.
type Logical struct {
	common.Entity

	// ODataContext is the odata context.
	ODataContext string `json:"@odata.context"`
	// ODataType is the odata type.
	ODataType string `json:"@odata.type"`
	// Description provides a description of this resource.
	Description string
	// Status is
	Status common.Status
	// CapacityBytes shall contain the size in bytes of the associated volume.
	CapacityBytes int
	// VolumeType shall contain the type of the associated Volume.
	MediaType MediaType

	BlockSizeBytes  int
	StripeSizeBytes int
	Raid            RAIDType
	// DrivesCount is the number of associated drives.
	DrivesCount int
	// drives contains references to associated drives.
	datadrive string
}

// UnmarshalJSON unmarshals a Volume object from the raw JSON.
func (logical *Logical) UnmarshalJSON(b []byte) error {
	type temp Logical
	type links struct {
		DataDrives common.Link
	}
	var t struct {
		temp
		Links       links
		CapacityMiB int
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}

	*logical = Logical(t.temp)

	// Extract the links to other entities for later
	if t.CapacityMiB > 0 {
		logical.CapacityBytes = t.CapacityMiB * 1e6
	}
	logical.datadrive = string(t.Links.DataDrives)

	return nil
}

// GetLogical will get a Volume instance from the service.
func GetLogical(c common.Client, uri string) (*Logical, error) {
	resp, err := c.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var volume Logical
	err = json.NewDecoder(resp.Body).Decode(&volume)
	if err != nil {
		return nil, err
	}

	physicaldrive, _ := ListReferencedPhysicalDrive(c, volume.datadrive)
	if physicaldrive != nil {
		volume.DrivesCount = physicaldrive.DrivesCount
	}

	volume.SetClient(c)
	return &volume, nil
}

// ListReferencedVolumes gets the collection of Volumes from a provided reference.
func ListReferencedLogical(c common.Client, link string) ([]*Logical, error) { //nolint:dupl
	var result []*Logical
	if link == "" {
		return result, nil
	}

	links, err := common.GetCollection(c, link)
	if err != nil {
		return result, err
	}

	collectionError := common.NewCollectionError()
	for _, volumeLink := range links.ItemLinks {
		volume, err := GetLogical(c, volumeLink)
		if err != nil {
			collectionError.Failures[volumeLink] = err
		} else if volume == nil {
			collectionError.Failures[volumeLink] = fmt.Errorf("volume %s not found", volumeLink)
		} else {
			result = append(result, volume)
		}
	}

	if collectionError.Empty() {
		return result, nil
	}

	return result, collectionError
}

// Drives references the Drives that this volume is associated with.
func (logical *Logical) Drives() ([]*DiskDrive, error) {
	var result []*DiskDrive

	physicaldrive, err := ListReferencedPhysicalDrive(logical.Client, logical.datadrive)
	if err != nil {
		return result, err
	}

	return physicaldrive.Drives()
}
