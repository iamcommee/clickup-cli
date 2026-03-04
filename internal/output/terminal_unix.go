//go:build !windows

package output

import (
	"os"
	"syscall"
	"unsafe"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// getTerminalWidth returns the terminal width, or a default if detection fails
func getTerminalWidth() int {
	ws := &winsize{}
	fd := os.Stdout.Fd()
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(ws)))
	if err != 0 || ws.Col == 0 {
		return 80
	}
	return int(ws.Col)
}
