package generate

import (
	"fmt"
	"time"

	"charm.land/lipgloss/v2"

	"github.com/JulienQNN/comai/internal/config"
	"github.com/JulienQNN/comai/internal/theme"
)

func renderPreview(t theme.Theme, cfg config.Config, elapsed time.Duration) string {
	titleSep := t.Muted.Render(fmt.Sprintf(" · Generated in %s · %s/%s",
		elapsed.Truncate(10*time.Millisecond), cfg.ProviderName, cfg.ModelName))

	header := "Commit" + titleSep
	return header
}

func renderFooter(t theme.Theme, formattedDate, authorName, authorEmail string) string {
	return t.Italic.Render(fmt.Sprintf("%s <%s> \n%s", authorName, authorEmail, formattedDate))
}

func renderFinal(t theme.Theme, msg, authorEmail, authorName, date string) {
	final := lipgloss.JoinVertical(
		lipgloss.Left,
		t.CommitBorder.Render(msg),
		t.Italic.Render(fmt.Sprintf(" %s <%s>", authorName, authorEmail)),
		t.Italic.PaddingBottom(1).Render(fmt.Sprintf(" %s", date)),
	)
	fmt.Println(final)
}
