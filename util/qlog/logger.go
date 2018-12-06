package qlog

import (
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

var _dir string = "log"

var _file *os.File

var _today time.Time
var _tomorrow time.Time

var _logger *log.Logger
var _loggerstdout *log.Logger

var _logflags = log.Ltime
var _log2stdout = true

var _level = INFO

const VERBOSE = 0
const DEBUG = 1
const INFO = 2
const WARN = 3
const ERROR = 4
const FATAL = 5

const TRACE = 100
const CODEINFO = 1000;

var _levelstr = [6]string{"VERBOSE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}



func LogInit(dir string, level int, logflags int, log2stdout bool) {
	_dir = dir
	_level = level
	_logflags = logflags
	_log2stdout = log2stdout
}

func Log_destroy() {
	if _file != nil {
		_file.Sync()
		_file.Close()
		_file = nil
	}
}


func Log(level int, v ...interface{}) {

	var trace= level >= TRACE
	if trace {
		level = level - TRACE
	}

	//var codeinfo= level >= CODEINFO;
	//if (codeinfo) {
	//	level = level - CODEINFO;
	//}

	if level < _level {
		return
	}

	if level >= ERROR {
		trace = true
	}

	var today= time.Now()

	if today.After(_tomorrow) {
		Log_destroy()
	}

	if _file == nil {
		_today = today
		_tomorrow = today.AddDate(0, 0, 1);
		var filename= today.Format("20060102") + ".log";
		var filepath= _dir + "/" + filename
		var err= os.MkdirAll("log", 0774)
		if err != nil {
			err.Error()
		}
		_file, err = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

		if err != nil {
			err.Error()
		}

		_logger = log.New(_file, "", _logflags)

		if _log2stdout {
			_loggerstdout = log.New(os.Stdout, "", _logflags|log.Ldate)
		}
	}
	var levelstr= _levelstr[level]


	var linenum = 0;
	var filename = "";
	var funcname = "";
	var stackstr = "";
	var pc uintptr;

	pc, filename, linenum, _= runtime.Caller(1)
	var slashindex = strings.LastIndex(filename, "/");
	filename = filename[slashindex+1:]
	funcname = runtime.FuncForPC(pc).Name()

	if trace {
		// adjust buffer size to be larger than expected stack
		var bytes = make([]byte, 8192)
		var stack = runtime.Stack(bytes, false)
		stackstr = string(bytes[:stack])
	}
	if trace {
		_logger.Println(levelstr, filename, linenum, funcname, v, stackstr)
	} else {
		_logger.Println(levelstr, filename, linenum, funcname, v)
	}

	if _log2stdout {
		if trace {
			_loggerstdout.Println(levelstr, filename, linenum,funcname, v, stackstr)
		} else {
			_loggerstdout.Println(levelstr, filename, linenum,funcname, v)
		}
	}

}

