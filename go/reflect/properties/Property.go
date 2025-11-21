package properties

import (
	"errors"
	"reflect"
	"strings"

	"github.com/saichler/l8reflect/go/reflect/helping"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8reflect"
	strings2 "github.com/saichler/l8utils/go/utils/strings"
)

type Property struct {
	parent    *Property
	node      *l8reflect.L8Node
	key       interface{}
	value     interface{}
	id        string
	displayId string
	isLeaf    bool
	resources ifs.IResources
}

func NewProperty(node *l8reflect.L8Node, parent *Property, key interface{}, value interface{}, resources ifs.IResources) *Property {
	property := &Property{}
	property.parent = parent
	property.node = node
	property.key = key
	property.value = value
	property.resources = resources
	property.isLeaf = true
	if parent != nil {
		parent.isLeaf = false
	}
	return property
}

func PropertyOf(propertyId string, resources ifs.IResources) (*Property, error) {
	propertyKey := helping.PropertyNodeKey(propertyId)
	node, ok := resources.Introspector().Node(propertyKey)
	if !ok {
		return nil, errors.New("Unknown attribute " + propertyKey)
	}
	return newProperty(node, propertyId, resources)
}

func (this *Property) Parent() ifs.IProperty {
	return this.parent
}

func (this *Property) Node() *l8reflect.L8Node {
	return this.node
}

func (this *Property) Key() interface{} {
	return this.key
}

func (this *Property) Value() interface{} {
	return this.value
}

func (this *Property) Resources() ifs.IResources {
	return this.resources
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

	if dIndex > beIndex {
		prefix := propertyId[0:dIndex]
		return prefix, nil
	}

	bsIndex := strings.LastIndex(propertyId, "<")
	if dIndex > bsIndex {
		id = propertyId[:bsIndex]
		dIndex = strings.LastIndex(id, ".")
	}
	prefix := propertyId[0:dIndex]
	suffix := propertyId[dIndex+1:]
	bbIndex := strings.LastIndex(suffix, "<")
	if bbIndex == -1 {
		return prefix, nil
	}

	v := suffix[bbIndex+1 : len(suffix)-1]
	k, e := strings2.FromString(v, this.resources.Registry())
	if e != nil {
		return "", e
	}
	this.key = k.Interface()
	return prefix, nil
}

func (this *Property) IsString() bool {
	if this.node.TypeName == reflect.String.String() {
		return true
	}
	return false
}

func (this *Property) PropertyId() (string, error) {
	if this.id != "" {
		return this.id, nil
	}
	buff := strings2.New()
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
		keyStr := strings2.New()
		keyStr.TypesPrefix = true
		buff.Add("<")
		buff.Add(keyStr.StringOf(this.key))
		buff.Add(">")
	}
	this.id = buff.String()
	return this.id, nil
}

func (this *Property) PropertyDisplayId() string {
	if this.displayId != "" {
		return this.displayId
	}
	buff := strings2.New()
	if this.parent == nil {
		buff.Add("[")
		buff.Add(this.node.TypeName)
		buff.Add(this.node.CachedKey)
		buff.Add("]")
		return buff.String()
	} else {
		pi := this.parent.PropertyDisplayId()
		buff.Add(pi)
		buff.Add("[")
		buff.Add(this.node.FieldName)
	}
	if this.key != nil {
		keyStr := strings2.New()
		buff.Add(" ")
		buff.Add(keyStr.StringOf(this.key))
	}
	buff.Add("]")
	this.displayId = buff.String()
	return this.displayId
}

func (this *Property) IsLeaf() bool {
	return this.isLeaf
}

func newProperty(node *l8reflect.L8Node, propertyPath string, resources ifs.IResources) (*Property, error) {
	property := &Property{}
	property.isLeaf = true
	property.node = node
	property.resources = resources
	if node.Parent != nil {
		prefix, err := property.setKeyValue(propertyPath)
		if err != nil {
			return nil, err
		}
		pi, err := newProperty(node.Parent, prefix, resources)
		if err != nil {
			return nil, err
		}
		property.parent = pi
		pi.isLeaf = false
	} else {
		index1 := strings.Index(propertyPath, "<")
		index2 := strings.Index(propertyPath, ">")
		if index1 != -1 && index2 != -1 && index2 > index1 {
			k, e := strings2.FromString(propertyPath[index1+1:index2], property.resources.Registry())
			if e != nil {
				return nil, e
			}
			property.key = k.Interface()
		}
	}
	return property, nil
}
