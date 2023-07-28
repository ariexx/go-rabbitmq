// consumer.go

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
)

var (
	queueName    = "notification_queue"
	templateName = "example_template.html"
)

type EmailData struct {
	To               string
	Title            string
	FirstRowMessage  string
	SecondRowMessage string
	ThirdRowMessage  string
	Footer           string
}

type Email struct {
	Email   string
	Message string
}

func main() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the exchange
	err = ch.ExchangeDeclare("notification_exchange", "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare the exchange: %v", err)
	}

	// Declare the queue
	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare the queue: %v", err)
	}

	// Bind the queue to the exchange
	err = ch.QueueBind(q.Name, "notification_key", "notification_exchange", false, nil)
	if err != nil {
		log.Fatalf("Failed to bind the queue to the exchange: %v", err)
	}

	// Consume messages from the queue
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Handle incoming messages
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			// marshal the data into a map
			var data Email

			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			// Process the message (send email)
			sendEmailWithTemplate(data.Email, data.Message)
			log.Printf("Received a notification: %s it sent to %s", data.Message, data.Email)
			d.Ack(false) // Acknowledge the message (remove it from the queue)
		}
	}()

	log.Println("Waiting for messages. To exit, press CTRL+C")
	<-forever
}

func sendEmailWithTemplate(to string, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@golang.org") // Replace with your email address
	m.SetHeader("To", to)                      // Replace with the recipient's email address
	m.SetHeader("Subject", body)               // Email subject
	tmpl, err := template.ParseFiles(templateName)
	if err != nil {
		fmt.Println("Failed to parse email template:", err)
		return
	}

	fmt.Println("First vim code")

	// Set the template data
	emailData := EmailData{
		To:               "Mr. Arief",
		Title:            body,
		FirstRowMessage:  "Selamat kamu mendapatkan hadiah sebesar Rp1.000.0000",
		SecondRowMessage: "Klik Disini Untuk Mengambil",
		ThirdRowMessage:  "Semoga dengan uang ini kamu bisa membeli apa yang kamu inginkan.",
		Footer:           "Good Luck, Stay Humble!",
	}

	// Execute the template with the data and store the output in a buffer
	var emailBodyContent string
	buffer := &bytes.Buffer{}
	err = tmpl.Execute(buffer, emailData)
	if err != nil {
		fmt.Println("Failed to execute email template:", err)
		return
	}

	// Get the content from the buffer and set it as the email body
	emailBodyContent = buffer.String()
	m.SetBody("text/html", emailBodyContent)
	d := gomail.NewDialer("smtp.mailtrap.io", 587, "27de36f4ba6439", "1752eeafc75a03") // Replace with your Mailtrap credentials
	d.TLSConfig = nil                                                                  // Disable SSL/TLS to work with Mailtrap

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Failed to send email:", err)
	}
}
