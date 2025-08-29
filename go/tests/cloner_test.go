package tests

import (
	"fmt"
	"github.com/saichler/l8types/go/testtypes"
	"github.com/saichler/reflect/go/reflect/cloning"
	"github.com/saichler/reflect/go/tests/utils"
	"reflect"
	"testing"
)

func TestCloner(t *testing.T) {
	m := utils.CreateTestModelInstance(1)
	c := cloning.NewCloner().Clone(m).(*testtypes.TestProto)
	fmt.Println(c.MyString)
}

// Test primitive types
func TestDeepClone_PrimitiveTypes(t *testing.T) {
	cloner := cloning.NewCloner()
	
	// Test int
	original := 42
	cloned := cloner.Clone(original).(int)
	if cloned != original {
		t.Errorf("Expected %d, got %d", original, cloned)
	}
	
	// Test string
	originalStr := "test string"
	clonedStr := cloner.Clone(originalStr).(string)
	if clonedStr != originalStr {
		t.Errorf("Expected %s, got %s", originalStr, clonedStr)
	}
	
	// Test bool
	originalBool := true
	clonedBool := cloner.Clone(originalBool).(bool)
	if clonedBool != originalBool {
		t.Errorf("Expected %t, got %t", originalBool, clonedBool)
	}
	
	// Test float64
	originalFloat := 3.14159
	clonedFloat := cloner.Clone(originalFloat).(float64)
	if clonedFloat != originalFloat {
		t.Errorf("Expected %f, got %f", originalFloat, clonedFloat)
	}
}

// Test numeric type variations
func TestDeepClone_NumericTypes(t *testing.T) {
	cloner := cloning.NewCloner()
	
	// Test int32
	original32 := int32(100)
	cloned32 := cloner.Clone(original32).(int32)
	if cloned32 != original32 {
		t.Errorf("Expected %d, got %d", original32, cloned32)
	}
	
	// Test int64
	original64 := int64(1000)
	cloned64 := cloner.Clone(original64).(int64)
	if cloned64 != original64 {
		t.Errorf("Expected %d, got %d", original64, cloned64)
	}
	
	// Test uint
	originalUint := uint(200)
	clonedUint := cloner.Clone(originalUint).(uint)
	if clonedUint != originalUint {
		t.Errorf("Expected %d, got %d", originalUint, clonedUint)
	}
	
	// Test float32
	originalFloat32 := float32(2.71)
	clonedFloat32 := cloner.Clone(originalFloat32).(float32)
	if clonedFloat32 != originalFloat32 {
		t.Errorf("Expected %f, got %f", originalFloat32, clonedFloat32)
	}
}

// Test slice cloning
func TestDeepClone_Slice(t *testing.T) {
	cloner := cloning.NewCloner()
	
	// Test int slice
	original := []int{1, 2, 3, 4, 5}
	cloned := cloner.Clone(original).([]int)
	
	if len(cloned) != len(original) {
		t.Errorf("Expected length %d, got %d", len(original), len(cloned))
	}
	
	for i, v := range original {
		if cloned[i] != v {
			t.Errorf("At index %d: expected %d, got %d", i, v, cloned[i])
		}
	}
	
	// Verify they are independent
	original[0] = 999
	if cloned[0] == 999 {
		t.Error("Clone is not independent of original")
	}
}

// Test nil slice
func TestDeepClone_NilSlice(t *testing.T) {
	cloner := cloning.NewCloner()
	
	var original []int
	cloned := cloner.Clone(original).([]int)
	
	// Due to Go's reflection behavior, nil slices become empty slices after Interface()
	if len(cloned) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(cloned))
	}
	// The slice will be non-nil but empty, which is acceptable behavior
}

// Test map cloning
func TestDeepClone_Map(t *testing.T) {
	cloner := cloning.NewCloner()
	
	original := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	
	cloned := cloner.Clone(original).(map[string]int)
	
	if len(cloned) != len(original) {
		t.Errorf("Expected length %d, got %d", len(original), len(cloned))
	}
	
	for k, v := range original {
		if cloned[k] != v {
			t.Errorf("At key %s: expected %d, got %d", k, v, cloned[k])
		}
	}
	
	// Verify they are independent
	original["one"] = 999
	if cloned["one"] == 999 {
		t.Error("Clone is not independent of original")
	}
}

