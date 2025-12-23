package progressbar

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/theme"
)

type Widget struct {
	orvyn.BaseWidget

	progress.Model

	TitleStyle lipgloss.Style

	title string

	MaxValue     int
	CurrentValue int

	showTitle                  bool
	showCurrentMaxValueInTitle bool
	showPercentage             bool
}

func New(title string) *Widget {
	w := new(Widget)

	t := orvyn.GetTheme()

	w.BaseWidget = orvyn.NewBaseWidget()

	w.Model = progress.New(progress.WithSolidFill(string(t.Color(theme.NormalFontColorID))))
	w.Model.ShowPercentage = false
	w.showPercentage = false

	w.TitleStyle = t.Style(theme.TitleStyleID).
		AlignHorizontal(lipgloss.Center)
	w.title = title

	w.showTitle = true
	w.showCurrentMaxValueInTitle = true

	return w
}

func (w *Widget) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case progress.FrameMsg:
		progressModel, cmd := w.Model.Update(msg)
		w.Model = progressModel.(progress.Model)
		return cmd
	}

	return nil
}

func (w *Widget) Render() string {
	var b strings.Builder

	if w.showTitle {
		switch {
		case w.showCurrentMaxValueInTitle && len(w.title) > 0:
			b.WriteString(w.TitleStyle.Render(
				fmt.Sprintf("%s (%d/%d)",
					w.title, w.CurrentValue, w.MaxValue)))
		case w.showCurrentMaxValueInTitle && len(w.title) == 0:
			b.WriteString(w.TitleStyle.Render(fmt.Sprintf("(%d/%d)",
				w.CurrentValue, w.MaxValue)))
		case !w.showCurrentMaxValueInTitle:
			b.WriteString(w.TitleStyle.Render(w.title))
		}

	}

	fmt.Fprintf(&b, "\n%s", w.Model.View())

	return b.String()
}

func (w *Widget) Resize(size orvyn.Size) {
	w.BaseWidget.Resize(size)

	w.Model.Width = size.Width

	w.TitleStyle = w.TitleStyle.Width(size.Width)
}

func (w *Widget) GetMinSize() orvyn.Size {
	titleHeight := orvyn.GetRenderSize(w.TitleStyle, w.title).Height
	return orvyn.NewSize(10, titleHeight+1)
}

func (w *Widget) GetPreferredSize() orvyn.Size {
	titleHeight := orvyn.GetRenderSize(w.TitleStyle, w.title).Height
	return orvyn.NewSize(30, titleHeight+1)
}

func (w *Widget) SetColor(color lipgloss.Color) {
	w.Model.FullColor = string(color)
}

func (w *Widget) SetTitleVisibility(b bool) {
	w.showTitle = b
}

func (w *Widget) SetTitleProgressVisibility(b bool) {
	w.showCurrentMaxValueInTitle = b
}

func (w *Widget) SetPercentageVisibility(b bool) {
	w.showPercentage = b
	w.Model.ShowPercentage = b
}
