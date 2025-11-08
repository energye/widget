package util

import (
	"github.com/energye/lcl/api/libname"
	"github.com/energye/lcl/tool"
	"os"
	"path/filepath"
)

func TestLoadLibPath() {
	var (
		name  string
		wd, _ = os.Getwd()
	)
	if tool.IsWindows() {
		name = "liblcl.dll"
	} else if tool.IsLinux() {
		name = "liblcl.so"
	} else if tool.IsDarwin() {
		name = "liblcl.dylib"
	}
	if name != "" {
		// 当前目录
		liblcl := filepath.Join(wd, name)
		if tool.IsExist(liblcl) {
			libname.LibName = liblcl
			return
		}
		// 测试编译输出目录
		if tool.IsWindows() {
			liblcl = filepath.Join("E:\\SWT\\gopath\\src\\github.com\\energye\\workspace\\gen\\gout", name)
		} else if tool.IsLinux() {
			liblcl = filepath.Join("/home/yanghy/app/projects/workspace/gen/gout", name)
		} else if tool.IsDarwin() {
			liblcl = filepath.Join("/Users/yanghy/app/workspace/gen/gout", name)
		}
		if tool.IsExist(liblcl) {
			libname.LibName = liblcl
			return
		}
	}
}
