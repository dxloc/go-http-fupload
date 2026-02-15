package main

import "go-http-fupload/cmd"

func init() {
	cmd.Version = ""
	cmd.Init()
}

func main() {
	cmd.Main()
}
