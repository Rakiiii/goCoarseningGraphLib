package coarseninggraph 

import (
	lspartitioninglib "github.com/Rakiiii/goBipartitonLocalSearch"
	gopair "github.com/Rakiiii/goPair"
	"math"	
)

type Graph struct{
	lspartitioninglib.Graph
	weightMatrix [][]float64
}

//SetEdgesWeight set weight for edges
func (g *Graph) SetEdgesWeight(newWeightMatrix [][]float64){
	g.weightMatrix = newWeightMatrix
}

//GetCoarseningGraph returns coarsening grpah with extra @n independent vertex
//@oldOrd must be nil for call
func (g *Graph)GetCoarseningGraph(n int,oldOrd []int)(*Graph,[]int){	
	if g.GetAmountOfIndependent() <= 0{
		return nil,nil
	}
	//if amount of extra independent vertex is 0 then stop recursion and return this graph
	if n <= 0 {
		return g,oldOrd
	}else{

		//lowest weight for vertex
		lowestWeight := math.MaxInt32
		//vertex with lowest weight
		vertex := -1

		//found vertex with lowest weight in dependent set
		//check all dependent vertex
		//[ graph vertex structure: ...setOfIndependentVertex (graph.GetAmountOfIndependent) setOfDependentVertex... ]
		for i := g.GetAmountOfIndependent();i < g.AmountOfVertex();i++{
			//get edges of vertex @i
			edges := g.GetEdges(i)
			//weigth of vertex
			weight := 0
			for _,edge := range edges{
				//weight depends on amount of edgew with independent vertex
				if edge < g.GetAmountOfIndependent(){
					weight ++
				}
			}
			//update lowestWweight and vertex
			if weight < lowestWeight && weight != 0{
				lowestWeight = weight
				vertex = i
			}
		}

		//if all vertex are independent return this graph
		if vertex == -1{
			return g,oldOrd
		}

		//create set of edges for deleting from graph
		edgesSet := make([]gopair.IntPair,lowestWeight)
		//number of edges in @edgesSet
		it := 0
		//check all edges
		for _,edge := range g.GetEdges(vertex){

			//we must delete edges with independent vertex
			//to increase amount of independent vertex 
			if edge < g.GetAmountOfIndependent(){
				edgesSet[it].First = vertex
				edgesSet[it].Second = edge
				it++
			}
		}
		//getting graph with out edges from edgesSet
		newGraph := g.GetGraphWithOutEdge(edgesSet...)

		//we need to update order of vertex,because we have new independent vertex
		//put new independent vertex to the end of independent set
		newOrder := make([]int,newGraph.AmountOfVertex())
		for i := 0 ; i < newGraph.AmountOfVertex(); i ++{
			switch{
			case i < g.GetAmountOfIndependent():
				newOrder[i] = i
			case i == g.GetAmountOfIndependent():
				newOrder[i] = vertex
				newOrder[i+1] = i
			case i != g.GetAmountOfIndependent() && i > g.GetAmountOfIndependent() && i < vertex:
				newOrder[i+1] = i	
			case i > vertex:
				newOrder[i] = i
			}
		}

		newGraph.RenumVertex(newOrder)
		//update amount of independent vertex
		newGraph.SetAmountOfIndependent(g.GetAmountOfIndependent()+1)

		if oldOrd == nil{
			oldOrd = newOrder
		}else{
			for i,v := range newOrder{
				newOrder[i] = oldOrd[v]
			}
		}
		//recursive call
		return newGraph.GetCoarseningGraph(n-1,append(oldOrd[:g.GetAmountOfIndependent()], newOrder[g.GetAmountOfIndependent():]...))
	}
}

//GetGraphWithOutEdge returns pointer to new graph that doesn't contain edges from @edgeSet
func (g *Graph)GetGraphWithOutEdge(edgeSet ...gopair.IntPair)*Graph{
	//init new void graph 
	var newGraph Graph
	newGraph.Init(g.AmountOfVertex(),0)

	//found edges to delete 
	//check all edges
	for i := 0;i < g.AmountOfVertex();i++{
		edges := g.GetEdges(i)
		for _,edge := range edgeSet{
			//delete edges
			switch i {
				case edge.First:
					for j,ed := range edges{
						if ed == edge.Second{
							edges = append(edges[:j],edges[j+1:]...)
							continue
						}
					}
				case edge.Second:
					for j,ed := range edges{
						if ed == edge.First{
							edges = append(edges[:j],edges[j+1:]...)
							continue
						}
					}
			}
		}
		newGraph.AddEdgesToVertex(i,edges)
	}
	return &newGraph
}
/*func (g *Graph)DoEdgesWeight()(ord []int){
	if g.AmountOfIndependent <= 0{
		ord = g.HungryNumIndependent()
	}else{
		for i := 0; i < g.AmountOfVertex();i++{
			ord = append(ord,i)
		}
	}

	g.weightMatrix = make([][]float64,g.AmountOfVertex())
	for i,_ := range g.weightMatrix{
		g.weightMatrix[i] = make([]float64,g.AmountOfVertex())
	}

	for i := 0; i < g.AmountOfVertex();i++{
		counter := 0
		for _,j := range g.GetEdges(i){
			if j < g.AmountOfIndependent{
				counter ++
			}
		}
		for _,j := range g.GetEdges(i){
			if j < g.AmountOfIndependent{
				weightMatrix[i][j] = 1/counter
			}
		}
	}
}*/