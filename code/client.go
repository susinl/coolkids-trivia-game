package code

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/susinl/coolkids-trivia-game/util"
)

type CheckReCaptchaClientFn func(token string) (*VerifyUserGoogleApiResponse, error)

func NewCheckReCaptchaClientFn(cli *http.Client) CheckReCaptchaClientFn {
	return func(token string) (*VerifyUserGoogleApiResponse, error) {
		data := url.Values{}
		data.Set(util.ReCaptchaSecretTag, os.Getenv("RECAPTCHA_V3_SECRET"))
		data.Set(util.ReCaptchaTokenTag, token)

		httpReq, err := http.NewRequest(http.MethodPost, viper.GetString("client.google-api.recaptchav3"), strings.NewReader(data.Encode()))
		if err != nil {
			return nil, err
		}
		httpReq.Header.Set(util.ContentType, util.ApplicationUrlEncoded)

		httpResp, err := cli.Do(httpReq)
		if err != nil {
			return nil, err
		}
		defer httpResp.Body.Close()

		respBody, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, err
		}

		var verifyUserGoogleApiResponse VerifyUserGoogleApiResponse
		if err := json.Unmarshal(respBody, &verifyUserGoogleApiResponse); err != nil {
			return nil, err
		}
		return &verifyUserGoogleApiResponse, nil
	}
}
