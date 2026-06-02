package logging

var LogLevelMap = map[string]int{
	"trace": 0,
	"debug": 1,
	"info":  2,
	"warn":  3,
	"error": 4,
	"fatal": 5,
}
var LogLevel = LogLevelMap["info"]

func Fatal(msg string) {
	Logger.Println("Fatal: " + msg)
	Logger.Fatal(msg)
}

func Error(msg string) {
	if LogLevelMap["error"] >= LogLevel {
		Logger.Println("Error: " + msg)
	}
}

func Warn(msg string) {
	if LogLevelMap["warn"] >= LogLevel {
		Logger.Println("Warn: " + msg)
	}
}

func Info(msg string) {
	if LogLevelMap["info"] >= LogLevel {
		Logger.Println("Info: " + msg)
	}
}

func Debug(msg string) {
	if LogLevelMap["debug"] >= LogLevel {
		Logger.Println("Debug: " + msg)
	}
}

func Trace(msg string) {
	if LogLevelMap["trace"] >= LogLevel {
		Logger.Println("Trace: " + msg)
	}
}
