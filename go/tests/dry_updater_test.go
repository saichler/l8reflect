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
	"github.com/saichler/l8reflect/go/reflect/updating"
	"github.com/saichler/l8reflect/go/tests/utils"
	"github.com/saichler/l8test/go/infra/t_resources"
	"github.com/saichler/l8types/go/testtypes"
)

func TestDryUpdater(t *testing.T) {
	res := newResources()

	aside := utils.CreateTestModelInstance(0)
	zside := t_resources.CloneTestModel(aside)
	zside.MyString = "updated"

	snapshot := cloning.NewCloner().Clone(aside).(*testtypes.TestProto)

	upd := updating.NewUpdater(res, false, false)
	err := upd.DryUpdate(aside, zside)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}

	changes := upd.Changes()
	if len(changes) != 1 {
		log.Fail(t, "Expected 1 change but got ", len(changes))
		return
	}

	// Verify old instance was NOT mutated
	if aside.MyString != snapshot.MyString {
		log.Fail(t, "DryUpdate mutated old: expected ", snapshot.MyString, " got ", aside.MyString)
		return
	}

	// Verify changes can still be applied manually
	for _, change := range changes {
		change.Apply(aside)
	}

	if aside.MyString != "updated" {
		log.Fail(t, "Apply after DryUpdate failed: expected updated got ", aside.MyString)
		return
	}
}

func TestDryUpdaterEnum(t *testing.T) {
	res := newResources()

	aside := utils.CreateTestModelInstance(0)
	zside := cloning.NewCloner().Clone(aside).(*testtypes.TestProto)
	zside.MyEnum = testtypes.TestEnum_ValueTwo

	snapshot := cloning.NewCloner().Clone(aside).(*testtypes.TestProto)

	upd := updating.NewUpdater(res, false, false)
	err := upd.DryUpdate(aside, zside)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}

	if aside.MyEnum != snapshot.MyEnum {
		log.Fail(t, "DryUpdate mutated enum: expected ", snapshot.MyEnum, " got ", aside.MyEnum)
		return
	}

	if len(upd.Changes()) == 0 {
		log.Fail(t, "Expected changes for enum update")
		return
	}
}

func TestDryUpdaterSubMap(t *testing.T) {
	res := newResources()

	aside := utils.CreateTestModelInstance(0)
	zside := t_resources.CloneTestModel(aside)
	for _, sub := range zside.MySingle.MySubs {
		sub.Int32Map[0]++
	}

	snapshot := cloning.NewCloner().Clone(aside).(*testtypes.TestProto)

	upd := updating.NewUpdater(res, false, false)
	err := upd.DryUpdate(aside, zside)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}

	if len(upd.Changes()) == 0 {
		log.Fail(t, "Expected changes")
		return
	}

	// Verify old instance was NOT mutated
	for k, snapSub := range snapshot.MySingle.MySubs {
		asideSub := aside.MySingle.MySubs[k]
		if snapSub.Int32Map[0] != asideSub.Int32Map[0] {
			log.Fail(t, "DryUpdate mutated sub map value")
			return
		}
	}
}

func TestDryUpdaterMatchesUpdate(t *testing.T) {
	res := newResources()

	aside := utils.CreateTestModelInstance(0)
	zside := t_resources.CloneTestModel(aside)
	zside.MyString = "updated"
	zside.MyInt32 = 999

	// DryUpdate
	dryUpd := updating.NewUpdater(res, false, false)
	err := dryUpd.DryUpdate(
		cloning.NewCloner().Clone(aside).(*testtypes.TestProto),
		cloning.NewCloner().Clone(zside).(*testtypes.TestProto),
	)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}

	// Regular Update
	regUpd := updating.NewUpdater(res, false, false)
	err = regUpd.Update(
		cloning.NewCloner().Clone(aside).(*testtypes.TestProto),
		cloning.NewCloner().Clone(zside).(*testtypes.TestProto),
	)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}

	dryChanges := dryUpd.Changes()
	regChanges := regUpd.Changes()

	if len(dryChanges) != len(regChanges) {
		log.Fail(t, "Change count mismatch: dry=", len(dryChanges), " reg=", len(regChanges))
		return
	}

	regSet := make(map[string]bool)
	for _, rc := range regChanges {
		regSet[rc.PropertyId()] = true
	}
	for _, dc := range dryChanges {
		if !regSet[dc.PropertyId()] {
			log.Fail(t, "DryUpdate produced change not in Update: ", dc.PropertyId())
			return
		}
	}
}
