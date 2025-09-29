// Package todayfile
package todayfile

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type TodayFile interface {
	CreateFile() (*os.File, error)
	PrintInfo()
	UpdateContents(*os.File) error
}

type dateInfo struct {
	year  string
	month string
	day   string
}

func prependZero(s string) string {
	if len(s) == 2 {
		return s
	}
	return fmt.Sprintf("0%s", s)
}

func (d *dateInfo) PrintInfo() {
	fmt.Printf("creating file with the following content")
	fmt.Printf("year: %s, month: %s, day: %s\n", d.year, d.month, d.day)
}

func NewTodayFile() TodayFile {
	t := time.Now()
	return &dateInfo{
		year:  strconv.Itoa(t.Year()),
		month: prependZero(strconv.Itoa(int(t.Month()))),
		day:   prependZero(strconv.Itoa(t.Day())),
	}
}

func newYesterdayFile() *dateInfo {
	yesterday := time.Now().AddDate(0, 0, -1)
	return &dateInfo{
		year:  strconv.Itoa(yesterday.Year()),
		month: prependZero(strconv.Itoa(int(yesterday.Month()))),
		day:   prependZero(strconv.Itoa(yesterday.Day())),
	}
}

func (d *dateInfo) getBasePath() string {
	return path.Join(d.year, d.month)
}

func (d *dateInfo) getShortYear() string {
	return strconv.Itoa(time.Now().Year() % 100)
}

func (d *dateInfo) getFilename() string {
	return fmt.Sprintf("%s_%s_%s.md", d.month, d.day, d.getShortYear())
}

func (d *dateInfo) getPath() string {
	return path.Join(d.getBasePath(), d.getFilename())
}

func (d *dateInfo) basePathExists() bool {
	_, err := os.Stat(d.getBasePath())
	return err == nil
}

func (d *dateInfo) createFileDirectory() error {
	if d.basePathExists() {
		return nil
	}

	return os.MkdirAll(d.getBasePath(), 0o755)
}

func (d *dateInfo) getContents(b []byte) ([]string, error) {
	return strings.Split(string(b), "\n"), nil
}

func (d *dateInfo) openFile() ([]byte, error) {
	return os.ReadFile(d.getPath())
}

func (d *dateInfo) openPreviousFile() ([]byte, error) {
	contents, err := os.ReadFile(d.getPath())
	if err == nil {
		return contents, nil
	}

	entries, err := os.ReadDir(d.getBasePath())
	if err != nil {
		return nil, err
	}
	var mostRecentFilename os.FileInfo
	for _, e := range entries {
		if e.Type().IsDir() {
			continue
		}

		entriePath := path.Join(d.getBasePath(), e.Name())
		fi, err := os.Stat(entriePath)
		if err != nil {
			panic("error loading file when it should exist")
		}

		if mostRecentFilename == nil {
			mostRecentFilename = fi
			continue
		}

		if fi.ModTime().After(mostRecentFilename.ModTime()) {
			mostRecentFilename = fi
		}
	}
	return os.ReadFile(path.Join(d.getBasePath(), mostRecentFilename.Name()))
}

func (d *dateInfo) CreateFile() (*os.File, error) {
	d.createFileDirectory()
	return os.OpenFile(d.getPath(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
}

func (d *dateInfo) readYesterdayFile() ([]string, error) {
	yesterdayFile, err := d.openPreviousFile()
	if err != nil {
		return nil, err
	}

	contents, err := d.getContents(yesterdayFile)
	if err != nil {
		return nil, err
	}

	if len(contents) == 0 {
		return nil, errors.New("file does not exist")
	}

	return contents, nil
}

func (d *dateInfo) getTodayHeader() string {
	return fmt.Sprintf("# %s_%s_%s\n", d.month, d.day, d.year)
}

func (d *dateInfo) UpdateContents(f *os.File) error {
	yesterdayContents, err := d.readYesterdayFile()
	if err != nil {
		_, _ = f.Write([]byte(d.getTodayHeader()))
		return nil
	}

	yesterdayContents[0] = d.getTodayHeader()
	_, err = f.Write([]byte(strings.Join(yesterdayContents, "\n")))
	return err
}
