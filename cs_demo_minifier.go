package main

import (
	min "github.com/markus-wa/cs-demo-minifier/csminify"
	"os"
	"runtime"
)

func main() {
	var demPath string
	if runtime.GOOS == "windows" {
		demPath = "C:\\Dev\\demo.dem"
	} else {
		demPath = "/home/markus/Downloads/demo.dem"
	}
	f, _ := os.Open(demPath)

	min.MinifyTo(f, os.Stdout)
}
