package main

import (
	"log"
	"time"

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

func snowflake(n *maelstrom.Node) uint64 {
	unix := uint64(time.Now().UnixMilli())
	nodeID := findIndex(n.NodeIDs(), n.ID())
	result := unix<<24 | ((uint64(nodeID) & 0xf) << 20) | (uint64(idCounter) & 0xfffff)
	idCounter++
	if idCounter == 1<<20 {
		idCounter = 0
	}
	return result
}

func main() {
	n := maelstrom.NewNode()
	n.Handle("generate", func(msg maelstrom.Message) error {
		body := make(map[string]interface{})
		body["type"] = "generate_ok"
		body["id"] = snowflake(n)

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
