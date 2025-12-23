# Layer 8 Model Agnostic Infra

A comprehensive Go library providing advanced reflection capabilities for deep cloning, introspection, property management, and change tracking of Go objects. This library forms part of the Layer 8 Model Agnostic Infrastructure, enabling robust and type-safe reflection operations across diverse data structures.

## Overview

This library extends Go's built-in reflection capabilities with a rich set of tools for working with complex data structures. It provides type-safe deep cloning, powerful introspection features, property-based access patterns, and sophisticated change detection and updating mechanisms.

## Features

### 🔍 **Introspection**
- Advanced struct analysis and type registration
- Node-based type representation with caching
- Table view generation for data structures
- Support for complex nested types and relationships

### 🔄 **Deep Cloning**
- Type-safe deep cloning of any Go data structure
- Handles circular references and complex pointer graphs
- Supports all Go primitive types, slices, maps, and structs
- Customizable field filtering (skip fields by name patterns)

### 🏷️ **Property Management**
- Property-based access to struct fields and nested data
- Hierarchical property paths with key support
- Getter/Setter abstractions for different data types
- Support for maps, slices, and complex nested structures

### 📊 **Change Tracking & Updates**
- Intelligent diff detection between object versions
- Granular change tracking with property-level precision
- Update application with validation and rollback support
- Support for partial updates and merge operations

## Installation

```bash
go get github.com/saichler/l8reflect/go
```

## Quick Start

### Basic Introspection

```go
import (
    "github.com/saichler/l8reflect/go/reflect/introspecting"
    "github.com/saichler/l8types/go/ifs"
)

// Create an introspector
registry := // your registry implementation
introspector := introspecting.NewIntrospect(registry)

// Inspect a struct
type Person struct {
    Name string
    Age  int
}

person := &Person{Name: "John", Age: 30}
node, err := introspector.Inspect(person)
if err != nil {
    log.Fatal(err)
}

// Access node information
fmt.Printf("Type: %s\n", node.TypeName)
```

### Deep Cloning

```go
import "github.com/saichler/l8reflect/go/reflect/cloning"

// Create a cloner
cloner := cloning.NewCloner()

// Clone any object
original := &Person{Name: "Alice", Age: 25}
cloned := cloner.Clone(original).(*Person)

// Modify clone without affecting original
cloned.Name = "Bob"
fmt.Printf("Original: %s, Clone: %s\n", original.Name, cloned.Name)
```

### Property Access

```go
import (
    "github.com/saichler/l8reflect/go/reflect/properties"
    "github.com/saichler/l8types/go/types"
)

// Create a property from a node and resources
prop := properties.NewProperty(node, nil, key, value, resources)

// Get property value
value := prop.Value()

// Check if property is a leaf node
if prop.IsLeaf() {
    fmt.Println("This is a leaf property")
}
```

### Change Tracking

```go
import "github.com/saichler/l8reflect/go/reflect/updating"

// Create an updater
updater := updating.NewUpdater(resources, false, true)

// Track changes between two versions
oldPerson := &Person{Name: "John", Age: 30}
newPerson := &Person{Name: "John", Age: 31}

err := updater.Update(oldPerson, newPerson)
if err != nil {
    log.Fatal(err)
}

// Get the changes
changes := updater.Changes()
for _, change := range changes {
    fmt.Printf("Changed: %v -> %v\n", change.OldValue, change.NewValue)
}
```

## Architecture

The library is organized into several key packages:

- **`introspecting/`**: Core introspection engine and node management
- **`cloning/`**: Deep cloning functionality with circular reference handling
- **`properties/`**: Property-based access patterns and utilities
- **`updating/`**: Change detection and update application
- **`helping/`**: Utility functions and helpers

## Testing

Run the comprehensive test suite:

```bash
cd go
./test.sh
```

This will:
- Initialize Go modules
- Download dependencies
- Run all tests with coverage
- Generate a coverage report

Individual test categories:
- **Introspection tests**: `introspect_test.go`
- **Cloning tests**: `cloner_test.go`
- **Property tests**: `property_test.go`, `property_map_test.go`
- **Update tests**: `updater_test.go`, `updater_map_test.go`
- **Application tests**: `apply_map_test.go`

