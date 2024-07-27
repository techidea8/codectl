package token

import (
	"errors"
	"net/http"

	"github.com/techidea8/restgo/utils"
)

func GenerateToken(values map[string]interface{}) (string, error) {
	return DefaultTokenManager.GenerateToken(values)
}
func ParseToken(in interface{}) (result map[string]interface{}, err error) {
	if token, ok := in.(string); ok {
		return DefaultTokenManager.ParseToken(token)
	} else if req, ok := in.(*http.Request); ok {
		token := utils.GetAuthorizationFromRequest(req)
		return DefaultTokenManager.ParseToken(token)
	} else {
		return nil, errors.New("不支持的数据类型")
	}
}
