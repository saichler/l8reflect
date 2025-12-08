package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/saichler/l8reflect/go/reflect/introspecting"
	"github.com/saichler/l8reflect/go/reflect/properties"
	"github.com/saichler/l8reflect/go/reflect/updating"
	"github.com/saichler/l8reflect/go/tests/utils"
	"github.com/saichler/l8test/go/infra/t_resources"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/testtypes"
	"github.com/saichler/l8utils/go/utils/registry"
	"github.com/saichler/l8utils/go/utils/resources"
)

func newResources() ifs.IResources {
	res := resources.NewResources(log)
	res.Set(registry.NewRegistry())
	in := introspecting.NewIntrospect(res.Registry())
	res.Set(in)
	addPrimary(res)
	return res
}

func propertyOf(id string, root interface{}, t *testing.T, res ifs.IResources) (interface{}, bool) {

	ins, err := properties.PropertyOf(id, res)
	if err != nil {
		log.Fail(t, "failed with id: ", id, err.Error())
		return nil, false
	}

	v, err := ins.Get(root)
	if err != nil {
		log.Fail(t, "failed with get: ", id, err.Error())
		return nil, false
	}
	return v, true
}

func addPrimary(res ifs.IResources) {
	res.Introspector().Decorators().AddPrimaryKeyDecorator(&testtypes.TestProto{}, "MyString")
}

func TestPrimaryKey(t *testing.T) {
	res := newResources()

	aside := utils.CreateTestModelInstance(1)
	zside := t_resources.CloneTestModel(aside)
	zside.MyEnum = testtypes.TestEnum_ValueTwo

	upd := updating.NewUpdater(res, false, false)
	err := upd.Update(aside, zside)
	if err != nil {
		log.Fail(t, "failed with update: ", err.Error())
		return
	}
	if len(upd.Changes()) != 1 {
		log.Fail(t, "wrong number of changes: ", len(upd.Changes()))
		return
	}

	pid := upd.Changes()[0].PropertyId()
	n := upd.Changes()[0].NewValue()

	p, e := properties.PropertyOf(pid, res)
	if e != nil {
		log.Fail(t, "failed with property: ", e.Error())
		return
	}

	_, root, e := p.Set(nil, n)
	if e != nil {
		log.Fail(t, "failed with set: ", e.Error())
		return
	}

	yside := root.(*testtypes.TestProto)
	if yside.MyEnum != aside.MyEnum {
		log.Fail(t, "wrong enum: ", yside.MyEnum)
		return
	}
	if yside.MyString != aside.MyString {
		log.Fail(t, "wrong string: ", yside.MyString)
		return
	}

	pid = "testproto.myenum"
	prod, err := properties.PropertyOf(pid, res)
	if err != nil {
		log.Fail(t, "failed with property: ", err.Error())
		return
	}

	res.Registry().RegisterEnums(testtypes.TestEnum_value)

	_, _, err = prod.Set(yside, "ValueOne")
	if err != nil {
		log.Fail(t, "failed with set: ", err.Error())
		return
	}

	if yside.MyEnum != testtypes.TestEnum_ValueOne {
		log.Fail(t, "wrong enum: ", yside.MyEnum)
		return
	}
}

func TestSetMap(t *testing.T) {
	res := newResources()

	aside := utils.CreateTestModelInstance(1)
	aside.MyString2ModelMap = nil
	pid := "testproto<{24}{24}root>.mystring2modelmap<{24}sub>.mystring"
	//m:=testtypes.TestProtoSub{}
	prop, err := properties.PropertyOf(pid, res)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}

	fmt.Println(prop.PropertyId())

	_, _, err = prop.Set(aside, "hhhh")
	if err != nil {
		log.Fail(t, err.Error())
		return
	}
	sub := aside.MyString2ModelMap["sub"]
	if sub == nil {
		log.Fail(t, "sub doesn't exist")
	}
	if sub.MyString != "hhhh" {
		log.Fail(t, "sub MyString exist")
		return
	}

	prop, _ = properties.PropertyOf("testproto.mystring2modelmap", res)
	m := aside.MyString2ModelMap
	_, _, err = prop.Set(aside, nil)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}

	if len(aside.MyString2ModelMap) != 0 {
		log.Fail(t, "expected map to be empty")
		return
	}

	_, _, err = prop.Set(aside, m)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}

	if len(aside.MyString2ModelMap) == 0 {
		log.Fail(t, "expected map to be non-empty")
		return
	}

}

func TestInstance(t *testing.T) {
	res := newResources()

	id := "testproto<{24}{24}Hello>"
	v, ok := propertyOf(id, nil, t, res)
	if !ok {
		return
	}

	mytest := v.(*testtypes.TestProto)
	if mytest.MyString != "Hello" {
		log.Fail(t, "wrong string: ", mytest.MyString)
		return
	}

	mytest.MyFloat64 = 128.128
	id = "testproto.myfloat64"
	v, ok = propertyOf(id, mytest, t, res)
	if !ok {
		return
	}

	f := v.(float64)
	if f != mytest.MyFloat64 {
		log.Fail(t, "wrong float64: ", f)
		return
	}

	mytest.MySingle = &testtypes.TestProtoSub{MyString: "Hello"}

	id = "testproto.mysingle.mystring"
	v, ok = propertyOf(id, mytest, t, res)
	if !ok {
		return
	}
	s := v.(string)
	if s != mytest.MySingle.MyString {
		log.Fail(t, "wrong string: ", s)
		return
	}

	/*
		myInstsnce:=model.MyTestModel{
			MyString: "Hello",
			MySingle: &model.MyTestSubModelSingle{MyString: "World"},
		}

		instance,_:=instance.propertyOf("mytestmodel.mysingle.mystring",introspect.DefaultIntrospect)

		//Getting a value
		v,_:=instance.Get(myInstsnce)
		//Creating another instance
		myOtherInstance:=model.MyTestModel{}
		//Setting the value we fetched from the original instance
		instance.Set(myOtherInstance,"Metadata")

	*/
}

func TestSubStructProperty(t *testing.T) {
	res := newResources()

	aside := &testtypes.TestProto{MyString: "Hello"}
	zside := &testtypes.TestProto{MyString: "Hello"}
	yside := &testtypes.TestProto{MyString: "Hello"}
	zside.MySingle = &testtypes.TestProtoSub{MyInt64: time.Now().Unix()}

	putUpdater := updating.NewUpdater(res, false, false)

	putUpdater.Update(aside, zside)

	changes := putUpdater.Changes()

	for _, change := range changes {
		change.Apply(yside)
	}
	fmt.Println(yside)
}
