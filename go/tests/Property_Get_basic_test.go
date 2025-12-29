// © 2025 Sharon Aicler (saichler@gmail.com)
//
// Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"reflect"
	"testing"

	"github.com/saichler/l8reflect/go/reflect/properties"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/probler/go/types"
)

// ============================================================================
// Helper Functions
// ============================================================================

// getProperty is a helper function to get a property value and handle errors
func getProperty(propertyId string, device *types.NetworkDevice, r ifs.IResources, t *testing.T) (interface{}, bool) {
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property for '"+propertyId+"': "+err.Error())
		return nil, false
	}
	value, err := prop.Get(device)
	if err != nil {
		r.Logger().Fail(t, "Failed to get value for '"+propertyId+"': "+err.Error())
		return nil, false
	}
	return value, true
}

// getPropertyExpectNil is a helper for cases where we expect nil result
func getPropertyExpectNil(propertyId string, device *types.NetworkDevice, r ifs.IResources, t *testing.T) bool {
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property for '"+propertyId+"': "+err.Error())
		return false
	}
	value, err := prop.Get(device)
	if err != nil {
		r.Logger().Fail(t, "Failed to get value for '"+propertyId+"': "+err.Error())
		return false
	}
	if value != nil {
		r.Logger().Fail(t, "Expected nil value for '"+propertyId+"', got non-nil")
		return false
	}
	return true
}

// ============================================================================
// Simple Attribute Tests (Root Level)
// ============================================================================

func Test_Get_Simple_String_Id(t *testing.T) {
	r, device, _, _, _ := createElems()

	value, ok := getProperty("networkdevice.id", device, r, t)
	if !ok {
		return
	}

	strValue, ok := value.(string)
	if !ok {
		r.Logger().Fail(t, "Expected string type for Id")
		return
	}
	if strValue != device.Id {
		r.Logger().Fail(t, "Expected Id to equal device.Id")
		return
	}
}

// ============================================================================
// Nested Struct Attribute Tests (EquipmentInfo)
// ============================================================================

func Test_Get_Nested_Struct_Equipmentinfo(t *testing.T) {
	r, device, _, _, _ := createElems()

	value, ok := getProperty("networkdevice.equipmentinfo", device, r, t)
	if !ok {
		return
	}

	equipInfo, ok := value.(*types.EquipmentInfo)
	if !ok {
		r.Logger().Fail(t, "Expected *EquipmentInfo type")
		return
	}
	if equipInfo.Vendor != device.Equipmentinfo.Vendor {
		r.Logger().Fail(t, "Expected Equipmentinfo.Vendor to match")
		return
	}
}

func Test_Get_Nested_Struct_Equipmentinfo_Vendor(t *testing.T) {
	r, device, _, _, _ := createElems()

	value, ok := getProperty("networkdevice.equipmentinfo.vendor", device, r, t)
	if !ok {
		return
	}

	strValue, ok := value.(string)
	if !ok {
		r.Logger().Fail(t, "Expected string type for Vendor")
		return
	}
	if strValue != device.Equipmentinfo.Vendor {
		r.Logger().Fail(t, "Expected Vendor to equal device.Equipmentinfo.Vendor")
		return
	}
}

func Test_Get_Nested_Struct_Equipmentinfo_Model(t *testing.T) {
	r, device, _, _, _ := createElems()

	value, ok := getProperty("networkdevice.equipmentinfo.model", device, r, t)
	if !ok {
		return
	}

	strValue, ok := value.(string)
	if !ok {
		r.Logger().Fail(t, "Expected string type for Model")
		return
	}
	if strValue != device.Equipmentinfo.Model {
		r.Logger().Fail(t, "Expected Model to equal device.Equipmentinfo.Model")
		return
	}
}

func Test_Get_Nested_Struct_Equipmentinfo_SerialNumber(t *testing.T) {
	r, device, _, _, _ := createElems()

	value, ok := getProperty("networkdevice.equipmentinfo.serialnumber", device, r, t)
	if !ok {
		return
	}

	strValue, ok := value.(string)
	if !ok {
		r.Logger().Fail(t, "Expected string type for SerialNumber")
		return
	}
	if strValue != device.Equipmentinfo.SerialNumber {
		r.Logger().Fail(t, "Expected SerialNumber to equal device.Equipmentinfo.SerialNumber")
		return
	}
}

