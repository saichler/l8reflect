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

	"github.com/saichler/l8reflect/go/reflect/cloning"
	"github.com/saichler/l8reflect/go/reflect/properties"
	l8api "github.com/saichler/l8types/go/types/l8api"
	"github.com/saichler/probler/go/types"
)

func TestTimeSeriesAppend(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	// Clear existing time series data
	device.Physicals[physicalKey].Performance.CpuUsagePercent = nil

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.performance.cpuusagepercent"

	// Append first point
	point1 := &l8api.L8TimeSeriesPoint{Stamp: 1, Value: 10.0}
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}
	_, _, err = prop.Set(device, point1)
	if err != nil {
		r.Logger().Fail(t, "Failed to set point1: "+err.Error())
		return
	}

	ts := device.Physicals[physicalKey].Performance.CpuUsagePercent
	if len(ts) != 1 {
		r.Logger().Fail(t, "Expected 1 point after first append")
		return
	}
	if ts[0].Value != 10.0 {
		r.Logger().Fail(t, "Expected first point value to be 10.0")
		return
	}

	// Append second point
	point2 := &l8api.L8TimeSeriesPoint{Stamp: 2, Value: 20.0}
	prop2, _ := properties.PropertyOf(propertyId, r)
	_, _, err = prop2.Set(device, point2)
	if err != nil {
		r.Logger().Fail(t, "Failed to set point2: "+err.Error())
		return
	}

	ts = device.Physicals[physicalKey].Performance.CpuUsagePercent
	if len(ts) != 2 {
		r.Logger().Fail(t, "Expected 2 points after second append")
		return
	}
	if ts[1].Value != 20.0 {
		r.Logger().Fail(t, "Expected second point value to be 20.0")
		return
	}
}

func TestTimeSeriesCapAt100(t *testing.T) {
	r, device, _, _, _ := createElems()

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	// Pre-fill with 100 points
	points := make([]*l8api.L8TimeSeriesPoint, 100)
	for i := 0; i < 100; i++ {
		points[i] = &l8api.L8TimeSeriesPoint{Stamp: int64(i), Value: float64(i)}
	}
	device.Physicals[physicalKey].Performance.CpuUsagePercent = points

	propertyId := "networkdevice.physicals<{24}" + physicalKey + ">.performance.cpuusagepercent"

	// Append one more point, should drop index 0
	newPoint := &l8api.L8TimeSeriesPoint{Stamp: 100, Value: 100.0}
	prop, err := properties.PropertyOf(propertyId, r)
	if err != nil {
		r.Logger().Fail(t, "Failed to create property: "+err.Error())
		return
	}
	_, _, err = prop.Set(device, newPoint)
	if err != nil {
		r.Logger().Fail(t, "Failed to set new point: "+err.Error())
		return
	}

	ts := device.Physicals[physicalKey].Performance.CpuUsagePercent
	if len(ts) != 100 {
		r.Logger().Fail(t, "Expected 100 points after cap")
		return
	}
	// Oldest should now be stamp=1 (stamp=0 was dropped)
	if ts[0].Stamp != 1 {
		r.Logger().Fail(t, "Expected oldest point stamp to be 1 after drop")
		return
	}
	// Newest should be the appended point
	if ts[99].Value != 100.0 {
		r.Logger().Fail(t, "Expected newest point value to be 100.0")
		return
	}
}

func TestTimeSeriesUpdater(t *testing.T) {
	r, device, _, _, updater := createElems()
	c := cloning.NewCloner()
	device2 := c.Clone(device)

	var physicalKey string
	for k := range device.Physicals {
		physicalKey = k
		break
	}

	// Clear time series on old device
	device.Physicals[physicalKey].Performance.CpuUsagePercent = nil

	// Set time series on new device
	device2.(*types.NetworkDevice).Physicals[physicalKey].Performance.CpuUsagePercent = []*l8api.L8TimeSeriesPoint{
		{Stamp: 1, Value: 10.0},
		{Stamp: 2, Value: 20.0},
	}

	err := updater.Update(device, device2)
	if err != nil {
		r.Logger().Fail(t, "Update failed: "+err.Error())
		return
	}

	changes := updater.Changes()
	if len(changes) == 0 {
		r.Logger().Fail(t, "Expected at least one change")
		return
	}

	// Verify no change drills into individual L8TimeSeriesPoint fields.
	// All time series changes should be at the slice level.
	for _, change := range changes {
		pid := change.PropertyId()
		// The change should be for the slice itself (e.g. ...cpuusagepercent),
		// not for fields inside a point (e.g. ...cpuusagepercent<0>.stamp)
		newVal := change.NewValue()
		if _, ok := newVal.(int64); ok {
			r.Logger().Fail(t, "Updater drilled into time series point scalar field: "+pid)
			return
		}
		if _, ok := newVal.(float64); ok {
			r.Logger().Fail(t, "Updater drilled into time series point scalar field: "+pid)
			return
		}
	}

	// Apply changes to a third device and verify
	device3 := c.Clone(device2).(*types.NetworkDevice)
	device3.Physicals[physicalKey].Performance.CpuUsagePercent = nil

	for _, change := range changes {
		pid := change.PropertyId()
		val := change.NewValue()
		prop, err := properties.PropertyOf(pid, r)
		if err != nil {
			r.Logger().Fail(t, "Failed to create property for apply: "+err.Error())
			return
		}
		_, _, err = prop.Set(device3, val)
		if err != nil {
			r.Logger().Fail(t, "Failed to apply change: "+err.Error())
			return
		}
	}

	ts := device3.Physicals[physicalKey].Performance.CpuUsagePercent
	if len(ts) < 2 {
		r.Logger().Fail(t, "Expected at least 2 points after applying changes, got ", len(ts))
		return
	}
}
