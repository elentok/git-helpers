package worktrees

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// placeOverlay composites fg on top of bg at position (x, y), where x and y
// are zero-based visual column and row offsets. Both bg and fg may contain
// arbitrary ANSI escape sequences; ansi.Truncate and ansi.TruncateLeft are
// used to splice each fg line into the corresponding bg line in a way that
// fully preserves the background's colors and styles outside the modal area.
func placeOverlay(bg, fg string, x, y int) string {
	bgLines := strings.Split(bg, "\n")
	fgLines := strings.Split(fg, "\n")

	for i, fgLine := range fgLines {
		bgY := y + i
		if bgY < 0 || bgY >= len(bgLines) {
			continue
		}
		bgLine := bgLines[bgY]
		fgW := ansi.StringWidth(fgLine)

		// Left portion of the background (columns 0..x-1).
		left := ansi.Truncate(bgLine, x, "")
		if leftW := ansi.StringWidth(left); leftW < x {
			left += strings.Repeat(" ", x-leftW)
		}

		// Right portion of the background (columns x+fgW onwards).
		right := ansi.TruncateLeft(bgLine, x+fgW, "")

		bgLines[bgY] = left + fgLine + right
	}
	return strings.Join(bgLines, "\n")
}