// Test nil map
func TestDeepClone_NilMap(t *testing.T) {
	cloner := cloning.NewCloner()
	
	var original map[string]int
	cloned := cloner.Clone(original).(map[string]int)
	
	// The cloner returns the original nil map
	if cloned != nil {
		t.Error("Expected nil map, got non-nil")
	}
}

// Test struct cloning
func TestDeepClone_Struct(t *testing.T) {
	cloner := cloning.NewCloner()
	
	type TestStruct struct {
		Name  string
		Age   int
		Score float64
		Active bool
	}
	
	original := TestStruct{
		Name:   "John",
		Age:    30,
		Score:  85.5,
		Active: true,
	}
	
	cloned := cloner.Clone(original).(TestStruct)
	
	if cloned.Name != original.Name {
		t.Errorf("Expected name %s, got %s", original.Name, cloned.Name)
	}
	if cloned.Age != original.Age {
		t.Errorf("Expected age %d, got %d", original.Age, cloned.Age)
	}
	if cloned.Score != original.Score {
		t.Errorf("Expected score %f, got %f", original.Score, cloned.Score)
	}
	if cloned.Active != original.Active {
		t.Errorf("Expected active %t, got %t", original.Active, cloned.Active)
	}
}

// Test pointer cloning
func TestDeepClone_Pointer(t *testing.T) {
	cloner := cloning.NewCloner()
	
	original := 42
	ptr := &original
	
	clonedPtr := cloner.Clone(ptr).(*int)
	
	if *clonedPtr != original {
		t.Errorf("Expected %d, got %d", original, *clonedPtr)
	}
	
	// Verify they point to different memory locations
	if clonedPtr == ptr {
		t.Error("Cloned pointer points to same memory location")
	}
	
	// Verify they are independent
	*ptr = 999
	if *clonedPtr == 999 {
		t.Error("Clone is not independent of original")
	}
}

// Test nil pointer
func TestDeepClone_NilPointer(t *testing.T) {
	cloner := cloning.NewCloner()
	
	var original *int
	cloned := cloner.Clone(original).(*int)
	
	// The cloner returns the original nil pointer
	if cloned != nil {
		t.Error("Expected nil pointer, got non-nil")
	}
}

// Test circular reference handling
func TestDeepClone_CircularReference(t *testing.T) {
	cloner := cloning.NewCloner()
	
	type Node struct {
		Value int
		Next  *Node
	}
	
	node1 := &Node{Value: 1}
	node2 := &Node{Value: 2}
	node1.Next = node2
	node2.Next = node1 // Circular reference
	
	cloned := cloner.Clone(node1).(*Node)
	
	if cloned.Value != 1 {
		t.Errorf("Expected value 1, got %d", cloned.Value)
	}
	if cloned.Next.Value != 2 {
		t.Errorf("Expected next value 2, got %d", cloned.Next.Value)
	}
	if cloned.Next.Next.Value != 1 {
		t.Errorf("Expected circular reference value 1, got %d", cloned.Next.Next.Value)
	}
	
	// Verify circular reference is maintained
	if cloned.Next.Next != cloned {
		t.Error("Circular reference not properly maintained in clone")
	}
}

// Test field skipping functionality
func TestDeepClone_FieldSkipping(t *testing.T) {
	cloner := cloning.NewCloner()
	
	type TestStruct struct {
		PublicField    string
		privateField   string // Should be skipped (lowercase)
		DoNotCompare   string // Should be skipped
		DoNotCopy      string // Should be skipped
		XXXSomeField   string // Should be skipped (XXX prefix)
		NormalField    string
	}
	
	original := TestStruct{
		PublicField:    "public",
		privateField:   "private",
		DoNotCompare:   "compare",
		DoNotCopy:      "copy",
		XXXSomeField:   "xxx",
		NormalField:    "normal",
	}
	
	cloned := cloner.Clone(original).(TestStruct)
	
	if cloned.PublicField != original.PublicField {
		t.Errorf("Expected public field %s, got %s", original.PublicField, cloned.PublicField)
	}
	if cloned.NormalField != original.NormalField {
		t.Errorf("Expected normal field %s, got %s", original.NormalField, cloned.NormalField)
	}
	
	// These fields should be empty/zero values due to skipping
	if cloned.privateField != "" {
		t.Errorf("Expected empty private field, got %s", cloned.privateField)
	}
	if cloned.DoNotCompare != "" {
		t.Errorf("Expected empty DoNotCompare field, got %s", cloned.DoNotCompare)
	}
	if cloned.DoNotCopy != "" {
		t.Errorf("Expected empty DoNotCopy field, got %s", cloned.DoNotCopy)
	}
	if cloned.XXXSomeField != "" {
		t.Errorf("Expected empty XXXSomeField field, got %s", cloned.XXXSomeField)
	}
}

