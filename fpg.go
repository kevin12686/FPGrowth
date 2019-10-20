package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
)

type Item struct {
	AttrName string
	Value    string
	_weight  int
}

type ItemList []*Item

type Data struct {
	Tid          string
	Items        ItemList
	SupportCount int
}

type Dataset []*Data

type Header struct {
	ItemPtr      *Item
	SupportCount int
	Nodes        FPNodes
}

type HeaderTable []*Header

type FPNode struct {
	ItemPtr      *Item
	Parent       *FPNode
	Sons         FPNodes
	SupportCount int
}

type FPNodes []*FPNode

type Rule struct {
	Base       *Data
	Candidate  *Data
	Confidence float64
}

type Rules []*Rule

func (list ItemList) Len() int {
	return len(list)
}

func (list ItemList) Less(i, j int) bool {
	return list[i]._weight < list[j]._weight
}

func (list ItemList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list ItemList) Find(attrName string, value string) *Item {
	for _, itemPtr := range list {
		if itemPtr.AttrName == attrName && itemPtr.Value == value {
			return itemPtr
		}
	}
	return nil
}

func (list ItemList) EqualTo(otherList ItemList) bool {
	if length := len(list); length == len(otherList) {
		tempList := append(ItemList{}, list...)
		for _, itemPtr := range otherList {
			found := false
			for _, tempPtr := range tempList {
				if itemPtr == tempPtr {
					found = true
					break
				}
			}
			if !found {
				tempList = append(tempList, itemPtr)
			}
		}
		if length == len(tempList) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (list ItemList) SupportCount(dataset Dataset) (support int) {
	length := len(list)
	for _, data := range dataset {
		sum := 0
		for _, itemPtr := range data.Items {
			for _, confItem := range list {
				if itemPtr == confItem {
					sum++
				}
			}
			if sum == length {
				support++
				break
			}
		}
	}
	return
}

func (list Dataset) Len() int {
	return len(list)
}

func (list Dataset) Less(i, j int) bool {
	return list[i].SupportCount < list[j].SupportCount
}

func (list Dataset) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list HeaderTable) Len() int {
	return len(list)
}

func (list HeaderTable) Less(i, j int) bool {
	return list[i].SupportCount < list[j].SupportCount
}

func (list HeaderTable) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list HeaderTable) Find(itemPtr *Item) *Header {
	for _, header := range list {
		if header.ItemPtr == itemPtr {
			return header
		}
	}
	return nil
}

func (list HeaderTable) IndexOf(itemPtr *Item) int {
	for i, header := range list {
		if header.ItemPtr == itemPtr {
			return i
		}
	}
	return -1
}

func (node *FPNode) Insert(itemPtr *Item, supportCount int) (fpnode *FPNode, created bool) {
	found := false
	for _, son := range node.Sons {
		if son.ItemPtr == itemPtr {
			found = true
			son.SupportCount += supportCount
			fpnode = son
			break
		}
	}
	if !found {
		fpnode = &FPNode{ItemPtr: itemPtr, Parent: node, SupportCount: supportCount}
		node.Sons = append(node.Sons, fpnode)
	}
	created = !found
	return
}

func (node *FPNode) Prefix() (list ItemList, supportCount int) {
	supportCount = node.SupportCount
	for cur := node; cur != nil && cur.ItemPtr != nil; cur = cur.Parent {
		if node != cur {
			list = append(ItemList{cur.ItemPtr}, list...)
		}
	}
	return
}

func (list Rules) Len() int {
	return len(list)
}

func (list Rules) Less(i, j int) bool {
	return list[i].Confidence < list[j].Confidence
}

func (list Rules) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list Rules) IndexOf(rule *Rule) int {
	for i, r := range list {
		if r.Confidence != rule.Confidence || r.Base.SupportCount != rule.Base.SupportCount || r.Candidate.SupportCount != rule.Candidate.SupportCount {
			continue
		} else {
			if r.Base.Items.EqualTo(rule.Base.Items) && r.Candidate.Items.EqualTo(rule.Candidate.Items) {
				return i
			}
		}
	}
	return -1
}

func readData(filename string) (dataset Dataset, dataSize int, itemList ItemList) {
	file, _ := os.Open(filename)
	defer file.Close()
	csvReader := csv.NewReader(file)
	csvData, _ := csvReader.ReadAll()
	attrs := csvData[0]
	for _, record := range csvData[1:] {
		data := &Data{Tid: record[0], SupportCount: 1}
		for i, attr := range record[1:] {
			itemPtr := itemList.Find(attrs[i+1], attr)
			if itemPtr == nil {
				itemPtr = &Item{AttrName: attrs[i+1], Value: attr}
				itemList = append(itemList, itemPtr)
			}
			data.Items = append(data.Items, itemPtr)
		}
		dataset = append(dataset, data)
	}
	dataSize = len(dataset)
	return
}

