package license

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/techidea8/codectl/infra/hardware"
	"github.com/techidea8/codectl/infra/stringx"
)

// licence 数据结构
type LicenseData struct {
	Biz   string    `json:"biz"` //业务编码
	Lid   string    `json:"lid"`
	Cid   string    `json:"cid"`
	Dvd   string    `json:"dvd"`
	Eat   time.Time `json:"eat"`
	Iat   time.Time `json:"iat"`
	Iss   string    `json:"iss"` //发行方
	Claim map[string]string
}

// 初始化数据
func NewLicenseData() *LicenseData {
	liceId := stringx.PKID()
	machine := hardware.Machine()
	return &LicenseData{
		Lid:   liceId,
		Dvd:   machine.MachineId,
		Claim: map[string]string{},
	}
}

// 添加其他数据
func (s *LicenseData) Add(key, value string) *LicenseData {
	s.Claim[key] = value
	return s
}

// 设置企业ID
func (s *LicenseData) CorpId(corpId string) *LicenseData {
	s.Cid = corpId
	return s
}

// 设置过期时间
func (s *LicenseData) BizCode(bizCode string) *LicenseData {
	s.Biz = bizCode
	return s
}

const TIMEFORMATE = "2006-01-02 15:04:05"

// 设置过期时间
func (s *LicenseData) ExpireAfter(tm any) *LicenseData {
	switch tm := tm.(type) {
	case time.Duration:
		s.Eat = time.Now().Add(tm)
	case time.Time:
		s.Eat = tm
	case string:
		s.Eat, _ = time.Parse(tm, TIMEFORMATE)
	case int64:
		s.Eat = time.Unix(int64(tm), 0)
	case uint64:
		s.Eat = time.Unix(int64(tm), 0)
	}
	s.Iat = time.Now()
	return s
}

// 全部bytes
func (s *LicenseData) Bytes() (bts []byte, err error) {
	err = s.Validate()
	if err != nil {
		return
	}
	return json.Marshal(s)
}

// 全部bytes
func (s *LicenseData) Validate() error {
	if s.Eat.IsZero() {
		return errors.New("缺少过期时间")
	}
	if s.Dvd == "" {
		return errors.New("缺少设备编码")
	}
	machine := hardware.Machine()
	if s.Dvd != machine.MachineId {
		return errors.New("该设备和license不匹配")
	}
	return nil

}

// 绑定设备ID
func (s *LicenseData) DeviceId(deviceId string) *LicenseData {
	s.Dvd = deviceId
	return s
}

// 绑定设备ID
func (s *LicenseData) Write(w io.Writer) error {
	bts, err := s.Bytes()
	if err != nil {
		return err
	}
	_, err = w.Write(bts)
	return err
}
