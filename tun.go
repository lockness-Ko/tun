package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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

	if len(os.Args) <= 1 {
		fmt.Print("Usage: tun <command>\n\nQuickly share files with friends using ngrok\n\n")
		fmt.Println("Commands:")
		fmt.Println("\tstart\tStart the tunnel")
		fmt.Println("\tstop\tStop the tunnel")
		fmt.Println("\turl\tGet the URL of the tunnel")
		fmt.Println("\tadd\tAdd a file to be served")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		if _, err := os.Stat(srv + "/running.lock"); err != nil {
			fmt.Println("Starting ngrok")
			runNgrok(srv)
		} else {
			fmt.Println("The tunnel is already running")
		}
	case "stop":
		fmt.Println("Stopping ngrok")
		stopNgrok(srv)
	case "url":
		// Get the URL of the tunnel
		resp, err := http.Get("http://localhost:4040/api/tunnels")
		if err != nil {
			fmt.Println("Tunnel is not started")
			os.Exit(1)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Tunnel is not started")
			os.Exit(1)
		}
		fmt.Println(strings.Split(strings.Split(strings.TrimSpace(string(body)), "public_url\":\"")[1], "\"")[0])
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
