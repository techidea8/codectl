package winkit

import (
	_ "embed"

	"github.com/getlantern/systray"
)

//go:embed tray.ico
var appIcon []byte

type TrayService struct {
	icon    []byte
	onReady func()
	onExit  func()
}
type Option func(*TrayService)

func UseIcon(icon []byte) Option {
	return func(s *TrayService) {
		s.icon = []byte{}
		s.icon = append(s.icon, icon...)
	}
}

type MenuItem struct {
	Title   string
	Tooltip string
	Attach  interface{}
	OnClick func(menu *MenuItem, item *systray.MenuItem, ch struct{})
}

func NewTrayService(options ...Option) *TrayService {
	result := &TrayService{
		onReady: func() {},
		onExit:  func() {},
	}
	//systray.SetIcon(appIcon)
	for _, v := range options {
		v(result)
	}
	return result
}

func (s *TrayService) AddMenuItem(item *MenuItem) (r *TrayService) {
	menu := systray.AddMenuItem(item.Title, item.Tooltip)
	go func(menu *systray.MenuItem) {
		for r := range menu.ClickedCh {
			item.OnClick(item, menu, r)
		}
	}(menu)
	return s
}

// 托盘程序
func (s *TrayService) OnReady(on func()) {
	s.onReady = on
}
func (s *TrayService) OnExit(on func()) {
	s.onExit = on
}
func (s *TrayService) run() {
	systray.Run(s.onReady, s.onExit)
}
func (s *TrayService) Start() {
	s.run()
}
