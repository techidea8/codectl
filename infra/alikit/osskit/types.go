package osskit

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type PolicyToken struct {
	AccessKeyId string                 `json:"accessid"`
	Host        string                 `json:"host"`
	Expire      int64                  `json:"expire"`
	Signature   string                 `json:"signature"`
	Policy      string                 `json:"policy"`
	Directory   string                 `json:"dir"`
	Callback    string                 `json:"callback"`
	Xvar        map[string]interface{} `json:"xvar"`
}

type CallbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

// // 请填写您的AccessKeyId。
// var accessKeyId string = "<yourAccessKeyId>"

// // 请填写您的AccessKeySecret。
// var accessKeySecret string = "<yourAccessKeySecret>"

// // host的格式为bucketname.endpoint，请替换为您的真实信息。
// var host string = "https://bucket-name.oss-cn-hangzhou.aliyuncs.com'"

// // callbackUrl为上传回调服务器的URL，请将下面的IP和Port配置为您自己的真实信息。
// var callbackUrl string = "http://192.0.2.0:8888";

// // 上传文件时指定的前缀。
// var upload_dir string = "user-dir-prefix/"

// // 上传策略Policy的失效时间，单位为秒。
// var expire_time int64 = 30

type OssConf struct {
	AccessKeyId     string
	AccessKeySecret string
	Host            string
	CallbackUrl     string
	UploadDir       string
	ExpireTime      int64
}
