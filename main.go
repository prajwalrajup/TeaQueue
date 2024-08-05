package main

import (
	"errors"
	"log"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var app *tview.Application
var textView *tview.TextView
var consumer *kafka.Consumer

var hosts = [][]string{
	{"localhost-9092", "localhost:9092", "local"},
	{"localhost-9093", "localhost:9093", "SSH"},
	{"localhost-9094", "localhost:9094", "Cloud"},
}

func main() {
	app = tview.NewApplication()
	textView = tview.NewTextView()

	// Define the callback function
	onItemSelected := func(index int, mainText string, secondaryText string, shortcut rune) {
		newItems := []string{"New Item 1", "New Item 2", "New Item 3"}
		body := listViewBlock("  New List () ", newItems, nil)
		baseFlex := buildBaseFlux(nil, body)
		app.SetRoot(baseFlex, true)
	}

	body := tableViewBlock("  Hosts () ", hosts, onItemSelected)

	baseFlex := buildBaseFlux(nil, body)

	// Set the root primitive to the flex layout
	app.SetRoot(baseFlex, true)

	// Run the application
	if err := app.Run(); err != nil {
		panic(err)
	}

	// kafkaMetadata := KafkaMetadata{
	// 	host:    "localhost",
	// 	groupId: "myGroup",
	// }
	// err := getConsumer(kafkaMetadata)
	// if err != nil {
	// 	log.Fatalf("Failed to create consumer: %s", err)
	// }
	// defer consumer.Close()

	// topics, err := getTopics()
	// if err != nil {
	// 	log.Fatalf("Failed to get topics: %s", err)
	// }
	// kafkaMetadata.topics = topics

	// fmt.Println(topics)

}

func listViewBlock(title string, items []string, onItemSelected func(int, string, string, rune)) tview.Primitive {
	listView := tview.NewList()
	listView.SetTitle(title)
	listView.SetBorder(true)
	listView.SetBorderPadding(1, 1, 1, 1)
	listView.SetMainTextColor(tcell.ColorLightSkyBlue)
	listView.SetSelectedBackgroundColor(tcell.ColorLightSkyBlue)

	// Add items to the list
	for index, item := range items {
		// Wrap the onItemSelected function in a closure
		listView.AddItem(item, "", 0, func(i int, mainText, secondaryText string, shortcut rune) func() {
			return func() {
				onItemSelected(i, mainText, secondaryText, shortcut)
			}
		}(index, item, "", 0))
	}

	return listView
}

func tableViewBlock(title string, items [][]string, onItemSelected func(int, string, string, rune)) tview.Primitive {
	tableView := tview.NewTable()
	tableView.SetTitle(title)
	tableView.SetBorder(true)
	tableView.SetBorderPadding(1, 1, 1, 1)
	tableView.SetSelectable(true, false)

	// Add header row
	headers := []string{"Name", "Broker Host", "Network"}
	for colIndex, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetSelectable(false)
		tableView.SetCell(0, colIndex, cell)
		cell.SetExpansion(1)
	}

	// Add items to the table
	for rowIndex, row := range items {
		for colIndex, cellText := range row {
			cell := tview.NewTableCell(cellText).
				SetTextColor(tcell.ColorLightSkyBlue).
				SetSelectable(true)
			tableView.SetCell(rowIndex+1, colIndex, cell)
		}
	}

	// Wrap the onItemSelected function in a closure
	tableView.SetSelectedFunc(func(row, column int) {
		if row > 0 { // Skip header row
			onItemSelected(row-1, items[row-1][0], items[row-1][1], 0)
		}
	})

	return tableView
}

func textViewBlock(title string, body string) tview.Primitive {

	textView = tview.NewTextView()

	textView.SetText(body)
	textView.SetDynamicColors(true)
	textView.SetTextColor(tview.Styles.PrimaryTextColor)
	textView.SetTitle(title)
	textView.SetBorderPadding(1, 1, 1, 1)
	textView.SetBorder(true)
	return textView
}

func buildBaseFlux(header, body tview.Primitive) *tview.Flex {
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
