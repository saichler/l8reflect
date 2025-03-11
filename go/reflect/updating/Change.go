package updating

import (
	"github.com/saichler/reflect/go/reflect/property"
	"github.com/saichler/shared/go/share/strings"
)

type Change struct {
	property *properties.Property
	oldValue interface{}
	newValue interface{}
}

func (this *Change) String() (string, error) {
	id, err := this.property.PropertyId()
	if err != nil {
		return "", err
	}
	str := strings.New(id)

	str.Add(" - Old=").Add(str.StringOf(this.oldValue)).
		Add(" New=").Add(str.StringOf(this.newValue))
	return str.String(), nil
}

func (this *Change) Apply(any interface{}) {
	this.property.Set(any, this.newValue)
}

func (this *Change) PropertyId() string {
	id, _ := this.property.PropertyId()
	return id
}

func (this *Change) OldValue() interface{} {
	return this.oldValue
}

func (this *Change) NewValue() interface{} {
	return this.newValue
}

func NewChange(old, new interface{}, property *properties.Property) *Change {
	change := &Change{}
	change.oldValue = old
	change.newValue = new
	change.property = property
	return change
}
