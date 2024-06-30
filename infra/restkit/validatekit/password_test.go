package validatekit

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

type AuthorizeReq struct {
	UserName    string `json:"userName" validate:"required,passwordscore=3,min=5" errmsg:"用户名格式不正确"`
	PassWord    string `json:"password"`
	CaptchaKey  string `json:"captchaKey"`
	CaptchaCode string `json:"captchaCode"`
}

func TestPasswordLevel(t *testing.T) {
	var validate *validator.Validate = validator.New()
	validate.RegisterValidation("passwordscore", PasswordScore)
	ptrreq := &AuthorizeReq{
		UserName: "test",
	}
	err := validate.Struct(ptrreq)

	err = ProcessError(ptrreq, err)
	if err != nil {
		t.Error(err)
	}

}
