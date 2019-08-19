package qerr

import (
	"fmt"
	"runtime"
	"strings"
)

type StackCut struct {
	Func  string
	Line  int
	File  string
	Stack []byte
}

func StackStringErr(skip int, stackline int, msg string, args ...interface{}) error {
	return fmt.Errorf(StackString(skip+1, stackline, msg, args...))
}

func StackErr(skip int, stackline int, err error, msg string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	if len(msg) == 0 {
		return fmt.Errorf(StackString(skip+1, stackline, err.Error()))
	} else {
		return fmt.Errorf(StackString(skip+1, stackline, msg+" | "+err.Error(), args...))
	}

}

func StackString(skip int, stackline int, msg string, args ...interface{}) string {
	if len(msg) > 0 {
		if args == nil || len(args) == 0 {
			msg = msg + " | "
		} else {
			msg = fmt.Sprintf(msg+" | ", args...)
		}
	}
	var cut = StackCutting(skip+1, stackline)
	if cut.Stack == nil {
		return fmt.Sprintf("%s%s | %s : %d", msg, cut.File, cut.Func, cut.Line)
	} else {
		return fmt.Sprintf("%s%s | %s : %d\n%s\n\n", msg, cut.File, cut.Func, cut.Line, string(cut.Stack))
	}
}

func StackCutting(skip int, stackline int) *StackCut {
	var pc, filename, linenum, _ = runtime.Caller(skip + 1)
	var slashindex = strings.LastIndex(filename, "/")
	filename = filename[slashindex+1:]
	var funcname = runtime.FuncForPC(pc).Name()
	// adjust buffer size to be larger than expected stack

	var stackbytes []byte
	if stackline >= 0 {

		skip = (skip * 2) + 6
		// stackline = stackline + 1
		stackbytes = make([]byte, 16*1024)
		var stacklen = runtime.Stack(stackbytes, false)
		var from = 0
		var to = stacklen
		var ncount = 0
		var linecount = 0
		for i := 0; i < stacklen; i++ {

			if stackbytes[i] != '\n' {
				continue
			}

			ncount++

			if ncount%2 == 0 {
				stackbytes[i] = ' '
			} else {

				if ncount <= skip {
					from = i + 1
					continue
				}

				linecount++
				if linecount >= stackline {
					to = i
					break
				}
			}
		}
		stackbytes = stackbytes[from:to]
	}

	return &StackCut{
		Func:  funcname,
		Line:  linenum,
		File:  filename,
		Stack: stackbytes,
	}

}

func StackCuttingMap(skip int, stackine int) map[string]interface{} {
	var cut = StackCutting(skip+1, stackine)
	var r = make(map[string]interface{})
	r["func"] = cut.Func
	r["line"] = cut.Line
	r["file"] = cut.File
	if cut.Stack != nil {
		r["stack"] = string(cut.Stack)
	}
	return r

}
