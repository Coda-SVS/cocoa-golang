package log

var (
	GLogger *Logger
)

func init() {
	GLogger = NewLogger()
}
