package anime_schedule

import "time"

type AnimeSchedule struct {
	ID                  string     `gorm:"column:id;primaryKey" json:"id"`
	AnimeID             string     `gorm:"column:anime_id" json:"anime_id"`
	AnimescheduleRoute  *string    `gorm:"column:animeschedule_route" json:"animeschedule_route"`
	JpnTime             *time.Time `gorm:"column:jpn_time" json:"jpn_time"`
	SubTime             *time.Time `gorm:"column:sub_time" json:"sub_time"`
	DubTime             *time.Time `gorm:"column:dub_time" json:"dub_time"`
	Notes               *string    `gorm:"column:notes" json:"notes"`
	DelayedTimetable    *string    `gorm:"column:delayed_timetable" json:"delayed_timetable"`
	SubDelayedTimetable *string    `gorm:"column:sub_delayed_timetable" json:"sub_delayed_timetable"`
	DubDelayedTimetable *string    `gorm:"column:dub_delayed_timetable" json:"dub_delayed_timetable"`
	LastSyncedAt        time.Time  `gorm:"column:last_synced_at" json:"last_synced_at"`
	CreatedAt           time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (AnimeSchedule) TableName() string {
	return "anime_schedule"
}
