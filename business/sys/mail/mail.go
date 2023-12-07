package mail

import (
	"log"

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
	// TODO: send email
}
