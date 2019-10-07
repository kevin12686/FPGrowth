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

type Dataset map[int][]string

func fpTreeCreadeNode(food string, parent *Node) *Node {
	return &Node{food, 1, parent, []*Node{}}
}

func (rootp *Node) fpTreeInsert(list Foodlist, food string) (node *Node) {
	for _, val := range rootp.Sons {
		if val.Food == food {
			node = val
			val.Number++
		}
	}
	if node == nil {
		link := list[food]
		node = fpTreeCreadeNode(food, rootp)
		rootp.Sons = append(rootp.Sons, node)
		link.Links = append(list[food].Links, node)
		list[food] = link
	}
	return
}

func (rootp *Node) fpTreeAppend(list Foodlist, foods ...string) {
	cur := rootp
	for _, food := range foods {
		cur = cur.fpTreeInsert(list, food)
	}
}

func removeItem(ds Dataset, item string) {
	for num, orders := range ds {
		idx := -1
		for i, v := range orders {
			if v == item {
				idx = i
				break
			}
		}
		if idx > -1 {
			ds[num] = append(orders[:idx], orders[idx+1:]...)
		}
	}
}

func handleSupport(ds Dataset, minSpt float64) Foodlist {
	headTable := make(Foodlist)
	total := float64(len(ds))
	for _, orderList := range ds {
		for _, food := range orderList {
			link := headTable[food]
			link.Count++
			headTable[food] = link
		}
	}
	for food, link := range headTable {
		if spt := float64(link.Count) / total; spt < minSpt {
			removeItem(ds, food)
		}
	}
	return headTable
}

func main() {
	root := Node{Food: "", Number: 0, Sons: []*Node{}}
	dataset := make(Dataset)
	file, err := ioutil.ReadFile("dataset.json")
	if err != nil {
		log.Fatalln("Unable to read dataset.json", err)
	}
	json.Unmarshal(file, &dataset)

	log.Println("Handle Support...")
	minSupport := 0.2
	headTable := handleSupport(dataset, minSupport)

	log.Println("Construct FP-Tree...")
	for _, val := range dataset {
		root.fpTreeAppend(headTable, val...)
	}

	log.Println("FP-Tree Mining...")
}
