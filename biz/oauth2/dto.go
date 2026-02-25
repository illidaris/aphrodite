package oauth2

import (
	"fmt"
	"time"
)

type AuthorizeParam struct {
	Verifier string `json:"verifier" form:"verifier" url:"verifier"`
	BizId    int64  `json:"bizId" form:"bizId" url:"bizId"`
	Expire   int64  `json:"expire" form:"expire" url:"expire"`
}

func (s AuthorizeParam) Valid(bizId int64) error {
	if s.Verifier == "" {
		return fmt.Errorf("[oauth2]code verifier is empty")
	}
	if s.BizId == 0 || s.BizId != bizId {
		return fmt.Errorf("[oauth2]biz id is empty or error")
	}
	if s.Expire < time.Now().Unix() {
		return fmt.Errorf("[oauth2]code verifier expired")
	}
	return nil
}
func (s AuthorizeParam) Encode(secret string) string {
	encoded, _ := AESEncode(s, secret)
	return encoded
}

func (s *AuthorizeParam) Decode(str, secret string) error {
	err := AESDecode(s, str, secret)
	if err != nil {
		return err
	}
	return nil
}

type OAuthCallbackParam struct {
	Code  string `json:"code" form:"code" url:"code"`
	State string `json:"state" form:"state" url:"state"`
}
