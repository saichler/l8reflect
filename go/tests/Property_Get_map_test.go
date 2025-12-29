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
// Map Attribute Tests (Physicals, Logicals)
// ============================================================================

func Test_Get_Map_Physicals(t *testing.T) {
	r, device, _, _, _ := createElems()

	value, ok := getProperty("networkdevice.physicals", device, r, t)
	if !ok {
		return
	}

	physicals, ok := value.(map[string]*types.Physical)
	if !ok {
		r.Logger().Fail(t, "Expected map[string]*Physical type")
		return
	}
	if len(physicals) != len(device.Physicals) {
		r.Logger().Fail(t, "Expected Physicals map length to match")
		return
	}
}

func Test_Get_Map_Logicals(t *testing.T) {
	r, device, _, _, _ := createElems()

	value, ok := getProperty("networkdevice.logicals", device, r, t)
	if !ok {
		return
	}

	logicals, ok := value.(map[string]*types.Logical)
	if !ok {
		r.Logger().Fail(t, "Expected map[string]*Logical type")
		return
	}
	if len(logicals) != len(device.Logicals) {
		r.Logger().Fail(t, "Expected Logicals map length to match")
		return
	}
}

func Test_Get_Map_Entry_By_Key(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	physicals, ok := value.(map[string]*types.Physical)
	if !ok {
		r.Logger().Fail(t, "Expected map[string]*Physical type for map with key")
		return
	}
	if len(physicals) != len(device.Physicals) {
		r.Logger().Fail(t, "Expected Physicals map to match")
		return
	}
}

func Test_Get_Map_Entry_Nested_Field(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.id"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	strValue, ok := value.(string)
	if !ok {
		r.Logger().Fail(t, "Expected string type for Physical.Id")
		return
	}
	if strValue != device.Physicals[physicalKey].Id {
		r.Logger().Fail(t, "Expected Physical.Id to match device.Physicals[key].Id")
		return
	}
}

func Test_Get_Map_NonExistent_Key(t *testing.T) {
	r, device, _, _, _ := createElems()

	propertyId := "networkdevice.physicals<{24}nonexistent-key>.id"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	values := prop.GetValue(reflect.ValueOf(device))
	if len(values) != 0 {
		r.Logger().Fail(t, "Expected empty values for non-existent map key")
		return
	}
}

func Test_Get_Map_With_Nil_Entry(t *testing.T) {
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

	values := prop.GetValue(reflect.ValueOf(device))
	if len(values) != 0 {
		r.Logger().Fail(t, "Expected empty values for nil map entry")
		return
	}
}

func Test_Get_Map_All_Values_Field(t *testing.T) {
	r, device, _, _, _ := createElems()

	propertyId := "networkdevice.physicals.id"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	values := prop.GetValue(reflect.ValueOf(device))
	if len(values) != len(device.Physicals) {
		r.Logger().Fail(t, "Expected values count to match physicals count")
		return
	}
}

func Test_Get_All_Physical_Ids(t *testing.T) {
	r, device, _, _, _ := createElems()

	propertyId := "networkdevice.physicals.id"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	values, ok := value.([]interface{})
	if !ok {
		_, ok = value.(string)
		if !ok {
			r.Logger().Fail(t, "Expected []interface{} or string type for all Physical Ids")
			return
		}
		return
	}

	if len(values) != len(device.Physicals) {
		r.Logger().Fail(t, "Expected number of Physical Ids to match map size")
		return
	}
}

func Test_Get_Multiple_Values(t *testing.T) {
	r, device, _, _, _ := createElems()

	if len(device.Physicals) < 1 {
		r.Logger().Fail(t, "Need at least 1 physical for this test")
		return
	}

	propertyId := "networkdevice.physicals.id"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	if len(device.Physicals) > 1 {
		values, ok := value.([]interface{})
		if !ok {
			_, ok = value.(string)
			if !ok {
				r.Logger().Fail(t, "Expected []interface{} or string for multiple physical IDs")
				return
			}
		} else {
			if len(values) != len(device.Physicals) {
				r.Logger().Fail(t, "Expected values count to match physicals count")
				return
			}
		}
	}
}

func Test_GetValue_Map(t *testing.T) {
	r, device, _, _, _ := createElems()

	prop, err := properties.PropertyOf("networkdevice.physicals", r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	values := prop.GetAsValues(device)
	if len(values) != len(device.Physicals) {
		r.Logger().Fail(t, "Expected GetAsValues to return all map entries")
		return
	}
}

func Test_GetAsValues_Map_Interface_Item(t *testing.T) {
	r, device, _, _, _ := createElems()

	propertyId := "networkdevice.physicals"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	values := prop.GetAsValues(device)
	if len(values) != len(device.Physicals) {
		r.Logger().Fail(t, "Expected GetAsValues to return all map entries")
		return
	}
}
