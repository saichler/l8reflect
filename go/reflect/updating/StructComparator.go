package updating

import (
	"errors"
	"github.com/saichler/l8types/go/types"
	"github.com/saichler/reflect/go/reflect/properties"
	"reflect"
)

func ptrUpdate(property *properties.Property, node *types.RNode, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.IsNil() && !newValue.IsNil() {
		updates.addUpdate(property, nil, newValue.Interface())
		oldValue.Set(newValue)
		return nil
	}
	if !oldValue.IsNil() && newValue.IsNil() && updates.isNilValid {
		updates.addUpdate(property, oldValue, nil)
		oldValue.Set(newValue)
		return nil
	}
	if oldValue.IsNil() && newValue.IsNil() {
		return nil
	}
	return update(property, node, oldValue.Elem(), newValue.Elem(), updates)
}

func structUpdate(property *properties.Property, node *types.RNode, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.Type().Name() != newValue.Type().Name() {
		return errors.New("Mismatch type, old=" + oldValue.Type().Name() + ", new=" + newValue.Type().Name())
	}
	for _, attr := range node.Attributes {
		oldFldValue := oldValue.FieldByName(attr.FieldName)
		newFldValue := newValue.FieldByName(attr.FieldName)
		subInstance := properties.NewProperty(attr, property, nil, oldFldValue, updates.resources)
		err := update(subInstance, attr, oldFldValue, newFldValue, updates)
		if err != nil {
			return err
		}
	}
	return nil
}
