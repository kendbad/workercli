package tview

import (
	"sync"
	"workercli/internal/domain/model"
	"workercli/internal/infrastructure/tui/tview/components"
	"workercli/pkg/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TViewProxyRenderer for ProxyCheck
type TViewProxyRenderer struct {
	logger     *utils.Logger
	results    *[]model.ProxyResult
	resultsMu  *sync.Mutex
	resultChan chan model.ProxyResult
	closeChan  chan struct{}
	tviewTable *tview.Table
	tviewApp   *tview.Application
}

// NewTViewProxyRenderer creates a new TViewProxyRenderer
func NewTViewProxyRenderer(logger *utils.Logger, results *[]model.ProxyResult, resultsMu *sync.Mutex, resultChan chan model.ProxyResult, closeChan chan struct{}) *TViewProxyRenderer {
	return &TViewProxyRenderer{
		logger:     logger,
		results:    results,
		resultsMu:  resultsMu,
		resultChan: resultChan,
		closeChan:  closeChan,
	}
}

func (r *TViewProxyRenderer) Start() error {
	r.tviewApp = tview.NewApplication()
	r.tviewTable = tview.NewTable().SetBorders(true).SetFixed(1, 0).SetSelectable(true, false)

	r.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			r.tviewApp.Stop()
		}
		return event
	})

	go func() {
		row := 1
		for {
			select {
			case result := <-r.resultChan:
				r.resultsMu.Lock()
				*r.results = append(*r.results, result)
				components.RenderProxyTable(r.tviewTable, r.results, row)
				row++
				r.resultsMu.Unlock()
				r.tviewApp.QueueUpdateDraw(func() {})
			case <-r.closeChan:
				r.tviewApp.QueueUpdateDraw(func() {})
				return
			}
		}
	}()

	if err := r.tviewApp.SetRoot(r.tviewTable, true).Run(); err != nil {
		r.logger.Errorf("Failed to start tview: %v", err)
		return err
	}
	return nil
}

func (r *TViewProxyRenderer) AddTaskResult(result model.Result) {
	// Not used in this renderer
}

func (r *TViewProxyRenderer) AddProxyResult(result model.ProxyResult) {
	select {
	case r.resultChan <- result:
		r.logger.Infof("Added proxy result to channel: %v", result)
	case <-r.closeChan:
		r.logger.Info("Proxy channel closed, ignoring proxy result")
	}
}

func (r *TViewProxyRenderer) Close() {
	// Closing is handled in usecase
}
