package tests

import (
	"fmt"
	"github.com/saichler/reflect/go/reflect/common"
	"github.com/saichler/reflect/go/reflect/inspect"
	"github.com/saichler/reflect/go/reflect/property"
	"github.com/saichler/reflect/go/types"
	"github.com/saichler/shared/go/share/registry"
	"github.com/saichler/shared/go/tests"
	"testing"
)

var _introspect common.IIntrospect

func propertyOf(id string, root interface{}, t *testing.T) (interface{}, bool) {
	ins, err := property.PropertyOf(id, _introspect)
	if err != nil {
		t.Fail()
		fmt.Println("failed with id: ", id, err)
		return nil, false
	}

	v, err := ins.Get(root)
	if err != nil {
		t.Fail()
		fmt.Println("failed with get: ", id, err)
		return nil, false
	}
	return v, true
}

func TestInstance(t *testing.T) {
	_introspect = inspect.NewIntrospect(registry.NewRegistry())
	node, err := _introspect.Inspect(&tests.TestProto{})
	if err != nil {
		fmt.Println("1", err)
		t.Fail()
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
		t.Fail()
		fmt.Println("Expected Hello but got ", mytest.MyString)
	}

	mytest.MyFloat64 = 128.128
	id = "testproto.myfloat64"
	v, ok = propertyOf(id, mytest, t)
	if !ok {
		return
	}

	f := v.(float64)
	if f != mytest.MyFloat64 {
		t.Fail()
		fmt.Println("float64 failed:", mytest.MyFloat64, "!=", f)
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
		t.Fail()
		fmt.Println("sum model string failed:", mytest.MySingle.MyString, "!=", f)
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
