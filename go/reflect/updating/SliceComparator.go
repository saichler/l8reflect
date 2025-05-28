package updating

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types"
	"github.com/saichler/reflect/go/reflect/properties"
	"reflect"
)

func sliceUpdate(instance *properties.Property, node *types.RNode, oldValue, newValue reflect.Value, updates *Updater) error {
	if oldValue.IsNil() && newValue.IsNil() {
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

	size := newValue.Len()
	if size > oldValue.Len() {
		size = oldValue.Len()
	}

	for i := 0; i < size; i++ {
		oldIndexValue := oldValue.Index(i)
		newIndexValue := newValue.Index(i)
		if !node.IsStruct {
			if oldIndexValue.IsValid() && deepEqual.Equal(oldIndexValue.Interface(), newIndexValue.Interface()) {
				continue
			}
			subProperty := properties.NewProperty(node, instance.Parent().(*properties.Property), i,
				newIndexValue.Interface(), updates.resources)
			updates.addUpdate(subProperty, nil, newIndexValue.Interface())
			oldIndexValue.Set(newIndexValue)
		} else if !oldIndexValue.IsValid() || oldIndexValue.IsNil() {
			subProperty := properties.NewProperty(node, instance.Parent().(*properties.Property),
				i, newIndexValue.Interface(), updates.resources)
			updates.addUpdate(subProperty, nil, newIndexValue.Interface())
			oldIndexValue.Set(newIndexValue)
		} else if oldIndexValue.IsValid() && newIndexValue.IsValid() {
			if deepEqual.Equal(oldIndexValue.Interface(), newIndexValue.Interface()) {
				continue
			}
			subProperty := properties.NewProperty(node, instance.Parent().(*properties.Property),
				i, newIndexValue.Interface(), updates.resources)
			err := structUpdate(subProperty, node, oldIndexValue.Elem(), newIndexValue.Elem(), updates)
			if err != nil {
				return err
			}
		}
	}

	vInfo, err := instance.Resources().Registry().Info(instance.Node().TypeName)
	if err != nil {
		return err
	}

	if size < oldValue.Len() {
		var newSlice reflect.Value
		if node.IsStruct {
			newSlice = reflect.MakeSlice(reflect.SliceOf(reflect.PointerTo(vInfo.Type())), size, size)
		} else {
			newSlice = reflect.MakeSlice(reflect.SliceOf(vInfo.Type()), size, size)
		}

		for i := 0; i < size; i++ {
			newSlice.Index(i).Set(oldValue.Index(i))
		}
		subProperty := properties.NewProperty(node, instance.Parent().(*properties.Property), size,
			nil, updates.resources)
		updates.addUpdate(subProperty, nil, ifs.Deleted_Entry)
		oldValue.Set(newSlice)
	} else if size > oldValue.Len() {
		newSlice := reflect.MakeSlice(reflect.SliceOf(reflect.PointerTo(vInfo.Type())), size, size)
		for i := 0; i < size; i++ {
			newSlice.Index(i).Set(oldValue.Index(i))
		}
		for i := size; i < newValue.Len(); i++ {
			newV := newValue.Index(i)
			newSlice.Index(i).Set(newV)
			subProperty := properties.NewProperty(node, instance.Parent().(*properties.Property), i,
				newV.Interface(), updates.resources)
			updates.addUpdate(subProperty, nil, newV.Interface())
		}
		oldValue.Set(newSlice)
	}

	return nil
}
