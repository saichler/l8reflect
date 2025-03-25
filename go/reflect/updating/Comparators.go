package updating

import (
	"errors"
	"github.com/saichler/reflect/go/reflect/introspecting"
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
	if newValue.Bool() == oldValue.Bool() {
		return nil
	}
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
	panic("Implement me")
}

func deepMapUpdate(instance *properties.Property, node *types.RNode, oldValue, newValue reflect.Value, updates *Updater) error {
	newKeys := newValue.MapKeys()
	for _, key := range newKeys {
		oldKeyValue := oldValue.MapIndex(key)
		newKeyValue := newValue.MapIndex(key)
		if !oldKeyValue.IsValid() ||
			(oldKeyValue.Kind() == reflect.Ptr && oldKeyValue.IsNil()) {
			subProperty := properties.NewProperty(node, instance, key.Interface(), newKeyValue.Interface(), updates.introspector)
			updates.addUpdate(subProperty, nil, newKeyValue.Interface())
			oldValue.SetMapIndex(key, newKeyValue)
		} else if oldKeyValue.IsValid() && newKeyValue.IsValid() {
			subProperty := properties.NewProperty(node, instance.Parent(), key.Interface(), newKeyValue.Interface(), updates.introspector)
			structUpdate(subProperty, node, oldKeyValue.Elem(), newKeyValue.Elem(), updates)
		}
	}
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
	if newValue.IsNil() && !updates.isNilValid {
		return nil
	}

	//If this is a struct, we need to check if we need to do deep update
	//and not just copy the new slice/map to the old slice/map
	if node.IsStruct {
		if introspecting.DeepDecorator(node) {
			if node.IsSlice {
				return deepSliceUpdate(instance, node, oldValue, newValue, updates)
			} else if node.IsMap {
				return deepMapUpdate(instance, node, oldValue, newValue, updates)
			}
		}
	}

	oldIns := oldValue.Interface()
	newIns := newValue.Interface()

	eq := reflect.DeepEqual(oldIns, newIns)
	if !eq {
		updates.addUpdate(instance, oldIns, newIns)
		oldValue.Set(newValue)
	}

	return nil
}
