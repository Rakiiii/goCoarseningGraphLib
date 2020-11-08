package coarseninggraph

import (
	"fmt"
	"log"
	"os"
	"testing"

	bipartitonlocalsearchlib "github.com/Rakiiii/goBipartitonLocalSearch"
	gopair "github.com/Rakiiii/goPair"
	permatchlib "github.com/Rakiiii/goPerfectMathcingLib"
	gotuple "github.com/Rakiiii/goTuple"
)

const (
	testdir    = "Testgraphs"
	benchGraph = "graph_bench_1"
)

func TestDir(t *testing.T) {
	fmt.Println("Start TestDir")
	if err := os.Chdir(testdir); err != nil {
		t.Error("Directory for graphs is not found")
		return
	} else {
		fmt.Println("TestDir=[ok]")
	}
}

func TestGetGraphWithOutEdge(t *testing.T) {

	fmt.Println("Start TestGetGraphWithOutEdge")

	graph := *NewGraphAny()
	//testgraph
	if err := graph.ParseGraph("GetGraphWithOutEdge"); err != nil {
		log.Println(err)
		return
	}

	edges := make([]gopair.IntPair, 4)

	edges[0].First = 0
	edges[0].Second = 1

	edges[1].First = 2
	edges[1].Second = 4

	edges[2].First = 2
	edges[2].Second = 5

	edges[3].First = 8
	edges[3].Second = 6

	graph2 := NewGraphAny()
	//testgraphed
	if err := graph2.ParseGraph("GetGraphWithOutEdgeResult"); err != nil {
		log.Println(err)
		return
	}

	newGraph := graph.GetGraphWithOutEdge(edges...)

	checkFlag := true

	for i := 0; i < newGraph.AmountOfVertex(); i++ {
		for j, e := range newGraph.GetEdges(i) {
			if e != graph2.GetEdges(i)[j] {
				t.Error("Wrong edge:", e, " expected:", graph2.GetEdges(i)[j])
				checkFlag = false
			}
		}
	}
	if checkFlag {
		fmt.Println("TestGetGraphWithOutEdge=[ok]")
	}
}

func TestGetCoarseningGraph(t *testing.T) {
	fmt.Println("Start TestGetCoarseningGraph")
	graphnc := NewGraphAny()
	//testgraphnc
	if err := graphnc.ParseGraph("GetCoarseningGraph"); err != nil {
		log.Println(err)
		return
	}
	graphc := NewGraphAny()
	//testgraphc
	if err := graphc.ParseGraph("GetCoarseningGraphResult"); err != nil {
		log.Println(err)
		return
	}

	graphnc.HungryNumIndependent()

	checkFlag := true

	rightorder := []int{2, 3, 0, 1, 5, 6, 7, 4}
	graphnc.RenumVertex(rightorder)

	testGraph, ord := graphnc.GetCoarseningGraph(2, nil)
	rightorder = []int{0, 1, 2, 3, 6, 4, 5, 7}
	for i, v := range ord {
		if rightorder[i] != v {
			t.Error("Wrong order:", v, " expexted:", rightorder[i])
			checkFlag = false
		}
	}

	for i := 0; i < testGraph.AmountOfVertex(); i++ {
		for j, e := range testGraph.GetEdges(i) {
			if e != graphc.GetEdges(i)[j] {
				t.Error("Wrong edge:", e, " expected:", graphc.GetEdges(i)[j])
				checkFlag = false
			}
		}
	}

	if checkFlag {
		fmt.Println("TestGetCoarseningGraph=[ok]")
	}
}

func TestAppenWithOutRepeat(t *testing.T) {
	fmt.Println("Start TestAppendWithOutRepeat")

	slice1 := []int{1, 2, 3, 4}
	slice2 := []int{3, 4, 1, 6, 5, 4}
	right := []int{1, 2, 3, 4, 6, 5}

	res := appendWithOutRepeat(slice1, slice2)

	checkFlag := true

	for i, j := range res {
		if j != right[i] {
			t.Error("Wrong value:", j, " at position:", i, " expected:", right[i])
			checkFlag = false
		}
	}

	if checkFlag {
		fmt.Println("TestAppendWithOutRepeat=[ok]")
	}
}

