package oskit

import (
	"fmt"
	"testing"

	"github.com/getlantern/systray"
)

func onReady() {
	systray.SetIcon(appIcon)
	systray.SetTitle("Awesome App")
	systray.SetTooltip("Pretty awesome超级棒")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// Sets the icon of a menu item. Only available on Mac and Windows.
	mQuit.SetIcon(appIcon)
}

func onExit() {
	// clean up here
}
func Test1(t *testing.T) {
	svc := NewTrayService().Icon(appIcon).OnExit(onExit)
	svc.OnReady(func() {
		svc.AddMenuItem(NewMenuItem().Title("QUIT").OnClick(func(menu *MenuItem) {
			fmt.Println("QUIR")
			systray.Quit()
		}))
	})
	svc.Start()
	fmt.Println("SS")
}
