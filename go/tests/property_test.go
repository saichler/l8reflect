package tests

import (
	"fmt"
	"github.com/saichler/reflect/go/reflect/inspect"
	"github.com/saichler/reflect/go/reflect/property"
	"github.com/saichler/reflect/go/reflect/updater"
	"github.com/saichler/reflect/go/tests/utils"
	"github.com/saichler/shared/go/share/registry"
	"github.com/saichler/shared/go/tests"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/types"
	"testing"
	"time"
)

var _introspect common.IIntrospector

func propertyOf(id string, root interface{}, t *testing.T) (interface{}, bool) {
	ins, err := property.PropertyOf(id, _introspect)
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

func TestPrimaryKey(t *testing.T) {
	_introspect = inspect.NewIntrospect(registry.NewRegistry())
	node, err := _introspect.Inspect(&tests.TestProto{})
	if err != nil {
		log.Fail(t, "failed with inspect: ", err.Error())
		return
	}
	_introspect.AddDecorator(types.DecoratorType_Primary, []string{"MyString"}, node)
	aside := utils.CreateTestModelInstance(1)
	zside := utils.CreateTestModelInstance(1)
	zside.MyEnum = tests.TestEnum_ValueTwo

	upd := updater.NewUpdater(_introspect, false)
	err = upd.Update(aside, zside)
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

	p, e := property.PropertyOf(pid, _introspect)
	if e != nil {
		log.Fail(t, "failed with property: ", e.Error())
		return
	}

	_, root, e := p.Set(nil, n)
	if e != nil {
		log.Fail(t, "failed with set: ", e.Error())
		return
	}

	yside := root.(*tests.TestProto)
	if yside.MyEnum != aside.MyEnum {
		log.Fail(t, "wrong enum: ", yside.MyEnum)
		return
	}
	if yside.MyString != aside.MyString {
		log.Fail(t, "wrong string: ", yside.MyString)
		return
	}

	pid = "testproto.myenum"
	prod, err := property.PropertyOf(pid, _introspect)
	if err != nil {
		log.Fail(t, "failed with property: ", err.Error())
		return
	}

	_introspect.Registry().RegisterEnums(tests.TestEnum_value)

	_, _, err = prod.Set(yside, "ValueOne")
	if err != nil {
		log.Fail(t, "failed with set: ", err.Error())
		return
	}

	if yside.MyEnum != tests.TestEnum_ValueOne {
		log.Fail(t, "wrong enum: ", yside.MyEnum)
		return
	}
}

func TestSetMap(t *testing.T) {
	_introspect := inspect.NewIntrospect(registry.NewRegistry())
	node, err := _introspect.Inspect(&tests.TestProto{})
	if err != nil {
		log.Fail(t, "failed with inspect: ", err.Error())
		return
	}
	_introspect.AddDecorator(types.DecoratorType_Primary, []string{"MyString"}, node)
	aside := utils.CreateTestModelInstance(1)
	aside.MyString2ModelMap = nil
	pid := "testproto<{24}root>.mystring2modelmap<{24}sub>.mystring"
	//m:=tests.TestProtoSub{}
	prop, err := property.PropertyOf(pid, _introspect)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}
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

	prop, _ = property.PropertyOf("testproto.mystring2modelmap", _introspect)
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
	_introspect = inspect.NewIntrospect(registry.NewRegistry())
	node, err := _introspect.Inspect(&tests.TestProto{})
	if err != nil {
		log.Fail(t, "failed with inspect: ", err.Error())
		return
	}
	_introspect.AddDecorator(types.DecoratorType_Primary, []string{"MyString"}, node)

	id := "testproto<{24}Hello>"
	v, ok := propertyOf(id, nil, t)
	if !ok {
		return
	}

	mytest := v.(*tests.TestProto)
	if mytest.MyString != "Hello" {
		log.Fail(t, "wrong string: ", mytest.MyString)
		return
	}

	mytest.MyFloat64 = 128.128
	id = "testproto.myfloat64"
	v, ok = propertyOf(id, mytest, t)
	if !ok {
		return
	}

	f := v.(float64)
	if f != mytest.MyFloat64 {
		log.Fail(t, "wrong float64: ", f)
		return
	}

	mytest.MySingle = &tests.TestProtoSub{MyString: "Hello"}

	id = "testproto.mysingle.mystring"
	v, ok = propertyOf(id, mytest, t)
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
	_introspect = inspect.NewIntrospect(registry.NewRegistry())
	node, err := _introspect.Inspect(&tests.TestProto{})
	if err != nil {
		log.Fail(t, "failed with inspect: ", err.Error())
		return
	}
	_introspect.AddDecorator(types.DecoratorType_Primary, []string{"MyString"}, node)

	aside := &tests.TestProto{MyString: "Hello"}
	zside := &tests.TestProto{MyString: "Hello"}
	yside := &tests.TestProto{MyString: "Hello"}
	zside.MySingle = &tests.TestProtoSub{MyInt64: time.Now().Unix()}

	putUpdater := updater.NewUpdater(_introspect, false)

	putUpdater.Update(aside, zside)

	changes := putUpdater.Changes()

	for _, change := range changes {
		change.Apply(yside)
	}
	fmt.Println(yside)
}
