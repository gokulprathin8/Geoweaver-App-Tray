package main

import (
	"bytes"
	"fmt"
	"geoweaver-systray/utils"
	"github.com/getlantern/systray"
	"image/png"
	"os"
	"path/filepath"
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
		_ = utils.KillGeoweaverProcesses() // no need to handle error here since the process may not be running
		systray.Quit()
	}()

	// handle server start and stop
	go func() {
		for {
			<-startServer.ClickedCh
			fmt.Print(startServer, serverStarted)
			if serverStarted {
				startServer.SetTitle("Start server")
				serverStarted = false
				err := utils.KillGeoweaverProcesses()
				if err != nil {
					panic("Failed to stop geoweaver server")
				}
			} else {
				startServer.SetTitle("Stop server")
				serverStarted = true
				homeDir, _ := os.UserHomeDir()
				err := utils.DownloadFile(utils.GEOWEAVER_JAR_URL, homeDir)
				if err != nil {
					panic("Unable to download file")
				} else {
					geoweaverJarPath := filepath.Join(homeDir, "geoweaver.jar")
					err := utils.RunJavaJar(geoweaverJarPath)
					if err != nil {
						panic("Failed to start geoweaver jar file")
					}
				}
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
}

func onExit() {
	// clean up here
}
