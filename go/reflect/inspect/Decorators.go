package inspect

import (
	"github.com/saichler/reflect/go/types"
	"github.com/saichler/shared/go/share/string_utils"
)

func (this *Introspector) AddDecorator(decoratorType types.DecoratorType, any interface{}, node *types.RNode) {
	s := string_utils.New()
	s.TypesPrefix = true
	str := s.StringOf(any)
	if node.Decorators == nil {
		node.Decorators = make(map[int32]string)
	}
	node.Decorators[int32(decoratorType)] = str
}

func (this *Introspector) DecoratorOf(decoratorType types.DecoratorType, node *types.RNode) interface{} {
	decValue := node.Decorators[int32(decoratorType)]
	v := string_utils.InstanceOf(decValue, this.registry)
	return v
}
