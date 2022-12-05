package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
)

func main() {
	graph := buildGraph()
	city := os.Args[1]
	dijkstra(graph, city)

	for _, node := range graph.Nodes {
		fmt.Printf("O tempo mais curto de %s para %s é de %d minutos\n",
			city, node.name, node.value)
		for n := node; n.through != nil; n = n.through {
			fmt.Print(n, " <- ")
		}
		fmt.Println(city)
		fmt.Println()
	}
}

func dijkstra(graph *WeightedGraph, city string) {
	visited := make(map[string]bool)
	heap := &Heap{}

	startNode := graph.GetNode(city)
	startNode.value = 0
	heap.Push(startNode)

	for heap.Size() > 0 {
		current := heap.Pop()
		visited[current.name] = true
		edges := graph.Edges[current.name]
		for _, edge := range edges {
			if !visited[edge.node.name] {
				heap.Push(edge.node)
				if current.value+edge.weight < edge.node.value {
					edge.node.value = current.value + edge.weight
					edge.node.through = current
				}
			}
		}
	}
}

func buildGraph() *WeightedGraph {
	graph := NewGraph()
	nodes := AddNodes(graph,
		"Brasilia",
		"Belo-Horizonte",
		"Vitoria",
		"Rio-de-Janeiro",
		"Sao-Paulo",
		"Salvador",
		"Curitiba",
		"Porto-Alegre",
		"Florianopolis",
		"Campo-Grande",
	)
	graph.AddEdge(nodes["Brasilia"], nodes["Belo-Horizonte"], 594)
	graph.AddEdge(nodes["Brasilia"], nodes["Rio-de-Janeiro"], 975)
	graph.AddEdge(nodes["Brasilia"], nodes["Vitoria"], 1080)
	graph.AddEdge(nodes["Belo-Horizonte"], nodes["Rio-de-Janeiro"], 394)
	graph.AddEdge(nodes["Belo-Horizonte"], nodes["Salvador"], 1204)
	graph.AddEdge(nodes["Rio-de-Janeiro"], nodes["Curitiba"], 660)
	graph.AddEdge(nodes["Rio-de-Janeiro"], nodes["Sao-Paulo"], 327)
	graph.AddEdge(nodes["Rio-de-Janeiro"], nodes["Vitoria"], 428)
	graph.AddEdge(nodes["Sao-Paulo"], nodes["Porto-Alegre"], 945)
	graph.AddEdge(nodes["Sao-Paulo"], nodes["Salvador"], 1680)
	graph.AddEdge(nodes["Sao-Paulo"], nodes["Curitiba"], 344)
	graph.AddEdge(nodes["Curitiba"], nodes["Vitoria"], 1093)
	graph.AddEdge(nodes["Curitiba"], nodes["Porto-Alegre"], 627)
	graph.AddEdge(nodes["Porto-Alegre"], nodes["Salvador"], 2580)
	graph.AddEdge(nodes["Porto-Alegre"], nodes["Campo-Grande"], 1234)
	graph.AddEdge(nodes["Florianopolis"], nodes["Curitiba"], 333)
	graph.AddEdge(nodes["Florianopolis"], nodes["Campo-Grande"], 1080)
	graph.AddEdge(nodes["Salvador"], nodes["Campo-Grande"], 1920)

	return graph
}

type Node struct {
	name    string
	value   int
	through *Node
}

type Edge struct {
	node   *Node
	weight int
}

type WeightedGraph struct {
	Nodes []*Node
	Edges map[string][]*Edge
	mutex sync.RWMutex
}

func NewGraph() *WeightedGraph {
	return &WeightedGraph{
		Edges: make(map[string][]*Edge),
	}
}

func (g *WeightedGraph) GetNode(name string) (node *Node) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	for _, n := range g.Nodes {
		if n.name == name {
			node = n
		}
	}
	return
}

func (g *WeightedGraph) AddNode(n *Node) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.Nodes = append(g.Nodes, n)
}

func AddNodes(graph *WeightedGraph, names ...string) (nodes map[string]*Node) {
	nodes = make(map[string]*Node)
	for _, name := range names {
		n := &Node{name, math.MaxInt, nil}
		graph.AddNode(n)
		nodes[name] = n
	}
	return
}

func (g *WeightedGraph) AddEdge(n1, n2 *Node, weight int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.Edges[n1.name] = append(g.Edges[n1.name], &Edge{n2, weight})
	g.Edges[n2.name] = append(g.Edges[n2.name], &Edge{n1, weight})
}

func (n *Node) String() string {
	return n.name
}

func (e *Edge) String() string {
	return e.node.String() + "(" + strconv.Itoa(e.weight) + ")"
}

func (g *WeightedGraph) String() (s string) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	for _, n := range g.Nodes {
		s = s + n.String() + " ->"
		for _, c := range g.Edges[n.name] {
			s = s + " " + c.node.String() + " (" + strconv.Itoa(c.weight) + ")"
		}
		s = s + "\n"
	}
	return
}

type Heap struct {
	elements []*Node
	mutex    sync.RWMutex
}

func (h *Heap) Size() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.elements)
}

func (h *Heap) Push(element *Node) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.elements = append(h.elements, element)
	i := len(h.elements) - 1
	for ; h.elements[i].value < h.elements[parent(i)].value; i = parent(i) {
		h.swap(i, parent(i))
	}
}

func (h *Heap) Pop() (i *Node) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	i = h.elements[0]
	h.elements[0] = h.elements[len(h.elements)-1]
	h.elements = h.elements[:len(h.elements)-1]
	h.rearrange(0)
	return
}

func (h *Heap) rearrange(i int) {
	smallest := i
	left, right, size := leftChild(i), rightChild(i), len(h.elements)
	if left < size && h.elements[left].value < h.elements[smallest].value {
		smallest = left
	}
	if right < size && h.elements[right].value < h.elements[smallest].value {
		smallest = right
	}
	if smallest != i {
		h.swap(i, smallest)
		h.rearrange(smallest)
	}
}

func (h *Heap) swap(i, j int) {
	h.elements[i], h.elements[j] = h.elements[j], h.elements[i]
}

func parent(i int) int {
	return (i - 1) / 2
}

func leftChild(i int) int {
	return 2*i + 1
}

func rightChild(i int) int {
	return 2*i + 2
}

func (h *Heap) String() (str string) {
	return fmt.Sprintf("%q\n", getNames(h.elements))
}

func getNames(nodes []*Node) (names []string) {
	for _, node := range nodes {
		names = append(names, node.name)
	}
	return
}
