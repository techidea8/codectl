package configure

import "strings"

type Naming struct {
	Namespace string
	Delim     string
	Env       string
	AppName   string
	Ext       string
}

func NewNaming(appName string) *Naming {
	tmp := &Naming{
		Namespace: "truing",
		Delim:     "-",
		Env:       "",
		AppName:   appName,
	}
	return tmp
}
func (n *Naming) Build(ext string) string {
	defaultfile := n.AppName + "." + ext
	arr := []string{}
	if n.Namespace != "" {
		arr = append(arr, n.Namespace)
	}
	if n.Env != "" {
		arr = append(arr, n.Env)
	}
	arr = append(arr, defaultfile)
	return strings.Join(arr, n.Delim)
}
func (n *Naming) Json() string {
	return n.Build("json")
}
func (n *Naming) Yaml() string {
	return n.Build("yml")
}
func (n *Naming) Xml() string {
	return n.Build("xml")
}
func (n *Naming) SetNameSpace(ns string) *Naming {
	n.Namespace = ns
	return n
}
func (n *Naming) SetDelim(dlm string) *Naming {
	n.Delim = dlm
	return n
}

func (n *Naming) UseEnv(env string) *Naming {
	n.Env = env
	return n
}

func (n *Naming) SetAppName(appName string) *Naming {
	n.AppName = appName
	return n
}
