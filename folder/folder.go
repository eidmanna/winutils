package folder

import (
	"syscall"
	"unsafe"
)

type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

var (
	FOLDERID_Fonts = GUID{0xFD228CB7, 0xAE11, 0x4AE3, [8]byte{0x86, 0x4C, 0x16, 0xF3, 0x91, 0x0A, 0xB8, 0xFE}}
	// {FD228CB7-AE11-4AE3-864C-16F3910AB8FE}
	FOLDERID_Documents = GUID{0xFDD39AD0, 0x238F, 0x46AF, [8]byte{0xAD, 0xB4, 0x6C, 0x85, 0x48, 0x03, 0x69, 0xC7}}
	// {FDD39AD0-238F-46AF-ADB4-6C85480369C7}
	FOLDERID_PublicDocuments = GUID{0xED4824AF, 0xDCE4, 0x45A8, [8]byte{0x81, 0xE2, 0xFC, 0x79, 0x65, 0x08, 0x36, 0x34}}
	// {ED4824AF-DCE4-45A8-81E2-FC7965083634}
)

var (
	modShell32               = syscall.NewLazyDLL("Shell32.dll")
	modOle32                 = syscall.NewLazyDLL("Ole32.dll")
	procSHGetKnownFolderPath = modShell32.NewProc("SHGetKnownFolderPath")
	procCoTaskMemFree        = modOle32.NewProc("CoTaskMemFree")
)

func shGetKnownFolderPath(rfid *GUID, dwFlags uint32, hToken syscall.Handle, pszPath *uintptr) (retval error) {
	r0, _, _ := syscall.Syscall6(procSHGetKnownFolderPath.Addr(), 4, uintptr(unsafe.Pointer(rfid)), uintptr(dwFlags), uintptr(hToken), uintptr(unsafe.Pointer(pszPath)), 0, 0)
	if r0 != 0 {
		retval = syscall.Errno(r0)
	}
	return
}

func coTaskMemFree(pv uintptr) {
	syscall.Syscall(procCoTaskMemFree.Addr(), 1, uintptr(pv), 0, 0)
	return
}

func Folder(rfid *GUID) (string, error) {
	var path uintptr
	err := shGetKnownFolderPath(rfid, 0, 0, &path)
	if err != nil {
		return "", err
	}
	defer coTaskMemFree(path)
	folder := syscall.UTF16ToString((*[1 << 16]uint16)(unsafe.Pointer(path))[:])
	return folder, nil
}
