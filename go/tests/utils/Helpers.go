package utils

import (
	"github.com/saichler/l8test/go/infra/t_resources"
	"github.com/saichler/l8types/go/testtypes"
)

func CreateTestModelInstance(index int) *testtypes.TestProto {
	return t_resources.CreateTestModelInstance(index)
}
