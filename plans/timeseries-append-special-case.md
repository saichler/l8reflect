# Plan: L8TimeSeriesPoint Append-Only Setter Special Case

## Problem
When `Set` is called on a property whose attribute type is `[]*l8api.L8TimeSeriesPoint`, the incoming value is a single `*l8api.L8TimeSeriesPoint`. Instead of replacing the slice by index (the current `sliceSet` behavior), the value should always be **appended** to the existing slice. When the slice reaches 100 elements, the oldest entry (index 0) should be dropped before appending.

## Current Flow
1. `Setter.go:Set()` (line 100-102) detects `this.node.IsSlice` and delegates to `sliceSet()`
2. `SliceSetter.go:sliceSet()` expects either a full slice replacement (`this.key == nil`) or an indexed set (`this.key` is an `int`). It does not support an "append" semantic.

## Detection
The node's `TypeName` for a `[]*l8api.L8TimeSeriesPoint` field will be `"L8TimeSeriesPoint"` (the element type name, as stored by the introspector). We can identify the special case by checking `this.node.TypeName == "L8TimeSeriesPoint"` combined with `this.node.IsSlice`.

## Plan

### Step 1: Add `timeSeriesAppend` method to `SliceSetter.go`

Add a new function in `SliceSetter.go`:

```go
const maxTimeSeriesPoints = 100

func (this *Property) timeSeriesAppend(myValue reflect.Value, newPoint reflect.Value) (interface{}, error) {
    // If the slice doesn't exist yet, create it
    if !myValue.IsValid() || myValue.IsNil() {
        info, err := this.resources.Registry().Info(this.node.TypeName)
        if err != nil {
            return nil, err
        }
        myValue.Set(reflect.MakeSlice(reflect.SliceOf(reflect.PointerTo(info.Type())), 0, 0))
    }

    // If slice is at capacity, drop the first element
    if myValue.Len() >= maxTimeSeriesPoints {
        myValue.Set(myValue.Slice(1, myValue.Len()))
    }

    // Append the new point
    myValue.Set(reflect.Append(myValue, newPoint))
    return myValue.Interface(), nil
}
```

### Step 2: Add the special case check in `Setter.go`

In the `Set()` method, **before** the existing `this.node.IsSlice` branch (line 100), add a check:

```go
} else if this.node.IsSlice && this.node.TypeName == "L8TimeSeriesPoint" {
    v, e := this.timeSeriesAppend(myValue, reflect.ValueOf(value))
    return v, any, e
}
```

This goes between the `this.node.IsMap` branch (line 97-99) and the existing `this.node.IsSlice` branch (line 100-102), so the time series case is handled first before falling through to normal slice logic.

### Step 3: Add a test

Create a test (or add to an existing test file, e.g. `updater_map_patch_test.go` or a new `timeseries_setter_test.go`) that:

1. Creates a `NetworkDevice` with an empty `CpuUsagePercent` slice
2. Calls `Set` with individual `*l8api.L8TimeSeriesPoint` values
3. Verifies append behavior (slice grows)
4. Fills the slice to 100 elements
5. Appends one more and verifies:
   - Slice length is still 100
   - The oldest point (original index 0) was dropped
   - The new point is at the end

### Files Changed

| File | Change |
|------|--------|
| `go/reflect/properties/Setter.go` | Add special case branch before `IsSlice` check |
| `go/reflect/properties/SliceSetter.go` | Add `timeSeriesAppend` method and `maxTimeSeriesPoints` const |
| `go/tests/` (new or existing test file) | Add test for append and cap-at-100 behavior |

### Considerations

- **Updater path**: The `SliceComparator.go` in the updater generates property changes that are later applied via `Set`. When a new `L8TimeSeriesPoint` arrives via the updater, it produces a change with the new point as the value. The `Set` method will now append instead of index-set, which is the desired behavior.
- **No impact on other slice types**: The check is gated on `TypeName == "L8TimeSeriesPoint"`, so all other slice types continue using the existing `sliceSet` logic.
- **The cap of 100 is hardcoded**: If this needs to be configurable in the future, it can be extracted, but for now a constant is sufficient.
