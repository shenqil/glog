package glog

import "github.com/sirupsen/logrus"

type Level logrus.Level

const (
	PanicLevel = (Level)(logrus.PanicLevel)
	FatalLevel = (Level)(logrus.FatalLevel)
	ErrorLevel = (Level)(logrus.ErrorLevel)
	WarnLevel  = (Level)(logrus.WarnLevel)
	InfoLevel  = (Level)(logrus.InfoLevel)
	DebugLevel = (Level)(logrus.DebugLevel)
	TraceLevel = (Level)(logrus.TraceLevel)
)

// Config 日志配置
type Config struct {
	DirName       string   // 日志目录名称
	NamePrefix    string   // 日志前置名称
	RetentionDays int64    // 日志保留时间
	FieldsOrder   []string // 字段顺序
	IsWriteToFile bool     // 是否写入文件中
	Level         Level    // 日志级别
}

// 初始一个默认值
var c = Config{
	DirName:       "log",
	NamePrefix:    "log",
	RetentionDays: 7 * 24 * 60 * 60,
	FieldsOrder:   []string{},
	IsWriteToFile: false,
	Level:         DebugLevel,
}

func setDirName(name string) {
	if name != "" {
		c.DirName = name
	}
}

func setNamePrefix(preName string) {
	if preName != "" {
		c.NamePrefix = preName
	}
}

func setRetentionDays(day int64) {
	if day > 0 {
		c.RetentionDays = day * 24 * 60 * 60
	}
}

func setFieldsOrder(FieldsOrder []string) {
	if FieldsOrder != nil {
		c.FieldsOrder = FieldsOrder
	}
}

func setIsWriteToFile(status bool) {
	c.IsWriteToFile = status
}

func setLevel(l Level) {
	if l >= PanicLevel && l <= TraceLevel {
		c.Level = l
	}
}