func TestCountSliceOverlap(t *testing.T) {
	fmt.Println("Start TestCountSliceOverlap")

	slice1 := []int{1, 2, 4, 5}
	slice2 := []int{4, 6, 8, 12, 1}

	if countSliceOverlap(slice1, slice2) != 2 {
		t.Error("Wrong size of slice overlap:", countSliceOverlap(slice1, slice2), " expected:2")
	} else {
		fmt.Println("TestCountSliceOverlap=[ok]")
	}
}

func TestRecursiveCheck(t *testing.T) {
	fmt.Println("Test TestRecursiveCheck")
	slice := [][]int{{0, 2, 3}, {2}, {4}, {3, 4, 5}, {0}}

	if recursiveCheck(slice, 1) != 0 {
		t.Error("Wrong recursive check:", recursiveCheck(slice, 1), " epected:0")
	} else {
		fmt.Println("TestRecursiveCheck=[ok]")
	}
}

func TestContractVertex(t *testing.T) {
	fmt.Println("Start TestContractVertex")
	slice := [][]int{{0, 2, 3}, {1, 2, 4}, {4}, {5}, {4, 6}, {5, 7}, {6, 9}, {8}, {8, 10}, {10}, {10, 12}, {11, 13}}
	resultSlice := [][]int{{0, 2, 3, 1, 4}, {0}, {4}, {5}, {4, 6, 5, 7}, {4}, {6, 9, 8, 10}, {8}, {6}, {10}, {10, 12, 11, 13}, {10}}

	contractVertex(slice, 0, 1)
	contractVertex(slice, 2, 3)
	contractVertex(slice, 6, 7)
	contractVertex(slice, 9, 11)

	checkFlag := true
	for i1, j1 := range slice {
		for i2, j2 := range j1 {
			if j2 != resultSlice[i1][i2] {
				t.Error("Wrong value:", j2, " at position:[", i1, "][", i2, "] expected:", resultSlice[i1][i2])
				checkFlag = false
			}
		}
	}

	if checkFlag {
		fmt.Println("TestContractVertex=[ok]")
	}
}

func TestIsContains(t *testing.T) {
	fmt.Println("Start TestIsContains")

	slice := []int{1, 2, 3, 4, 5, 6, 2}

	switch {
	case !isContains(slice, 2):
		t.Error("slice contains 2")
	case !isContains(slice, 3):
		t.Error("slice contains 3")
	case !isContains(slice, 6):
		t.Error("slice contains 6")
	case isContains(slice, 10):
		t.Error("slice doesn't contain 10")
	default:
		fmt.Println("TestIsContains=[ok]")
	}
}

func TestRemoveRepeat(t *testing.T) {
	fmt.Println("Start TestRemoveRepeat")

	slice := []int{1, 2, 4, 5, 1, 4, 2, 6}
	right := []int{1, 2, 4, 5, 6}

	slice = removeRepeat(slice)

	checkFlag := true
	for i, j := range slice {
		if j != right[i] {
			t.Error("Wrong value:", j, " at position:", i, " expected:", right[i])
			checkFlag = false
		}
	}
	if checkFlag {
		fmt.Println("TestRemoveRepeat=[ok]")
	}
}

func TestContractVertexGraph(t *testing.T) {
	fmt.Println("Start TestContractVertexGraph")

	graphStart := NewGraphAny()
	if err := graphStart.ParseGraph("ContractVertex"); err != nil {
		log.Println(err)
		return
	}
	graphRes := NewGraphAny()
	if err := graphRes.ParseGraph("ContractVertexResult"); err != nil {
		log.Println(err)
		return
	}

	ordRes := [][]int{{0, 6, 5}, {1}, {2, 3}, {4}, {7}, {8}, {9}}

	set := [][]int{{0, 6, 5}, {1}, {2, 3}, {2}, {4}, {0}, {0}, {7}, {8}, {9}}

	testGraph, testOrd := graphStart.contractVertex(set)

	checkFlag := true

	for i := 0; i < testGraph.AmountOfVertex(); i++ {
		for j, v := range testGraph.GetEdges(i) {
			if v != graphRes.GetEdges(i)[j] {
				t.Error("Wrong edge:(", i, ",", v, ") expected:(", i, ",", graphRes.GetEdges(i)[j], ")")
				checkFlag = false
			}
		}
	}

	for i, s := range testOrd {
		for j, v := range s {
			if v != ordRes[i][j] {
				t.Error("Wrong vertex subset:", v, " at position:[", i, ",", j, "] expected:", ordRes[i][j])
				checkFlag = false
			}
		}
	}

	if checkFlag {
		fmt.Println("ContractVertexGraph=[ok]")
	}
}

