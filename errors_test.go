package errors

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		err  string
		want error
	}{
		{"", fmt.Errorf("")},
		{"foo", fmt.Errorf("foo")},
		{"foo", New("foo")},
		{"string with format specifiers: %v", errors.New("string with format specifiers: %v")},
	}

	for _, tt := range tests {
		got := New(tt.err)
		if got.Error() != tt.want.Error() {
			t.Errorf("New.Error(): got: %q, want %q", got, tt.want)
		}
	}
}

func TestWrapNil(t *testing.T) {
	got := Wrap(nil, "no error")
	if got != nil {
		t.Errorf("Wrap(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{Wrap(io.EOF, "read error"), "client error", "client error: read error: EOF"},
	}

	for _, tt := range tests {
		got := Wrap(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("Wrap(%v, %q): got: %v, want %v", tt.err, tt.message, got, tt.want)
		}
	}
}

type nilError struct{}

func (nilError) Error() string { return "nil error" }

func TestCause(t *testing.T) {
	topErr := Codef(errEOF, "topErr")
	x := New("error")
	tests := []struct {
		err  error
		want error
	}{{
		// nil error is nil
		err:  nil,
		want: nil,
	},
		{
			// explicit nil error is nil
			err:  (error)(nil),
			want: nil,
		},
		{
			// typed nil is nil
			err:  (*nilError)(nil),
			want: (*nilError)(nil),
		},
		{
			// uncaused error is unaffected
			err:  io.EOF,
			want: io.EOF,
		},
		{
			// caused error returns cause
			err:  Wrap(io.EOF, "ignored"),
			want: io.EOF,
		},
		{
			err:  x, // return from errors.New
			want: x,
		},
		{
			WithMessage(nil, "whoops"),
			nil,
		},
		{
			WithMessage(io.EOF, "whoops"),
			io.EOF,
		},
		{
			WithCodef(Wrap(topErr, "err2"), errConfigurationNotValid, "err1"),
			topErr,
		},
		{
			WithStack(nil),
			nil,
		},
		{
			WithStack(io.EOF),
			io.EOF,
		},
	}

	for i, tt := range tests {
		got := Cause(tt.err)
		assert.Equalf(t, got, tt.want, "test %d: got %#v, want %#v", i+1, got, tt.want)
	}
}

func TestWrapfNil(t *testing.T) {
	got := Wrapf(nil, "no error")
	if got != nil {
		t.Errorf("Wrapf(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWrapf(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{Wrapf(io.EOF, "read error without format specifiers"), "client error", "client error: read error without format specifiers: EOF"},
		{Wrapf(io.EOF, "read error with %d format specifier", 1), "client error", "client error: read error with 1 format specifier: EOF"},
		{Codef(errEOF, "EOF"), "Codef", "Codef"},
	}

	for _, tt := range tests {
		got := Wrapf(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("Wrapf(%v, %q): got: %v, want %v", tt.err, tt.message, got, tt.want)
		}
	}
}

func TestErrorf(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{Errorf("read error without format specifiers"), "read error without format specifiers"},
		{Errorf("read error with %d format specifier", 1), "read error with 1 format specifier"},
	}

	for _, tt := range tests {
		got := tt.err.Error()
		if got != tt.want {
			t.Errorf("Errorf(%v): got: %q, want %q", tt.err, got, tt.want)
		}
	}
}

func TestWithStackNil(t *testing.T) {
	got := WithStack(nil)
	if got != nil {
		t.Errorf("WithStack(nil): got %#v, expected nil", got)
	}
}

func TestWithStack(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{io.EOF, "EOF"},
		{WithStack(io.EOF), "EOF"},
		{Codef(errEOF, "EOF"), "EOF"},
	}

	for _, tt := range tests {
		got := WithStack(tt.err).Error()
		if got != tt.want {
			t.Errorf("WithStack(%v): got: %v, want %v", tt.err, got, tt.want)
		}
	}
}

