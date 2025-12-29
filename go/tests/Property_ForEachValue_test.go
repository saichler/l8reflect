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
// ForEachValue Basic Tests
// ============================================================================

func Test_ForEachValue_Simple_String(t *testing.T) {
	r, device, _, _, _ := createElems()

	prop, err := properties.PropertyOf("networkdevice.id", r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	var foundValue string
	count := 0
	prop.ForEachValue(reflect.ValueOf(device), func(v reflect.Value) bool {
		foundValue = v.String()
		count++
		return true
	})

	if count != 1 {
		r.Logger().Fail(t, "Expected 1 value from ForEachValue")
		return
	}
	if foundValue != device.Id {
		r.Logger().Fail(t, "Expected value to equal device.Id")
		return
	}
}

func Test_ForEachValue_Root_Property(t *testing.T) {
	r, device, _, _, _ := createElems()

	prop, err := properties.PropertyOf("networkdevice", r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	count := 0
	prop.ForEachValue(reflect.ValueOf(device), func(v reflect.Value) bool {
		count++
		return true
	})

	if count != 1 {
		r.Logger().Fail(t, "Expected 1 value for root property")
		return
	}
}

func Test_ForEachValue_Invalid_Value(t *testing.T) {
	r, _, _, _, _ := createElems()

	prop, err := properties.PropertyOf("networkdevice.id", r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	count := 0
	prop.ForEachValue(reflect.Value{}, func(v reflect.Value) bool {
		count++
		return true
	})

	if count != 0 {
		r.Logger().Fail(t, "Expected 0 values for invalid reflect.Value")
		return
	}
}

func Test_ForEachValue_Nil_Pointer(t *testing.T) {
	r, _, _, _, _ := createElems()

	prop, err := properties.PropertyOf("networkdevice.id", r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	var device *types.NetworkDevice = nil
	count := 0
	prop.ForEachValue(reflect.ValueOf(device), func(v reflect.Value) bool {
		count++
		return true
	})

	if count != 0 {
		r.Logger().Fail(t, "Expected 0 values for nil pointer")
		return
	}
}

// ============================================================================
// ForEachValue Map Tests
// ============================================================================

func Test_ForEachValue_Map_All_Entries(t *testing.T) {
	r, device, _, _, _ := createElems()

	prop, err := properties.PropertyOf("networkdevice.physicals.id", r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	count := 0
	prop.ForEachValue(reflect.ValueOf(device), func(v reflect.Value) bool {
		count++
		return true
	})

	if count != len(device.Physicals) {
		r.Logger().Fail(t, "Expected count to match physicals count")
		return
	}
}

func Test_ForEachValue_Map_With_Key(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.id"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	var foundValue string
	count := 0
	prop.ForEachValue(reflect.ValueOf(device), func(v reflect.Value) bool {
		foundValue = v.String()
		count++
		return true
	})

	if count != 1 {
		r.Logger().Fail(t, "Expected 1 value for map with key")
		return
	}
	if foundValue != device.Physicals[physicalKey].Id {
		r.Logger().Fail(t, "Expected value to match Physical.Id")
		return
	}
}

func Test_ForEachValue_Map_NonExistent_Key(t *testing.T) {
	r, device, _, _, _ := createElems()

	propertyId := "networkdevice.physicals<{24}nonexistent-key>.id"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	count := 0
	prop.ForEachValue(reflect.ValueOf(device), func(v reflect.Value) bool {
		count++
		return true
	})

	if count != 0 {
		r.Logger().Fail(t, "Expected 0 values for non-existent map key")
		return
	}
}

func Test_ForEachValue_Map_Nil_Entry(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}
	device.Physicals[physicalKey] = nil

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.id"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	count := 0
	prop.ForEachValue(reflect.ValueOf(device), func(v reflect.Value) bool {
		count++
		return true
	})

	if count != 0 {
		r.Logger().Fail(t, "Expected 0 values for nil map entry")
		return
	}
}

// ============================================================================
// ForEachValue Slice Tests
// ============================================================================

func Test_ForEachValue_Slice_All_Elements(t *testing.T) {
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

	count := 0
	prop.ForEachValue(reflect.ValueOf(device), func(v reflect.Value) bool {
		count++
		return true
	})

	expectedLen := len(device.Physicals[physicalKey].Chassis)
	if count != expectedLen {
		r.Logger().Fail(t, "Expected count to match chassis count")
		return
	}
}

func Test_ForEachValue_Slice_With_Index(t *testing.T) {
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

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.chassis<{2}0>.id"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	var foundValue string
	count := 0
	prop.ForEachValue(reflect.ValueOf(device), func(v reflect.Value) bool {
		foundValue = v.String()
		count++
		return true
	})

	if count != 1 {
		r.Logger().Fail(t, "Expected 1 value for slice with index")
		return
	}
	if foundValue != device.Physicals[physicalKey].Chassis[0].Id {
		r.Logger().Fail(t, "Expected value to match Chassis[0].Id")
		return
	}
}

func Test_ForEachValue_Slice_Nil_Element(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	originalLen := len(device.Physicals[physicalKey].Chassis)
	device.Physicals[physicalKey].Chassis = append(device.Physicals[physicalKey].Chassis, nil)

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.chassis.id"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	count := 0
	prop.ForEachValue(reflect.ValueOf(device), func(v reflect.Value) bool {
		count++
		return true
	})

	if count != originalLen {
		r.Logger().Fail(t, "Expected nil elements to be skipped")
		return
	}
}

// ============================================================================
// ForEachValue Early Termination Tests
// ============================================================================

func Test_ForEachValue_Early_Termination(t *testing.T) {
	r, device, _, _, _ := createElems()

	prop, err := properties.PropertyOf("networkdevice.physicals.id", r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	count := 0
	prop.ForEachValue(reflect.ValueOf(device), func(v reflect.Value) bool {
		count++
		return false // Stop after first value
	})

	if count != 1 {
		r.Logger().Fail(t, "Expected early termination after first value")
		return
	}
}

func Test_ForEachValue_Early_Termination_Slice(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	if len(device.Physicals[physicalKey].Chassis) < 2 {
		// Add another chassis for testing
		device.Physicals[physicalKey].Chassis = append(
			device.Physicals[physicalKey].Chassis,
			&types.Chassis{Id: "test-chassis-2"},
		)
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.chassis.id"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	count := 0
	prop.ForEachValue(reflect.ValueOf(device), func(v reflect.Value) bool {
		count++
		return false // Stop after first value
	})

	if count != 1 {
		r.Logger().Fail(t, "Expected early termination after first slice element")
		return
	}
}

// ============================================================================
// ForEachValue Comparison with GetValue
// ============================================================================

func Test_ForEachValue_Matches_GetValue(t *testing.T) {
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

	// Collect values using ForEachValue
	var forEachValues []string
	prop.ForEachValue(reflect.ValueOf(device), func(v reflect.Value) bool {
		forEachValues = append(forEachValues, v.String())
		return true
	})

	// Collect values using GetValue
	getValues := prop.GetValue(reflect.ValueOf(device))
	var getValueStrings []string
	for _, v := range getValues {
		getValueStrings = append(getValueStrings, v.String())
	}

	// Compare counts
	if len(forEachValues) != len(getValueStrings) {
		r.Logger().Fail(t, "ForEachValue and GetValue returned different counts")
		return
	}

	// Compare values (order may differ for maps, but for slices should match)
	for i, v := range forEachValues {
		if v != getValueStrings[i] {
			r.Logger().Fail(t, "ForEachValue and GetValue values don't match at index")
			return
		}
	}
}
