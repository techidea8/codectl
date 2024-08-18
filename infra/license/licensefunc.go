package license

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/techidea8/codectl/infra/cryptor"
	"github.com/techidea8/codectl/infra/slicekit"
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
		file: "LICENSE",
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

	if !slicekit.Contains([]int{16, 24, 32}, len(ctrl.key)) {
		return errors.New("KEY长度不对")
	}

	data, err := ctrl.data.Bytes()
	if err != nil {
		return err
	}
	data, err = cryptor.AesEcbEncrypt(data, []byte(ctrl.key))
	if err != nil {
		return err
	}
	b64data := cryptor.Base64StdEncode(data)
	return os.WriteFile(ctrl.file, []byte(b64data), 0664)
}

// licence 解析
func (ctrl *LicenseCtrl) Parse() (licenseData *LicenseData, err error) {
	// base64 解码
	filedata, err := os.ReadFile(ctrl.file)
	if err != nil {
		return
	}
	data, err := cryptor.Base64StdDecode(string(filedata))
	if err != nil {
		return
	}
	//进行解密
	result, err := cryptor.AesEcbDecrypt(data, []byte(ctrl.key))
	if err != nil {
		return
	}
	licenseData = NewLicenseData()
	err = json.Unmarshal(result, licenseData)
	return licenseData, err
}
