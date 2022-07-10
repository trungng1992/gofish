package redfish

import (
	"encoding/json"
	"fmt"

	"github.com/trungng1992/gofish/common"
)

// SmartStorage is used to represent resources that represent a smart storage
// subsystem in the Redfish specification.
type SmartStorage struct {
	common.Entity

	// ODataContext is the odata context.
	ODataContext string `json:"@odata.context"`
	// ODataType is the odata type.
	ODataType string `json:"@odata.type"`
	// Description provides a description of this resource.
	Description string
	// ArrayController
	arraycontroller string
	Status          common.Status
}

// UnmarshalJSON unmarshals a Storage object from the raw JSON.
func (storage *SmartStorage) UnmarshalJSON(b []byte) error {
	type temp SmartStorage
	type mylinks struct {
		HostBusAdapters  common.Link
		ArrayControllers common.Link
	}

	var t struct {
		temp
		Links   mylinks
		Members common.Links
	}

	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}

	*storage = SmartStorage(t.temp)
	storage.arraycontroller = string(t.Links.ArrayControllers)

	if len(t.Members.ToStrings()) > 0 {
		storage.arraycontroller = t.Members.ToStrings()[0]
	}

	return nil
}

// GetSmartStorage will get a Storage instance from the service.
func GetSmartStorage(c common.Client, uri string) (*SmartStorage, error) {
	resp, err := c.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var smartStorage SmartStorage
	err = json.NewDecoder(resp.Body).Decode(&smartStorage)
	if err != nil {
		return nil, err
	}

	smartStorage.SetClient(c)
	return &smartStorage, nil
}

// ListReferencedSmartStorages gets the collection of Storage from a provided
// reference.
func ListReferencedSmartStorages(c common.Client, link string) (*SmartStorage, error) { //nolint:dupl
	var result *SmartStorage
	if link == "" {
		return result, nil
	}
	collectionError := common.NewCollectionError()

	smartStorage, err := GetSmartStorage(c, link)
	if err != nil {
		collectionError.Failures[link] = err
	}

	if collectionError.Empty() {
		return smartStorage, nil
	}

	return smartStorage, collectionError
}

// ArrayController gets the Array attached to the storage controllers that this
// resource represents.
func (smartstorage *SmartStorage) ArrayControllers() ([]*ArrayController, error) {
	return ListReferencedArrayControllers(smartstorage.Client, smartstorage.arraycontroller)
}

func (smartstorage *SmartStorage) DebugTest() {
	fmt.Println(smartstorage.arraycontroller)
}
