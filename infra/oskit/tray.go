package oskit

import (
	_ "embed"

	"github.com/getlantern/systray"
	"github.com/techidea8/codectl/infra/logger"
)

//go:embed tray.ico
var appIcon []byte

type TrayService struct {
	icon    []byte
	tooltip string
	title   string
	onReady func()
	onExit  func()
	menus   []*MenuItem
	logger  logger.ILogger
}
type Option func(*TrayService)

func UseIcon(icon []byte) Option {
	return func(s *TrayService) {
		s.icon = []byte{}
		s.icon = append(s.icon, icon...)
	}
}

type MenuItem struct {
	title   string
	tooltip string
	attach  any
	icon    []byte
	onClick func(menu *MenuItem)
}

func NewMenuItem() *MenuItem {
	return &MenuItem{
		icon: make([]byte, 0),
	}
}
func (s *MenuItem) Title(title string) *MenuItem {
	s.title = title
	return s
}

func (s *MenuItem) Tooltip(tooltip string) *MenuItem {
	s.tooltip = tooltip
	return s
}

func (s *MenuItem) Attach(data any) *MenuItem {
	s.attach = data
	return s
}
func (s *MenuItem) Icon(icon []byte) *MenuItem {
	s.icon = append(s.icon, icon...)
	return s
}
func (s *MenuItem) OnClick(fun func(menu *MenuItem)) *MenuItem {
	s.onClick = fun
	return s
}

func NewTrayService(options ...Option) *TrayService {
	result := &TrayService{
		onReady: func() {},
		icon:    appIcon,
		tooltip: "",
		title:   "",
		menus:   make([]*MenuItem, 0),
		onExit:  func() {},
		logger:  logger.DefaultLogger,
	}
	//systray.SetIcon(appIcon)
	for _, v := range options {
		v(result)
	}
	return result
}

func (s *TrayService) AddMenuItem(item ...*MenuItem) (r *TrayService) {
	s.menus = append(s.menus, item...)
	return s
}
func (s *TrayService) Quit() {
	systray.Quit()
}
func (s *TrayService) Logger(logger logger.ILogger) (r *TrayService) {
	s.logger = logger
	return s
}
func (s *TrayService) Tooltip(tool string) (r *TrayService) {
	s.tooltip = tool
	return s
}
func (s *TrayService) Title(title string) (r *TrayService) {
	s.title = title
	return s
}

// 托盘程序
func (s *TrayService) OnReady(on func()) *TrayService {
	s.onReady = on
	return s
}
func (s *TrayService) OnExit(on func()) *TrayService {
	s.onExit = on
	return s
}
func (s *TrayService) Icon(appIcon []byte) *TrayService {
	s.icon = appIcon
	return s
}
func (s *TrayService) run() {
	systray.Run(func() {
		s.onReady()
		if s.tooltip != "" {
			systray.SetTooltip(s.tooltip)
		}
		if s.title != "" {
			systray.SetTitle(s.title)
		}
		if len(s.icon) > 0 {
			systray.SetIcon(s.icon)
		}
		for i := range s.menus {
			item := s.menus[i]
			menu := systray.AddMenuItem(item.title, item.tooltip)
			if len(item.icon) > 0 {
				menu.SetIcon(item.icon)
			}
			go func(menu *systray.MenuItem) {
				for range menu.ClickedCh {
					s.logger.Debug("click", menu.String())
					item.onClick(item)
				}
			}(menu)
		}
	}, s.onExit)
}
func (s *TrayService) Start() {
	s.run()
}
