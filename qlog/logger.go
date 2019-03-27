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
const CODEINFO = 1000

var LEVELSTRS = [6]string{"VERBOSE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

type Logi struct {
	key      string
	today    *time.Time
	todayday int
	writers  []io.Writer
	agents   []*log.Logger
	//logChannel chan * string;
	lock       sync.RWMutex
	Dir        string
	Level      int
	ToStdout   bool
	LogFlag    int
	LogPrefix  string
	FilePrefix string
	FileSuffix string
}

var _loggerManager = &LogManager{
	def: &Logi{
		key:        "",
		Dir:        "log",
		FileSuffix: ".log",
		Level:      INFO,
		ToStdout:   true,
		LogFlag:    log.Ltime,
	},
	loggers: map[string]*Logi{},
}

type LogManager struct {
	def     *Logi
	mutex   sync.RWMutex
	loggers map[string]*Logi
}

func GetLogManager() *LogManager {
	return _loggerManager
}

func (o *LogManager) GetDef() *Logi {
	return o.def
}

func (o *LogManager) Get(key string) *Logi {
	if len(key) == 0 {
		return o.def
	}
	o.mutex.RLock()
	var logger = o.loggers[key]
	o.mutex.Unlock()
	return logger
}

func (o *LogManager) Set(key string, logi *Logi) {
	if len(key) == 0 {
		o.def = logi
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.loggers[key] = logi
}

func (o *LogManager) New(key string, dir string, prefix string, suffix string, level string, stdout bool, logflag int) *Logi {
	var logger = o.Get(key)
	if logger != nil {
		return logger
	}
	var inst = &Logi{}
	inst.key = key
	inst.Dir = dir
	inst.FilePrefix = prefix
	inst.FileSuffix = suffix
	inst.Level = o.LevelInt(level)
	inst.ToStdout = stdout
	if inst.FileSuffix == "" {
		inst.FileSuffix = ".log"
	}
	if logflag > 0 {
		inst.LogFlag = logflag
	}
	o.Set(key, inst)
	return inst
}

func (o *LogManager) LevelStr(level int) string {
	if level < 0 || level >= len(LEVELSTRS) {
		return "VERBOSE"
	}
	return LEVELSTRS[level]
}

func (o *LogManager) LevelInt(level string) int {
	level = strings.ToLower(level)
	switch level {
	case "verbose":
		return VERBOSE
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	}
	return VERBOSE
}

/* ===================================== Logi ======================================= */

func (o *Logi) Destroy(lock bool) {
	if lock {
		o.lock.Lock()
		defer o.lock.Unlock()
	}
	if o.writers != nil {
		for _, writer := range o.writers {
			if writer != nil {
				if writer != os.Stdout && writer != os.Stderr {
					var closer, ok = writer.(io.Closer)
					if ok {
						closer.Close()
					}
				}
			}
		}
		o.writers = nil
		o.agents = nil
	}
	o.agents = nil
}

func (o *Logi) AddWriter(writer io.Writer, prefix string, flag int, lock bool) {
	if writer != nil {
		if flag <= 0 {
			flag = o.LogFlag
		}
		if len(prefix) == 0 {
			prefix = o.LogPrefix
		}
		if lock {
			o.lock.Lock()
			defer o.lock.Unlock()
		}
		if o.writers == nil {
			o.writers = make([]io.Writer, 1)
		}
		if o.agents == nil {
			o.agents = make([]*log.Logger, 1)
		}
		var nospace = true
		var agent = log.New(writer, prefix, flag)
		for i := 0; i < len(o.writers); i++ {
			if o.writers[i] == nil {
				o.writers[i] = writer
				nospace = false
				break
			}
		}
		if nospace {
			o.writers = append(o.writers, writer)
		}
		for i := 0; i < len(o.agents); i++ {
			if o.agents[i] == nil {
				o.agents[i] = agent
				nospace = false
				break
			}
		}
		if nospace {
			o.agents = append(o.agents, agent)
		}
	}
}

func (o *Logi) InitWriter(today *time.Time) {
	o.lock.Lock()
	defer o.lock.Unlock()
	if o.writers != nil {
		return
	}
	o.agents = nil

	o.today = today
	o.todayday = today.Day()
	var filename = o.FilePrefix + today.Format("20060102") + o.FileSuffix
	var filepath = o.Dir + "/" + filename
	if err := os.MkdirAll("log", 0774); err != nil {
		panic(err)
	}

	var file, err = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	o.AddWriter(file, o.LogPrefix, o.LogFlag, false)
	if o.ToStdout {
		o.AddWriter(os.Stdout, o.LogPrefix, o.LogFlag, false)
	}
}

func (o *Logi) SetWriters(writers []io.Writer) {
	o.lock.Lock()
	defer o.lock.Unlock()
	if o.writers != nil {
		o.Destroy(false)
	}
	for _, writer := range writers {
		o.AddWriter(writer, "", o.LogFlag, false)
	}
}

func (o *Logi) LogEx(level int, stackSkip int, v ...interface{}) {
	var trace = o.Level >= TRACE
	if trace {
		o.Level = o.Level - TRACE
	}

	if level < o.Level {
		return
	}

	if level >= ERROR {
		trace = true
	}

	var today = time.Now()
	var todayday = today.Day()

	if todayday != o.todayday {
		o.Destroy(true)
	}

	if o.writers == nil {
		o.InitWriter(&today)
	}
	var levelstr = LEVELSTRS[level]

	var stackstr = ""
	var line = ""
	var vs = util.SliceToString(" ", v...)
	if stackSkip >= 0 {
		var linenum = 0
		var filename = ""
		var funcname = ""
		var pc uintptr
		pc, filename, linenum, _ = runtime.Caller(stackSkip)
		var slashindex = strings.LastIndex(filename, "/")
		filename = filename[slashindex+1:]
		funcname = runtime.FuncForPC(pc).Name()

		if trace {
			// adjust buffer size to be larger than expected stack
			var bytes = make([]byte, 8192)
			var stack = runtime.Stack(bytes, false)
			stackstr = string(bytes[:stack])
		}
		line = fmt.Sprintf("%s %s %d %s   %s", levelstr, filename, linenum, funcname, vs)
	} else {
		line = fmt.Sprintf("%s   %s", levelstr, vs)
	}

	o.Print(line, stackstr)
}

func (o *Logi) Print(line string, stackstr string) {
	for _, agent := range o.agents {
		if agent != nil {
			agent.Println(line)
			if len(stackstr) > 0 {
				agent.Println(stackstr)
			}
		}
	}
}

func (o *Logi) Printf(format string, v ...interface{}) {
	for _, agent := range o.agents {
		if agent != nil {
			agent.Printf(format, v)
		}
	}
}

func (o *Logi) Error(skipStack int, v ...interface{}) {
	o.LogEx(ERROR, 2+skipStack, v...)
}

func (o *Logi) Log(level int, v ...interface{}) {
	o.LogEx(level, 2, v...)
}

func LogEx(level int, stackSkip int, v ...interface{}) {
	_loggerManager.def.LogEx(level, stackSkip, v...)
}

func Error(stackSkip int, v ...interface{}) {
	_loggerManager.def.LogEx(ERROR, 2+stackSkip, v...)
}

func Log(level int, v ...interface{}) {
	_loggerManager.def.LogEx(level, 2, v...)
}

func LevelInt(level string) int {
	return _loggerManager.LevelInt(level)
}

func LevelStr(level int) string {
	return _loggerManager.LevelStr(level)
}
