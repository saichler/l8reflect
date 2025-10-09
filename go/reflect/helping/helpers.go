package helping

import (
	"reflect"
	"strings"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8reflect"
	strings2 "github.com/saichler/l8utils/go/utils/strings"
)

func ValueAndType(any interface{}) (reflect.Value, reflect.Type) {
	v := reflect.ValueOf(any)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	return v, t
}

func IsLeaf(node *l8reflect.L8Node) bool {
	if node.Attributes == nil || len(node.Attributes) == 0 {
		return true
	}
	return false
}

func IsRoot(node *l8reflect.L8Node) bool {
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

func NodeCacheKey(node *l8reflect.L8Node) string {
	if node.CachedKey != "" {
		return node.CachedKey
	}
	if node.Parent == nil {
		return strings.ToLower(node.TypeName)
	}
	buff := strings2.New()
	buff.Add(NodeCacheKey(node.Parent))
	buff.Add(".")
	buff.Add(strings.ToLower(node.FieldName))
	node.CachedKey = buff.String()
	return node.CachedKey
}

func PrimaryDecorator(node *l8reflect.L8Node, value reflect.Value, registry ifs.IRegistry) interface{} {
	fields := PrimaryDecoratorFields(node, registry)
	if fields == nil || len(fields) == 0 {
		return nil
	}
	str := strings2.New()
	str.TypesPrefix = true
	first := true
	for _, field := range fields {
		if !first {
			str.TypesPrefix = false
			str.Add("::")
			str.TypesPrefix = true
		}
		v := value.FieldByName(field).Interface()
		v2 := str.StringOf(v)
		str.Add(v2)
		first = false
	}
	return str.String()
}

func PrimaryDecoratorFields(node *l8reflect.L8Node, registry ifs.IRegistry) []string {
	decValue := node.Decorators[int32(l8reflect.L8DecoratorType_Primary)]
	v, _ := strings2.InstanceOf(decValue, registry)
	fields, ok := v.([]string)
	if !ok {
		return nil
	}
	return fields
}
