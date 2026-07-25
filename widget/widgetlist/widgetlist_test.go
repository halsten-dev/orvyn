package widgetlist

import (
	"fmt"
	"testing"

	"github.com/halsten-dev/orvyn"
)

func newTestList(t *testing.T, count int, size orvyn.Size) *Widget[string] {
	t.Helper()

	orvyn.Init()

	w := New(SimpleListItemConstructor)
	w.SetFilterable(false)
	w.Resize(size)

	items := make([]string, 0, count)

	for i := range count {
		items = append(items, fmt.Sprintf("item %d", i))
	}

	w.SetItems(items)

	return w
}

func checkBounds(t *testing.T, w *Widget[string]) {
	t.Helper()

	length := len(w.listItems)

	if w.filterState == FilterApplied {
		length = len(w.filteredListItems)
	}

	start, end := w.paginator.GetSliceBounds(length)

	if start > end {
		t.Fatalf("slice bounds out of range [%d:%d] (len %d, page %d, perPage %d, totalPages %d)",
			start, end, length, w.paginator.Page, w.paginator.PerPage, w.paginator.TotalPages)
	}
}

// TestSetItemsShrinksList covers the reported crash: the list is replaced by a
// shorter one while the cursor sits on a page that no longer exists.
func TestSetItemsShrinksList(t *testing.T) {
	w := newTestList(t, 20, orvyn.NewSize(20, 10))

	for range 15 {
		w.NextItem()
	}

	if w.paginator.Page == 0 {
		t.Fatalf("test setup: cursor never left the first page")
	}

	w.SetItems([]string{"a", "b", "c", "d", "e", "f", "g"})

	checkBounds(t, w)

	if w.globalIndex != 6 {
		t.Errorf("global index = %d, want 6 (last item of the new list)", w.globalIndex)
	}

	w.Render()
}

// TestSetItemsEmptiesList covers replacing the list with an empty one.
func TestSetItemsEmptiesList(t *testing.T) {
	w := newTestList(t, 20, orvyn.NewSize(20, 10))

	for range 15 {
		w.NextItem()
	}

	w.SetItems([]string{})

	checkBounds(t, w)

	w.Render()
}

// TestRemoveItemsDownToOne removes items one by one from the last page.
func TestRemoveItemsDownToOne(t *testing.T) {
	w := newTestList(t, 20, orvyn.NewSize(20, 10))

	for range 19 {
		w.NextItem()
	}

	for w.Length() > 1 {
		w.RemoveItem(w.Length() - 1)

		checkBounds(t, w)

		w.Render()
	}
}

// TestResizeChangesPerPage covers the other way the page index goes stale: the
// item count is unchanged but the per-page capacity grows or shrinks.
func TestResizeChangesPerPage(t *testing.T) {
	heights := []int{40, 6, 25, 5, 60, 7}

	for _, from := range heights {
		for _, to := range heights {
			w := newTestList(t, 17, orvyn.NewSize(20, from))

			for range 16 {
				w.NextItem()
			}

			w.Resize(orvyn.NewSize(20, to))

			checkBounds(t, w)

			if w.globalIndex != 16 {
				t.Errorf("resize %d -> %d: global index = %d, want 16",
					from, to, w.globalIndex)
			}

			w.Render()
		}
	}
}

// TestFilteredListShrinks covers the filtered branch of Render, where the
// bounds are computed against the filtered items.
func TestFilteredListShrinks(t *testing.T) {
	w := newTestList(t, 20, orvyn.NewSize(20, 10))
	w.SetFilterable(true)

	w.filter("item 1")

	for range 10 {
		w.NextItem()
	}

	checkBounds(t, w)

	w.SetItems([]string{"item 1", "item 2"})
	w.filter("item 1")

	checkBounds(t, w)

	w.Render()

	w.clearFilter()

	checkBounds(t, w)

	w.Render()
}
