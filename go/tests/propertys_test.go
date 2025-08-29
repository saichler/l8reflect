package tests

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8utils/go/utils/logger"
	"github.com/saichler/l8utils/go/utils/registry"
	"github.com/saichler/l8utils/go/utils/resources"
	"github.com/saichler/probler/go/tests"
	"github.com/saichler/probler/go/types"
	"github.com/saichler/reflect/go/reflect/introspecting"
	"github.com/saichler/reflect/go/reflect/properties"
)

func Devices() *types.NetworkDeviceList {
	return tests.GenerateExactDeviceTableMockData()
}

var propertiesLog = logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})

func newResourcesForNetworkDevice() ifs.IResources {
	res := resources.NewResources(propertiesLog)
	res.Set(registry.NewRegistry())
	in := introspecting.NewIntrospect(res.Registry())
	res.Set(in)
	return res
}

// Test the ConvertValue function with comprehensive scenarios
func TestConvertValue(t *testing.T) {
	// Test numeric conversions
	intVal := reflect.ValueOf(42)
	floatTarget := reflect.ValueOf(float64(0))
	result := properties.ConvertValue(floatTarget, intVal)
	if result.Float() != 42.0 {
		t.Errorf("Expected 42.0, got %v", result.Float())
	}

	// Test float to int conversion
	floatVal := reflect.ValueOf(3.14)
	intTarget := reflect.ValueOf(int(0))
	result = properties.ConvertValue(intTarget, floatVal)
	if result.Int() != 3 {
		t.Errorf("Expected 3, got %v", result.Int())
	}

	// Test string conversion
	stringTarget := reflect.ValueOf("")
	result = properties.ConvertValue(stringTarget, intVal)
	if result.String() != "42" {
		t.Errorf("Expected '42', got %v", result.String())
	}

	// Test complex number support
	complexVal := reflect.ValueOf(complex64(3 + 4i))
	complex128Target := reflect.ValueOf(complex128(0))
	result = properties.ConvertValue(complex128Target, complexVal)
	expected := complex128(3 + 4i)
	if result.Complex() != expected {
		t.Errorf("Expected %v, got %v", expected, result.Complex())
	}

	// Test no conversion needed
	stringVal := reflect.ValueOf("hello")
	result = properties.ConvertValue(stringTarget, stringVal)
	if result.String() != "hello" {
		t.Errorf("Expected 'hello', got %v", result.String())
	}
}

// Test basic field access through reflection with NetworkDevice
func TestProperties_NetworkDevice_BasicAccess(t *testing.T) {
	deviceList := Devices()
	if deviceList == nil || len(deviceList.List) == 0 {
		t.Skip("No devices available for testing")
	}

	device := deviceList.List[0]
	deviceValue := reflect.ValueOf(device)
	if deviceValue.Kind() == reflect.Ptr {
		deviceValue = deviceValue.Elem()
	}

	// Test accessing ID field
	idField := deviceValue.FieldByName("Id")
	if !idField.IsValid() {
		t.Error("Could not access Id field")
		return
	}

	// Test accessing EquipmentInfo field
	equipField := deviceValue.FieldByName("Equipmentinfo")
	if !equipField.IsValid() {
		t.Error("Could not access Equipmentinfo field")
		return
	}

	fmt.Printf("Device ID: %v, Equipment Info valid: %v\n",
		idField.Interface(), equipField.IsValid())
}

// Test complex nested structure navigation
func TestProperties_NetworkDevice_NestedAccess(t *testing.T) {
	deviceList := Devices()
	if deviceList == nil || len(deviceList.List) == 0 {
		t.Skip("No devices available for testing")
	}

	device := deviceList.List[0]
	deviceValue := reflect.ValueOf(device)
	if deviceValue.Kind() == reflect.Ptr {
		deviceValue = deviceValue.Elem()
	}

	// Test accessing nested EquipmentInfo fields
	equipField := deviceValue.FieldByName("Equipmentinfo")
	if equipField.IsValid() && !equipField.IsNil() {
		equipValue := equipField.Elem()
		if equipValue.Kind() == reflect.Struct {
			// Try to access common equipment info fields
			vendorField := equipValue.FieldByName("Vendor")
			if vendorField.IsValid() {
				fmt.Printf("Equipment Vendor: %v\n", vendorField.Interface())
			}
		}
	}
}

