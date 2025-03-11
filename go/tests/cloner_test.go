package tests

import (
	"fmt"
	"github.com/saichler/reflect/go/reflect/cloning"
	"github.com/saichler/reflect/go/tests/utils"
	"github.com/saichler/types/go/testtypes"
	"testing"
)

func TestCloner(t *testing.T) {
	m := utils.CreateTestModelInstance(1)
	c := cloning.NewCloner().Clone(m).(*testtypes.TestProto)
	fmt.Println(c.MyString)
}
