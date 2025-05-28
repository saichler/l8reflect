package tests

import (
	"github.com/saichler/l8types/go/testtypes"
	"github.com/saichler/reflect/go/reflect/updating"
	"github.com/saichler/reflect/go/tests/utils"
	"testing"
)

func patchUpdateApply(o, n, z *testtypes.TestProto, t *testing.T) bool {
	res := newResources()

	_, err := res.Introspector().Inspect(&testtypes.TestProto{})
	if err != nil {
		log.Fail(t, err.Error())
		return false
	}

	u := updating.NewUpdater(res, false, true)
	err = u.Update(o, n)
	if err != nil {
		log.Fail(t, err.Error())
		return false
	}

	for i, c := range u.Changes() {
		flog.Debug(i, " - ", c.PropertyId(), " - ", c.NewValue())
		c.Apply(z)
	}

	return true
}

func TestMapPrimitiveApplySetFromNil(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	o.MyString2StringMap = nil
	z.MyString2StringMap = nil

	if !patchUpdateApply(o, n, z, t) {
		return
	}

	if !checkPrimitive(z, n, t) {
		return
	}
}

func TestMapPrimitiveApplySetFromEmpty(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	o.MyString2StringMap = make(map[string]string)
	z.MyString2StringMap = make(map[string]string)

	if !patchUpdateApply(o, n, z, t) {
		return
	}

	if !checkPrimitive(z, n, t) {
		return
	}
}

func TestMapPrimitiveApplyChangeValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2StringMap {
		n.MyString2StringMap[k] = n.MyString2StringMap[k] + "C"
	}

	if !patchUpdateApply(o, n, z, t) {
		return
	}

	if !checkPrimitive(z, n, t) {
		return
	}
}

func TestMapPrimitiveApplyAddValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	n.MyString2StringMap["new"] = "new"

	if !patchUpdateApply(o, n, z, t) {
		return
	}

	if !checkPrimitive(z, n, t) {
		return
	}
}

func TestMapPrimitiveApplyDelValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2StringMap {
		delete(n.MyString2StringMap, k)
		break
	}

	if !patchUpdateApply(o, n, z, t) {
		return
	}

	if !checkPrimitive(z, n, t) {
		return
	}
}

func TestMapPrimitiveApplyAddDelValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2StringMap {
		delete(n.MyString2StringMap, k)
		break
	}
	n.MyString2StringMap["new"] = "new"

	if !patchUpdateApply(o, n, z, t) {
		return
	}

	if !checkPrimitive(z, n, t) {
		return
	}
}