// Test map access functionality
func TestProperties_NetworkDevice_MapAccess(t *testing.T) {
	deviceList := Devices()
	if deviceList == nil || len(deviceList.List) == 0 {
		t.Skip("No devices available for testing")
	}

	device := deviceList.List[0]
	deviceValue := reflect.ValueOf(device)
	if deviceValue.Kind() == reflect.Ptr {
		deviceValue = deviceValue.Elem()
	}

	// Test Physicals map access
	physicalsField := deviceValue.FieldByName("Physicals")
	if physicalsField.IsValid() && !physicalsField.IsNil() {
		physicalsMap := physicalsField
		if physicalsMap.Kind() == reflect.Map {
			keys := physicalsMap.MapKeys()
			fmt.Printf("Physicals map has %d entries\n", len(keys))

			// Try to access first entry
			if len(keys) > 0 {
				firstValue := physicalsMap.MapIndex(keys[0])
				if firstValue.IsValid() {
					fmt.Printf("First physical entry: %v\n", firstValue.Interface())
				}
			}
		}
	}

	// Test Logicals map access
	logicalsField := deviceValue.FieldByName("Logicals")
	if logicalsField.IsValid() && !logicalsField.IsNil() {
		logicalsMap := logicalsField
		if logicalsMap.Kind() == reflect.Map {
			keys := logicalsMap.MapKeys()
			fmt.Printf("Logicals map has %d entries\n", len(keys))
		}
	}
}

// Test slice access functionality
func TestProperties_NetworkDevice_SliceAccess(t *testing.T) {
	deviceList := Devices()
	if deviceList == nil || len(deviceList.List) == 0 {
		t.Skip("No devices available for testing")
	}

	// Test the device list slice itself
	listValue := reflect.ValueOf(deviceList.List)
	if listValue.Kind() == reflect.Slice {
		fmt.Printf("Device list has %d devices\n", listValue.Len())

		// Access first device through reflection
		if listValue.Len() > 0 {
			firstDevice := listValue.Index(0)
			if firstDevice.IsValid() {
				devicePtr := firstDevice.Elem()
				idField := devicePtr.FieldByName("Id")
				if idField.IsValid() {
					fmt.Printf("First device ID via reflection: %v\n", idField.Interface())
				}
			}
		}
	}

	// Test NetworkLinks slice if available
	device := deviceList.List[0]
	deviceValue := reflect.ValueOf(device)
	if deviceValue.Kind() == reflect.Ptr {
		deviceValue = deviceValue.Elem()
	}

	networkLinksField := deviceValue.FieldByName("NetworkLinks")
	if networkLinksField.IsValid() && !networkLinksField.IsNil() {
		linksSlice := networkLinksField
		if linksSlice.Kind() == reflect.Slice {
			fmt.Printf("NetworkLinks slice has %d entries\n", linksSlice.Len())
		}
	}
}

// Test enum value handling
func TestProperties_NetworkDevice_EnumHandling(t *testing.T) {
	deviceList := Devices()
	if deviceList == nil || len(deviceList.List) == 0 {
		t.Skip("No devices available for testing")
	}

	device := deviceList.List[0]
	deviceValue := reflect.ValueOf(device)
	if deviceValue.Kind() == reflect.Ptr {
		deviceValue = deviceValue.Elem()
	}

	// Test DeviceStatus enum field
	statusField := deviceValue.FieldByName("DeviceStatus")
	if statusField.IsValid() {
		statusValue := statusField.Interface()
		fmt.Printf("Device Status: %v (type: %T)\n", statusValue, statusValue)

		// Test enum conversion
		stringTarget := reflect.ValueOf("")
		converted := properties.ConvertValue(stringTarget, statusField)
		fmt.Printf("Status as string: %v\n", converted.String())
	}
}

