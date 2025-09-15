package list

import (
	"github.com/halsten-dev/orvyn/widget/textinput"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/theme"
	"github.com/halsten-dev/orvyn/widget"
)

type IListItem interface {
	orvyn.Focusable
	orvyn.Renderable
	FilterValue() string
}

// ItemConstructor defines the signature of the item constructor.
// T type represents the type of the item data.
type ItemConstructor[T any] func(*T) IListItem

type filteredItem struct {
	index int // corresponding global index
	item  *IListItem
}

type filteredItems []filteredItem

// FilterState Taken from github.com/charmbracelet/bubbles/list/list.go
// FilterState describes the current filtering state on the model.
type FilterState int

// Possible filter states.
const (
	Unfiltered    FilterState = iota // no filter set
	Filtering                        // user is actively setting a filter
	FilterApplied                    // a filter is applied and user is not editing filter
)

// String returns a human-readable string of the current filter state.
func (f FilterState) String() string {
	return [...]string{
		"unfiltered",
		"filtering",
		"filter applied",
	}[f]
}

type keybinds struct {
	cursorUp    key.Binding
	cursorDown  key.Binding
	enterFilter key.Binding
	clearFilter key.Binding
	applyFilter key.Binding
}

// Widget defines a list widget.
// T type represents the type of the item data.
type Widget[T any] struct {
	orvyn.BaseWidget
	orvyn.BaseFocusable

	InfiniteScroll bool
	filterable     bool
	filterState    FilterState

	cursor      int
	globalIndex int

	MinSize       orvyn.Size
	PreferredSize orvyn.Size

	listItems         []IListItem
	filteredListItems filteredItems
	items             []T

	tiFilter *textinput.Widget

	paginator paginator.Model

	focusManager *orvyn.FocusManager

	itemConstructor ItemConstructor[T]

	style lipgloss.Style

	contentSize   orvyn.Size
	maxItemHeight int

	keybinds keybinds

	CursorMovedCallback func(int)
}

// New creates a new *Widget list and takes an itemConstructor as parameter.
// T type represents the type of the item data.
func New[T any](itemConstructor ItemConstructor[T]) *Widget[T] {
	w := new(Widget[T])

	w.BaseWidget = orvyn.NewBaseWidget()

	w.keybinds = keybinds{
		cursorUp:    key.NewBinding(key.WithKeys("up")),
		cursorDown:  key.NewBinding(key.WithKeys("down")),
		enterFilter: key.NewBinding(key.WithKeys("/")),
		applyFilter: key.NewBinding(key.WithKeys("enter")),
		clearFilter: key.NewBinding(key.WithKeys("esc")),
	}

	w.itemConstructor = itemConstructor

	w.InfiniteScroll = false
	w.filterable = true
	w.filterState = Unfiltered

	w.cursor = 0

	w.tiFilter = textinput.New()
	w.tiFilter.Placeholder = "Press '/' to filter"

	w.paginator = paginator.New()
	w.paginator.Type = paginator.Dots

	w.focusManager = orvyn.NewFocusManager()
	w.focusManager.ManageFocusNextPrevKeybind = false
	w.focusManager.PreviousFocusKeybind = w.keybinds.cursorUp
	w.focusManager.NextFocusKeybind = w.keybinds.cursorDown
	w.focusManager.Focus(0)

	w.MinSize = orvyn.NewSize(10, 5)
	w.PreferredSize = orvyn.NewSize(20, 10)

	w.maxItemHeight = 1

	w.OnBlur()

	return w
}

func (w *Widget[T]) Init() tea.Cmd {
	w.clearFilter()

	w.focusManager.FocusFirst()

	return nil
}