func TestGetHungryContractedGraphNI(t *testing.T) {
	fmt.Println("Start GetHungryContractedGraphNI")

	graphStart := NewGraphAny()
	if err := graphStart.ParseGraph("GetHungryContractedGraphNI"); err != nil {
		log.Println(err)
		return
	}
	graphRes := NewGraphAny()
	if err := graphRes.ParseGraph("GetHungryContractedGraphNIResult"); err != nil {
		log.Println(err)
		return
	}

	ordRes := [][]int{{0, 1}, {2}, {3}, {4, 5}, {6}, {7}}

	testGraph, testOrd := graphStart.GetHungryContractedGraphNI(2)

	checkFlag := true

	for i := 0; i < testGraph.AmountOfVertex(); i++ {
		for j, v := range testGraph.GetEdges(i) {
			if v != graphRes.GetEdges(i)[j] {
				t.Error("Wrong edge:(", i, ",", v, ") expected:(", i, ",", graphRes.GetEdges(i)[j], ")")
				checkFlag = false
			}
		}
	}

	for i, s := range testOrd {
		for j, v := range s {
			if v != ordRes[i][j] {
				t.Error("Wrong vertex subset:", v, " at position:[", i, ",", j, "] expected:", ordRes[i][j])
				checkFlag = false
			}
		}
	}

	if checkFlag {
		fmt.Println("GetHungryContractedGraphNI=[ok]")
	}
}

func TestUncontractedGraphBipartition(t *testing.T) {
	fmt.Println("Start TestUncontractedGraphBipartition")

	set := [][]int{{0, 1, 2}, {3, 4, 5}, {6}, {7}}
	testVec := []bool{true, false, true, false}
	resultVec := []bool{true, true, true, false, false, false, true, false}

	testVec = UncontractedGraphBipartition(set, testVec)

	checkFlag := true

	for i, j := range testVec {
		if j != resultVec[i] {
			t.Error("Wrong value:", j, " at position:", i, " expected:", resultVec[i])
			checkFlag = false
		}
	}

	if checkFlag {
		fmt.Println("TestUncontractedGraphBipartition=[ok]")
	}
}

func TestHungryFixBipartitionDisbalance(t *testing.T) {
	fmt.Println("Start TestHungryFixBipartitionDisbalance")

	graphStart := NewGraphAny()
	if err := graphStart.ParseGraph("HungryFixBipartitionDisbalance"); err != nil {
		log.Println(err)
		return
	}

	startVec := []bool{false, false, true, true, true, true, true, false}
	resultVec := []bool{false, false, true, true, false, true, true, false}

	testVec := graphStart.HungryFixBipartitionDisbalance(startVec, 4)

	checkFlag := true

	for i, j := range testVec {
		if j != resultVec[i] {
			t.Error("Wrong value:", j, " at position:", i, " expected:", resultVec[i])
			checkFlag = false
		}
	}

	if checkFlag {
		fmt.Println("TestHungryFixBipartitionDisbalance=[ok]")
	}
}

