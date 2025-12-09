package introspecting

import (
	"errors"
	"reflect"

	"github.com/saichler/l8reflect/go/reflect/helping"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8reflect"
	strings2 "github.com/saichler/l8utils/go/utils/strings"
)

func (this *Introspector) Decorators() ifs.IDecorators {
	return this
}

func (this *Introspector) AddPrimaryKeyDecorator(any interface{}, fields ...string) error {
	node, _, err := this.NodeFor(any)
	if err != nil || node == nil {
		node, _ = this.Inspect(any)
	}
	addDecorator(l8reflect.L8DecoratorType_Primary, fields, node)
	return nil
}

func (this *Introspector) AddUniqueKeyDecorator(any interface{}, fields ...string) error {
	node, _, err := this.NodeFor(any)
	if err != nil || node == nil {
		return err
	}
	addDecorator(l8reflect.L8DecoratorType_Unique, fields, node)
	return nil
}

func (this *Introspector) AddAlwayOverwriteDecorator(nodeId string) error {
	node, ok := this.Node(nodeId)
	if !ok {
		return errors.New(strings2.New("Node for ID ", nodeId, " not found").String())
	}
	addAlwayOverwriteDecorator(node)
	return nil
}

func (this *Introspector) AddNoNestedInspection(any interface{}) error {
	node, _, err := this.NodeFor(any)
	if err != nil {
		return err
	}
	addNoNestedInspection(node)
	return nil
}

func (this *Introspector) NodeFor(any interface{}) (*l8reflect.L8Node, reflect.Value, error) {
	if any == nil {
		panic("Node For a nil interface")
	}
	v, e := helping.PtrValue(any)
	if e != nil {
		return nil, v, e
	}
	node, ok := this.Node(v.Type().Name())
	if !ok {
		node, e = this.Inspect(any)
		if e != nil {
			return nil, v, e
		}
	}
	return node, v, nil
}

func (this *Introspector) PrimaryKeyDecoratorValue(any interface{}) (string, *l8reflect.L8Node, error) {
	node, v, err := this.NodeFor(any)
	if err != nil {
		return "", node, err
	}
	return this.PrimaryKeyDecoratorFromValue(node, v)
}

func (this *Introspector) UniqueKeyDecoratorValue(any interface{}) (string, *l8reflect.L8Node, error) {
	node, v, err := this.NodeFor(any)
	if err != nil {
		return "", node, err
	}
	return this.uniqueKeyDecoratorValue(node, v)
}

func (this *Introspector) uniqueKeyDecoratorValue(node *l8reflect.L8Node, value reflect.Value) (string, *l8reflect.L8Node, error) {
	return this.decoratorKey(node, l8reflect.L8DecoratorType_Unique, value)
}

func (this *Introspector) PrimaryKeyDecoratorFromValue(node *l8reflect.L8Node, value reflect.Value) (string, *l8reflect.L8Node, error) {
	return this.decoratorKey(node, l8reflect.L8DecoratorType_Primary, value)
}

func (this *Introspector) decoratorKey(node *l8reflect.L8Node, decoratorType l8reflect.L8DecoratorType, value reflect.Value) (string, *l8reflect.L8Node, error) {
	fields, err := this.Fields(node, decoratorType)
	if err != nil {
		return "", node, err
	}

	str := strings2.New()
	str.TypesPrefix = true
	first := true
	for _, field := range fields {
		if !first {
			str.TypesPrefix = false
			str.Add("::")
			str.TypesPrefix = true
		}
		v := value.FieldByName(field).Interface()
		v2 := str.StringOf(v)
		str.Add(v2)
		first = false
	}
	return str.String(), node, nil
}

func (this *Introspector) Fields(node *l8reflect.L8Node, decoratorType l8reflect.L8DecoratorType) ([]string, error) {
	if node == nil {
		return nil, errors.New("Node is nil")
	}
	decValue := node.Decorators[int32(decoratorType)]
	if decValue == nil {
		return nil, errors.New(strings2.New("Decorator Not Found in ", node.TypeName).String())
	}
	return decValue.Fields, nil
}

func (this *Introspector) KeyForValue(fields []string, value reflect.Value, typeName string, returnError bool) (string, error) {
	if fields == nil || len(fields) == 0 {
		if returnError {
			return "", errors.New(strings2.New("Primary Key Decorator is empty for type ", typeName).String())
		}
		return "", nil
	}
	switch len(fields) {
	case 1:
		return strings2.New(value.FieldByName(fields[0]).Interface()).String(), nil
	case 2:
		return strings2.New(value.FieldByName(fields[0]).Interface(), value.FieldByName(fields[1]).Interface()).String(), nil
	case 3:
		return strings2.New(value.FieldByName(fields[0]).Interface(),
			value.FieldByName(fields[1]).Interface(),
			value.FieldByName(fields[2]).Interface()).String(), nil
	default:
		result := strings2.New()
		for i := 0; i < len(fields); i++ {
			result.Add(result.StringOf(value.FieldByName(fields[i]).Interface()))
		}
		return result.String(), nil
	}
	return "", errors.New("Unexpected code")
}

func addDecorator(decoratorType l8reflect.L8DecoratorType, fields []string, node *l8reflect.L8Node) {
	if node.Decorators == nil {
		node.Decorators = make(map[int32]*l8reflect.L8Decorator)
	}
	node.Decorators[int32(decoratorType)] = &l8reflect.L8Decorator{Fields: fields}
}

func addNoNestedInspection(rnode *l8reflect.L8Node) {
	addDecorator(l8reflect.L8DecoratorType_NoNestedInspection, []string{}, rnode)
}

func (this *Introspector) NoNestedInspection(any interface{}) bool {
	return this.BoolDecoratorValueFor(any, l8reflect.L8DecoratorType_NoNestedInspection)
}

func addAlwayOverwriteDecorator(rnode *l8reflect.L8Node) {
	addDecorator(l8reflect.L8DecoratorType_AlwaysFull, []string{}, rnode)
}

func (this *Introspector) AlwaysFullDecorator(any interface{}) bool {
	return this.BoolDecoratorValueFor(any, l8reflect.L8DecoratorType_AlwaysFull)
}

func (this *Introspector) BoolDecoratorValueFor(any interface{}, typ l8reflect.L8DecoratorType) bool {
	node, _, err := this.NodeFor(any)
	if err != nil {
		return false
	}
	return this.BoolDecoratorValueForNode(node, typ)
}

func (this *Introspector) BoolDecoratorValueForNode(node *l8reflect.L8Node, typ l8reflect.L8DecoratorType) bool {
	if node == nil {
		return false
	}
	_, err := this.Fields(node, typ)
	if err != nil {
		return false
	}
	return true
}
