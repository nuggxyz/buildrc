package install

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func InstallSelfAs(ctx context.Context, name string) error {

	switch runtime.GOOS {
	case "darwin", "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err // or handle the error as you see fit
		}

		nameDir := filepath.Join(homeDir, "."+name)

		// create a ~/.name directory
		if err := os.MkdirAll(nameDir, 0755); err != nil {
			return err
		}

		// move self to ~/.name/name
		if err := os.Rename(os.Args[0], filepath.Join(nameDir, name)); err != nil {
			return err
		}

		// check if name is in path
		// if not, add it
		path := os.Getenv("PATH")
		if !strings.Contains(path, nameDir) {
			if err := os.Setenv("PATH", nameDir+":"+path); err != nil {
				return err
			}
		}

		fmt.Println("installed name to " + filepath.Join(nameDir, name))

	case "windows":
		fmt.Println("installing for windows")
		fmt.Println("")
		// move self to $LOCALAPPDATA and then suffixed with \Microsoft\WindowsApps\og.exe"
		if err := os.Rename(os.Args[0], "$LOCALAPPDATA\\Microsoft\\WindowsApps\\"+name+".exe"); err != nil {
			return err
		}

		fmt.Println("installed og to $LOCALAPPDATA\\Microsoft\\WindowsApps\\" + name + ".exe")

	default:
		fmt.Println("unsupported platform")
	}

	return nil

}
