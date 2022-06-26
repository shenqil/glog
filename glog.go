package glog

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/zput/zxcTool/ztLog"
	"github.com/zput/zxcTool/ztLog/zt_formatter"
	"os"
	"path"
	"runtime"
)

type Fields logrus.Fields

var Logger *logrus.Logger

// InitLog 初始化日志
func InitLog(conf *Config) {
	// 初始化配置
	setDirName(conf.DirName)
	setNamePrefix(conf.NamePrefix)
	setRetentionDays(conf.RetentionDays)
	setFieldsOrder(conf.FieldsOrder)
	setIsWriteToFile(conf.IsWriteToFile)

	// 创建日志目录
	if !isExistedDir(c.DirName) {
		if err := os.MkdirAll(c.DirName, 0777); err != nil {
			panic(err)
		}
	}

	// 开始日志清理任务
	err := startRotatingTask()
	if err != nil {
		panic(err)
	}

	// 创建日志实体
	createLog()
}

func createLog() {
	f := &zt_formatter.ZtFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		Formatter: nested.Formatter{
			FieldsOrder:     c.FieldsOrder,
			TimestampFormat: "2006-01-02T15:04:05.999999999Z07:00",
		},
	}

	Logger = logrus.New()

	if c.IsWriteToFile {
		ztLog.SetupLogsCanExpand(Logger, f, true, fmt.Sprintf("./%s/%s", c.DirName, c.NamePrefix), int(logrus.DebugLevel))
	} else {
		Logger.SetLevel(logrus.DebugLevel)
		Logger.SetReportCaller(true)
		Logger.SetFormatter(f)
	}
}
