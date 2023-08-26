package install

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/afero"
)

func InstallSelfAs(ctx context.Context, afos afero.Fs, fls afero.Fs, path string, name string) error {

	fle, err := fls.Open(path)
	if err != nil {
		return err
	}

	stat, err := fle.Stat()
	if err != nil {
		return err
	}

	switch runtime.GOOS {
	case "darwin", "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err // or handle the error as you see fit
		}

		nameDir := filepath.Join(homeDir, "."+name)

		// create a ~/.name directory
		if err := afos.MkdirAll(nameDir, 0755); err != nil {
			return err
		}

		if stat.IsDir() {
			res, err := fle.Readdir(-1)
			if err != nil {
				return err
			}
			for _, f := range res {
				flee, err := fls.Open(filepath.Join(path, f.Name()))
				if err != nil {
					return err
				}
				defer flee.Close()
				err = afero.WriteReader(afos, filepath.Join(nameDir, f.Name()), flee)
				if err != nil {
					return err
				}
			}
		} else {
			flee, err := fls.Open(path)
			if err != nil {
				return err
			}

			defer flee.Close()

			err = afero.WriteReader(afos, filepath.Join(nameDir, flee.Name()), flee)
			if err != nil {
				return err
			}

		}

		// check if name is in path
		// if not, add it
		ptf := os.Getenv("PATH")
		if !strings.Contains(ptf, nameDir) {
			if err := os.Setenv("PATH", nameDir+":"+ptf); err != nil {
				return err
			}
		}

		fmt.Println("installed name to " + filepath.Join(nameDir, name))

	case "windows":
		fmt.Println("installing for windows")
		fmt.Println("")

		var ref string

		if stat.IsDir() {
			res, err := fle.Readdir(-1)
			if err != nil {
				return err
			}

			for _, f := range res {
				if f.Name() == name+".exe" {
					ref = filepath.Join(path, stat.Name(), f.Name())
				}
			}

		} else {
			if stat.Name() == name+".exe" {
				ref = filepath.Join(path, stat.Name())
			}
		}

		if ref == "" {
			return fmt.Errorf("could not find %s.exe in %s", name, path)
		}

		flee, err := fls.Open(ref)
		if err != nil {
			return err
		}
		defer flee.Close()
		err = afero.WriteReader(afos, "$LOCALAPPDATA\\Microsoft\\WindowsApps\\"+name+".exe", flee)
		if err != nil {
			return err
		}

		fmt.Println("installed og to $LOCALAPPDATA\\Microsoft\\WindowsApps\\" + name + ".exe")

	default:
		fmt.Println("unsupported platform")
	}

	return nil

}
