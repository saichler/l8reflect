package tests

import (
	"fmt"
	"github.com/saichler/reflect/go/reflect/clone"
	"github.com/saichler/reflect/go/tests/utils"
	"github.com/saichler/shared/go/tests"
	"testing"
)

func TestCloner(t *testing.T) {
	m := utils.CreateTestModelInstance(1)
	c := clone.NewCloner().Clone(m).(*tests.TestProto)
	fmt.Println(c.MyString)
}
