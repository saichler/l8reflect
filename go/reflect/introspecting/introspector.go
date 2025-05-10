package introspecting

import (
	"errors"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types"
	"github.com/saichler/l8utils/go/utils/maps"
	"github.com/saichler/reflect/go/reflect/cloning"
	"github.com/saichler/reflect/go/reflect/helping"
	"reflect"
	"strings"
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

func (this *Introspector) Inspect(any interface{}) (*types.RNode, error) {
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

func (this *Introspector) Node(path string) (*types.RNode, bool) {
	return this.pathToNode.Get(strings.ToLower(path))
}

func (this *Introspector) NodeByValue(any interface{}) (*types.RNode, bool) {
	val := reflect.ValueOf(any)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return this.NodeByType(val.Type())
}

func (this *Introspector) NodeByType(typ reflect.Type) (*types.RNode, bool) {
	return this.NodeByTypeName(typ.Name())
}

func (this *Introspector) NodeByTypeName(name string) (*types.RNode, bool) {
	return this.typeToNode.Get(name)
}

func (this *Introspector) Nodes(onlyLeafs, onlyRoots bool) []*types.RNode {
	filter := func(any interface{}) bool {
		n := any.(*types.RNode)
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

func (this *Introspector) Kind(node *types.RNode) reflect.Kind {
	info, err := this.registry.Info(node.TypeName)
	if err != nil {
		panic(err.Error())
	}
	return info.Type().Kind()
}

func (this *Introspector) Clone(any interface{}) interface{} {
	return this.cloner.Clone(any)
}

func (this *Introspector) addTableView(node *types.RNode) {
	tv := &types.TableView{Table: node, Columns: make([]*types.RNode, 0), SubTables: make([]*types.RNode, 0)}
	for _, attr := range node.Attributes {
		if helping.IsLeaf(attr) {
			tv.Columns = append(tv.Columns, attr)
		} else {
			tv.SubTables = append(tv.SubTables, attr)
		}
	}
	this.tableViews.Put(node.TypeName, tv)
}

func (this *Introspector) TableView(name string) (*types.TableView, bool) {
	tv, ok := this.tableViews.Get(name)
	if !ok {
		return nil, ok
	}
	return tv.(*types.TableView), ok
}

func (this *Introspector) TableViews() []*types.TableView {
	list := this.tableViews.ValuesAsList(reflect.TypeOf(&types.TableView{}), nil)
	return list.([]*types.TableView)
}

func (this *Introspector) Clean(typeName string) {
	node, ok := this.NodeByTypeName(typeName)
	if !ok {
		return
	}
	this.clean(node)
}

func (this *Introspector) clean(node *types.RNode) {
	if node.Attributes != nil {
		for _, attr := range node.Attributes {
			this.clean(attr)
		}
	}
	this.typeToNode.Del(node.TypeName)
	this.pathToNode.Del(helping.NodeCacheKey(node))
}
