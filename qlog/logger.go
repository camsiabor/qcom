package qlog

import (
	"bitbucket.org/avd/go-ipc/sync"
	"fmt"
	"github.com/camsiabor/qcom/util"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

const VERBOSE = 0
const DEBUG = 1
const INFO = 2
const WARN = 3
const ERROR = 4
const FATAL = 5

const TRACE = 100
const CODEINFO = 1000;

var LevelStr = [6]string{"VERBOSE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

type Logi struct {
	Dir      string;
	file     * os.File;
	today    time.Time;
	tomorrow time.Time;
	agents   [] * log.Logger;
	Flag     int;
	Level    int;
	ToStdout bool;
	lock     sync.RWMutex;
	FilePrefix string;
	FileSuffix string;
}


var _loggerManager = &LogManager{
	def : &Logi{
		Dir : "log",
		FileSuffix: ".log",
		Level : INFO,
		ToStdout: true,
		Flag : log.Ltime,
	},
	loggers : map[string] * Logi {},
}

type LogManager struct {
	def * Logi;
	loggers map[string] * Logi;
}

func GetLogManager() ( * LogManager) {
	return _loggerManager;
}

func (o * LogManager) GetDef() (* Logi) {
	return o.def;
}

func (o * LogManager) Get(key string) ( * Logi ) {
	if (len(key) == 0) {
		return o.def;
	}
	return o.loggers[key];
}

func (o * LogManager) Set(key string, logi * Logi) {
	if (len(key) == 0) {
		o.def = logi;
	}
	o.loggers[key] = logi;
}

func (o * Logi) Destroy () {
	o.lock.Lock();
	defer o.lock.Unlock();
	if (o.file != nil) {
		o.file.Sync();
		o.file.Close();
		o.file = nil;
	}
	o.agents = nil;
}

func (o * Logi) LogEx(level int, stackSkip int, v ... interface{}) {
	var trace= o.Level >= TRACE
	if trace {
		o.Level = o.Level - TRACE
	}

	if level < o.Level {
		return
	}

	if level >= ERROR {
		trace = true
	}

	var today= time.Now()

	if today.After(o.tomorrow) {
		o.Destroy();
	}

	if o.file == nil {
		func() {
			o.lock.Lock();
			defer o.lock.Lock();
			if (o.file != nil) {
				return;
			}

			o.today = today
			o.tomorrow = today.AddDate(0, 0, 1);
			var filename = o.FilePrefix + today.Format("20060102") + o.FileSuffix;
			var filepath = o.Dir + "/" + filename
			var err= os.MkdirAll("log", 0774)
			if err != nil {
				err.Error()
			}
			o.file, err = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

			if err != nil {
				err.Error()
			}
			o.agents = make([] * log.Logger, 2);
			o.agents = append(o.agents, log.New(o.file, "", o.Flag | log.Ltime));

			if o.ToStdout {
				o.agents = append(o.agents, log.New(os.Stdout, "", o.Flag | log.Ltime | log.Ldate));
			}
		}();


	}
	var levelstr = LevelStr[level]


	var linenum = 0;
	var filename = "";
	var funcname = "";
	var stackstr = "";
	var pc uintptr;

	pc, filename, linenum, _= runtime.Caller(stackSkip)
	var slashindex = strings.LastIndex(filename, "/");
	filename = filename[slashindex+1:]
	funcname = runtime.FuncForPC(pc).Name()

	if trace {
		// adjust buffer size to be larger than expected stack
		var bytes = make([]byte, 8192)
		var stack = runtime.Stack(bytes, false)
		stackstr = string(bytes[:stack]);
	}

	var vs = util.SliceToString(" ", v...);
	var line = fmt.Sprintf("%s %s %d %s   %s", levelstr, filename, linenum, funcname, vs);


	for _, agent := range o.agents {
		if (agent != nil) {
			agent.Println(line);
			if trace {
				agent.Println(stackstr);
			}
		}
	}
}

func (o * Logi) Error(skipStack int, v ... interface {}) {
	o.LogEx(ERROR, 2 + skipStack, v...);
}

func (o * Logi) Log(level int, v... interface{}) {
	o.LogEx(level, 2, v...);
}

func LogEx(level int, stackSkip int, v ... interface{}) {
	_loggerManager.def.LogEx(level, stackSkip, v...);
}

func Error(stackSkip int, v ... interface {}) {
	_loggerManager.def.LogEx(ERROR, 2 + stackSkip, v...);
}

func Log(level int, v ...interface{}) {
	_loggerManager.def.LogEx(level, 2, v...);
}

