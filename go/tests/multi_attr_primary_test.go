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
	"github.com/saichler/l8reflect/go/reflect/updating"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/types/l8services"
)

func TestMultiAttrPrimary(t *testing.T) {

	serviceName := "Test2"
	serviceArea := byte(0)

	aside := &l8services.L8ReplicationIndex{}
	aside.ServiceName = serviceName
	aside.ServiceArea = int32(serviceArea)

	cloner := cloning.NewCloner()

	zside := cloner.Clone(aside).(*l8services.L8ReplicationIndex)
	zside.Keys = make(map[string]*l8services.L8ReplicationKey)
	zside.Keys["test"] = &l8services.L8ReplicationKey{}
	zside.Keys["test"].Location = make(map[string]int32)

	yside := cloner.Clone(aside).(*l8services.L8ReplicationIndex)

	patchUpdateIndex(aside, zside, yside, t)
}

func patchUpdateIndex(o, n, z *l8services.L8ReplicationIndex, t *testing.T) bool {
	res := newResources()
	res.Introspector().Decorators().AddPrimaryKeyDecorator(&l8services.L8ReplicationIndex{}, "ServiceName", "ServiceArea")

	u := updating.NewUpdater(res, false, false)
	err := u.Update(o, n)
	if err != nil {
		log.Fail(t, err.Error())
		return false
	}

	for _, c := range u.Changes() {
		pid := c.PropertyId()
		oObj := object.NewEncode()
		oObj.Add(c.OldValue())
		nObj := object.NewEncode()
		nObj.Add(c.NewValue())
		prop, err := properties.PropertyOf(pid, res)
		if err != nil {
			log.Fail(t, err.Error())
			return false
		}
		pObj := object.NewDecode(nObj.Data(), 0, res.Registry())
		v, _ := pObj.Get()
		_, _, err = prop.Set(z, v)
	}

	return true
}
