package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//正则变量
var (
	imgReg=`https?://[^"]+?(\.((jpg)|(png)|(jpeg)|(gif)|(bmp)))`
)

//task 任务实体
type task struct {
	f func()
}

//workPool 协程池实体
type workPool struct {
	jobChan chan *task
	entryChan chan *task
	maxWorkerNum int
}
//NewWorkPool 实例化一个协程池
func NewWorkPool(maxNum int)(pool *workPool){
	pool=&workPool{
		jobChan: make(chan *task),
		entryChan: make(chan *task),
		maxWorkerNum: maxNum,
	}
	return
}

func NewTask(fun func())(t *task){
	t=&task{
		f: fun,
	}
	return
}
//Excute 执行具体任务
func (t *task)Excute(){
	t.f()
}

//Work 每个worker的具体任务
func (w *workPool)Work(i int){
	for task:= range w.jobChan{
		task.Excute()
		fmt.Println("worker ID:", i, "has execute a task")
	}
}

//Run 协程池运行
func (w *workPool)Run(){
	for i:=0;i<w.maxWorkerNum;i++ {
		go w.Work(i)
	}
	for job := range w.entryChan{
			w.jobChan<-job
	}

}
//HandleError 处理错误
func HandleError(err error,why string){
	if err != nil {
		fmt.Println("Error:",why)
	}
}
//PointObj 指定网站拿数据并读取页面内容
func PointObj(url string)(fileStr string){
	//1.指定网站拿数据
	resp, err := http.Get(url)
	HandleError(err,"http.Get url")
	defer resp.Body.Close()
	//2.读取页面内容
	fileByte, err := ioutil.ReadAll(resp.Body)
	HandleError(err,"ioutil.ReadAll")
	fileStr=string(fileByte)
	return
}
//GetImgName 获取图片名称
func GetImgName(urlStr string)(fileName string){
	index := strings.LastIndex(urlStr, "/")
	fileName=urlStr[index+1:]
	fileNamePrefix:=time.Now().UnixNano()
	fileName=strconv.Itoa(int(fileNamePrefix))+"_"+fileName
	return
}

//DownImg 下载图片到本地
func DownImg(linkUrl string){
	fileName := GetImgName(linkUrl)
	resp, err := http.Get(linkUrl)
	HandleError(err,"http.get")
	defer resp.Body.Close()
	imgByte, err := ioutil.ReadAll(resp.Body)
	HandleError(err,"ioutil.readall")
	fileName="D://www/paChong_go/img/"+fileName
	err= ioutil.WriteFile(fileName, imgByte, 0666)
	HandleError(err,"ioutil.writeFile func: downImg")
}

//f 根据爬虫目标的地址来下载图片到本地
func PaImgToLocal(url string){
	fileStr := PointObj(url)
	re := regexp.MustCompile(imgReg)
	results := re.FindAllStringSubmatch(fileStr, -1)
	for _,result:= range results{
		DownImg(result[0])
	}

}
//fTrue 将函数包装成需要的函数格式
func fTrue(fun func(string),url string)(fT func ()){
	fT=func(){
		fun(url)
	}
	return
}

func main(){
	p:=NewWorkPool(20)
	for i:=1;i<30;i++ {
		url:="http://wmtp.net/page/"
		itoa := strconv.Itoa(i)
		url=url+itoa
		fT := fTrue(PaImgToLocal, url)
		newTask := NewTask(fT)
		go func(){
			p.entryChan<-newTask
		}()
	}
	p.Run()
}


