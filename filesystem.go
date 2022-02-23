package vz

/*
#cgo darwin CFLAGS: -x objective-c -fno-objc-arc
#cgo darwin LDFLAGS: -lobjc -framework Foundation -framework Virtualization
# include "virtualization.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type DirectorySharingDeviceConfiguration interface {
	NSObject

	directorySharingDeviceConfiguration()
}

type baseDirectorySharingDeviceConfiguration struct{}

func (*baseDirectorySharingDeviceConfiguration) directorySharingDeviceConfiguration() {}

type VirtioFileSystemDeviceConfiguration struct {
	pointer

	*baseDirectorySharingDeviceConfiguration
}

type SharedDirectory struct {
	HostPath  string
	GuestPath string
	ReadOnly  bool
}

// NewVirtioSocketDeviceConfiguration creates a new VirtioSocketDeviceConfiguration.
func NewVirtioFileSystemDeviceConfiguration(directories []SharedDirectory, tag string) *VirtioFileSystemDeviceConfiguration {
	tagChar := charWithGoString(tag)

	var dirList []*C.SharedDirectory

	for _, d := range directories {
		dir := &C.SharedDirectory{}
		dir = (*C.SharedDirectory)(C.malloc(C.size_t(unsafe.Sizeof(C.SharedDirectory{}))))

		dir.hostpath = charWithGoString(d.HostPath).CString()
		dir.guestpath = charWithGoString(d.GuestPath).CString()
		dir.readonly = C.bool(d.ReadOnly)
		dirList = append(dirList, dir)
	}

	sharedDirs := C.MultipleSharedDirectory{}
	if len(directories) > 0 {
		sharedDirs.directories = &dirList[0]
	}
	sharedDirs.n_directories = C.uint(len(directories))

	config := &VirtioFileSystemDeviceConfiguration{
		pointer: pointer{
			ptr: C.newVZVirtioFileSystemDeviceConfiguration(
				sharedDirs,
				tagChar.CString()),
		},
	}

	runtime.SetFinalizer(config, func(self *VirtioFileSystemDeviceConfiguration) {
		self.Release()
	})

	return config
}

type VirtioDirectorySharingDevice struct {
	pointer
}

func newVirtioDirectorySharingDevice(ptr unsafe.Pointer) *VirtioDirectorySharingDevice {
	directorySharingDevice := &VirtioDirectorySharingDevice{
		pointer: pointer{
			ptr: ptr,
		},
	}

	return directorySharingDevice
}

// SetSharedDirectories updates the list of directories shared by a virtiofs device.
func (d *VirtioDirectorySharingDevice) SetSharedDirectories(directories []SharedDirectory) {
	var dirList []*C.SharedDirectory

	for _, d := range directories {
		dir := &C.SharedDirectory{}
		dir = (*C.SharedDirectory)(C.malloc(C.size_t(unsafe.Sizeof(C.SharedDirectory{}))))

		dir.hostpath = charWithGoString(d.HostPath).CString()
		dir.guestpath = charWithGoString(d.GuestPath).CString()
		dir.readonly = C.bool(d.ReadOnly)
		dirList = append(dirList, dir)
	}

	sharedDirs := C.MultipleSharedDirectory{}
	sharedDirs.directories = &dirList[0]
	sharedDirs.n_directories = C.uint(len(directories))

	C.VZVirtioFileSystemDevice_setShare(d.Ptr(), sharedDirs)
}
