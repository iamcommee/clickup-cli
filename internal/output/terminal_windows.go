//go:build windows

package output

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32                       = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
)

type coord struct {
	X int16
	Y int16
}

type smallRect struct {
	Left   int16
	Top    int16
	Right  int16
	Bottom int16
}

type consoleScreenBufferInfo struct {
	Size              coord
	CursorPosition    coord
	Attributes        uint16
	Window            smallRect
	MaximumWindowSize coord
}

// getTerminalWidth returns the terminal width, or a default if detection fails
func getTerminalWidth() int {
	var info consoleScreenBufferInfo
	handle := os.Stdout.Fd()
	r, _, _ := procGetConsoleScreenBufferInfo.Call(handle, uintptr(unsafe.Pointer(&info)))
	if r == 0 {
		return 80
	}
	width := int(info.Window.Right-info.Window.Left) + 1
	if width <= 0 {
		return 80
	}
	return width
}
