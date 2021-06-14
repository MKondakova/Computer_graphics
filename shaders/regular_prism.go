package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"io/ioutil"
	"log"
	"math"
	"os"
	"runtime"
	"time"

	"strings"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const SIZE = 600
const HEIGHT = 0.5

type SaveStruct struct {
	Alpha                   float32
	LastXpos                float64
	LastYpos                float64
	Yaw                     float64
	Pitch                   float64
	Scale                   float64
	SetPolygonMode          bool
	SetInfinityDistantLight bool
	AmbientMode             int
	DiffuseMode             int
	SpecularMode            int

	T          float64
	Phase      int
	TextureMod int
}

var (
	RADIUS  float64 = math.Sqrt(0.5) * 0.5
	CORNERS int     = 6

	alpha    float32 = 0
	lastXpos float64 = SIZE / 2
	lastYpos float64 = SIZE / 2
	yaw      float64 = -90
	pitch    float64 = 0
	scale    float64 = 1

	setPolygonMode          bool = false
	setInfinityDistantLight bool = false

	ambientMode  int         = 0
	ambient      [][]float32 = [][]float32{{0, 0, 0, 1}, {1, 1, 1, 0.5}, {1, 1, 1, 1}, {0.5, 0.5, 0.5, 1}, {0, 1, 0, 1}}
	diffuseMode  int         = 0
	diffuse      [][]float32 = [][]float32{{1, 1, 1, 1}, {0, 0, 0, 1}, {1, 1, 1, 0.5}, {0.5, 0.5, 0.5, 1}, {0, 1, 0, 1}}
	specularMode int         = 0
	specular     [][]float32 = [][]float32{{1, 1, 1, 1}, {0, 0, 0, 1}, {1, 1, 1, 0.5}, {0.5, 0.5, 0.5, 1}, {0, 1, 0, 1}}

	lightPosition []float32 = []float32{0, 0, 1, 1}

	t              float64      = 0.0
	curveTicker    *time.Ticker = time.NewTicker(50 * time.Millisecond)
	animationSpeed float64      = 0.01
	phase          int          = 0

	POINT1 []float64 = []float64{0, 0.6, 0}
	POINT2 []float64 = []float64{0.6, 0.6, 0}
	POINT3 []float64 = []float64{0.6, 0, 0}

	generatedTexture uint32 = 0
	loadedTexture    uint32 = 0
	textureMod       int    = 0

	program uint32 = 0
)

func drawBase(vertexes [][2]float64, normals [][3]float64, z float64) {
	offset := 0
	if z < 0 {
		offset = len(vertexes)
	}
	gl.Color3d(z, 1, 1)
	vertexArray := make([][3]float64, len(vertexes))
	normalArray := make([][3]float64, len(vertexes))
	for i, vertex := range vertexes {
		normalArray[i] = [3]float64{normals[offset+i][0], normals[offset+i][1], normals[offset+i][2]}
		vertexArray[i] = [3]float64{vertex[0], vertex[1], z}
	}

	name := "isTexture\000"
	isTexture := gl.GetUniformLocation(program, gl.Str(name))
	gl.Uniform1f(isTexture, 0)

	gl.EnableClientState(gl.NORMAL_ARRAY)
	gl.EnableClientState(gl.VERTEX_ARRAY)

	gl.VertexPointer(3, gl.DOUBLE, 0, gl.Ptr(&vertexArray[0][0]))
	gl.NormalPointer(gl.DOUBLE, 0, gl.Ptr(&normalArray[0][0]))

	gl.DrawArrays(gl.POLYGON, 0, int32(len(vertexArray)))

	gl.DisableClientState(gl.NORMAL_ARRAY)
	gl.DisableClientState(gl.VERTEX_ARRAY)
}

