package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Node struct {
	Food   string
	Number int
	Parent *Node
	Sons   []*Node
}

type FoodLink struct {
	Count int
	Links []*Node
}
type Foodlist map[string]FoodLink

func fptreeCreadeNode(food string, parent *Node) *Node {
	return &Node{food, 1, parent, []*Node{}}
}

func (rootp *Node) fptreeInsert(listp *Foodlist, food string) (node *Node) {
	list := *listp
	link := list[food]
	link.Count++
	for _, val := range rootp.Sons {
		if val.Food == food {
			node = val
			val.Number++
		}
	}
	if node == nil {
		node = fptreeCreadeNode(food, rootp)
		rootp.Sons = append(rootp.Sons, node)
		link.Links = append(list[food].Links, node)
	}
	list[food] = link
	return
}

func (rootp *Node) fptreeAppend(listp *Foodlist, foods ...string) {
	cur := rootp
	for _, food := range foods {
		cur = cur.fptreeInsert(listp, food)
	}
}

func main() {
	root := Node{Food: "", Number: 0, Sons: []*Node{}}
	nodeLink := make(Foodlist)
	dataset := make(map[int][]string)
	file, err := ioutil.ReadFile("dataset.json")
	if err != nil {
		log.Fatalln("Unable to read dataset.json", err)
	}
	json.Unmarshal(file, &dataset)

	log.Println("Construct FP-Tree...")
	for _, val := range dataset {
		root.fptreeAppend(&nodeLink, val...)
	}

	log.Println("FP-Tree Mining...")
}
