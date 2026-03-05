# Plan: Time Series Short-Circuit in Updater

## Problem

When the Updater compares two objects containing `[]*L8TimeSeriesPoint` slices, it drills down into individual points and their scalar fields (stamp, value). This produces per-field changes like `temperature<0>.value = 35` (int64). When these changes are later applied via `Change.Apply` → `Property.Set`, the setter's time series branch at `Setter.go:100` receives the raw scalar value (int64) instead of an `*L8TimeSeriesPoint`, causing:

```
panic: reflect.Set: value of type int64 is not assignable to type *l8api.L8TimeSeriesPoint
```

## Root Cause

The Updater treats `[]*L8TimeSeriesPoint` like any other struct slice — it compares element-by-element, then field-by-field inside each element. Time series data is append-only; the Updater should never drill into individual points.

## Solution

Short-circuit in `structUpdate` (StructComparator.go): when iterating a struct's attributes, if an attribute node is a slice with `TypeName == "L8TimeSeriesPoint"`, skip the normal `update()` dispatch. Instead, treat the new slice value as the whole change — record a single update with `newValue` as the entire `[]*L8TimeSeriesPoint` slice and let the Setter's existing `timeSeriesAppend` handle append + capacity logic.

## Changes

### 1. `go/reflect/updating/StructComparator.go` — `structUpdate()`

Inside the attribute loop (line 70-78), before calling `update()`, add a check:

```go
for _, attr := range node.Attributes {
    oldFldValue := oldValue.FieldByName(attr.FieldName)
    newFldValue := newValue.FieldByName(attr.FieldName)

    // Time series slices are append-only: record the whole new slice
    // as a single change and let timeSeriesAppend handle merging.
    if attr.IsSlice && attr.TypeName == "L8TimeSeriesPoint" {
        if !newFldValue.IsValid() || newFldValue.IsNil() || newFldValue.Len() == 0 {
            continue
        }
        subInstance := properties.NewProperty(attr, property, nil, oldFldValue, updates.resources)
        updates.addUpdate(subInstance, nil, newFldValue.Interface())
        oldFldValue.Set(newFldValue)
        continue
    }

    subInstance := properties.NewProperty(attr, property, nil, oldFldValue, updates.resources)
    err := update(subInstance, attr, oldFldValue, newFldValue, updates)
    if err != nil {
        return err
    }
}
```

**What this does:**
- When the new time series slice is nil/empty, skip it (no change).
- Otherwise, record a single change with the whole `[]*L8TimeSeriesPoint` slice as `newValue`.
- Copy the new slice into old (so the old object is updated).
- The Setter's `timeSeriesAppend` will receive the full `[]*L8TimeSeriesPoint` slice and handle appending + capacity enforcement.

### 2. `go/tests/timeseries_setter_test.go` — Add updater test

Add a test `TestTimeSeriesUpdater` that:
1. Creates two device instances with different time series data.
2. Runs `Updater.Update(old, new)`.
3. Verifies the changes list contains a single change for the time series property (not per-field changes).
4. Applies the changes to a third instance and verifies the time series data is correct.

## What Does NOT Change

- `SliceSetter.go` (`timeSeriesAppend`) — already handles both single `*L8TimeSeriesPoint` and full `[]*L8TimeSeriesPoint` slice values correctly.
- `Setter.go` — no changes needed. The time series branch at line 100 will now receive the correct type (`[]*L8TimeSeriesPoint` from the updater change) instead of a raw scalar.
- `SliceComparator.go` (`sliceUpdate`) — the slice comparator is never reached for time series because the short-circuit happens one level up in `structUpdate`.