func drawSideFaces(vertexes [][2]float64, normals [][3]float64, height float64) {
	name := "isTexture\000"
	isTexture := gl.GetUniformLocation(program, gl.Str(name))
	gl.Uniform1f(isTexture, 0)

	if textureMod == 1 {
		gl.BindTexture(gl.TEXTURE_2D, generatedTexture)
		gl.TexEnvi(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.MODULATE)
	}
	if textureMod == 2 {
		gl.BindTexture(gl.TEXTURE_2D, loadedTexture)
		gl.TexEnvi(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.MODULATE)
	}
	vertexArray := [][3]float64{}
	normalArray := [][3]float64{}
	textureArray := [][2]float64{}
	gl.Color4d(1, 1, 1, 1)
	gl.EnableClientState(gl.VERTEX_ARRAY)
	gl.EnableClientState(gl.NORMAL_ARRAY)
	if textureMod > 0 {
		gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)
		gl.Uniform1f(isTexture, 1)
	}

	for i := 0; i < len(vertexes); i++ {

		normalArray = append(normalArray,
			[3]float64{normals[i+len(vertexes)][0], normals[i+len(vertexes)][1], normals[i+len(vertexes)][2]})
		vertexArray = append(vertexArray, [3]float64{vertexes[i][0], vertexes[i][1], height / -2})

		normalArray = append(normalArray, [3]float64{normals[i][0], normals[i][1], normals[i][2]})
		vertexArray = append(vertexArray, [3]float64{vertexes[i][0], vertexes[i][1], height / 2})

		normalArray = append(normalArray, [3]float64{normals[(i+1)%len(vertexes)][0], normals[(i+1)%len(vertexes)][1], normals[(i+1)%len(vertexes)][2]})
		vertexArray = append(vertexArray, [3]float64{vertexes[(i+1)%len(vertexes)][0], vertexes[(i+1)%len(vertexes)][1], height / 2})

		normalArray = append(normalArray, [3]float64{normals[(i+1)%len(vertexes)+len(vertexes)][0],
			normals[(i+1)%len(vertexes)+len(vertexes)][1],
			normals[(i+1)%len(vertexes)+len(vertexes)][2]})
		vertexArray = append(vertexArray, [3]float64{vertexes[(i+1)%len(vertexes)][0], vertexes[(i+1)%len(vertexes)][1], height / -2})

		textureArray = append(textureArray, [2]float64{0, 0}, [2]float64{1, 0}, [2]float64{1, 1}, [2]float64{0, 1})
	}

	texture := gl.GetUniformLocation(program, gl.Str("texture\000"))
	gl.Uniform1i(texture, 0)

	gl.TexCoordPointer(2, gl.DOUBLE, 0, gl.Ptr(&textureArray[0][0]))
	gl.VertexPointer(3, gl.DOUBLE, 0, gl.Ptr(&vertexArray[0][0]))
	gl.NormalPointer(gl.DOUBLE, 0, gl.Ptr(&normalArray[0][0]))

	gl.DrawArrays(gl.QUADS, 0, int32(len(vertexArray)))

	gl.DisableClientState(gl.NORMAL_ARRAY)
	gl.DisableClientState(gl.VERTEX_ARRAY)
	gl.DisableClientState(gl.TEXTURE_COORD_ARRAY)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.Uniform1f(isTexture, 0)

}

func drawPrism(n int) {
	vertexes := [][2]float64{}
	normals := make([][3]float64, n*2)
	for i := 0; i < n; i++ {
		vertexes = append(vertexes, [2]float64{
			RADIUS * math.Cos(2*math.Pi*float64(i)/float64(n)+math.Pi*45/180),
			RADIUS * math.Sin(2*math.Pi*float64(i)/float64(n)+math.Pi*45/180),
		})
	}
	for i := 0; i < n; i++ {
		normals[i][0] = -1 * (vertexes[(n+i-1)%n][0] + vertexes[(n+i+1)%n][0] - 2*vertexes[i][0])
		normals[i][1] = -1 * (vertexes[(n+i-1)%n][1] + vertexes[(n+i+1)%n][1] - 2*vertexes[i][1])
		normals[i][2] = HEIGHT / 2
		normals[n+i][0] = normals[i][0]
		normals[n+i][1] = normals[i][1]
		normals[n+i][2] = HEIGHT / -2
	}
	drawBase(vertexes, normals, HEIGHT/-2)
	drawBase(vertexes, normals, HEIGHT/2)
	drawSideFaces(vertexes, normals, HEIGHT)
}

