package utils

import (
	"github.com/saichler/types/go/testtypes"
	"strconv"
	"time"
)

func CreateTestModelInstance(index int) *testtypes.TestProto {
	tag := strconv.Itoa(index)
	sub := &testtypes.TestProtoSub{
		MyString: "string-sub-" + tag,
		MyInt64:  time.Now().Unix(),
		MySubs:   make(map[string]*testtypes.TestProtoSubSub),
	}
	sub.MySubs["sub"] = &testtypes.TestProtoSubSub{MyString: "sub", Int32Map: make(map[int32]int32)}
	sub.MySubs["sub"].Int32Map[0] = 0
	sub.MySubs["sub"].Int32Map[1] = 0

	sub1 := &testtypes.TestProtoSub{
		MyString: "string-sub-1-" + tag,
		MyInt64:  time.Now().Unix(),
	}
	sub2 := &testtypes.TestProtoSub{
		MyString: "string-sub-2-" + tag,
		MyInt64:  time.Now().Unix(),
		MySubs:   make(map[string]*testtypes.TestProtoSubSub),
	}
	sub2.MySubs["sub2"] = &testtypes.TestProtoSubSub{MyString: "sub2-string-sub", Int32Map: make(map[int32]int32)}
	sub2.MySubs["sub2"].Int32Map[0] = 0
	sub2.MySubs["sub2"].Int32Map[1] = 0
	i := &testtypes.TestProto{
		MyString:           "string-" + tag,
		MyFloat64:          123456.123456,
		MyBool:             true,
		MyFloat32:          123.123,
		MyInt32:            int32(index),
		MyInt64:            int64(index * 10),
		MyInt32Slice:       []int32{1, 2, 3, int32(index)},
		MyStringSlice:      []string{"a", "b", "c", "d", tag},
		MyInt32ToInt64Map:  map[int32]int64{1: 11, 2: 22, 3: 33, 4: 44, int32(index): int64(index * 10)},
		MyString2StringMap: map[string]string{"a": "aa", "b": "bb", "c": "cc", tag: tag + tag},
		MySingle:           sub,
		MyModelSlice:       []*testtypes.TestProtoSub{sub1, sub2},
		MyString2ModelMap:  map[string]*testtypes.TestProtoSub{sub1.MyString: sub1, sub2.MyString: sub2},
		MyEnum:             testtypes.TestEnum_ValueOne,
	}
	return i
}
