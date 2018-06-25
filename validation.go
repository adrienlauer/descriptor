package model

import (
	"encoding/json"
	"log"
)

type ErrorType int

const (
	Warning ErrorType = 0
	Error   ErrorType = 1
)

type ValidationErrors struct {
	Errors []ValidationError
}

type ValidationError struct {
	ErrorType ErrorType
	Message   string
	Location  string
}

func (ve ValidationErrors) Error() string {
	return "Validation errors or warnings have occurred"
}

// JSonContent returns the serialized content of all validations
// errors as JSON
func (ve ValidationErrors) JSonContent() (b []byte, e error) {
	b, e = json.MarshalIndent(ve.Errors, "", "    ")
	return
}

func (t ErrorType) String() string {
	names := [...]string{
		"Warning",
		"Error"}
	if t < Warning || t > Error {
		return "Unknown"
	} else {
		return names[t]
	}
}

// HasErrors returns true if the ValidationErrors contains at least on error
func (ve ValidationErrors) HasErrors() bool {
	for _, v := range ve.Errors {
		if v.ErrorType == Error {
			return true
		}
	}
	return false
}

// HasWarnings returns true if the ValidationErrors contains at least on warning
func (ve ValidationErrors) HasWarnings() bool {
	for _, v := range ve.Errors {
		if v.ErrorType == Warning {
			return true
		}
	}
	return false
}

func (ve *ValidationErrors) AddError(err error, loc string) {
	ve.Errors = append(ve.Errors, ValidationError{
		Location:  loc,
		Message:   err.Error(),
		ErrorType: Error})
}

func (ve *ValidationErrors) AddWarning(message string, loc string) {
	ve.Errors = append(ve.Errors, ValidationError{
		Location:  loc,
		Message:   message,
		ErrorType: Warning})
}

func (ve ValidationErrors) Log(logger *log.Logger) {
	for _, err := range ve.Errors {
		logger.Println(err.ErrorType.String() + ": " + err.Message + " @" + err.Location)
	}
}
