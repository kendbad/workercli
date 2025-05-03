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

// BubbleTeaProxyRenderer for ProxyCheck
type BubbleTeaProxyRenderer struct {
	logger     *utils.Logger
	results    *[]model.ProxyResult
	resultsMu  *sync.Mutex
	resultChan chan model.ProxyResult
	closeChan  chan struct{}
	teaProgram *tea.Program
}

// ProxyRendererModel holds the state for the BubbleTea proxy renderer
type ProxyRendererModel struct {
	renderer    *BubbleTeaProxyRenderer
	status      *components.StatusComponent
	selectedRow int
}

// ProxyResultMsg represents a proxy result message
type ProxyResultMsg struct {
	Result model.ProxyResult
}

func NewProxyRendererModel(renderer *BubbleTeaProxyRenderer) ProxyRendererModel {
	return ProxyRendererModel{
		renderer:    renderer,
		status:      components.NewStatusComponent(renderer.logger),
		selectedRow: 0,
	}
}

func (m ProxyRendererModel) Init() tea.Cmd {
	return nil
}

func (m ProxyRendererModel) ensureSelectedRowVisible() {
	visibleHeight := m.status.Viewport.Height
	if m.selectedRow < m.status.Viewport.YOffset {
		m.status.Viewport.SetYOffset(m.selectedRow)
	} else if m.selectedRow >= m.status.Viewport.YOffset+visibleHeight {
		m.status.Viewport.SetYOffset(m.selectedRow - visibleHeight + 1)
	}
}

func (m ProxyRendererModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.status.Viewport.SetContent(components.RenderProxyTable(m.renderer.results, m.selectedRow))
		case "down":
			m.renderer.resultsMu.Lock()
			if m.selectedRow < len(*m.renderer.results)-1 {
				m.selectedRow++
			}
			m.renderer.resultsMu.Unlock()
			m.ensureSelectedRowVisible()
			m.status.Viewport.SetContent(components.RenderProxyTable(m.renderer.results, m.selectedRow))
		case "pgup":
			m.status.Viewport.PageUp()
		case "pgdown":
			m.status.Viewport.PageDown()
		}
	case tea.WindowSizeMsg:
		m.status.UpdateViewport(msg.Width, msg.Height, components.RenderProxyTable(m.renderer.results, m.selectedRow))
	case ProxyResultMsg:
		m.renderer.resultsMu.Lock()
		*m.renderer.results = append(*m.renderer.results, msg.Result)
		m.renderer.resultsMu.Unlock()
		if m.status.Ready {
			tableContent := components.RenderProxyTable(m.renderer.results, m.selectedRow)
			m.renderer.logger.Infof("Bubbletea updated table: %s", tableContent)
			m.status.Viewport.SetContent(tableContent)
		}
	}
	return m, nil
}

func (m ProxyRendererModel) View() string {
	return m.status.View()
}

// NewBubbleTeaProxyRenderer creates a new BubbleTeaProxyRenderer
func NewBubbleTeaProxyRenderer(logger *utils.Logger, results *[]model.ProxyResult, resultsMu *sync.Mutex, resultChan chan model.ProxyResult, closeChan chan struct{}) *BubbleTeaProxyRenderer {
	return &BubbleTeaProxyRenderer{
		logger:     logger,
		results:    results,
		resultsMu:  resultsMu,
		resultChan: resultChan,
		closeChan:  closeChan,
	}
}

func (r *BubbleTeaProxyRenderer) Start() error {
	r.teaProgram = tea.NewProgram(NewProxyRendererModel(r), tea.WithOutput(os.Stdout))
	go func() {
		for {
			select {
			case result := <-r.resultChan:
				r.logger.Infof("Bubbletea sending proxy result: %v", result)
				r.teaProgram.Send(ProxyResultMsg{Result: result})
			case <-r.closeChan:
				r.logger.Info("Bubbletea closing proxy renderer")
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

func (r *BubbleTeaProxyRenderer) AddTaskResult(result model.Result) {
	// Not used in this renderer
}

func (r *BubbleTeaProxyRenderer) AddProxyResult(result model.ProxyResult) {
	select {
	case r.resultChan <- result:
		r.logger.Infof("Added proxy result to channel: %v", result)
	case <-r.closeChan:
		r.logger.Info("Proxy channel closed, ignoring proxy result")
	}
}

func (r *BubbleTeaProxyRenderer) Close() {
	// Closing is handled in usecase
}
