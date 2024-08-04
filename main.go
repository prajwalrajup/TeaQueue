// package main

// import (
// 	"fmt"

// 	"github.com/confluentinc/confluent-kafka-go/kafka"
// 	"github.com/rivo/tview"
// 	"time"
// 	// "kafka-ui/utils"
// )

// func logsView(host string, topic string) tview.Primitive {

// 	textView = tview.NewTextView()

// 	textView.SetText("sdf asf\nasdfasf")
// 	textView.SetDynamicColors(true)
// 	textView.SetTextColor(tview.Styles.PrimaryTextColor)
// 	textView.SetTitle(fmt.Sprintf("  Logs (%s:%s) ", host, topic))
// 	textView.SetBorderPadding(1, 1, 1, 1)
// 	textView.SetBorder(true)
// 	return textView
// }

// var app *tview.Application
// var textView *tview.TextView

// func main() {
// 	// Create a new application
// 	app = tview.NewApplication()

// 	// Create a layout to center the TextView
// 	flex := tview.NewFlex().
// 		SetDirection(tview.FlexRow).
// 		AddItem(nil, 0, 1, false).
// 		AddItem(logsView("local", "myTopic"), 0, 10, true)

// 	// Set the root primitive to the flex layout
// 	app.SetRoot(flex, true)

// 	// Run the application
// 	if err := app.Run(); err != nil {
// 		panic(err)
// 	}

// 	// Kafka consumer
// 	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
// 		"bootstrap.servers": "localhost:9092",
// 		"group.id":          "foo",
// 		"auto.offset.reset": "earliest"})
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = consumer.SubscribeTopics([]string{"myTopic"}, nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	go func() {
// 		run := true
// 		for run {
// 			msg, err := consumer.ReadMessage(time.Second)
// 			if err == nil {
// 				app.QueueUpdateDraw(func() {
// 					textView.SetText(string(msg.Value))
// 				})
// 				// fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
// 			}
// 			// else if !err.(kafka.Error).Error()() {
// 			// 	// The client will automatically try to recover from all errors.
// 			// 	// Timeout is not considered an error because it is raised by
// 			// 	// ReadMessage in absence of messages.
// 			// 	fmt.Printf("Consumer error: %v (%v)\n", err, msg)
// 			// }
// 		}
// 	}()
// 	// for _, line := range text {
// 	// 	fmt.Fprintln(textView, line)
// 	// }
// 	consumer.Close()
// }

package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rivo/tview"
	"log"
)

var app *tview.Application
var textView *tview.TextView
var consumer *kafka.Consumer

func main() {
	app = tview.NewApplication()
	textView = tview.NewTextView()

	// Initialize Kafka consumer
	var err error
	consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}

	err = consumer.Subscribe("myTopic", nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s", err)
	}

	go consumeMessages()

	if err := app.SetRoot(textView, true).Run(); err != nil {
		log.Fatalf("Error running application: %s", err)
	}
}

func consumeMessages() {
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			log.Printf("Received message: %s", string(msg.Value))
			// Process the message
			updateUI(string(msg.Value))
		} else {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func updateUI(message string) {
	app.QueueUpdateDraw(func() {
		log.Printf("Updating UI with message: %s", message)
		textView.SetText(message)
	})
}
