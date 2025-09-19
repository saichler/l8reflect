package tests

import (
	"github.com/saichler/l8reflect/go/tests/utils"
	"github.com/saichler/l8types/go/testtypes"
	"testing"
)

func TestMapModelApplySetFromNil(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	o.MyString2ModelMap = nil

	if !patchUpdateApply(o, n, z, t) {
		return
	}

	if !checkStruct(z, n, t) {
		return
	}
}

func TestMapModelApplySetFromEmpty(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)

	o.MyString2ModelMap = make(map[string]*testtypes.TestProtoSub)

	if !patchUpdateApply(o, n, z, t) {
		return
	}

	if !checkStruct(z, n, t) {
		return
	}
}

func TestMapModelApplyChangeValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2ModelMap {
		n.MyString2ModelMap[k] = &testtypes.TestProtoSub{MyString: k + "-Hello"}
	}
	//This is because the pointer for this element is used in multiple attributes
	//so to avoid double changed from othe rproperties.
	for k, v := range n.MyString2ModelMap {
		o.MyString2ModelMap[k] = &testtypes.TestProtoSub{MyString: v.MyString}
		z.MyString2ModelMap[k] = &testtypes.TestProtoSub{MyString: v.MyString}
	}
	if !patchUpdateApply(o, n, z, t) {
		return
	}

	if !checkStruct(z, n, t) {
		return
	}
}

func TestMapModelChangeApplyInternalValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2ModelMap {
		n.MyString2ModelMap[k].MyString = k + "changed"
	}

	if !patchUpdateApply(o, n, z, t) {
		return
	}

	if !checkStruct(z, n, t) {
		return
	}
}

func TestMapAddModelValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	n.MyString2ModelMap["new"] = &testtypes.TestProtoSub{MyString: "new"}

	if !patchUpdateApply(o, n, z, t) {
		return
	}

	if !checkStruct(z, n, t) {
		return
	}
}

func TestMapModelDelApplyValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2ModelMap {
		delete(n.MyString2ModelMap, k)
		break
	}

	if !patchUpdateApply(o, n, z, t) {
		return
	}

	if !checkStruct(z, n, t) {
		return
	}
}

func TestMapStructApplyAddDelValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2ModelMap {
		delete(n.MyString2ModelMap, k)
		break
	}
	n.MyString2ModelMap["new"] = &testtypes.TestProtoSub{MyString: "new"}

	if !patchUpdateApply(o, n, z, t) {
		return
	}

	if !checkStruct(z, n, t) {
		return
	}
}
