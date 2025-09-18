package introspecting

import (
	"github.com/saichler/l8types/go/types/l8reflect"
	"github.com/saichler/l8utils/go/utils/strings"
)

func addDecorator(decoratorType l8reflect.L8DecoratorType, any interface{}, node *l8reflect.L8Node) {
	s := strings.New()
	s.TypesPrefix = true
	str := s.StringOf(any)
	if node.Decorators == nil {
		node.Decorators = make(map[int32]string)
	}
	node.Decorators[int32(decoratorType)] = str
}

func decoratorOf(decoratorType l8reflect.L8DecoratorType, node *l8reflect.L8Node) interface{} {
	decValue := node.Decorators[int32(decoratorType)]
	v, err := strings.InstanceOf(decValue, nil)
	if err != nil {
		panic(err)
	}
	return v
}

func AddPrimaryKeyDecorator(rnode *l8reflect.L8Node, fields ...string) {
	addDecorator(l8reflect.L8DecoratorType_Primary, fields, rnode)
}

func PrimaryKeyDecorator(rnode *l8reflect.L8Node) interface{} {
	return decoratorOf(l8reflect.L8DecoratorType_Primary, rnode)
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
