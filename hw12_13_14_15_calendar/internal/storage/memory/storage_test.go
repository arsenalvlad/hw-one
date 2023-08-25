package memorystorage

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/app"
	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/model"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want app.Storage
	}{
		{
			name: "Test create new memory storage",
			want: &Storage{
				data: nil,
				mu:   &sync.RWMutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_AddEvent(t *testing.T) {
	type fields struct {
		data map[int]model.Event
		mu   *sync.RWMutex
	}
	type args struct {
		in0  context.Context
		data model.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Event
		wantErr bool
	}{
		{
			name: "test memory storage add event",
			fields: fields{
				data: nil,
				mu:   &sync.RWMutex{},
			},
			args: args{
				in0: context.Background(),
				data: model.Event{
					ID:          1,
					Title:       "qwe",
					EventTime:   time.Now(),
					Duration:    3000,
					Description: "asdq",
					UserID:      1,
				},
			},
			want: &model.Event{
				ID:          1,
				Title:       "qwe",
				EventTime:   time.Now(),
				Duration:    3000,
				Description: "asdq",
				UserID:      1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				data: tt.fields.data,
				mu:   tt.fields.mu,
			}
			err := s.Connect()
			require.NoError(t, err, "AddEvent() error = %v", err)

			got, err := s.AddEvent(tt.args.in0, tt.args.data)
			require.NoError(t, err, "AddEvent() error = %v", err)
			require.Equal(t, got, tt.want, "AddEvent() got = %v, want %v", got, tt.want)
		})
	}
}

func TestStorage_Close(t *testing.T) {
	type fields struct {
		data map[int]model.Event
		mu   *sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test memory storage close",
			fields: fields{
				data: nil,
				mu:   &sync.RWMutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				data: tt.fields.data,
				mu:   tt.fields.mu,
			}
			err := s.Close()
			require.NoError(t, err, "Close() error = %v", err)
		})
	}
}

func TestStorage_Connect(t *testing.T) {
	type fields struct {
		data map[int]model.Event
		mu   *sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test memory storage connect",
			fields: fields{
				data: nil,
				mu:   &sync.RWMutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				data: tt.fields.data,
				mu:   tt.fields.mu,
			}
			err := s.Connect()
			require.NoError(t, err, "Connect() error = %v", err)
			require.Equal(t, 0, len(s.data), "Connect() error = %v", err)
		})
	}
}
