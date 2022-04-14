//go:build go1.13
// +build go1.13

package errors

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAggregate(t *testing.T) {
	type args struct {
		errList []error
	}
	tests := []struct {
		name string
		args args
		want Aggregate
	}{
		{
			name: "error list is empty",
			args: args{
				errList: []error{},
			},
			want: nil,
		},
		{
			name: "error is nil",
			args: args{
				errList: []error{nil, nil},
			},
			want: nil,
		},
		{
			name: "errList is nil",
			args: args{
				errList: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewAggregate(tt.args.errList...), "NewAggregate(%v)", tt.args.errList)
		})
	}
}

func Test_aggregate_Error(t *testing.T) {
	tests := []struct {
		name string
		agg  Aggregate
		want string
	}{
		{
			name: "normal",
			agg:  NewAggregate(New("err")),
			want: "err",
		},
		{
			name: "multiple errors",
			agg:  NewAggregate(New("err-1"), New("err-2")),
			want: "[err-1, err-2]",
		},
		{
			name: "multiple errors, there are duplicate errors",
			agg:  NewAggregate(New("err-1"), New("err-1")),
			want: "err-1",
		},
		{
			name: "multiple errors, there are duplicate errors",
			agg:  NewAggregate(New("err-1"), New("err-2"), New("err-2")),
			want: "[err-1, err-2]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.agg.Error(), "Error()")
		})
	}
}

func Test_aggregate_Error_panic(t *testing.T) {
	tests := []struct {
		name string
		agg  aggregate
	}{
		{
			name: "panic",
			agg:  nil,
		},
		{
			name: "panic",
			agg:  aggregate{},
		},
		{
			name: "panic",
			agg:  aggregate([]error{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() { _ = tt.agg.Error() })
		})
	}
}

type mockAggregate []error

func (m mockAggregate) Error() string {
	return "mockAggregate"
}

func (m mockAggregate) Errors() []error {
	return m
}

func (m mockAggregate) Is(target error) bool {
	return m.visit(func(err error) bool {
		return Is(err, target)
	})
}

func (m mockAggregate) visit(f func(err error) bool) bool {
	for _, err := range m {
		switch err := err.(type) {
		case Aggregate:
			for _, nestedErr := range err.Errors() {
				if match := f(nestedErr); match {
					return match
				}
			}
		default:
			if match := f(err); match {
				return match
			}
		}
	}

	return false
}

var _ Aggregate = mockAggregate{}

func Test_aggregate_Is(t *testing.T) {
	err1 := New("err-1")
	err2 := New("err-2")
	err3 := New("err-3")
	type args struct {
		target error
	}
	tests := []struct {
		name string
		agg  aggregate
		args args
		want bool
	}{
		{
			name: "normal",
			agg:  aggregate{err1, err2, err3},
			args: args{
				target: err3,
			},
			want: true,
		},
		{
			name: "normal",
			agg:  aggregate{err1, aggregate{err2, err3}},
			args: args{
				target: err3,
			},
			want: true,
		},
		{
			name: "normal",
			agg:  aggregate{err1, mockAggregate{err2, err3}},
			args: args{
				target: err3,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.agg.Is(tt.args.target), "Is(%v)", tt.args.target)
		})
	}
}

func Test_aggregate_Errors(t *testing.T) {
	errs := []error{loadConfig()}
	tests := []struct {
		name string
		agg  aggregate
		want []error
	}{
		{
			name: "normal",
			agg:  aggregate(errs),
			want: errs,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.agg.Errors(), "Errors()")
		})
	}
}

func TestFilterOut(t *testing.T) {
	type args struct {
		err error
		fns []Matcher
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				err: NewAggregate(io.EOF, loadConfig()),
				fns: []Matcher{
					func(err error) bool {
						return Is(err, io.EOF)
					},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "service configuration could not be loaded")
				return true
			},
		},
		{
			name: "err is nil",
			args: args{
				err: nil,
				fns: nil,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, FilterOut(tt.args.err, tt.args.fns...), fmt.Sprintf("FilterOut(%v, ...Matcher)", tt.args.err))
		})
	}
}

func TestFlatten(t *testing.T) {
	type args struct {
		agg Aggregate
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{
				agg: NewAggregate(NewAggregate(io.EOF, loadConfig())),
			},
			want: "[EOF, service configuration could not be loaded]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Flatten(tt.args.agg).Error(), "Flatten(%v)", tt.args.agg)
		})
	}
}

func TestFlatten_nil(t *testing.T) {
	type args struct {
		agg Aggregate
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "normal",
			args: args{
				agg: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Flatten(tt.args.agg), "Flatten(%v)", tt.args.agg)
		})
	}
}

func TestAggregateGoroutines(t *testing.T) {
	type args struct {
		funcs []func() error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{
				funcs: []func() error{
					func() error {
						return loadConfig()
					},
					func() error {
						return loadConfig()
					},
				},
			},
			want: NewAggregate(loadConfig(), loadConfig()).Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, AggregateGoroutines(tt.args.funcs...).Error(), "AggregateGoroutines()")
		})
	}
}

func TestAggregateGoroutines_nil(t *testing.T) {
	type args struct {
		funcs []func() error
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "normal",
			args: args{
				funcs: []func() error{
					func() error {
						return nil
					},
					func() error {
						return nil
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, AggregateGoroutines(tt.args.funcs...), "AggregateGoroutines()")
		})
	}
}
