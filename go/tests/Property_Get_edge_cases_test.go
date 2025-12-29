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
// Nil/Empty Value Edge Cases
// ============================================================================

func Test_Get_Nil_Equipmentinfo(t *testing.T) {
	r, device, _, _, _ := createElems()
	device.Equipmentinfo = nil

	if !getPropertyExpectNil("networkdevice.equipmentinfo", device, r, t) {
		return
	}
}

func Test_Get_Nested_Field_With_Nil_Parent(t *testing.T) {
	r, device, _, _, _ := createElems()
	device.Equipmentinfo = nil

	prop, err := properties.PropertyOf("networkdevice.equipmentinfo.vendor", r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	values := prop.GetValue(reflect.ValueOf(device))
	if len(values) != 0 {
		r.Logger().Fail(t, "Expected GetValue to return empty slice when parent is nil")
		return
	}
}

func Test_Get_Empty_Map(t *testing.T) {
	r, device, _, _, _ := createElems()
	device.Physicals = map[string]*types.Physical{}

	value, ok := getProperty("networkdevice.physicals", device, r, t)
	if !ok {
		return
	}

	physicals, ok := value.(map[string]*types.Physical)
	if !ok {
		r.Logger().Fail(t, "Expected map[string]*Physical type")
		return
	}
	if len(physicals) != 0 {
		r.Logger().Fail(t, "Expected empty Physicals map")
		return
	}
}

func Test_Get_Nil_Map(t *testing.T) {
	r, device, _, _, _ := createElems()
	device.Physicals = nil

	value, ok := getProperty("networkdevice.physicals", device, r, t)
	if !ok {
		return
	}

	physicals, ok := value.(map[string]*types.Physical)
	if !ok {
		r.Logger().Fail(t, "Expected map[string]*Physical type")
		return
	}
	if physicals != nil {
		r.Logger().Fail(t, "Expected nil Physicals map")
		return
	}
}

func Test_Get_Nil_Device(t *testing.T) {
	r, _, _, _, _ := createElems()

	prop, err := properties.PropertyOf("networkdevice<{24}{24}test-id>", r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	value, err := prop.Get(nil)
	if err != nil {
		r.Logger().Fail(t, "Failed to get value with nil device: "+err.Error())
		return
	}

	device, ok := value.(*types.NetworkDevice)
	if !ok {
		r.Logger().Fail(t, "Expected *NetworkDevice type when creating from nil")
		return
	}
	if device.Id != "test-id" {
		r.Logger().Fail(t, "Expected device.Id to be 'test-id', got: "+device.Id)
		return
	}
}

func Test_Get_Returns_Nil_Pointer(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	device.Physicals[physicalKey].Performance = nil

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.performance"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	if value != nil {
		r.Logger().Fail(t, "Expected nil for nil pointer field")
		return
	}
}

func Test_GetAsValues_With_Nil_Result(t *testing.T) {
	r, device, _, _, _ := createElems()
	device.Equipmentinfo = nil

	propertyId := "networkdevice.equipmentinfo"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	values := prop.GetAsValues(device)
	if len(values) != 1 {
		r.Logger().Fail(t, "Expected GetAsValues to return 1 value when result is nil")
		return
	}
}

func Test_Get_With_Primary_Key(t *testing.T) {
	r, _, _, _, _ := createElems()

	propertyId := "networkdevice<{24}{24}my-test-id>"
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}

	value, err := prop.Get(nil)
	if err != nil {
		r.Logger().Fail(t, "Failed to get value: "+err.Error())
		return
	}

	device, ok := value.(*types.NetworkDevice)
	if !ok {
		r.Logger().Fail(t, "Expected *NetworkDevice type")
		return
	}
	if device.Id != "my-test-id" {
		r.Logger().Fail(t, "Expected Id to be 'my-test-id', got: "+device.Id)
		return
	}
}