func Test_Get_Nested_Struct_Equipmentinfo_Latitude(t *testing.T) {
	r, device, _, _, _ := createElems()

	value, ok := getProperty("networkdevice.equipmentinfo.latitude", device, r, t)
	if !ok {
		return
	}

	floatValue, ok := value.(float64)
	if !ok {
		r.Logger().Fail(t, "Expected float64 type for Latitude")
		return
	}
	if floatValue != device.Equipmentinfo.Latitude {
		r.Logger().Fail(t, "Expected Latitude to equal device.Equipmentinfo.Latitude")
		return
	}
}

// ============================================================================
// Enum Type Tests
// ============================================================================

func Test_Get_Enum_DeviceType(t *testing.T) {
	r, device, _, _, _ := createElems()

	if device.Equipmentinfo == nil {
		r.Logger().Fail(t, "No EquipmentInfo for testing")
		return
	}

	propertyId := "networkdevice.equipmentinfo.devicetype"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	enumValue, ok := value.(types.DeviceType)
	if !ok {
		r.Logger().Fail(t, "Expected DeviceType enum type")
		return
	}
	if enumValue != device.Equipmentinfo.DeviceType {
		r.Logger().Fail(t, "Expected DeviceType to match")
		return
	}
}

// ============================================================================
// Invalid Property Path Tests
// ============================================================================

func Test_Get_Invalid_Property_Path(t *testing.T) {
	r, _, _, _, _ := createElems()

	_, err := properties.PropertyOf("networkdevice.nonexistent", r)
	if err == nil {
		r.Logger().Fail(t, "Expected error for invalid property path")
		return
	}
}

func Test_Get_Invalid_Nested_Property_Path(t *testing.T) {
	r, _, _, _, _ := createElems()

	_, err := properties.PropertyOf("networkdevice.equipmentinfo.nonexistent", r)
	if err == nil {
		r.Logger().Fail(t, "Expected error for invalid nested property path")
		return
	}
}

// ============================================================================
// GetValue Tests
// ============================================================================

func Test_GetValue_Simple_String(t *testing.T) {
	r, device, _, _, _ := createElems()

	prop, err := properties.PropertyOf("networkdevice.id", r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	values := prop.GetAsValues(device)
	if len(values) != 1 {
		r.Logger().Fail(t, "Expected 1 value from GetAsValues")
		return
	}

	if values[0].String() != device.Id {
		r.Logger().Fail(t, "Expected value to equal device.Id")
		return
	}
}

func Test_GetValue_With_Invalid_Value(t *testing.T) {
	r, _, _, _, _ := createElems()

	propertyId := "networkdevice.id"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	values := prop.GetValue(reflect.Value{})
	if len(values) != 0 {
		r.Logger().Fail(t, "Expected empty values for invalid reflect.Value")
		return
	}
}

func Test_GetValue_With_Nil_Pointer(t *testing.T) {
	r, _, _, _, _ := createElems()

	propertyId := "networkdevice.id"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	var device *types.NetworkDevice = nil
	values := prop.GetValue(reflect.ValueOf(device))
	if len(values) != 0 {
		r.Logger().Fail(t, "Expected empty values for nil pointer")
		return
	}
}

func Test_GetValue_Root_Property(t *testing.T) {
	r, device, _, _, _ := createElems()

	propertyId := "networkdevice"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	values := prop.GetValue(reflect.ValueOf(device))
	if len(values) != 1 {
		r.Logger().Fail(t, "Expected 1 value for root property")
		return
	}
}

func Test_GetValue_Pointer_Parent(t *testing.T) {
	r, device, _, _, _ := createElems()

	propertyId := "networkdevice.equipmentinfo.vendor"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	values := prop.GetValue(reflect.ValueOf(device))
	if len(values) != 1 {
		r.Logger().Fail(t, "Expected 1 value for equipmentinfo.vendor")
		return
	}

	if values[0].String() != device.Equipmentinfo.Vendor {
		r.Logger().Fail(t, "Expected vendor to match")
		return
	}
}
