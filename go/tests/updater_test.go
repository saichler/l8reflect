package tests

import (
	"fmt"
	"github.com/saichler/reflect/go/reflect/cloning"
	"github.com/saichler/reflect/go/reflect/introspecting"
	"github.com/saichler/reflect/go/reflect/updating"
	"github.com/saichler/reflect/go/tests/utils"
	"github.com/saichler/shared/go/share/registry"
	"github.com/saichler/types/go/testtypes"
	"testing"
)

func TestUpdater(t *testing.T) {
	in := introspecting.NewIntrospect(registry.NewRegistry())
	_, err := in.Inspect(&testtypes.TestProto{})
	if err != nil {
		log.Fail(t, err.Error())
		return
	}
	upd := updating.NewUpdater(in, false)
	aside := utils.CreateTestModelInstance(0)
	zside := &testtypes.TestProto{MyString: "updated"}
	uside := in.Clone(aside).(*testtypes.TestProto)
	err = upd.Update(aside, zside)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}

	changes := upd.Changes()

	if len(changes) != 1 {
		t.Fail()
		fmt.Println("Expected 1 change but got ", len(upd.Changes()))
		for _, c := range changes {
			fmt.Println(c.String())
		}
		return
	}

	if aside.MyString != zside.MyString {
		t.Fail()
		fmt.Println("1 Expected ", zside.MyString, " got ", aside.MyString)
		return
	}

	for _, change := range changes {
		change.Apply(uside)
	}

	if uside.MyString != aside.MyString {
		fmt.Println("2 Expected ", aside.MyString, " got ", uside.MyString)
		t.Fail()
		return
	}
}

func TestEnum(t *testing.T) {
	in := introspecting.NewIntrospect(registry.NewRegistry())
	_, err := in.Inspect(&testtypes.TestProto{})
	if err != nil {
		log.Fail(t, err.Error())
		return
	}
	upd := updating.NewUpdater(in, false)
	aside := utils.CreateTestModelInstance(0)
	zside := cloning.NewCloner().Clone(aside).(*testtypes.TestProto)
	zside.MyEnum = testtypes.TestEnum_ValueTwo

	err = upd.Update(aside, zside)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}
	if aside.MyEnum != zside.MyEnum {
		log.Fail(t, aside.MyEnum)
		return
	}
}

func TestSubMap(t *testing.T) {
	in := introspecting.NewIntrospect(registry.NewRegistry())
	_, err := in.Inspect(&testtypes.TestProto{})
	if err != nil {
		log.Fail(t, err.Error())
		return
	}
	upd := updating.NewUpdater(in, false)
	aside := utils.CreateTestModelInstance(0)
	zside := cloning.NewCloner().Clone(aside).(*testtypes.TestProto)
	zside.MySingle.MySubs["sub"].Int32Map[0]++

	err = upd.Update(aside, zside)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}

	if zside.MySingle.MySubs["sub"].Int32Map[0] != aside.MySingle.MySubs["sub"].Int32Map[0] {
		log.Fail(t, aside.MySingle.MySubs["sub"].Int32Map[0])
		return
	}

	if len(upd.Changes()) == 0 {
		log.Fail(t, "Expected changes")
	}
}
