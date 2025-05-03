package tview

import (
	"sync"
	"workercli/internal/domain/model"
	"workercli/internal/infrastructure/tui/tview/components"
	"workercli/pkg/utils"

	"github.com/rivo/tview"
)

// TViewRenderer for BatchTask
type TViewRenderer struct {
	logger     *utils.Logger
	results    *[]model.Result
	resultsMu  *sync.Mutex
	resultChan chan model.Result
	closeChan  chan struct{}
	tviewTable *tview.Table
	tviewApp   *tview.Application
}

// NewTViewRenderer creates a new TViewRenderer
func NewTViewRenderer(logger *utils.Logger, results *[]model.Result, resultsMu *sync.Mutex, resultChan chan model.Result, closeChan chan struct{}) *TViewRenderer {
	return &TViewRenderer{
		logger:     logger,
		results:    results,
		resultsMu:  resultsMu,
		resultChan: resultChan,
		closeChan:  closeChan,
	}
}

func (r *TViewRenderer) Start() error {
	r.tviewApp = tview.NewApplication()
	r.tviewTable = tview.NewTable().SetBorders(true).SetFixed(1, 0).SetSelectable(true, false)

	go func() {
		row := 1
		for {
			select {
			case result := <-r.resultChan:
				r.resultsMu.Lock()
				*r.results = append(*r.results, result)
				components.RenderTaskTable(r.tviewTable, r.results, row)
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

func (r *TViewRenderer) AddTaskResult(result model.Result) {
	select {
	case r.resultChan <- result:
		r.logger.Infof("Added task result to channel: %v", result)
	case <-r.closeChan:
		r.logger.Info("Task channel closed, ignoring result")
	}
}

func (r *TViewRenderer) AddProxyResult(result model.ProxyResult) {
	// Not used in this renderer
}

func (r *TViewRenderer) Close() {
	// Closing is handled in usecase
}
