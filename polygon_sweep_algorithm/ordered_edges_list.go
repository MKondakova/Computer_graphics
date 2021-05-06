package main

import (
	"log"
	"math"
	"runtime"
	"sort"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const (
	BUILDING_POLYGON = 1
	POLYGON_BUILDED  = 2
	RASTERISATION    = 3
	FILTRATION       = 4
	SIZE             = 1000
)

type point struct {
	x float64
	y float64
}

var (
	mouse      point
	stage      int = BUILDING_POLYGON
	points     []point
	sizeX      int
	sizeY      int
	pixels     []uint8
	tempPixels []uint8
	list       map[int][]int
	edges      [][2]point
)

func makeEdges() {
	edges = make([][2]point, len(points))
	for i, p := range points {
		nextP := points[(i+1)%len(points)]
		edges[i] = [2]point{p, nextP}
	}
}

func isExtrema(y, y1, y2 float64) bool {
	return (y > y1 && y > y2) || (y < y1 && y < y2)
}

func vertexCountTwice(i, j int) bool {
	l := len(edges)
	return isExtrema(
		edges[i][j].y,
		edges[i][(j+1)%2].y,
		edges[(i-1+l)%l][j].y)
}

func addToList(x, y float64) {
	list[int(math.Floor(y))] = append(list[int(math.Floor(y))], int(math.Floor(x)))
}

func DDA() {
	for i, edge := range edges {
		if vertexCountTwice(i, 0) {
			addToList(edge[0].x, edge[0].y)
		}
		addToList(edge[1].x, edge[1].y)

		dy := edge[1].y - edge[0].y //разница между вершинами
		dx := edge[1].x - edge[0].x

		if dy == 0 {
			continue
		}

		count := int(math.Ceil(math.Abs(dy)))
		dy = dy / float64(count) //дельта отступа
		dx = dx / float64(count)

		checkEndFunc := func(i float64) bool {
			if dy > 0 {
				return edge[0].y+i*dy < edge[1].y
			} else {
				return edge[0].y+i*dy > edge[1].y
			}
		}

		for i := float64(1); checkEndFunc(i); i++ {
			addToList(edge[0].x+i*dx, edge[0].y+i*dy)
		}
	}
}

func drawLine(y, x1, x2 int) {
	for i := x1; i <= x2; i++ {
		pixels[(sizeY-y)*sizeX+i] = 255
	}
}

func fill() {
	for y := range list {
		sort.Ints(list[y])
		for i := 0; i < len(list[y]); i += 2 {
			drawLine(y, list[y][i], list[y][i+1])
		}
	}
}

func rasterisation() {
	pixels = make([]uint8, sizeY*sizeX)
	list = make(map[int][]int)
	makeEdges()
	DDA()
	fill()
}

func getNeighborsSum(i, j int) (int, int) {
	result := int(0)
	counter := int(0)
	for k := i - 1; k-(i-1) < 3; k++ {
		for l := j - 1; l-(j-1) < 3; l++ {
			if k >= 0 && k < sizeY && l >= 0 && l < sizeX {
				result += int(pixels[sizeX*k+l])
				counter++
			}
		}
	}
	result -= int(pixels[sizeX*i+j])
	counter--
	return result, counter
}

func filtrate() {
	tempPixels = make([]uint8, sizeY*sizeX)
	for i := 0; i < sizeY; i++ {
		for j := 0; j < sizeX; j++ {
			sum, counter := getNeighborsSum(i, j)
			tempPixels[i*sizeX+j] = uint8(sum / counter)
		}
	}
	pixels = tempPixels
}

func drawPolygon() {
	if stage == BUILDING_POLYGON || stage == POLYGON_BUILDED {
		if len(points) > 0 {
			gl.Begin(gl.LINE_LOOP)
			for _, p := range points {
				gl.Vertex2d(p.x, p.y)
			}
			if stage == BUILDING_POLYGON {
				gl.Vertex2d(mouse.x, mouse.y)
			}
			gl.End()
		}
	} else {
		gl.DrawPixels(int32(sizeX), int32(sizeY), gl.BLUE, gl.UNSIGNED_BYTE, unsafe.Pointer(&pixels[0]))

	}
}

func cycleInit(w *glfw.Window) {
	mouse.x, mouse.y = w.GetCursorPos()

}
func closeWindowCallback(w *glfw.Window) {
	log.Println("ESC")
	w.SetShouldClose(true)
}
func sizeCallback(w *glfw.Window, width int, height int) {
	sizeX, sizeY = width, height
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(width), float64(height), 0, 0, 1)

	gl.Viewport(0, 0, int32(width), int32(height))

	points = []point{}
	pixels = []uint8{}
	stage = BUILDING_POLYGON
}

func changeStateCallback() {
	if len(points) > 2 {
		stage = stage + 1
		if stage > FILTRATION {
			stage = BUILDING_POLYGON
		}
		if stage == RASTERISATION {
			rasterisation()
		}
		if stage == FILTRATION {
			filtrate()
		}
		log.Println(stage)
	}
}
func clear() {
	points = []point{}
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		if key == glfw.KeyEscape {
			closeWindowCallback(w)
		}
		if key == glfw.KeySpace {
			changeStateCallback()
		}
		if key == glfw.KeyDelete {
			clear()
		}
	}
}
func makePoint(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if stage == BUILDING_POLYGON {
		x, y := w.GetCursorPos()
		log.Println(x, y, " :mouse")
		points = append(points, point{x, y})
	}
}

func deletePoint(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if stage == BUILDING_POLYGON && len(points) > 0 {
		points = points[:len(points)-1]
	}
}

func mouseCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft && action == glfw.Press {
		makePoint(w, button, action, mod)
	}
	if button == glfw.MouseButtonRight && action == glfw.Press {
		deletePoint(w, button, action, mod)
	}

}

func initWindow() *glfw.Window {
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	window, err := glfw.CreateWindow(SIZE, SIZE, "LAB_4", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	return window
}

func main() {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	window := initWindow()

	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize gl:", err)
	}

	window.SetFramebufferSizeCallback(glfw.FramebufferSizeCallback(sizeCallback))
	window.SetKeyCallback(glfw.KeyCallback(keyCallback))
	window.SetMouseButtonCallback(glfw.MouseButtonCallback(mouseCallback))

	w, h := window.GetFramebufferSize()
	sizeCallback(window, w, h)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		cycleInit(window)

		drawPolygon()

		glfw.WaitEvents()
		window.SwapBuffers()
	}
}
