package helping

import (
	strings2 "github.com/saichler/l8utils/go/utils/strings"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types"
	"reflect"
	"strings"
)

func ValueAndType(any interface{}) (reflect.Value, reflect.Type) {
	v := reflect.ValueOf(any)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	return v, t
}

func IsLeaf(node *types.RNode) bool {
	if node.Attributes == nil || len(node.Attributes) == 0 {
		return true
	}
	return false
}

func IsRoot(node *types.RNode) bool {
	if node.Parent == nil {
		return true
	}
	return false
}

func IgnoreName(fieldName string) bool {
	if fieldName == "DoNotCompare" {
		return true
	}
	if fieldName == "DoNotCopy" {
		return true
	}
	if len(fieldName) > 3 && fieldName[0:3] == "XXX" {
		return true
	}
	if fieldName[0:1] == strings.ToLower(fieldName[0:1]) {
		return true
	}
	return false
}

func PropertyNodeKey(instanceId string) string {
	buff := strings2.New()
	open := false
	for _, c := range instanceId {
		if c == '<' {
			open = true
		} else if c == '>' {
			open = false
		} else if !open {
			buff.Add(string(c))
		}
	}
	return buff.String()
}

func InspectNodeKey(node *types.RNode) string {
	if node.CachedKey != "" {
		return node.CachedKey
	}
	if node.Parent == nil {
		return strings.ToLower(node.TypeName)
	}
	buff := strings2.New()
	buff.Add(InspectNodeKey(node.Parent))
	buff.Add(".")
	buff.Add(strings.ToLower(node.FieldName))
	node.CachedKey = buff.String()
	return node.CachedKey
}

func PrimaryDecorator(node *types.RNode, value reflect.Value, registry ifs.IRegistry) interface{} {
	fields := PrimaryDecoratorFields(node, registry)
	if fields == nil || len(fields) == 0 {
		return nil
	}
	str := strings2.New()
	for _, field := range fields {
		v := value.FieldByName(field).Interface()
		v2 := str.StringOf(v)
		str.Add(v2)
	}
	return str.String()
}

func PrimaryDecoratorFields(node *types.RNode, registry ifs.IRegistry) []string {
	decValue := node.Decorators[int32(types.DecoratorType_Primary)]
	v, _ := strings2.InstanceOf(decValue, registry)
	fields, ok := v.([]string)
	if !ok {
		return nil
	}
	return fields
}
