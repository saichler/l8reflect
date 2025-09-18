package introspecting

import (
	"github.com/saichler/l8utils/go/utils/strings"
)

func addDecorator(decoratorType types.DecoratorType, any interface{}, node *types.RNode) {
	s := strings.New()
	s.TypesPrefix = true
	str := s.StringOf(any)
	if node.Decorators == nil {
		node.Decorators = make(map[int32]string)
	}
	node.Decorators[int32(decoratorType)] = str
}

func decoratorOf(decoratorType types.DecoratorType, node *types.RNode) interface{} {
	decValue := node.Decorators[int32(decoratorType)]
	v, err := strings.InstanceOf(decValue, nil)
	if err != nil {
		panic(err)
	}
	return v
}

func AddPrimaryKeyDecorator(rnode *types.RNode, fields ...string) {
	addDecorator(types.DecoratorType_Primary, fields, rnode)
}

func PrimaryKeyDecorator(rnode *types.RNode) interface{} {
	return decoratorOf(types.DecoratorType_Primary, rnode)
}

func AddNoNestedInspection(rnode *types.RNode) {
	addDecorator(types.DecoratorType_NoNestedInspection, "t", rnode)
}

func NoNestedInspection(rnode *types.RNode) bool {
	dec, _ := decoratorOf(types.DecoratorType_NoNestedInspection, rnode).(string)
	if dec != "" {
		return true
	}
	return false
}
