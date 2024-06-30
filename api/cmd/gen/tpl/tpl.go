package tpl

import (
	"embed"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/techidea8/codectl/app/gen/conf"
	"github.com/techidea8/codectl/infra/logger"
)

//go:embed tpl-vue3-go-v2/*
var tpldir embed.FS

func Release() error {
	//读取文件夹dir
	//var filemap map[string][]string = map[string][]string{}
	dirs, err := fs.ReadDir(tpldir, ".")
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		dstdir := filepath.Join(conf.DirTpldata, dir.Name())
		err = os.MkdirAll(dstdir, os.ModeDir)
		fs.WalkDir(tpldir, dir.Name(), func(fpath string, d fs.DirEntry, e error) error {
			if d.IsDir() {
				return nil
			}
			bts, err := fs.ReadFile(tpldir, fpath)
			if err != nil {
				return err
			}
			dstfile := path.Join(conf.DirTpldata, fpath)
			os.RemoveAll(dstfile)
			//logger.Debugf("release tpl file %s", dstfile)
			return os.WriteFile(dstfile, bts, 0644)
		})
	}

	if err != nil {
		logger.Error("Error walking the file system:", err)
		os.Exit(1)
	}
	return nil
}
