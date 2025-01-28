package common

import (
	"github.com/saichler/reflect/go/types"
	"github.com/saichler/shared/go/share/interfaces"
	"reflect"
)

type IIntrospect interface {
	Inspect(interface{}) (*types.RNode, error)
	Node(string) (*types.RNode, bool)
	NodeByType(p reflect.Type) (*types.RNode, bool)
	NodeByTypeName(string) (*types.RNode, bool)
	NodeByValue(interface{}) (*types.RNode, bool)
	Nodes(bool, bool) []*types.RNode
	Registry() interfaces.IRegistry
	Kind(*types.RNode) reflect.Kind
	Clone(interface{}) interface{}
	AddDecorator(types.DecoratorType, interface{}, *types.RNode)
	DecoratorOf(types.DecoratorType, *types.RNode) interface{}
	TableView(string) (*types.TableView, bool)
	TableViews() []*types.TableView
}
