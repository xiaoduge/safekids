/*
 * @Author: dcj
 * @Date: 2020-04-02 10:29:30
 * @LastEditTime: 2020-04-02 16:02:42
 * @Description: 用于测试指定网址是否可以访问
 */

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"safekids/email"
	"time"
)

const version string = "V0.0.1"

var emailInfo *email.EmailInfo

var logfile *os.File

func init() {
	var err error
	logfile, err = os.OpenFile("safikids.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logfile)
}

/**
 * @Author: dcj
 * @Date: 2020-04-02 15:48:45
 * @Description: 根据状态码确定是否需要发送信息，不需要发送则只记录日志
 * @Param : title[邮件标题]; text[邮件正文描述]; code[http状态码]
 * @Return:
 */
func SendMesssage(title, text string, code int) {
	now := time.Now()
	strNow := now.Format("2006/01/02 15:04:05")
	content := fmt.Sprintf("  %s: %d", text, code)
	log.Println(title, content)

	//请求状态码不是200，则发送邮件通知
	if code != 200 {
		body := fmt.Sprintf("<h1> %s </h1><p><pre> %s	%s </pre></p>", title, strNow, content)
		email.SendEmail("网页测试信息", body, emailInfo)
	}

}

/**
 * @Author: dcj
 * @Date: 2020-04-02 15:51:38
 * @Description: 处理一些不明确的状态码
 * @Param : code[http状态码]
 * @Return:
 */
func handlOtherCode(code int) {
	if code < 200 && code >= 100 {
		SendMesssage("测试网址", "服务器收到请求，需要请求者继续执行操作:", code)
	} else if code >= 200 && code < 300 {
		SendMesssage("测试网址", "成功，操作被成功接收并处理, 但是存在问题, 需要人工验证一下:", code)
	} else if code >= 300 && code < 400 {
		SendMesssage("测试网址", "重定向，需要进一步的操作以完成请求:", code)
	} else if code >= 400 && code < 500 {
		SendMesssage("测试网址", "客户端错误，请求包含语法错误或无法完成请求:", code)
	} else if code >= 500 && code < 600 {
		SendMesssage("测试网址", "服务器错误，服务器在处理请求的过程中发生了错误:", code)
	} else {
		SendMesssage("测试网址", "咱也不知道这个状态码是什么:", code)
	}
}

/**
 * @Author: dcj
 * @Date: 2020-04-02 15:52:42
 * @Description: 测试任务,用于测试指定网址是否可访问
 * @Param :
 * @Return:
 */
func CheckTask() {
	//发送get请求
	resp, err := http.Get("http://www.safekidschina.org/")
	if err != nil {
		SendMesssage("网址错误", err.Error(), 0)
		return
	}
	defer resp.Body.Close()

	//响应状态
	switch resp.StatusCode {
	case 200:
		SendMesssage("测试网址", "请求成功", resp.StatusCode)
	case 301:
		SendMesssage("测试网址", "资源（网页等）被永久转移到其它URL", resp.StatusCode)
	case 404:
		SendMesssage("测试网址", "请求的资源（网页等）不存在", resp.StatusCode)
	case 500:
		SendMesssage("测试网址", "内部服务器错误", resp.StatusCode)
	case 502:
		SendMesssage("测试网址", "从远程服务器接收到了一个无效的响应", resp.StatusCode)
	default:
		handlOtherCode(resp.StatusCode)
	}
}

/**
 * @Author: dcj
 * @Date: 2020-04-02 15:53:55
 * @Description: 启动定时任务
 * @Param : f[测试任务]
 * @Return:
 */
func startTask(f func()) {
	f() //开始运行则执行一次

	//新建计时器，go触发计时器的方法比较特别，就是在计时器的channel中发送值
	tick := time.NewTicker(10 * time.Second)

	for {
		select {
		//此处在等待channel中的信号，因此执行此段代码时会阻塞xxx时间
		case <-tick.C:
			f() //执行我们想要的操作
		}
	}
}

func main() {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("Panic err: ", err)
		}
		logfile.Close()
	}()

	log.Println("服务程序版本：", version)
	log.Println("服务程序启动...")

	reclist := []string{"邮箱地址"} //收件人列表
	cclist := []string{"邮箱地址"}  //抄送列表

	emailInfo = &email.EmailInfo{
		ServerHost: "smtp.qq.com",
		ServerPort: 25,
		FromEmail:  "邮箱地址", //发件人
		FromPasswd: "授权码",  //授权码
		Recipient:  reclist,
		CC:         cclist,
	}

	startTask(CheckTask)

	log.Println("服务程序结束...")
}
