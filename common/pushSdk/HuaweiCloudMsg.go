package pushSdk

type HWCloudMsgConfig struct {
	AppInfoKey    string
	AppInfoSecret string
	ApiAddress    string
	Sender        string
	TemplateId    string
	Signature     string
	Receiver      string
}

func NewTestConfig() HWCloudMsgConfig {
	return HWCloudMsgConfig{
		//自己写 保密
	}
}
