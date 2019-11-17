package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
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

	for !win.ShouldClose() {
		draw()
		win.SwapBuffers()
		glfw.PollEvents()
	}
}

func draw() {
	gl.ClearColor(0, 0.5, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}
