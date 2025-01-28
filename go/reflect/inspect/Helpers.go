package inspect

import (
	"github.com/saichler/reflect/go/reflect/common"
	"github.com/saichler/reflect/go/types"
	"reflect"
)

func (this *Introspector) addAttribute(node *types.RNode, _type reflect.Type, _fieldName string) *types.RNode {
	this.registry.RegisterType(_type)
	if node != nil && node.Attributes == nil {
		node.Attributes = make(map[string]*types.RNode)
	}

	subNode := &types.RNode{}
	subNode.TypeName = _type.Name()
	subNode.Parent = node
	subNode.FieldName = _fieldName

	if node != nil {
		node.Attributes[subNode.FieldName] = subNode
	}
	return subNode
}

func (this *Introspector) addNode(_type reflect.Type, _parent *types.RNode, _fieldName string) (*types.RNode, bool) {
	exist, ok := this.typeToNode.Get(_type.Name())
	if ok && !common.IsLeaf(exist) {
		clone := this.cloner.Clone(exist).(*types.RNode)
		clone.Parent = _parent
		clone.FieldName = _fieldName
		clone.CachedKey = ""
		nodePath := common.InspectNodeKey(clone)
		this.pathToNode.Put(nodePath, clone)
		return clone, true
	}

	node := this.addAttribute(_parent, _type, _fieldName)
	nodePath := common.InspectNodeKey(node)
	_, ok = this.pathToNode.Get(nodePath)
	if ok {
		return nil, false
	}
	this.pathToNode.Put(nodePath, node)
	if _type.Kind() == reflect.Struct {
		this.typeToNode.Put(node.TypeName, node)
	}
	return node, false
}

func (this *Introspector) inspectStruct(_type reflect.Type, _parent *types.RNode, _fieldName string) *types.RNode {
	localNode, isClone := this.addNode(_type, _parent, _fieldName)
	if isClone {
		return localNode
	}
	this.registry.RegisterType(_type)
	for index := 0; index < _type.NumField(); index++ {
		field := _type.Field(index)
		if common.IgnoreName(field.Name) {
			continue
		}
		if field.Type.Kind() == reflect.Slice {
			this.inspectSlice(field.Type, localNode, field.Name)
		} else if field.Type.Kind() == reflect.Map {
			this.inspectMap(field.Type, localNode, field.Name)
		} else if field.Type.Kind() == reflect.Ptr {
			subnode := this.inspectPtr(field.Type.Elem(), localNode, field.Name)
			this.typeToNode.Put(subnode.TypeName, subnode)
		} else {
			this.addNode(field.Type, localNode, field.Name)
		}
	}
	this.addTableView(localNode)
	return localNode
}

func (this *Introspector) inspectPtr(_type reflect.Type, _parent *types.RNode, _fieldName string) *types.RNode {
	switch _type.Kind() {
	case reflect.Struct:
		return this.inspectStruct(_type, _parent, _fieldName)
	}
	panic("unknown ptr kind " + _type.Kind().String())
}

func (this *Introspector) inspectMap(_type reflect.Type, _parent *types.RNode, _fieldName string) *types.RNode {
	if _type.Elem().Kind() == reflect.Ptr && _type.Elem().Elem().Kind() == reflect.Struct {
		subNode := this.inspectStruct(_type.Elem().Elem(), _parent, _fieldName)
		subNode.IsMap = true
		_parent.Attributes[_fieldName] = subNode
		return subNode
	} else {
		subNode, _ := this.addNode(_type.Elem(), _parent, _fieldName)
		subNode.IsMap = true
		return subNode
	}
}

func (this *Introspector) inspectSlice(_type reflect.Type, _parent *types.RNode, _fieldName string) *types.RNode {
	if _type.Elem().Kind() == reflect.Ptr && _type.Elem().Elem().Kind() == reflect.Struct {
		subNode := this.inspectStruct(_type.Elem().Elem(), _parent, _fieldName)
		subNode.IsSlice = true
		_parent.Attributes[_fieldName] = subNode
		return subNode
	} else {
		subNode, _ := this.addNode(_type.Elem(), _parent, _fieldName)
		subNode.IsSlice = true
		return subNode
	}
}
