package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
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
		runNgrok(srv)
	}

	if len(os.Args) <= 1 {
		fmt.Println("Usage: tun <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		fmt.Println("Starting ngrok")
		runNgrok(srv)
	case "stop":
		fmt.Println("Stopping ngrok")
		stopNgrok(srv)
	case "":
		os.Exit(1)
	default:
		fmt.Println("Serving file")

		splited := strings.Split(srv+"/"+os.Args[2], "/")

		// Copy the [2]th argument to the srv directory
		copy(os.Args[2], srv+"/"+splited[len(splited)-1])
	}

}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func stopNgrok(srv string) {
	// Delete the running.lock file
	os.Remove(srv + "/running.lock")

	// Kill the ngrok process
	if runtime.GOOS == "windows" {
		exec.Command("taskkill", "/F", "/IM", "ngrok.exe").Run()
	} else {
		exec.Command("pkill", "ngrok").Run()
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
	cmd.Start()
}
