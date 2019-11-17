// https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl
package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	vertexShaderSource = `
    #version 410
    in vec3 vp;
    void main() {
      gl_Position = vec4(vp, 1.0);
    }
  ` + "\x00"

	fragmentShaderSource = `
    #version 410
    out vec4 frag_color;
    void main() {
      frag_color = vec4(1, 1, 1, 1);
    }
  ` + "\x00"
)

func init() {
	runtime.LockOSThread()
}

func initGlfw(width, height int) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	win, err := glfw.CreateWindow(width, height, "Blue Sky", nil, nil)

	if err != nil {
		panic(fmt.Errorf("could not create opengl renderer: %v", err))
	}

	//MakeContextCurrent is what actually creates the OpenGL context
	//within the platform window that gets created by GLFW.
	//Now, we’re ready to open the window.
	win.MakeContextCurrent()

	return win
}

func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
}

func main() {
	win := initGlfw(800, 600)
	defer glfw.Terminate()

	initOpenGL()

	triangle := makeTriangle()
	vao := makeVao(triangle)

	program := makeProgram(vertexShaderSource, fragmentShaderSource)

	for !win.ShouldClose() {
		draw(vao, int32(len(triangle)), program)
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func draw(vao uint32, len int32, program uint32) {
	gl.ClearColor(0, 0.5, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, len/3)
}

func makeTriangle() []float32 {
	return []float32{
		0, 0.5, 0, // top
		-0.5, -0.5, 0, // left
		0.5, -0.5, 0, // right
	}
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func makeProgram(vertex_s, frag_s string) uint32 {
	program := gl.CreateProgram()

	vertexShader, err := compileShader(vertex_s, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(frag_s, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	return program
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
