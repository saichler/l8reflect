# Plan: DryUpdate Function

## Goal
Create a `DryUpdate` method on `Updater` that compares old and new instances and compiles a list of `Change` objects **without mutating the old instance**.

## Analysis

The current `Update` flow works as follows:
1. `Updater.Update()` resolves primary key, creates a root `Property`, calls `update()`
2. `update()` dispatches to type-specific comparators via `comparators[kind]`
3. Each comparator (int, string, bool, float, uint, ptr, struct, slice, map) does two things:
   - Detects differences and calls `updates.addUpdate(...)` to record changes
   - Mutates the old value via `oldValue.Set(newValue)` to apply the change

The mutation (`oldValue.Set(...)`) is scattered across **all 6 comparator files** in ~20 locations:
- `Comparators.go`: 5 primitives (int, uint, string, bool, float) — each has one `oldValue.Set(newValue)`
- `StructComparator.go`: `ptrUpdate` has 2 sets, `structUpdate` has 2 sets + 1 time-series set
- `MapComparator.go`: 5 `oldValue.Set(...)` / `oldValue.SetMapIndex(...)` calls + 1 delete
- `SliceComparator.go`: 5 `oldValue.Set(...)` / `oldIndexValue.Set(...)` calls

## Approach: Add a `dryRun` flag to `Updater`

The simplest and least invasive approach is adding a `dryRun bool` field to the `Updater` struct. When `dryRun` is true, all `oldValue.Set(...)` calls are skipped. The change detection and `addUpdate` calls remain unchanged.

### Why this approach
- **No code duplication** — we don't copy all comparators into a parallel set
- **No new comparator registry** — the same `comparators` map is used
- **Minimal diff** — each `oldValue.Set(...)` gets wrapped in `if !updates.dryRun {}`
- **Same Change list output** — identical behavior except no mutation

### Alternative considered: duplicate comparators
Creating a parallel set of "dry" comparators would avoid the conditional but would duplicate ~200 lines of comparison logic. Any future comparator change would need to be mirrored in two places. Rejected for maintainability.

## Changes

### 1. `Updater.go` — Add `dryRun` field and `DryUpdate` method

- Add `dryRun bool` field to the `Updater` struct
- Add `DryUpdate(old, new interface{}) error` method — identical to `Update` but sets `dryRun = true` before calling `update()`, then resets it
- Alternatively, `DryUpdate` can simply set the flag and delegate to `Update` since the flag controls behavior

```go
func (this *Updater) DryUpdate(old, new interface{}) error {
    this.dryRun = true
    defer func() { this.dryRun = false }()
    return this.Update(old, new)
}
```

### 2. `Comparators.go` — Guard primitive sets

Wrap each `oldValue.Set(newValue)` in the 5 primitive comparators:

```go
// Before:
oldValue.Set(newValue)

// After:
if !updates.dryRun {
    oldValue.Set(newValue)
}
```

Affected functions: `intUpdate`, `uintUpdate`, `stringUpdate`, `boolUpdate`, `floatUpdate`

### 3. `StructComparator.go` — Guard ptr and struct sets

- `ptrUpdate`: 2 `oldValue.Set(...)` calls (nil-to-value, value-to-nil)
- `structUpdate`: 2 `oldValue.Set(...)` calls (invalid-to-valid, valid-to-invalid) + 1 time-series `oldFldValue.Set(...)`

### 4. `MapComparator.go` — Guard map mutations

- `oldValue.Set(newValue)` (nil-to-value, value-to-nil, alwaysFull)
- `oldValue.SetMapIndex(key, newKeyValue)` (new key, modified non-struct key)
- `oldValue.SetMapIndex(key, reflect.Value{})` (deleted key)

Total: 6 guarded calls

### 5. `SliceComparator.go` — Guard slice mutations

- `oldValue.Set(newValue)` (nil-to-value, value-to-nil)
- `oldIndexValue.Set(newIndexValue)` (element updates)
- `oldValue.Set(newSlice)` (shrink and grow)

Total: 5 guarded calls

### 6. Tests — New test in `go/tests/`

Create a test that:
1. Creates two instances with known differences
2. Calls `DryUpdate(old, new)`
3. Asserts the change list is identical to what `Update` would produce
4. Asserts the old instance is **unchanged** (deep equal to its state before the call)

## Summary

| File | Change |
|------|--------|
| `Updater.go` | Add `dryRun` field, add `DryUpdate` method (~5 lines) |
| `Comparators.go` | Guard 5 `Set` calls |
| `StructComparator.go` | Guard 5 `Set` calls |
| `MapComparator.go` | Guard 6 `Set`/`SetMapIndex` calls |
| `SliceComparator.go` | Guard 5 `Set` calls |
| `go/tests/` | New dry update test file |

Total: ~21 `Set` calls wrapped in `if !updates.dryRun {}`, 1 new method, 1 new test file.
