package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

const kmPerDegreeLongitude float64 = 71.5
const kmPerDegreeLatitude float64 = 111

const fileNameNodes string = "Data/smallOSMNode.csv"
const fileNameWays string = "Data/smallOSMWayNode.csv"
const fileNameTags string = "Data/smallOSMWayTag.csv"

type OSM_NODE struct {
	NODE_ID int     `csv:"NODE_ID"`
	LON     float64 `csv:"LON"`
	LAT     float64 `csv:"LAT"`
}

type OSM_WAY_NODE struct {
	WAY_ID  int `csv:"WAY_ID"`
	NODE_ID int `csv:"NODE_ID"`
	SEQ_NR  int `csv:"SEQ_NR"`
}

type OSM_WAY_TAG struct {
	WAY_ID int    `csv:"WAY_ID"`
	KEY    string `csv:"KEY"`
	VALUE  string `csv:"VALUE"`
}

var Nodes []OSM_NODE = []OSM_NODE{}
var WayNodes []OSM_WAY_NODE = []OSM_WAY_NODE{}
var WayTags []OSM_WAY_TAG = []OSM_WAY_TAG{}

func main() {
	GetCSVData() // Retrieve data from CSV files

	minX, maxX, minY, maxY := GetMinMaxPoints() // Get "borders" and display them
	fmt.Printf("LON borders: [%f, %f]\nLAT borders: [%f, %f]\n", minX, maxX, minY, maxY)

	input := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Query a node ID (ex. 591147201): ")
		input.Scan()
		idQuery, err := strconv.Atoi(input.Text())
		if err == nil {
			id := slices.IndexFunc(Nodes, func(n OSM_NODE) bool { return n.NODE_ID == idQuery })
			if id == -1 {
				fmt.Println("No node for index", idQuery)
			} else {
				fmt.Println(id)
				fmt.Printf("Node %d: [LON: %f, LAT: %f]\n", idQuery, Nodes[id].LON, Nodes[id].LAT)
				fmt.Println("Ways related:")
				for i := range WayNodes {
					if WayNodes[i].NODE_ID == idQuery {
						fmt.Printf("Way(%d), Seq_nr(%d)\n", WayNodes[i].WAY_ID, WayNodes[i].SEQ_NR)
					}
				}
			}
		} else {
			fmt.Println("Please input a integer as index")
		}
	}
}

func GetMinMaxPoints() (minX, maxX, minY, maxY float64) {
	minX, maxX = Nodes[0].LON, Nodes[0].LON
	minY, maxY = Nodes[0].LAT, Nodes[0].LAT
	for _, v := range Nodes {
		if v.LON < minX {
			minX = v.LON
		}
		if v.LON > maxX {
			maxX = v.LON
		}
		if v.LAT < minY {
			minY = v.LAT
		}
		if v.LAT > maxY {
			maxY = v.LAT
		}
	}
	return
}

func GetCSVData() {
	getNodes()
	getWays()
	getTags()
}

func getNodes() {
	file, err := os.Open(fileNameNodes)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	skipFirstLine := true
	for scanner.Scan() {
		if !skipFirstLine {
			splited := strings.Split(scanner.Text(), ";")
			id, err := strconv.Atoi(splited[0])
			if err != nil {
				panic(err)
			}
			lon, err := strconv.ParseFloat(strings.Replace(splited[1], ",", ".", -1), 64)
			if err != nil {
				panic(err)
			}
			lat, err := strconv.ParseFloat(strings.Replace(splited[2], ",", ".", -1), 64)
			if err != nil {
				panic(err)
			}
			Nodes = append(Nodes, OSM_NODE{id, lon, lat})
		} else {
			skipFirstLine = false // First line skipped
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func getWays() {
	file, err := os.Open(fileNameWays)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	skipFirstLine := true
	for scanner.Scan() {
		if !skipFirstLine {
			splited := strings.Split(scanner.Text(), ";")
			way_id, err := strconv.Atoi(splited[0])
			if err != nil {
				panic(err)
			}
			node_id, err := strconv.Atoi(splited[1])
			if err != nil {
				panic(err)
			}
			seq_nr, err := strconv.Atoi(splited[2])
			if err != nil {
				panic(err)
			}
			WayNodes = append(WayNodes, OSM_WAY_NODE{way_id, node_id, seq_nr})
		} else {
			skipFirstLine = false // First line skipped
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func getTags() {
	file, err := os.Open(fileNameTags)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	skipFirstLine := true
	for scanner.Scan() {
		if !skipFirstLine {
			splited := strings.Split(scanner.Text(), ";")
			id, err := strconv.Atoi(splited[0])
			if err != nil {
				panic(err)
			}
			key := splited[1][1 : len(splited[1])-1]   // Remove the double quotes
			value := splited[2][1 : len(splited[2])-1] // Remove the double quotes
			WayTags = append(WayTags, OSM_WAY_TAG{id, key, value})
		} else {
			skipFirstLine = false // First line skipped
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
