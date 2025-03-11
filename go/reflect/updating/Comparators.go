package updating

import (
	"errors"
	"github.com/saichler/reflect/go/reflect/helping"
	"github.com/saichler/reflect/go/reflect/properties"
	"github.com/saichler/types/go/types"
	"reflect"
)

var comparators map[reflect.Kind]func(*properties.Property, *types.RNode, reflect.Value, reflect.Value, *Updater) error

func init() {
	comparators = make(map[reflect.Kind]func(*properties.Property, *types.RNode, reflect.Value, reflect.Value, *Updater) error)
	comparators[reflect.Int] = intUpdate
	comparators[reflect.Int32] = intUpdate
	comparators[reflect.Int64] = intUpdate

	comparators[reflect.Uint] = uintUpdate
	comparators[reflect.Uint32] = uintUpdate
	comparators[reflect.Uint64] = uintUpdate

	comparators[reflect.String] = stringUpdate

	comparators[reflect.Bool] = boolUpdate

	comparators[reflect.Float32] = floatUpdate
	comparators[reflect.Float64] = floatUpdate

	comparators[reflect.Ptr] = ptrUpdate

	comparators[reflect.Struct] = structUpdate

	comparators[reflect.Slice] = sliceOrMapUpdate

	comparators[reflect.Map] = sliceOrMapUpdate
}

func intUpdate(property *properties.Property, node *types.RNode, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.Int() != newValue.Int() && (newValue.Int() != 0 || updates.isNilValid) {
		updates.addUpdate(property, oldValue.Interface(), newValue.Interface())
		oldValue.Set(newValue)
	}
	return nil
}

func uintUpdate(instance *properties.Property, node *types.RNode, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.Uint() != newValue.Uint() && (newValue.Uint() != 0 || updates.isNilValid) {
		updates.addUpdate(instance, oldValue.Interface(), newValue.Interface())
		oldValue.Set(newValue)
	}
	return nil
}

func stringUpdate(instance *properties.Property, node *types.RNode, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.String() != newValue.String() && (newValue.String() != "" || updates.isNilValid) {
		updates.addUpdate(instance, oldValue.Interface(), newValue.Interface())
		oldValue.Set(newValue)
	}
	return nil
}

func boolUpdate(instance *properties.Property, node *types.RNode, oldValue, newValue reflect.Value, updates *Updater) error {
	if newValue.Bool() && !oldValue.Bool() || updates.isNilValid {
		updates.addUpdate(instance, oldValue.Interface(), newValue.Interface())
		oldValue.Set(newValue)
	}
	return nil
}

func floatUpdate(instance *properties.Property, node *types.RNode, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.Float() != newValue.Float() && (newValue.Float() != 0 || updates.isNilValid) {
		updates.addUpdate(instance, oldValue.Interface(), newValue.Interface())
		oldValue.Set(newValue)
	}
	return nil
}

func ptrUpdate(instance *properties.Property, node *types.RNode, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.IsNil() && !newValue.IsNil() {
		updates.addUpdate(instance, nil, newValue.Interface())
		oldValue.Set(newValue)
		return nil
	}
	if !oldValue.IsNil() && newValue.IsNil() && updates.isNilValid {
		updates.addUpdate(instance, oldValue, nil)
		oldValue.Set(newValue)
		return nil
	}
	if oldValue.IsNil() && newValue.IsNil() {
		return nil
	}
	return update(instance, node, oldValue.Elem(), newValue.Elem(), updates)
}

func structUpdate(prop *properties.Property, node *types.RNode, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.Type().Name() != newValue.Type().Name() {
		return errors.New("Mismatch type, old=" + oldValue.Type().Name() + ", new=" + newValue.Type().Name())
	}
	for _, attr := range node.Attributes {
		oldFldValue := oldValue.FieldByName(attr.FieldName)
		newFldValue := newValue.FieldByName(attr.FieldName)
		subInstance := properties.NewProperty(attr, prop, nil, oldFldValue, updates.introspector)
		err := update(subInstance, attr, oldFldValue, newFldValue, updates)
		if err != nil {
			return err
		}
	}
	return nil
}

func deepSliceUpdate(instance *properties.Property, node *types.RNode, oldValue, newValue reflect.Value, updates *Updater) error {
	//TODO - implement deep slice update
	return nil
}

func deepMapUpdate(instance *properties.Property, node *types.RNode, oldValue, newValue reflect.Value, updates *Updater) error {
	//TODO - implement deep map update
	return nil
}

func sliceOrMapUpdate(instance *properties.Property, node *types.RNode, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.IsNil() && !newValue.IsNil() {
		updates.addUpdate(instance, nil, newValue.Interface())
		oldValue.Set(newValue)
		return nil
	}
	if oldValue.IsNil() && !newValue.IsNil() {
		updates.addUpdate(instance, nil, newValue.Interface())
		oldValue.Set(newValue)
		return nil
	}
	if !oldValue.IsNil() && newValue.IsNil() && updates.isNilValid {
		updates.addUpdate(instance, oldValue, nil)
		oldValue.Set(newValue)
		return nil
	}
	if oldValue.IsNil() && newValue.IsNil() {
		return nil
	}

	//If this is a struct, we need to check if we need to do deep update
	//and not just copy the new slice/map to the old slice/map
	if updates.introspector.Kind(node) == reflect.Struct {
		if helping.DeepDecorator(node) {
			if node.IsSlice {
				return deepSliceUpdate(instance, node, oldValue, newValue, updates)
			} else if node.IsMap {
				return deepMapUpdate(instance, node, oldValue, newValue, updates)
			}
		}
	}

	eq := reflect.DeepEqual(oldValue.Interface(), newValue.Interface())
	if !eq {
		updates.addUpdate(instance, oldValue, nil)
		oldValue.Set(newValue)
	}

	return nil
}
