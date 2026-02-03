package models

import "errors"

// Common errors for walkthrough validation
var (
	ErrMissingTitle  = errors.New("walkthrough must have a title")
	ErrNoSteps       = errors.New("walkthrough must have at least one step")
	ErrInvalidStepID = errors.New("step IDs must be sequential starting from 1")
	ErrMissingFile   = errors.New("step must specify a file")
)
