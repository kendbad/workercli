package bubbletea

import (
	"fmt"
	"os"
	"sync"
	"workercli/internal/domain/model"
	"workercli/internal/infrastructure/tui/bubbletea/components"
	"workercli/pkg/utils"

	tea "github.com/charmbracelet/bubbletea"
)

// BubbleTeaRenderer for BatchTask
type BubbleTeaRenderer struct {
	logger     *utils.Logger
	results    *[]model.Result
	resultsMu  *sync.Mutex
	resultChan chan model.Result
	closeChan  chan struct{}
	teaProgram *tea.Program
}

// RendererModel holds the state for the BubbleTea task renderer
type RendererModel struct {
	renderer    *BubbleTeaRenderer
	status      *components.StatusComponent
	selectedRow int
}

// ResultMsg represents a task result message
type ResultMsg struct {
	Result model.Result
}

func NewRendererModel(renderer *BubbleTeaRenderer) RendererModel {
	return RendererModel{
		renderer:    renderer,
		status:      components.NewStatusComponent(renderer.logger),
		selectedRow: 0,
	}
}

func (m RendererModel) Init() tea.Cmd {
	return nil
}

func (m RendererModel) ensureSelectedRowVisible() {
	visibleHeight := m.status.Viewport.Height
	if m.selectedRow < m.status.Viewport.YOffset {
		m.status.Viewport.SetYOffset(m.selectedRow)
	} else if m.selectedRow >= m.status.Viewport.YOffset+visibleHeight {
		m.status.Viewport.SetYOffset(m.selectedRow - visibleHeight + 1)
	}
}

func (m RendererModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
			m.renderer.resultsMu.Lock()
			if m.selectedRow > 0 {
				m.selectedRow--
			}
			m.renderer.resultsMu.Unlock()
			m.ensureSelectedRowVisible()
			m.status.Viewport.SetContent(components.RenderTaskTable(m.renderer.results, m.selectedRow))
		case "down":
			m.renderer.resultsMu.Lock()
			if m.selectedRow < len(*m.renderer.results)-1 {
				m.selectedRow++
			}
			m.renderer.resultsMu.Unlock()
			m.ensureSelectedRowVisible()
			m.status.Viewport.SetContent(components.RenderTaskTable(m.renderer.results, m.selectedRow))
		case "pgup":
			m.status.Viewport.PageUp()
		case "pgdown":
			m.status.Viewport.PageDown()
		}
	case tea.WindowSizeMsg:
		m.status.UpdateViewport(msg.Width, msg.Height, components.RenderTaskTable(m.renderer.results, m.selectedRow))
	case ResultMsg:
		m.renderer.resultsMu.Lock()
		*m.renderer.results = append(*m.renderer.results, msg.Result)
		m.renderer.resultsMu.Unlock()
		if m.status.Ready {
			tableContent := components.RenderTaskTable(m.renderer.results, m.selectedRow)
			m.renderer.logger.Infof("Bubbletea updated table: %s", tableContent)
			m.status.Viewport.SetContent(tableContent)
		}
	}
	return m, nil
}

func (m RendererModel) View() string {
	return m.status.View()
}

// NewBubbleTeaRenderer creates a new BubbleTeaRenderer
func NewBubbleTeaRenderer(logger *utils.Logger, results *[]model.Result, resultsMu *sync.Mutex, resultChan chan model.Result, closeChan chan struct{}) *BubbleTeaRenderer {
	return &BubbleTeaRenderer{
		logger:     logger,
		results:    results,
		resultsMu:  resultsMu,
		resultChan: resultChan,
		closeChan:  closeChan,
	}
}

func (r *BubbleTeaRenderer) Start() error {
	r.teaProgram = tea.NewProgram(NewRendererModel(r), tea.WithOutput(os.Stdout))
	go func() {
		for {
			select {
			case result := <-r.resultChan:
				r.logger.Infof("Bubbletea sending task result: %v", result)
				r.teaProgram.Send(ResultMsg{Result: result})
			case <-r.closeChan:
				r.logger.Info("Bubbletea closing task renderer")
				r.teaProgram.Quit()
				return
			}
		}
	}()
	if _, err := r.teaProgram.Run(); err != nil {
		r.logger.Errorf("Failed to start bubbletea: %v", err)
		return fmt.Errorf("could not start bubbletea: %w", err)
	}
	return nil
}

func (r *BubbleTeaRenderer) AddTaskResult(result model.Result) {
	select {
	case r.resultChan <- result:
		r.logger.Infof("Added task result to channel: %v", result)
	case <-r.closeChan:
		r.logger.Info("Task channel closed, ignoring result")
	}
}

func (r *BubbleTeaRenderer) AddProxyResult(result model.ProxyResult) {
	// Not used in this renderer
}

func (r *BubbleTeaRenderer) Close() {
	// Closing is handled in usecase
}
