package pushSdk

type HWCloudMsgConfig struct {
	AppInfoKey    string `json:"app_info_Key"`
	AppInfoSecret string `json:"app_info_secret"`
	ApiAddress    string `json:"api_address"`
	Sender        string `json:"sender"`
	TemplateId    string `json:"template_id"`
	Signature     string `json:"signature"`
}
