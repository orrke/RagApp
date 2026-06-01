package logging

var LogLevelMap = map[string]int{
	"trace": 0,
	"debug": 1,
	"info":  2,
	"warn":  3,
	"error": 4,
	"fatal": 5,
}
var LogLevel int = LogLevelMap["trace"]

func Fatal(msg string) {
	Logger.Fatal(msg)
}

func Error(msg string) {
	if LogLevelMap["error"] >= LogLevel {
		return
	}
	Logger.Println("Error: " + msg)
}

func Warn(msg string) {
	if LogLevelMap["warn"] >= LogLevel {
		return
	}
	Logger.Println("Warn: " + msg)
}

func Info(msg string) {
	if LogLevelMap["info"] >= LogLevel {
		return
	}
	Logger.Println("Info: " + msg)
}

func Debug(msg string) {
	if LogLevelMap["debug"] >= LogLevel {
		return
	}
	Logger.Println("Debug: " + msg)
}

func Trace(msg string) {
	if LogLevelMap["trace"] >= LogLevel {
		return
	}
	Logger.Println("Trace: " + msg)
}
