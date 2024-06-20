package main

import (
	"encoding/json"
	"log"
	"strconv"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var (
	idCounter uint64 = 0
)

func findIndex(arr []string, id string) int {
	for i, v := range arr {
		if v == id {
			return i
		}
	}
	return -1
}

func main() {
	n := maelstrom.NewNode()
	n.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		nodeID := strconv.Itoa(findIndex(n.NodeIDs(), n.ID()))
		msgID := strconv.Itoa(int(body["msg_id"].(float64)))

		body = make(map[string]any)
		body["type"] = "generate_ok"
		body["id"] = nodeID + "-" + msgID

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
