package pushSdk

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	core "huaweicloud.com/apig/signer"

	"github.com/opensourceways/message-push/models/bo"
	"github.com/opensourceways/message-push/models/dto"
)

type MsgConfig struct {
	AppInfoKey    string `json:"app_info_Key"`
	AppInfoSecret string `json:"app_info_secret"`
	ApiAddress    string `json:"api_address"`
	Sender        string `json:"sender"`
	TemplateId    string `json:"template_id"`
	Signature     string `json:"signature"`
}

func SendHWCloudMessage(msgConfig MsgConfig, templateParas []string, recipient bo.RecipientPushConfig) dto.PushResult {
	//必填,请参考"开发准备"获取如下数据,替换为实际值
	appInfo := core.Signer{
		// 认证用的appKey和appSecret硬编码到代码中或者明文存储都有很大的安全风险，建议在配置文件或者环境变量中密文存放，使用时解密，确保安全；
		Key:    msgConfig.AppInfoKey,    //App Key
		Secret: msgConfig.AppInfoSecret, //App Secret
	}
	apiAddress := msgConfig.ApiAddress //APP接入地址(在控制台"应用管理"页面获取)+接口访问URI
	sender := msgConfig.Sender         //国内短信签名通道号
	templateId := msgConfig.TemplateId //模板ID

	//条件必填,国内短信关注,当templateId指定的模板类型为通用模板时生效且必填,必须是已审核通过的,与模板类型一致的签名名称

	signature := msgConfig.Signature //签名名称

	//必填,全局号码格式(包含国家码),示例:+86151****6789,多个号码之间用英文逗号分隔
	receiver := recipient.Message //短信接收人号码

	//选填,短信状态报告接收地址,推荐使用域名,为空或者不填表示不接收状态报告
	statusCallBack := ""

	/*
	 * 选填,使用无变量模板时请赋空值 string templateParas = "";
	 * 单变量模板示例:模板内容为"您的验证码是${1}"时,templateParas可填写为"[\"369751\"]"
	 * 双变量模板示例:模板内容为"您有${1}件快递请到${2}领取"时,templateParas可填写为"[\"3\",\"人民公园正门\"]"
	 * 模板中的每个变量都必须赋值，且取值不能为空
	 * 查看更多模板规范和变量规范:产品介绍>短信模板须知和短信变量须知
	 */

	templateParasString := "[\"" + strings.Join(templateParas[:], "\",\"") + "\"]"

	body := buildRequestBody(sender, receiver, templateId, templateParasString, statusCallBack, signature)
	_, err := post(apiAddress, []byte(body), appInfo)
	if err != nil {
		return dto.PushResult{Res: dto.Failed, Remark: err.Error(), RecipientId: recipient.RecipientId,
			PushAddress: recipient.Phone}
	}
	return dto.PushResult{
		Res:         dto.Succeed,
		Time:        time.Now(),
		Remark:      "succeed",
		RecipientId: recipient.RecipientId,
		PushAddress: recipient.Phone,
		PushType:    "message",
	}
}

/**
 * sender,receiver,templateId不能为空
 */
func buildRequestBody(sender, receiver, templateId, templateParas, statusCallBack, signature string) string {
	param := "from=" + url.QueryEscape(sender) + "&to=" + url.QueryEscape(receiver) + "&templateId=" + url.QueryEscape(templateId)
	if templateParas != "" {
		param += "&templateParas=" + url.QueryEscape(templateParas)
	}
	if statusCallBack != "" {
		param += "&statusCallback=" + url.QueryEscape(statusCallBack)
	}
	if signature != "" {
		param += "&signature=" + url.QueryEscape(signature)
	}
	return param
}

func post(url string, param []byte, appInfo core.Signer) (string, error) {
	if len(param) == 0 || appInfo == (core.Signer{}) {
		return "", nil
	}

	// 代码样例为了简便，设置了不进行证书校验，请在商用环境自行开启证书校验。
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(param))
	if err != nil {
		return "", err
	}

	// 对请求增加内容格式，固定头域
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// 对请求进行HMAC算法签名，并将签名结果设置到Authorization头域。
	appInfo.Sign(req)

	fmt.Println(req.Header)
	// 发送短信请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if resp.StatusCode != 200 {
		logrus.Error("发送短信失败", param)
		return "", errors.New(resp.Status)
	}
	fmt.Println(resp)

	// 获取短信响应
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
