package updater

import (
	"errors"
	"github.com/saichler/reflect/go/reflect/common"
	"github.com/saichler/reflect/go/reflect/property"
	"github.com/saichler/reflect/go/types"
	"reflect"
)

type Updater struct {
	nilIsValid   bool
	property     *property.Property
	changes      []*Change
	introspector common.IIntrospect
}

func NewUpdater(introspector common.IIntrospect, nilIsValid bool) *Updater {
	updates := &Updater{}
	updates.changes = make([]*Change, 0)
	updates.introspector = introspector
	updates.nilIsValid = nilIsValid
	return updates
}

func (this *Updater) Changes() []*Change {
	return this.changes
}

func (this *Updater) Update(old, new interface{}, introspect common.IIntrospect) error {
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

	prop := property.NewProperty(node, nil, common.PrimaryDecorator(node, oldValue, this.introspector.Registry()), oldValue, this.introspector)

	err := update(prop, node, oldValue, newValue, this)
	return err
}

func update(instance *property.Property, node *types.RNode, oldValue, newValue reflect.Value, updates *Updater) error {
	if !newValue.IsValid() {
		return nil
	}
	if newValue.Kind() == reflect.Ptr && newValue.IsNil() && !updates.nilIsValid {
		return nil
	}

	kind := oldValue.Kind()
	comparator := comparators[kind]
	if comparator == nil {
		panic("No comparator for kind:" + kind.String() + ", please add it!")
	}
	return comparator(instance, node, oldValue, newValue, updates)
}

func (this *Updater) addUpdate(instance *property.Property, node *types.RNode, oldValue, newValue interface{}) {
	if !this.nilIsValid && newValue == nil {
		return
	}
	this.changes = append(this.changes, NewChange(oldValue, newValue, instance))
}
