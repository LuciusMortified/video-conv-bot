package ent

import "io"

type ConvertParams struct {
	URL string
}

type ConvertStatus string

const (
	ConvertDownloading ConvertStatus = "downloading"
	ConvertConverting  ConvertStatus = "converting"
	ConvertDone        ConvertStatus = "done"
	ConvertError       ConvertStatus = "error"
	ConvertUnsupported ConvertStatus = "unsupported"
)

func (s ConvertStatus) Valid() bool {
	switch s {
	case
		ConvertDownloading,
		ConvertConverting,
		ConvertDone,
		ConvertError,
		ConvertUnsupported:
		return true
	default:
		return false
	}
}

type ConvertResult struct {
	Filename string
	Data     io.ReadCloser
}

type ConvertState struct {
	Status ConvertStatus
	Result *ConvertResult
	Error  *string
}

func (s *ConvertState) Cleanup() error {
	if s.Result != nil && s.Result.Data != nil {
		return s.Result.Data.Close()
	}
	return nil
}

type ConvertStateChan chan ConvertState
