package tests

import (
	"testing"

	"github.com/saichler/probler/go/types"
	"github.com/saichler/reflect/go/reflect/updating"
)

func TestPatchMapItem(t *testing.T) {
	res := newResources()
	res.Introspector().Inspect(&types.NetworkDevice{})
	aside := types.NetworkDevice{Physicals: map[string]*types.Physical{"1": &types.Physical{Ports: []*types.Port{&types.Port{Id: "id"}}}}}
	zside := types.NetworkDevice{Physicals: map[string]*types.Physical{"1": &types.Physical{Performance: &types.PerformanceMetrics{CpuUsagePercent: 88.0}}}}

	updater := updating.NewUpdater(res, false, false)

	err := updater.Update(aside, zside)
	if err != nil {
		res.Logger().Fail(t, err.Error())
		return
	}

	if len(aside.Physicals["1"].Ports) == 0 {
		res.Logger().Fail(t, "Expected ports")
		return
	}
}
