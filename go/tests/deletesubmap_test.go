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
	"fmt"
	"github.com/saichler/l8reflect/go/tests/utils"
	"github.com/saichler/l8types/go/testtypes"
	"testing"
)

func TestDeleteSubSubMap(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)

	deleted := false
	for k1, v := range n.MyString2ModelMap {
		fmt.Println(k1)
		for _, v1 := range v.MySubs {
			var key int32
			if len(v1.Int32Map) != 4 {
				continue
			}
			for k, _ := range v1.Int32Map {
				key = k
				break
			}
			fmt.Println("Deleted key ", key)
			delete(v1.Int32Map, key)
			deleted = true
			break
		}
		if deleted {
			break
		}
	}

	if !patchUpdateProperty(o, n, z, t) {
		return
	}

	if !checkSubSubMap(z, n, t) {
		return
	}
}

func checkSubSubMap(o, n *testtypes.TestProto, t *testing.T) bool {
	if o.MyString2ModelMap == nil {
		log.Fail(t, "Expected map to not be nil")
		return false
	}
	if len(o.MyString2ModelMap) != len(n.MyString2ModelMap) {
		log.Fail(t, "maps are not the same len")
		return false
	}
	for nk, nv := range n.MyString2ModelMap {
		ov := o.MyString2ModelMap[nk]
		if ov == nil {
			log.Fail(t, "Expected key to exist in old map")
			return false
		}
		if ov.MyString != nv.MyString {
			log.Fail(t, "Expected values to match for key")
			return false
		}
		for nk1, nv1 := range nv.MySubs {
			ov1 := ov.MySubs[nk1]
			if ov1 == nil {
				log.Fail(t, "Expected sub to exist in old map")
				return false
			}
			if ov1.MyString != nv1.MyString {
				log.Fail(t, "expected sub value to be eq")
				return false
			}
			if len(ov1.Int32Map) != len(nv1.Int32Map) {
				log.Fail(t, "Expected maps to be eq")
				return false
			}
		}
	}
	return true
}
