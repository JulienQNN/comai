package generate

import (
	"fmt"
	"time"

	"charm.land/lipgloss/v2"

	"github.com/JulienQNN/comai/internal/config"
	"github.com/JulienQNN/comai/internal/git"
	"github.com/JulienQNN/comai/internal/theme"
)

func renderPreview(
	t theme.Theme,
	cfg config.Config,
	msg string,
	author git.AuthorInfo,
	date string,
	elapsed time.Duration,
) {
	titleSep := t.Muted.Render(fmt.Sprintf(" · Generated in %s · %s/%s",
		elapsed.Truncate(10*time.Millisecond), cfg.ProviderName, cfg.ModelName))

	header := lipgloss.JoinHorizontal(lipgloss.Center, t.Title.Render("Commit"), titleSep)

	finalOutput := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		t.CommitBorder.Render(msg),
		t.Italic.Render(fmt.Sprintf(" %s <%s>", author.Name, author.Email)),
		t.Italic.PaddingBottom(1).Render(fmt.Sprintf(" %s", date)),
	)
	fmt.Println(finalOutput)
}
