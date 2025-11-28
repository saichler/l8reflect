package helping

import (
	"reflect"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8reflect"
	strings2 "github.com/saichler/l8utils/go/utils/strings"
)

func PrimaryKeyDecoratorValue(node *l8reflect.L8Node, value reflect.Value, registry ifs.IRegistry) interface{} {
	return decorator(node, l8reflect.L8DecoratorType_Primary, value, registry)
}

func UniqueKeyDecoratorValue(node *l8reflect.L8Node, value reflect.Value, registry ifs.IRegistry) interface{} {
	return decorator(node, l8reflect.L8DecoratorType_Unique, value, registry)
}

func decorator(node *l8reflect.L8Node, decoratorType l8reflect.L8DecoratorType, value reflect.Value, registry ifs.IRegistry) interface{} {
	fields := decoratorFields(node, decoratorType, registry)
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
	return decoratorFields(node, l8reflect.L8DecoratorType_Primary, registry)
}

func UniqueDecoratorFields(node *l8reflect.L8Node, registry ifs.IRegistry) []string {
	return decoratorFields(node, l8reflect.L8DecoratorType_Unique, registry)
}

func decoratorFields(node *l8reflect.L8Node, decoratorType l8reflect.L8DecoratorType, registry ifs.IRegistry) []string {
	decValue := node.Decorators[int32(decoratorType)]
	v, _ := strings2.InstanceOf(decValue, registry)
	fields, ok := v.([]string)
	if !ok {
		return nil
	}
	return fields
}

func addDecorator(decoratorType l8reflect.L8DecoratorType, any interface{}, node *l8reflect.L8Node) {
	s := strings2.New()
	s.TypesPrefix = true
	str := s.StringOf(any)
	if node.Decorators == nil {
		node.Decorators = make(map[int32]string)
	}
	node.Decorators[int32(decoratorType)] = str
}

func decoratorOf(decoratorType l8reflect.L8DecoratorType, node *l8reflect.L8Node) interface{} {
	decValue := node.Decorators[int32(decoratorType)]
	v, err := strings2.InstanceOf(decValue, nil)
	if err != nil {
		panic(err)
	}
	return v
}

func PrimaryKeyDecorator(rnode *l8reflect.L8Node) interface{} {
	return decoratorOf(l8reflect.L8DecoratorType_Primary, rnode)
}

func UniqueKeyDecorator(rnode *l8reflect.L8Node) interface{} {
	return decoratorOf(l8reflect.L8DecoratorType_Unique, rnode)
}

func AddPrimaryKeyDecorator(rnode *l8reflect.L8Node, fields ...string) {
	addDecorator(l8reflect.L8DecoratorType_Primary, fields, rnode)
}

func AddUniqueKeyDecorator(rnode *l8reflect.L8Node, fields ...string) {
	addDecorator(l8reflect.L8DecoratorType_Unique, fields, rnode)
}

func AddNoNestedInspection(rnode *l8reflect.L8Node) {
	addDecorator(l8reflect.L8DecoratorType_NoNestedInspection, "t", rnode)
}

func NoNestedInspection(rnode *l8reflect.L8Node) bool {
	dec, _ := decoratorOf(l8reflect.L8DecoratorType_NoNestedInspection, rnode).(string)
	if dec != "" {
		return true
	}
	return false
}

func AddAlwayOverwriteDecorator(rnode *l8reflect.L8Node) {
	addDecorator(l8reflect.L8DecoratorType_AlwaysFull, true, rnode)
}

func AlwaysFullDecorator(rnode *l8reflect.L8Node) bool {
	val := decoratorOf(l8reflect.L8DecoratorType_AlwaysFull, rnode)
	b, ok := val.(bool)
	if !ok {
		return false
	}
	return b
}
