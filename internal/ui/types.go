package ui

import (
	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textinput"
)

type WizardResult struct {
	ModelName      string
	Language       string
	MaxLength      string
	IncludeBody    bool
	PromptAddition string
}

type model struct {
	quitting   bool
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	isGlobal   bool
}
