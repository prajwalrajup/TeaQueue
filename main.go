package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rivo/tview"
)

var app *tview.Application
var textView *tview.TextView
var consumer *kafka.Consumer

func main() {
	app = tview.NewApplication()
	textView = tview.NewTextView()
	body := View(fmt.Sprintf("  Logs () "), "Hello World")

	baseFlex := buildBaseFlux(nil, body)

	// Set the root primitive to the flex layout
	app.SetRoot(baseFlex, true)

	// Run the application
	if err := app.Run(); err != nil {
		panic(err)
	}

	kafkaMetadata := KafkaMetadata{
		host:    "localhost",
		groupId: "myGroup",
	}
	err := getConsumer(kafkaMetadata)
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()

	topics, err := getTopics()
	if err != nil {
		log.Fatalf("Failed to get topics: %s", err)
	}
	kafkaMetadata.topics = topics

	fmt.Println(topics)

}

func View(title string, body string) tview.Primitive {

	textView = tview.NewTextView()

	textView.SetText(body)
	textView.SetDynamicColors(true)
	textView.SetTextColor(tview.Styles.PrimaryTextColor)
	textView.SetTitle(title)
	textView.SetBorderPadding(1, 1, 1, 1)
	textView.SetBorder(true)
	return textView
}

func buildBaseFlux(header tview.Primitive, body tview.Primitive) *tview.Flex {
	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 0, 1, false).
		AddItem(body, 0, 10, true)
}

type KafkaMetadata struct {
	host    string
	groupId string
	topics  []string
}

func getConsumer(kafkaMetadat KafkaMetadata) error {
	var err error
	consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaMetadat.host,
		"group.id":          kafkaMetadat.groupId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
		return err
	}

	return nil
}

func getTopics() ([]string, error) {
	if consumer == nil {
		return nil, errors.New("consumer not initialized")
	}
	metadata, err := consumer.GetMetadata(nil, true, 10000)
	if err != nil {
		return nil, err
	}

	// Extract and print topics
	var topicsList []string
	for topic := range metadata.Topics {
		if !strings.HasPrefix(topic, "__") {
			topicsList = append(topicsList, topic)
		}
	}
	return topicsList, nil
}
