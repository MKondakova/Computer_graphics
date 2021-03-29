package main

import (
	"log"
	"math"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const SIZE = 600

var RADIUS float64 = math.Sqrt(0.5)
var HEIGHT float64 = 0.5

func drowBase(vertexes [][2]float64, z float64) {

	gl.Begin(gl.POLYGON)
	gl.Color3d(z, 0.2, 0.2)
	for _, vertex := range vertexes {
		gl.Vertex3d(vertex[0], vertex[1], z)
		log.Println(vertex[0], vertex[1], z)
	}
	gl.End()
}

func drowSideFaces(vertexes [][2]float64, height float64) {
	for i := 0; i < len(vertexes); i++ {
		gl.Begin(gl.QUADS)
		gl.Color3d(0.5, 0.3, float64(i)/5)
		gl.Vertex3d(vertexes[i][0], vertexes[i][1], 0)
		gl.Vertex3d(vertexes[i][0], vertexes[i][1], height)
		gl.Vertex3d(vertexes[(i+1)%len(vertexes)][0], vertexes[(i+1)%len(vertexes)][1], height)
		gl.Vertex3d(vertexes[(i+1)%len(vertexes)][0], vertexes[(i+1)%len(vertexes)][1], 0)
		gl.End()
	}

}

func drowPrism(n int) {
	vertexes := [][2]float64{}
	for i := 0; i < n; i++ {
		vertexes = append(vertexes, [2]float64{
			RADIUS * math.Cos(2*math.Pi*float64(i)/float64(n)),
			RADIUS * math.Sin(2*math.Pi*float64(i)/float64(n)),
		})
	}
	drowBase(vertexes, 0)
	drowBase(vertexes, HEIGHT)
	drowSideFaces(vertexes, HEIGHT)
}

func closeWindowCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		log.Print("ESC")
		w.SetShouldClose(true)
	}
}

func initWindow() *glfw.Window {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	window, err := glfw.CreateWindow(SIZE, SIZE, "LAB_2/3", nil, nil)
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

	window.SetKeyCallback(glfw.KeyCallback(closeWindowCallback))
	gl.Enable(gl.DEPTH_TEST)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.LoadIdentity()
		gl.Rotatef(50, 1, 1, 0.5)
		drowPrism(7)
		glfw.WaitEvents()
		window.SwapBuffers()
	}

}
