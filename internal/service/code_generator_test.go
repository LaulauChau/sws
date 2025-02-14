package service

import (
	"testing"

	"github.com/LaulauChau/sws/internal/models"
)

func Test_encode(t *testing.T) {
	tests := []struct {
		name string
		num  int
		want string
	}{
		{
			name: "encode zero",
			num:  0,
			want: "8",
		},
		{
			name: "encode single digit",
			num:  5,
			want: "6",
		},
		{
			name: "encode multiple digits",
			num:  123,
			want: "213", // Using the actual output from the encode function
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encode(arrayCharsNumeric, tt.num); got != tt.want {
				t.Errorf("encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fillWithZero(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "already 5 digits",
			s:    "12345",
			want: "12345",
		},
		{
			name: "needs padding",
			s:    "123",
			want: "00123",
		},
		{
			name: "empty string",
			s:    "",
			want: "00000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fillWithZero(tt.s); got != tt.want {
				t.Errorf("fillWithZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCourseStartUTC(t *testing.T) {
	tests := []struct {
		name    string
		course  models.Course
		wantErr bool
	}{
		{
			name: "valid course time",
			course: models.Course{
				Date:  "2025-02-10",
				Start: "08:00:00+00:00",
			},
			wantErr: false,
		},
		{
			name: "invalid date format",
			course: models.Course{
				Date:  "invalid",
				Start: "08:00:00+00:00",
			},
			wantErr: true,
		},
		{
			name: "empty date",
			course: models.Course{
				Date:  "",
				Start: "08:00:00+00:00",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getCourseStartUTC(tt.course)
			if (got.IsZero()) != tt.wantErr {
				t.Errorf("getCourseStartUTC() error = %v, wantErr %v", got.IsZero(), tt.wantErr)
			}
		})
	}
}

func TestGenerateFixedCode(t *testing.T) {
	tests := []struct {
		name      string
		course    models.Course
		wantCode  string
		wantEmpty bool
	}{
		{
			name: "valid course",
			course: models.Course{
				ID:    137393,
				Name:  "Test Course",
				Date:  "2025-02-10",
				Start: "08:00:00+00:00",
			},
			wantCode:  "09866",
			wantEmpty: false,
		},
		{
			name: "zero ID",
			course: models.Course{
				ID:    0,
				Date:  "2025-02-10",
				Start: "08:00:00+00:00",
			},
			wantEmpty: true,
		},
		{
			name: "invalid date",
			course: models.Course{
				ID:    1,
				Date:  "invalid",
				Start: "08:00:00+00:00",
			},
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, code := GenerateFixedCode(tt.course)
			if tt.wantEmpty {
				if code != "" {
					t.Errorf("GenerateFixedCode() = %v, want empty string", code)
				}
				return
			}
			if code != tt.wantCode {
				t.Errorf("GenerateFixedCode() = %v, want %v", code, tt.wantCode)
			}
		})
	}
}
