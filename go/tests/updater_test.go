package tests

import (
	"fmt"
	"github.com/saichler/reflect/go/reflect/cloning"
	"github.com/saichler/reflect/go/reflect/introspecting"
	"github.com/saichler/reflect/go/reflect/properties"
	"github.com/saichler/reflect/go/reflect/updating"
	"github.com/saichler/reflect/go/tests/utils"
	"github.com/saichler/serializer/go/serialize/object"
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
		log.Fail(t, "Expected 1 change but got ", len(upd.Changes()))
		for _, c := range changes {
			log.Info(c.String())
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

	fmt.Println(upd.Changes()[0].PropertyId())
}

func TestSubMapDeep(t *testing.T) {
	in := introspecting.NewIntrospect(registry.NewRegistry())
	_, err := in.Inspect(&testtypes.TestProto{})

	if err != nil {
		log.Fail(t, err.Error())
		return
	}
	upd := updating.NewUpdater(in, false)
	aside := utils.CreateTestModelInstance(0)
	zside := cloning.NewCloner().Clone(aside).(*testtypes.TestProto)
	yside := cloning.NewCloner().Clone(aside).(*testtypes.TestProto)
	zside.MyString2ModelMap["newone"] = &testtypes.TestProtoSub{MyString: "newone"}

	err = upd.Update(aside, zside)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}

	newone := aside.MyString2ModelMap["newone"]
	if newone == nil {
		log.Fail(t, "new one is nil")
		return
	}

	for _, chg := range upd.Changes() {
		fmt.Println(chg.PropertyId())
		chg.Apply(yside)
	}

	newone = yside.MyString2ModelMap["newone"]
	if newone == nil {
		log.Fail(t, "new one is nil")
		return
	}

	aside = utils.CreateTestModelInstance(0)
	zside = cloning.NewCloner().Clone(aside).(*testtypes.TestProto)
	yside = cloning.NewCloner().Clone(aside).(*testtypes.TestProto)
	zside.MyString2ModelMap["newone"] = &testtypes.TestProtoSub{MyString: "newone"}
	aside.MyString2ModelMap["newone"] = &testtypes.TestProtoSub{MyString: "newone"}
	yside.MyString2ModelMap["newone"] = &testtypes.TestProtoSub{MyString: "newone"}

	zside.MyString2ModelMap["newone"].MyString = "newer"
	upd = updating.NewUpdater(in, false)
	err = upd.Update(aside, zside)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}

	newone = aside.MyString2ModelMap["newone"]
	if newone.MyString != "newer" {
		log.Fail(t, "Expected newer")
		return
	}

	for _, chg := range upd.Changes() {
		chg.Apply(yside)
	}

	newone = yside.MyString2ModelMap["newone"]
	if newone.MyString != "newer" {
		log.Fail(t, "expected newer")
		return
	}

	aside = utils.CreateTestModelInstance(0)
	zside = cloning.NewCloner().Clone(aside).(*testtypes.TestProto)
	yside = cloning.NewCloner().Clone(aside).(*testtypes.TestProto)
	xside := cloning.NewCloner().Clone(aside).(*testtypes.TestProto)

	zside.MyString2ModelMap["newone"] = &testtypes.TestProtoSub{MyString: "newone"}
	aside.MyString2ModelMap["newone"] = &testtypes.TestProtoSub{MyString: "newone"}
	yside.MyString2ModelMap["newone"] = &testtypes.TestProtoSub{MyString: "newone"}
	xside.MyString2ModelMap["newone"] = &testtypes.TestProtoSub{MyString: "newone"}

	zside.MyString2ModelMap["newone"].MySubs = make(map[string]*testtypes.TestProtoSubSub)
	aside.MyString2ModelMap["newone"].MySubs = make(map[string]*testtypes.TestProtoSubSub)
	yside.MyString2ModelMap["newone"].MySubs = make(map[string]*testtypes.TestProtoSubSub)
	xside.MyString2ModelMap["newone"].MySubs = make(map[string]*testtypes.TestProtoSubSub)

	zside.MyString2ModelMap["newone"].MySubs["newsub"] = &testtypes.TestProtoSubSub{}
	zside.MyString2ModelMap["newone"].MySubs["newsub"].MyString = "newsub"
	aside.MyString2ModelMap["newone"].MySubs["newsub"] = &testtypes.TestProtoSubSub{}
	aside.MyString2ModelMap["newone"].MySubs["newsub"].MyString = "newsub"
	yside.MyString2ModelMap["newone"].MySubs["newsub"] = &testtypes.TestProtoSubSub{}
	yside.MyString2ModelMap["newone"].MySubs["newsub"].MyString = "newsub"
	xside.MyString2ModelMap["newone"].MySubs["newsub"] = &testtypes.TestProtoSubSub{}
	xside.MyString2ModelMap["newone"].MySubs["newsub"].MyString = "newsub"

	zside.MyString2ModelMap["newone"].MySubs["newsub"].MyString = "newersub"

	upd = updating.NewUpdater(in, false)
	err = upd.Update(aside, zside)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}

	val := aside.MyString2ModelMap["newone"].MySubs["newsub"].MyString
	if val != "newersub" {
		log.Fail(t, "expected newersub")
		return
	}

	pid := ""

	for _, chg := range upd.Changes() {
		pid = chg.PropertyId()
		chg.Apply(yside)
	}

	val = yside.MyString2ModelMap["newone"].MySubs["newsub"].MyString
	if val != "newersub" {
		log.Fail(t, "expected newersub")
		return
	}

	prop, err := properties.PropertyOf(pid, in)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}

	prop.Set(xside, "newersub")
	val = xside.MyString2ModelMap["newone"].MySubs["newsub"].MyString
	if val != "newersub" {
		log.Fail(t, "expected newersub")
		return
	}
}

