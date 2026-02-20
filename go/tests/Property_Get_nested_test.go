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
	"testing"

	l8api "github.com/saichler/l8types/go/types/l8api"
	"github.com/saichler/probler/go/types"
)

// ============================================================================
// Deeply Nested Value Tests
// ============================================================================

func Test_Get_Deeply_Nested_Chassis_Id(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	if len(device.Physicals[physicalKey].Chassis) == 0 {
		r.Logger().Fail(t, "No chassis found in Physical for testing")
		return
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.chassis<{2}0>.id"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	strValue, ok := value.(string)
	if !ok {
		r.Logger().Fail(t, "Expected string type for Chassis.Id")
		return
	}
	if strValue != device.Physicals[physicalKey].Chassis[0].Id {
		r.Logger().Fail(t, "Expected Chassis.Id to match device.Physicals[key].Chassis[0].Id")
		return
	}
}

func Test_Get_Deeply_Nested_Chassis_SerialNumber(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	if len(device.Physicals[physicalKey].Chassis) == 0 {
		r.Logger().Fail(t, "No chassis found in Physical for testing")
		return
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.chassis<{2}0>.serialnumber"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	strValue, ok := value.(string)
	if !ok {
		r.Logger().Fail(t, "Expected string type for Chassis.SerialNumber")
		return
	}
	if strValue != device.Physicals[physicalKey].Chassis[0].SerialNumber {
		r.Logger().Fail(t, "Expected Chassis.SerialNumber to match")
		return
	}
}

func Test_Get_Deeply_Nested_Chassis_Model(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	if len(device.Physicals[physicalKey].Chassis) == 0 {
		r.Logger().Fail(t, "No chassis found in Physical for testing")
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

func Test_Get_Deeply_Nested_Chassis_Temperature(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	if len(device.Physicals[physicalKey].Chassis) == 0 {
		r.Logger().Fail(t, "No chassis found in Physical for testing")
		return
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.chassis<{2}0>.temperature"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	tsValue, ok := value.([]*l8api.L8TimeSeriesPoint)
	if !ok {
		r.Logger().Fail(t, "Expected []*L8TimeSeriesPoint type for Chassis.Temperature")
		return
	}
	expected := device.Physicals[physicalKey].Chassis[0].Temperature
	if len(tsValue) != len(expected) {
		r.Logger().Fail(t, "Expected Chassis.Temperature length to match")
		return
	}
	if len(tsValue) > 0 && tsValue[0].Value != expected[0].Value {
		r.Logger().Fail(t, "Expected Chassis.Temperature value to match")
		return
	}
}

// ============================================================================
// Performance Metrics (Nested Struct within Map Entry)
// ============================================================================

func Test_Get_Performance_Metrics(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	if device.Physicals[physicalKey].Performance == nil {
		r.Logger().Fail(t, "No Performance found in Physical for testing")
		return
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.performance"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	perf, ok := value.(*types.PerformanceMetrics)
	if !ok {
		r.Logger().Fail(t, "Expected *PerformanceMetrics type")
		return
	}
	if len(perf.CpuUsagePercent) != len(device.Physicals[physicalKey].Performance.CpuUsagePercent) {
		r.Logger().Fail(t, "Expected CpuUsagePercent length to match")
		return
	}
	if len(perf.CpuUsagePercent) > 0 && perf.CpuUsagePercent[0].Value != device.Physicals[physicalKey].Performance.CpuUsagePercent[0].Value {
		r.Logger().Fail(t, "Expected CpuUsagePercent value to match")
		return
	}
}

func Test_Get_Performance_CpuUsagePercent(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	if device.Physicals[physicalKey].Performance == nil {
		r.Logger().Fail(t, "No Performance found in Physical for testing")
		return
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.performance.cpuusagepercent"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	tsValue, ok := value.([]*l8api.L8TimeSeriesPoint)
	if !ok {
		r.Logger().Fail(t, "Expected []*L8TimeSeriesPoint type for CpuUsagePercent")
		return
	}
	expected := device.Physicals[physicalKey].Performance.CpuUsagePercent
	if len(tsValue) != len(expected) {
		r.Logger().Fail(t, "Expected CpuUsagePercent length to match")
		return
	}
	if len(tsValue) > 0 && tsValue[0].Value != expected[0].Value {
		r.Logger().Fail(t, "Expected CpuUsagePercent value to match")
		return
	}
}

func Test_Get_Performance_MemoryUsagePercent(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	if device.Physicals[physicalKey].Performance == nil {
		r.Logger().Fail(t, "No Performance found in Physical for testing")
		return
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.performance.memoryusagepercent"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	tsValue, ok := value.([]*l8api.L8TimeSeriesPoint)
	if !ok {
		r.Logger().Fail(t, "Expected []*L8TimeSeriesPoint type for MemoryUsagePercent")
		return
	}
	expected := device.Physicals[physicalKey].Performance.MemoryUsagePercent
	if len(tsValue) != len(expected) {
		r.Logger().Fail(t, "Expected MemoryUsagePercent length to match")
		return
	}
	if len(tsValue) > 0 && tsValue[0].Value != expected[0].Value {
		r.Logger().Fail(t, "Expected MemoryUsagePercent value to match")
		return
	}
}

// ============================================================================
// Port and Interface Tests (Deeper Nesting)
// ============================================================================

func Test_Get_Port_From_Chassis(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	if len(device.Physicals[physicalKey].Chassis) == 0 {
		r.Logger().Fail(t, "No chassis found in Physical for testing")
		return
	}
	if len(device.Physicals[physicalKey].Chassis[0].Ports) == 0 {
		r.Logger().Fail(t, "No ports found in Chassis for testing")
		return
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.chassis<{2}0>.ports"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	ports, ok := value.([]*types.Port)
	if !ok {
		r.Logger().Fail(t, "Expected []*Port type")
		return
	}
	if len(ports) != len(device.Physicals[physicalKey].Chassis[0].Ports) {
		r.Logger().Fail(t, "Expected Ports slice length to match")
		return
	}
}

func Test_Get_Port_Id_From_Chassis(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	if len(device.Physicals[physicalKey].Chassis) == 0 {
		r.Logger().Fail(t, "No chassis found in Physical for testing")
		return
	}
	if len(device.Physicals[physicalKey].Chassis[0].Ports) == 0 {
		r.Logger().Fail(t, "No ports found in Chassis for testing")
		return
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.chassis<{2}0>.ports<{2}0>.id"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	strValue, ok := value.(string)
	if !ok {
		r.Logger().Fail(t, "Expected string type for Port.Id")
		return
	}
	if strValue != device.Physicals[physicalKey].Chassis[0].Ports[0].Id {
		r.Logger().Fail(t, "Expected Port.Id to match")
		return
	}
}

func Test_Get_Interface_From_Port(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	if len(device.Physicals[physicalKey].Chassis) == 0 {
		r.Logger().Fail(t, "No chassis found in Physical for testing")
		return
	}
	if len(device.Physicals[physicalKey].Chassis[0].Ports) == 0 {
		r.Logger().Fail(t, "No ports found in Chassis for testing")
		return
	}
	if len(device.Physicals[physicalKey].Chassis[0].Ports[0].Interfaces) == 0 {
		r.Logger().Fail(t, "No interfaces found in Port for testing")
		return
	}

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.chassis<{2}0>.ports<{2}0>.interfaces<{2}0>.name"
	value, ok := getProperty(propertyId, device, r, t)
	if !ok {
		return
	}

	strValue, ok := value.(string)
	if !ok {
		r.Logger().Fail(t, "Expected string type for Interface.Name")
		return
	}
	if strValue != device.Physicals[physicalKey].Chassis[0].Ports[0].Interfaces[0].Name {
		r.Logger().Fail(t, "Expected Interface.Name to match")
		return
	}
}
