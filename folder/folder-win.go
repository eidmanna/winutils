// +build windows

package folder

import (
	"syscall"
	"unsafe"
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