func constructFPTree(dataset Dataset, supportCount int) (fproot *FPNode, table HeaderTable) {
	tempTable := make(HeaderTable, 0)
	for _, data := range dataset {
		for _, itemPtr := range data.Items {
			if header := tempTable.Find(itemPtr); header == nil {
				tempTable = append(tempTable, &Header{ItemPtr: itemPtr, SupportCount: data.SupportCount,})
			} else {
				header.SupportCount += data.SupportCount
			}
		}
	}

	for _, header := range tempTable {
		if header.SupportCount >= supportCount {
			table = append(table, header)
		}
	}

	if len(table) > 0 {
		fproot = &FPNode{}

		sort.Sort(table)
		for _, header := range table {
			header.ItemPtr._weight = table.IndexOf(header.ItemPtr)
		}

		for _, data := range dataset {
			cur := fproot
			sort.Sort(sort.Reverse(data.Items))
			for _, itemPtr := range data.Items {
				if header := table.Find(itemPtr); header != nil {
					node, created := cur.Insert(itemPtr, data.SupportCount)
					if created {
						header.Nodes = append(header.Nodes, node)
					}
					cur = node
				}
			}
		}
	}

	return
}

func mineFPTree(table HeaderTable, minSupportCount int, prefix ItemList, frequentItemSet *Dataset) {
	for _, header := range table {
		newPrefix := append(append(ItemList{}, prefix...), header.ItemPtr)
		*frequentItemSet = append(*frequentItemSet, &Data{Items: newPrefix, SupportCount: header.SupportCount})
		conditionPatternBases := Dataset{}
		for _, node := range header.Nodes {
			pattern, count := node.Prefix()
			conditionPatternBases = append(conditionPatternBases, &Data{Items: pattern, SupportCount: count,})
		}
		conditionFPTree, conditionHeaderTable := constructFPTree(conditionPatternBases, minSupportCount)
		if conditionFPTree != nil && len(conditionHeaderTable) > 0 {
			mineFPTree(conditionHeaderTable, minSupportCount, newPrefix, frequentItemSet)
		}
	}
}

func generateRules(dataset Dataset, supportCount int, frequentItem ItemList, subset ItemList, minConfidence float64, rules *Rules) {
	if len(frequentItem) > 1 {
		for i, _ := range frequentItem {
			newFrequentItem := append(ItemList{}, frequentItem[:i]...)
			newFrequentItem = append(newFrequentItem, frequentItem[i+1:]...)
			newSubset := append(subset, frequentItem[i])
			newSupportCount := newFrequentItem.SupportCount(dataset)
			confidence := float64(supportCount) / float64(newSupportCount)
			if confidence >= minConfidence {
				rule := &Rule{Base: &Data{Items: newFrequentItem, SupportCount: newSupportCount},
					Candidate:  &Data{Items: newSubset, SupportCount: supportCount,},
					Confidence: confidence}
				if rules.IndexOf(rule) == -1 {
					*rules = append(*rules, rule)
				}
				generateRules(dataset, supportCount, newFrequentItem, newSubset, minConfidence, rules)
			}
		}
	}
}

func main() {
	minSupport := 0.7
	minConfidence := 0.935
	dataset, dataSize, _ := readData("zoo.csv")
	minSupportCount := int(float64(dataSize) * minSupport)
	_, headerTable := constructFPTree(dataset, minSupportCount)
	frequentItemSet := &Dataset{}
	mineFPTree(headerTable, minSupportCount, ItemList{}, frequentItemSet)
	sort.Sort(sort.Reverse(frequentItemSet))
	fmt.Printf("Data Size: %d\n", dataSize)
	fmt.Printf("Minimal Support: %.2f\n", minSupport)
	fmt.Printf("Minimal Support Count: %d\n", minSupportCount)
	fmt.Printf("Minimal Confidence: %.2f\n", minConfidence)
	fmt.Println("\nFrequent Itemset:")
	for i, data := range *frequentItemSet {
		fmt.Print(i+1, "\t")
		for _, itemPtr := range data.Items {
			fmt.Printf("%s=%s ", itemPtr.AttrName, itemPtr.Value)
		}
		fmt.Printf("(Support Count: %d)\n", data.SupportCount)
	}

	rules := &Rules{}
	for _, frequentItem := range *frequentItemSet {
		generateRules(dataset, frequentItem.SupportCount, frequentItem.Items, ItemList{}, minConfidence, rules)
	}
	sort.Sort(sort.Reverse(rules))
	fmt.Println("\nRules:")
	for i, rule := range *rules {
		fmt.Printf("%d\t", i+1)
		for _, itemPtr := range rule.Base.Items {
			fmt.Printf("%s=%s ", itemPtr.AttrName, itemPtr.Value)
		}
		fmt.Printf("%d ==> ", rule.Base.SupportCount)
		for _, itemPtr := range rule.Candidate.Items {
			fmt.Printf("%s=%s ", itemPtr.AttrName, itemPtr.Value)
		}
		fmt.Printf("%d (Confidence: %.2f)\n", rule.Candidate.SupportCount, rule.Confidence)
	}
}
