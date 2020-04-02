/*
 * @Author: dcj
 * @Date: 2020-04-02 12:11:15
 * @LastEditTime: 2020-04-02 15:58:51
 * @Description: 提供邮件发送功能
 */

package email

import (
	"log"

	"github.com/go-gomail/gomail"
)

type EmailInfo struct {
	ServerHost string // ServerHost 邮箱服务器地址，如腾讯企业邮箱为smtp.exmail.qq.com
	ServerPort int    // ServerPort 邮箱服务器端口，如腾讯企业邮箱为465

	FromEmail  string // FromEmail　发件人邮箱地址
	FromPasswd string //发件人邮箱密码（注意，这里是明文形式)

	Recipient []string //收件人邮箱
	CC        []string //抄送
}

var emailMessage *gomail.Message

/**
 * @Author: dcj
 * @Date: 2020-04-02 15:45:55
 * @Description: 发送邮件
 * @Param : subject[主题]、body[内容]、emailInfo[发邮箱需要的信息(参考EmailInfo)]
 * @Return:
 */
func SendEmail(subject, body string, emailInfo *EmailInfo) {
	if len(emailInfo.Recipient) == 0 {
		log.Print("收件人列表为空")
		return
	}

	emailMessage = gomail.NewMessage()
	//设置收件人
	emailMessage.SetHeader("To", emailInfo.Recipient...)
	//设置抄送列表
	if len(emailInfo.CC) != 0 {
		emailMessage.SetHeader("Cc", emailInfo.CC...)
	}
	// 第三个参数为发件人别名，如"dcj"，可以为空（此时则为邮箱名称）
	emailMessage.SetAddressHeader("From", emailInfo.FromEmail, "dcj")

	//主题
	emailMessage.SetHeader("Subject", subject)

	//正文
	emailMessage.SetBody("text/html", body)

	d := gomail.NewPlainDialer(emailInfo.ServerHost, emailInfo.ServerPort,
		emailInfo.FromEmail, emailInfo.FromPasswd)
	err := d.DialAndSend(emailMessage)
	if err != nil {
		log.Println("发送邮件失败： ", err)
	} else {
		log.Println("已成功发送邮件到指定邮箱")
	}
}
