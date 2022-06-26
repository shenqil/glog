因为[logrus](https://github.com/Sirupsen/logrus)原生不支持日志分割和自动清理，需要第三方插件或者自己实现，我们对[logrus](https://github.com/Sirupsen/logrus)进行简单封装，并增加了日志分割和自动清理。方便我们在后续项目中直接使用.
***
# 1.新建一个配置文件
```go
// Config 日志配置
type Config struct {
	DirName       string   // 日志目录名称
	NamePrefix    string   // 日志前置名称
	RetentionDays int64    // 日志保留时间
	FieldsOrder   []string // 字段顺序
	IsWriteToFile bool     // 是否写入文件中
}
```
+ 将常用配置给暴露出来，其他的配置都使用默认配置
***
# 2.读取配置的日志目录，清理过期的日志
```go
// clearExpiredLogs 清除过期日志任务
func clearExpiredLogs() {
	timeTemplate := "2006-01-2"
	entry, err := os.ReadDir(c.DirName)
	if err != nil {
		fmt.Println(err)
		return
	}

	retentionSeconds := c.RetentionDays * 24 * 60 * 60
	for _, v := range entry {
		if !v.IsDir() {
			name := v.Name()
			if strings.HasPrefix(name, c.NamePrefix) {
				timeStr := strings.Replace(name, c.NamePrefix+"_", "", 1)
				timeStr = strings.Replace(timeStr, ".log", "", 1)
				timeStr = strings.Replace(timeStr, "_", "-", 1)

				stamp, err := time.Parse(timeTemplate, timeStr)
				if err != nil {
					fmt.Println(err)
					continue
				}

				if time.Now().Unix()-stamp.Unix() >= retentionSeconds {
					os.Remove(fmt.Sprintf("./%s/%s", c.DirName, v.Name()))
				}
			}
		}
	}
}
```
+ 这里根据目录后面的时间戳与当前时间做对比，超出给定的保留时间，就执行删除
***
# 3.使用[cron](https://github.com/robfig/cron)创建一个每天凌晨执行的任务，去清理日志
```go
func startRotatingTask() error {
	c := cron.New()

	_, err := c.AddFunc("@daily", func() {
		fmt.Println("Every day")
		clearExpiredLogs()
	}) //  每天运行一次，午夜 | 0 0 * * *
	clearExpiredLogs()

	if err != nil {
		return err
	}

	c.Start()

	return nil
}
```
***
# 4.暴露一个初始化函数，用于给定配置
```go
// 日志实例
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
	if Logger != nil {
		err := startRotatingTask()
		if err != nil {
			panic(err)
		}
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
```
***
# 5.在其他项目使用
## 安装
```
go get -u  github.com/shenqil/glog
```

## 使用
```go
func main() {
	glog.InitLog(&glog.Config{
		IsWriteToFile: true,
	})

	l := glog.Logger
	l.Infof("this is %v demo", "glog")

	lWebServer := l.WithField("component", "web-server")
	lWebServer.Info("starting...")

	lWebServerReq := lWebServer.WithFields(logrus.Fields{
		"req":   "GET /api/stats",
		"reqId": "#1",
	})

	lWebServerReq.Info("params: startYear=2048")
	lWebServerReq.Error("response: 400 Bad Request")

	lDbConnector := l.WithField("category", "db-connector")
	lDbConnector.Info("connecting to db on 10.10.10.13...")
	lDbConnector.Warn("connection took 10s")

	l.Info("demo end.")
}
```

## 输出
![image.png](https://upload-images.jianshu.io/upload_images/25820166-475053b132fc3ded.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

***
[源码](http://github.com/shenqil/glog)