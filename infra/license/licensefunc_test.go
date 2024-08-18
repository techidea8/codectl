package license

import (
	"fmt"
	"testing"
	"time"
)

func TestCreateLicence(t *testing.T) {
	var lice *LicenseData = NewLicenseData()
	var licectrl *LicenseCtrl = NewLicenseCtrl("codectl")
	fmt.Println("输入CorpId")
	fmt.Scanln(&lice.Cid)
	fmt.Println("输入ExpireAt")
	var expireAt string
	fmt.Scanln(&expireAt)
	lice.Eat, _ = time.Parse(expireAt, "2006-01-02")
	//lice.ExpireAt = time.
	licectrl.Release()
	fmt.Println("============")
	data, _ := licectrl.Parse()
	fmt.Printf("%+v", data)
}

func TestCreateLicence2(t *testing.T) {
	var lice *LicenseData = NewLicenseData()
	var licectrl *LicenseCtrl = NewLicenseCtrl("codectl")

	lice.Cid = "10086"
	lice.Eat, _ = time.Parse("2025-08-14", "2006-01-02")
	lice.Dvd = "02390KLSKDSALK201"
	lice.Biz = "standalone"
	lice.Biz = "saas"
	lice.ExpireAfter(time.Now().AddDate(1, 0, 0))
	licectrl.WithData(lice)
	licectrl.Release()
	fmt.Println("============")
	data, _ := licectrl.Parse()
	fmt.Printf("%+v", data)
}
