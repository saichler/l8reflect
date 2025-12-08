package tests

import (
	"testing"

	"github.com/saichler/l8reflect/go/reflect/properties"
	"github.com/saichler/l8reflect/go/reflect/updating"
	"github.com/saichler/l8reflect/go/tests/utils"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/testtypes"
)

func patchUpdateProperty(o, n, z *testtypes.TestProto, t *testing.T) bool {
	res := newResources()

	u := updating.NewUpdater(res, false, true)
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

func TestMapPrimitivePropertySetFromNil(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	o.MyString2StringMap = nil
	z.MyString2StringMap = nil

	if !patchUpdateProperty(o, n, z, t) {
		return
	}

	if !checkPrimitive(z, n, t) {
		return
	}
}

func TestMapPrimitivePropertySetFromEmpty(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	o.MyString2StringMap = make(map[string]string)
	z.MyString2StringMap = make(map[string]string)

	if !patchUpdateProperty(o, n, z, t) {
		return
	}

	if !checkPrimitive(z, n, t) {
		return
	}
}

func TestMapPrimitivePropertyChangeValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2StringMap {
		n.MyString2StringMap[k] = n.MyString2StringMap[k] + "C"
	}

	if !patchUpdateProperty(o, n, z, t) {
		return
	}

	if !checkPrimitive(z, n, t) {
		return
	}
}

func TestMapPrimitivePropertyAddValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	n.MyString2StringMap["new"] = "new"

	if !patchUpdateProperty(o, n, z, t) {
		return
	}

	if !checkPrimitive(z, n, t) {
		return
	}
}

func TestMapPrimitivePropertyDelValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2StringMap {
		delete(n.MyString2StringMap, k)
		break
	}

	if !patchUpdateProperty(o, n, z, t) {
		return
	}

	if !checkPrimitive(z, n, t) {
		return
	}
}

func TestMapPrimitivePropertyAddDelValue(t *testing.T) {
	o := utils.CreateTestModelInstance(1)
	n := utils.CreateTestModelInstance(1)
	z := utils.CreateTestModelInstance(1)
	for k, _ := range o.MyString2StringMap {
		delete(n.MyString2StringMap, k)
		break
	}
	n.MyString2StringMap["new"] = "new"

	if !patchUpdateProperty(o, n, z, t) {
		return
	}

	if !checkPrimitive(z, n, t) {
		return
	}
}
