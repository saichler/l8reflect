package properties

import (
	"reflect"

	"github.com/saichler/l8reflect/go/reflect/helping"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8reflect"
)

func Collect(root interface{}, r ifs.IResources, typeName string) map[string]interface{} {
	typ := reflect.ValueOf(root).Elem().Type()
	node, _ := r.Introspector().NodeByType(typ)
	result := make(map[string]interface{}, 0)
	rootKey := helping.PrimaryKeyDecoratorValue(node, reflect.ValueOf(root).Elem(), r.Registry())
	collect(root, node, typeName, nil, rootKey, result, r)
	return result
}

func collect(any interface{}, node *l8reflect.L8Node, typeName string,
	parent *Property, key interface{}, elems map[string]interface{}, r ifs.IResources) {
	if any == nil {
		return
	}
	val := reflect.ValueOf(any)
	if !val.IsValid() {
		return
	}
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}
	typ := val.Type()
	myProperty := NewProperty(node, parent, key, any, r)
	if typ.Name() == typeName {
		id, _ := myProperty.PropertyId()
		elems[id] = any
		return
	}

	if node.Attributes != nil {
		for _, attr := range node.Attributes {
			if attr.IsMap {
				value := val.FieldByName(attr.FieldName)
				if value.IsValid() {
					keys := value.MapKeys()
					for i := 0; i < len(keys); i++ {
						collect(value.MapIndex(keys[i]).Interface(), attr, typeName, myProperty, keys[i].Interface(), elems, r)
					}
				}
			} else if attr.IsSlice {
				value := val.FieldByName(attr.FieldName)
				if value.IsValid() {
					for i := 0; i < value.Len(); i++ {
						collect(value.Index(i).Interface(), attr, typeName, myProperty, i, elems, r)
					}
				}
			} else if attr.IsStruct {
				value := val.FieldByName(attr.FieldName)
				if value.IsValid() {
					collect(value.Interface(), attr, typeName, myProperty, nil, elems, r)
				}
			}
		}
	}
}
