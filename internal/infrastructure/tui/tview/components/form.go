package components

import "github.com/rivo/tview"

// FormComponent will handle user input forms in the future
type FormComponent struct {
	form *tview.Form
}

// NewFormComponent creates a new form component
func NewFormComponent() *FormComponent {
	return &FormComponent{
		form: tview.NewForm(),
	}
}

// RenderForm renders the form (placeholder for future implementation)
func (f *FormComponent) RenderForm() *tview.Form {
	// Add form fields as needed
	return f.form
}