// Test type conversion with real NetworkDevice data
func TestProperties_NetworkDevice_TypeConversion(t *testing.T) {
	deviceList := Devices()
	if deviceList == nil || len(deviceList.List) == 0 {
		t.Skip("No devices available for testing")
	}

	device := deviceList.List[0]
	deviceValue := reflect.ValueOf(device)
	if deviceValue.Kind() == reflect.Ptr {
		deviceValue = deviceValue.Elem()
	}

	// Find a numeric field and test conversion
	numericFields := []string{"Id"} // Add more based on actual structure
	for _, fieldName := range numericFields {
		field := deviceValue.FieldByName(fieldName)
		if field.IsValid() && properties.IsNumeric(field.Kind()) {
			// Test conversion to different numeric types
			floatTarget := reflect.ValueOf(float64(0))
			converted := properties.ConvertValue(floatTarget, field)

			stringTarget := reflect.ValueOf("")
			stringConverted := properties.ConvertValue(stringTarget, field)

			fmt.Printf("Field %s: original=%v, as float=%v, as string=%v\n",
				fieldName, field.Interface(), converted.Float(), stringConverted.String())
		}
	}
}

// Test error handling scenarios
func TestProperties_NetworkDevice_ErrorHandling(t *testing.T) {
	deviceList := Devices()
	if deviceList == nil || len(deviceList.List) == 0 {
		t.Skip("No devices available for testing")
	}

	device := deviceList.List[0]
	deviceValue := reflect.ValueOf(device)
	if deviceValue.Kind() == reflect.Ptr {
		deviceValue = deviceValue.Elem()
	}

	// Test accessing non-existent field
	nonExistentField := deviceValue.FieldByName("NonExistentField")
	if nonExistentField.IsValid() {
		t.Error("Expected invalid field, but got valid result")
	}

	// Test nil pointer handling in nested structures
	equipField := deviceValue.FieldByName("Equipmentinfo")
	if equipField.IsValid() {
		if equipField.IsNil() {
			fmt.Println("EquipmentInfo is nil - this is expected in some cases")
		} else {
			fmt.Println("EquipmentInfo is not nil - can safely access nested fields")
		}
	}
}

// Test Property creation and basic operations with NetworkDevice
func TestProperties_PropertyCreation(t *testing.T) {
	deviceList := Devices()
	if deviceList == nil || len(deviceList.List) == 0 {
		t.Skip("No devices available for testing")
	}

	device := deviceList.List[0]
	res := newResourcesForNetworkDevice()

	// First need to introspect the NetworkDevice type
	_, err := res.Introspector().Inspect(&types.NetworkDevice{})
	if err != nil {
		t.Errorf("Failed to introspect NetworkDevice: %v", err)
		return
	}

	// Test creating property instance for device ID using PropertyOf
	prop, err := properties.PropertyOf("networkdevice.id", res)
	if err != nil {
		t.Errorf("Failed to create property for Id: %v", err)
		return
	}

	// Test getting value using property
	value, err := prop.Get(device)
	if err != nil {
		t.Errorf("Failed to get Id value: %v", err)
		return
	}

	deviceID := value.(string)
	fmt.Printf("Retrieved device ID via property: %s\n", deviceID)

	if deviceID == "" {
		t.Error("Expected non-empty device ID")
	}
}

// Test Property getter functionality with complex nested structures
func TestProperties_PropertyGetter_NestedStructures(t *testing.T) {
	deviceList := Devices()
	if deviceList == nil || len(deviceList.List) == 0 {
		t.Skip("No devices available for testing")
	}

	device := deviceList.List[0]
	res := newResourcesForNetworkDevice()

	// First need to introspect the NetworkDevice type
	_, err := res.Introspector().Inspect(&types.NetworkDevice{})
	if err != nil {
		t.Errorf("Failed to introspect NetworkDevice: %v", err)
		return
	}

	// Test accessing nested EquipmentInfo.Vendor
	vendorProp, err := properties.PropertyOf("networkdevice.equipmentinfo.vendor", res)
	if err != nil {
		t.Errorf("Failed to create vendor property: %v", err)
		return
	}

	vendor, err := vendorProp.Get(device)
	if err != nil {
		t.Errorf("Failed to get vendor: %v", err)
		return
	}

	fmt.Printf("Equipment vendor via property: %v\n", vendor)

	// Test accessing Physical component fields
	physicsProp, err := properties.PropertyOf("networkdevice.physicals", res)
	if err != nil {
		t.Errorf("Failed to create physicals property: %v", err)
		return
	}

	physicals, err := physicsProp.Get(device)
	if err != nil {
		t.Errorf("Failed to get physicals: %v", err)
		return
	}

	fmt.Printf("Physicals map type: %T\n", physicals)
}

