package updating

import (
	"errors"
	"github.com/saichler/reflect/go/reflect/cloning"
	"github.com/saichler/reflect/go/reflect/introspecting"
	"github.com/saichler/reflect/go/reflect/properties"
	"github.com/saichler/types/go/types"
	"reflect"
)

var comparators map[reflect.Kind]func(*properties.Property, *types.RNode, reflect.Value, reflect.Value, *Updater) error
var deepEqual = cloning.NewDeepEqual()

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
	if oldValue.Len() != newValue.Len() {
		oldValue.Set(newValue)
		updates.addUpdate(instance, oldValue, newValue)
		return nil
	}
	for i := 0; i < oldValue.Len(); i++ {
		oldIndexValue := oldValue.Index(i)
		newIndexValue := newValue.Index(i)
		if oldIndexValue.IsValid() && newIndexValue.IsValid() {
			if deepEqual.Equal(oldIndexValue.Interface(), newIndexValue.Interface()) {
				continue
			}
			subProperty := properties.NewProperty(node, instance.Parent().(*properties.Property), i, newIndexValue.Interface(), updates.introspector)
			err := structUpdate(subProperty, node, oldIndexValue.Elem(), newIndexValue.Elem(), updates)
			return err
		}
	}
	return nil
}

func deepMapUpdate(instance *properties.Property, node *types.RNode, oldValue, newValue reflect.Value, updates *Updater) error {
	newKeys := newValue.MapKeys()
	for _, key := range newKeys {
		oldKeyValue := oldValue.MapIndex(key)
		newKeyValue := newValue.MapIndex(key)
		if !oldKeyValue.IsValid() ||
			(oldKeyValue.Kind() == reflect.Ptr && oldKeyValue.IsNil()) {
			subProperty := properties.NewProperty(node, instance.Parent().(*properties.Property), key.Interface(), newKeyValue.Interface(), updates.introspector)
			updates.addUpdate(subProperty, nil, newKeyValue.Interface())
			oldValue.SetMapIndex(key, newKeyValue)
		} else if oldKeyValue.IsValid() && newKeyValue.IsValid() {
			if deepEqual.Equal(oldKeyValue.Interface(), newKeyValue.Interface()) {
				continue
			}
			subProperty := properties.NewProperty(node, instance.Parent().(*properties.Property), key.Interface(), newKeyValue.Interface(), updates.introspector)
			err := structUpdate(subProperty, node, oldKeyValue.Elem(), newKeyValue.Elem(), updates)
			return err
		}
	}

	oldKeys := oldValue.MapKeys()
	for _, key := range oldKeys {
		newKeyValue := newValue.MapIndex(key)
		oldKeyValue := oldValue.MapIndex(key)
		if !newKeyValue.IsValid() {
			subProperty := properties.NewProperty(node, instance.Parent().(*properties.Property), key.Interface(), nil, updates.introspector)
			updates.forceUpdate(subProperty, oldKeyValue.Interface(), nil)
			oldValue.SetMapIndex(key, reflect.Value{})
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

	if node.IsStruct {
		//By default we do nested inspection of map & slices when the node is struct
		//it is possible to disable that with a decorator so map & slice
		//will be copied fully if not eq
		noDeepInspection := introspecting.NoNestedInspection(node)
		if !noDeepInspection {
			if node.IsSlice {
				return deepSliceUpdate(instance, node, oldValue, newValue, updates)
			} else if node.IsMap {
				return deepMapUpdate(instance, node, oldValue, newValue, updates)
			}
		}
	}

	oldIns := oldValue.Interface()
	newIns := newValue.Interface()

	if !deepEqual.Equal(oldIns, newIns) {
		updates.addUpdate(instance, oldIns, newIns)
		oldValue.Set(newValue)
	}

	return nil
}
