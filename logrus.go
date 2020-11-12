package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func InitLog(path string, logLevel int, fileName string, separatorTime int, expiredTime int) {
	filePath := path + string(os.PathSeparator) + fileName
	writer, err := rotatelogs.New(
		filePath+".%Y%m%d%H%M",
		// rotatelogs.WithLinkName(filePath), // 生成软链，指向最新日志文件
		// rotatelogs.WithMaxAge(time.Duration(s.router.container.Config.LogExpiredTime)*24*time.Hour),         // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Duration(separatorTime)*24*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &LogFormatter{})
	logrus.AddHook(lfHook)
	logrus.SetFormatter(&LogFormatter{})
	logrus.SetLevel(logrus.Level(logLevel - 2))
	go removeDir(path, fileName, separatorTime, expiredTime)
}

func removeDir(path string, fileName string, separatorTime int, expiredTime int) {
	for {
		fs, _ := ioutil.ReadDir(path)
		for _, v := range fs {
			if strings.Contains(v.Name(), fileName) {
				createTime := v.ModTime()
				expiredTime := time.Now().Add(-time.Duration(expiredTime)*24*time.Hour).Unix() * 1000
				if createTime.Unix()*1000 < expiredTime {
					os.RemoveAll(path + string(os.PathSeparator) + v.Name())
				}
			}
		}
		time.Sleep(time.Duration(separatorTime) * 24 * time.Hour)
	}
}

type LogFormatter struct {
}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	logLevel, _ := entry.Level.MarshalText()
	level := strings.ToUpper(string(logLevel))
	strList := findCaller(3)
	strFmt := fmt.Sprintf("[%s][%d][%s]", entry.Time.Format("2006-01-02 03:04:05.000"), -1, level)
	str := ""
	for _, v := range strList {
		if str == "" {
			str = strFmt + v
		} else {
			str = str + "\n" + strFmt + v
		}
	}
	b := []byte(str)
	if len(entry.Data) != 0 {
		serialized, err := json.Marshal(entry.Data)
		if err != nil {
			return serialized, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
		}
		b = append(b, serialized...)
	}
	b = append(b, []byte(entry.Message)...)
	return append(b, '\n'), nil
}

func findCaller(skip int) []string {
	strList := []string{}
	for i := 0; i < 20; i++ {
		file := ""
		line := 0
		var pc uintptr
		file, line, pc = getCaller(skip + i)
		if line == 0 {
			break
		}
		fullFnName := runtime.FuncForPC(pc)
		fnName := fullFnName.Name()
		if !strings.Contains(fnName, "dtsp_go/service") || strings.Contains(fnName, "PanicHandler") || strings.Contains(fnName, "CheckHandler") {
			continue
		}
		str := fmt.Sprintf("[%s][%s][%d]", file, fnName, line)
		strList = append(strList, str)
	}

	return strList
}

func getCaller(skip int) (string, int, uintptr) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0, pc
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line, pc
}
