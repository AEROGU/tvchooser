package mounted

import (
	"syscall"
)

// Get a list of mounted devices on Windows.
// like "wmic logicaldisk get name"
func GetWindowsDriveLetters() ([]string, error) {
	kernel32, err := syscall.LoadLibrary("kernel32.dll")
	defer syscall.FreeLibrary(kernel32)

	if err != nil {
		return nil, err
	}

	getLogicalDrivesHandle, err := syscall.GetProcAddress(kernel32, "GetLogicalDrives")

	if err != nil {
		return nil, err
	}

	var drives []string

	// if ret, _, callErr := syscall.Syscall(uintptr(getLogicalDrivesHandle), 0, 0, 0, 0); callErr != 0 { // deprecated
	if ret, _, callErr := syscall.SyscallN(uintptr(getLogicalDrivesHandle), 0, 0, 0, 0); callErr != 0 {
		// handle error
	} else {
		drives = bitsToDrives(uint32(ret))
	}

	return drives, nil
}

func bitsToDrives(bitMap uint32) (drives []string) {
	availableDrives := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	for i := range availableDrives {
		if bitMap&1 == 1 {
			drives = append(drives, availableDrives[i])
		}
		bitMap >>= 1
	}

	return
}
