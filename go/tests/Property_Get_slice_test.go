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
	"github.com/saichler/probler/go/types"
)

// ============================================================================
// Slice Attribute Tests (Chassis, Ports within Physical)
// ============================================================================

func Test_Get_Slice_Chassis(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.chassis"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	chassis, ok := value.([]*types.Chassis)
	if !ok {
		r.Logger().Fail(t, "Expected []*Chassis type")
		return
	}
	if len(chassis) != len(device.Physicals[physicalKey].Chassis) {
		r.Logger().Fail(t, "Expected Chassis slice length to match")
		return
	}
}

func Test_Get_Slice_Ports(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.ports"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	ports, ok := value.([]*types.Port)
	if !ok {
		r.Logger().Fail(t, "Expected []*Port type")
		return
	}
	if len(ports) != len(device.Physicals[physicalKey].Ports) {
		r.Logger().Fail(t, "Expected Ports slice length to match")
		return
	}
}

func Test_Get_Slice_All_Elements_Field(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.chassis.id"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	values := prop.GetValue(reflect.ValueOf(device))
	expectedLen := len(device.Physicals[physicalKey].Chassis)
	if len(values) != expectedLen {
		r.Logger().Fail(t, "Expected values count to match chassis count")
		return
	}
}

func Test_Get_Slice_With_Nil_Element(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	device.Physicals[physicalKey].Chassis = append(device.Physicals[physicalKey].Chassis, nil)
	originalLen := len(device.Physicals[physicalKey].Chassis) - 1

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.chassis.id"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	values := prop.GetValue(reflect.ValueOf(device))
	if len(values) != originalLen {
		r.Logger().Fail(t, "Expected nil elements to be skipped")
		return
	}
}

func Test_Get_Slice_Index_Pointer_Type(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	if len(device.Physicals[physicalKey].Chassis) == 0 {
		r.Logger().Fail(t, "No chassis for testing")
		return
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.chassis<{2}0>.model"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	strValue, ok := value.(string)
	if !ok {
		r.Logger().Fail(t, "Expected string type for Chassis.Model")
		return
	}
	if strValue != device.Physicals[physicalKey].Chassis[0].Model {
		r.Logger().Fail(t, "Expected Chassis.Model to match")
		return
	}
}

func Test_Get_Invalid_First_Value(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	device.Physicals[physicalKey].Chassis = []*types.Chassis{}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.chassis.id"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	values := prop.GetValue(reflect.ValueOf(device))
	if len(values) != 0 {
		r.Logger().Fail(t, "Expected empty values for empty slice")
		return
	}
}

func Test_GetAsValues_Slice(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.chassis"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	values := prop.GetAsValues(device)
	expectedLen := len(device.Physicals[physicalKey].Chassis)
	if len(values) != expectedLen {
		r.Logger().Fail(t, "Expected GetAsValues to return all slice elements")
		return
	}
}

func Test_GetAsValues_Slice_Interface_Item(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.ports"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	values := prop.GetAsValues(device)
	expectedLen := len(device.Physicals[physicalKey].Ports)
	if len(values) != expectedLen {
		r.Logger().Fail(t, "Expected GetAsValues to return all ports")
		return
	}
}

// ============================================================================
// PowerSupply and Fan Tests (Additional Slices)
// ============================================================================

func Test_Get_PowerSupplies(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.powersupplies"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	supplies, ok := value.([]*types.PowerSupply)
	if !ok {
		r.Logger().Fail(t, "Expected []*PowerSupply type")
		return
	}
	if len(supplies) != len(device.Physicals[physicalKey].PowerSupplies) {
		r.Logger().Fail(t, "Expected PowerSupplies slice length to match")
		return
	}
}

func Test_Get_Fans(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.fans"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	fans, ok := value.([]*types.Fan)
	if !ok {
		r.Logger().Fail(t, "Expected []*Fan type")
		return
	}
	if len(fans) != len(device.Physicals[physicalKey].Fans) {
		r.Logger().Fail(t, "Expected Fans slice length to match")
		return
	}
}
