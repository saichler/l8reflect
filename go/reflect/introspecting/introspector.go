package introspecting

import (
	"errors"
	"reflect"
	"strings"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8reflect"
	"github.com/saichler/l8utils/go/utils/maps"
	"github.com/saichler/reflect/go/reflect/cloning"
	"github.com/saichler/reflect/go/reflect/helping"
)

type Introspector struct {
	pathToNode *RNodeMap
	typeToNode *RNodeMap
	registry   ifs.IRegistry
	cloner     *cloning.Cloner
	tableViews *maps.SyncMap
}

func NewIntrospect(registry ifs.IRegistry) *Introspector {
	instrospector := &Introspector{}
	instrospector.registry = registry
	instrospector.cloner = cloning.NewCloner()
	instrospector.pathToNode = NewIntrospectNodeMap()
	instrospector.typeToNode = NewIntrospectNodeMap()
	instrospector.tableViews = maps.NewSyncMap()
	return instrospector
}

func (this *Introspector) Registry() ifs.IRegistry {
	return this.registry
}

func (this *Introspector) Inspect(any interface{}) (*l8reflect.L8Node, error) {
	if any == nil {
		return nil, errors.New("Cannot introspect a nil value")
	}

	_, t := helping.ValueAndType(any)
	if t.Kind() == reflect.Slice && t.Kind() == reflect.Map {
		t = t.Elem().Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, errors.New("Cannot introspect a value that is not a struct")
	}
	localNode, ok := this.pathToNode.Get(strings.ToLower(t.Name()))
	if ok {
		return localNode, nil
	}
	return this.inspectStruct(t, nil, ""), nil
}

func (this *Introspector) Node(path string) (*l8reflect.L8Node, bool) {
	return this.pathToNode.Get(strings.ToLower(path))
}

func (this *Introspector) NodeByValue(any interface{}) (*l8reflect.L8Node, bool) {
	val := reflect.ValueOf(any)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return this.NodeByType(val.Type())
}

func (this *Introspector) NodeByType(typ reflect.Type) (*l8reflect.L8Node, bool) {
	return this.NodeByTypeName(typ.Name())
}

func (this *Introspector) NodeByTypeName(name string) (*l8reflect.L8Node, bool) {
	return this.typeToNode.Get(name)
}

func (this *Introspector) Nodes(onlyLeafs, onlyRoots bool) []*l8reflect.L8Node {
	filter := func(any interface{}) bool {
		n := any.(*l8reflect.L8Node)
		if onlyLeafs && !helping.IsLeaf(n) {
			return false
		}
		if onlyRoots && !helping.IsRoot(n) {
			return false
		}
		return true
	}

	return this.pathToNode.NodesList(filter)
}

func (this *Introspector) Kind(node *l8reflect.L8Node) reflect.Kind {
	info, err := this.registry.Info(node.TypeName)
	if err != nil {
		panic(err.Error())
	}
	return info.Type().Kind()
}

func (this *Introspector) Clone(any interface{}) interface{} {
	return this.cloner.Clone(any)
}

func (this *Introspector) addTableView(node *l8reflect.L8Node) {
	tv := &l8reflect.L8TableView{Table: node, Columns: make([]*l8reflect.L8Node, 0), SubTables: make([]*l8reflect.L8Node, 0)}
	for _, attr := range node.Attributes {
		if helping.IsLeaf(attr) {
			tv.Columns = append(tv.Columns, attr)
		} else {
			tv.SubTables = append(tv.SubTables, attr)
		}
	}
	this.tableViews.Put(node.TypeName, tv)
}

func (this *Introspector) TableView(name string) (*l8reflect.L8TableView, bool) {
	tv, ok := this.tableViews.Get(name)
	if !ok {
		return nil, ok
	}
	return tv.(*l8reflect.L8TableView), ok
}

func (this *Introspector) TableViews() []*l8reflect.L8TableView {
	list := this.tableViews.ValuesAsList(reflect.TypeOf(&l8reflect.L8TableView{}), nil)
	return list.([]*l8reflect.L8TableView)
}

func (this *Introspector) Clean(typeName string) {
	node, ok := this.NodeByTypeName(typeName)
	if !ok {
		return
	}
	this.clean(node)
}

func (this *Introspector) clean(node *l8reflect.L8Node) {
	if node.Attributes != nil {
		for _, attr := range node.Attributes {
			this.clean(attr)
		}
	}
	this.typeToNode.Del(node.TypeName)
	this.pathToNode.Del(helping.NodeCacheKey(node))
}
