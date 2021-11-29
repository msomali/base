package base

import (
	"io"
)

type Params struct {
	DebugMode bool
	Logger    io.Writer
}

type OptionFunc func(params *Params)

func DebugModeOption(mode bool) OptionFunc {
	return func(params *Params) {
		params.DebugMode = mode
	}
}

func LoggerOption(writer io.Writer) OptionFunc {
	return func(params *Params) {
		params.Logger = writer
	}
}
