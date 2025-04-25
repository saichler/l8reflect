package updating

import (
	"github.com/saichler/reflect/go/reflect/cloning"
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

	comparators[reflect.Slice] = sliceUpdate

	comparators[reflect.Map] = mapUpdate
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
