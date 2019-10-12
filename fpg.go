package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"reflect"
	"sort"
	"strconv"
)

type Dataset map[int][]string

type Node struct {
	Item   *Item
	Number int
	Parent *Node
	Sons   []*Node
}

type Item struct {
	Name  string
	Count int
	Links []*Node
}

type ItemList []*Item

type SupportedData []ItemList

type Patterns struct {
	Support int
	Item    []string
}

func (list ItemList) Len() int {
	return len(list)
}

func (list ItemList) Less(i, j int) bool {
	return list[i].Count < list[j].Count
}

func (list ItemList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list ItemList) indexByName(name string) int {
	for i, link := range list {
		if link.Name == name {
			return i
		}
	}
	return -1
}

func fpTreeCreadeNode(item *Item, parent *Node, num int) *Node {
	return &Node{item, num, parent, []*Node{}}
}

func (rootp *Node) fpTreeLeafPath(paths *[]ItemList, curPath ...*Item) {
	for _, node := range rootp.Sons {
		node.fpTreeLeafPath(paths, append(curPath, rootp.Item)...)
	}
	if len(rootp.Sons) == 0 {
		*paths = append(*paths, append(ItemList{}, curPath...))
	}
}

func (rootp *Node) fpTreeInsert(list ItemList, item *Item, count int) (node *Node) {
	for _, val := range rootp.Sons {
		if val.Item == item {
			node = val
			val.Number += count
		}
	}
	if node == nil {
		node = fpTreeCreadeNode(item, rootp, count)
		rootp.Sons = append(rootp.Sons, node)
		item.Links = append(item.Links, node)
	}
	return
}

func (rootp *Node) fpTreeAppend(list ItemList, count bool, items ...*Item) {
	cur := rootp
	for _, item := range items {
		if count {
			cur = cur.fpTreeInsert(list, item, 1)
		} else {
			cur = cur.fpTreeInsert(list, item, item.Count)
		}
	}
}

func removeItem(ds Dataset, itemName string) {
	for num, orders := range ds {
		idx := -1
		for i, v := range orders {
			if v == itemName {
				idx = i
				break
			}
		}
		if idx != -1 {
			ds[num] = append(orders[:idx], orders[idx+1:]...)
		}
	}
}

func handleSupport(ds Dataset, spt float64) (headerTable ItemList, sptData SupportedData, total int, minSptc int) {
	headerTable = make(ItemList, 0)
	sptData = make(SupportedData, 0)
	tempTable := make(ItemList, 0)
	total = len(ds)
	minSptc = int(float64(total) * spt)
	for _, orderList := range ds {
		for _, item := range orderList {
			idx := tempTable.indexByName(item)
			if idx == -1 {
				tempTable = append(tempTable, &Item{item, 1, []*Node{}})
			} else {
				item := tempTable[idx]
				item.Count++
			}
		}
	}
	for _, link := range tempTable {
		if link.Count < minSptc {
			removeItem(ds, link.Name)
		} else {
			headerTable = append(headerTable, link)
		}
	}

	sort.Sort(headerTable)
	for _, orders := range ds {
		temp := make(ItemList, 0)
		for _, item := range orders {
			temp = append(temp, headerTable[headerTable.indexByName(item)])
		}
		if temp.Len() > 0 {
			sort.Sort(sort.Reverse(temp))
			sptData = append(sptData, temp)
		}
	}
	return
}

func freqtPatn(list ItemList, item string) (patterns []Patterns) {
	_len := len(list)
	n := int(math.Pow(2.0, float64(_len)))
	for i := 1; i < n; i++ {
		bin := strconv.FormatInt(int64(i), 2)
		for zero := _len - len(bin); zero > 0; zero-- {
			bin = "0" + bin
		}
		temp := make([]string, 0)
		minSup := -1
		for j, b := range bin {
			if string(b) == "1" {
				if minSup == -1 || list[j].Count < minSup {
					minSup = list[j].Count
				}
				temp = append(temp, list[j].Name)
			}
		}
		temp = append(temp, item)
		patterns = append(patterns, Patterns{
			Support: minSup,
			Item:    temp,
		})
	}
	return patterns
}

func mining(headTable ItemList, minSup int) (patterns []Patterns) {
	patterns = make([]Patterns, 0)
	for _, item := range headTable {
		temp := Node{Item: nil, Number: 0, Sons: []*Node{}}
		pattern := make([]ItemList, 0)
		for _, link := range item.Links {
			nodes := ItemList{link.Item}
			for link.Parent != nil {
				link = link.Parent
				if link.Item != nil {
					nodes = append(append(ItemList{}, link.Item), nodes...)
				}
			}
			temp.fpTreeAppend(ItemList{}, false, nodes...)
		}
		temp.fpTreeLeafPath(&pattern)

		// handle minSupport
		for i, pat := range pattern {
			idx := len(pat)
			for j, node := range pat[1:] {
				if node.Count < minSup {
					idx = j
					break
				}
			}
			pattern[i] = pat[1:idx]
			for _, p := range freqtPatn(pattern[i], item.Name) {
				addFlag := true
				for k, _ := range patterns {
					if reflect.DeepEqual(patterns[k].Item, p.Item) {
						patterns[k].Support += p.Support
						addFlag = false
						break
					}
				}
				if addFlag {
					patterns = append(patterns, p)
				}
			}
		}
	}
	return
}

func main() {
	root := Node{Item: nil, Number: 0, Sons: []*Node{}}
	dataset := make(Dataset)
	file, err := ioutil.ReadFile("dataset.json")
	if err != nil {
		log.Fatalln("Unable to read dataset.json", err)
	}
	json.Unmarshal(file, &dataset)

	log.Println("Handle Support...")
	minSupport := 0.2
	headTable, sptDS, _, minSup := handleSupport(dataset, minSupport)

	log.Println("Construct FP-Tree...")
	for _, val := range sptDS {
		root.fpTreeAppend(headTable, true, val...)
	}

	log.Println("FP-Tree Mining...")
	patterns := mining(headTable, minSup)

	log.Println("Write File...")
	jbyte, err := json.MarshalIndent(patterns, "    ", "    ")
	ioutil.WriteFile("result.json", jbyte, 644)
}
