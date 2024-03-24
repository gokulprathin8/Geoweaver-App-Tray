package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const GEOWEAVER_JAR_URL string = "https://github.com/ESIPFed/Geoweaver/releases/download/latest/geoweaver.jar"

func KillGeoweaverProcesses() error {
	cmd := exec.Command("pgrep", "java")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run pgrep: %v", err)
	}
	pids := strings.Fields(out.String())
	if len(pids) == 0 {
		fmt.Println("No Java processes found.")
		return nil
	}
	for _, pid := range pids {
		err := exec.Command("kill", pid).Run()
		if err != nil {
			return fmt.Errorf("failed to kill process %s: %v", pid, err)
		}
	}
	return nil
}

func DownloadFile(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func RunJavaJar(jarPath string) error {
	cmd := exec.Command("java", "-jar", jarPath)
	cmd.Dir = filepath.Dir(jarPath)
	return cmd.Run()
}

func OpenURLInBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = exec.Command("xdg-open", url).Start() // Fallback for other UNIX-like OSes
	}
	return err
}
