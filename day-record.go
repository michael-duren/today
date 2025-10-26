package main

import "time"

type DayRecordErr string

func (d DayRecordErr) Error() string {
	return string(d)
}

const MissingDate DayRecordErr = "Unable to read file date contents"

type DayRecord struct {
	date     time.Time
	previous []string
	current  []string
}

// NewDayRecord creates a DayRecord from
// a []byte typically read from a file
func NewDayRecord(filecontents []byte) (*DayRecord, error) {
	return &DayRecord{}, nil
}
