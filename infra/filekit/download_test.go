package filekit

import "testing"

func TestDownload(t *testing.T) {
	DownloadAs("appdata", "http://a.b.c/a.png?123=5")
}
