package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

const kmPerDegreeLongitude float64 = 71.5
const kmPerDegreeLatitude float64 = 111

var maxDistanceKM float64 = 5.0
var maxDistanceTest float64 = 0.006

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

// Structured list of Ways
var Ways map[int]Way = map[int]Way{}

type Way struct {
	properties map[string]string
	list       []OSM_NODE
}

// List of id of way that matches the query
var QueriedWaysId []int = []int{}

func main() {
	GetCSVData() // Retrieve data from CSV files

	BuildWays() // Build a list of ways that is a list of Node

	minX, maxX, minY, maxY := GetMinMaxPoints() // Get "borders" and display them
	fmt.Printf("LON borders: [%f, %f]\nLAT borders: [%f, %f]\n", minX, maxX, minY, maxY)

	go scannerLoop() //Infinite loop of node id input

	openWindow() // Display the map with some graphics
}

func BuildWays() {
	// 1) Grouping all nodes together to form ways
	for _, v := range WayNodes {
		_, ok := Ways[v.WAY_ID]
		if !ok {
			key := ""
			value := ""
			if id := getIndexWays(v.WAY_ID); id != -1 {
				key = WayTags[id].KEY
				value = WayTags[id].VALUE
			}
			Ways[v.WAY_ID] = Way{
				map[string]string{
					key: value,
				},
				[]OSM_NODE{},
			}
		}
		if id := getIndexNodes(v.NODE_ID); id != -1 {
			curWay, _ := Ways[v.WAY_ID]
			key := ""
			value := ""
			if id := getIndexWays(v.WAY_ID); id != -1 {
				key = WayTags[id].KEY
				value = WayTags[id].VALUE
			}
			curWay.list = append(curWay.list, Nodes[id])
			curWay.properties[key] = value
			Ways[v.WAY_ID] = curWay
		}
	}
}

func getIndexNodes(id int) int {
	for p, v := range Nodes {
		if v.NODE_ID == id {
			return p
		}
	}
	return -1
}

func getIndexWays(id int) int {
	for p, v := range WayTags {
		if v.WAY_ID == id {
			return p
		}
	}
	return -1
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

func scannerLoop() {
	QueriedWaysId = []int{}
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
				fmt.Println("Ways containing this node:")
				for i := range WayNodes {
					if WayNodes[i].NODE_ID == idQuery {
						fmt.Printf("Way(%d), Seq_nr(%d)\n", WayNodes[i].WAY_ID, WayNodes[i].SEQ_NR)
					}
				}
			}

			fmt.Println("\nList of ways close to this node:")
			for key, way := range Ways {
				for id, _ := range way.list {
					if id < len(way.list)-1 {
						// Coordinates of center of way
						wayX := min(way.list[id].LON, way.list[id+1].LON) + math.Abs(way.list[id].LON-way.list[id+1].LON)
						wayY := min(way.list[id].LAT, way.list[id+1].LAT) + math.Abs(way.list[id].LAT-way.list[id+1].LAT)
						// fmt.Println(wayX, wayY)
						// Distance between the query point and the center of way
						dist := math.Sqrt(
							math.Pow(math.Abs(wayX-Nodes[id].LON), 2.0) +
								math.Pow(math.Abs(wayY-Nodes[id].LAT), 2.0))
						if dist < maxDistanceTest {
							// fmt.Printf("Way(%d)\n", key)
							QueriedWaysId = append(QueriedWaysId, key)
						}
					}
				}
			}

		} else {
			fmt.Println("Please input an integer as index")
		}
	}
}

func min(a, b float64) float64 {
	if a > b {
		return b
	} else {
		return a
	}
}
