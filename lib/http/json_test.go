package lib

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeJSON(t *testing.T) {
	type args struct {
		b io.Reader
		v any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				b: io.NopCloser(bytes.NewBuffer([]byte(`{"name":"test"}`))),
				v: &struct {
					Name string `json:"name"`
				}{},
			},
			wantErr: false,
		},
		{
			name: "Error",
			args: args{
				b: io.NopCloser(bytes.NewBuffer([]byte(`{"name":"test"`))),
				v: &struct {
					Name string `json:"name"`
				}{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DecodeJSON(tt.args.b, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("DecodeJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRespJSON(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		status int
		data   any
	}
	tests := []struct {
		name     string
		args     args
		wantBody string
	}{
		{
			name: "Success",
			args: args{
				w:      httptest.NewRecorder(),
				status: http.StatusOK,
				data: map[string]string{
					"name": "test",
				},
			},
			wantBody: `{"name":"test"}`,
		},
		{
			name: "InvalidData",
			args: args{
				w:      httptest.NewRecorder(),
				status: http.StatusNotFound,
				data:   nil,
			},
			wantBody: `null`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RespJSON(tt.args.w, tt.args.status, tt.args.data)
			got := tt.args.w.(*httptest.ResponseRecorder).Body.String()
			got = strings.TrimSuffix(got, "\n") // Trim the newline character
			if got != tt.wantBody {
				t.Errorf("RespJSON() = %v, want %v", got, tt.wantBody)
			}
		})
	}
}
