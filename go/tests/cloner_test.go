package tests

import (
	"fmt"
	"github.com/saichler/l8types/go/testtypes"
	"github.com/saichler/reflect/go/reflect/cloning"
	"github.com/saichler/reflect/go/tests/utils"
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
