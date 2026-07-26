package widgetlist

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
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

// TestSetItemsShrinksListKeepsFocusInRange covers the focus side of the same
// crash: SetItems hands a shorter widget list to the focus manager, which keeps
// its old tab index and then indexes past the new list on the next Update.
func TestSetItemsShrinksListKeepsFocusInRange(t *testing.T) {
	w := newTestList(t, 20, orvyn.NewSize(20, 10))

	// Move through Update, not NextItem: the focus manager only follows the
	// cursor from Update, and it is its index that goes stale.
	for range 15 {
		w.Update(tea.KeyMsg{Type: tea.KeyDown})
	}

	if index := w.focusManager.TabIndex(); index != 15 {
		t.Fatalf("test setup: tab index = %d, want 15", index)
	}

	w.SetItems([]string{"a", "b"})

	if index := w.focusManager.TabIndex(); index < 0 || index >= w.Length() {
		t.Errorf("tab index = %d, want inside [0, %d)", index, w.Length())
	}

	// The panic happens on the first message reaching the focus manager.
	w.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
}

// TestSetItemsEmptiesListThenUpdate covers the same path down to an empty list,
// where the focus manager has no widget left to index at all.
func TestSetItemsEmptiesListThenUpdate(t *testing.T) {
	w := newTestList(t, 20, orvyn.NewSize(20, 10))

	for range 15 {
		w.Update(tea.KeyMsg{Type: tea.KeyDown})
	}

	w.SetItems([]string{})

	w.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	w.Update(tea.KeyMsg{Type: tea.KeyDown})
	w.Update(tea.KeyMsg{Type: tea.KeyUp})
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
