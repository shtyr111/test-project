package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func StartLogger() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel) // уровень логирования
	log.SetReportCaller(true)

	//log.SetFormatter(&log.TextFormatter{
	//	FullTimestamp:   true,                      // полный формат времени
	//	TimestampFormat: "2006-01-02 15:04:05.000", // формат времени по Go layout
	//	ForceColors:     true,                      // цветной вывод в консоли
	//	CallerPrettyfier: func(f *runtime.Frame) (string, string) {
	//		// Красивый вывод имени функции и файла со строкой
	//		return "", fmt.Sprintf(" [ %s:%d ]:", filepath.Base(f.File), f.Line)
	//	},
	//	DisableLevelTruncation: true,
	//})

	log.Info("Started logger")
}
