# L8Reflect - Model-Agnostic Reflection for Go

Most software needs a blueprint of your data before it can do anything with it. L8Reflect reads the blueprint automatically, so it can clone, compare, update, and navigate any data structure — even ones it has never seen before.

A Go library providing deep cloning, type introspection, property path navigation, and differential change tracking. Part of the **Layer 8 Model-Agnostic Infrastructure**.

## Features

**Introspection** — Struct analysis, node-based type trees with caching, decorator system (primary keys, unique fields, always-overwrite), thread-safe registry with mutex protection.

**Deep Cloning** — Type-safe deep clone of any Go struct, including pointers, slices, maps, and nested types. Automatic filtering of internal fields (`XXX`, `DoNotCompare`, `DoNotCopy`). Deep equality comparison.

**Property Navigation** — Path-based get/set on arbitrarily nested structs (e.g. `"order.lines<key>.amount"`). Supports maps, slices, and struct fields. Collect and ForEachValue traversal.

**Change Tracking & Updates** — Differential comparison between two instances of the same type. Records property-level `Change` objects. Supports dry-run mode (detect changes without mutating the original), nil-is-valid semantics, and full-replacement mode.

**Time Series Support** — Append-only slice semantics for time series fields. Thread-safe setters with mutex protection. Memory-leak-safe lifecycle management.

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

```
go/reflect/
  introspecting/  — Type introspection, L8Node tree, decorator registry (thread-safe)
  cloning/        — Deep clone and deep equality
  properties/     — Path-based get/set, collect, ForEachValue traversal
  updating/       — Differential update, dry-run, change recording
  helping/        — Value extraction, filtering utilities
```

## Testing

```bash
cd go && go test ./...
```

The test suite (~6.2K LOC, 26 test files) covers introspection, cloning, property get/set, map operations, updater diffs, dry-run updates, time series setters, and boolean edge cases.

## Dependencies

| Library | Purpose |
|---------|---------|
| `l8types` | Core type definitions and interfaces |
| `l8utils` | Utility functions (maps, strings, caching) |
| `l8srlz` | Serialization support |
| `l8test` | Testing utilities |
| `probler` | Protobuf reflection helpers |

## Advanced Features

### Dry-Run Updates

Detect changes without modifying the original object:

```go
updater := updating.NewUpdater(resources, false, true)
err := updater.DryUpdate(oldPerson, newPerson)
// oldPerson is unchanged; updater.Changes() contains the diff
```

### Decorators

```go
// Primary key — used by the updater to identify instances
introspecting.AddPrimaryKeyDecorator(node, "ServiceName", "ServiceArea")

// AlwaysOverwrite — replace entire map instead of merging keys
introspecting.AddAlwayOverwriteDecorator(node)

// Non — exclude a field from cloning/comparison
// Unique — mark a field as unique within a collection
```

### Property Path Navigation

Navigate complex object hierarchies using dot-delimited paths with map key syntax:

```go
// "person.address<123>.street" — traverse struct -> map key -> field
property, err := properties.PropertyOf("person.address<123>.street", resources)
```

### Field Filtering

Fields automatically skipped during cloning: `DoNotCompare`, `DoNotCopy`, `XXX*` prefixed, and unexported fields.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

© 2024-2026 Sharon Aicler (saichler@gmail.com)

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

### March 2026
- **Dry-Run Updates**: `DryUpdate()` detects changes without mutating the original instance
- **Time Series Support**: Append-only slice setters for time series data with thread-safe mutex protection
- **Memory Leak Fix**: Fixed lifecycle management for time series data structures
- **Byte Slice Handling**: Fixed cloning and comparison of `[]byte` fields
- **Boolean Edge Cases**: Fixed bool field comparison and update issues
- **Thread Safety**: Added mutex protection to introspector and decorator registry
- **Crash Prevention**: Multiple panic guards across cloning, properties, and time series paths

### December 2025
- **Non Decorator**: Exclude fields from operations
- **Unique Decorator**: Identify unique instances within collections
- **Display ID**: Better property identification and debugging
- **Introspection Refactoring**: Major refactoring for performance and maintainability
- **Decorator Refactoring**: Cleaner architecture for the decorator system

### December 2024
- **AlwaysOverwrite Decorator**: Force full updates on map structures
- **Multiple Primary Keys**: Composite primary key support
- **Performance Optimization**: Node cache key for faster introspection
- **Go Modules**: Added `go.mod`, vendored dependencies

### October 2024
- **Repository Rename**: Layer 8 Model Agnostic Infrastructure
- **Slice Handling**: Fixed slice operations for robust data handling
- **Updater Enhancements**: Better error handling