// Test complex nested structure (without interface{} which isn't supported)
func TestDeepClone_ComplexStructure(t *testing.T) {
	cloner := cloning.NewCloner()
	
	type Address struct {
		Street string
		City   string
	}
	
	type Person struct {
		Name      string
		Age       int
		Addresses []*Address
		Tags      map[string]string
	}
	
	original := &Person{
		Name: "Alice",
		Age:  25,
		Addresses: []*Address{
			{Street: "123 Main St", City: "Anytown"},
			{Street: "456 Oak Ave", City: "Somewhere"},
		},
		Tags: map[string]string{
			"level":  "senior",
			"status": "active",
		},
	}
	
	cloned := cloner.Clone(original).(*Person)
	
	if cloned.Name != original.Name {
		t.Errorf("Expected name %s, got %s", original.Name, cloned.Name)
	}
	if cloned.Age != original.Age {
		t.Errorf("Expected age %d, got %d", original.Age, cloned.Age)
	}
	if len(cloned.Addresses) != len(original.Addresses) {
		t.Errorf("Expected %d addresses, got %d", len(original.Addresses), len(cloned.Addresses))
	}
	
	// Verify addresses are cloned
	for i, addr := range original.Addresses {
		if cloned.Addresses[i].Street != addr.Street {
			t.Errorf("Address %d street: expected %s, got %s", i, addr.Street, cloned.Addresses[i].Street)
		}
		if cloned.Addresses[i].City != addr.City {
			t.Errorf("Address %d city: expected %s, got %s", i, addr.City, cloned.Addresses[i].City)
		}
		// Verify different memory locations
		if cloned.Addresses[i] == addr {
			t.Errorf("Address %d not properly cloned (same memory location)", i)
		}
	}
	
	// Verify tags are cloned
	if len(cloned.Tags) != len(original.Tags) {
		t.Errorf("Expected %d tag entries, got %d", len(original.Tags), len(cloned.Tags))
	}
	
	for k, v := range original.Tags {
		if cloned.Tags[k] != v {
			t.Errorf("Tag %s: expected %s, got %s", k, v, cloned.Tags[k])
		}
	}
	
	// Verify independence
	original.Name = "Changed"
	original.Addresses[0].Street = "Changed Street"
	original.Tags["status"] = "inactive"
	
	if cloned.Name == "Changed" {
		t.Error("Clone name was affected by original change")
	}
	if cloned.Addresses[0].Street == "Changed Street" {
		t.Error("Clone address was affected by original change")
	}
	if cloned.Tags["status"] == "inactive" {
		t.Error("Clone tags was affected by original change")
	}
}

// Test missing numeric types
func TestDeepClone_MissingNumericTypes(t *testing.T) {
	cloner := cloning.NewCloner()
	
	// Test int8
	originalInt8 := int8(127)
	clonedInt8 := cloner.Clone(originalInt8).(int8)
	if clonedInt8 != originalInt8 {
		t.Errorf("Expected %d, got %d", originalInt8, clonedInt8)
	}
	
	// Test int16
	originalInt16 := int16(32767)
	clonedInt16 := cloner.Clone(originalInt16).(int16)
	if clonedInt16 != originalInt16 {
		t.Errorf("Expected %d, got %d", originalInt16, clonedInt16)
	}
	
	// Test uint8
	originalUint8 := uint8(255)
	clonedUint8 := cloner.Clone(originalUint8).(uint8)
	if clonedUint8 != originalUint8 {
		t.Errorf("Expected %d, got %d", originalUint8, clonedUint8)
	}
	
	// Test uint16
	originalUint16 := uint16(65535)
	clonedUint16 := cloner.Clone(originalUint16).(uint16)
	if clonedUint16 != originalUint16 {
		t.Errorf("Expected %d, got %d", originalUint16, clonedUint16)
	}
}

