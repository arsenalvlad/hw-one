package serializers

import (
	"time"

	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/model"
	sqlModel "github.com/arsenalvlad/hw12_13_14_15_calendar/internal/storage/sql/model"
)

func ToModelEvent(data *sqlModel.CalendarEvent) *model.Event {
	return &model.Event{
		ID:          int(data.ID),
		Title:       data.Title,
		EventTime:   data.EventTime.Time,
		Duration:    time.Duration(data.Duration),
		Description: data.Description,
		UserID:      int(data.UserID),
	}
}

func ToModelEventSlice(data []*sqlModel.CalendarEvent) []*model.Event {
	items := make([]*model.Event, 0, len(data))

	for _, item := range data {
		items = append(items, ToModelEvent(item))
	}

	return items
}
