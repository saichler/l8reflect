package tests

import (
	"testing"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/probler/go/tests"
	"github.com/saichler/probler/go/types"
	"github.com/saichler/l8reflect/go/reflect/cloning"
	"github.com/saichler/l8reflect/go/reflect/properties"
	"github.com/saichler/l8reflect/go/reflect/updating"
)

func createElems() (ifs.IResources, *types.NetworkDevice, *types.NetworkDevice, *types.NetworkDevice, *updating.Updater) {
	r := newResources()
	r.Introspector().Inspect(&types.NetworkDevice{})
	deviceList := tests.GenerateExactDeviceTableMockData()
	c := cloning.NewCloner()
	for _, device := range deviceList.List {
		if device.Equipmentinfo.Model == "Cisco ASR 9000" {
			device2 := c.Clone(device).(*types.NetworkDevice)
			device3 := c.Clone(device).(*types.NetworkDevice)
			return r, device, device2, device3, updating.NewUpdater(r, true, true)
		}
	}
	panic("No device found")
}

func updateElems(updater *updating.Updater, aside, zside, yside *types.NetworkDevice, r ifs.IResources, t *testing.T) bool {
	err := updater.Update(aside, zside)
	if err != nil {
		r.Logger().Fail(t, err.Error())
		return false
	}
	changes := updater.Changes()
	if len(changes) != 1 {
		r.Logger().Fail(t, "Expected 1 change")
		return false
	}
	propertyId := changes[0].PropertyId()
	value := changes[0].NewValue()
	instance, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, err.Error())
		return false
	}
	_, _, err = instance.Set(yside, value)
	if err != nil {
		r.Logger().Fail(t, err.Error())
		return false
	}
	return true
}

func Test_NetworkDevice_Id_Set(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	zside.Id = "other"
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Id != zside.Id {
		r.Logger().Fail(t, "Expected zside.Id to equal yside.Id")
		return
	}
}

func Test_NetworkDevice_Nested_map(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	physicalKey := ""
	for k, _ := range zside.Physicals {
		physicalKey = k
		break
	}
	zside.Physicals[physicalKey].Id = "other"
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Physicals[physicalKey].Id != zside.Physicals[physicalKey].Id {
		r.Logger().Fail(t, "Expected zside.physoicals.Id to equal yside.physicals.Id")
		return
	}
}

func Test_NetworkDevice_Nested_slice(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	zside.NetworkLinks[0].LinkStatus = types.LinkStatus_LINK_STATUS_INACTIVE
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.NetworkLinks[0].LinkStatus != zside.NetworkLinks[0].LinkStatus {
		r.Logger().Fail(t, "Expected zside.NetworkLinks[0].LinkStatus to equal yside.NetworkLinks[0].LinkStatus")
		return
	}
}
