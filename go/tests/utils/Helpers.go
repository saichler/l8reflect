package utils

import (
	"github.com/saichler/shared/go/tests"
	"strconv"
	"time"
)

func CreateTestModelInstance(index int) *tests.TestProto {
	tag := strconv.Itoa(index)
	sub := &tests.TestProtoSub{
		MyString: "string-sub-" + tag,
		MyInt64:  time.Now().Unix(),
	}
	sub1 := &tests.TestProtoSub{
		MyString: "string-sub-1-" + tag,
		MyInt64:  time.Now().Unix(),
	}
	sub2 := &tests.TestProtoSub{
		MyString: "string-sub-2-" + tag,
		MyInt64:  time.Now().Unix(),
	}
	i := &tests.TestProto{
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
		MyModelSlice:       []*tests.TestProtoSub{sub1, sub2},
		MyString2ModelMap:  map[string]*tests.TestProtoSub{sub1.MyString: sub1, sub2.MyString: sub2},
		MyEnum:             tests.TestEnum_ValueOne,
	}
	return i
}
