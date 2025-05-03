package tui

// TUIMode defines the type of TUI interface
type TUIMode string

const (
	TUIModeNone      TUIMode = ""
	TUIModeTView     TUIMode = "tview"
	TUIModeBubbleTea TUIMode = "bubbletea"
)