func checkEQ(aside, zside interface{}, t *testing.T) bool {
	de := cloning.NewDeepEqual()
	if de.Equal(aside, zside) {
		log.Info("Equale")
		return true
	} else {
		log.Fail(t, "Not Equale")
		log.Error(aside)
		log.Error(zside)
	}
	return false
}

func TestSubMapDeepAlwaysChanging(t *testing.T) {
	in := introspecting.NewIntrospect(registry.NewRegistry())
	_, err := in.Inspect(&testtypes.TestProto{})

	if err != nil {
		log.Fail(t, err.Error())
		return
	}

	aside := utils.CreateTestModelInstance(0)
	zside := cloning.NewCloner().Clone(aside).(*testtypes.TestProto)
	yside := cloning.NewCloner().Clone(aside).(*testtypes.TestProto)
	xside := cloning.NewCloner().Clone(aside).(*testtypes.TestProto)

	zside.MyString2ModelMap["newone"] = &testtypes.TestProtoSub{MyString: "newone"}
	aside.MyString2ModelMap["newone"] = &testtypes.TestProtoSub{MyString: "newone"}
	yside.MyString2ModelMap["newone"] = &testtypes.TestProtoSub{MyString: "newone"}
	xside.MyString2ModelMap["newone"] = &testtypes.TestProtoSub{MyString: "newone"}

	zside.MyString2ModelMap["newone"].MySubs = make(map[string]*testtypes.TestProtoSubSub)
	aside.MyString2ModelMap["newone"].MySubs = make(map[string]*testtypes.TestProtoSubSub)
	yside.MyString2ModelMap["newone"].MySubs = make(map[string]*testtypes.TestProtoSubSub)
	xside.MyString2ModelMap["newone"].MySubs = make(map[string]*testtypes.TestProtoSubSub)

	zside.MyString2ModelMap["newone"].MySubs["newsub"] = &testtypes.TestProtoSubSub{}
	zside.MyString2ModelMap["newone"].MySubs["newsub"].MyString = "newsub"
	aside.MyString2ModelMap["newone"].MySubs["newsub"] = &testtypes.TestProtoSubSub{}
	aside.MyString2ModelMap["newone"].MySubs["newsub"].MyString = "newsub"
	yside.MyString2ModelMap["newone"].MySubs["newsub"] = &testtypes.TestProtoSubSub{}
	yside.MyString2ModelMap["newone"].MySubs["newsub"].MyString = "newsub"
	xside.MyString2ModelMap["newone"].MySubs["newsub"] = &testtypes.TestProtoSubSub{}
	xside.MyString2ModelMap["newone"].MySubs["newsub"].MyString = "newsub"

	if !checkEQ(zside.MyString2ModelMap, xside.MyString2ModelMap, t) {
		log.Fail(t, "Init data isn't good")
		return
	}

	zside.MyString2ModelMap["newone"].MySubs["newsub"].MyString = "newersub"

	upd := updating.NewUpdater(in, false)
	err = upd.Update(aside, zside)

	if len(upd.Changes()) == 0 {
		log.Fail(t, "Expected changes")
		return
	}

	for _, chg := range upd.Changes() {
		chg.Apply(yside)
		prop, e := properties.PropertyOf(chg.PropertyId(), in)
		if e != nil {
			panic(e)
		}
		obj := object.NewEncode([]byte{}, 0)
		obj.Add(chg.NewValue())
		data := obj.Data()
		obj = object.NewDecode(data, 0, "", in.Registry())
		val, e := obj.Get()
		fmt.Println(val)
		_, _, e = prop.Set(xside, val)
		if e != nil {
			panic(e)
		}
	}

	if !checkEQ(zside.MyString2ModelMap, aside.MyString2ModelMap, t) {
		log.Fail(t, "Not EQ aside")
		return
	}

	if !checkEQ(zside.MyString2ModelMap, yside.MyString2ModelMap, t) {
		log.Fail(t, "Not EQ yside")
		return
	}

	if !checkEQ(zside.MyString2ModelMap, xside.MyString2ModelMap, t) {
		log.Fail(t, "Not EQ xside")
		return
	}
}
