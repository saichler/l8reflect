package updater

import (
	"errors"
	"github.com/saichler/my.simple/go/common"
	"github.com/saichler/my.simple/go/instance"
	"github.com/saichler/my.simple/go/introspect/model"
	"reflect"
)

type Updater struct {
	nilIsValid bool
	instance   *instance.Instance
	changes    []*Change
	introspect common.IIntrospect
}

func NewUpdater(introspect common.IIntrospect, nilIsValid bool) *Updater {
	updates := &Updater{}
	updates.changes = make([]*Change, 0)
	updates.introspect = introspect
	updates.nilIsValid = nilIsValid
	return updates
}

func (updates *Updater) Changes() []*Change {
	return updates.changes
}

func (updates *Updater) Update(old, new interface{}, introspect common.IIntrospect) error {
	oldValue := reflect.ValueOf(old)
	newValue := reflect.ValueOf(new)
	if !oldValue.IsValid() || !newValue.IsValid() {
		return errors.New("either old or new are nil or invalid")
	}
	if oldValue.Kind() == reflect.Ptr {
		oldValue = oldValue.Elem()
		newValue = newValue.Elem()
	}
	node, _ := introspect.Node(oldValue.Type().Name())
	if node == nil {
		return errors.New("cannot find node for type " + oldValue.Type().Name() + ", please register it")
	}

	instance := instance.NewInstance(node, nil, common.PrimaryDecorator(node, oldValue), oldValue, introspect)

	err := update(instance, node, oldValue, newValue, updates)
	return err
}

func update(instance *instance.Instance, node *model.Node, oldValue, newValue reflect.Value, updates *Updater) error {
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

func (updates *Updater) addUpdate(instance *instance.Instance, node *model.Node, oldValue, newValue interface{}) {
	if !updates.nilIsValid && newValue == nil {
		return
	}
	updates.changes = append(updates.changes, NewChange(oldValue, newValue, instance))
}
