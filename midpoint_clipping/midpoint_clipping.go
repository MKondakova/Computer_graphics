package main

import (
	"log"
	"math"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const (
	PLOTTING_SEGMENTS        = 1
	SEGMENTS_PLOTTED         = 2
	CLIPPING                 = 3
	SIZE                     = 1000
	ZONE_PADDING_COEFFICIENT = 0.2
	ACCURACY                 = math.Sqrt2
)

type point struct {
	x float64
	y float64
}

var (
	mouse     point
	stage     int = PLOTTING_SEGMENTS
	sizeX     int
	sizeY     int
	segments  [][2]point
	points    []point
	zoneLeft  float64
	zoneRight float64
	zoneCeil  float64
	zoneFloor float64
)

func clipping() {
	tempSegments := makeSegments()
	for _, segment := range tempSegments {
		midpointClipping(segment, 1)
	}
}

func midpointClipping(segment [2]point, count int) {
	firstCode := getCode(segment[0])
	secondCode := getCode(segment[1])
	if firstCode+secondCode == 0 {
		segments = append(segments, segment)
		return
	}
	if firstCode&secondCode != 0 {
		return
	}
	if count > 2 {
		segments = append(segments, segment)
		return
	}
	firstPoint := segment[0]
	if secondCode == 0 {
		segment[1], segment[0] = firstPoint, segment[1]
		count++
		midpointClipping(segment, count)
		return
	}
	for {
		if math.Hypot(segment[0].x-segment[1].x, segment[0].y-segment[1].y) <= ACCURACY {
			segment[1], segment[0] = firstPoint, segment[1]
			count++
			midpointClipping(segment, count)
			return
		}
		midpoint := point{(segment[0].x + segment[1].x) / 2, (segment[0].y + segment[1].y) / 2}
		memoizedPoint := segment[0]
		segment[0] = midpoint
		firstCode = getCode(segment[0])
		if firstCode&secondCode != 0 {
			segment[0], segment[1] = memoizedPoint, midpoint
		}
	}
}

func getCode(p point) int {
	code := 0
	if p.y > zoneFloor {
		code++
	}
	code *= 2
	if p.x > zoneRight {
		code++
	}
	code *= 2
	if p.y < zoneCeil {
		code++
	}
	code *= 2
	if p.x < zoneLeft {
		code++
	}
	return code
}

func makeSegments() [][2]point {
	temp := [][2]point{}
	for i := 0; i < len(points)/2; i++ {
		temp = append(temp, [2]point{points[2*i], points[2*i+1]})
	}
	return temp
}

func drawZone() {
	zoneLeft = float64(sizeX) * ZONE_PADDING_COEFFICIENT
	zoneRight = float64(sizeX) * (1 - ZONE_PADDING_COEFFICIENT)
	zoneFloor = float64(sizeY) * (1 - ZONE_PADDING_COEFFICIENT)
	zoneCeil = float64(sizeY) * ZONE_PADDING_COEFFICIENT

	gl.Begin(gl.LINE_LOOP)
	gl.Vertex2d(zoneLeft, zoneFloor)
	gl.Vertex2d(zoneLeft, zoneCeil)
	gl.Vertex2d(zoneRight, zoneCeil)
	gl.Vertex2d(zoneRight, zoneFloor)
	gl.End()
}

func drawSegments() {
	if stage == PLOTTING_SEGMENTS || stage == SEGMENTS_PLOTTED {
		gl.Begin(gl.LINES)
		for _, p := range points {
			gl.Vertex2d(p.x, p.y)
		}
		if stage == PLOTTING_SEGMENTS {
			gl.Vertex2d(mouse.x, mouse.y)
		}
		gl.End()
	} else {
		gl.Begin(gl.LINES)
		for _, segment := range segments {
			gl.Vertex2d(segment[0].x, segment[0].y)
			gl.Vertex2d(segment[1].x, segment[1].y)
		}
		gl.End()
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

	clear()
}

func changeStateCallback() {
	if len(points) >= 2 {
		stage = stage + 1
		if stage > CLIPPING {
			stage = PLOTTING_SEGMENTS
		}
		if stage == CLIPPING {
			clipping()
		}
		log.Println(stage)
	}
}

func clear() {
	segments = [][2]point{}
	points = []point{}
	stage = PLOTTING_SEGMENTS
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
	if stage == PLOTTING_SEGMENTS {
		x, y := w.GetCursorPos()
		log.Println(x, y, " :mouse")
		points = append(points, point{x, y})
	}
}

func deletePoint(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if stage == PLOTTING_SEGMENTS && len(points) > 0 {
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
	window, err := glfw.CreateWindow(SIZE, SIZE, "LAB_5", nil, nil)
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

		drawZone()
		drawSegments()

		glfw.WaitEvents()
		window.SwapBuffers()
	}
}
