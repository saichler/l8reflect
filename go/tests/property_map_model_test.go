package tests

import (
	"github.com/saichler/l8reflect/go/tests/utils"
	"github.com/saichler/l8types/go/testtypes"
	"testing"
)

func TestMapModelPropertySetFromNil(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	o.MyString2ModelMap = nil

	if !patchUpdateProperty(o, n, z, t) {
		return
	}

	if !checkStruct(z, n, t) {
		return
	}
}

func TestMapModelPropertySetFromEmpty(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)

	o.MyString2ModelMap = make(map[string]*testtypes.TestProtoSub)

	if !patchUpdateProperty(o, n, z, t) {
		return
	}

	if !checkStruct(z, n, t) {
		return
	}
}

func TestMapModelPropertyChangeValue(t *testing.T) {
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
	if !patchUpdateProperty(o, n, z, t) {
		return
	}

	if !checkStruct(z, n, t) {
		return
	}
}

func TestMapModelChangePropertyInternalValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2ModelMap {
		n.MyString2ModelMap[k].MyString = k + "changed"
	}

	if !patchUpdateProperty(o, n, z, t) {
		return
	}

	if !checkStruct(z, n, t) {
		return
	}
}

func TestMapAddPropertyModelValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	n.MyString2ModelMap["new"] = &testtypes.TestProtoSub{MyString: "new"}

	if !patchUpdateProperty(o, n, z, t) {
		return
	}

	if !checkStruct(z, n, t) {
		return
	}
}

func TestMapModelDelPropertyValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2ModelMap {
		delete(n.MyString2ModelMap, k)
		break
	}

	if !patchUpdateProperty(o, n, z, t) {
		return
	}

	if !checkStruct(z, n, t) {
		return
	}
}

func TestMapStructPropertyAddDelValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2ModelMap {
		delete(n.MyString2ModelMap, k)
		break
	}
	n.MyString2ModelMap["new"] = &testtypes.TestProtoSub{MyString: "new"}

	if !patchUpdateProperty(o, n, z, t) {
		return
	}

	if !checkStruct(z, n, t) {
		return
	}
}
