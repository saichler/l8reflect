package properties

import (
	"fmt"
	"github.com/saichler/types/go/common"
	"reflect"
)

func (this *Property) mapSet(myMapValue reflect.Value, newMapValue reflect.Value) (interface{}, error) {
	var vInfo common.IInfo
	var kInfo common.IInfo
	var err error

	vInfo, err = this.introspector.Registry().Info(this.node.TypeName)
	if err != nil {
		return nil, err
	}

	kInfo, err = this.introspector.Registry().Info(this.node.KeyTypeName)
	if err != nil {
		return nil, err
	}

	//create the map if it is nil
	if !myMapValue.IsValid() || myMapValue.IsNil() {
		if this.node.IsStruct {
			myMapValue.Set(reflect.MakeMap(reflect.MapOf(kInfo.Type(), reflect.PointerTo(vInfo.Type()))))
		} else {
			myMapValue.Set(reflect.MakeMap(reflect.MapOf(kInfo.Type(), vInfo.Type())))
		}
	}

	//This means the entire map is new
	if this.key == nil {
		if !newMapValue.IsValid() {
			myMapValue.SetZero()
		} else {
			myMapValue.Set(newMapValue)
		}
		return myMapValue.Interface(), nil
	}

	mapKey := reflect.ValueOf(this.key)
	oKeyValue := myMapValue.MapIndex(mapKey)
	//in this case, the newMapValue isn't a map, it is a value
	//this.value = newMapValue.Interface()
	nKeyValue := newMapValue

	//This map entry was marked for deletion so delete it
	if nKeyValue.Kind() == reflect.String && nKeyValue.String() == common.Deleted_Entry {
		myMapValue.SetMapIndex(mapKey, reflect.Value{})
		return myMapValue.Interface(), err
	}

	//If this node is a struct & this property is not the leaf
	//we need to return the old struct instance to keep drilling down to the updated property.
	if this.node.IsStruct && !this.IsLeaf() {
		//if the old value is not valid, create it
		if !oKeyValue.IsValid() {
			typeName := newMapValue.Type().Name()
			if newMapValue.Kind() == reflect.Ptr {
				typeName = newMapValue.Elem().Type().Name()
			}
			if typeName == vInfo.Type().Name() {
				myMapValue.SetMapIndex(mapKey, newMapValue)
				oKeyValue = newMapValue
			} else {
				fmt.Println(newMapValue.Interface(), newMapValue.Kind())
				o, _ := vInfo.NewInstance()
				oKeyValue = reflect.ValueOf(o)
				myMapValue.SetMapIndex(mapKey, oKeyValue)
			}
		}
		return oKeyValue.Interface(), nil
	}

	myMapValue.SetMapIndex(mapKey, nKeyValue)

	return myMapValue.Interface(), err
}
