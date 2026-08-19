package domain_task

import "time"

type FieldChange struct {
	Old any `json:"old"`
	New any `json:"new"`
}

type Changes map[string]FieldChange

type History struct {
	ID        int64
	TaskID    int64
	ChangedBy int64
	Changes   Changes
	CreatedAt time.Time
}
