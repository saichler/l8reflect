package tests

import (
	"github.com/saichler/reflect/go/reflect/introspecting"
	"github.com/saichler/l8utils/go/utils/logger"
	"github.com/saichler/l8utils/go/utils/registry"
	"github.com/saichler/l8types/go/testtypes"
	"testing"
	"time"
)

var log = logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})
var flog = logger.NewLoggerDirectImpl(logger.NewFileLogMethod("./test.log"))

func TestIntrospect(t *testing.T) {
	defer time.Sleep(time.Second)
	m := &testtypes.TestProto{}

	in := introspecting.NewIntrospect(registry.NewRegistry())

	_, err := in.Inspect(m)
	if err != nil {
		log.Fail(t, err.Error())
		return
	}

	nodes := in.Nodes(false, false)
	expectedNodes := 33
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
	if len(nodes) != 26 {
		log.Fail(t, "Expected length to be 26 leafs but got ", len(nodes))
		return
	}

	_, ok := in.Node("testproto.myint32toint64map")
	if !ok {
		log.Fail(t, "Could not fetch node")
		return
	}

	_, ok = in.NodeByValue(&testtypes.TestProtoSub{})
	if !ok {
		log.Fail(t, "Could not fetch node by type")
		return
	}

	_, ok = in.Node("testproto.mystring2modelmap.mystring")
	if !ok {
		log.Fail(t, "Could not fetch node")
		return
	}
}
