package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCode(t *testing.T) {
	type args struct {
		err  error
		code int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Is code",
			args: args{
				err:  Codef(unknownCoder.Code(), "internal error message"),
				code: unknownCoder.Code(),
			},
			want: true,
		},
		{
			name: "Is code cause",
			args: args{
				err:  WithCodef(Codef(unknownCoder.Code(), "test"), errEOF, "test"),
				code: unknownCoder.Code(),
			},
			want: true,
		},
		{
			name: "Error code is not the target",
			args: args{
				err:  Codef(errConfigurationNotValid, "Configuration Not Valid"),
				code: unknownCoder.Code(),
			},
			want: false,
		},
		{
			name: "Not code err",
			args: args{
				err:  errors.New("test"),
				code: unknownCoder.Code(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, IsCode(tt.args.err, tt.args.code), "IsCode(%v, %v)", tt.args.err, tt.args.code)
		})
	}
}

func TestMustRegister(t *testing.T) {
	type args struct {
		coder Coder
	}
	MustRegister(defaultCoder{1000, 500, "ConfigurationNotValid error"})
	tests := []struct {
		name string
		args args
	}{
		{
			name: "MustRegister",
			args: args{
				coder: defaultCoder{-1, 500, "ConfigurationNotValid error"},
			},
		},
		{
			name: "MustRegister Panic coder == 1",
			args: args{
				coder: defaultCoder{1, 500, "ConfigurationNotValid error"},
			},
		},
		{
			name: "MustRegister Panic coder repeat",
			args: args{
				coder: defaultCoder{1000, 500, "ConfigurationNotValid error"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					_, ok := codes[tt.args.coder.Code()]
					assert.True(t, (0 <= tt.args.coder.Code() && tt.args.coder.Code() <= 100) || ok)
				}
			}()
			MustRegister(tt.args.coder)
		})
	}
}

func TestParseCoder(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		wantHTTPCode int
		wantString   string
		wantCode     int
		wantNil      bool
	}{
		{
			name:         "fmt err",
			err:          fmt.Errorf("yes error"),
			wantHTTPCode: 500,
			wantString:   "内部服务器错误",
			wantCode:     1,
			wantNil:      false,
		},
		{
			name:         "Codef",
			err:          Codef(unknownCoder.Code(), "internal error message"),
			wantHTTPCode: 500,
			wantString:   "内部服务器错误",
			wantCode:     1,
			wantNil:      false,
		},
		{
			name:    "wantNil",
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coder := ParseCoder(tt.err)
			if tt.wantNil {
				assert.Nil(t, coder)
				return
			}
			assert.Equalf(t, coder.HTTPStatus(), tt.wantHTTPCode, "TestCodeParse: got %q, want: %q", coder.HTTPStatus(), tt.wantHTTPCode)
			assert.Equalf(t, coder.String(), tt.wantString, "TestCodeParse: got %q, want: %q", coder.String(), tt.wantString)
			assert.Equalf(t, coder.Code(), tt.wantCode, "TestCodeParse: got %q, want: %q", coder.Code(), tt.wantCode)
		})
	}
}

func TestRegister(t *testing.T) {
	type args struct {
		coder Coder
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Register",
			args: args{
				coder: defaultCoder{-1, 500, "ConfigurationNotValid error"},
			},
		},
		{
			name: "Register",
			args: args{
				coder: defaultCoder{1, 500, "ConfigurationNotValid error"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					assert.True(t, 0 <= tt.args.coder.Code() && tt.args.coder.Code() <= 100)
				}
			}()
			Register(tt.args.coder)
		})
	}
}

func Test_defaultCoder_Code(t *testing.T) {
	type fields struct {
		C    int
		HTTP int
		Ext  string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "defaultCoder_Code",
			fields: fields{
				C: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coder := defaultCoder{
				C:    tt.fields.C,
				HTTP: tt.fields.HTTP,
				Ext:  tt.fields.Ext,
			}
			assert.Equalf(t, tt.want, coder.Code(), "Code()")
		})
	}
}

func Test_defaultCoder_HTTPStatus(t *testing.T) {
	type fields struct {
		C    int
		HTTP int
		Ext  string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "defaultCoder_HTTPStatus",
			fields: fields{
				HTTP: 500,
			},
			want: 500,
		},
		{
			name: "defaultCoder_HTTPStatus",
			fields: fields{
				HTTP: 0,
			},
			want: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coder := defaultCoder{
				C:    tt.fields.C,
				HTTP: tt.fields.HTTP,
				Ext:  tt.fields.Ext,
			}
			assert.Equalf(t, tt.want, coder.HTTPStatus(), "HTTPStatus()")
		})
	}
}

func Test_defaultCoder_String(t *testing.T) {
	type fields struct {
		C    int
		HTTP int
		Ext  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "defaultCoder_String",
			fields: fields{
				Ext: "defaultCoder_String",
			},
			want: "defaultCoder_String",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coder := defaultCoder{
				C:    tt.fields.C,
				HTTP: tt.fields.HTTP,
				Ext:  tt.fields.Ext,
			}
			assert.Equalf(t, tt.want, coder.String(), "String()")
		})
	}
}