func (w *Widget[T]) Update(msg tea.Msg) tea.Cmd {
	if w.filterState == Filtering {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, w.keybinds.applyFilter):
				w.basicFilter(w.tiFilter.Value())

				return nil

			case key.Matches(msg, w.keybinds.clearFilter):
				w.clearFilter()

				return nil
			}
		}

		cmd := w.tiFilter.Update(msg)

		return cmd
	}

	isInputting := w.checkInputting()

	if !isInputting {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, w.focusManager.PreviousFocusKeybind):
				w.PreviousItem()

				if w.CursorMovedCallback != nil {
					w.CursorMovedCallback(w.globalIndex)
				}

			case key.Matches(msg, w.focusManager.NextFocusKeybind):
				w.NextItem()

				if w.CursorMovedCallback != nil {
					w.CursorMovedCallback(w.globalIndex)
				}

			case key.Matches(msg, w.keybinds.enterFilter):
				if w.filterable {
					w.enterFilter()

					return nil
				}

			case key.Matches(msg, w.keybinds.clearFilter):
				if w.filterState == FilterApplied {
					w.clearFilter()

					return nil
				}
			}
		}
	}

	cmd := w.focusManager.Update(msg)

	w.focusManager.Focus(w.globalIndex)

	return cmd
}

func (w *Widget[T]) Resize(size orvyn.Size) {
	w.maxItemHeight = 1

	w.BaseWidget.Resize(size)

	size.Width -= w.style.GetHorizontalFrameSize()
	size.Height -= w.style.GetVerticalFrameSize()

	w.tiFilter.Resize(size)

	w.contentSize = size

	w.paginatorUpdate()
}

func (w *Widget[T]) paginatorUpdate() {
	var perPage int

	for _, li := range w.listItems {
		li.Resize(w.contentSize)

		w.maxItemHeight = max(w.maxItemHeight, li.GetSize().Height)
	}

	calcHeight := w.contentSize.Height - 1 // paginator

	if w.filterable {
		calcHeight -= w.tiFilter.GetSize().Height
	}

	perPage = calcHeight / w.maxItemHeight
	perPage = max(perPage, 1)

	w.paginator.PerPage = perPage

	w.paginator.TotalPages = 0

	if w.filterState == FilterApplied {
		w.paginator.SetTotalPages(len(w.filteredListItems))
	} else {
		w.paginator.SetTotalPages(len(w.listItems))
	}
}

func (w *Widget[T]) Render() string {
	var elements []string
	var b strings.Builder
	var view string
	var start, end int

	elements = make([]string, 0)

	if w.filterState == FilterApplied {
		start, end = w.paginator.GetSliceBounds(len(w.filteredListItems))

		for i, li := range w.filteredListItems[start:end] {
			if i > 0 {
				b.WriteString("\n")
			}

			item := *li.item
			b.WriteString(item.Render())
		}
	} else {
		start, end = w.paginator.GetSliceBounds(len(w.listItems))

		for i, li := range w.listItems[start:end] {
			if i > 0 {
				b.WriteString("\n")
			}

			b.WriteString(li.Render())
		}
	}

	if w.filterable {
		elements = append(elements, w.tiFilter.Render())
	}

	elements = append(elements, b.String())

	if w.paginator.TotalPages > 1 {
		elements = append(elements, w.paginator.View())
	}

	view = lipgloss.JoinVertical(lipgloss.Center,
		elements...)

	return w.style.
		Width(w.contentSize.Width).
		Height(w.contentSize.Height).
		Render(view)
}

func (w *Widget[T]) GetMinSize() orvyn.Size {
	return w.MinSize
}

func (w *Widget[T]) GetPreferredSize() orvyn.Size {
	return w.PreferredSize
}

func (w *Widget[T]) OnFocus() {
	w.style = orvyn.GetTheme().Style(theme.FocusedWidgetStyleID)

	widget.UpdatePaginatorTheme(&w.paginator)
}

func (w *Widget[T]) OnBlur() {
	w.style = orvyn.GetTheme().Style(theme.BlurredWidgetStyleID)

	widget.UpdatePaginatorTheme(&w.paginator)
}

func (w *Widget[T]) OnEnterInput() {}

func (w *Widget[T]) OnExitInput() {}

func (w *Widget[T]) IsInputting() bool {
	inputting := w.checkInputting()

	if inputting {
		return true
	}

	return w.BaseFocusable.IsInputting()
}

