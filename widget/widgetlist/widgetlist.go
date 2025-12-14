package widgetlist

import (
	"math"
	"strings"

	"github.com/halsten-dev/orvyn/widget/textinput"
	"github.com/sahilm/fuzzy"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/widget"
)

type ListItem[T any] interface {
	orvyn.Focusable
	orvyn.Renderable
	FilterValue() string
	UpdateData(data T)
	GetData() T
}

// ItemConstructor defines the signature of the item constructor.
// T type represents the type of the item data.
type ItemConstructor[T any] func(T) ListItem[T]

type FilteredItem struct {
	Index int // corresponding global index
}

type FilteredItems []FilteredItem

type ListFilter[T any] func(items *[]ListItem[T], s string) FilteredItems

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

// Widget defines a widgetlist widget.
// T type represents the type of the item data.
type Widget[T any] struct {
	orvyn.BaseWidget
	orvyn.BaseFocusable

	InfiniteScroll            bool
	AutoFocusNewItem          bool
	filterable                bool
	blockCursorMovingCallback bool
	filterState               FilterState

	cursor      int
	globalIndex int

	listItems         []ListItem[T]
	filteredListItems FilteredItems

	tiFilter *textinput.Widget

	paginator paginator.Model

	focusManager *orvyn.FocusManager

	itemConstructor ItemConstructor[T]

	maxItemHeight int

	keybinds keybinds

	CursorMovingCallback func(int)
	CursorMovedCallback  func(int)

	Filter ListFilter[T]
}

