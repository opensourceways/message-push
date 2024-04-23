package pushSdk

type HWCloudMsgConfig struct {
	AppInfoKey    string
	AppInfoSecret string
	ApiAddress    string
	Sender        string
	TemplateId    string
	Signature     string
}

func NewTestConfig() HWCloudMsgConfig {
	return HWCloudMsgConfig{
		AppInfoKey:    "C5I7EH4iPr4xPXRNa65Qg70bATxP",
		AppInfoSecret: "slpTMyXRJdAOldFm2bw6fxdtpoQy",
		ApiAddress:    "https://smsapi.cn-south-1.myhuaweicloud.com:443/sms/batchSendSms/v1",
		Sender:        "8824041123653",
		TemplateId:    "33fabc75615647a886a87daa7064f766",
		Signature:     "开源社区",
	}
}
