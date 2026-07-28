package converter

import (
	"database/sql"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var testTime = time.Date(2026, time.July, 28, 12, 30, 45, 0, time.UTC)

func TestTimeToProto(t *testing.T) {
	got := TimeToProto(testTime)

	if got == nil {
		t.Fatal("TimeToProto() = nil, want timestamp")
	}

	if !got.AsTime().Equal(testTime) {
		t.Errorf("TimeToProto() = %v, want %v", got.AsTime(), testTime)
	}
}

func TestNullTimeToProto(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullTime
		wantNil  bool
		wantTime time.Time
	}{
		{
			name:     "valid",
			input:    sql.NullTime{Time: testTime, Valid: true},
			wantNil:  false,
			wantTime: testTime,
		},
		{
			name:    "null",
			input:   sql.NullTime{},
			wantNil: true,
		},
		{
			// Valid = false важнее непустого Time: колонка всё равно NULL.
			name:    "null with time set",
			input:   sql.NullTime{Time: testTime, Valid: false},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NullTimeToProto(tt.input)

			if tt.wantNil {
				if got != nil {
					t.Errorf("NullTimeToProto() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("NullTimeToProto() = nil, want timestamp")
			}

			if !got.AsTime().Equal(tt.wantTime) {
				t.Errorf("NullTimeToProto() = %v, want %v", got.AsTime(), tt.wantTime)
			}
		})
	}
}

func TestProtoToNullTime(t *testing.T) {
	tests := []struct {
		name  string
		input *timestamppb.Timestamp
		want  sql.NullTime
	}{
		{
			name:  "valid",
			input: timestamppb.New(testTime),
			want:  sql.NullTime{Time: testTime, Valid: true},
		},
		{
			name:  "nil",
			input: nil,
			want:  sql.NullTime{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProtoToNullTime(tt.input)

			if got.Valid != tt.want.Valid {
				t.Errorf("ProtoToNullTime().Valid = %v, want %v", got.Valid, tt.want.Valid)
			}

			if !got.Time.Equal(tt.want.Time) {
				t.Errorf("ProtoToNullTime().Time = %v, want %v", got.Time, tt.want.Time)
			}
		})
	}
}

func TestNullTimeRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input sql.NullTime
	}{
		{name: "valid", input: sql.NullTime{Time: testTime, Valid: true}},
		{name: "null", input: sql.NullTime{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProtoToNullTime(NullTimeToProto(tt.input))

			if got.Valid != tt.input.Valid {
				t.Errorf("round trip Valid = %v, want %v", got.Valid, tt.input.Valid)
			}

			if !got.Time.Equal(tt.input.Time) {
				t.Errorf("round trip Time = %v, want %v", got.Time, tt.input.Time)
			}
		})
	}
}

func TestNullStringToProto(t *testing.T) {
	tests := []struct {
		name    string
		input   sql.NullString
		wantNil bool
		want    string
	}{
		{
			name:  "valid",
			input: sql.NullString{String: "dc.host", Valid: true},
			want:  "dc.host",
		},
		{
			// Пустая строка и NULL — разные значения, и после конвертации
			// должны остаться разными.
			name:  "valid empty",
			input: sql.NullString{String: "", Valid: true},
			want:  "",
		},
		{
			name:    "null",
			input:   sql.NullString{},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NullStringToProto(tt.input)

			if tt.wantNil {
				if got != nil {
					t.Errorf("NullStringToProto() = %q, want nil", *got)
				}
				return
			}

			if got == nil {
				t.Fatal("NullStringToProto() = nil, want value")
			}

			if *got != tt.want {
				t.Errorf("NullStringToProto() = %q, want %q", *got, tt.want)
			}
		})
	}
}

func TestProtoToNullString(t *testing.T) {
	value := "dc.host"
	empty := ""

	tests := []struct {
		name  string
		input *string
		want  sql.NullString
	}{
		{name: "valid", input: &value, want: sql.NullString{String: value, Valid: true}},
		{name: "valid empty", input: &empty, want: sql.NullString{String: "", Valid: true}},
		{name: "nil", input: nil, want: sql.NullString{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProtoToNullString(tt.input)

			if got != tt.want {
				t.Errorf("ProtoToNullString() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
