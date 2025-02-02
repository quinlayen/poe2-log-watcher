// //

// package main

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/fsnotify/fsnotify"
// )

// const testFilePath = "/Users/peterfaso/Programming/Go/poe-log-watcher/tests/test.txt" // Change to your actual home directory

// func main() {
// 	// Ensure the file exists
// 	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
// 		log.Fatal("❌ File does not exist:", testFilePath)
// 	}

// 	// Create new watcher
// 	watcher, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		log.Fatal("❌ Error creating file watcher:", err)
// 	}
// 	defer watcher.Close()

// 	// Add file to watcher
// 	err = watcher.Add(testFilePath)
// 	if err != nil {
// 		log.Fatal("❌ Error watching file:", err)
// 	}

// 	fmt.Println("✅ Watching:", testFilePath)

// 	// Listen for events
// 	for {
// 		select {
// 		case event := <-watcher.Events:
// 			fmt.Println("📌 File event detected:", event)
// 		case err := <-watcher.Errors:
// 			fmt.Println("❌ Watcher error:", err)
// 		}
// 	}
// }

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fsnotify/fsnotify"
)

const logFilePath = "C:/Program Files (x86)/Grinding Gear Games/Path of Exile 2/logs/Client.txt"

// var keywords = []string{"your item is ready", "special offer", "new trade message"}

func watchLogFile() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("❌ Error creating file watcher:", err)
	}
	defer watcher.Close()

	err = watcher.Add(logFilePath)
	if err != nil {
		log.Fatal("❌ Error watching file:", err)
	}

	log.Println("✅ Watching:", logFilePath)

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("📌 File updated:", event.Name)
				processFile(event.Name)
			}
		case err := <-watcher.Errors:
			log.Println("❌ Watcher error:", err)
		}
	}
}

// ✅ Reopen the file every time it's updated
func processFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("❌ Failed to open file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("📜 Read line:", line) // Debugging output

		// ✅ Only trigger notifications if '@' is present in the line
		if strings.Contains(line, "@") {
			log.Println("🔔 Match found:", line)
			go sendWindowsNotification(line)
		}
	}
}

// ✅ Windows Notifications
func sendWindowsNotification(message string) {
	cmd := exec.Command("powershell", "-Command", `
	[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null;
	$toastXml = New-Object Windows.Data.Xml.Dom.XmlDocument;
	$toastXml.LoadXml('<toast><visual><binding template="ToastGeneric"><text>PoE2 Alert</text><text>`+message+`</text></binding></visual></toast>');
	$toast = [Windows.UI.Notifications.ToastNotification]::new($toastXml);
	[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("PoE2 Alert").Show($toast);
	`)
	err := cmd.Run()
	if err != nil {
		log.Println("❌ Failed to send Windows notification:", err)
	} else {
		log.Println("✅ Windows notification sent!")
	}
}

func main() {
	watchLogFile()
}
