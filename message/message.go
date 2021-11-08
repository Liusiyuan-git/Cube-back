package message

import (
	"Cube-back/log"
	"Cube-back/models/common/configure"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

var Message *M

var CodeBox = map[string]string{
	"register": "1177892",
	"login":    "1181535",
	"password": "1181541",
}

type Conf struct {
	SecretId  string
	SecretKey string
}

type M struct {
	client *sms.Client
}

func (m *M) SendMessage(mode, value, phone string) string {
	request := sms.NewSendSmsRequest()

	request.PhoneNumberSet = common.StringPtrs([]string{phone})
	request.SmsSdkAppId = common.StringPtr("1400590793")
	request.SignName = common.StringPtr("魔方技术")
	request.TemplateId = common.StringPtr(CodeBox[mode])
	request.TemplateParamSet = common.StringPtrs([]string{value, "2"})

	response, err := m.client.SendSms(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		log.Error(err)
	}
	if err != nil {
		log.Error(err)
	}
	return response.ToJsonString()
}

func init() {
	Message = new(M)
	conf := new(Conf)
	configure.Get(&conf)
	credential := common.NewCredential(
		conf.SecretId,
		conf.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	Message.client, _ = sms.NewClient(credential, "ap-guangzhou", cpf)
}
