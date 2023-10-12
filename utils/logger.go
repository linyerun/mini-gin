package utils

import (
	"github.com/linyerun/Mini-Gin/gee"
	"io"
	"log"
	"os"
	"path"
	"sync"
)

var myLog *log.Logger
var once sync.Once

func Logger() *log.Logger {
	once.Do(func() {
		output, err := getOutputFile()
		if err != nil {
			panic(err)
		}
		multiWriter := io.MultiWriter(output, os.Stdout)
		if w := gee.GetMiniGinLogOutputWriter(); w != nil {
			multiWriter = io.MultiWriter(w, os.Stdout)
		}
		myLog = log.New(multiWriter, "[Mini-Gin-Log]", log.LstdFlags|log.Llongfile)
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
	fileName := "gee" + ".log"
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
