package coarseninggraph

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	modelPath  = "model.pkl"
	resultPath = "contractedPairs"
	dataPath   = "pairs.csv"
	linergPath = "GetCoarsedVertex.py"
)

type data struct {
	Overlap  int
	Diff     int
	Incedent bool
}

var GlobalWriter *csv.Writer

func (d data) GetSlice() []string {
	res := make([]string, 3)
	res[0] = strconv.FormatFloat(float64(d.Overlap), 'f', -1, 64)
	res[1] = strconv.FormatFloat(float64(d.Diff), 'f', -1, 64)
	if d.Incedent {
		res[2] = "1.0"
	} else {
		res[2] = "0.0"
	}
	return res
}

func createRawResultFile() {
	f, err := os.Create(dataPath)
	if err != nil {
		log.Panic(err)
	} else {
		GlobalWriter = csv.NewWriter(f)
		GlobalWriter.Write([]string{"diff", "overlap", "incedent"})
	}
}

func IsIncedent(graph *Graph, fv, sv int) bool {
	for _, v := range graph.GetEdges(fv) {
		if v == sv {
			return true
		}
	}
	return false
}

func collectRawData(graph *Graph) {
	d := data{}
	for fv := 0; fv < graph.AmountOfVertex(); fv++ {
		for sv := fv + 1; sv < graph.AmountOfVertex(); sv++ {
			d.Diff = countSliceDiff(graph.GetEdges(fv), graph.GetEdges(sv))
			d.Overlap = countSliceOverlap(graph.GetEdges(fv), graph.GetEdges(sv))
			d.Incedent = IsIncedent(graph, fv, sv)
			GlobalWriter.Write(d.GetSlice())
		}
	}
	GlobalWriter.Flush()
}

func readCoarsing(av int) ([][]int, error) {
	file, err := os.Open(resultPath)
	if err != nil {
		return nil, err
	}

	//constract set for contracting vertex
	result := make([][]int, av)
	for i, _ := range result {
		result[i] = make([]int, 1)
		result[i][0] = i
	}

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	partition := strings.Fields(scanner.Text())

	pos := 0
	for i := 0; i < av; i++ {
		for j := i + 1; j < av; j++ {
			if partition[pos+(j-i-1)] == "1" {
				contractVertex(result, i, j)
			}
		}
		pos += av - i
	}
	return result, nil
}
