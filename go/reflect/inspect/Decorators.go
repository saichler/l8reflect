package inspect

import (
	"github.com/saichler/reflect/go/types"
	"github.com/saichler/shared/go/share/strings"
)

func (this *Introspector) AddDecorator(decoratorType types.DecoratorType, any interface{}, node *types.RNode) {
	s := strings.New()
	s.TypesPrefix = true
	str := s.StringOf(any)
	if node.Decorators == nil {
		node.Decorators = make(map[int32]string)
	}
	node.Decorators[int32(decoratorType)] = str
}

func (this *Introspector) DecoratorOf(decoratorType types.DecoratorType, node *types.RNode) interface{} {
	decValue := node.Decorators[int32(decoratorType)]
	v, _ := strings.InstanceOf(decValue, this.registry)
	return v
}
