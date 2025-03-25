package introspecting

import (
	"github.com/saichler/shared/go/share/strings"
	"github.com/saichler/types/go/types"
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

func AddDeepDecorator(rnode *types.RNode) {
	addDecorator(types.DecoratorType_Deep, "true", rnode)
}

func DeepDecorator(rnode *types.RNode) bool {
	dec, _ := decoratorOf(types.DecoratorType_Deep, rnode).(string)
	if dec != "" {
		return true
	}
	return false
}
