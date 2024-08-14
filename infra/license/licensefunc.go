package license

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/techidea8/codectl/infra/cryptor"
)

// 一个随机的key,不能随机
var _key = cryptor.Md5String("winlion@techidea8.com")

type LicenseCtrl struct {
	file string
	key  string
	data *LicenseData
}

func NewLicenseCtrl(appName string) *LicenseCtrl {
	return &LicenseCtrl{
		file: "app.license",
		key:  cryptor.Md5String(_key + appName),
	}
}

// 设置输出文件
func (ctrl *LicenseCtrl) SetFileName(file string) *LicenseCtrl {
	ctrl.file = file
	return ctrl
}

// 设置输出文件
func (ctrl *LicenseCtrl) SetKey(key string) *LicenseCtrl {
	ctrl.key = key
	return ctrl
}

// 设置输出文件
func (ctrl *LicenseCtrl) WithData(data *LicenseData) *LicenseCtrl {

	ctrl.data = data
	return ctrl
}

// 设置输出文件
func (ctrl *LicenseCtrl) Release() error {
	if ctrl.data == nil {
		return errors.New("请配置参数")
	}

	data, err := ctrl.data.Bytes()
	if err != nil {
		return err
	}
	data = cryptor.AesEcbEncrypt(data, []byte(ctrl.key))
	return os.WriteFile(ctrl.file, data, 0664)
}

// licence 解析
func (ctrl *LicenseCtrl) Parse() (licenseData *LicenseData, err error) {
	filedata, err := os.ReadFile(ctrl.file)
	if err != nil {
		return
	}
	//进行解密
	result, err := cryptor.AesEcbDecrypt(filedata, []byte(ctrl.key))
	if err != nil {
		return
	}
	licenseData = NewLicenseData()
	err = json.Unmarshal(result, licenseData)
	return licenseData, err
}
