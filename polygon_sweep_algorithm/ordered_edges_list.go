package main

import (
	"log"
	"math"
	"runtime"
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
type line struct {
	y float64
	x []float64
}

var (
	mouse      point
	stage      int = BUILDING_POLYGON
	points     []point
	sizeX      int
	sizeY      int
	pixels     []uint8
	tempPixels []uint8
	list       map[float64]line
	edges      [][2]point
)

func makeEdges() {
	edges = make([][2]point, len(points))
	for i, p := range points {
		nextP := points[(i+1)%len(points)]
		if p.y > nextP.y {
			p, nextP = nextP, p
		}
		edges[i] = [2]point{p, nextP}
	}
}
func eqVertex(v1, v2 point) bool {
	return v1.x == v2.x && v1.y == v2.y
}
func vertexCountTwice(i, j int) bool {
	vertex := edges[i][j]
	fDiffVertex := edges[i][(j+1)%2]
	var sDiffVertex point
	if eqVertex(edges[(i+1)%len(edges)][0], vertex) {
		sDiffVertex = edges[(i+1)%len(edges)][1]
	}
	if eqVertex(edges[(i+1)%len(edges)][1], vertex) {
		sDiffVertex = edges[(i+1)%len(edges)][1]

	}
	if eqVertex(edges[(len(edges)-(i-1))%len(edges)][0], vertex) {
		sDiffVertex = edges[(i+1)%len(edges)][1]

	}
	if eqVertex(edges[(len(edges)-i+1)%len(edges)][1], vertex) {
		sDiffVertex = edges[(i+1)%len(edges)][1]
	}
	if (vertex.y > fDiffVertex.y && vertex.y > sDiffVertex.y) ||
		(vertex.y < fDiffVertex.y && vertex.y < sDiffVertex.y) {
		return true
	}
	return false
}

func DDA() {
	for i, edge := range edges {
		dx := edge[1].x - edge[0].x
		x := edge[0].x
		dy := edge[1].y - edge[0].y
		y := edge[0].y
		counter := math.Max(dx, dy)
		dy = dy / counter
		dx = dx / counter
		for i := float64(0); i < counter; i++ {
			list[int(math.Round(float64(sizeY)-(y+i*dy)))] = int(math.Round(x + i*dx))
		}
		if !vertexCountTwice(i, 1) {
			list[int(math.Round(edge[1].y))] = int(math.Round(edge[1].x))
		}
		if !vertexCountTwice(i, 0) {
			list[int(math.Round(edge[0].y))] = int(math.Round(edge[0].x))
		}
	}
}

func fill() {
	previousPixel := false
	previousZone := false
	for i := 0; i < sizeY; i++ {
		for j := 0; j < sizeX; j++ {
			if pixels[i*sizeX+j] == 0 && previousPixel {
				if previousZone {
					previousZone = false
				} else {
					for ; pixels[i*sizeX+j] == 0 && j < sizeX; j++ {
						pixels[i*sizeX+j] = 255
					}
					previousPixel = false
					previousZone = true
				}
			} else if pixels[i*sizeX+j] > 0 {
				previousPixel = true
			}
		}
	}
}

func rasterisation() {
	pixels = make([]uint8, sizeY*sizeX)
	list = make(map[int]int)
	makeEdges()
	DDA()
	fill()
}

func getNeighborsSum(i, j int) (uint8, uint8) {
	result := uint8(0)
	counter := uint8(0)
	for k := i - 1; k+1-i < 3; k++ {
		for l := j - 1; l+1-j < 3; l++ {
			if k >= 0 && k < sizeY && l >= 0 && l < sizeX {
				result += pixels[int(sizeX*k+l)]
				counter++
			}
		}
	}
	return result, counter
}

func filtrate() {
	tempPixels = make([]uint8, sizeY*sizeX)
	for i := 0; i < sizeY; i++ {
		for j := 0; j < sizeX; j++ {
			newColour, counter := getNeighborsSum(i, j)
			tempPixels[i*sizeX+j] = newColour/counter/2 + pixels[i*sizeX+j]/2
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
		gl.DrawPixels(int32(sizeX), int32(sizeY), gl.GREEN, gl.UNSIGNED_BYTE, unsafe.Pointer(&pixels[0]))
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
