package main

import (
	"fmt"
	"io/fs"
	"path"
	"strconv"
	"time"
)

const dailyDir = "daily"

type FileSystemErr string

func (f FileSystemErr) Error() string {
	return string(f)
}

const NoPreviousDay FileSystemErr = "No previous day exists"

type FileSystem interface {
	fs.FS
	Mkdir(name string) error
}

func setLocation(fileSystem FileSystem) error {
	dirs, err := fs.ReadDir(fileSystem, ".")
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if dir.Name() == dailyDir && dir.IsDir() {
			return nil
		}
	}

	return fileSystem.Mkdir(dailyDir)
}

// getLastDay - finds the previous days file path
// if it does not exist returns error
func getLastDayPath(currentDate time.Time, fileSystem FileSystem) (string, error) {
	prevPath := getPreviousDayPath(currentDate)
	if _, err := fileSystem.Open(prevPath); err != nil {
		return "", NoPreviousDay
	}
	return prevPath, nil
}

func getPreviousDayPath(currentDate time.Time) string {
	prev := currentDate.AddDate(0, 0, -1)

	year := strconv.Itoa(prev.Year())
	month := prependZero(strconv.Itoa(int(prev.Month())))
	day := prependZero(strconv.Itoa(int(prev.Day())))
	shortYear := prev.Format("06")
	filename := fmt.Sprintf("%s_%s_%s.md", month, day, shortYear)

	return path.Join(year, month, filename)
}

func prependZero(s string) string {
	if len(s) == 2 {
		return s
	}
	return fmt.Sprintf("0%s", s)
}
