// +build linux,!cgo

package folder

func Folder(rfid *GUID) (string, error) {
	return "pub", nil
}
