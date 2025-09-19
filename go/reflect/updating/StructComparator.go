package updating

import (
	"errors"
	"reflect"

	"github.com/saichler/l8types/go/types/l8reflect"
	"github.com/saichler/l8reflect/go/reflect/properties"
)

func ptrUpdate(property *properties.Property, node *l8reflect.L8Node, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.IsNil() && !newValue.IsNil() {
		updates.addUpdate(property, nil, newValue.Interface())
		oldValue.Set(newValue)
		return nil
	}
	if !oldValue.IsNil() && newValue.IsNil() && updates.nilIsValid {
		updates.addUpdate(property, oldValue, nil)
		oldValue.Set(newValue)
		return nil
	}
	if oldValue.IsNil() && newValue.IsNil() {
		return nil
	}
	return update(property, node, oldValue.Elem(), newValue.Elem(), updates)
}

func structUpdate(property *properties.Property, node *l8reflect.L8Node, oldValue, newValue reflect.Value, updates *Updater) error {
	if !oldValue.IsValid() && newValue.IsValid() {
		oldValue.Set(newValue)
		updates.addUpdate(property, nil, newValue.Interface())
		return nil
	}
	if oldValue.IsValid() && !newValue.IsValid() && updates.nilIsValid {
		newValue.Set(reflect.New(oldValue.Type()).Elem())
		updates.addUpdate(property, oldValue.Interface(), newValue.Interface())
		return nil
	}
	if !oldValue.IsValid() && !newValue.IsValid() {
		return nil
	}

	if !newValue.IsValid() && oldValue.IsValid() {
		newValue = reflect.New(oldValue.Type()).Elem()
		newValue.Set(reflect.New(oldValue.Type()).Elem())
	}

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
