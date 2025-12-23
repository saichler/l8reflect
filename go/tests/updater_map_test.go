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

	"github.com/saichler/l8reflect/go/reflect/updating"
	"github.com/saichler/l8reflect/go/tests/utils"
	"github.com/saichler/l8types/go/testtypes"
)

func patchUpdate(o, n *testtypes.TestProto, t *testing.T) bool {
	res := newResources()

	u := updating.NewUpdater(res, false, true)
	err := u.Update(o, n)
	if err != nil {
		log.Fail(t, err.Error())
		return false
	}
	return true
}

func checkPrimitive(o, n *testtypes.TestProto, t *testing.T) bool {
	if o.MyString2StringMap == nil {
		log.Fail(t, "Expected map to not be nil")
		return false
	}
	if len(o.MyString2StringMap) != len(n.MyString2StringMap) {
		log.Fail(t, "maps are not the same len")
		return false
	}
	for k, v := range n.MyString2StringMap {
		vo, ok := o.MyString2StringMap[k]
		if !ok {
			log.Fail(t, "Expected key to exist in old map")
			return false
		}
		if vo != v {
			log.Fail(t, "Expected values to match for key")
			return false
		}
	}
	return true
}

func TestMapPrimitiveSetFromNil(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	o.MyString2StringMap = nil

	if !patchUpdate(o, n, t) {
		return
	}

	if !checkPrimitive(o, n, t) {
		return
	}
}

func TestMapPrimitiveSetFromEmpty(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	o.MyString2StringMap = make(map[string]string)

	if !patchUpdate(o, n, t) {
		return
	}

	if !checkPrimitive(o, n, t) {
		return
	}
}

func TestMapPrimitiveChangeValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2StringMap {
		n.MyString2StringMap[k] = n.MyString2StringMap[k] + "C"
	}

	if !patchUpdate(o, n, t) {
		return
	}

	if !checkPrimitive(o, n, t) {
		return
	}
}

func TestMapPrimitiveAddValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	n.MyString2StringMap["new"] = "new"

	if !patchUpdate(o, n, t) {
		return
	}

	if !checkPrimitive(o, n, t) {
		return
	}
}

func TestMapPrimitiveDelValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2StringMap {
		delete(n.MyString2StringMap, k)
		break
	}

	if !patchUpdate(o, n, t) {
		return
	}

	if !checkPrimitive(o, n, t) {
		return
	}
}

func TestMapPrimitiveAddDelValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2StringMap {
		delete(n.MyString2StringMap, k)
		break
	}
	n.MyString2StringMap["new"] = "new"

	if !patchUpdate(o, n, t) {
		return
	}

	if !checkPrimitive(o, n, t) {
		return
	}
}
