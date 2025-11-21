package util

import (
	"github.com/energye/lcl/api/libname"
	"github.com/energye/lcl/tool"
	"os"
	"path/filepath"
)

func TestLoadLibPath() {
	var (
		name  = libname.GetDLLName()
		wd, _ = os.Getwd()
	)
	if name != "" {
		// 当前目录
		liblcl := filepath.Join(wd, name)
		if tool.IsExist(liblcl) {
			libname.LibName = liblcl
			return
		}
		// 测试编译输出目录
		if tool.IsWindows() {
			liblcl = filepath.Join("C:\\app\\workspace\\gen\\gout", name)
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
