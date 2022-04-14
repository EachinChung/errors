//go:build go1.13
// +build go1.13

package errors

//goland:noinspection SpellCheckingInspection
import (
	mapset "github.com/deckarep/golang-set"
)

// Aggregate 表示包含多个错误的对象, 但不一定具有单一的语义含义。
// Aggregate 可以与 errors.Is() 一起使用, 以检查特定错误类型的发生。
// errors.As() 不支持, 因为调用者可能关心与给定类型匹配的潜在多个特定错误。
type Aggregate interface {
	error
	Errors() []error
	Is(error) bool
}

// NewAggregate 将 errList 转换为 Aggregate, Aggregate 本身就是 errors 接口的实现。
// 如果 slice 为空, 则返回 nil。
// 它将检查输入 errList 的任何元素是否为 nil, 以避免调用 Error() 时出现 nil panic。
func NewAggregate(errList ...error) Aggregate {
	if len(errList) == 0 {
		return nil
	}
	// 确保 errList 不包含 nil
	var errs []error
	for _, e := range errList {
		if e != nil {
			errs = append(errs, e)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return aggregate(errs)
}

// aggregate 实现了 error 与 Aggregate 接口。
type aggregate []error

// Error error 接口的一部分
func (agg aggregate) Error() string {
	if len(agg) == 0 {
		panic("error slice is empty")
	}
	if len(agg) == 1 {
		return agg[0].Error()
	}

	var result string
	// 不需要线程安全的 set
	seenErrs := mapset.NewThreadUnsafeSet()

	agg.visit(func(err error) bool {
		msg := err.Error()
		if seenErrs.Contains(msg) {
			return false
		}
		seenErrs.Add(msg)
		if seenErrs.Cardinality() > 1 {
			result += ", "
		}
		result += msg
		return false
	})
	if seenErrs.Cardinality() == 1 {
		return result
	}
	return "[" + result + "]"
}

func (agg aggregate) Is(target error) bool {
	return agg.visit(func(err error) bool {
		return Is(err, target)
	})
}

func (agg aggregate) visit(f func(err error) bool) bool {
	for _, err := range agg {
		switch err := err.(type) {
		case aggregate:
			if match := err.visit(f); match {
				return match
			}
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

// Errors 为 Aggregate 的一部分。
func (agg aggregate) Errors() []error { return agg }

// Matcher 用于匹配 errors。如果 errors 匹配, 则返回true。
type Matcher func(error) bool

// FilterOut 从输入错误中删除与 Matcher 匹配的错误。
// 如果输入是非 Aggregate error, 则仅测试该错误。
// 如果输入 Aggregate error, 错误列表将被递归处理。
//
// 例如, 这可以用于从错误列表中删除已知的错误 (例如 io.EOF 或 os.PathNotFound )。
func FilterOut(err error, fns ...Matcher) error {
	if err == nil {
		return nil
	}
	if agg, ok := err.(Aggregate); ok {
		return NewAggregate(filterErrors(agg.Errors(), fns...)...)
	}
	if !matchesError(err, fns...) {
		return err
	}
	return nil
}

// matchesError 如果有 Matcher 返回 true, 则返回 true
func matchesError(err error, fns ...Matcher) bool {
	for _, fn := range fns {
		if fn(err) {
			return true
		}
	}
	return false
}

// filterErrors 返回所有 fns 返回false 的任何 error (或嵌套错误, 如果列表包含嵌套错误)。
// 如果没有 error, 则返回 nil。副作用会使所有嵌套的切片变平
func filterErrors(list []error, fns ...Matcher) []error {
	var result []error
	for _, err := range list {
		r := FilterOut(err, fns...)
		if r != nil {
			result = append(result, r)
		}
	}
	return result
}

// Flatten 将可能嵌套 Aggregate 的 Aggregate 全部递归地压平为一个 Aggregate。
func Flatten(agg Aggregate) Aggregate {
	var result []error
	if agg == nil {
		return nil
	}
	for _, err := range agg.Errors() {
		if a, ok := err.(Aggregate); ok {
			r := Flatten(a)
			if r != nil {
				result = append(result, r.Errors()...)
			}
		} else {
			if err != nil {
				result = append(result, err)
			}
		}
	}
	return NewAggregate(result...)
}

// AggregateGoroutines 协程 error 收集器, 将所有非 nil error 填充到返回的 Aggregate 中。
// 如果所有函数均成功完成, 则返回 nil。
func AggregateGoroutines(funcs ...func() error) Aggregate {
	errChan := make(chan error, len(funcs))
	for _, f := range funcs {
		go func(f func() error) { errChan <- f() }(f)
	}
	errs := make([]error, 0)
	for i := 0; i < cap(errChan); i++ {
		if err := <-errChan; err != nil {
			errs = append(errs, err)
		}
	}
	return NewAggregate(errs...)
}
