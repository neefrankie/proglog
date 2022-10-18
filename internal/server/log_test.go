package server

import (
	"reflect"
	"testing"
)

func TestLog_Append(t *testing.T) {
	c := NewLog()

	type args struct {
		record Record
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name: "append hello",
			args: args{
				record: Record{
					Value: []byte("hello"),
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "append world",
			args: args{
				record: Record{
					Value: []byte("world"),
				},
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Append(tt.args.record)
			if (err != nil) != tt.wantErr {
				t.Errorf("Append() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Append() got = %v, want %v", got, tt.want)
			}

			t.Logf("%v", c.records)
		})
	}
}

func TestLog_Read(t *testing.T) {
	c := NewLog()
	c.Append(Record{Value: []byte("hello")})
	c.Append(Record{Value: []byte("world")})

	type args struct {
		offset uint64
	}
	tests := []struct {
		name    string
		args    args
		want    Record
		wantErr bool
	}{
		{
			name: "Read hello",
			args: args{
				offset: 0,
			},
			want: Record{
				Value:  []byte("hello"),
				Offset: 0,
			},
			wantErr: false,
		},
		{
			name: "Read world",
			args: args{
				offset: 1,
			},
			want: Record{
				Value:  []byte("world"),
				Offset: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := c.Read(tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}
