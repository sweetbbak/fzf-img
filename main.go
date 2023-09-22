package main

import (
	"bufio"
	"fmt"
	fzf "github.com/ktr0731/go-fuzzyfinder"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
	"strings"
)

func show_image(img string) {
	fmt.Printf(img)

	cc := "bash"
	dash_c := "-c"
	stwing := fmt.Sprintf("--place=%vx%v@%vx0", cols, cols, cols)
	// kit := strings.Join([]string{"kitten icat --transfer-mode=memory --place=30x30@50x0 --stdin=no", img}, " ")
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

var t = get_size()
var cols = int(float64(t.Col) * 0.5)
var rows = int(float64(t.Row) * 0.5)

func main() {

	fmt.Println(rows, cols)
	// -----------------------------------------------------
	// This command works
	// img := exec.Command("bash", "-c", "kitten icat --transfer-mode=memory --stdin=no /home/sweet/Pictures/anime-icons/coffee.jpg > /dev/pts/0")
	// stdout, err := img.Output()
	// if err != nil {
	// 	fmt.Println("err: ", err)
	// }

	// pts, err := os.Open("/dev/tty")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// io.WriteString(pts, string(stdout))
	// -----------------------------------------------------

	a := []string{}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		a = append(a, s.Text())
	}
	fmt.Println()

	idx, err := fzf.Find(
		a,
		func(i int) string {
			return a[i]
		},
		fzf.WithPromptString(">"),
		// fzf.WithPreviewWindow()
		fzf.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}

			go show_image(a[i])
			return ""
		}),
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a[idx])
}