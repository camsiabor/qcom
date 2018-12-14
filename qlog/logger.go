package qlog

import (
	"fmt"
	"github.com/camsiabor/qcom/util"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
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
	today      time.Time;
	tomorrow   time.Time;
	writers    [] io.WriteCloser;
	agents     [] * log.Logger;
	lock       sync.RWMutex;
	Dir        string;
	Level      int;
	ToStdout   bool;
	LogFlag    int;
	LogPrefix  string;
	FilePrefix string;
	FileSuffix string;
}


var _loggerManager = &LogManager{
	def : &Logi{
		Dir :       "log",
		FileSuffix: ".log",
		Level :     INFO,
		ToStdout:   true,
		LogFlag:    log.Ltime,
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
	if (o.writers != nil) {
		for _, writer := range o.writers {
			writer.Close();
		}
		o.writers = nil;
		o.agents = nil;
	}
	o.agents = nil;
}

func (o * Logi) AddWriter(writer io.WriteCloser, prefix string, flag int, lock bool) {
	if (writer != nil) {
		if (flag <= 0) {
			flag = o.LogFlag;
		}
		if (len(prefix) == 0) {
			prefix = o.LogPrefix;
		}
		if (lock) {
			o.lock.Lock();
			defer o.lock.Unlock();
		}
		if (o.writers == nil) {
			o.writers = make([]io.WriteCloser, 2);
		}
		if (o.agents == nil) {
			o.agents = make([] * log.Logger, 2);
		}
		o.writers = append(o.writers, writer);
		o.agents = append(o.agents, log.New(writer, prefix, flag));
	}
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

	if o.writers == nil {
		func() {
			o.lock.Lock();
			defer o.lock.Unlock();
			if (o.writers != nil) {
				return;
			}
			o.agents = nil;

			o.today = today
			o.tomorrow = today.AddDate(0, 0, 1);
			var filename = o.FilePrefix + today.Format("20060102") + o.FileSuffix;
			var filepath = o.Dir + "/" + filename
			if err := os.MkdirAll("log", 0774); err != nil {
				panic(err);
			}

			var file, err = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				panic(err);
			}

			o.AddWriter(file, o.LogPrefix, o.LogFlag, false);
			if o.ToStdout {
				o.AddWriter(os.Stdout, o.LogPrefix, o.LogFlag, false);
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
	go o.Print(line, stackstr);
}

func (o * Logi) Print(line string, stackstr string) {
	for _, agent := range o.agents {
		if (agent != nil) {
			agent.Println(line);
			if len(stackstr) > 0 {
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

