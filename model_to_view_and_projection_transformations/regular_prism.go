package main

import (
	"log"
	"math"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const SIZE = 1000
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

	lastXpos       float64 = SIZE / 2
	lastYpos       float64 = SIZE / 2
	yaw            float64 = -90
	pitch          float64 = 0
	scale          float64 = 1
	setProjection  bool    = false
	setPolygonMode bool    = false
	CORNERS        int     = 6
)

func drawBase(vertexes [][2]float64, z float64) {

	gl.Begin(gl.POLYGON)
	gl.Color3d(z, 1, 1)
	for _, vertex := range vertexes {
		gl.Vertex3d(vertex[0], vertex[1], z)
	}
	gl.End()
}

func drawSideFaces(vertexes [][2]float64, height float64) {
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

func drawPrism(n int) {
	vertexes := [][2]float64{}
	for i := 0; i < n; i++ {
		vertexes = append(vertexes, [2]float64{
			RADIUS * math.Cos(2*math.Pi*float64(i)/float64(n)+math.Pi*45/180),
			RADIUS * math.Sin(2*math.Pi*float64(i)/float64(n)+math.Pi*45/180),
		})
	}
	drawBase(vertexes, HEIGHT/-2)
	drawBase(vertexes, HEIGHT/2)
	drawSideFaces(vertexes, HEIGHT)
}

func closeWindowCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		log.Println("ESC")
		w.SetShouldClose(true)
	}
}

func changeNumberOfCorners(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		if key == glfw.KeyMinus {
			if CORNERS != 3 {
				CORNERS -= 1
			}
		}
		if key == glfw.KeyEqual {
			if CORNERS < 100 {
				CORNERS += 1
			}
		}
		log.Println(CORNERS)
	}
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	closeWindowCallback(w, key, scancode, action, mods)
	changeNumberOfCorners(w, key, scancode, action, mods)
}

func makeProjection(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft && action == glfw.Press {
		if setProjection == true {
			setProjection = false
		} else {
			setProjection = true
		}

	}
}

func makeModePolygon(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if button == glfw.MouseButtonRight && action == glfw.Press {
		if setPolygonMode == true {
			setPolygonMode = false
		} else {
			setPolygonMode = true
		}
	}
}

func mouseCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	makeProjection(w, button, action, mod)
	makeModePolygon(w, button, action, mod)
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

	window.SetKeyCallback(glfw.KeyCallback(keyCallback))
	window.SetCursorPosCallback(glfw.CursorPosCallback(mouseCursorCallback))
	window.SetScrollCallback(glfw.ScrollCallback(mouseScrollCallback))
	window.SetMouseButtonCallback(glfw.MouseButtonCallback(mouseCallback))

	gl.Enable(gl.DEPTH_TEST)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.MatrixMode(gl.PROJECTION)
		if setProjection {
			gl.LoadMatrixd(&isometricProjectionMatrix[0])
		} else {
			gl.LoadIdentity()
		}
		gl.MatrixMode(gl.MODELVIEW)

		if setPolygonMode {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		} else {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		}

		gl.LoadIdentity()
		width, height := window.GetSize()
		aspect := float64(width) / float64(height)
		gl.Viewport(0, 0, int32(width), int32(height))
		gl.Ortho(-1*aspect, aspect, -1, 1, 1.0, -1.0)

		gl.Translated(-0.4, -0.4, 0)
		drawPrism(4)

		gl.PushMatrix()

		gl.Translated(0.8, 0.8, 0)
		gl.Rotated(yaw, 0, 1, 0)
		gl.Rotated(pitch, 1, 0, 0)
		gl.Scaled(scale, scale, scale)
		drawPrism(CORNERS)

		gl.PopMatrix()

		glfw.WaitEvents()
		window.SwapBuffers()
	}

}
