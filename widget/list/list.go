package list

import (
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

type IListItem[T any] interface {
	orvyn.Focusable
	orvyn.Renderable
	GetData() T
}

// ItemConstructor defines the signature of the item constructor.
// T type represents the type of the item data.
type ItemConstructor[T any] func(T) IListItem[T]

// Widget defines a list widget.
// T type represents the type of the item data.
type Widget[T any] struct {
	orvyn.BaseWidget
	orvyn.BaseFocusable

	InfiniteScroll bool

	cursor      int
	globalIndex int

	listItems []IListItem[T]
	items     []T

	paginator paginator.Model

	focusManager *orvyn.FocusManager

	itemConstructor ItemConstructor[T]

	style lipgloss.Style

	contentSize orvyn.Size

	CursorMovedCallback func(int)
}

// New creates a new *Widget list and takes an itemConstructor as parameter.
// T type represents the type of the item data.
func New[T any](itemConstructor ItemConstructor[T]) *Widget[T] {
	w := new(Widget[T])

	w.BaseWidget = orvyn.NewBaseWidget()

	w.itemConstructor = itemConstructor

	w.InfiniteScroll = false

	w.cursor = 0

	w.paginator = paginator.New()
	w.paginator.Type = paginator.Dots

	w.focusManager = orvyn.NewFocusManager()
	w.focusManager.ManageFocusNextPrevKeybind = false
	w.focusManager.PreviousFocusKeybind = key.NewBinding(key.WithKeys("up"))
	w.focusManager.NextFocusKeybind = key.NewBinding(key.WithKeys("down"))
	w.focusManager.Focus(0)

	w.OnBlur()

	return w
}

func (w *Widget[T]) Update(msg tea.Msg) tea.Cmd {
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

			}
		}
	}

	cmd := w.focusManager.Update(msg)

	w.focusManager.Focus(w.globalIndex)

	return cmd
}

func (w *Widget[T]) Resize(size orvyn.Size) {
	var perPage int

	maxItemHeight := 1

	w.BaseWidget.Resize(size)

	size.Width -= w.style.GetHorizontalFrameSize()
	size.Height -= w.style.GetVerticalFrameSize()

	for _, li := range w.listItems {
		li.Resize(size)

		maxItemHeight = max(maxItemHeight, li.GetSize().Height)
	}

	perPage = (size.Height - 1) / maxItemHeight
	perPage = max(perPage, 1)

	w.paginator.PerPage = perPage
	w.paginator.SetTotalPages(len(w.listItems))

	w.contentSize = size
}

func (w *Widget[T]) Render() string {
	var b strings.Builder
	var view string

	count := 0
	start, end := w.paginator.GetSliceBounds(len(w.listItems))

	for i, li := range w.listItems[start:end] {
		if i > 0 {
			b.WriteString("\n")
		}

		b.WriteString(li.Render())

		count++
	}

	if w.paginator.TotalPages > 1 {
		view = lipgloss.JoinVertical(
			lipgloss.Center, b.String(),
			w.paginator.View())
	} else {
		view = lipgloss.JoinVertical(
			lipgloss.Center, b.String(),
		)
	}

	return w.style.
		Width(w.contentSize.Width).
		Height(w.contentSize.Height).
		Render(view)
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

// NextItem manages the focus of the next item.
func (w *Widget[T]) NextItem() {
	if len(w.listItems) == 0 {
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

func (w *Widget[T]) moveCursor(index int) {
	// based on the global index set the cursor and the current page.

	itemsOnPage := w.paginator.PerPage

	page := int(math.Floor(float64(index) / float64(itemsOnPage)))
	cursor := index % itemsOnPage

	w.paginator.Page = page
	w.cursor = cursor
}

func (w *Widget[T]) GetGlobalIndex() int {
	return w.globalIndex
}

// SetItems takes a []T (slice of data) and instantiate all items
// based on it.
func (w *Widget[T]) SetItems(items []T) {
	w.items = items

	w.listItems = make([]IListItem[T], 0)
	focusableList := make([]orvyn.Focusable, 0)

	for _, i := range w.items {
		item := w.itemConstructor(i)
		w.listItems = append(w.listItems,
			item)
		focusableList = append(focusableList,
			item)
	}

	w.focusManager.SetWidgets(focusableList)
}

func (w *Widget[T]) GetItems() []T {
	var items []T

	items = make([]T, 0)

	for _, item := range w.listItems {
		items = append(items, item.GetData())
	}

	return items
}

func (w *Widget[T]) FocusItem(index int) {
	w.focusManager.Focus(index)
	w.globalIndex = index
	w.moveCursor(index)
}

func (w *Widget[T]) BlurCurrent() {
	w.focusManager.BlurCurrent()
}
