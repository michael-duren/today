package main

import (
	"path/filepath"
	"testing"
	"testing/fstest"
	"time"
)

type MockFileSystem struct {
	fstest.MapFS
	createdDirs []string
}

func (m *MockFileSystem) Mkdir(name string) error {
	m.createdDirs = append(m.createdDirs, name)
	return nil
}

func NewMockFileSystem(files map[string]*fstest.MapFile) *MockFileSystem {
	return &MockFileSystem{
		MapFS:       files,
		createdDirs: []string{},
	}
}

func TestSetLocation(t *testing.T) {
	t.Run("detects when in correct directory", func(t *testing.T) {
		mockFS := NewMockFileSystem(fstest.MapFS{
			"daily/note.txt": {Data: []byte("test")},
		})

		err := setLocation(mockFS)

		if err != nil {
			t.Fatalf("expected err to be nil, got %v", err)
		}
		if len(mockFS.createdDirs) != 0 {
			t.Fatalf("expected no directories to be created, got %v", mockFS.createdDirs)
		}
	})

	t.Run("creates daily directory when not present", func(t *testing.T) {
		mockFS := NewMockFileSystem(fstest.MapFS{
			"other/file.txt": {Data: []byte("test")},
		})

		err := setLocation(mockFS)

		if err != nil {
			t.Fatalf("expected err to be nil, got %v", err)
		}
		if len(mockFS.createdDirs) != 1 {
			t.Fatalf("expected 1 directory to be created, got %v", len(mockFS.createdDirs))
		}
		if mockFS.createdDirs[0] != "daily" {
			t.Errorf("expected 'daily' to be created, got %q", mockFS.createdDirs[0])
		}
	})
}

func TestGetLastDayPath(t *testing.T) {
	t.Run("returns correct path when previous day is yesterday", func(t *testing.T) {
		previousDayPath := filepath.Join("2025", "01", "01_01_25.md")
		mockFS := NewMockFileSystem(fstest.MapFS{
			previousDayPath: {Data: []byte("test")},
			filepath.Join("2025", "01", "01_02_25.md"): {Data: []byte("test")},
			filepath.Join("2025", "01", "01_03_25.md"): {Data: []byte("test")},
		})
		currentDate := time.Date(2025, 01, 02, 12, 30, 0, 0, time.Local)
		result, err := getLastDayPath(currentDate, mockFS)
		if result != previousDayPath {
			t.Fatalf("expected %v, got %v", previousDayPath, result)
		}

		if err != nil {
			t.Fatalf("err was not nil, %v", err)
		}

	})
	t.Run("returns correct path when previous day is last month", func(t *testing.T) {
		lastMonth := filepath.Join("2025", "01", "01_31_25.md")
		mockFS := NewMockFileSystem(fstest.MapFS{
			lastMonth: {Data: []byte("test")},
		})
		currentDate := time.Date(2025, 02, 01, 12, 30, 0, 0, time.Local)
		result, err := getLastDayPath(currentDate, mockFS)

		AssertNoError(t, err)
		AssertEqual(t, lastMonth, result)
	})
}