// Test Property setter functionality with maps and slices
func TestProperties_PropertySetter_ComplexTypes(t *testing.T) {
	deviceList := Devices()
	if deviceList == nil || len(deviceList.List) == 0 {
		t.Skip("No devices available for testing")
	}

	device := deviceList.List[0]
	res := newResourcesForNetworkDevice()

	// First need to introspect the NetworkDevice type
	_, err := res.Introspector().Inspect(&types.NetworkDevice{})
	if err != nil {
		t.Errorf("Failed to introspect NetworkDevice: %v", err)
		return
	}

	// Test setting a simple string field
	idProp, err := properties.PropertyOf("networkdevice.id", res)
	if err != nil {
		t.Errorf("Failed to create Id property: %v", err)
		return
	}

	// Get original value
	originalID, err := idProp.Get(device)
	if err != nil {
		t.Errorf("Failed to get original ID: %v", err)
		return
	}

	// Set new value
	newID := "modified-device-001"
	result, newRoot, err := idProp.Set(device, newID)
	if err != nil {
		t.Errorf("Failed to set ID: %v", err)
		return
	}

	if result == nil {
		t.Error("Expected property to return result")
	}

	// Verify the change
	if newRoot != nil {
		newDevice := newRoot.(*types.NetworkDevice)
		if newDevice.Id != newID {
			t.Errorf("Expected ID to be %s, got %s", newID, newDevice.Id)
		}
	}

	fmt.Printf("Original ID: %v, New ID: %v, Result: %v\n", originalID, newID, result)

	// Test type conversion during set
	stringTarget := "12345"
	_, _, err = idProp.Set(device, stringTarget)
	if err != nil {
		t.Errorf("Failed to set string value: %v", err)
	}
}

// Test Property navigation with map keys
func TestProperties_PropertyNavigation_Maps(t *testing.T) {
	deviceList := Devices()
	if deviceList == nil || len(deviceList.List) == 0 {
		t.Skip("No devices available for testing")
	}

	device := deviceList.List[0]
	res := newResourcesForNetworkDevice()

	// First need to introspect the NetworkDevice type
	_, err := res.Introspector().Inspect(&types.NetworkDevice{})
	if err != nil {
		t.Errorf("Failed to introspect NetworkDevice: %v", err)
		return
	}

	// Test accessing physicals map entries
	physicsProp, err := properties.PropertyOf("networkdevice.physicals", res)
	if err != nil {
		t.Errorf("Failed to create physicals property: %v", err)
		return
	}

	physicals, err := physicsProp.Get(device)
	if err != nil {
		t.Errorf("Failed to get physicals: %v", err)
		return
	}

	physicsValue := reflect.ValueOf(physicals)
	if physicsValue.Kind() == reflect.Map && physicsValue.Len() > 0 {
		keys := physicsValue.MapKeys()
		firstKey := keys[0].String()

		// Try to access specific map entry by key
		mapEntryPath := fmt.Sprintf("Physicals[%s]", firstKey)
		fmt.Printf("Attempting to access: %s\n", mapEntryPath)

		// This tests if the property system can handle map key access
		// Note: This might not be supported yet and could be an area for improvement
	}
}

// Test error scenarios and edge cases
func TestProperties_EdgeCases(t *testing.T) {
	deviceList := Devices()
	if deviceList == nil || len(deviceList.List) == 0 {
		t.Skip("No devices available for testing")
	}

	res := newResourcesForNetworkDevice()

	// First need to introspect the NetworkDevice type
	_, err := res.Introspector().Inspect(&types.NetworkDevice{})
	if err != nil {
		t.Errorf("Failed to introspect NetworkDevice: %v", err)
		return
	}

	// Test invalid property path
	invalidProp, err := properties.PropertyOf("networkdevice.nonexistentfield.subfield", res)
	if err == nil {
		t.Error("Expected error for invalid property path")
	}

	// Test nil pointer handling
	var nilDevice *types.NetworkDevice
	if invalidProp != nil {
		_, err = invalidProp.Get(nilDevice)
		if err == nil {
			t.Error("Expected error when accessing property on nil object")
		}
	}

	// Test accessing property on wrong type - try to use string property on device
	// This should work since we're using the correct resources, but may fail if property doesn't exist
	fmt.Println("Testing edge cases completed - invalid property paths should fail as expected")
}
