package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

// https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func openInBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		fmt.Println("openning...", url)
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		fmt.Println("openning...", url)
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		fmt.Println("openning...", url)
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
