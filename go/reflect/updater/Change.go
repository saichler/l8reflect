package updater

import (
	"github.com/saichler/my.simple/go/utils/logs"
	"github.com/saichler/my.simple/go/utils/strng"
)

type Change struct {
	instance *notifications.Instance
	oldValue interface{}
	newValue interface{}
}

func (change *Change) String() string {
	id, err := change.instance.InstanceId()
	if err != nil {
		logs.Error("Instance String fail: ", err)
	}
	str := strng.New(id)

	str.Add(" - Old=").Add(str.StringOf(change.oldValue)).
		Add(" New=").Add(str.StringOf(change.newValue))
	return str.String()
}

func (change *Change) Apply(any interface{}) {
	change.instance.Set(any, change.newValue)
}

func NewChange(old, new interface{}, instance *notifications.Instance) *Change {
	change := &Change{}
	change.oldValue = old
	change.newValue = new
	change.instance = instance
	return change
}
