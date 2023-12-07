package mail

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/TSMC-Uber/server/business/sys/mq"
	"github.com/streadway/amqp"
)

func StartSendEmailWorker() {
	msgsReceiver := mq.NewDelayMsgsReceiver()
	for msg := range msgsReceiver {
		sendEmail(msg)
	}
}

func sendEmail(msg amqp.Delivery) {
	log.Printf("send email to %s", string(msg.Body))

	from := "tsmcuber@gmail.com"      // sender email address
	password := "muoq gjoo plxk ejhj" // email password
	to := string(msg.Body)            // recipient email address
	subject := "Hello"
	body := "This is a test email from Go!"
	// SMTP server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// 電子郵件消息格式
	message := []byte(
		"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, strings.Split(to, ","), message)
	if err != nil {
		fmt.Println("err:", err)
	}
}
