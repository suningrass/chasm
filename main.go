package main

import (
	"bufio"
	prompt "chasm/components"
	"fmt"
	"io"
	"os"
	"runtime"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-tty"
	"github.com/muesli/termenv"
	"golang.design/x/clipboard"
)

func main() {

	/* // 1. 判断是否有管道数据
	fi, _ := os.Stdin.Stat()
	isPipe := (fi.Mode()&os.ModeNamedPipe != 0) || (fi.Mode()&os.ModeCharDevice == 0)

	if isPipe {
		// 2. 有管道 → 直接透传
		_, _ = io.Copy(os.Stdout, os.Stdin)
		return
	}
	*/
	//读入
	rawline := readRawLines()

	mypmpt := prompt.InitialModel()
	var linesStr []string
	for _, b := range rawline {
		if utf8.Valid(b) {
			linesStr = append(linesStr, string(b))
		}
	}
	mypmpt.Argchoose.ArgList = linesStr
	mypmpt.Argchoose.SetInput()

	// 打开 TTY
	in, out, close, err := openTTY()
	if err != nil {
		fmt.Fprintln(os.Stderr, "tty:", err)
		os.Exit(1)
	}
	defer close()

	p := tea.NewProgram(mypmpt, tea.WithAltScreen(), tea.WithInput(in), tea.WithOutput(out))

	mInterface, err := p.Run()
	if err != nil {
		fmt.Printf("bubble program error: %v\n", err)
		os.Exit(1)
	}

	mypmpt = mInterface.(prompt.Promptmodel)

	// 初始化一次即可
	clipboard.Init()
	// 把结果字符串写入系统剪贴板（UTF-8）
	clipboard.Write(clipboard.FmtText, []byte(mypmpt.Input.Value()))

	//fmt.Println(mypmpt.Input.Value())
	fmt.Fprintln(os.Stdout, mypmpt.Input.Value())

}

// readRawLines 返回 [][]byte，保留原始字节
func readRawLines() [][]byte {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return nil // 无管道
	}

	var lines [][]byte
	r := bufio.NewReader(os.Stdin)
	for {
		line, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			// 只保留已读内容
			break
		}
		// 去掉末尾换行符（可选）
		if len(line) > 0 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
		if len(line) > 0 {
			lines = append(lines, line) // 原样追加
		}
		if err == io.EOF {
			break
		}
	}
	return lines
}

// openTTY 返回 (输入, 输出, 关闭函数)
func openTTY() (in io.Reader, out io.Writer, close func() error, err error) {
	if runtime.GOOS == "windows" {
		t, err := tty.Open()
		if err != nil {
			return nil, nil, nil, err
		}
		// Windows 用 go-tty 提供的句柄
		in, out = t.Input(), t.Output()
		close = t.Close
		lipgloss.SetColorProfile(termenv.NewOutput(out).Profile)
		return in, out, close, nil
	}
	// 类 Unix：/dev/tty
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return os.Stdin, os.Stdout, func() error { return nil }, nil // 降级
	}
	in, out = tty, tty
	close = tty.Close
	lipgloss.SetColorProfile(termenv.NewOutput(out).Profile)
	return in, out, close, nil
}