func getBezierPosition(t float64, p1, p2, p3 float64) float64 {
	t -= math.Floor(t)
	if t > 0 {
		return (1-t)*(1-t)*p1 + 2*t*(1-t)*p2 + t*t*p3
	}
	return 0
}

func drawMovingPrism() {
	gl.PushMatrix()

	gl.Rotated(yaw, 0, 1, 0)
	gl.Rotated(pitch, 1, 0, 0)
	gl.Scaled(scale, scale, scale)
	drawPrism(CORNERS)

	gl.PopMatrix()

	gl.Begin(gl.POINTS)

	gl.PointSize(5)
	gl.Vertex3dv(&POINT1[0])
	gl.Vertex3dv(&POINT2[0])
	gl.Vertex3dv(&POINT3[0])

	gl.PointSize(10)
	gl.Vertex3d(
		getBezierPosition(t, POINT1[0], POINT2[0], POINT3[0]),
		getBezierPosition(t, POINT1[1], POINT2[1], POINT3[1]),
		getBezierPosition(t, POINT1[2], POINT2[2], POINT3[2]))

	gl.End()
}

func setLight() {

	if setInfinityDistantLight {
		lightPosition[3] = 0
	} else {
		lightPosition[3] = 1
	}
	gl.Lightfv(gl.LIGHT0, gl.POSITION, &lightPosition[0])
	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, &ambient[ambientMode][0])
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, &diffuse[diffuseMode][0])
	gl.Lightfv(gl.LIGHT0, gl.SPECULAR, &specular[specularMode][0])

	gl.Color3d(1, 1, 1)
	gl.PointSize(10)
	gl.Normal3b(0, 0, -1)

	gl.Begin(gl.POINTS)
	gl.Vertex3fv(&lightPosition[0])
	gl.End()
}
func setUniformVariables() {
	uniformLightPos := gl.GetUniformLocation(program, gl.Str("lightPos\000"))
	gl.Uniform3fv(uniformLightPos, 1, &lightPosition[0])

	uniformAmbient := gl.GetUniformLocation(program, gl.Str("ambient\000"))
	gl.Uniform4fv(uniformAmbient, 1, &ambient[ambientMode][0])
	uniformDiffuse := gl.GetUniformLocation(program, gl.Str("diffuse\000"))
	gl.Uniform4fv(uniformDiffuse, 1, &diffuse[diffuseMode][0])

}

