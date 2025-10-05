package updating

import (
	"reflect"

	"github.com/saichler/l8reflect/go/reflect/introspecting"
	"github.com/saichler/l8reflect/go/reflect/properties"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8reflect"
)

func mapUpdate(instance *properties.Property, node *l8reflect.L8Node, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.IsNil() && newValue.IsNil() {
		return nil
	}
	if oldValue.IsNil() && !newValue.IsNil() {
		updates.addUpdate(instance, nil, newValue.Interface())
		oldValue.Set(newValue)
		return nil
	}
	if !oldValue.IsNil() && newValue.IsNil() && updates.nilIsValid {
		updates.addUpdate(instance, oldValue.Interface(), nil)
		oldValue.Set(newValue)
		return nil
	}

	if newValue.IsValid() && !newValue.IsNil() && introspecting.AlwaysFullDecorator(node) {
		updates.addUpdate(instance, nil, newValue.Interface())
		oldValue.Set(newValue)
		return nil
	}

	newKeys := newValue.MapKeys()
	for _, key := range newKeys {
		oldKeyValue := oldValue.MapIndex(key)
		newKeyValue := newValue.MapIndex(key)

		if !oldKeyValue.IsValid() {
			subProperty := properties.NewProperty(node, instance.Parent().(*properties.Property), key.Interface(),
				newKeyValue.Interface(), updates.resources)
			updates.addUpdate(subProperty, nil, newKeyValue.Interface())
			oldValue.SetMapIndex(key, newKeyValue)
			continue
		}

		if !node.IsStruct {
			if deepEqual.Equal(oldKeyValue.Interface(), newKeyValue.Interface()) {
				continue
			}
			subProperty := properties.NewProperty(node, instance.Parent().(*properties.Property), key.Interface(), newKeyValue.Interface(), updates.resources)
			updates.addUpdate(subProperty, nil, newKeyValue.Interface())
			oldValue.SetMapIndex(key, newKeyValue)
		} else if oldKeyValue.IsValid() && newKeyValue.IsValid() {
			if deepEqual.Equal(oldKeyValue.Interface(), newKeyValue.Interface()) {
				continue
			}
			subProperty := properties.NewProperty(node, instance.Parent().(*properties.Property), key.Interface(), newKeyValue.Interface(), updates.resources)
			err := structUpdate(subProperty, node, oldKeyValue.Elem(), newKeyValue.Elem(), updates)
			if err != nil {
				return err
			}
		}
	}

	if updates.newItemIsFull {
		oldKeys := oldValue.MapKeys()
		for _, key := range oldKeys {
			newKeyValue := newValue.MapIndex(key)
			oldKeyValue := oldValue.MapIndex(key)
			if !newKeyValue.IsValid() {
				subProperty := properties.NewProperty(node, instance.Parent().(*properties.Property), key.Interface(), nil, updates.resources)
				updates.addUpdate(subProperty, oldKeyValue.Interface(), ifs.Deleted_Entry)
				oldValue.SetMapIndex(key, reflect.Value{})
			}
		}
	}
	return nil
}
