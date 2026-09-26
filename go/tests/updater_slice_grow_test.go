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

	"github.com/saichler/l8reflect/go/reflect/cloning"
	"github.com/saichler/l8reflect/go/reflect/updating"
	"github.com/saichler/l8reflect/go/tests/utils"
	"github.com/saichler/l8types/go/testtypes"
)

// growSlice updates aside with zside, where zside has a longer slice, and
// checks that aside ends up with zside's slice.
func growSlice(t *testing.T, grow func(*testtypes.TestProto), get func(*testtypes.TestProto) interface{}) {
	t.Helper()
	res := newResources()
	aside := utils.CreateTestModelInstance(0)
	zside := cloning.NewCloner().Clone(aside).(*testtypes.TestProto)
	grow(zside)

	upd := updating.NewUpdater(res, false, false)
	if err := upd.Update(aside, zside); err != nil {
		log.Fail(t, err.Error())
		return
	}
	if len(upd.Changes()) == 0 {
		log.Fail(t, "Expected changes for the grown slice")
		return
	}
	if !reflect.DeepEqual(get(aside), get(zside)) {
		log.Fail(t, "Expected ", get(zside), " got ", get(aside))
	}
}

func TestSliceGrowString(t *testing.T) {
	growSlice(t, func(p *testtypes.TestProto) {
		p.MyStringSlice = append(p.MyStringSlice, "added-1", "added-2")
	}, func(p *testtypes.TestProto) interface{} { return p.MyStringSlice })
}

func TestSliceGrowInt32(t *testing.T) {
	growSlice(t, func(p *testtypes.TestProto) {
		p.MyInt32Slice = append(p.MyInt32Slice, 1001, 1002)
	}, func(p *testtypes.TestProto) interface{} { return p.MyInt32Slice })
}

func TestSliceGrowModel(t *testing.T) {
	growSlice(t, func(p *testtypes.TestProto) {
		p.MyModelSlice = append(p.MyModelSlice, &testtypes.TestProtoSub{MyString: "added", MyInt64: 7})
	}, func(p *testtypes.TestProto) interface{} {
		out := []string{}
		for _, s := range p.MyModelSlice {
			out = append(out, s.MyString)
		}
		return out
	})
}
