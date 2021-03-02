package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type point struct {
	x float32
	y float32
}

var points []point

const SIZE = 600

func closeWindowCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		log.Print("ESC")
		w.SetShouldClose(true)
	}
}

func drowPoint(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft && action == glfw.Press {
		x, y := w.GetCursorPos()
		log.Println((x/SIZE-0.5)*2, -2*(y/SIZE-0.5), " :mouse")
		points = append(points, point{float32(x/SIZE-0.5) * 2, -2 * float32(y/SIZE-0.5)})

	}
}

func clearScene(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if button == glfw.MouseButtonRight && action == glfw.Press {
		points = []point{}
		log.Println("CLEAR!!")
	}
}

func mouseCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	clearScene(w, button, action, mod)
	drowPoint(w, button, action, mod)
}

func makePointPool() {

	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.PointSize(10)
	gl.Begin(gl.POINTS)
	gl.Color3d(116/125, 185/125, 255/125)
	for _, p := range points {
		gl.Vertex2f(p.x, p.y)
	}
	gl.End()

	gl.Begin(gl.TRIANGLES)
	gl.Color3f(71/125, 38/125, 134/125)
	gl.Vertex2f(-0.7, -0.7)
	gl.Color3d(1, 0, 0.75)
	gl.Vertex2f(0, 0)
	gl.Color3f(0, 0.5, 0.5)
	gl.Vertex2f(0.7, -0.7)
	gl.End()

}

func initWindow() *glfw.Window {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	window, err := glfw.CreateWindow(SIZE, SIZE, "LAB_1", nil, nil)
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
	window.SetMouseButtonCallback(glfw.MouseButtonCallback(mouseCallback))

	for !window.ShouldClose() {
		makePointPool()

		glfw.WaitEvents()
		window.SwapBuffers()
	}

}
