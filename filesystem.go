package vz

/*
#cgo darwin CFLAGS: -x objective-c -fno-objc-arc
#cgo darwin LDFLAGS: -lobjc -framework Foundation -framework Virtualization
# include "virtualization.h"
*/
import "C"
import (
	"runtime"
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

// NewVirtioSocketDeviceConfiguration creates a new VirtioSocketDeviceConfiguration.
func NewVirtioFileSystemDeviceConfiguration(hostPath string, readOnly bool, tag string) *VirtioFileSystemDeviceConfiguration {
	hostPathChar := charWithGoString(hostPath)
	tagChar := charWithGoString(tag)

	config := &VirtioFileSystemDeviceConfiguration{
		pointer: pointer{
			ptr: C.newVZVirtioFileSystemDeviceConfiguration(
				hostPathChar.CString(),
				C.bool(readOnly),
				tagChar.CString()),
		},
	}

	runtime.SetFinalizer(config, func(self *VirtioFileSystemDeviceConfiguration) {
		self.Release()
	})

	return config
}