// New creates a new *Widget widgetlist and takes an itemConstructor as parameter.
// T type represents the type of the item data.
func New[T any](itemConstructor ItemConstructor[T]) *Widget[T] {
	w := new(Widget[T])

	w.BaseWidget = orvyn.NewBaseWidget()
	w.BaseFocusable = orvyn.NewBaseFocusable(w)

	w.keybinds = keybinds{
		cursorUp:    key.NewBinding(key.WithKeys("up", "k")),
		cursorDown:  key.NewBinding(key.WithKeys("down", "j")),
		enterFilter: key.NewBinding(key.WithKeys("/")),
		applyFilter: key.NewBinding(key.WithKeys("enter")),
		clearFilter: key.NewBinding(key.WithKeys("esc")),
	}

	w.itemConstructor = itemConstructor

	w.InfiniteScroll = false
	w.AutoFocusNewItem = false
	w.filterable = true
	w.blockCursorMovingCallback = false
	w.filterState = Unfiltered
	w.Filter = FuzzyFilter

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

	w.BaseRenderable.SetMinSize(orvyn.NewSize(10, 5))
	w.BaseRenderable.SetPreferredSize(orvyn.NewSize(20, 10))

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
				w.filter(w.tiFilter.Value())

				return nil

			case key.Matches(msg, w.keybinds.clearFilter):
				w.clearFilter()

				w.FocusFirst()

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

			case key.Matches(msg, w.focusManager.NextFocusKeybind):
				w.NextItem()

			case key.Matches(msg, w.keybinds.enterFilter):
				if w.filterable {
					w.enterFilter()

					return nil
				}

			case key.Matches(msg, w.keybinds.clearFilter):
				if w.filterState == FilterApplied {
					w.clearFilter()

					w.FocusFirst()

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
	w.tiFilter.Resize(w.GetContentSize())

	w.paginatorUpdate()
}

func (w *Widget[T]) paginatorUpdate() {
	var perPage int

	contentSize := w.GetContentSize()
	total := 0
	calcHeight := 0

	for _, li := range w.listItems {
		li.Resize(contentSize)

		height := li.GetSize().Height

		w.maxItemHeight = max(w.maxItemHeight, height)

		total += height
	}

	calcHeight = contentSize.Height

	if total > contentSize.Height {
		calcHeight -= 1 // paginator
	}

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

			item := w.listItems[li.Index]
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

	contentSize := w.GetContentSize()

	return w.GetStyle().
		Width(contentSize.Width).
		Height(contentSize.Height).
		Render(view)
}

func (w *Widget[T]) OnFocus() {
	w.BaseFocusable.OnFocus()
	widget.UpdatePaginatorTheme(&w.paginator)
}

func (w *Widget[T]) OnBlur() {
	w.BaseFocusable.OnBlur()
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

	w.callCursorMovingCallback(w.globalIndex)

	if w.filterState == FilterApplied {
		w.previousFilteredItem()
		return
	}

	w.globalIndex--

	if w.globalIndex < 0 {
		if w.InfiniteScroll {
			w.globalIndex = len(w.listItems) - 1
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

	w.callCursorMovingCallback(w.globalIndex)

	if w.filterState == FilterApplied {
		w.nextFilteredItem()
		return
	}

	w.globalIndex++

	if w.globalIndex > len(w.listItems)-1 {
		if w.InfiniteScroll {
			w.globalIndex = 0
			w.moveCursor(0)
			return
		}

		w.globalIndex = len(w.listItems) - 1
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

func (w *Widget[T]) moveCursor(globalIndex int) {
	var cursor int

	// based on the global index set the cursor and the current page.
	if globalIndex < 0 {
		return
	}

	itemsOnPage := w.paginator.PerPage
	index := globalIndex

	if w.FilterState() == FilterApplied {
		for i, fi := range w.filteredListItems {
			if fi.Index == globalIndex {
				cursor = i
				index = i
				break
			}
		}
	} else {
		cursor = globalIndex % itemsOnPage
	}

	page := int(math.Floor(float64(index) / float64(itemsOnPage)))

	w.paginator.Page = page
	w.cursor = cursor

	if w.CursorMovedCallback != nil {
		w.CursorMovedCallback(w.globalIndex)
	}
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

// SetItems takes a []T (slice of data) and instantiate all items
// based on it.
func (w *Widget[T]) SetItems(items []T) {
	w.listItems = make([]ListItem[T], 0)
	focusableList := make([]orvyn.Focusable, 0)

	for i := range items {
		item := w.itemConstructor(items[i])
		w.listItems = append(w.listItems,
			item)
		focusableList = append(focusableList,
			item)
	}

	w.focusManager.SetWidgets(focusableList)

	w.focusManager.Focus(w.globalIndex)

	w.paginatorUpdate()
}

func (w *Widget[T]) SetCursorMovementKeybinds(cursorUp, cursorDown key.Binding) {
	w.keybinds.cursorUp = cursorUp
	w.keybinds.cursorDown = cursorDown
	w.focusManager.PreviousFocusKeybind = cursorUp
	w.focusManager.NextFocusKeybind = cursorDown
}

func (w *Widget[T]) GetItems() []T {
	var data []T

	for _, li := range w.listItems {
		data = append(data, li.GetData())
	}

	return data
}

func (w *Widget[T]) GetSelectedItem() T {
	var none T

	if w.globalIndex < 0 || w.globalIndex >= len(w.listItems) {
		return none
	}

	return w.listItems[w.globalIndex].GetData()
}

func (w *Widget[T]) GetItem(index int) T {
	var none T

	if index < 0 || index >= len(w.listItems) {
		return none
	}

	return w.listItems[index].GetData()
}

func (w *Widget[T]) SetItem(index int, data T) {
	if index < 0 || index >= len(w.listItems) {
		return
	}

	w.listItems[index].UpdateData(data)

	if w.filterState == FilterApplied {
		w.filter(w.tiFilter.Value())
	}
}

func (w *Widget[T]) AppendItem(data T) {
	w.clearFilter()

	index := len(w.listItems)

	widget := w.itemConstructor(data)

	w.listItems = append(w.listItems, widget)
	w.focusManager.Add(widget)

	if w.filterState == FilterApplied {
		w.filter(w.tiFilter.Value())
	}

	w.paginatorUpdate()

	if w.AutoFocusNewItem {
		w.callCursorMovingCallback(w.globalIndex)
		w.globalIndex = index
		w.moveCursor(w.globalIndex)
		w.focusManager.Focus(w.globalIndex)
	}
}

func (w *Widget[T]) InsertItem(index int, data T) {
	w.clearFilter()

	length := len(w.listItems)

	if length == 0 || index >= length {
		w.AppendItem(data)
		return
	}

	widget := w.itemConstructor(data)

	w.listItems = append(w.listItems[:index+1], w.listItems[index:]...)
	w.listItems[index] = widget
	w.focusManager.Insert(index, widget)

	if w.filterState == FilterApplied {
		w.filter(w.tiFilter.Value())
	}

	w.paginatorUpdate()

	if w.AutoFocusNewItem {
		w.callCursorMovingCallback(w.globalIndex)
		w.globalIndex = index
		w.moveCursor(w.globalIndex)
		w.focusManager.Focus(w.globalIndex)
	} else {
		if index <= w.globalIndex {
			w.NextItem()
		}
	}
}

func (w *Widget[T]) MoveItem(startIndex, destIndex int) {
	w.clearFilter()

	if startIndex < 0 || startIndex >= len(w.listItems) {
		return
	}

	if destIndex < 0 || destIndex > len(w.listItems) {
		return
	}

	item := w.listItems[startIndex]
	w.removeItem(startIndex)

	autoFocus := w.AutoFocusNewItem

	w.AutoFocusNewItem = true
	w.blockCursorMovingCallback = true

	w.InsertItem(destIndex, item.GetData())

	w.AutoFocusNewItem = autoFocus
	w.blockCursorMovingCallback = false
}

func (w *Widget[T]) RemoveItem(index int) {
	if index < 0 || index >= len(w.listItems) {
		return
	}

	w.removeItem(index)

	if w.filterState == FilterApplied {
		w.filter(w.tiFilter.Value())
	}

	w.paginatorUpdate()

	w.blockCursorMovingCallback = true

	w.PreviousItem()

	w.blockCursorMovingCallback = false

	w.focusManager.Focus(w.globalIndex)
}

func (w *Widget[T]) removeItem(index int) {
	if index < 0 || index >= len(w.listItems) {
		return
	}

	w.listItems = append(w.listItems[:index], w.listItems[index+1:]...)
	w.focusManager.Remove(index)
}

func (w *Widget[T]) FocusFirst() {
	w.focusManager.FocusFirst()

	if w.filterState == FilterApplied {
		if len(w.filteredListItems) > 0 {
			w.globalIndex = w.filteredListItems[0].Index
			w.cursor = 0
		} else {
			w.cursor = -1
			w.globalIndex = -1
		}

		w.paginator.Page = 0
	} else {
		w.globalIndex = 0
	}

	w.moveCursor(w.globalIndex)
}

func (w *Widget[T]) BlurCurrent() {
	w.focusManager.BlurCurrent()
}

func (w *Widget[T]) FilterState() FilterState {
	return w.filterState
}

func (w *Widget[T]) filter(s string) {
	if s == "" {
		w.clearFilter()
	}

	w.tiFilter.OnBlur()

	w.filteredListItems = w.Filter(&w.listItems, s)

	w.filterState = FilterApplied

	w.paginatorUpdate()

	w.FocusFirst()
}

func BasicFilter[T any](items *[]ListItem[T], s string) FilteredItems {
	var filteredItems FilteredItems

	for i, v := range *items {
		if strings.Contains(strings.ToLower(v.FilterValue()), strings.ToLower(s)) {
			filteredItems = append(filteredItems, FilteredItem{
				Index: i,
			})
		}
	}

	return filteredItems
}

func FuzzyFilter[T any](items *[]ListItem[T], s string) FilteredItems {
	var data []string
	var filteredItems FilteredItems

	for _, v := range *items {
		data = append(data, v.FilterValue())
	}

	matches := fuzzy.Find(s, data)

	for _, m := range matches {
		filteredItems = append(filteredItems, FilteredItem{
			Index: m.Index,
		})
	}

	return filteredItems
}

// Length returns the count of items in the list.
func (w *Widget[T]) Length() int {
	return len(w.listItems)
}

func (w *Widget[T]) clearFilter() {
	w.tiFilter.SetValue("")
	w.tiFilter.OnBlur()

	w.filteredListItems = make(FilteredItems, 0)

	for _, v := range w.listItems {
		v.SetActive(true)
	}

	w.filterState = Unfiltered

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
	return w.filteredListItems[index].Index
}

func (w *Widget[T]) callCursorMovingCallback(index int) {
	if w.blockCursorMovingCallback {
		return
	}

	if index < 0 || index >= len(w.listItems) {
		return
	}

	if w.CursorMovingCallback != nil {
		w.CursorMovingCallback(index)
	}
}
