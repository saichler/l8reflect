package introspecting

import (
	"github.com/saichler/reflect/go/reflect/helping"
	"github.com/saichler/types/go/types"
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

func (this *Introspector) fixClone(clone *types.RNode, parent *types.RNode, fieldName string) {
	clone.Parent = parent
	clone.FieldName = fieldName
	clone.CachedKey = ""
	nodePath := helping.InspectNodeKey(clone)
	this.pathToNode.Put(nodePath, clone)
	if clone.Attributes != nil {
		for k, v := range clone.Attributes {
			this.fixClone(v, clone, k)
		}
	}
}

func (this *Introspector) addNode(_type reflect.Type, _parent *types.RNode, _fieldName string) (*types.RNode, bool) {
	exist, ok := this.typeToNode.Get(_type.Name())
	if ok && !helping.IsLeaf(exist) {
		clone := this.cloner.Clone(exist).(*types.RNode)
		this.fixClone(clone, _parent, _fieldName)
		return clone, true
	}

	node := this.addAttribute(_parent, _type, _fieldName)
	nodePath := helping.InspectNodeKey(node)
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
	localNode.IsStruct = true
	this.registry.RegisterType(_type)
	for index := 0; index < _type.NumField(); index++ {
		field := _type.Field(index)
		if helping.IgnoreName(field.Name) {
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
		subNode.IsStruct = true
		subNode.KeyTypeName = _type.Key().Name()
		if _parent.Attributes == nil {
			_parent.Attributes = make(map[string]*types.RNode)
		}
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
		subNode.IsStruct = true
		if _parent.Attributes == nil {
			_parent.Attributes = make(map[string]*types.RNode)
		}
		_parent.Attributes[_fieldName] = subNode
		return subNode
	} else {
		subNode, _ := this.addNode(_type.Elem(), _parent, _fieldName)
		subNode.IsSlice = true
		return subNode
	}
}
