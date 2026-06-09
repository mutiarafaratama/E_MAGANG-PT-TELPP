package service

import (
	"os"
	"time"
)

func parseDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func mkdirAll(path string) error {
	return os.MkdirAll(path, 0755)
}
