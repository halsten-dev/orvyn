package layout

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/halsten-dev/orvyn"
)

// CenterLayout centers the given element in the available size.
type CenterLayout struct {
	orvyn.BaseLayout
}

func NewCenterLayout(element orvyn.Renderable) *CenterLayout {
	l := new(CenterLayout)

	l.BaseLayout = orvyn.NewBaseLayout(element)

	return l
}

func (l *CenterLayout) Render() string {
	if len(l.GetElements()) == 0 {
		return ""
	}

	size := l.GetSize()

	l.GetElements()[0].Resize(size)

	return lipgloss.Place(
		size.Width, size.Height,
		lipgloss.Center, lipgloss.Center,
		l.GetElements()[0].Render(),
	)
}

func (l *CenterLayout) GetMinSize() orvyn.Size {
	return l.GetElements()[0].GetMinSize()
}

func (l *CenterLayout) GetPreferredSize() orvyn.Size {
	return l.GetElements()[0].GetPreferredSize()
}
