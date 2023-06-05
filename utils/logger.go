package utils

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"sync"
)

var myLog *logrus.Logger
var once sync.Once

func Logger() *logrus.Logger {
	once.Do(func() {
		myLog = logrus.New()              //创建logrus
		myLog.SetLevel(logrus.DebugLevel) //设置日志级别
		myLog.SetFormatter(
			&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"},
		) //设置时间格式
		output, err := getOutputFile()
		if err != nil {
			panic(err)
		}
		myLog.SetOutput(output)
		myLog.AddHook(new(myHook))
	})
	return myLog
}

func getOutputFile() (*os.File, error) {
	//获取绝对路径
	rootDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	filePath := rootDir + "/logs"
	_, err = os.Stat(filePath) //判断这个文件夹是否存在
	if os.IsNotExist(err) {
		//不存在这个文件夹就创建
		err := os.MkdirAll(filePath, 0666) //可读写，不可执行
		if err != nil {
			return nil, err
		}
	}

	//目录存在，那就直接执行下面的就行了
	//开始创建文件
	fileName := "Gee" + ".log"
	pathName := path.Join(filePath, fileName)
	_, err = os.Stat(pathName)
	if os.IsNotExist(err) {
		//文件不存在就创建文件
		//Create采用模式0666（任何人都可读写，不可执行）创建一个名为name的文件，如果文件已存在会截断它（为空文件）
		file, err := os.Create(pathName)
		if err != nil {
			return nil, err
		}
		return file, nil
	}

	//存在就不需要再打开了吧,是需要的?假如存在,但是程序调试开关多次,还是当天,但是logrusObj=nil,因为还没初始化
	return os.OpenFile(pathName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
}

type myHook struct {
}

func (m *myHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	}
}

func (m *myHook) Fire(entry *logrus.Entry) error {
	log.Println("==>", entry.Message)
	return nil
}
