package properties

import (
	"github.com/saichler/l8types/go/ifs"
	"reflect"
)

func (this *Property) sliceSet(myValue reflect.Value, newSliceValue reflect.Value) (interface{}, error) {
	//Replace all the slice
	if this.key == nil || !myValue.IsValid() || myValue.IsNil() || !newSliceValue.IsValid() {
		myValue.Set(newSliceValue)
		return myValue.Interface(), nil
	}

	index := this.key.(int)
	info, err := this.introspector.Registry().Info(this.node.TypeName)
	if err != nil {
		return nil, err
	}

	//If this is a new slice
	if !myValue.IsValid() || myValue.IsNil() {
		if this.node.IsStruct {
			myValue.Set(reflect.MakeSlice(reflect.SliceOf(reflect.PointerTo(info.Type())), index+1, index+1))
		} else {
			myValue.Set(reflect.MakeSlice(reflect.SliceOf(info.Type()), index+1, index+1))
		}
	}

	//If elements were delete from the slice,
	//reduce the size of the slice
	if newSliceValue.Kind() == reflect.String && newSliceValue.String() == ifs.Deleted_Entry {
		var newSlice reflect.Value
		if this.node.IsStruct {
			newSlice = reflect.MakeSlice(reflect.SliceOf(reflect.PointerTo(info.Type())), index, index)
		} else {
			newSlice = reflect.MakeSlice(reflect.SliceOf(info.Type()), index, index)
		}
		for i := 0; i < index; i++ {
			newSlice.Index(i).Set(myValue.Index(i))
		}
		myValue.Set(newSlice)
		return myValue.Interface(), nil
	}

	//If the index is larger than the current slice, enlarge it
	if index >= myValue.Len() {
		var newSlice reflect.Value
		if this.node.IsStruct {
			newSlice = reflect.MakeSlice(reflect.SliceOf(reflect.PointerTo(info.Type())), index+1, index+1)
		} else {
			newSlice = reflect.MakeSlice(reflect.SliceOf(info.Type()), index+1, index+1)
		}
		for i := 0; i < myValue.Len(); i++ {
			newSlice.Index(i).Set(myValue.Index(i))
		}
		myValue.Set(newSlice)
	}

	oIndexValue := myValue.Index(index)

	if this.node.IsStruct && !this.IsLeaf() {
		return oIndexValue.Interface(), nil
	}

	nIndexValue := newSliceValue.Index(index)

	//If this is not a leaf property
	//We need to continue drilling down
	if this.node.IsStruct && !this.IsLeaf() {
		if !oIndexValue.IsValid() {
			o, _ := info.NewInstance()
			oIndexValue.Set(reflect.ValueOf(o))
		}
		return oIndexValue.Interface(), nil
	}

	oIndexValue.Set(nIndexValue)

	return oIndexValue.Interface(), err
}
