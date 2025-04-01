package cloning

import (
	"reflect"
)

type DeepEqual struct {
	comparators map[reflect.Kind]func(reflect.Value, reflect.Value) bool
}

func NewDeepEqual() *DeepEqual {
	de := &DeepEqual{}
	de.initCloners()
	return de
}

func (this *DeepEqual) initCloners() {
	this.comparators = make(map[reflect.Kind]func(reflect.Value, reflect.Value) bool)
	this.comparators[reflect.Int] = this.intComp
	this.comparators[reflect.Int32] = this.intComp
	this.comparators[reflect.Int64] = this.intComp

	this.comparators[reflect.Uint] = this.uintComp
	this.comparators[reflect.Uint32] = this.uintComp
	this.comparators[reflect.Uint64] = this.uintComp

	this.comparators[reflect.String] = this.stringComp

	this.comparators[reflect.Bool] = this.boolComp

	this.comparators[reflect.Float32] = this.floatComp
	this.comparators[reflect.Float64] = this.floatComp

	this.comparators[reflect.Ptr] = this.ptrComp

	this.comparators[reflect.Struct] = this.structComp

	this.comparators[reflect.Slice] = this.sliceComp

	this.comparators[reflect.Map] = this.mapComp

}

func (this *DeepEqual) Equal(aSide, zSide interface{}) bool {
	aSideValue := reflect.ValueOf(aSide)
	zSideValue := reflect.ValueOf(zSide)
	return this.equal(aSideValue, zSideValue)
}

func (this *DeepEqual) equal(aSideValue, zSideValue reflect.Value) bool {
	if aSideValue.IsValid() && !zSideValue.IsValid() {
		return false
	}
	if !aSideValue.IsValid() && zSideValue.IsValid() {
		return false
	}
	if !aSideValue.IsValid() && !zSideValue.IsValid() {
		return true
	}
	if aSideValue.Kind() != zSideValue.Kind() {
		return false
	}

	kind := aSideValue.Kind()
	comparator := this.comparators[kind]
	if comparator == nil {
		panic("No comparator for kind:" + kind.String() + ", please add it!")
	}
	return comparator(aSideValue, zSideValue)
}

//---------------------------------------------------------------------------

func (this *DeepEqual) intComp(aSideValue, zSideValue reflect.Value) bool {
	return aSideValue.Int() == zSideValue.Int()
}

func (this *DeepEqual) uintComp(aSideValue, zSideValue reflect.Value) bool {
	return aSideValue.Uint() == zSideValue.Uint()
}

func (this *DeepEqual) stringComp(aSideValue, zSideValue reflect.Value) bool {
	return aSideValue.String() == zSideValue.String()
}

func (this *DeepEqual) boolComp(aSideValue, zSideValue reflect.Value) bool {
	return aSideValue.Bool() == zSideValue.Bool()
}

func (this *DeepEqual) floatComp(aSideValue, zSideValue reflect.Value) bool {
	return aSideValue.Float() == zSideValue.Float()
}

func (this *DeepEqual) ptrComp(aSideValue, zSideValue reflect.Value) bool {
	if aSideValue.IsNil() && !zSideValue.IsNil() {
		return false
	}
	if !aSideValue.IsNil() && zSideValue.IsNil() {
		return false
	}
	if aSideValue.IsNil() && zSideValue.IsNil() {
		return true
	}
	return this.equal(aSideValue.Elem(), zSideValue.Elem())
}

func (this *DeepEqual) structComp(aSideValue, zSideValue reflect.Value) bool {
	if aSideValue.Type().Name() != zSideValue.Type().Name() {
		return false
	}
	for i := 0; i < aSideValue.Type().NumField(); i++ {
		fieldName := aSideValue.Type().Field(i).Name
		if SkipFieldByName(fieldName) {
			continue
		}
		aFieldValue := aSideValue.Field(i)
		zFieldValue := zSideValue.Field(i)
		eq := this.equal(aFieldValue, zFieldValue)
		if !eq {
			return false
		}
	}
	return true
}

func (this *DeepEqual) sliceComp(aSideValue, zSideValue reflect.Value) bool {
	if aSideValue.IsNil() && !zSideValue.IsNil() {
		return false
	}
	if !aSideValue.IsNil() && zSideValue.IsNil() {
		return false
	}
	if aSideValue.IsNil() && zSideValue.IsNil() {
		return true
	}

	if aSideValue.Len() != zSideValue.Len() {
		return false
	}

	for i := 0; i < aSideValue.Len(); i++ {
		aSideCel := aSideValue.Index(i)
		zSideCel := zSideValue.Index(i)
		eq := this.equal(aSideCel, zSideCel)
		if !eq {
			return false
		}
	}
	return true
}

func (this *DeepEqual) mapComp(aSideValue, zSideValue reflect.Value) bool {
	if aSideValue.IsNil() && !zSideValue.IsNil() {
		return false
	}
	if !aSideValue.IsNil() && zSideValue.IsNil() {
		return false
	}
	if aSideValue.IsNil() && zSideValue.IsNil() {
		return true
	}
	mapKeysAside := aSideValue.MapKeys()
	mapKeysZside := zSideValue.MapKeys()

	if len(mapKeysAside) != len(mapKeysZside) {
		return false
	}

	for _, key := range mapKeysAside {
		aSideV := aSideValue.MapIndex(key)
		zSideV := zSideValue.MapIndex(key)
		eq := this.equal(aSideV, zSideV)
		if !eq {
			return false
		}
	}
	return true
}
