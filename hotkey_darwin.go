package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Carbon

#include <Cocoa/Cocoa.h>
#include <Carbon/Carbon.h>

static EventHotKeyRef hotKeyRef = NULL;

extern void hotkeyTriggered();

static OSStatus hotkeyHandler(EventHandlerCallRef nextHandler, EventRef event, void *userData) {
	hotkeyTriggered();
	return noErr;
}

static void registerHotKeyNative(unsigned int keyCode, unsigned int modifiers) {
	if (hotKeyRef != NULL) {
		UnregisterEventHotKey(hotKeyRef);
		hotKeyRef = NULL;
	}
	
	EventHotKeyID hotKeyID;
	hotKeyID.signature = 'dstc';
	hotKeyID.id = 1;
	
	EventTypeSpec eventType;
	eventType.eventClass = kEventClassKeyboard;
	eventType.eventKind = kEventHotKeyPressed;
	
	InstallApplicationEventHandler(
		NewEventHandlerUPP(hotkeyHandler),
		1,
		&eventType,
		NULL,
		NULL
	);
	
	RegisterEventHotKey(
		keyCode,
		modifiers,
		hotKeyID,
		GetApplicationEventTarget(),
		0,
		&hotKeyRef
	);
}

static void unregisterHotKeyNative() {
	if (hotKeyRef != NULL) {
		UnregisterEventHotKey(hotKeyRef);
		hotKeyRef = NULL;
	}
}
*/
import "C"

import (
	"strings"
)

// macOS 虚拟键码映射
func keyCodeFor(key string) uint {
	key = strings.ToLower(key)
	keyCodes := map[string]uint{
		"a": 0x00, "s": 0x01, "d": 0x02, "f": 0x03, "h": 0x04,
		"g": 0x05, "z": 0x06, "x": 0x07, "c": 0x08, "v": 0x09,
		"b": 0x0B, "q": 0x0C, "w": 0x0D, "e": 0x0E, "r": 0x0F,
		"y": 0x10, "t": 0x11, "1": 0x12, "2": 0x13, "3": 0x14,
		"4": 0x15, "6": 0x16, "5": 0x17, "9": 0x19, "7": 0x1A,
		"8": 0x1C, "0": 0x1D, "o": 0x1F, "u": 0x20, "i": 0x22,
		"p": 0x23, "l": 0x25, "j": 0x26, "k": 0x28, "n": 0x2D,
		"m": 0x2E,
	}
	if code, ok := keyCodes[key]; ok {
		return code
	}
	return 0x08 // default to 'c'
}

// macOS 修饰键映射
func modifiersFor(hk *Hotkey) uint {
	var mods uint = 0
	if hk.Cmd {
		mods |= 0x100 // cmdKey
	}
	if hk.Shift {
		mods |= 0x200 // shiftKey
	}
	if hk.Alt {
		mods |= 0x400 // optionKey
	}
	if hk.Ctrl {
		mods |= 0x800 // controlKey
	}
	return mods
}

func registerGlobalHotkey(hk *Hotkey) bool {
	if hk == nil {
		return false
	}
	keyCode := keyCodeFor(hk.Key)
	modifiers := modifiersFor(hk)
	C.registerHotKeyNative(C.uint(keyCode), C.uint(modifiers))
	return true
}

func unregisterGlobalHotkey() {
	C.unregisterHotKeyNative()
}

//export hotkeyTriggered
func hotkeyTriggered() {
	if appInstance != nil && appInstance.mainWin != nil {
		if appInstance.mainWin.IsVisible() {
			appInstance.mainWin.Hide()
		} else {
			appInstance.mainWin.Show()
			appInstance.mainWin.Focus()
		}
	}
}