func loadTexture() {
	imgFile, err := os.Open("../textures/square.png")
	if err != nil {
		log.Panicln("texture not found on disk: ", err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Panicln("unsupported stride", err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		log.Panicln("unsupported stride", err)

	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	gl.GenTextures(2, &loadedTexture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, loadedTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
}

func generateTexture() {
	data := [2][2][4]uint8{{{255, 0, 0, 0}, {255, 255, 0, 0}}, {{0, 255, 0, 0}, {0, 0, 255, 0}}}
	gl.GenTextures(1, &generatedTexture)
	gl.BindTexture(gl.TEXTURE_2D, generatedTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 2, 2, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&data[0][0][0]))
	log.Println(generatedTexture)
	gl.BindTexture(gl.TEXTURE_2D, 0)

}
func saveState() {

	file, _ := json.MarshalIndent(SaveStruct{alpha, lastXpos, lastYpos, yaw, pitch, scale, setPolygonMode,
		setInfinityDistantLight, ambientMode, diffuseMode,
		specularMode, t, phase, textureMod}, "", " ")

	_ = ioutil.WriteFile("test.json", file, 0644)
}
func loadState() {
	jsonFile, err := os.Open("test.json")
	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var state SaveStruct
	json.Unmarshal(byteValue, &state)
	alpha = state.Alpha
	lastXpos = state.LastXpos
	lastYpos = state.LastYpos
	yaw = state.Yaw
	pitch = state.Pitch
	scale = state.Scale
	setPolygonMode = state.SetPolygonMode
	setInfinityDistantLight = state.SetInfinityDistantLight
	ambientMode = state.AmbientMode
	diffuseMode = state.DiffuseMode
	specularMode = state.SpecularMode
	t = state.T
	phase = state.Phase
	textureMod = state.TextureMod
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		if key == glfw.KeyI {
			setInfinityDistantLight = !setInfinityDistantLight
		}
		if key == glfw.KeyA {
			ambientMode = (ambientMode + 1) % len(ambient)
			log.Println("ambient: ", ambient[ambientMode])
		}
		if key == glfw.KeyD {
			diffuseMode = (diffuseMode + 1) % len(diffuse)
			log.Println("diffuse: ", diffuse[diffuseMode])
		}
		if key == glfw.KeyS {
			specularMode = (specularMode + 1) % len(specular)
			log.Println("specular: ", specular[specularMode])
		}
		if key == glfw.KeyT {
			textureMod = (textureMod + 1) % 3
		}
		if key == glfw.KeyP {
			saveState()
		}
		if key == glfw.KeyL {
			loadState()
		}
		if key == glfw.KeyMinus {
			if CORNERS != 3 {
				CORNERS -= 1
			}
			log.Println(CORNERS)
		}
		if key == glfw.KeyEqual {
			if CORNERS < 100 {
				CORNERS += 1
			}
			log.Println(CORNERS)
		}
		if key == glfw.KeyEscape {
			log.Println("ESC")
			w.SetShouldClose(true)
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

func tick(ticker *time.Ticker, f func(), isEnd func() bool, stop chan bool) {
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if isEnd() {
				stop <- true
			}
			f()
		case <-stop:
			return
		}
	}
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program = gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func initWindow() *glfw.Window {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	window, err := glfw.CreateWindow(SIZE, SIZE, "LAB_6", nil, nil)
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

	fragmentShader, err := ioutil.ReadFile("frag.glsl")
	if err != nil {
		log.Fatal(err)
	}
	vertexShader, err := ioutil.ReadFile("vert.glsl")
	if err != nil {
		log.Fatal(err)
	}

	program, err := newProgram(string(vertexShader), string(fragmentShader))
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.NORMALIZE)
	gl.Enable(gl.COLOR_MATERIAL)
	gl.Enable(gl.TEXTURE_2D)
	generateTexture()
	defer gl.DeleteTextures(1, &generatedTexture)
	loadTexture()
	defer gl.DeleteTextures(2, &loadedTexture)

	gl.Enable(gl.LIGHTING)
	gl.Enable(gl.LIGHT0)

	backLight := []float32{0.3, 0.3, 0.3, 1}
	gl.LightModelfv(gl.LIGHT_MODEL_AMBIENT, &backLight[0])

	go tick(curveTicker, func() {
		t += float64(phase*-1*2+1) * animationSpeed
		if t < 0 || t > 1 {
			phase = (1 + phase) % 2
			t += float64(phase*-1*2+1) * animationSpeed
		}
	},
		func() bool { return false }, make(chan bool, 1))
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.LoadIdentity()

		if setPolygonMode {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		} else {
			gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		}

		width, height := window.GetSize()
		gl.Viewport(0, 0, int32(width), int32(height))

		setUniformVariables()
		setLight()
		drawMovingPrism()

		glfw.PollEvents()
		window.SwapBuffers()
	}

}
