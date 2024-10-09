package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procSetWindowsHookExW   = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procGetMessageW         = user32.NewProc("GetMessageW")
	procSendInput           = user32.NewProc("SendInput")
)

const (
	WH_KEYBOARD_LL = 13
	WM_KEYDOWN     = 0x0100
	WM_KEYUP       = 0x0101
	VK_A           = 0x41
	VK_D           = 0x44
	VK_W           = 0x57
	VK_S           = 0x53
)

var (
	aHeld, dHeld, wHeld, sHeld     bool
	aScrip, dScrip, wScrip, sScrip bool
	hookID                         uintptr
)

// KBDLLHOOKSTRUCT is the structure used by the Windows API to store keyboard hook information.
type KBDLLHOOKSTRUCT struct {
	VKCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

// INPUT and KEYBDINPUT structures for SendInput API
type INPUT struct {
	Type uint32
	Ki   KEYBDINPUT
}

type KEYBDINPUT struct {
	WVk         uint16
	WScan       uint16
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

func main() {
	fmt.Println("Starting low-level keyboard control...")

	// Set up a low-level keyboard hook
	hookID = SetWindowsHookEx(WH_KEYBOARD_LL, syscall.NewCallback(keyboardProc), 0, 0)
	if hookID == 0 {
		fmt.Println("Failed to install hook!")
		return
	}
	defer UnhookWindowsHookEx(hookID)

	// Message loop to keep the hook active
	var msg struct {
		hwnd    uintptr
		message uint32
		wParam  uintptr
		lParam  uintptr
		time    uint32
		pt      struct{ x, y int32 }
	}
	for {
		procGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
	}
}

// keyboardProc is the hook procedure called by the system when a keyboard event is caught
func keyboardProc(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode == 0 { // 0 means the hook procedure must process the message
		kbdStruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		switch wParam {
		case WM_KEYDOWN:
			handleKeyPress(kbdStruct.VKCode)
		case WM_KEYUP:
			handleKeyRelease(kbdStruct.VKCode)
		}
	}
	// Call the next hook in the chain
	ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
	return ret
}

// handleKeyPress handles the press event of the specified virtual key code
func handleKeyPress(vkCode uint32) {
	switch vkCode {
	case VK_A:
		if !aHeld { // Avoid duplicate keydown events
			aHeld = true
			if dScrip {
				dScrip = false
				sendKey(VK_D, "up")
			}
			aScrip = true
			sendKey(VK_A, "down")
		}
	case VK_D:
		if !dHeld {
			dHeld = true
			if aScrip {
				aScrip = false
				sendKey(VK_A, "up")
			}
			dScrip = true
			sendKey(VK_D, "down")
		}
	case VK_W:
		if !wHeld {
			wHeld = true
			if sScrip {
				sScrip = false
				sendKey(VK_S, "up")
			}
			wScrip = true
			sendKey(VK_W, "down")
		}
	case VK_S:
		if !sHeld {
			sHeld = true
			if wScrip {
				wScrip = false
				sendKey(VK_W, "up")
			}
			sScrip = true
			sendKey(VK_S, "down")
		}
	}
}

// handleKeyRelease handles the release event of the specified virtual key code
func handleKeyRelease(vkCode uint32) {
	switch vkCode {
	case VK_A:
		if aHeld {
			aHeld = false
			aScrip = false
			sendKey(VK_A, "up")
			if dHeld && !dScrip {
				dScrip = true
				sendKey(VK_D, "down")
			}
		}
	case VK_D:
		if dHeld {
			dHeld = false
			dScrip = false
			sendKey(VK_D, "up")
			if aHeld && !aScrip {
				aScrip = true
				sendKey(VK_A, "down")
			}
		}
	case VK_W:
		if wHeld {
			wHeld = false
			wScrip = false
			sendKey(VK_W, "up")
			if sHeld && !sScrip {
				sScrip = true
				sendKey(VK_S, "down")
			}
		}
	case VK_S:
		if sHeld {
			sHeld = false
			sScrip = false
			sendKey(VK_S, "up")
			if wHeld && !wScrip {
				wScrip = true
				sendKey(VK_W, "down")
			}
		}
	}
}

// sendKey sends a simulated key event using SendInput
func sendKey(vkCode uint32, action string) {
	var input INPUT
	input.Type = 1 // INPUT_KEYBOARD
	input.Ki.WVk = uint16(vkCode)

	if action == "up" {
		input.Ki.DwFlags = 2 // KEYEVENTF_KEYUP
	}

	procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
}

// SetWindowsHookEx sets a low-level keyboard hook
func SetWindowsHookEx(idHook int, lpfn uintptr, hmod uintptr, dwThreadId uint32) uintptr {
	ret, _, _ := procSetWindowsHookExW.Call(uintptr(idHook), lpfn, hmod, uintptr(dwThreadId))
	return ret
}

// UnhookWindowsHookEx removes a previously set hook
func UnhookWindowsHookEx(hhk uintptr) {
	procUnhookWindowsHookEx.Call(hhk)
}
