package errors

import (
	"bytes"
	"fmt"
)

// formatInfo contains all the error information.
type formatInfo struct {
	code    int
	message string
	err     string
	stack   *stack
}

//goland:noinspection GoUnhandledErrorResult
func format(k int, jsonData []map[string]interface{}, str *bytes.Buffer, info *formatInfo,
	sep string, flagDetail, flagTrace, modeJSON bool) ([]map[string]interface{}, *bytes.Buffer) {
	if modeJSON {
		data := map[string]interface{}{}
		if flagDetail || flagTrace {
			data = map[string]interface{}{
				"message": info.message,
				"code":    info.code,
				"error":   info.err,
			}

			caller := fmt.Sprintf("#%d", k)
			if info.stack != nil {
				f := Frame((*info.stack)[0])
				caller = fmt.Sprintf("%s %s:%d (%s)",
					caller,
					f.file(),
					f.line(),
					f.name(),
				)
			}
			data["caller"] = caller
		} else {
			data["error"] = info.message
		}
		jsonData = append(jsonData, data)
	} else {
		if flagDetail || flagTrace {
			if info.stack != nil {
				f := Frame((*info.stack)[0])
				fmt.Fprintf(str, "%s%s - #%d [%s:%d (%s)] (%d) %s",
					sep,
					info.err,
					k,
					f.file(),
					f.line(),
					f.name(),
					info.code,
					info.message,
				)
			} else {
				fmt.Fprintf(str, "%s%s - #%d %s", sep, info.err, k, info.message)
			}

		} else {
			fmt.Fprintf(str, info.err)
		}
	}

	return jsonData, str
}

// list 将错误堆栈转换为一个简单的数组
func list(e error) []error {
	var ret []error

	if e != nil {
		if w, ok := e.(interface{ Unwrap() error }); ok {
			ret = append(ret, e)
			ret = append(ret, list(w.Unwrap())...)
		} else {
			ret = append(ret, e)
		}
	}

	return ret
}

func buildFormatInfo(e error) *formatInfo {
	var info *formatInfo

	switch err := e.(type) {
	case *fundamental:
		info = &formatInfo{
			code:    unknownCoder.Code(),
			message: err.msg,
			err:     err.msg,
			stack:   err.stack,
		}
	case *withStack:
		info = &formatInfo{
			code:    unknownCoder.Code(),
			message: err.Error(),
			err:     err.Error(),
			stack:   err.stack,
		}
	case *withCode:
		coder, ok := codes[err.code]
		if !ok {
			coder = unknownCoder
		}

		extMsg := coder.String()
		if extMsg == "" {
			extMsg = err.msg
		}

		info = &formatInfo{
			code:    coder.Code(),
			message: extMsg,
			err:     err.msg,
			stack:   err.stack,
		}
	default:
		info = &formatInfo{
			code:    unknownCoder.Code(),
			message: err.Error(),
			err:     err.Error(),
		}
	}

	return info
}
