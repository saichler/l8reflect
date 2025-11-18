package properties

import (
	"reflect"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8reflect"
)

func Collect(root interface{}, r ifs.IResources, typeName string) []interface{} {
	typ := reflect.ValueOf(root).Elem().Type()
	node, _ := r.Introspector().NodeByType(typ)
	list := make([]interface{}, 0)
	collect(root, node, typeName, &list)
	return list
}

func collect(any interface{}, node *l8reflect.L8Node, typeName string, list *[]interface{}) {
	val := reflect.ValueOf(any)
	if val.Kind() != reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Elem().Type()
	if typ.Name() == typeName {
		*list = append(*list, any)
		return
	}
	if node.Attributes != nil {
		for _, attr := range node.Attributes {
			if attr.IsStruct {
				value := val.FieldByName(attr.FieldName)
				if value.IsValid() {
					collect(value.Interface(), attr, attr.TypeName, list)
				}
			} else if attr.IsSlice {
				value := val.FieldByName(attr.FieldName)
				if value.IsValid() {
					for i := 0; i < value.Len(); i++ {
						collect(value.Index(i).Interface(), attr, attr.TypeName, list)
					}
				}
			} else if attr.IsMap {
				value := val.FieldByName(attr.FieldName)
				if value.IsValid() {
					keys := value.MapKeys()
					for i := 0; i < len(keys); i++ {
						collect(value.MapIndex(keys[i]).Interface(), attr, attr.TypeName, list)
					}
				}
			}
		}
	}
}
