package properties

import (
	"errors"
	"reflect"

	"github.com/saichler/l8types/go/ifs"
)

func (this *Property) mapSet(myMapValue reflect.Value, newMapValue reflect.Value) (interface{}, error) {
	var vInfo ifs.IInfo
	var kInfo ifs.IInfo
	var err error

	vInfo, err = this.resources.Registry().Info(this.node.TypeName)
	if err != nil {
		return nil, err
	}

	kInfo, err = this.resources.Registry().Info(this.node.KeyTypeName)
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
			if newMapValue.Kind() != reflect.Map {
				return nil, errors.New("invalid map type " + newMapValue.Kind().String() + " for map " + myMapValue.Type().String())
			}
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
	if this.isLeaf && nKeyValue.Kind() == reflect.String && nKeyValue.String() == ifs.Deleted_Entry {
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
