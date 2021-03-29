package main

import (
	"log"
	"math"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const SIZE = 600
const HEIGHT = 0.5

var isometricProjectionMatrix []float64 = []float64{

	math.Cos(phi), 0.0, math.Sin(phi), 0.0,
	math.Sin(phi) * math.Sin(theta), math.Cos(theta), -1 * math.Cos(phi) * math.Sin(theta), 0.0,
	math.Sin(phi) * math.Cos(theta), -1 * math.Sin(theta), -1 * math.Cos(phi) * math.Cos(theta), 0.0,
	0.0, 0.0, 0.0, 1.0,
}

var (
	theta  float64 = 45 * math.Pi / 180.0
	phi    float64 = 35.26 * math.Pi / 180.0
	RADIUS float64 = math.Sqrt(0.5) * 0.5

	lastXpos float64 = SIZE / 2
	lastYpos float64 = SIZE / 2
	yaw      float64 = -90
	pitch    float64 = 0
	scale    float64 = 1
)

func drowBase(vertexes [][2]float64, z float64) {

	gl.Begin(gl.POLYGON)
	gl.Color3d(z, 1, 1)
	for _, vertex := range vertexes {
		gl.Vertex3d(vertex[0], vertex[1], z)
	}
	gl.End()
}

func drowSideFaces(vertexes [][2]float64, height float64) {
	for i := 0; i < len(vertexes); i++ {
		gl.Begin(gl.QUADS)
		gl.Color3d(0.2, 0.2, float64(i)/3)
		gl.Vertex3d(vertexes[i][0], vertexes[i][1], height/-2)
		gl.Vertex3d(vertexes[i][0], vertexes[i][1], height/2)
		gl.Vertex3d(vertexes[(i+1)%len(vertexes)][0], vertexes[(i+1)%len(vertexes)][1], height/2)
		gl.Vertex3d(vertexes[(i+1)%len(vertexes)][0], vertexes[(i+1)%len(vertexes)][1], height/-2)
		gl.End()
	}

}

func drowPrism(n int) {
	vertexes := [][2]float64{}
	for i := 0; i < n; i++ {
		vertexes = append(vertexes, [2]float64{
			RADIUS * math.Cos(2*math.Pi*float64(i)/float64(n)+math.Pi*45/180),
			RADIUS * math.Sin(2*math.Pi*float64(i)/float64(n)+math.Pi*45/180),
		})
	}
	drowBase(vertexes, HEIGHT/-2)
	drowBase(vertexes, HEIGHT/2)
	drowSideFaces(vertexes, HEIGHT)
}

func closeWindowCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		log.Print("ESC")
		w.SetShouldClose(true)
	}
}

func mouseCursorCallback(w *glfw.Window, xpos float64, ypos float64) {
	xOffset := xpos - lastXpos
	yOffset := lastYpos - ypos

	lastXpos = xpos
	lastYpos = ypos

	sensetivity := 0.2
	xOffset *= sensetivity
	yOffset *= sensetivity

	yaw += xOffset
	pitch += yOffset

}

func mouseScrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	sensetivity := 0.05
	scale -= yoff * sensetivity
	if scale < 0.05 {
		scale = 0.05
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
	window.SetCursorPosCallback(glfw.CursorPosCallback(mouseCursorCallback))
	window.SetScrollCallback(glfw.ScrollCallback(mouseScrollCallback))

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	gl.Enable(gl.DEPTH_TEST)
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	gl.MatrixMode(gl.PROJECTION)
	//gl.LoadMatrixd(&isometricProjectionMatrix[0])
	gl.MatrixMode(gl.MODELVIEW)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.LoadIdentity()
		gl.Translated(-0.6, -0.6, 0)
		drowPrism(4)
		gl.PushMatrix()
		gl.Translated(0.6, 0.6, 0)
		gl.Rotated(yaw, 0, 1, 0)
		gl.Rotated(pitch, 1, 0, 0)
		gl.Scaled(scale, scale, scale)
		drowPrism(10)
		gl.PopMatrix()
		glfw.WaitEvents()
		window.SwapBuffers()
	}

}
