package updating

import (
	"errors"
	"github.com/saichler/reflect/go/reflect/helping"
	"github.com/saichler/reflect/go/reflect/properties"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/types"
	"reflect"
)

type Updater struct {
	changes      []*Change
	introspector common.IIntrospector
	isNilValid   bool
}

func NewUpdater(introspector common.IIntrospector, isNilValid bool) *Updater {
	upd := &Updater{}
	upd.introspector = introspector
	upd.isNilValid = isNilValid
	return upd
}

func (this *Updater) Changes() []*Change {
	return this.changes
}

func (this *Updater) Update(old, new interface{}) error {
	oldValue := reflect.ValueOf(old)
	newValue := reflect.ValueOf(new)
	if !oldValue.IsValid() || !newValue.IsValid() {
		return errors.New("either old or new are nil or invalid")
	}
	if oldValue.Kind() == reflect.Ptr {
		oldValue = oldValue.Elem()
		newValue = newValue.Elem()
	}
	node, _ := this.introspector.Node(oldValue.Type().Name())
	if node == nil {
		return errors.New("cannot find node for type " + oldValue.Type().Name() + ", please register it")
	}

	pKey := helping.PrimaryDecorator(node, oldValue, this.introspector.Registry())
	prop := properties.NewProperty(node, nil, pKey, oldValue, this.introspector)
	err := update(prop, node, oldValue, newValue, this)
	return err
}

func update(instance *properties.Property, node *types.RNode, oldValue, newValue reflect.Value, updates *Updater) error {
	if !newValue.IsValid() {
		return nil
	}
	if newValue.Kind() == reflect.Ptr && newValue.IsNil() && !updates.isNilValid {
		return nil
	}

	kind := oldValue.Kind()
	comparator := comparators[kind]
	if comparator == nil {
		panic("No comparator for kind:" + kind.String() + ", please add it!")
	}
	return comparator(instance, node, oldValue, newValue, updates)
}

func (this *Updater) addUpdate(prop *properties.Property, oldValue, newValue interface{}) {
	if !this.isNilValid && newValue == nil {
		return
	}
	if this.changes == nil {
		this.changes = make([]*Change, 0)
	}
	this.changes = append(this.changes, NewChange(oldValue, newValue, prop))
}

func (this *Updater) forceUpdate(prop *properties.Property, oldValue, newValue interface{}) {
	if newValue == nil {
		newValue = "nil"
	}
	if this.changes == nil {
		this.changes = make([]*Change, 0)
	}
	this.changes = append(this.changes, NewChange(oldValue, newValue, prop))
}