func TestWithMessageNil(t *testing.T) {
	got := WithMessage(nil, "no error")
	if got != nil {
		t.Errorf("WithMessage(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWithMessage(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{WithMessage(io.EOF, "read error"), "client error", "client error: read error: EOF"},
	}

	for _, tt := range tests {
		got := WithMessage(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("WithMessage(%v, %q): got: %q, want %q", tt.err, tt.message, got, tt.want)
		}
	}
}

func TestWithMessagefNil(t *testing.T) {
	got := WithMessagef(nil, "no error")
	if got != nil {
		t.Errorf("WithMessage(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWithMessagef(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{WithMessagef(io.EOF, "read error without format specifier"), "client error", "client error: read error without format specifier: EOF"},
		{WithMessagef(io.EOF, "read error with %d format specifier", 1), "client error", "client error: read error with 1 format specifier: EOF"},
	}

	for _, tt := range tests {
		got := WithMessagef(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("WithMessage(%v, %q): got: %q, want %q", tt.err, tt.message, got, tt.want)
		}
	}
}

func TestCode(t *testing.T) {
	tests := []struct {
		code       int
		message    string
		wantType   string
		wantCode   int
		wangString string
	}{
		{errConfigurationNotValid, "Configuration error", "*withCode", errConfigurationNotValid, "Configuration error"},
	}

	for _, tt := range tests {
		got := Code(tt.code, tt.message)
		err, ok := got.(*withCode)
		if !ok {
			t.Errorf("Codef(%v, %q): error type got: %T, want %s", tt.code, tt.message, got, tt.wantType)
		}

		if err.code != tt.wantCode {
			t.Errorf("Codef(%v, %q): got: %v, want %v", tt.code, tt.message, err.code, tt.wantCode)
		}

		if got.Error() != tt.wangString {
			t.Errorf("Codef(%v, %q): got: %v, want %v", tt.code, tt.message, got.Error(), tt.wangString)
		}
	}
}

func TestCodef(t *testing.T) {
	tests := []struct {
		code       int
		format     string
		args       []interface{}
		wantType   string
		wantCode   int
		wangString string
	}{
		{errConfigurationNotValid, "Configuration error", nil, "*withCode", errConfigurationNotValid, "Configuration error"},
		{errConfigurationNotValid, "Configuration %s", []interface{}{"failed"}, "*withCode", errConfigurationNotValid, "Configuration failed"},
	}

	for _, tt := range tests {
		got := Codef(tt.code, tt.format, tt.args...)
		err, ok := got.(*withCode)
		if !ok {
			t.Errorf("Codef(%v, %q %q): error type got: %T, want %s", tt.code, tt.format, tt.args, got, tt.wantType)
		}

		if err.code != tt.wantCode {
			t.Errorf("Codef(%v, %q %q): got: %v, want %v", tt.code, tt.format, tt.args, err.code, tt.wantCode)
		}

		if got.Error() != tt.wangString {
			t.Errorf("Codef(%v, %q %q): got: %v, want %v", tt.code, tt.format, tt.args, got.Error(), tt.wangString)
		}
	}
}

func TestWithCode(t *testing.T) {
	type args struct {
		err     error
		code    int
		message string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{
				err:     New("err"),
				code:    errEOF,
				message: "err is nil",
			},
		},
		{
			name: "err is nil",
			args: args{
				err:     nil,
				code:    errEOF,
				message: "err is nil",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { _ = WithCode(tt.args.err, tt.args.code, tt.args.message) })
	}
}

func TestWithCodef(t *testing.T) {
	type args struct {
		err    error
		code   int
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{
				err:    New("err"),
				code:   errEOF,
				format: "err %s",
				args:   []interface{}{"test"},
			},
		},
		{
			name: "err is nil",
			args: args{
				err:    nil,
				code:   errEOF,
				format: "err is nil",
				args:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = WithCodef(tt.args.err, tt.args.code, tt.args.format, tt.args.args...)
		})
	}
}

// errors.New, etc values are not expected to be compared by value
// but the change in errors#27 made them incomparable. Assert that
// various kinds of errors have a functional equality operator, even
// if the result of that equality is always false.
func TestErrorEquality(t *testing.T) {
	vals := []error{
		nil,
		io.EOF,
		errors.New("EOF"),
		New("EOF"),
		Errorf("EOF"),
		Wrap(io.EOF, "EOF"),
		Wrapf(io.EOF, "EOF%d", 2),
		WithMessage(nil, "whoops"),
		WithMessage(io.EOF, "whoops"),
		WithStack(io.EOF),
		WithStack(nil),
	}

	for i := range vals {
		for j := range vals {
			_ = vals[i] == vals[j] // mustn't panic
		}
	}
}
