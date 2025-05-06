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
	boGhiNhatKy *utils.Logger
	ketQua      *[]model.KetQuaTrungGian
	khoaKetQua  *sync.Mutex
	kenhKetQua  chan model.KetQuaTrungGian
	kenhDong    chan struct{}
	teaProgram  *tea.Program
}

// ProxyRendererModel holds the state for the BubbleTea proxy renderer
type ProxyRendererModel struct {
	renderer    *BubbleTeaProxyRenderer
	status      *components.StatusComponent
	selectedRow int
}

// KetQuaTrungGianMsg represents a proxy result message
type KetQuaTrungGianMsg struct {
	KetQua model.KetQuaTrungGian
}

func NewProxyRendererModel(renderer *BubbleTeaProxyRenderer) ProxyRendererModel {
	return ProxyRendererModel{
		renderer:    renderer,
		status:      components.NewStatusComponent(renderer.boGhiNhatKy),
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
			m.renderer.khoaKetQua.Lock()
			if m.selectedRow > 0 {
				m.selectedRow--
			}
			m.renderer.khoaKetQua.Unlock()
			m.ensureSelectedRowVisible()
			m.status.Viewport.SetContent(components.RenderProxyTable(m.renderer.ketQua, m.selectedRow))
		case "down":
			m.renderer.khoaKetQua.Lock()
			if m.selectedRow < len(*m.renderer.ketQua)-1 {
				m.selectedRow++
			}
			m.renderer.khoaKetQua.Unlock()
			m.ensureSelectedRowVisible()
			m.status.Viewport.SetContent(components.RenderProxyTable(m.renderer.ketQua, m.selectedRow))
		case "pgup":
			m.status.Viewport.PageUp()
		case "pgdown":
			m.status.Viewport.PageDown()
		}
	case tea.WindowSizeMsg:
		m.status.UpdateViewport(msg.Width, msg.Height, components.RenderProxyTable(m.renderer.ketQua, m.selectedRow))
	case KetQuaTrungGianMsg:
		m.renderer.khoaKetQua.Lock()
		*m.renderer.ketQua = append(*m.renderer.ketQua, msg.KetQua)
		m.renderer.khoaKetQua.Unlock()
		if m.status.Ready {
			tableContent := components.RenderProxyTable(m.renderer.ketQua, m.selectedRow)
			m.renderer.boGhiNhatKy.Infof("Bubbletea updated table: %s", tableContent)
			m.status.Viewport.SetContent(tableContent)
		}
	}
	return m, nil
}

func (m ProxyRendererModel) View() string {
	return m.status.View()
}

// NewBubbleTeaProxyRenderer creates a new BubbleTeaProxyRenderer
func NewBubbleTeaProxyRenderer(boGhiNhatKy *utils.Logger, ketQua *[]model.KetQuaTrungGian, khoaKetQua *sync.Mutex, kenhKetQua chan model.KetQuaTrungGian, kenhDong chan struct{}) *BubbleTeaProxyRenderer {
	return &BubbleTeaProxyRenderer{
		boGhiNhatKy: boGhiNhatKy,
		ketQua:      ketQua,
		khoaKetQua:  khoaKetQua,
		kenhKetQua:  kenhKetQua,
		kenhDong:    kenhDong,
	}
}

func (r *BubbleTeaProxyRenderer) Start() error {
	r.teaProgram = tea.NewProgram(NewProxyRendererModel(r), tea.WithOutput(os.Stdout))
	go func() {
		for {
			select {
			case ketQua := <-r.kenhKetQua:
				r.boGhiNhatKy.Infof("Bubbletea sending proxy result: %v", ketQua)
				r.teaProgram.Send(KetQuaTrungGianMsg{KetQua: ketQua})
			case <-r.kenhDong:
				r.boGhiNhatKy.Info("Bubbletea closing proxy renderer")
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

func (r *BubbleTeaProxyRenderer) AddTaskResult(ketQua model.KetQua) {
	// Not used in this renderer
}

func (r *BubbleTeaProxyRenderer) AddProxyResult(ketQua model.KetQuaTrungGian) {
	select {
	case r.kenhKetQua <- ketQua:
		r.boGhiNhatKy.Infof("Added proxy result to channel: %v", ketQua)
	case <-r.kenhDong:
		r.boGhiNhatKy.Info("Proxy channel closed, ignoring proxy result")
	}
}

func (r *BubbleTeaProxyRenderer) Close() {
	// Closing is handled in usecase
}