func (w *Widget[T]) checkInputting() bool {
	for _, item := range w.listItems {
		if item.IsInputting() {
			return true
		}
	}

	return false
}

// Public API

// PreviousItem manages the focus of the previous item.
func (w *Widget[T]) PreviousItem() {
	if len(w.listItems) == 0 {
		return
	}

	if w.filterState == FilterApplied {
		w.previousFilteredItem()
		return
	}

	w.globalIndex--

	if w.globalIndex < 0 {
		if w.InfiniteScroll {
			w.globalIndex = len(w.items) - 1
			w.moveCursor(w.globalIndex)
			return
		}

		w.globalIndex = 0
		w.moveCursor(0)
		return
	}

	w.moveCursor(w.globalIndex)
}

func (w *Widget[T]) previousFilteredItem() {
	if len(w.filteredListItems) == 0 {
		return
	}

	w.cursor--

	if w.cursor < 0 && w.paginator.Page == 0 {
		if w.InfiniteScroll {
			w.paginator.Page = w.paginator.TotalPages - 1
			w.cursor = w.paginator.ItemsOnPage(len(w.filteredListItems)) - 1
			w.globalIndex = w.getFilteredGlobalIndex()
			return
		}

		w.cursor = 0
		w.globalIndex = w.getFilteredGlobalIndex()
		return
	}

	w.globalIndex = w.getFilteredGlobalIndex()

	if w.cursor >= 0 {
		return
	}

	w.paginator.PrevPage()
	w.cursor = w.paginator.ItemsOnPage(len(w.filteredListItems)) - 1
	w.globalIndex = w.getFilteredGlobalIndex()
}

// NextItem manages the focus of the next item.
func (w *Widget[T]) NextItem() {
	if len(w.listItems) == 0 {
		return
	}

	if w.filterState == FilterApplied {
		w.nextFilteredItem()
		return
	}

	w.globalIndex++

	if w.globalIndex > len(w.items)-1 {
		if w.InfiniteScroll {
			w.globalIndex = 0
			w.moveCursor(0)
			return
		}

		w.globalIndex = len(w.items) - 1
		w.moveCursor(w.globalIndex)
		return
	}

	w.moveCursor(w.globalIndex)
}

func (w *Widget[T]) nextFilteredItem() {
	if len(w.filteredListItems) == 0 {
		return
	}

	itemsOnPage := w.paginator.ItemsOnPage(len(w.filteredListItems))

	w.cursor++

	if w.cursor >= itemsOnPage && w.paginator.OnLastPage() {
		if w.InfiniteScroll {
			w.paginator.Page = 0
			w.cursor = 0
			w.globalIndex = w.getFilteredGlobalIndex()
			return
		}

		w.cursor = itemsOnPage - 1
		w.globalIndex = w.getFilteredGlobalIndex()
		return
	}

	w.globalIndex = w.getFilteredGlobalIndex()

	if w.cursor <= itemsOnPage-1 {
		return
	}

	w.paginator.NextPage()
	w.cursor = 0
	w.globalIndex = w.getFilteredGlobalIndex()
}

func (w *Widget[T]) moveCursor(index int) {
	// based on the global index set the cursor and the current page.

	itemsOnPage := w.paginator.PerPage

	page := int(math.Floor(float64(index) / float64(itemsOnPage)))
	cursor := index % itemsOnPage

	w.paginator.Page = page
	w.cursor = cursor
}

func (w *Widget[T]) SetFilterable(filterable bool) {
	w.filterable = filterable
	w.tiFilter.SetActive(filterable)
}

func (w *Widget[T]) SetFilterPlaceholder(s string) {
	w.tiFilter.Placeholder = s
}

func (w *Widget[T]) GetGlobalIndex() int {
	return w.globalIndex
}

func (w *Widget[T]) RemoveItem(index int) {
	if index < 0 || index >= len(w.items) {
		return
	}

	w.items = append(w.items[:index], w.items[index+1:]...)
	w.listItems = append(w.listItems[:index], w.listItems[index+1:]...)
	w.focusManager.Remove(index)

	if w.filterState == FilterApplied {
		w.basicFilter(w.tiFilter.Value())
	}
}

