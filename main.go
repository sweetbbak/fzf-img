package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	fzf "github.com/ktr0731/go-fuzzyfinder"
	"golang.org/x/sys/unix"
)

func show_image(img string) {
	fmt.Printf(img)

	cc := "bash"
	dash_c := "-c"
	stwing := fmt.Sprintf("--place=%vx%v@%vx0", cols, cols, cols)
	kit := strings.Join([]string{"kitten icat --transfer-mode=memory --stdin=no", stwing, img}, " ")
	im := exec.Command(cc, dash_c, kit)
	stdout, err := im.Output()
	if err != nil {
		fmt.Println("err: ", err)
	}

	fmt.Printf("%s", stdout)
}

type Termsize struct {
	Col uint16
	Row uint16
}

func get_size() Termsize {
	var err error
	var f *os.File
	var t Termsize
	if f, err = os.OpenFile("/dev/tty", unix.O_NOCTTY|unix.O_CLOEXEC|unix.O_NDELAY|unix.O_RDWR, 0666); err == nil {
		var sz *unix.Winsize
		if sz, err = unix.IoctlGetWinsize(int(f.Fd()), unix.TIOCGWINSZ); err == nil {
			// fmt.Printf("rows: %v columns: %v width: %v height %v\n", sz.Row, sz.Col, sz.Xpixel, sz.Ypixel)
			t = Termsize{sz.Col, sz.Row}
			return t
		}
	}

	fmt.Fprintln(os.Stderr, err)
	return t
}

func pause(pid int) {
	err := syscall.Kill(pid, syscall.SIGSTOP)
	if err != nil {
		fmt.Println(err)
	}

}

func System(cmd string) int {
	c := exec.Command("sh", "-c", cmd)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()

	if err == nil {
		return 0
	}

	// Figure out the exit code
	if ws, ok := c.ProcessState.Sys().(syscall.WaitStatus); ok {
		if ws.Exited() {
			return ws.ExitStatus()
		}

		if ws.Signaled() {
			return -int(ws.Signal())
		}
	}

	return -1
}

func is_stdin_open() bool {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// fmt.Println("data is being piped to stdin")
		return true
	} else {
		// fmt.Println("stdin is from a terminal")
		return false
	}
}
func find_images(root string) []string {
	// array := make(map[string][]string)
	// use this because idk what that make map shit was lol
	var array []string
	filepath.WalkDir(root, func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !file.IsDir() {
			ext := filepath.Ext(path)
			// match file extension
			switch ext {
			case ".jpg":
				array = append(array, path)
			case ".png":
				array = append(array, path)
			case ".webp":
				array = append(array, path)
			case ".gif":
				array = append(array, path)
			}
		}
		return nil
	})
	// return nil
	return array
}

var t = get_size()
var cols = int(float64(t.Col) * 0.5)
var rows = int(float64(t.Row) * 0.5)

func main() {
	a := []string{}
	if is_stdin_open() == true {
		// fmt.Println("STDIN OPEN")
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			a = append(a, s.Text())
		}
	} else {
		// fmt.Println("STDIN IS NOT OPEN")
		a = find_images(os.Getenv("HOME"))
	}

	idx, err := fzf.Find(
		a,
		func(i int) string {
			return a[i]
		},
		fzf.WithPromptString(">"),
		fzf.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}

			// go show_image(a[i])
			stwing := fmt.Sprintf("--place=%vx%v@%vx0", cols, cols, cols)
			kit := strings.Join([]string{"kitten icat --transfer-mode=memory --clear --stdin=no", stwing, a[i]}, " ")
			go System(kit)
			return ""
		}),
	)
	if err == fzf.ErrAbort {
		os.Exit(0)
	}

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a[idx])
}
