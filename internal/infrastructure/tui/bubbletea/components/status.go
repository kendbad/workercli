// workercli/internal/infrastructure/tui/bubbletea/components/status.go
package components

import (
	"workercli/pkg/utils"

	"github.com/charmbracelet/bubbles/viewport"
)

// StatusComponent manages generic status rendering
type StatusComponent struct {
	logger   *utils.Logger
	Ready    bool
	Viewport viewport.Model
}

// NewStatusComponent creates a new status component
func NewStatusComponent(logger *utils.Logger) *StatusComponent {
	return &StatusComponent{
		logger: logger,
		Ready:  false,
	}
}

// UpdateViewport updates the viewport size
func (s *StatusComponent) UpdateViewport(width, height int, content string) {
	if !s.Ready {
		s.Viewport = viewport.New(width, height-2)
		s.Viewport.SetContent(content)
		s.Ready = true
	} else {
		s.Viewport.Width = width
		s.Viewport.Height = height - 2
	}
}

// View renders the status view
func (s *StatusComponent) View() string {
	if !s.Ready {
		s.logger.Info("Status component not ready")
		return "Loading..."
	}
	s.logger.Infof("Rendering status viewport: %s", s.Viewport.View())
	return s.Viewport.View() + "\nPress 'q' to quit."
}
