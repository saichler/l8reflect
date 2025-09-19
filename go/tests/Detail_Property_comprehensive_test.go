package tests

import (
	"testing"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/probler/go/types"
	"github.com/saichler/l8reflect/go/reflect/properties"
	"github.com/saichler/l8reflect/go/reflect/updating"
)

// ============================================================================
// Basic String Attribute Tests
// ============================================================================

// Note: Test_NetworkDevice_Id_Set already exists in Detail_Property_test.go

// Helper function to update with multiple changes for complex object replacement
func updateElemsMultiple(updater *updating.Updater, aside, zside, yside *types.NetworkDevice, r ifs.IResources, t *testing.T) bool {
	err := updater.Update(aside, zside)
	if err != nil {
		r.Logger().Fail(t, err.Error())
		return false
	}
	changes := updater.Changes()
	if len(changes) == 0 {
		r.Logger().Fail(t, "Expected at least one change")
		return false
	}
	
	// Apply all changes
	for _, change := range changes {
		propertyId := change.PropertyId()
		value := change.NewValue()
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
	}
	return true
}

// ============================================================================
// Pointer Attribute Tests (Nested Structures)
// ============================================================================

func Test_NetworkDevice_Equipmentinfo_Set_NonNil(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	
	// Create new EquipmentInfo with different values
	newEquipment := &types.EquipmentInfo{
		Vendor:       "Modified Cisco",
		Series:       "Modified Series",
		Family:       "Modified Family",
		Software:     "Modified Software",
		Hardware:     "Modified Hardware",
		Version:      "Modified Version",
		Model:        "Modified Model",
		SerialNumber: "Modified-SN-123",
	}
	zside.Equipmentinfo = newEquipment
	
	if !updateElemsMultiple(updater, aside, zside, yside, r, t) {
		return
	}
	
	if yside.Equipmentinfo == nil {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo to be non-nil")
		return
	}
	if yside.Equipmentinfo.Vendor != zside.Equipmentinfo.Vendor {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo.Vendor to equal zside.Equipmentinfo.Vendor")
		return
	}
}

func Test_NetworkDevice_Equipmentinfo_Set_Nil(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	zside.Equipmentinfo = nil
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Equipmentinfo != nil {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo to be nil")
		return
	}
}