// SetItems takes a []T (slice of data) and instantiate all items
// based on it.
func (w *Widget[T]) SetItems(items []T) {
	w.items = items

	w.listItems = make([]IListItem, 0)
	focusableList := make([]orvyn.Focusable, 0)

	for i := range w.items {
		item := w.itemConstructor(&w.items[i])
		w.listItems = append(w.listItems,
			item)
		focusableList = append(focusableList,
			item)
	}

	w.focusManager.SetWidgets(focusableList)

	w.paginatorUpdate()
}

func (w *Widget[T]) SetCursorMovementKeybinds(cursorUp, cursorDown key.Binding) {
	w.keybinds.cursorUp = cursorUp
	w.keybinds.cursorDown = cursorDown
	w.focusManager.PreviousFocusKeybind = cursorUp
	w.focusManager.NextFocusKeybind = cursorDown
}

func (w *Widget[T]) GetItems() []T {
	return w.items
}

func (w *Widget[T]) GetSelectedItem() T {
	return w.items[w.globalIndex]
}

func (w *Widget[T]) SetItem(index int, data T) {
	w.items[index] = data
}

func (w *Widget[T]) AppendItem(data T) {
	w.clearFilter()

	w.items = append(w.items, data)

	w.SetItems(w.items)

	w.globalIndex = len(w.items) - 1
	w.moveCursor(w.globalIndex)
	w.focusManager.Focus(w.globalIndex)
}

func (w *Widget[T]) InsertItem(index int, data T) {
	w.clearFilter()

	w.items = append(w.items[:index+1], w.items[index:]...)

	w.items[index] = data

	w.SetItems(w.items)

	w.globalIndex = index
	w.moveCursor(w.globalIndex)
	w.focusManager.Focus(w.globalIndex)
}

func (w *Widget[T]) FocusFirst() {
	w.focusManager.FocusFirst()

	if w.filterState == FilterApplied {
		if len(w.filteredListItems) > 0 {
			w.globalIndex = w.filteredListItems[0].index
			w.cursor = 0
			w.paginator.Page = 0
		}
	} else {
		w.globalIndex = 0
		w.moveCursor(w.globalIndex)
	}

	if w.CursorMovedCallback != nil {
		w.CursorMovedCallback(w.globalIndex)
	}

}

func (w *Widget[T]) BlurCurrent() {
	w.focusManager.BlurCurrent()
}

func (w *Widget[T]) FilterState() FilterState {
	return w.filterState
}

func (w *Widget[T]) basicFilter(s string) {
	if s == "" {
		w.clearFilter()
	}

	w.tiFilter.OnBlur()

	w.filteredListItems = make(filteredItems, 0)

	for i, v := range w.listItems {
		if strings.Contains(v.FilterValue(), s) {
			w.filteredListItems = append(w.filteredListItems, filteredItem{
				index: i,
				item:  &w.listItems[i],
			})
			v.SetActive(true)
			continue
		}

		v.SetActive(false)
	}

	w.filterState = FilterApplied

	w.FocusFirst()

	w.paginatorUpdate()
}

func (w *Widget[T]) clearFilter() {
	w.tiFilter.SetValue("")
	w.tiFilter.OnBlur()

	w.filteredListItems = make(filteredItems, 0)

	for _, v := range w.listItems {
		v.SetActive(true)
	}

	w.filterState = Unfiltered

	w.FocusFirst()

	w.paginatorUpdate()
}

func (w *Widget[T]) enterFilter() {
	w.focusManager.BlurCurrent()
	w.tiFilter.OnFocus()
	w.filterState = Filtering

	w.paginatorUpdate()
}

func (w *Widget[T]) getFilteredGlobalIndex() int {
	index := (w.paginator.Page * w.paginator.PerPage) + w.cursor
	return w.filteredListItems[index].index
}
