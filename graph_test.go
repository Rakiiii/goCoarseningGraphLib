package coarseninggraph

import(
	"testing"
	"fmt"
	"log"
	gopair "github.com/Rakiiii/goPair"
)

func TestGetGraphWithOutEdge(t *testing.T){
	
	fmt.Println("Start TestGetGraphWithOutEdge")
	var graph Graph
	if err := graph.ParseGraph("testgraph"); err != nil {
		log.Println(err)
		return
	}

	edges := make([]gopair.IntPair,4)

	edges[0].First = 0
	edges[0].Second = 1

	edges[1].First = 2
	edges[1].Second = 4

	edges[2].First = 2
	edges[2].Second = 5

	edges[3].First = 8
	edges[3].Second = 6

	var graph2 Graph
	if err := graph2.ParseGraph("testgraphed"); err != nil {
		log.Println(err)
		return
	}

	newGraph := graph.GetGraphWithOutEdge(edges...)

	checkFlag := true

	for i := 0; i < newGraph.AmountOfVertex(); i++{
		for j,e := range newGraph.GetEdges(i){
			if e != graph2.GetEdges(i)[j]{
				t.Error("Wrong edge:",e," expected:",graph2.GetEdges(i)[j])
				checkFlag = false
			}
		}
	}
	if checkFlag{
		fmt.Println("TestGetGraphWithOutEdge=[ok]")
	}
}

func TestGetCoarseningGraph( t *testing.T){
	fmt.Println("Start TestGetCoarseningGraph")
	var graphnc Graph
	if err := graphnc.ParseGraph("testgraphnc"); err != nil {
		log.Println(err)
		return
	}
	var graphc Graph
	if err := graphc.ParseGraph("testgraphc"); err != nil {
		log.Println(err)
		return
	}

	graphnc.HungryNumIndependent()

	checkFlag := true

	rightorder := []int{2,3,0,1,5,6,7,4}
	graphnc.RenumVertex(rightorder)

	testGraph,ord := graphnc.GetCoarseningGraph(2,nil)
	rightorder = []int{0,1,2,3,6,4,5,7}
	for i,v := range ord{
		if rightorder[i] != v{
			t.Error("Wrong order:",v," expexted:",rightorder[i])
			checkFlag = false
		}
	}

	

	for i := 0 ; i < testGraph.AmountOfVertex(); i++{
		for j,e := range testGraph.GetEdges(i){
			if e != graphc.GetEdges(i)[j]{
				t.Error("Wrong edge:",e," expected:",graphc.GetEdges(i)[j])
				checkFlag = false
			}
		}
	}

	if checkFlag{
		fmt.Println("TestGetCoarseningGraph=[ok]")
	}
}