func Test_NetworkDevice_Topology_Set_NonNil(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	
	// Create new NetworkTopology with only basic fields to avoid nested slice issues
	// Note: Setting nested complex fields (nodes, edges, etc.) to nil causes property navigation issues
	newTopology := &types.NetworkTopology{
		TopologyId:  "modified-topology-id", 
		Name:        "Modified Topology",
		LastUpdated: "2023-12-01T10:00:00Z",
		// Preserve complex nested fields from original to avoid nil navigation issues
		TopologyType:    aside.Topology.TopologyType,
		Nodes:           aside.Topology.Nodes,
		Edges:           aside.Topology.Edges,
		GeographicBounds: aside.Topology.GeographicBounds,
		Statistics:      aside.Topology.Statistics,
		HealthStatus:    aside.Topology.HealthStatus,
	}
	zside.Topology = newTopology
	
	if !updateElemsMultiple(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Topology == nil {
		r.Logger().Fail(t, "Expected yside.Topology to be non-nil")
		return
	}
	if yside.Topology.TopologyId != zside.Topology.TopologyId {
		r.Logger().Fail(t, "Expected yside.Topology.TopologyId to equal zside.Topology.TopologyId")
		return
	}
}

func Test_NetworkDevice_Topology_Set_Nil(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	zside.Topology = nil
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Topology != nil {
		r.Logger().Fail(t, "Expected yside.Topology to be nil")
		return
	}
}

func Test_NetworkDevice_NetworkHealth_Set_NonNil(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	
	// Create new NetworkHealth with basic fields, preserving complex nested fields
	// Note: Setting nested complex fields (alerts, etc.) to nil causes property navigation issues
	newHealth := &types.NetworkHealth{
		OverallStatus:   types.HealthStatus_HEALTH_STATUS_DEGRADED,  // Changed to test modification
		TotalDevices:    120,  // Changed to test modification
		OnlineDevices:   95,
		OfflineDevices:  5,
		// Preserve complex nested fields from original to avoid nil navigation issues
		WarningDevices:           aside.NetworkHealth.WarningDevices,
		TotalLinks:               aside.NetworkHealth.TotalLinks,
		ActiveLinks:              aside.NetworkHealth.ActiveLinks,
		InactiveLinks:            aside.NetworkHealth.InactiveLinks,
		WarningLinks:             aside.NetworkHealth.WarningLinks,
		NetworkAvailabilityPercent: aside.NetworkHealth.NetworkAvailabilityPercent,
		Alerts:                   aside.NetworkHealth.Alerts,
		LastHealthCheck:          aside.NetworkHealth.LastHealthCheck,
	}
	zside.NetworkHealth = newHealth
	
	if !updateElemsMultiple(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.NetworkHealth == nil {
		r.Logger().Fail(t, "Expected yside.NetworkHealth to be non-nil")
		return
	}
	if yside.NetworkHealth.OverallStatus != zside.NetworkHealth.OverallStatus {
		r.Logger().Fail(t, "Expected yside.NetworkHealth.OverallStatus to equal zside.NetworkHealth.OverallStatus")
		return
	}
}

func Test_NetworkDevice_NetworkHealth_Set_Nil(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	zside.NetworkHealth = nil
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.NetworkHealth != nil {
		r.Logger().Fail(t, "Expected yside.NetworkHealth to be nil")
		return
	}
}

// ============================================================================
// Map Attribute Tests
// ============================================================================

func Test_NetworkDevice_Physicals_Set_NonNil(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	
	// Create new Physicals map
	newPhysical := &types.Physical{
		Id: "modified-physical-001",
		Chassis: []*types.Chassis{
			{
				Id:           "modified-chassis-001",
				SerialNumber: "MOD-12345",
				Model:        "Modified Model",
				Description:  "Modified Physical Device",
			},
		},
	}
	zside.Physicals = map[string]*types.Physical{
		"modified-key": newPhysical,
	}
	
	if !updateElemsMultiple(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Physicals == nil {
		r.Logger().Fail(t, "Expected yside.Physicals to be non-nil")
		return
	}
	if len(yside.Physicals) != len(zside.Physicals) {
		r.Logger().Fail(t, "Expected yside.Physicals length to equal zside.Physicals length")
		return
	}
	if yside.Physicals["modified-key"] == nil {
		r.Logger().Fail(t, "Expected yside.Physicals['modified-key'] to exist")
		return
	}
	if yside.Physicals["modified-key"].Id != zside.Physicals["modified-key"].Id {
		r.Logger().Fail(t, "Expected yside.Physicals['modified-key'].Id to equal zside.Physicals['modified-key'].Id")
		return
	}
}

func Test_NetworkDevice_Physicals_Set_Empty(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	zside.Physicals = map[string]*types.Physical{}
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Physicals == nil {
		r.Logger().Fail(t, "Expected yside.Physicals to be non-nil but empty")
		return
	}
	if len(yside.Physicals) != 0 {
		r.Logger().Fail(t, "Expected yside.Physicals to be empty")
		return
	}
}

func Test_NetworkDevice_Physicals_Set_Nil(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	zside.Physicals = nil
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Physicals != nil {
		r.Logger().Fail(t, "Expected yside.Physicals to be nil")
		return
	}
}

func Test_NetworkDevice_Logicals_Set_NonNil(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	
	// Create new Logicals map
	newLogical := &types.Logical{
		Id: "modified-logical-001",
	}
	zside.Logicals = map[string]*types.Logical{
		"modified-logical-key": newLogical,
	}
	
	if !updateElemsMultiple(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Logicals == nil {
		r.Logger().Fail(t, "Expected yside.Logicals to be non-nil")
		return
	}
	if len(yside.Logicals) != len(zside.Logicals) {
		r.Logger().Fail(t, "Expected yside.Logicals length to equal zside.Logicals length")
		return
	}
	if yside.Logicals["modified-logical-key"] == nil {
		r.Logger().Fail(t, "Expected yside.Logicals['modified-logical-key'] to exist")
		return
	}
	if yside.Logicals["modified-logical-key"].Id != zside.Logicals["modified-logical-key"].Id {
		r.Logger().Fail(t, "Expected yside.Logicals['modified-logical-key'].Id to equal zside.Logicals['modified-logical-key'].Id")
		return
	}
}

func Test_NetworkDevice_Logicals_Set_Empty(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	zside.Logicals = map[string]*types.Logical{}
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Logicals == nil {
		r.Logger().Fail(t, "Expected yside.Logicals to be non-nil but empty")
		return
	}
	if len(yside.Logicals) != 0 {
		r.Logger().Fail(t, "Expected yside.Logicals to be empty")
		return
	}
}

func Test_NetworkDevice_Logicals_Set_Nil(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	zside.Logicals = nil
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Logicals != nil {
		r.Logger().Fail(t, "Expected yside.Logicals to be nil")
		return
	}
}

// ============================================================================
// Slice Attribute Tests
// ============================================================================

func Test_NetworkDevice_NetworkLinks_Set_NonNil(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	
	// Create new NetworkLinks slice with same length as original (1 element)
	// Note: Property system cannot change slice length, only modify existing elements
	newLink1 := &types.NetworkLink{
		LinkId:   "modified-link-001",
		Name:     "Modified Link 1",
		FromNode: "node-001",
		ToNode:   "node-002",
	}
	zside.NetworkLinks = []*types.NetworkLink{newLink1}
	
	if !updateElemsMultiple(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.NetworkLinks == nil {
		r.Logger().Fail(t, "Expected yside.NetworkLinks to be non-nil")
		return
	}
	if len(yside.NetworkLinks) != len(zside.NetworkLinks) {
		r.Logger().Fail(t, "Expected yside.NetworkLinks length to equal zside.NetworkLinks length")
		return
	}
	if len(yside.NetworkLinks) < 1 {
		r.Logger().Fail(t, "Expected at least one NetworkLink")
		return
	}
	if yside.NetworkLinks[0].LinkId != zside.NetworkLinks[0].LinkId {
		r.Logger().Fail(t, "Expected yside.NetworkLinks[0].LinkId to equal zside.NetworkLinks[0].LinkId")
		return
	}
}

func Test_NetworkDevice_NetworkLinks_Set_Empty(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	zside.NetworkLinks = []*types.NetworkLink{}
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.NetworkLinks == nil {
		r.Logger().Fail(t, "Expected yside.NetworkLinks to be non-nil but empty")
		return
	}
	if len(yside.NetworkLinks) != 0 {
		r.Logger().Fail(t, "Expected yside.NetworkLinks to be empty")
		return
	}
}

func Test_NetworkDevice_NetworkLinks_Set_Nil(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	zside.NetworkLinks = nil
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.NetworkLinks != nil {
		r.Logger().Fail(t, "Expected yside.NetworkLinks to be nil")
		return
	}
}

// ============================================================================
// Nested Structure Attribute Tests (EquipmentInfo fields)
// ============================================================================

func Test_NetworkDevice_EquipmentInfo_Vendor_Set(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	if zside.Equipmentinfo == nil {
		zside.Equipmentinfo = &types.EquipmentInfo{}
	}
	zside.Equipmentinfo.Vendor = "Modified Vendor"
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Equipmentinfo == nil {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo to be non-nil")
		return
	}
	if yside.Equipmentinfo.Vendor != zside.Equipmentinfo.Vendor {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo.Vendor to equal zside.Equipmentinfo.Vendor")
		return
	}
}

func Test_NetworkDevice_EquipmentInfo_Series_Set(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	if zside.Equipmentinfo == nil {
		zside.Equipmentinfo = &types.EquipmentInfo{}
	}
	zside.Equipmentinfo.Series = "Modified Series"
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Equipmentinfo == nil {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo to be non-nil")
		return
	}
	if yside.Equipmentinfo.Series != zside.Equipmentinfo.Series {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo.Series to equal zside.Equipmentinfo.Series")
		return
	}
}

func Test_NetworkDevice_EquipmentInfo_Family_Set(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	if zside.Equipmentinfo == nil {
		zside.Equipmentinfo = &types.EquipmentInfo{}
	}
	zside.Equipmentinfo.Family = "Modified Family"
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Equipmentinfo == nil {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo to be non-nil")
		return
	}
	if yside.Equipmentinfo.Family != zside.Equipmentinfo.Family {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo.Family to equal zside.Equipmentinfo.Family")
		return
	}
}

func Test_NetworkDevice_EquipmentInfo_Software_Set(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	if zside.Equipmentinfo == nil {
		zside.Equipmentinfo = &types.EquipmentInfo{}
	}
	zside.Equipmentinfo.Software = "Modified Software v2.0"
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Equipmentinfo == nil {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo to be non-nil")
		return
	}
	if yside.Equipmentinfo.Software != zside.Equipmentinfo.Software {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo.Software to equal zside.Equipmentinfo.Software")
		return
	}
}

func Test_NetworkDevice_EquipmentInfo_Hardware_Set(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	if zside.Equipmentinfo == nil {
		zside.Equipmentinfo = &types.EquipmentInfo{}
	}
	zside.Equipmentinfo.Hardware = "Modified Hardware Rev C"
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Equipmentinfo == nil {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo to be non-nil")
		return
	}
	if yside.Equipmentinfo.Hardware != zside.Equipmentinfo.Hardware {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo.Hardware to equal zside.Equipmentinfo.Hardware")
		return
	}
}

func Test_NetworkDevice_EquipmentInfo_Version_Set(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	if zside.Equipmentinfo == nil {
		zside.Equipmentinfo = &types.EquipmentInfo{}
	}
	zside.Equipmentinfo.Version = "Modified Version 3.1.4"
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Equipmentinfo == nil {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo to be non-nil")
		return
	}
	if yside.Equipmentinfo.Version != zside.Equipmentinfo.Version {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo.Version to equal zside.Equipmentinfo.Version")
		return
	}
}

func Test_NetworkDevice_EquipmentInfo_Model_Set(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	if zside.Equipmentinfo == nil {
		zside.Equipmentinfo = &types.EquipmentInfo{}
	}
	zside.Equipmentinfo.Model = "Modified Model XYZ-2000"
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Equipmentinfo == nil {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo to be non-nil")
		return
	}
	if yside.Equipmentinfo.Model != zside.Equipmentinfo.Model {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo.Model to equal zside.Equipmentinfo.Model")
		return
	}
}

func Test_NetworkDevice_EquipmentInfo_SerialNumber_Set(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	if zside.Equipmentinfo == nil {
		zside.Equipmentinfo = &types.EquipmentInfo{}
	}
	zside.Equipmentinfo.SerialNumber = "MOD-SN-98765"
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Equipmentinfo == nil {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo to be non-nil")
		return
	}
	if yside.Equipmentinfo.SerialNumber != zside.Equipmentinfo.SerialNumber {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo.SerialNumber to equal zside.Equipmentinfo.SerialNumber")
		return
	}
}

func Test_NetworkDevice_EquipmentInfo_SysName_Set(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	if zside.Equipmentinfo == nil {
		zside.Equipmentinfo = &types.EquipmentInfo{}
	}
	zside.Equipmentinfo.SysName = "modified-sys-name"
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Equipmentinfo == nil {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo to be non-nil")
		return
	}
	if yside.Equipmentinfo.SysName != zside.Equipmentinfo.SysName {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo.SysName to equal zside.Equipmentinfo.SysName")
		return
	}
}

func Test_NetworkDevice_EquipmentInfo_SysOid_Set(t *testing.T) {
	r, aside, zside, yside, updater := createElems()
	if zside.Equipmentinfo == nil {
		zside.Equipmentinfo = &types.EquipmentInfo{}
	}
	zside.Equipmentinfo.SysOid = "1.3.6.1.4.1.9999.1.1.1"
	
	if !updateElems(updater, aside, zside, yside, r, t) {
		return
	}
	if yside.Equipmentinfo == nil {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo to be non-nil")
		return
	}
	if yside.Equipmentinfo.SysOid != zside.Equipmentinfo.SysOid {
		r.Logger().Fail(t, "Expected yside.Equipmentinfo.SysOid to equal zside.Equipmentinfo.SysOid")
		return
	}
}