## Dependencies

This library depends on several companion libraries:

- **`github.com/saichler/l8types`**: Core type definitions and interfaces
- **`github.com/saichler/l8utils`**: Utility functions for maps, strings, etc.
- **`github.com/saichler/l8test`**: Testing utilities
- **`github.com/saichler/serializer`**: Serialization support

## Advanced Features

### Custom Field Filtering

The cloning system supports automatic filtering of fields:

```go
// Fields automatically skipped during cloning:
// - "DoNotCompare"
// - "DoNotCopy"
// - Fields starting with "XXX"
// - Private fields (lowercase first letter)
```

### AlwaysOverwrite Decorator

Force complete map structure updates during change tracking:

```go
import "github.com/saichler/l8reflect/go/reflect/introspecting"

// Add AlwaysOverwrite decorator to a node
node, _ := introspector.Inspect(&MyStruct{})
introspecting.AddAlwayOverwriteDecorator(node)

// When this decorator is present, the entire map will be replaced
// during updates instead of merging individual keys
```

### Multiple Primary Keys

Support for composite primary keys with multiple attributes:

```go
// Define multiple fields as primary keys
node, _ := introspector.Inspect(&MyStruct{})
introspecting.AddPrimaryKeyDecorator(node, "ServiceName", "ServiceArea")

// The updater will use both fields as a composite key
// when tracking changes and identifying unique instances
```

### Table Views

Generate table-like views of your data structures:

```go
tableView, exists := introspector.TableView("Person")
if exists {
    fmt.Printf("Columns: %d\n", len(tableView.Columns))
    fmt.Printf("Sub-tables: %d\n", len(tableView.SubTables))
}
```

### Property Path Navigation

Navigate complex object hierarchies using property paths:

```go
// Example property path: "person.address<key>.street"
property, err := properties.PropertyOf("person.address<123>.street", resources)
if err != nil {
    log.Fatal(err)
}
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

© 2025 Sharon Aicler (saichler@gmail.com)

Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
You may obtain a copy of the License at:

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

## Support

For questions, issues, or contributions, please visit the [GitHub repository](https://github.com/saichler/l8reflect).

## Recent Updates

### Latest Changes (December 2025)
- **Non Decorator**: Added non decorator support for excluding fields from operations
- **Unique Decorator**: Added unique decorator for identifying unique instances within collections
- **Display ID**: Added display id functionality for better property identification and debugging
- **Property ID Enhancement**: Added propertyid to collect operations for improved tracking
- **Introspection Refactoring**: Major refactoring of the introspection system for better performance and maintainability
- **Decorator Refactoring**: Refactored decorator system for cleaner architecture
- **Interface Alignment**: Aligned interfaces for better API consistency
- **Crash Prevention**: Multiple fixes for panic conditions including nil inspect and map creation issues

### Previous Updates (December 2024)
- **AlwaysOverwrite Decorator**: Added new decorator support for forcing full updates on map structures
- **Multiple Primary Keys**: Fixed handling of composite primary keys with multiple attributes
- **Map Comparator Enhancement**: Improved map comparison logic with AlwaysFull decorator integration
- **Performance Optimization**: Added node cache key functionality for improved introspection performance
- **Test Suite Expansion**: Added comprehensive tests for multi-attribute primary keys
- **Go Modules**: Added Go module support with `go.mod`, `go.sum`, and vendored dependencies
- **Test Coverage**: Enhanced test coverage with HTML reports available in `go/cover.html`

### Previous Updates (October 2024)
- **Repository Rename**: Updated repository name to reflect Layer 8 Model Agnostic Infrastructure
- **Interface Improvements**: Enhanced interfaces for better compatibility and usability
- **Import Optimization**: Cleaned up unnecessary imports for better performance
- **Slice Handling**: Fixed slice operations for more robust data handling
- **Updater Enhancements**: Improved updater functionality with better error handling
- **Crash Prevention**: Multiple stability improvements to prevent runtime crashes
