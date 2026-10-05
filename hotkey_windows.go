//go:build windows

package main

import (
	"runtime"
	"strings"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	hotkeyID        = 1
	wmHotkey        = 0x0312
	modControl      = 0x0002
	modAlt          = 0x0001
	modShift        = 0x0004
	hwndMessage     = ^uintptr(2)
	wmRegisterKey   = 0x8001 // Custom message for registering hotkey
	wmUnregisterKey = 0x8002 // Custom message for unregistering hotkey
)

var (
	user32           = windows.NewLazySystemDLL("user32.dll")
	registerHotKey   = user32.NewProc("RegisterHotKey")
	unregisterHotKey = user32.NewProc("UnregisterHotKey")
	createWindowEx   = user32.NewProc("CreateWindowExW")
	defWindowProc    = user32.NewProc("DefWindowProcW")
	getMessage       = user32.NewProc("GetMessageW")
	translateMessage = user32.NewProc("TranslateMessage")
	dispatchMessage  = user32.NewProc("DispatchMessageW")
	registerClassEx  = user32.NewProc("RegisterClassExW")
	postQuitMessage  = user32.NewProc("PostQuitMessage")
	postMessage      = user32.NewProc("PostMessageW")
)

type WNDCLASSEXW struct {
	cbSize        uint32
	style         uint32
	lpfnWndProc   uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     windows.Handle
	hIcon         windows.Handle
	hCursor       windows.Handle
	hbrBackground windows.Handle
	lpszMenuName  *uint16
	lpszClassName *uint16
	hIconSm       windows.Handle
}

type MSG struct {
	hwnd    uintptr
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      POINT
}

type POINT struct {
	x int32
	y int32
}

var (
	messageWindow  windows.Handle
	hotkeyMu       sync.Mutex
	msgLoopReady   chan struct{}
	quitMsgLoop    chan struct{}
	wndProcPtr     uintptr
	pendingHotkey  *Hotkey
	registerResult bool
	resultReady    chan struct{}
)

func wndProc(hwnd windows.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	if msg == wmHotkey && wParam == hotkeyID {
		hotkeyTriggered()
		return 0
	}
	if msg == wmRegisterKey {
		if pendingHotkey != nil {
			mods := getModifiers(pendingHotkey.Ctrl, pendingHotkey.Alt, pendingHotkey.Shift, pendingHotkey.Cmd)
			vk := getKeyCode(pendingHotkey.Key)
			ret, _, _ := registerHotKey.Call(uintptr(messageWindow), hotkeyID, uintptr(mods), uintptr(vk))
			registerResult = ret != 0
		} else {
			registerResult = false
		}
		if resultReady != nil {
			close(resultReady)
		}
		return 0
	}
	if msg == wmUnregisterKey {
		unregisterHotKey.Call(uintptr(messageWindow), hotkeyID)
		return 0
	}
	ret, _, _ := defWindowProc.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	return ret
}

func createMessageWindow() error {
	className, _ := windows.UTF16PtrFromString("DesktopCalHotkeyClass")
	wndProcPtr = windows.NewCallback(wndProc)

	var wc WNDCLASSEXW
	wc.cbSize = uint32(unsafe.Sizeof(wc))
	wc.lpfnWndProc = wndProcPtr
	wc.hInstance = windows.Handle(getModuleHandle())
	wc.lpszClassName = className

	ret, _, err := registerClassEx.Call(uintptr(unsafe.Pointer(&wc)))
	if ret == 0 {
		return err
	}

	hwnd, _, err := createWindowEx.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		0,
		0,
		0, 0, 0, 0,
		hwndMessage,
		0,
		uintptr(wc.hInstance),
		0,
	)

	if hwnd == 0 {
		return err
	}

	messageWindow = windows.Handle(hwnd)
	return nil
}

func getModuleHandle() uintptr {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	proc := kernel32.NewProc("GetModuleHandleW")
	ret, _, _ := proc.Call(0)
	return ret
}

func getModifiers(ctrl, alt, shift, cmd bool) uint32 {
	var mods uint32
	if ctrl || cmd {
		mods |= modControl
	}
	if alt {
		mods |= modAlt
	}
	if shift {
		mods |= modShift
	}
	return mods
}

func getKeyCode(key string) uint32 {
	if len(key) == 0 {
		return 0
	}
	k := strings.ToLower(key)
	if k[0] >= 'a' && k[0] <= 'z' {
		return uint32(k[0] - 'a' + 'A')
	}
	if k[0] >= '0' && k[0] <= '9' {
		return uint32(k[0])
	}
	return 0
}

func registerGlobalHotkey(hk *Hotkey) bool {
	if hk == nil {
		return false
	}

	hotkeyMu.Lock()

	if quitMsgLoop == nil {
		msgLoopReady = make(chan struct{})
		quitMsgLoop = make(chan struct{})
		go messageLoop()
		hotkeyMu.Unlock()
		<-msgLoopReady
		hotkeyMu.Lock()
	}

	if messageWindow == 0 {
		hotkeyMu.Unlock()
		return false
	}

	// Post message to register hotkey on the message loop thread
	pendingHotkey = hk
	resultReady = make(chan struct{})
	
	ret, _, _ := postMessage.Call(
		uintptr(messageWindow),
		wmRegisterKey,
		0,
		0,
	)
	
	if ret == 0 {
		hotkeyMu.Unlock()
		return false
	}
	
	// Wait for the result from the message loop thread
	hotkeyMu.Unlock()
	<-resultReady
	hotkeyMu.Lock()
	
	result := registerResult
	pendingHotkey = nil
	resultReady = nil
	hotkeyMu.Unlock()

	return result
}

func unregisterGlobalHotkey() {
	hotkeyMu.Lock()
	defer hotkeyMu.Unlock()

	if messageWindow != 0 {
		postMessage.Call(
			uintptr(messageWindow),
			wmUnregisterKey,
			0,
			0,
		)
	}
}

func messageLoop() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := createMessageWindow(); err != nil {
		close(msgLoopReady)
		return
	}

	close(msgLoopReady)

	var m MSG
	for {
		ret, _, _ := getMessage.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if ret == 0 || ret == ^uintptr(0) {
			return
		}
		translateMessage.Call(uintptr(unsafe.Pointer(&m)))
		dispatchMessage.Call(uintptr(unsafe.Pointer(&m)))
	}
}

func hotkeyTriggered() {
	if appInstance != nil && appInstance.mainWin != nil {
		if appInstance.mainWin.IsVisible() {
			appInstance.mainWin.Hide()
			setDockIconVisible(false)
		} else {
			appInstance.mainWin.Show()
			appInstance.mainWin.Focus()
			setDockIconVisible(true)
		}
	}
}
