package property

import (
	"errors"
	"github.com/saichler/reflect/go/reflect/common"
	"github.com/saichler/reflect/go/types"
	"github.com/saichler/shared/go/share/string_utils"
	"strings"
)

type Property struct {
	parent       *Property
	node         *types.RNode
	key          interface{}
	value        interface{}
	id           string
	introspector common.IIntrospect
}

func NewProperty(node *types.RNode, parent *Property, key interface{}, value interface{}, introspector common.IIntrospect) *Property {
	property := &Property{}
	property.parent = parent
	property.node = node
	property.key = key
	property.value = value
	property.introspector = introspector
	return property
}

func PropertyOf(propertyId string, introspector common.IIntrospect) (*Property, error) {
	propertyKey := common.NodeKey(propertyId)
	node, ok := introspector.Node(propertyKey)
	if !ok {
		return nil, errors.New("Unknown attribute " + propertyKey)
	}
	return newProperty(node, propertyId, introspector)
}

func (this *Property) Parent() *Property {
	return this.parent
}

func (this *Property) Node() *types.RNode {
	return this.node
}

func (this *Property) Key() interface{} {
	return this.key
}

func (this *Property) Value() interface{} {
	return this.value
}

func (this *Property) setKeyValue(propertyId string) (string, error) {
	id := propertyId
	dIndex := strings.LastIndex(propertyId, ".")
	if dIndex == -1 {
		return "", nil
	}
	beIndex := strings.LastIndex(propertyId, ">")
	if beIndex == -1 {
		return "", nil
	}
	for dIndex < beIndex {
		id = id[0:beIndex]
		dIndex = strings.LastIndex(id, ".")
		beIndex = strings.LastIndex(id, ">")
	}
	prefix := propertyId[0:dIndex]
	suffix := propertyId[dIndex+1:]
	bbIndex := strings.LastIndex(suffix, "<")
	if bbIndex == -1 {
		return prefix, nil
	}

	v := suffix[bbIndex+1 : len(suffix)-1]
	this.key = string_utils.FromString(v, this.introspector.Registry()).Interface()
	return prefix, nil
}

func (this *Property) PropertyId() (string, error) {
	if this.id != "" {
		return this.id, nil
	}
	buff := string_utils.New()
	if this.parent == nil {
		buff.Add(strings.ToLower(this.node.TypeName))
		buff.Add(this.node.CachedKey)
	} else {
		pi, err := this.parent.PropertyId()
		if err != nil {
			return "", err
		}
		buff.Add(pi)
		buff.Add(".")
		buff.Add(strings.ToLower(this.node.FieldName))
	}

	if this.key != nil {
		keyStr := string_utils.New()
		keyStr.TypesPrefix = true
		buff.Add("<")
		buff.Add(keyStr.StringOf(this.key))
		buff.Add(">")
	}
	this.id = buff.String()
	return this.id, nil
}

func newProperty(node *types.RNode, propertyPath string, introspector common.IIntrospect) (*Property, error) {
	property := &Property{}
	property.node = node
	property.introspector = introspector
	if node.Parent != nil {
		prefix, err := property.setKeyValue(propertyPath)
		if err != nil {
			return nil, err
		}
		pi, err := newProperty(node.Parent, prefix, introspector)
		if err != nil {
			return nil, err
		}
		property.parent = pi
	} else {
		index1 := strings.Index(propertyPath, "<")
		index2 := strings.Index(propertyPath, ">")
		if index1 != -1 && index2 != -1 && index2 > index1 {
			property.key = string_utils.FromString(propertyPath[index1+1:index2], property.introspector.Registry()).Interface()
		}
	}
	return property, nil
}
