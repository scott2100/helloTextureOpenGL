package shader

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"io/ioutil"
	"strings"
)

type shader struct {
	vertexShaderPath   string
	fragmentShaderPath string
	vertexShaderCode   string
	fragmentShaderCode string
	VertexShaderCompiled uint32
	FragmentShaderCompiled uint32
}

func New(vertexPath string, fragmentPath string) shader {
	s := shader{vertexPath, fragmentPath, "","",0,0}
	s.vertexShaderCode = s.readFile(vertexPath)
	s.fragmentShaderCode = s.readFile(fragmentPath)
	vertexShaderCompiled, err := s.compileShader(s.vertexShaderCode, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShaderCompiled, err := s.compileShader(s.fragmentShaderCode, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	s.VertexShaderCompiled = vertexShaderCompiled
	s.FragmentShaderCompiled = fragmentShaderCompiled

	return s
}

func (s shader) readFile(path string) string {
	dat, err := ioutil.ReadFile(path)
	check(err)
	return string(dat)
}

func (s shader) compileShader(source string, shaderType uint32) (uint32, error) {
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

func (s shader) PrintSource() string {
	vertexShaderCode := s.vertexShaderCode
	fragmentShaderCode := s.fragmentShaderCode
	return vertexShaderCode + fragmentShaderCode
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