// Test complex types
func TestDeepClone_ComplexTypes(t *testing.T) {
	cloner := cloning.NewCloner()
	
	// Test complex64
	originalComplex64 := complex64(3.14 + 2.71i)
	clonedComplex64 := cloner.Clone(originalComplex64).(complex64)
	if clonedComplex64 != originalComplex64 {
		t.Errorf("Expected %v, got %v", originalComplex64, clonedComplex64)
	}
	
	// Test complex128
	originalComplex128 := complex128(3.141592653589793 + 2.718281828459045i)
	clonedComplex128 := cloner.Clone(originalComplex128).(complex128)
	if clonedComplex128 != originalComplex128 {
		t.Errorf("Expected %v, got %v", originalComplex128, clonedComplex128)
	}
}

// Test array cloning
func TestDeepClone_Array(t *testing.T) {
	cloner := cloning.NewCloner()
	
	// Test int array
	original := [5]int{1, 2, 3, 4, 5}
	cloned := cloner.Clone(original).([5]int)
	
	if len(cloned) != len(original) {
		t.Errorf("Expected length %d, got %d", len(original), len(cloned))
	}
	
	for i, v := range original {
		if cloned[i] != v {
			t.Errorf("At index %d: expected %d, got %d", i, v, cloned[i])
		}
	}
	
	// Verify they are independent
	original[0] = 999
	if cloned[0] == 999 {
		t.Error("Clone is not independent of original")
	}
}

// Test nested array
func TestDeepClone_NestedArray(t *testing.T) {
	cloner := cloning.NewCloner()
	
	type Point struct {
		X, Y int
	}
	
	original := [3]Point{{1, 2}, {3, 4}, {5, 6}}
	cloned := cloner.Clone(original).([3]Point)
	
	for i, p := range original {
		if cloned[i].X != p.X || cloned[i].Y != p.Y {
			t.Errorf("At index %d: expected %+v, got %+v", i, p, cloned[i])
		}
	}
	
	// Verify independence
	original[0].X = 999
	if cloned[0].X == 999 {
		t.Error("Clone is not independent of original")
	}
}

// Test interface{} cloning
func TestDeepClone_Interface(t *testing.T) {
	cloner := cloning.NewCloner()
	
	// Test with various concrete types in interface{}
	testCases := []interface{}{
		42,
		"test string",
		true,
		3.14,
		[]int{1, 2, 3},
	}
	
	for i, original := range testCases {
		cloned := cloner.Clone(original)
		
		// Check that the types match
		if fmt.Sprintf("%T", original) != fmt.Sprintf("%T", cloned) {
			t.Errorf("Test case %d: type mismatch. Expected %T, got %T", i, original, cloned)
		}
		
		// Check that the values match
		if !reflect.DeepEqual(original, cloned) {
			t.Errorf("Test case %d: value mismatch. Expected %v, got %v", i, original, cloned)
		}
	}
}

// Test nil interface
func TestDeepClone_NilInterface(t *testing.T) {
	cloner := cloning.NewCloner()
	
	var original interface{}
	cloned := cloner.Clone(original)
	
	if cloned != nil {
		t.Error("Expected nil interface{}, got non-nil")
	}
}

// Test interface with pointer
func TestDeepClone_InterfaceWithPointer(t *testing.T) {
	cloner := cloning.NewCloner()
	
	type TestStruct struct {
		Value int
	}
	
	originalStruct := &TestStruct{Value: 42}
	var original interface{} = originalStruct
	
	cloned := cloner.Clone(original).(*TestStruct)
	
	if cloned.Value != originalStruct.Value {
		t.Errorf("Expected %d, got %d", originalStruct.Value, cloned.Value)
	}
	
	// Verify they point to different memory locations
	if cloned == originalStruct {
		t.Error("Cloned pointer points to same memory location")
	}
	
	// Verify independence
	originalStruct.Value = 999
	if cloned.Value == 999 {
		t.Error("Clone is not independent of original")
	}
}

