package introspecting

import (
	"reflect"

	"github.com/saichler/l8types/go/types/l8reflect"
	"github.com/saichler/l8utils/go/utils/maps"
)

var _node *l8reflect.L8Node
var _nodeType = reflect.TypeOf(_node)

type RNodeMap struct {
	impl *maps.SyncMap
}

func NewIntrospectNodeMap() *RNodeMap {
	m := &RNodeMap{}
	m.impl = maps.NewSyncMap()
	return m
}

func (this *RNodeMap) Put(key string, value *l8reflect.L8Node) bool {
	return this.impl.Put(key, value)
}

func (this *RNodeMap) Get(key string) (*l8reflect.L8Node, bool) {
	value, ok := this.impl.Get(key)
	if value != nil {
		return value.(*l8reflect.L8Node), ok
	}
	return nil, ok
}

func (this *RNodeMap) Contains(key string) bool {
	return this.impl.Contains(key)
}

func (this *RNodeMap) NodesList(filter func(v interface{}) bool) []*l8reflect.L8Node {
	return this.impl.ValuesAsList(_nodeType, filter).([]*l8reflect.L8Node)
}

func (this *RNodeMap) Iterate(do func(k, v interface{})) {
	this.impl.Iterate(do)
}

func (this *RNodeMap) Del(key string) bool {
	_, ok := this.impl.Delete(key)
	return ok
}
