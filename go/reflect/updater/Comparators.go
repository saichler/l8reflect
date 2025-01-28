package updater

import (
	"errors"
	"github.com/saichler/my.simple/go/common"
	"github.com/saichler/my.simple/go/instance"
	"github.com/saichler/my.simple/go/introspect/model"
	"reflect"
)

var comparators map[reflect.Kind]func(*instance.Instance, *model.Node, reflect.Value, reflect.Value, *Updater) error

func init() {
	comparators = make(map[reflect.Kind]func(*instance.Instance, *model.Node, reflect.Value, reflect.Value, *Updater) error)
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

func intUpdate(instance *instance.Instance, node *model.Node, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.Int() != newValue.Int() && (newValue.Int() != 0 || updates.nilIsValid) {
		updates.addUpdate(instance, node, oldValue.Interface(), newValue.Interface())
		oldValue.Set(newValue)
	}
	return nil
}

func uintUpdate(instance *instance.Instance, node *model.Node, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.Uint() != newValue.Uint() && (newValue.Uint() != 0 || updates.nilIsValid) {
		updates.addUpdate(instance, node, oldValue.Interface(), newValue.Interface())
		oldValue.Set(newValue)
	}
	return nil
}

func stringUpdate(instance *instance.Instance, node *model.Node, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.String() != newValue.String() && (newValue.String() != "" || updates.nilIsValid) {
		updates.addUpdate(instance, node, oldValue.Interface(), newValue.Interface())
		oldValue.Set(newValue)
	}
	return nil
}

func boolUpdate(instance *instance.Instance, node *model.Node, oldValue, newValue reflect.Value, updates *Updater) error {
	if newValue.Bool() && !oldValue.Bool() || updates.nilIsValid {
		updates.addUpdate(instance, node, oldValue.Interface(), newValue.Interface())
		oldValue.Set(newValue)
	}
	return nil
}

func floatUpdate(instance *instance.Instance, node *model.Node, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.Float() != newValue.Float() && (newValue.Float() != 0 || updates.nilIsValid) {
		updates.addUpdate(instance, node, oldValue.Interface(), newValue.Interface())
		oldValue.Set(newValue)
	}
	return nil
}

func ptrUpdate(instance *instance.Instance, node *model.Node, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.IsNil() && !newValue.IsNil() {
		updates.addUpdate(instance, node, nil, newValue.Interface())
		oldValue.Set(newValue)
		return nil
	}
	if !oldValue.IsNil() && newValue.IsNil() && updates.nilIsValid {
		updates.addUpdate(instance, node, oldValue, nil)
		oldValue.Set(newValue)
		return nil
	}
	if oldValue.IsNil() && newValue.IsNil() {
		return nil
	}
	return update(instance, node, oldValue.Elem(), newValue.Elem(), updates)
}

func structUpdate(inst *instance.Instance, node *model.Node, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.Type().Name() != newValue.Type().Name() {
		return errors.New("Mismatch type, old=" + oldValue.Type().Name() + ", new=" + newValue.Type().Name())
	}
	for _, attr := range node.Attributes {
		oldFldValue := oldValue.FieldByName(attr.FieldName)
		newFldValue := newValue.FieldByName(attr.FieldName)
		subInstance := instance.NewInstance(attr, inst, "", oldFldValue, updates.introspect)
		err := update(subInstance, attr, oldFldValue, newFldValue, updates)
		if err != nil {
			return err
		}
	}
	return nil
}

func deepSliceUpdate(instance *instance.Instance, node *model.Node, oldValue, newValue reflect.Value, updates *Updater) error {
	//TODO - implement deep slice update
	return nil
}

func deepMapUpdate(instance *instance.Instance, node *model.Node, oldValue, newValue reflect.Value, updates *Updater) error {
	//TODO - implement deep map update
	return nil
}

func sliceOrMapUpdate(instance *instance.Instance, node *model.Node, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.IsNil() && !newValue.IsNil() {
		updates.addUpdate(instance, node, nil, newValue.Interface())
		oldValue.Set(newValue)
		return nil
	}
	if oldValue.IsNil() && !newValue.IsNil() {
		updates.addUpdate(instance, node, nil, newValue.Interface())
		oldValue.Set(newValue)
		return nil
	}
	if !oldValue.IsNil() && newValue.IsNil() && updates.nilIsValid {
		updates.addUpdate(instance, node, oldValue, nil)
		oldValue.Set(newValue)
		return nil
	}
	if oldValue.IsNil() && newValue.IsNil() {
		return nil
	}

	//If this is a struct, we need to check if we need to do deep update
	//and not just copy the new slice/map to the old slice/map
	if updates.introspect.Kind(node) == reflect.Struct {
		if common.DeepDecorator(node) {
			if node.IsSlice {
				return deepSliceUpdate(instance, node, oldValue, newValue, updates)
			} else if node.IsMap {
				return deepMapUpdate(instance, node, oldValue, newValue, updates)
			}
		}
	}

	eq := reflect.DeepEqual(oldValue.Interface(), newValue.Interface())
	if !eq {
		updates.addUpdate(instance, node, oldValue, nil)
		oldValue.Set(newValue)
	}

	return nil
}
