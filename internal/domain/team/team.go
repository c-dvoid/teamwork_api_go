package domain_team

import "time"

type Team struct {
	ID        int64
	Name      string
	CreatedBy int64
	CreatedAt time.Time
}