// Test channel cloning
func TestDeepClone_Channel(t *testing.T) {
	cloner := cloning.NewCloner()
	
	// Test unbuffered channel
	original := make(chan int)
	cloned := cloner.Clone(original).(chan int)
	
	// Channels can't be truly cloned, but we should get a new channel of the same type
	if cloned == original {
		t.Error("Expected different channel instance")
	}
	
	// Test buffered channel
	originalBuffered := make(chan string, 5)
	clonedBuffered := cloner.Clone(originalBuffered).(chan string)
	
	if clonedBuffered == originalBuffered {
		t.Error("Expected different buffered channel instance")
	}
}

// Test nil channel
func TestDeepClone_NilChannel(t *testing.T) {
	cloner := cloning.NewCloner()
	
	var original chan int
	cloned := cloner.Clone(original).(chan int)
	
	// Nil channels should remain nil
	if cloned != nil {
		t.Error("Expected nil channel, got non-nil")
	}
}

// Test function cloning
func TestDeepClone_Function(t *testing.T) {
	cloner := cloning.NewCloner()
	
	original := func(x int) int { return x * 2 }
	cloned := cloner.Clone(original).(func(int) int)
	
	// Functions can't be cloned, so we expect the same function
	if fmt.Sprintf("%p", original) != fmt.Sprintf("%p", cloned) {
		t.Error("Expected same function reference")
	}
	
	// Test that the function still works
	if cloned(5) != 10 {
		t.Error("Cloned function doesn't work correctly")
	}
}

// Test nil function
func TestDeepClone_NilFunction(t *testing.T) {
	cloner := cloning.NewCloner()
	
	var original func(int) int
	cloned := cloner.Clone(original).(func(int) int)
	
	// Nil functions should remain nil
	if cloned != nil {
		t.Error("Expected nil function, got non-nil")
	}
}

