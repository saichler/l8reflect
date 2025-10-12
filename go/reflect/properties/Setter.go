package properties

import (
	"errors"
	"reflect"
	"strings"

	"github.com/saichler/l8reflect/go/reflect/introspecting"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8reflect"
	strings2 "github.com/saichler/l8utils/go/utils/strings"
)

func (this *Property) Set(any interface{}, value interface{}) (interface{}, interface{}, error) {
	if this == nil {
		return nil, nil, errors.New("property is nil, cannot instantiate")
	}
	if this.parent == nil {
		if any == nil {
			info, err := this.resources.Registry().Info(this.node.TypeName)
			if err != nil {
				return nil, nil, err
			}
			newAny, err := info.NewInstance()
			if err != nil {
				return nil, nil, err
			}
			any = newAny
		}
		if this.key != nil {
			this.SetPrimaryKey(this.node, any, this.key)
		}
		return any, any, nil
	}
	parent, root, err := this.parent.Set(any, value)
	if err != nil {
		return nil, nil, err
	}
	if any == nil {
		any = root
	}
	parentValue := reflect.ValueOf(parent)
	if parentValue.Kind() == reflect.Ptr {
		parentValue = parentValue.Elem()
	}

	//Special case for setting a value to the map
	if this.node.IsMap && parentValue.Kind() == reflect.Map {
		if this.IsLeaf() {
			parentValue.SetMapIndex(reflect.ValueOf(this.key), reflect.ValueOf(this.value))
		}
		return this.value, any, nil
	} else if parentValue.Kind() == reflect.Map {
		parentValue = parentValue.MapIndex(reflect.ValueOf(this.key))
	}

	//Special case where the model is setting the same reference
	//in different attributes, which is incorrect.
	if parentValue.Kind() == reflect.Slice {
		pid, _ := this.PropertyId()
		strValue, ok := value.(string)
		if ok && strValue == ifs.Deleted_Entry {
			this.resources.Logger().Error("The model contain same reference in a map and a slice, pid=" + pid)
			return nil, nil, nil
		}
	}

	myValue := parentValue.FieldByName(this.node.FieldName)
	info, err := this.resources.Registry().Info(this.node.TypeName)
	if err != nil {
		return nil, nil, err
	}
	typ := info.Type()
	if this.node.IsMap {
		v, e := this.mapSet(myValue, reflect.ValueOf(value))
		return v, any, e
	} else if this.node.IsSlice {
		v, e := this.sliceSet(myValue, reflect.ValueOf(value))
		return v, any, e
	} else if this.resources.Introspector().Kind(this.node) == reflect.Struct {
		// Handle setting to nil
		if value == nil {
			if myValue.IsValid() && myValue.CanSet() {
				myValue.Set(reflect.Zero(myValue.Type()))
			}
			return nil, any, err
		}

		if !myValue.IsValid() || myValue.IsNil() {
			v := reflect.ValueOf(value)
			if v.Kind() == reflect.Ptr &&
				!v.IsNil() && v.Elem().Type().Name() == typ.Name() {
				myValue.Set(reflect.ValueOf(value))
			} else {
				newInstance := reflect.New(typ)
				if v.Kind() == reflect.String {
					serializer := info.Serializer(ifs.STRING)
					if serializer != nil {
						inst, _ := serializer.Unmarshal([]byte(v.String()), this.Resources())
						if inst != nil {
							newInstance = reflect.ValueOf(inst)
						}
					}
				}
				if myValue.CanSet() {
					myValue.Set(newInstance)
				} else {
					p, _ := this.PropertyId()
					return nil, any, errors.New("Cannot set value to " + p)
				}
			}
		} else {
			// Handle replacing existing struct pointer with new value
			v := reflect.ValueOf(value)
			if v.Kind() == reflect.Ptr &&
				!v.IsNil() && v.Elem().Type().Name() == typ.Name() {
				myValue.Set(reflect.ValueOf(value))
			}
		}
		return myValue.Interface(), any, err
	} else if reflect.ValueOf(value).Kind() == reflect.Int32 || myValue.Kind() == reflect.Int32 {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.String {
			value = this.resources.Registry().Enum(value.(string))
		}
		myValue.SetInt(reflect.ValueOf(value).Int())
		return value, any, err
	} else {
		if value != nil {
			v := reflect.ValueOf(value)
			if v.Kind() != myValue.Kind() {
				v = ConvertValue(myValue, v)
			}
			myValue.Set(v)
		}
		return value, any, err
	}
}

func (this *Property) SetPrimaryKey(node *l8reflect.L8Node, any interface{}, anyKey interface{}) {
	if anyKey == nil {
		return
	}
	keyString := anyKey.(string)
	tokens := strings.Split(keyString, "::")
	fieldsValues := make([]interface{}, len(tokens))
	for i, token := range tokens {
		vv, _ := strings2.FromString(token, nil)
		fieldsValues[i] = vv.Interface()
	}

	value := reflect.ValueOf(any)
	if !value.IsValid() {
		return
	}
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return
		}
		value = value.Elem()
	}

	f := introspecting.PrimaryKeyDecorator(node)
	if f != nil {
		fields := f.([]string)
		for i, attr := range fields {
			fld := value.FieldByName(attr)
			fld.Set(reflect.ValueOf(fieldsValues[i]))
		}
	}
}
