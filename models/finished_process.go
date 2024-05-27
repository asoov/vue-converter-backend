package models

import "time"

type FinishedProcess struct {
	Timestamp  time.Time
	BucketName string
	FileCount  int
	SignedUrls []string
	Name       string
}
