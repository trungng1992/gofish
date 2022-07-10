package redfish

import "github.com/trungng1992/gofish/common"

type CSOem struct {
	Hpe CSHpe
	Hp  CSHpe
}

type CSHpe struct {
	Links HpeLink // ilo5
}

type HpeLink struct {
	SmartStorage       common.Link
	PCIDevices         common.Link
	PCISlots           common.Link
	NetworkAdapter     common.Link
	USBPorts           common.Link
	USBDevices         common.Link
	EthernetInterfaces common.Link
	Memory             common.Link
}

type PSOem struct {
	Hpe PSHpe
	Hp  PSHpe
}

type PSHpe struct {
	// ODataContext is the odata context.
	ODataContext string `json:"@odata.context"`
	// ODataType is the odata type.
	ODataType               string `json:"@odata.type"`
	AveragePowerOutputWatts float64
	BayNumber               int
	HotPluggable            bool
	MaxPowerOutputWatts     float64
}