func TestCuatom(t *testing.T) {
	t.Skip()
	fmt.Println("Start TestCustom")
	graph := NewGraphAny()
	if err := graph.ParseGraph("../graph"); err != nil {
		log.Println(err)
		return
	}

	sol := &bipartitonlocalsearchlib.Solution{Value: -1, Vector: make([]bool, graph.AmountOfVertex()), Gr: graph.IGraph}
	sol.Init(graph.IGraph)

	depSize := 9

	fmt.Println("start mark serach with hungry ni contracted graph")
	contractgraph := Graph{IGraph: graph.IGraph}
	coarsedNumber := graph.AmountOfVertex() - graph.GetAmountOfIndependent() - depSize
	fmt.Println("coarsed number is:", coarsedNumber)
	if coarsedNumber > 0 {
		contractedGraph, contr := contractgraph.GetHungryContractedGraphI(coarsedNumber)

		groupSize := contractedGraph.AmountOfVertex()/2 + contractedGraph.AmountOfVertex()%2

		fmt.Println("groupSize:", groupSize)
		contractedGraph.Print()

		subOrd := contractedGraph.HungryNumIndependent()

		mark := bipartitonlocalsearchlib.LSPartiotionAlgorithmNonRec(contractedGraph.IGraph, nil, groupSize)

		if mark == nil {
			fmt.Println("mark is nil")
		}
		if subOrd == nil {
			fmt.Println("subOrd is nil")
		}

		mark.Vector = bipartitonlocalsearchlib.TranslateResultVector(mark.Vector, subOrd)

		fmt.Println("cont:", contr)

		//need to fix disbalance
		sol.Vector = contractgraph.HungryFixBipartitionDisbalance(UncontractedGraphBipartition(contr, mark.Vector), groupSize)
	}
}

func TestGetPerfectlyContractedGraph(t *testing.T) {
	fmt.Println("Start TestGetPerfectlyContractedGraph")

	graphStart := NewGraphAny()
	if err := graphStart.ParseGraph("GetHungryContractedGraphNI"); err != nil {
		log.Println(err)
		return
	}

	graphStart.Print()

	fixedVertexes := []gopair.IntPair{{First: 0, Second: 6}, {First: 3, Second: 4}}
	matcher := permatchlib.NewRandomMathcerWithFixedVertexes(fixedVertexes)
	graph, ord, err := graphStart.GetPerfectlyContractedGraph(matcher, matcher)
	if err != nil {
		fmt.Println(err)
		t.Error(err)
		return
	}

	fmt.Println("Contracted graph")
	graph.Print()
	fmt.Println("new order")
	for _, i := range ord {
		fmt.Print(i, " ")
	}

	fmt.Println("TestGetPerfectlyContractedGraph=[ok]")

}

func TestConutSliceDiff(t *testing.T) {
	fmt.Println("Start TestConutSliceDiff")

	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := []int{4, 5, 6, 7, 8}

	if countSliceDiff(slice1, slice2) != 6 {
		t.Error("Wrong slice differense size:", countSliceDiff(slice1, slice2), " expected:6")
	} else {
		fmt.Println("TestConutSliceDiff=[ok]")
	}
}

func TestCheckTupleSetContainment(t *testing.T) {
	fmt.Println("Start CheckTupleSetContainment")

	set := []gotuple.IntTuple{gotuple.IntTuple{First: 1, Second: 2, Third: 0}, gotuple.IntTuple{First: 3, Second: 4, Third: 0},
		gotuple.IntTuple{First: 1, Second: 3, Third: 0}, gotuple.IntTuple{First: 5, Second: 8, Third: 0}}

	checkFlag := true

	if !checkTupleSetContainment(set, 1, 2) {
		t.Error("Wrong result at [1,2]")
		checkFlag = false
	}

	if checkTupleSetContainment(set, 1, 5) {
		t.Error("Wrong result at [1,5]")
		checkFlag = false
	}

	if !checkTupleSetContainment(set, 4, 3) {
		t.Error("Wrong result at [4,3]")
		checkFlag = false
	}

	if checkTupleSetContainment(set, 3, 8) {
		t.Error("Wrong result at [3,8]")
		checkFlag = false
	}

	if checkFlag {
		fmt.Println("TestCheckTupleSetContainment=[ok]")
	}
}

func TestLinRegCoarsing(t *testing.T) {
	t.Skip()
	fmt.Println("start lin reg test")
	graphStart := NewGraphAny()
	if err := graphStart.ParseGraph("TestLinRegGraph"); err != nil {
		log.Println(err)
		return
	}

	_, _, err := graphStart.GetContractedWithLinRegGraph()
	if err != nil {
		log.Println(err)
	}
}
