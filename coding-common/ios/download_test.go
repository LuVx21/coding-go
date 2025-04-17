package ios

import (
	"os"
	"testing"
)

const (
	_url = ""
)

func Test_download_00(t *testing.T) {
	home := os.TempDir()
	DownloadFile(_url, home+"/1.jpg")
	DownloadWithProgress(_url, home+"/2.jpg")
	DownloadWithResume(_url, home+"/3.jpg")
	DownloadConcurrent(_url, home+"/4.jpg", 3)
}
