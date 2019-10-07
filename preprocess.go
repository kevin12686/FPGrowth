package main

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
)

func preprocess(filenames ...string) map[int][]string {
	rlt := make(map[int][]string)
	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalln("Unable to open \""+filename+"\".", err)
		}
		csvReader := csv.NewReader(file)
		records, err := csvReader.ReadAll()
		if err != nil {
			log.Fatalln("Unable to read \""+filename+"\".", err)
		}
		for idx, val := range records {
			if idx == 0 || val[2] == "" {
				continue
			}
			orderNum, _ := strconv.Atoi(val[0])
			rlt[orderNum] = append(rlt[orderNum], val[2])
			sort.Strings(rlt[orderNum])
		}
	}
	return rlt
}

func main() {
	names := []string{"restaurant-1-orders.csv", "restaurant-2-orders.csv"}
	data := preprocess(names...)
	jbytes, err := json.MarshalIndent(data, "    ", "    ")
	if err != nil {
		log.Fatalln("Unable to dump json.", err)
	}
	ioutil.WriteFile("dataset.json", jbytes, 0644)
}
