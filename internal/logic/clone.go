package logic

import "techidea8.com/codectl/internal/util"

func Clone(fpath string) (str string, err error) {
	return util.ExecQuit("git", "clone", fpath)
}
