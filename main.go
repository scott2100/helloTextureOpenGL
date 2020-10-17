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
	0.5,  0.5, 0.0,    1.0, 0.0, 0.0,   1.0, 1.0,   // top right
	0.5, -0.5, 0.0,    0.0, 1.0, 0.0,   1.0, 0.0,   // bottom right
	-0.5, -0.5, 0.0,   0.0, 0.0, 1.0,   0.0, 0.0,   // bottom let
	-0.5,  0.5, 0.0,   1.0, 1.0, 0.0,   0.0, 1.0,   // top let
}

func main() {

	runtime.LockOSThread()

	glfw.Init()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, 1)

	window, err := glfw.CreateWindow(640, 480, "Scott Window", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	initOpenGL()

	log.Println("Creating and Compiling Shaders")
	vertexShaderPath := "./shader/shaders/transformShader.vs"
	fragmentShaderPath := "./shader/shaders/transformShader.fs"
	shader := shader.New(vertexShaderPath, fragmentShaderPath)

	vao2 := makeVao(vertices)

	imgFile, err := os.Open("bricks.png")
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

	// load and create a texture
	// -------------------------
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture) // all upcoming GL_TEXTURE_2D operations now have effect on this texture object
	// set the texture wrapping parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)	// set texture wrapping to GL_REPEAT (default wrapping method)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	// set texture filtering parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

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
	gl.GenerateMipmap(gl.TEXTURE_2D)

	prog := gl.CreateProgram()
	gl.AttachShader(prog, shader.VertexShaderCompiled)
	gl.AttachShader(prog, shader.FragmentShaderCompiled)
	gl.LinkProgram(prog)

	//render loop
	for !window.ShouldClose() {
		//process commands
		processInput(window)

		//render commands
		gl.ClearColor(.2, .3, .3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		//gl.UseProgram(prog)
		//gl.BindVertexArray(vao2)
		//

		gl.BindTexture(gl.TEXTURE_2D, texture)

		// render container
		gl.UseProgram(prog)
		gl.BindVertexArray(vao2)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

		glfw.PollEvents()
		window.SwapBuffers()

	}
	glfw.Terminate()
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}

// initOpenGL initializes OpenGL
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
func makeVao(points []float32) uint32 {
	var vbo uint32
	var vao uint32
	var stride int32

	//points only 9
	//points and colors 18
	//points color textures 32
	stride = int32(4 * len(points) / 3)
	println("stride: ", stride)

	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 32, gl.PtrOffset(0))
	log.Println("In if")
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 32, gl.PtrOffset(3*4))
	// texture coord attribute
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 32, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	return vao
}
