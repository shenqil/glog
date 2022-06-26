package glog

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"os"
	"strings"
	"time"
)

// isExistedDir 判断目录是否存在
func isExistedDir(name string) bool {
	if info, err := os.Stat(name); err == nil {
		return info.IsDir()
	}

	return false
}

// startRotatingTask 开启一个每天0点清理一次日志任务
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
