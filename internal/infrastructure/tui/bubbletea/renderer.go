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
	boGhiNhatKy *utils.Logger
	ketQua      *[]model.KetQua
	khoaKetQua  *sync.Mutex
	kenhKetQua  chan model.KetQua
	kenhDong    chan struct{}
	teaProgram  *tea.Program
}

// RendererModel holds the state for the BubbleTea task renderer
type RendererModel struct {
	renderer    *BubbleTeaRenderer
	status      *components.StatusComponent
	selectedRow int
}

// KetQuaMsg represents a task result message
type KetQuaMsg struct {
	KetQua model.KetQua
}

func NewRendererModel(renderer *BubbleTeaRenderer) RendererModel {
	return RendererModel{
		renderer:    renderer,
		status:      components.NewStatusComponent(renderer.boGhiNhatKy),
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
			m.renderer.khoaKetQua.Lock()
			if m.selectedRow > 0 {
				m.selectedRow--
			}
			m.renderer.khoaKetQua.Unlock()
			m.ensureSelectedRowVisible()
			m.status.Viewport.SetContent(components.RenderTaskTable(m.renderer.ketQua, m.selectedRow))
		case "down":
			m.renderer.khoaKetQua.Lock()
			if m.selectedRow < len(*m.renderer.ketQua)-1 {
				m.selectedRow++
			}
			m.renderer.khoaKetQua.Unlock()
			m.ensureSelectedRowVisible()
			m.status.Viewport.SetContent(components.RenderTaskTable(m.renderer.ketQua, m.selectedRow))
		case "pgup":
			m.status.Viewport.PageUp()
		case "pgdown":
			m.status.Viewport.PageDown()
		}
	case tea.WindowSizeMsg:
		m.status.UpdateViewport(msg.Width, msg.Height, components.RenderTaskTable(m.renderer.ketQua, m.selectedRow))
	case KetQuaMsg:
		m.renderer.khoaKetQua.Lock()
		*m.renderer.ketQua = append(*m.renderer.ketQua, msg.KetQua)
		m.renderer.khoaKetQua.Unlock()
		if m.status.Ready {
			tableContent := components.RenderTaskTable(m.renderer.ketQua, m.selectedRow)
			m.renderer.boGhiNhatKy.Infof("Bubbletea updated table: %s", tableContent)
			m.status.Viewport.SetContent(tableContent)
		}
	}
	return m, nil
}

func (m RendererModel) View() string {
	return m.status.View()
}

// NewBubbleTeaRenderer creates a new BubbleTeaRenderer
func NewBubbleTeaRenderer(boGhiNhatKy *utils.Logger, ketQua *[]model.KetQua, khoaKetQua *sync.Mutex, kenhKetQua chan model.KetQua, kenhDong chan struct{}) *BubbleTeaRenderer {
	return &BubbleTeaRenderer{
		boGhiNhatKy: boGhiNhatKy,
		ketQua:      ketQua,
		khoaKetQua:  khoaKetQua,
		kenhKetQua:  kenhKetQua,
		kenhDong:    kenhDong,
	}
}

func (r *BubbleTeaRenderer) Start() error {
	r.teaProgram = tea.NewProgram(NewRendererModel(r), tea.WithOutput(os.Stdout))
	go func() {
		for {
			select {
			case ketQua := <-r.kenhKetQua:
				r.boGhiNhatKy.Infof("Bubbletea sending task result: %v", ketQua)
				r.teaProgram.Send(KetQuaMsg{KetQua: ketQua})
			case <-r.kenhDong:
				r.boGhiNhatKy.Info("Bubbletea closing task renderer")
				r.teaProgram.Quit()
				return
			}
		}
	}()
	if _, err := r.teaProgram.Run(); err != nil {
		r.boGhiNhatKy.Errorf("Failed to start bubbletea: %v", err)
		return fmt.Errorf("could not start bubbletea: %w", err)
	}
	return nil
}

func (r *BubbleTeaRenderer) AddTaskResult(ketQua model.KetQua) {
	select {
	case r.kenhKetQua <- ketQua:
		r.boGhiNhatKy.Infof("Added task result to channel: %v", ketQua)
	case <-r.kenhDong:
		r.boGhiNhatKy.Info("Task channel closed, ignoring result")
	}
}

func (r *BubbleTeaRenderer) AddProxyResult(ketQua model.KetQuaProxy) {
	// Not used in this renderer
}

func (r *BubbleTeaRenderer) Close() {
	// Closing is handled in usecase
}
