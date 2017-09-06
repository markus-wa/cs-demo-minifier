package main

import (
	min "github.com/markus-wa/cs-demo-minifier/csminify"
	"os"
	"runtime"
)

func main() {
	var demPath, outPath string
	if runtime.GOOS == "windows" {
		demPath = "C:\\Dev\\demo.dem"
		outPath = "C:\\Dev\\demo.mrpv2"
	} else {
		demPath = "/home/markus/Downloads/demo.dem"
		outPath = "/home/markus/Downloads/demo.mrpv2"
	}
	f, _ := os.Open(demPath)
	fo, _ := os.Create(outPath)

	min.MinifyTo(f, fo, 0.5)
	f.Close()
}
