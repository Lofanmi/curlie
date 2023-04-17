package formatter

import (
	"io"
	"os"
	"testing"
)

func TestJSON(t *testing.T) {
	var stdout io.Writer = os.Stdout
	scheme := ColorScheme{}
	stdout = &JSON{
		Out:              stdout,
		Scheme:           scheme,
		ParseJsonUnicode: true,
	}
	var data string
	data = `{"a":1}`
	_, _ = stdout.Write([]byte(data))
	data = `{"a":"\u4f60\u597d"}`
	_, _ = stdout.Write([]byte(data))
	data = `{"code":0,"msg":"\u6210\u529f","data":[{"id":1065,"name":"\u6d4b\u8bd5","status":3}]}`
	_, _ = stdout.Write([]byte(data))
}

func TestJSON_Write(t *testing.T) {
	j1 := &JSON{Out: os.Stdout, Scheme: ColorScheme{}, ParseJsonUnicode: true}
	j2 := &JSON{Out: os.Stdout, Scheme: ColorScheme{}, ParseJsonUnicode: false}
	t1, t2, t3 := []byte(`{"a":1}`), []byte(`{"a":"\u4f60\u597d"}`), []byte(`{"code":0,"msg":"\u6210\u529f","data":[{"id":1065,"name":"\u6d4b\u8bd5","status":3}]}`)
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		j       *JSON
		args    args
		wantErr bool
	}{
		{"j1:t1", j1, args{p: t1}, false},
		{"j1:t2", j1, args{p: t2}, false},
		{"j1:t3", j1, args{p: t3}, false},

		{"j2:t1", j2, args{p: t1}, false},
		{"j2:t2", j2, args{p: t2}, false},
		{"j2:t3", j2, args{p: t3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.j.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
