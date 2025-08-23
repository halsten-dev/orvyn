package widget

import (
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/halsten-dev/orvyn"
	"github.com/halsten-dev/orvyn/theme"
)

func UpdatePaginatorTheme(paginator *paginator.Model) {
	paginator.ActiveDot = orvyn.GetTheme().
		Style(theme.PaginatorActiveStyleName).Render("•")
	paginator.InactiveDot = orvyn.GetTheme().
		Style(theme.PaginatorInactiveStyleName).Render("•")
}
