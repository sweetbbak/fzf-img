package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	fzf "github.com/ktr0731/go-fuzzyfinder"
	"golang.org/x/sys/unix"
)

func show_image(img string, cols int) {
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

func winch(c chan os.Signal) {
	for {
		s := <-c
		if s == syscall.SIGWINCH {
			fmt.Println("Got signal winch")
			t = get_size()
			cols = int(float64(t.Col) * 0.5)
			rows = int(float64(t.Row) * 0.5)
			break
		}
	}
	t = get_size()
	cols = int(float64(t.Col) * 0.5)
	rows = int(float64(t.Row) * 0.5)
	signal.Notify(c, syscall.SIGWINCH)
}

func set_vars(c chan Termsize) (int, int) {
	var cols = int(float64(t.Col) * 0.5)
	var rows = int(float64(t.Row) * 0.5)
	return cols, rows
}

func monitor_term(ch chan os.Signal) int {
	for {
		<-ch // Wait for the SIGWINCH signal
		t := get_size()
		cols = int(float64(t.Col) * 0.5)
		rows = int(float64(t.Row) * 0.5)
		return cols
	}
}

func logit(text string) {
	filename := "terminal.log"
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	text = fmt.Sprintf("[%s]\n", text)
	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

// get terminal size to guess image size
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

	// -----------------------------------
	t = get_size()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)

	go monitor_term(ch)

	// go func() {
	// 	for {
	// 		s := <-ch
	// 		if s == syscall.SIGWINCH {
	// 			os.WriteFile("winch.log", []byte(s.String()), 0644)
	// 			fmt.Println("Got signal winch")
	// 			t = get_size()
	// 			co <- get_size()
	// 			set_vars(co)
	// 			cols = int(float64(t.Col) * 0.5)
	// 			rows = int(float64(t.Row) * 0.5)
	// 		}
	// 	}
	// 	close(done)
	// }()

	// -----------------------------------

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
			t = get_size()
			cols = int(float64(t.Col) * 0.5)

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
