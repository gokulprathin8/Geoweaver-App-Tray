package main

import (
	"bytes"
	"fmt"
	"geoweaver-systray/utils"
	"github.com/getlantern/systray"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var serverStarted bool = false

func main() {
	systray.Run(onReady, onExit)
}

func loadImageAsByteSlice(filePath string) ([]byte, error) {
	// Open the file specified by the filePath
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Decode the PNG image
	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding PNG: %w", err)
	}

	// Encode the image to PNG format and store it in a byte buffer
	var buffer bytes.Buffer
	err = png.Encode(&buffer, img)
	if err != nil {
		return nil, fmt.Errorf("error encoding PNG to byte slice: %w", err)
	}

	// Return the bytes of the encoded image
	return buffer.Bytes(), nil
}

func onReady() {

	// Load Geoweaver logo
	geoweaverLogo, err := loadImageAsByteSlice("./icon/geoweaver.png")
	if err != nil {
		return
	}
	systray.SetIcon(geoweaverLogo)
	systray.SetTooltip("Geoweaver App for Workflow Management")

	startServer := systray.AddMenuItem("Start Server", "Start Geoweaver server on the background")
	openBrowser := systray.AddMenuItem("Open Browser", "Open geoweaver on your default browser")
	restartServer := systray.AddMenuItem("Restart Server", "Restart Geoweaver server on the background")

	systray.AddSeparator()

	quitMenu := systray.AddMenuItem("Quit", "Quit the whole app")

	// Sets the icon of a menu item. Only available on Mac and Windows.
	browserIcon, err := loadImageAsByteSlice("./icon/browser.png")
	if err != nil {
		return
	}
	openBrowser.SetIcon(browserIcon)

	restartServerIcon, err := loadImageAsByteSlice("./icon/restart.png")
	if err != nil {
		return
	}
	restartServer.SetIcon(restartServerIcon)

	startServerIcon, err := loadImageAsByteSlice("./icon/start-server.png")
	if err != nil {
		return
	}
	startServer.SetIcon(startServerIcon)

	quitIcon, err := loadImageAsByteSlice("./icon/close.png")
	if err != nil {
		return
	}
	quitMenu.SetIcon(quitIcon)

	// handle quit
	go func() {
		<-quitMenu.ClickedCh
		fmt.Print("Requesting quit")
		_ = utils.KillGeoweaverProcesses()
		systray.Quit()
	}()

	// handle server start and stop
	go func() {
		var mutex sync.Mutex

		for {
			<-startServer.ClickedCh

			mutex.Lock()
			currentState := serverStarted
			mutex.Unlock()

			if currentState {
				startServer.SetTitle("Start server")

				err := utils.KillGeoweaverProcesses()
				if err != nil {
					log.Printf("Failed to stop geoweaver server: %v", err)
					continue
				}

				mutex.Lock()
				serverStarted = false
				mutex.Unlock()
			} else {
				startServer.SetTitle("Stop server")

				homeDir, err := os.UserHomeDir()
				if err != nil {
					log.Printf("Failed to get home directory: %v", err)
					continue
				}

				geoweaverJarPath := filepath.Join(homeDir, "geoweaver.jar")
				err = utils.DownloadFile(utils.GEOWEAVER_JAR_URL, geoweaverJarPath)
				if err != nil {
					log.Printf("Unable to download geoweaver.jar: %v", err)
					continue
				}

				err = utils.RunJavaJar(geoweaverJarPath)
				if err != nil {
					log.Printf("Failed to start geoweaver jar file: %v", err)
					continue
				}

				mutex.Lock()
				serverStarted = true
				mutex.Unlock()
			}
		}
	}()

	// handle URL open
	go func() {
		for {
			<-openBrowser.ClickedCh
			err := utils.OpenURLInBrowser("http://localhost:8070/Geoweaver")
			if err != nil {
				panic("Failed to open browser")
			}
		}
	}()

	// handle restart
	go func() {
		for {
			<-restartServer.ClickedCh

			_ = utils.KillGeoweaverProcesses()

			homeDir, _ := os.UserHomeDir()
			geoweaverJarPath := filepath.Join(homeDir, "geoweaver.jar")
			err := utils.RunJavaJar(geoweaverJarPath)
			if err != nil {
				panic("Unable to start geoweaver")
			}
		}
	}()
}

func onExit() {
	// clean up here
}
