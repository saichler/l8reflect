package properties

import (
	"reflect"

	"github.com/saichler/l8types/go/ifs"
)

func (this *Property) sliceSet(myValue reflect.Value, newSliceValue reflect.Value) (interface{}, error) {
	//Replace all the slice
	if this.key == nil {
		// Handle setting slice to nil or a new slice
		if newSliceValue.Kind() == reflect.Slice || !newSliceValue.IsValid() {
			// Check if myValue is valid and settable
			if myValue.IsValid() && myValue.CanSet() {
				if !newSliceValue.IsValid() {
					// Setting to nil - create a zero value of the appropriate slice type
					sliceType := myValue.Type()
					nilSlice := reflect.Zero(sliceType)
					myValue.Set(nilSlice)
					return nil, nil
				} else {
					myValue.Set(newSliceValue)
					return myValue.Interface(), nil
				}
			} else {
				// If we can't set the value, just return the new value
				if newSliceValue.IsValid() {
					return newSliceValue.Interface(), nil
				}
				return nil, nil
			}
		}
	}

	// Check if this.key is nil before casting
	if this.key == nil {
		return nil, nil // Return nil for setting nil on slice without index
	}

	index := this.key.(int)
	info, err := this.resources.Registry().Info(this.node.TypeName)
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

	if this.node.IsStruct && (!oIndexValue.IsValid() || oIndexValue.IsNil()) {
		oIndexValue.Set(reflect.New(info.Type()))
	}

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
