package coarseninggraph 

import (
	lspartitioninglib "github.com/Rakiiii/goBipartitonLocalSearch"
	gosort "github.com/Rakiiii/goSort"
	gopair "github.com/Rakiiii/goPair"
	gotuple "github.com/Rakiiii/goTuple"
	"math"	
	//"fmt"
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

//contractVertex returns graph with vertex from @set connected and matrix of reconnecting vertex
//@set must have such structure: first cord is number of vertex in previous graph,
//and this slice contains group of vertex for connections to vertex with
//number first cord, if no vertex must be connected to vertex slice dhould contain only this vertex number
//if vertex must be connected to another vertex slice must contain number of vertex to which it must be connected
func (g *Graph)contractVertex(set [][]int)(*Graph,[][]int){
	n := 0

	//count amount of vertex that will be contarcted
	for i,j := range set{
		if j[0] == i{
			n++
		}
	}


	//init slice for recontacting vertex
	fixed := make([][]int,n)
	//init slice for reordering vertex
	ord := make([]int,g.AmountOfVertex())
	//itterator for right pos in @fixed
	it := 0

	//look whole @set
	//remove vertex which will be contarcted
	for i,j := range set{
		//if vertex poininting on it self, then other vertex will be contacted to it
		if j[0] == i{
			//we must save it
			fixed[it] = j
			//it must be have its own number
			ord[i] = it
			//move itterator
			it++
		}else{
			//if vertex will be contarcted with another one, it will have non personal number later
			ord[i] = -1
		}
	} 

	//set for contracted vertex thei numbers
	//they must have the same number as vertex to which they will be contarcted
	for i,j := range ord{
		if j == -1{
			ord[i] = ord[recursiveCheck(set,i)]
		}
	}


	//init graph with contracted vertex
	newGraph := new(Graph)
	newGraph.Init(len(fixed),0)

	//add edges to this grpah
	for i,j := range fixed{
		//init slice of edges
		edges := make([]int,0)

		//check all edges of vertex that will be contracted
		for _,v := range j{
		//edges of vertex @v in source graph
		sourceEdges := g.GetEdges(v)
		aped := make([]int,0)
		//remove all edges with contacted vertex
		for _,e := range sourceEdges{
			if !isContains(j,e){
				aped = append(aped,e)
			}
		}

		//renum vertex of prev graph with new order
		for i1,j1 := range aped{
			aped[i1] = ord[j1]
		}

		//remove vertex repeat
		aped = removeRepeat(aped)

		//add reworked edges of vertex @v from source graph to set of edges of vertex @i in new graph
		edges = gosort.QuicksortInt(appendWithOutRepeat(edges,aped))
		}

		//add edges to graph
		newGraph.AddEdgesToVertex(i,edges)
	}

	return newGraph,fixed
}

//GetHungryContractedGraphNI returns pointer to graph of type Graph that composed from set of contaracted vertex
//and matrix for uncotractiong vertex matrix strucutre:line number is num of vertex in new graph,and line it self contains 
//number of vertex of source grpah that composed this vertex
//@n is amount of vertex that will be contarcted
func (g *Graph)GetHungryContractedGraphNI(n int)(*Graph,[][]int){
	//constract slice of tuples of this struct:@s.First is number of first vertex ,@s.Second is number of second vertex,
	//@s.Third is size of edges overlap for vertex in tuple
	//it contains all vertex pairs
	pairSet := make([]gotuple.IntTuple,(g.AmountOfVertex()*(g.AmountOfVertex() - 1 ))/2)
	it := 0

	//fill set
	for fv := 0 ; fv < g.AmountOfVertex(); fv++{
		for sv := fv + 1 ; sv < g.AmountOfVertex();sv++{
			//count edges overlap
			counter := countSliceOverlap(g.GetEdges(fv),g.GetEdges(sv))
			pairSet[it].First = fv
			pairSet[it].Second = sv
			pairSet[it].Third = counter
			it++
		}
	}


	//sort from low to high to hungry work
	pairSet = gotuple.QuicksortIntTupleThird(pairSet)
	
	/*for _,i := range pairSet{
		fmt.Print("[",i.First,",",i.Second,"]:",i.Third)
	}*/


	//constract set for contracting vertex
	result := make([][]int,g.AmountOfVertex())
	it = 0
	for i,_ := range result{
		result[i] = make([]int,1)
		result[i][0] = i
	}

	//complite contract vertex in set
	for i := len(pairSet) - 1; i > len(pairSet)-n-1;i--{
		contractVertex(result,pairSet[i].First,pairSet[i].Second)
	}

	//constract graph with contracted vertex
	return g.contractVertex(result)
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

//HungryFixBipartitionDisbalance
func (g *Graph)HungryFixBipartitionDisbalance(vec []bool,groupSize int)[]bool{

	tg := 0
	fg := 0
	for _,b := range vec{
		if b{
			tg++
		}else{
			fg++
		}
	}

	var flag bool
	switch  {
	case tg>fg:
		if groupSize - tg >= 0{
			return vec
		}
		flag = true
	case tg<fg:
		if groupSize - fg >= 0{
			return vec
		}
		flag = false
	case fg == tg:
		return vec
	}

	newVec := make([]bool,len(vec))
	for i,j := range vec{
		newVec[i] = j
	}

	vertex := -1
	bestWeight := math.Inf(1)

		for i := 0; i < g.AmountOfVertex(); i++{
			if vec[i] == flag{
				inEdges := 0
				outEdges := 0
				for _,v := range g.GetEdges(i){
					if vec[i] == vec[v]{
						inEdges ++ 
					}else{
						outEdges++
					}
				}
				if bestWeight > float64(inEdges - outEdges){
					bestWeight = float64(inEdges - outEdges)
					vertex = i
				}
			}
		}

	newVec[vertex] = !newVec[vertex]

	return g.HungryFixBipartitionDisbalance(newVec,groupSize)
}

//UncontractedGraphBipartition takes @contr which is matrix that returns by any contarction method and bipartition vector
//returns bipartition vector for uncontracted graph
func UncontractedGraphBipartition(contr [][]int,vec []bool)[]bool{
	if len(contr) != len(vec){
		return nil
	}

	n := 0
	for _,j := range contr{
		for i := 0; i < len(j);i++{
			n++
		}
	}

	result := make([]bool,n)

	for i,r := range vec{
		for _,v := range contr[i]{
			result[v] = r
		}
	}

	return result
}

func countSliceOverlap(f []int,s []int)int{
	counter := 0
	for _,i := range s{
		for _,j := range f{
			if i == j{
				//log.Println("counter increased")
				counter++
			}
		}
	}
	return counter
}

func recursiveCheck(slice [][]int,check int)int{
	if slice[check][0] == check{
		return check
	}else{
		return recursiveCheck(slice,slice[check][0])
	}
}

func contractVertex(set [][]int,v1 int,v2 int){
	switch{
	case set[v1][0] == v1 && set[v2][0] == v2:
		set[v1] = appendWithOutRepeat(set[v1],set[v2])
		set[v2] = []int{v1}
	case set[v1][0] != v1 && set[v2][0] != v2:
		contractVertex(set,set[v1][0],set[v2][0])
	case set[v1][0] != v1 && set[v2][0] == v2 :
		contractVertex(set,set[v1][0],v2)
	case set[v1][0] == v1 && set[v2][0] != v2:
		contractVertex(set,v1,set[v2][0])
	}
}

func appendWithOutRepeat(s1 []int,s2 []int)[]int{
	res := make([]int,len(s1))
	copy(res,s1)
	for _,j := range s2{
		flag := true
		for _,i := range s1{
			if j == i{
				flag = false
			}
		}

		if flag {
			res = append(res,j)
		}
	}

	return res
}

func isContains(slice []int,a int)bool{
	for _,j := range slice{
		if j == a{
			return true
		}
	}
	return false
}

func removeRepeat(slice []int)[]int{
	newSlice := make([]int,0)
	for _,j := range slice{
		flag := true
		for _,i := range newSlice{
			if i == j{
				flag = false
			}
		}
		if flag{
			newSlice = append(newSlice,j)
		}
	}
	return newSlice
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