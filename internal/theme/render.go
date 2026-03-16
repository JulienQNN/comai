package theme

import (
	"fmt"
	"path/filepath"
)

func (t Theme) RenderFileChange(statusChar, fullPath string) string {
	var badge string

	switch statusChar {
	case "M":
		badge = t.BadgeMod.Render("~ MOD")
	case "A", "?":
		badge = t.BadgeAdd.Render("+ ADD")
	case "D":
		badge = t.BadgeDel.Render("- DEL")
	case "R":
		badge = t.BadgeMod.Render("» REN")
	case "C":
		badge = t.BadgeAdd.Render("= CPY")
	default:
		badge = t.BadgeUnk.Render("? UNK")
	}

	dir := filepath.Dir(fullPath)
	file := filepath.Base(fullPath)

	var styledPath string
	if dir == "." {
		styledPath = t.FileName.Render(file)
	} else {
		styledPath = t.FileDir.Render(dir+"/") + t.FileName.Render(file)
	}

	return fmt.Sprintf("%s %s", badge, styledPath)
}
