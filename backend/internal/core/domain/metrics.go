package domain

import "time"

type ProcessStat struct {
	Process string
	Time time.Time
}

type UnitStat struct {
	Module string
	Processes []ProcessStat
	TotalTime time.Time
}

type Metrics struct {
	Stats []UnitStat
	DataQuantity string
	Time time.Time
}
