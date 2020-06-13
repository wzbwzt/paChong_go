package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

var  (
	//email163=`(\d+)@163.com`
	emailReg=`\w+@\w+\.\w+?`
	linkReg=`href="(https?://[\s\S]+?)"`
	phoneReg=`1[3456789]\d\s?\d{4}\s?\d{4}` //手机号有或者没有空格
	idcardReg=`[1-9]\d{5}((19\d{2})|(20[01]\d))((0[1-9])|(1[012]))((0[1-9])|([12]\d)|(3[01]))\d{3}[\dXx]`//身份证号
	imgReg=`https?://[^"]+?(\.((jpg)|(png)|(jpeg)|(gif)|(bmp)))`
)

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


//GetEmail 爬邮箱
func GetEmail(url string){
	fileStr := PointObj(url)
	//3.根据正则来过滤数据
	re := regexp.MustCompile(emailReg)
	results := re.FindAllStringSubmatch(fileStr, -1) //-1：全部；可以指定取的数量
	fmt.Println(results)
	//for _,result:= range results{
	//	fmt.Println("email:",result[0])
	//	fmt.Println("account:",result[1])
	//}

}
//GetLink 爬链接
func GetLink(url string){
	fileStr := PointObj(url)
	re := regexp.MustCompile(linkReg)
	results := re.FindAllStringSubmatch(fileStr, -1)
	//fmt.Println(results)
	for _,result:= range results{
		fmt.Println(result[1])
	}
}
//GetPhone 爬手机号
func GetPhone(url string){
	fileStr := PointObj(url)
	re := regexp.MustCompile(phoneReg)
	results := re.FindAllStringSubmatch(fileStr, -1)
	fmt.Println(results)

}
//GetIdCard 爬身份证号
func GetIdCard(url string){
	fileStr := PointObj(url)
	re := regexp.MustCompile(idcardReg)
	results := re.FindAllStringSubmatch(fileStr, -1)
	fmt.Println(results)
}
//GetImg 爬图片链接
func GetImg(url string){
	fileStr := PointObj(url)
	re := regexp.MustCompile(imgReg)
	results := re.FindAllStringSubmatch(fileStr, -1)
	for _,result:= range results{
		fmt.Println(result[0])

	}
}

func main(){
	//GetEmail("https://tieba.baidu.com/p/6714901467")
	//GetLink("https://tieba.baidu.com/p/6714901467")
	//GetPhone("http://www.1686888.com/?bd_vid=7193869533034878773")
	//GetIdCard("https://henan.qq.com/a/20171107/069413.htm")
	GetImg("https://tieba.baidu.com/f?ie=utf-8&kw=%E7%9C%BC%E9%95%9C%E5%A8%98")
}