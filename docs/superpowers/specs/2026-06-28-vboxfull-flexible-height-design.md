# VBoxFull flexible-height elements — design

Date: 2026-06-28

## Goal

Let `VBoxFullLayout` distribute height across **N flexible-height elements**, the way
`DefinedWidthVerticalLayout` already does, instead of expanding only a single
`growIndex` element. Expose this as **opt-in** via new constructors. Existing
constructors and their single-grow behavior stay unchanged.

## Background

Current behavior (`layout/vboxfull.go`):

- `NewVBoxFullLayout(margin, growIndex, elements...)` and
  `NewMaxWidthVBoxFullLayout(margin, growIndex, elements...)`.
- One designated `growWidget` (chosen by `growIndex`) takes all leftover vertical
  space. Every other element is rendered at its `GetMinSize().Height`.

Reference behavior (`layout/definedwidthvertical.go`):

- Classifies each element at construction time:
  - `pref.Height == 0` → fixed
  - `min.Height == pref.Height` → fixed
  - otherwise (`min.Height < pref.Height`) → flexible
- Distributes leftover height across flexible elements proportionally to their
  preferred heights.

## Decisions

- Detection: **auto-detect** via `min.Height < pref.Height` (mirror DefinedWidthVertical).
- Width: provide **both** variants (normal + max-width), paralleling the existing pair.
- Remainder: leftover pixels from rounding go to the **last flexible element** so the
  layout fully fills its height (it is "VBoxFull"). This is better than the reference,
  which can under-fill.
- DefinedWidthVertical ratio bug (`definedwidthvertical.go` ratio denominator): **leave
  as-is**. Out of scope.
- Tests: **skip** (repo currently has no layout tests).

## API

New constructors in `layout/vboxfull.go` (old two untouched):

```go
// Auto-detects flexible-height elements (min.Height < pref.Height) and shares
// leftover height among them proportionally to their preferred heights.
func NewFlexibleVBoxFullLayout(margin orvyn.Size, elements ...orvyn.Renderable) *VBoxFullLayout

// Same, but ignores width min/preferred constraints (full width).
func NewMaxWidthFlexibleVBoxFullLayout(margin orvyn.Size, elements ...orvyn.Renderable) *VBoxFullLayout
```

## Struct changes

```go
type VBoxFullLayout struct {
    orvyn.BaseLayout

    margin     orvyn.Size
    growWidget orvyn.Renderable // single-grow mode only (nil in auto-flex mode)
    maxWidth   bool

    autoFlex               bool
    fixedHeightElements    []orvyn.Renderable
    flexibleHeightElements []flexibleHeightElement // type reused from definedwidthvertical.go (same package)
}
```

The new constructors set `autoFlex = true`, leave `growWidget = nil`, and classify
`elements` into `fixedHeightElements` / `flexibleHeightElements` using the same switch
as DefinedWidthVertical.

## Render algorithm (auto-flex path)

`Render()` branches on `l.autoFlex`. Single-grow path is unchanged.

1. Fit width (shared helper, see below).
2. Resize every fixed element to its own height (`GetMinSize().Height`, falling back to
   current size when min is 0 — match how `calculateGrowSize` already reads heights).
3. `remaining = layoutHeight − Σ(fixed heights) − margin.Height`, floored at 0.
4. `totalFlexPref = Σ(flexible elements' pref.Height)`.
5. For each flexible element except the last:
   `h = round(remaining × pref_i / totalFlexPref)`, resize, subtract `h` from a running
   `left` counter (initialized to `remaining`).
6. Last flexible element gets the running `left` (absorbs rounding remainder).
7. If `totalFlexPref == 0` (no flexible elements), nothing grows; leftover space is
   unused. Same as reference.
8. Render elements in original order, `\n`-joined.

## Refactor

Extract the width-fit block (`vboxfull.go:55-64`) into a small helper, e.g.
`fitWidth(layoutWidth int) int`, returning the element width after applying
`maxWidth` / min / preferred constraints. Both render paths call it. No behavior change
for the existing single-grow path.

`GetMinSize` / `GetPreferredSize` are unchanged — they already sum element heights and
max element widths, which is correct for both modes.

## Backward compatibility

- `NewVBoxFullLayout` / `NewMaxWidthVBoxFullLayout` signatures and behavior unchanged.
- New fields default to zero values (`autoFlex=false`, empty slices) for old
  constructors, so the single-grow path is selected exactly as before.

## Out of scope

- DefinedWidthVertical ratio bug.
- Tests.
- Any explicit-index or configurable flexible-set API.
