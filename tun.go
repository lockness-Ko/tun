package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
)

func main() {
	srv := ""
	switch runtime.GOOS {
	case "windows":
		srv = "C:\\Windows\\Tasks\\tun"
	case "linux":
		srv = "/tmp/tun"
	}

	// If the srv directory doesn't exist, create it
	if _, err := os.Stat(srv); os.IsNotExist(err) {
		os.Mkdir(srv, 0755)
	}

	if _, err := os.Stat(srv + "/running.lock"); err != nil {
		fmt.Println("Starting ngrok")
		runNgrok(srv)
	}
}

func runNgrok(srv string) {
	// Create an empty file to indicate that the server is running
	f, _ := os.Create(srv + "/running.lock")
	f.Close()

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		user, _ := user.Current()
		cmd = exec.Command("cmd", "/C", "start", "C:\\Users\\"+user.Username+"\\ngrok.exe", "http", "-region=au", "file://"+srv)
	case "linux":
		cmd = exec.Command("ngrok", "http", "-region=au", "file://"+srv)
	}
	go cmd.Run()
}
