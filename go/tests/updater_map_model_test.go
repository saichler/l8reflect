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
	"github.com/saichler/l8reflect/go/tests/utils"
	"github.com/saichler/l8types/go/testtypes"
	"testing"
)

func checkStruct(o, n *testtypes.TestProto, t *testing.T) bool {
	if o.MyString2ModelMap == nil {
		log.Fail(t, "Expected map to not be nil")
		return false
	}
	if len(o.MyString2ModelMap) != len(n.MyString2ModelMap) {
		log.Fail(t, "maps are not the same len")
		return false
	}
	for k, v := range n.MyString2ModelMap {
		vo, ok := o.MyString2ModelMap[k]
		if !ok {
			log.Fail(t, "Expected key to exist in old map")
			return false
		}
		if vo.MyString != v.MyString {
			log.Fail(t, "Expected values to match for key")
			return false
		}
	}
	return true
}

func TestMapModelSetFromNil(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	o.MyString2ModelMap = nil

	if !patchUpdate(o, n, t) {
		return
	}

	if !checkStruct(o, n, t) {
		return
	}
}

func TestMapModelSetFromEmpty(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	o.MyString2ModelMap = make(map[string]*testtypes.TestProtoSub)

	if !patchUpdate(o, n, t) {
		return
	}

	if !checkStruct(o, n, t) {
		return
	}
}

func TestMapModelChangeValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2ModelMap {
		n.MyString2ModelMap[k] = &testtypes.TestProtoSub{MyString: k + "-Hello"}
	}

	if !patchUpdate(o, n, t) {
		return
	}

	if !checkStruct(o, n, t) {
		return
	}
}

func TestMapModelChangeInternalValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2ModelMap {
		n.MyString2ModelMap[k].MyString = k + "changed"
	}

	if !patchUpdate(o, n, t) {
		return
	}

	if !checkStruct(o, n, t) {
		return
	}
}

func TestMapAddValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	n.MyString2ModelMap["new"] = &testtypes.TestProtoSub{MyString: "new"}

	if !patchUpdate(o, n, t) {
		return
	}

	if !checkStruct(o, n, t) {
		return
	}
}

func TestMapModelDelValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2ModelMap {
		delete(n.MyString2ModelMap, k)
		break
	}

	if !patchUpdate(o, n, t) {
		return
	}

	if !checkStruct(o, n, t) {
		return
	}
}

func TestMapStructAddDelValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2ModelMap {
		delete(n.MyString2ModelMap, k)
		break
	}
	n.MyString2ModelMap["new"] = &testtypes.TestProtoSub{MyString: "new"}

	if !patchUpdate(o, n, t) {
		return
	}

	if !checkStruct(o, n, t) {
		return
	}
}
