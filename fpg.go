package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sort"
)

type Node struct {
	Food   *Food
	Number int
	Parent *Node
	Sons   []*Node
}

type Food struct {
	Name  string
	Count int
	Links []*Node
}

type FoodList []Food

type Dataset map[int][]string

type FoodptrList []*Food

type SupportData []FoodptrList

func (list FoodList) Len() int {
	return len(list)
}

func (list FoodList) Less(i, j int) bool {
	return list[i].Count < list[j].Count
}

func (list FoodList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list FoodptrList) Len() int {
	return len(list)
}

func (list FoodptrList) Less(i, j int) bool {
	return list[i].Count < list[j].Count
}

func (list FoodptrList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list FoodList) index(name string) int {
	for i, link := range list {
		if link.Name == name {
			return i
		}
	}
	return -1
}

func fpTreeCreadeNode(food *Food, parent *Node) *Node {
	return &Node{food, 1, parent, []*Node{}}
}

func (rootp *Node) fpTreeInsert(list FoodList, food *Food) (node *Node) {
	for _, val := range rootp.Sons {
		if val.Food == food {
			node = val
			val.Number++
		}
	}
	if node == nil {
		node = fpTreeCreadeNode(food, rootp)
		rootp.Sons = append(rootp.Sons, node)
	}
	return
}

func (rootp *Node) fpTreeAppend(list FoodList, foods ...*Food) {
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
		if idx != -1 {
			ds[num] = append(orders[:idx], orders[idx+1:]...)
		}
	}
}

func handleSupport(ds Dataset, minSpt float64) (headTable FoodList, sptData SupportData) {
	headTable = make(FoodList, 0)
	sptData = make(SupportData, 0)
	tempTable := make(FoodList, 0)
	total := float64(len(ds))
	for _, orderList := range ds {
		for _, food := range orderList {
			idx := tempTable.index(food)
			if idx == -1 {
				tempTable = append(tempTable, Food{food, 1, []*Node{}})
			} else {
				link := &tempTable[idx]
				link.Count++
			}
		}
	}
	for _, link := range tempTable {
		if spt := float64(link.Count) / total; spt < minSpt {
			removeItem(ds, link.Name)
		} else {
			headTable = append(headTable, link)
		}
	}

	sort.Sort(sort.Reverse(headTable))
	for _, orders := range ds {
		temp := make(FoodptrList, 0)
		for _, food := range orders {
			temp = append(temp, &(headTable[headTable.index(food)]))
		}
		if temp.Len() > 0{
			sort.Sort(sort.Reverse(temp))
			sptData = append(sptData, temp)
		}
	}
	return
}

func main() {
	root := Node{Food: nil, Number: 0, Sons: []*Node{}}
	dataset := make(Dataset)
	file, err := ioutil.ReadFile("dataset.json")
	if err != nil {
		log.Fatalln("Unable to read dataset.json", err)
	}
	json.Unmarshal(file, &dataset)

	log.Println("Handle Support...")
	minSupport := 0.2
	headTable, sptDS := handleSupport(dataset, minSupport)

	log.Println("Construct FP-Tree...")
	for _, val := range sptDS {
		root.fpTreeAppend(headTable, val...)
	}

	log.Println("FP-Tree Mining...")
}