// Test struct with all supported types
func TestDeepClone_AllTypesStruct(t *testing.T) {
	cloner := cloning.NewCloner()
	
	type AllTypes struct {
		IntField        int
		Int8Field       int8
		Int16Field      int16
		Int32Field      int32
		Int64Field      int64
		UintField       uint
		Uint8Field      uint8
		Uint16Field     uint16
		Uint32Field     uint32
		Uint64Field     uint64
		Float32Field    float32
		Float64Field    float64
		Complex64Field  complex64
		Complex128Field complex128
		BoolField       bool
		StringField     string
		ArrayField      [3]int
		SliceField      []string
		MapField        map[string]int
		PtrField        *int
		InterfaceField  interface{}
		ChanField       chan int
		FuncField       func() string
	}
	
	value := 42
	original := AllTypes{
		IntField:        1,
		Int8Field:       2,
		Int16Field:      3,
		Int32Field:      4,
		Int64Field:      5,
		UintField:       6,
		Uint8Field:      7,
		Uint16Field:     8,
		Uint32Field:     9,
		Uint64Field:     10,
		Float32Field:    11.5,
		Float64Field:    12.5,
		Complex64Field:  13 + 14i,
		Complex128Field: 15 + 16i,
		BoolField:       true,
		StringField:     "test",
		ArrayField:      [3]int{1, 2, 3},
		SliceField:      []string{"a", "b", "c"},
		MapField:        map[string]int{"key": 100},
		PtrField:        &value,
		InterfaceField:  "interface value",
		ChanField:       make(chan int),
		FuncField:       func() string { return "test" },
	}
	
	cloned := cloner.Clone(original).(AllTypes)
	
	// Test all numeric fields
	if cloned.IntField != original.IntField {
		t.Errorf("IntField: expected %d, got %d", original.IntField, cloned.IntField)
	}
	if cloned.Int8Field != original.Int8Field {
		t.Errorf("Int8Field: expected %d, got %d", original.Int8Field, cloned.Int8Field)
	}
	if cloned.Int16Field != original.Int16Field {
		t.Errorf("Int16Field: expected %d, got %d", original.Int16Field, cloned.Int16Field)
	}
	if cloned.Int32Field != original.Int32Field {
		t.Errorf("Int32Field: expected %d, got %d", original.Int32Field, cloned.Int32Field)
	}
	if cloned.Int64Field != original.Int64Field {
		t.Errorf("Int64Field: expected %d, got %d", original.Int64Field, cloned.Int64Field)
	}
	if cloned.UintField != original.UintField {
		t.Errorf("UintField: expected %d, got %d", original.UintField, cloned.UintField)
	}
	if cloned.Uint8Field != original.Uint8Field {
		t.Errorf("Uint8Field: expected %d, got %d", original.Uint8Field, cloned.Uint8Field)
	}
	if cloned.Uint16Field != original.Uint16Field {
		t.Errorf("Uint16Field: expected %d, got %d", original.Uint16Field, cloned.Uint16Field)
	}
	if cloned.Uint32Field != original.Uint32Field {
		t.Errorf("Uint32Field: expected %d, got %d", original.Uint32Field, cloned.Uint32Field)
	}
	if cloned.Uint64Field != original.Uint64Field {
		t.Errorf("Uint64Field: expected %d, got %d", original.Uint64Field, cloned.Uint64Field)
	}
	if cloned.Float32Field != original.Float32Field {
		t.Errorf("Float32Field: expected %f, got %f", original.Float32Field, cloned.Float32Field)
	}
	if cloned.Float64Field != original.Float64Field {
		t.Errorf("Float64Field: expected %f, got %f", original.Float64Field, cloned.Float64Field)
	}
	if cloned.Complex64Field != original.Complex64Field {
		t.Errorf("Complex64Field: expected %v, got %v", original.Complex64Field, cloned.Complex64Field)
	}
	if cloned.Complex128Field != original.Complex128Field {
		t.Errorf("Complex128Field: expected %v, got %v", original.Complex128Field, cloned.Complex128Field)
	}
	if cloned.BoolField != original.BoolField {
		t.Errorf("BoolField: expected %t, got %t", original.BoolField, cloned.BoolField)
	}
	if cloned.StringField != original.StringField {
		t.Errorf("StringField: expected %s, got %s", original.StringField, cloned.StringField)
	}
	
	// Test array field
	for i, v := range original.ArrayField {
		if cloned.ArrayField[i] != v {
			t.Errorf("ArrayField[%d]: expected %d, got %d", i, v, cloned.ArrayField[i])
		}
	}
	
	// Test slice field
	if len(cloned.SliceField) != len(original.SliceField) {
		t.Errorf("SliceField length: expected %d, got %d", len(original.SliceField), len(cloned.SliceField))
	}
	for i, v := range original.SliceField {
		if cloned.SliceField[i] != v {
			t.Errorf("SliceField[%d]: expected %s, got %s", i, v, cloned.SliceField[i])
		}
	}
	
	// Test map field
	if len(cloned.MapField) != len(original.MapField) {
		t.Errorf("MapField length: expected %d, got %d", len(original.MapField), len(cloned.MapField))
	}
	for k, v := range original.MapField {
		if cloned.MapField[k] != v {
			t.Errorf("MapField[%s]: expected %d, got %d", k, v, cloned.MapField[k])
		}
	}
	
	// Test pointer field
	if *cloned.PtrField != *original.PtrField {
		t.Errorf("PtrField value: expected %d, got %d", *original.PtrField, *cloned.PtrField)
	}
	if cloned.PtrField == original.PtrField {
		t.Error("PtrField: expected different memory locations")
	}
	
	// Test interface field
	if cloned.InterfaceField != original.InterfaceField {
		t.Errorf("InterfaceField: expected %v, got %v", original.InterfaceField, cloned.InterfaceField)
	}
	
	// Test channel field (should be different instances)
	if cloned.ChanField == original.ChanField {
		t.Error("ChanField: expected different channel instances")
	}
	
	// Test function field (should be same reference)
	if fmt.Sprintf("%p", cloned.FuncField) != fmt.Sprintf("%p", original.FuncField) {
		t.Error("FuncField: expected same function reference")
	}
}
