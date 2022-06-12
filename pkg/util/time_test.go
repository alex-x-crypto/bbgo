package util

import (
	"reflect"
	"testing"
	"time"
)

func TestParseTimeWithFormats(t *testing.T) {
	var looseTimeFormats = []string{
		time.RFC3339,
		time.RFC822,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02",
	}

	type args struct {
		strTime string
		formats []string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name:    "date",
			args:    args{strTime: "2022-01-04", formats: looseTimeFormats},
			want:    time.Date(2022, time.January, 4, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "date",
			args:    args{strTime: "2022-01-04T11:31", formats: looseTimeFormats},
			want:    time.Date(2022, time.January, 4, 11, 31, 0, 0, time.UTC),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTimeWithFormats(tt.args.strTime, tt.args.formats)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimeWithFormats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTimeWithFormats() got = %v, want %v", got, tt.want)
			}
		})
	}
}
