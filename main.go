package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"helloOpenGLWindow/shader"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"runtime"
)

var vertices = []float32{
	// positions       // colors        // texture coords
	0.5,  0.5, 0.0,    1.0, 0.0, 0.0,   2, 2,   // top right
	0.5, -0.5, 0.0,    0.0, 1.0, 0.0,   2, 0,   // bottom right
	-0.5, -0.5, 0.0,   0.0, 0.0, 1.0,   0, 0,   // bottom let
	-0.5,  0.5, 0.0,   1.0, 1.0, 0.0,   0, 2,   // top let
}

var triangles = []uint32 {
	0, 1, 3, // first triangle
	1, 2, 3,  // second triangle
}

const texture1Path = "bricks.jpg"
//var texture2Path = "Badger.jpg"
const texture2Path = "vines.png"
//var texture1Path = "container.jpg"
//var texture2Path = "awesomeface.png"

var mixValue float32 = 0

func main() {

	runtime.LockOSThread()

	glfw.Init()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, 1)

	window, err := glfw.CreateWindow(1024, 768, "Scott Window", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	initOpenGL()

	log.Println("Creating and Compiling Shaders")
	vertexShaderPath := "./shader/shaders/transformShader.vs"
	fragmentShaderPath := "./shader/shaders/transformShader.fs"
	shader := shader.New(vertexShaderPath, fragmentShaderPath)

	vao2 := makeVao(vertices, triangles)

	rgba1 := loadImage(texture1Path)
	textureWrap1 := gl.CLAMP_TO_EDGE
	texture1:= initTexture(rgba1, textureWrap1)

	rgba2 := loadImage(texture2Path)
	textureWrap2 := gl.REPEAT
	texture2:= initTexture(rgba2, textureWrap2)

	prog := gl.CreateProgram()
	gl.AttachShader(prog, shader.VertexShaderCompiled)
	gl.AttachShader(prog, shader.FragmentShaderCompiled)
	gl.LinkProgram(prog)
	gl.UseProgram(prog)
	gl.Uniform1i(gl.GetUniformLocation(prog, gl.Str("texture1\x00")), int32(0))
	if e := gl.GetError(); e != 0 {
		log.Fatalf("ERROR: %X", e)
	}
	gl.Uniform1i(gl.GetUniformLocation(prog, gl.Str("texture2\x00")), int32(1))
	if e := gl.GetError(); e != 0 {
		log.Fatalf("ERROR: %X", e)
	}
	//render loop
	for !window.ShouldClose() {
		//process commands
		processInput(window)

		//render commands
		gl.ClearColor(.2, .3, .3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.ActiveTexture(gl.TEXTURE0)
		if e := gl.GetError(); e != 0 {
			log.Fatalf("ERROR: %X", e)
		}
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		if e := gl.GetError(); e != 0 {
			log.Fatalf("ERROR: %X", e)
		}

		gl.ActiveTexture(gl.TEXTURE1)
		if e := gl.GetError(); e != 0 {
			log.Fatalf("ERROR: %X", e)
		}
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		if e := gl.GetError(); e != 0 {
			log.Fatalf("ERROR: %X", e)
		}

		// render container
		gl.UseProgram(prog)

		mixLocation := gl.GetUniformLocation(prog, gl.Str("mixValue\x00"))
		gl.UseProgram(prog)
		gl.Uniform1f(mixLocation, mixValue)

		gl.BindVertexArray(vao2)
		//gl.DrawArrays(gl.TRIANGLES, 0, 3) //draw 1 triangle
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

		glfw.PollEvents()
		window.SwapBuffers()
	}
	glfw.Terminate()
}

func initTexture(rgba *image.RGBA, wrap int) uint32 {
	var texture uint32
	// load and create a texture
	// -------------------------
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture) // all upcoming GL_TEXTURE_2D operations now have effect on this texture object
	// set the texture wrapping parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, int32(wrap)) // set texture wrapping to GL_REPEAT (default wrapping method)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, int32(wrap))
	// set texture filtering parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// flip image: first pixel is lower left corner
	imgWidth, imgHeight := rgba.Rect.Bounds().Dx(), rgba.Rect.Bounds().Dy()
	data := make([]byte, imgWidth*imgHeight*4)
	lineLen := imgWidth * 4
	dest := len(data) - lineLen
	for src := 0; src < len(rgba.Pix); src += rgba.Stride {
		copy(data[dest:dest+lineLen], rgba.Pix[src:src+rgba.Stride])
		dest -= lineLen
	}

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Bounds().Dx()),
		int32(rgba.Rect.Bounds().Dy()),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return texture
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
	if window.GetKey(glfw.KeyUp) == glfw.Press{
		mixValue = mixValue + 0.1
	}
	if window.GetKey(glfw.KeyDown) == glfw.Press{
		mixValue = mixValue - 0.1
	}
}

func initOpenGL() {
	var nrAttributes int32
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))

	gl.GetIntegerv(gl.MAX_VERTEX_ATTRIBS, &nrAttributes)
	log.Println("nrAttributes", nrAttributes)
	log.Println("OpenGL version", version)
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32, triangles []uint32) uint32 {
	var vbo uint32
	var ebo uint32
	var vao uint32

	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &ebo)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(triangles), gl.Ptr(triangles), gl.STATIC_DRAW)

	// vertex coord attribute
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 32, gl.PtrOffset(0))
	// color coord attribute
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 32, gl.PtrOffset(3*4))
	// texture coord attribute
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 32, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	return vao
}

func loadImage(imagePath string) *image.RGBA {
	imgFile, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("os.Open: ", err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal("image.Decode: ", err)
	}
	rgba := image.NewRGBA(img.Bounds())

	if rgba.Stride != rgba.Rect.Size().X*4 {
		log.Fatal("image.NewRGBA: ", err)
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return rgba
}
