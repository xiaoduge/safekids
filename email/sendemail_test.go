/*
 * @Author: dcj
 * @Date: 2020-04-02 12:39:08
 * @LastEditTime: 2020-04-02 15:48:05
 * @Description: 发送邮件功能测试
 */
package email

import "testing"

func TestSendEmail(t *testing.T) {
	reclist := []string{"收件人地址"}

	info := &EmailInfo{
		"smtp.163.com",
		25,
		"发件人地址",
		"邮箱授权码",
		reclist,
		nil,
	}

	SendEmail("网页测试信息", "<h1>测试信息：</h1><p>您收到一条测试信息</p>", info)
}
