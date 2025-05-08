package introspecting

import (
	"github.com/saichler/l8utils/go/utils/maps"
	"github.com/saichler/l8types/go/types"
	"reflect"
)

var _node *types.RNode
var _nodeType = reflect.TypeOf(_node)

type RNodeMap struct {
	impl *maps.SyncMap
}

func NewIntrospectNodeMap() *RNodeMap {
	m := &RNodeMap{}
	m.impl = maps.NewSyncMap()
	return m
}

func (this *RNodeMap) Put(key string, value *types.RNode) bool {
	return this.impl.Put(key, value)
}

func (this *RNodeMap) Get(key string) (*types.RNode, bool) {
	value, ok := this.impl.Get(key)
	if value != nil {
		return value.(*types.RNode), ok
	}
	return nil, ok
}

func (this *RNodeMap) Contains(key string) bool {
	return this.impl.Contains(key)
}

func (this *RNodeMap) NodesList(filter func(v interface{}) bool) []*types.RNode {
	return this.impl.ValuesAsList(_nodeType, filter).([]*types.RNode)
}

func (this *RNodeMap) Iterate(do func(k, v interface{})) {
	this.impl.Iterate(do)
}
