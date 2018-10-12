package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

/*
	Построить случайный граф ПО, для которого число выходящих из вер-шины ребер
	определяется датчиком случайных чисел. Определить зна-чение структурного параметра
	всех узлов /число висячих уз-лов. Привести гистограмму для полученных при построении
	графа зна-чений (m-1). Убедится в правильности построения гистограммы.
	В тексте программы и в отчете привести «утверждение», подтвер-ждающее правильность.
	Правило остановки построения графа (А или В) приведено в задании.
	Вариант: m = 3; N = 100; R = 100
*/

func init() {
	// trying to draw
	/*
		file, err := os.Create("draw.png")
		if err != nil {
			fmt.Println("some shit occurred")
			return
		}
		defer file.Close()
		graph := chart.Chart{
			Series: []chart.Series{
				chart.ContinuousSeries{
					XValues: []float64{1.0, 2.0, 3.0, 4.0},
					YValues: []float64{1.0, 2.0, 3.0, 4.0},
				},
			},
		}

		buffer := bytes.NewBuffer([]byte{})
		err = graph.Render(chart.PNG, buffer)
		file.Write(buffer.Bytes())
		if err != nil {}
	*/
}

func main() {

	var (
		m, N, R, exitCode int
		tempGraph         *graph
		grArr             []graph
		averageA          float64
	)
	// ввод данных
	fmt.Print("m : ")
	fmt.Scanf("%d", &m)
	fmt.Print("N : ")
	fmt.Scanf("%d", &N)
	fmt.Print("R : ")
	fmt.Scanf("%d", &R)

	grArr = make([]graph, 0, 0)
	for r := 0; r < R; r++ { // цикл, создающий R графов, выводит информацию о каждом
		tempGraph, exitCode = buildGraph(m, N)
		if exitCode == 1 {
			R++
		} else {
			grArr = append(grArr, *tempGraph)
			fmt.Println("№", r, tempGraph.getInfo())
			averageA += tempGraph.calcA()
		}
	}
	fmt.Println("Average a = ", averageA/float64(len(grArr)))

	// отрисовка первого графа
	fmt.Println(grArr[0].toString())

	// создание файла new.txt с матрицей смежности графа, для отрисовки
	grArr[0].printAdjacencyMatrix()

	// исследование схождения структурного параметра альфа:
	alphaArr := calcAfor(5, 200, 1, m)
	fmt.Println("Alpha convergence research:")
	for i := 0; i < len(alphaArr); i++ {
		fmt.Println(alphaArr[i])
	}

}

type graph struct {
	rootNode    node
	nodes       []node
	nodeCounter int
	hangNodes   int
}

type node struct {
	id             int
	hierarchyLevel int
	childs         []node
	parent         *node
}

func buildGraph(m, N int) (*graph, int) {
	var ( // создание объекта графа с корневым узлом
		gr = graph{
			node{
				1,
				0,
				make([]node, 0, 0),
				nil,
			},
			make([]node, 0, 0),
			1,
			0,
		}
		ran, beg, prevLen, addedToLvl int
		tempNode                      node
	)
	gr.nodes = append(gr.nodes, gr.rootNode)
	/*
		Цикл, генерирующий новый уровень иерархии до тех пор,
		пока количество узлов граца не будет > или = N
	*/
	for len(gr.nodes) < N {
		prevLen = len(gr.nodes)
		for i := beg; i < prevLen; i++ { // генерация детей для узлов текущего уровня иерархии
			rand.Seed(time.Now().UnixNano())
			ran = rand.Intn(m)
			for r := 0; r < ran; r++ { // генерация случайного кол-ва детей для i-того узла
				gr.nodeCounter++
				addedToLvl++
				tempNode = node{
					gr.nodeCounter,
					gr.nodes[i].hierarchyLevel + 1,
					make([]node, 0, 0),
					&gr.nodes[i],
				}
				gr.nodes = append(gr.nodes, tempNode)
				gr.nodes[i].childs = append(gr.nodes[i].childs, tempNode)
			}
		}
		/*
			Обрабатываем случай, когда на уровне иерархии все узлы сгенерированны
			без детей, завершение работы функции с кодом ошибки 1 (или 0 если достаточно узлов)
		*/
		if addedToLvl == 0 {
			if len(gr.nodes) < N {
				return &gr, 1
			} else {
				return &gr, 0
			}
		} else {
			addedToLvl = 0
		}
		beg = prevLen
	}
	// цикл подсчета висячих узлов
	for i := 0; i < len(gr.nodes); i++ {
		if len(gr.nodes[i].childs) == 0 {
			gr.hangNodes++
		}
	}
	return &gr, 0
}

func (n *node) toString() string {
	var s string
	s += "(" + strconv.Itoa(n.id)
	//s += strconv.Itoa(n.hierarchyLevel)
	if n.parent != nil {
		s += "-" + strconv.Itoa(n.parent.id)
	} else {
		s += "-null"
	}
	s += ")"
	return s
}

func (g *graph) toString() string {
	var (
		s   string
		lvl int
	)
	fmt.Print("0 ")
	for i := 0; i < len(g.nodes); i++ {
		if g.nodes[i].hierarchyLevel > lvl {
			lvl++
			s += "\n" + strconv.Itoa(g.nodes[i].hierarchyLevel) + " "
		}
		s += g.nodes[i].toString() + " "
	}
	return s
}

func (g *graph) getInfo() string {
	var s string
	s += "Nodes :" + strconv.Itoa(g.nodeCounter)
	s += " Hang :" + strconv.Itoa(g.hangNodes)
	s += " a = " + strconv.FormatFloat(g.calcA(), 'f', 6, 64)
	return s
}

func (g *graph) calcA() float64 {
	var a float64
	a = float64(g.nodeCounter) / float64(g.hangNodes)
	return a
}

// построение и вывод матрици смежности в файл new.txt
// http://graphonline.ru/create_graph_by_matrix
func (g *graph) printAdjacencyMatrix() {
	file, err := os.Create("out.txt")
	if err != nil {
		fmt.Println("some shit occurred")
		return
	}
	defer file.Close()
	matrix := make([][]int, len(g.nodes), len(g.nodes))

	for i := 0; i < len(g.nodes); i++ {
		matrix[i] = make([]int, len(g.nodes), len(g.nodes))
	}

	for i := 0; i < len(g.nodes); i++ {
		for j := 0; j < len(g.nodes[i].childs); j++ {
			matrix[i][g.nodes[i].childs[j].id-1] = 1
		}
	}

	for i := 0; i < len(g.nodes); i++ {
		var b bytes.Buffer
		for j := 0; j < len(matrix[i]); j++ {
			b.WriteString(strconv.Itoa(matrix[i][j]))
			if j == (len(matrix[i]) - 1) {
				b.WriteString("\n")
			} else {
				b.WriteString(", ")
			}
		}
		file.Write(b.Bytes())
	}

}

func calcAfor(begin, end, step, m int) (alphaArr []float64) {
	arr := make([]float64, 0, 0)
	var g *graph
	var err int
	for ; begin < end; begin += step {
		for {
			g, err = buildGraph(m, begin)
			if err == 0 {
				arr = append(arr, g.calcA())
				break
			}
		}
	}
	return arr
}
