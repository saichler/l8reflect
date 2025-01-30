package tests

import (
	"github.com/saichler/reflect/go/reflect/inspect"
	"github.com/saichler/shared/go/share/logger"
	"github.com/saichler/shared/go/share/registry"
	"github.com/saichler/shared/go/tests"
	"testing"
	"time"
)

var log = logger.NewLoggerImpl(&logger.FmtLogMethod{})

func TestIntrospect(t *testing.T) {
	defer time.Sleep(time.Second)
	m := &tests.TestProto{}

	in := inspect.NewIntrospect(registry.NewRegistry())

	_, err := in.Inspect(m)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}

	nodes := in.Nodes(false, false)
	expectedNodes := 17
	if len(nodes) != expectedNodes {
		log.Fail(t, "Expected length to be ", expectedNodes, " but got ", len(nodes))
		return
	}

	nodes = in.Nodes(false, true)
	if len(nodes) != 1 {
		log.Fail(t, "Expected length to be 1 roots but got ", len(nodes))
		return
	}

	nodes = in.Nodes(true, false)
	if len(nodes) != 13 {
		log.Fail(t, "Expected length to be 13 leafs but got ", len(nodes))
		return
	}

	_, ok := in.Node("testproto.myint32toint64map")
	if !ok {
		log.Fail(t, "Could not fetch node")
		return
	}

	_, ok = in.NodeByValue(&tests.TestProtoSub{})
	if !ok {
		log.Fail(t, "Could not fetch node by type")
		return
	}
}
