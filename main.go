package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fsnotify/fsnotify"
)

const logFilePath = "Path to POE2 log.config"

var keywords = []string{"slice of things to look for in log file"}

var currentOS = runtime.GOOS

func watchLogFiles() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Error creating watcher:", err)
	}
	defer watcher.Close()

	file, err := os.Open(logFilePath)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer file.Close()

	file.Seek(0, io.SeekEnd)

	log.Println("Monitory PoE2 log file for incoming whispers...")

	err = watcher.Add(logFilePath)
	if err != nil {
		log.Fatal("Erro adding file to watcher:", err)
	}

	scanner := bufio.NewScanner(file)

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				for scanner.Scan() {
					line := scanner.Text()
					for _, keyword := range keywords {
						if strings.Contains(strings.ToLower(line),
							strings.ToLower(keyword)) {
							log.Println("🔔 Match found:", line)
							go sendNotification(line)
						}
					}
				}
			}
		}
	}

}

func sendNotification(message string) {
	switch currentOS {
	case "darwin":
		exec.Command("osascript", "-e", `dispaly notification "`+message+`" with title "Poe2 Alert"`).Run()
	case "windows":
		exec.Command("powershell", "-Command", `[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("PoE2 Alert").Show((New-Object Windows.UI.Notifications.ToastNotification (New-Object Windows.Data.Xml.Dom.XmlDocument)))`).Run()
	default:
		fmt.Println("🔔 PoE2 Alert:", message)
	}
}

func main() {
	watchLogFiles()
}
