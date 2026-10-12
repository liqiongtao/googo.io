package goomessage

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"runtime"
	"sync"

	goo_request "github.com/liqiongtao/googo.io/goo-request"
	gooutils "github.com/liqiongtao/googo.io/goo-util"
)

var (
	__qiWeiCH   chan struct{}
	__qiWeiOnce sync.Once
)

// QiWeiArticle 企微图文消息条目
type QiWeiArticle struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
	PicURL      string `json:"picurl,omitempty"`
}

// QiWei 发送文本消息
// 文档: https://developer.work.weixin.qq.com/document/path/99110
func QiWei(hookUrl, text string) error {
	return QiWeiText(hookUrl, text, nil, nil)
}

// QiWeiText 发送文本消息，可 @指定成员（userid / 手机号，传 @all 表示全员）
func QiWeiText(hookUrl, text string, mentionedUserIds, mentionedMobiles []string) error {
	if text == "" {
		return nil
	}
	body := map[string]any{
		"content": text,
	}
	if len(mentionedUserIds) > 0 {
		body["mentioned_list"] = mentionedUserIds
	}
	if len(mentionedMobiles) > 0 {
		body["mentioned_mobile_list"] = mentionedMobiles
	}
	return qiWeiSend(hookUrl, map[string]any{
		"msgtype": "text",
		"text":    body,
	})
}

// QiWeiMarkdown 发送 markdown 消息（支持 <@userid>、字体颜色等）
func QiWeiMarkdown(hookUrl, content string) error {
	if content == "" {
		return nil
	}
	return qiWeiSend(hookUrl, map[string]any{
		"msgtype": "markdown",
		"markdown": map[string]any{
			"content": content,
		},
	})
}

// QiWeiMarkdownV2 发送 markdown_v2 消息（不支持 @成员、字体颜色）
func QiWeiMarkdownV2(hookUrl, content string) error {
	if content == "" {
		return nil
	}
	return qiWeiSend(hookUrl, map[string]any{
		"msgtype": "markdown_v2",
		"markdown_v2": map[string]any{
			"content": content,
		},
	})
}

// QiWeiImage 发送图片消息（JPG/PNG，原始数据不超过 2M）
func QiWeiImage(hookUrl string, image []byte) error {
	if len(image) == 0 {
		return nil
	}
	sum := md5.Sum(image)
	return qiWeiSend(hookUrl, map[string]any{
		"msgtype": "image",
		"image": map[string]any{
			"base64": base64.StdEncoding.EncodeToString(image),
			"md5":    hex.EncodeToString(sum[:]),
		},
	})
}

// QiWeiNews 发送图文消息（1~8 条）
func QiWeiNews(hookUrl string, articles ...QiWeiArticle) error {
	if len(articles) == 0 {
		return nil
	}
	if len(articles) > 8 {
		return errors.New("qiwei news articles max 8")
	}
	return qiWeiSend(hookUrl, map[string]any{
		"msgtype": "news",
		"news": map[string]any{
			"articles": articles,
		},
	})
}

// QiWeiFile 上传并发送文件附件（不超过 20M）
func QiWeiFile(hookUrl, fileName string, r io.Reader) error {
	mediaId, err := qiWeiUpload(hookUrl, "file", fileName, r)
	if err != nil {
		return err
	}
	return qiWeiSend(hookUrl, map[string]any{
		"msgtype": "file",
		"file": map[string]any{
			"media_id": mediaId,
		},
	})
}

// QiWeiVoice 上传并发送语音（AMR，不超过 2M / 60s）
func QiWeiVoice(hookUrl, fileName string, r io.Reader) error {
	mediaId, err := qiWeiUpload(hookUrl, "voice", fileName, r)
	if err != nil {
		return err
	}
	return qiWeiSend(hookUrl, map[string]any{
		"msgtype": "voice",
		"voice": map[string]any{
			"media_id": mediaId,
		},
	})
}

func qiWeiSend(hookUrl string, data map[string]any) error {
	qiWeiAcquire()
	defer qiWeiRelease()

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("[goo-msg][2001]", err)
		return err
	}

	buf, err := goo_request.PostJson(hookUrl, b)
	if err != nil {
		fmt.Println("[goo-msg][2002]", err)
		return err
	}
	return qiWeiCheck(buf)
}

func qiWeiUpload(hookUrl, typ, fileName string, r io.Reader) (string, error) {
	if r == nil {
		return "", errors.New("qiwei upload reader is nil")
	}
	if fileName == "" {
		return "", errors.New("qiwei upload filename is empty")
	}

	key, err := qiWeiKey(hookUrl)
	if err != nil {
		return "", err
	}

	qiWeiAcquire()
	defer qiWeiRelease()

	uploadUrl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media?key=%s&type=%s",
		url.QueryEscape(key), url.QueryEscape(typ))
	buf, err := goo_request.Upload(uploadUrl, "media", fileName, r, nil)
	if err != nil {
		fmt.Println("[goo-msg][2012]", fileName, err)
		return "", err
	}
	if err = qiWeiCheck(buf); err != nil {
		return "", err
	}

	rst, err := gooutils.Byte(buf).Params()
	if err != nil {
		fmt.Println("[goo-msg][2013]", fileName, string(buf), err)
		return "", err
	}
	mediaId := rst.Get("media_id").String()
	if mediaId == "" {
		return "", errors.New("qiwei upload empty media_id")
	}
	return mediaId, nil
}

func qiWeiKey(hookUrl string) (string, error) {
	u, err := url.Parse(hookUrl)
	if err != nil {
		return "", err
	}
	key := u.Query().Get("key")
	if key == "" {
		return "", errors.New("qiwei webhook missing key")
	}
	return key, nil
}

func qiWeiCheck(buf []byte) error {
	if len(buf) == 0 {
		return errors.New("empty qiwei response")
	}

	rst, err := gooutils.Byte(buf).Params()
	if err != nil {
		fmt.Println("[goo-msg][2003]", string(buf), err)
		return err
	}

	errcode := rst.Get("errcode").Int()
	if errcode == 0 {
		return nil
	}

	errmsg := rst.Get("errmsg").String()
	fmt.Println("[goo-msg][2004]", string(buf))
	if errmsg == "" {
		return fmt.Errorf("qiwei error: %d", errcode)
	}
	return errors.New(errmsg)
}

func qiWeiAcquire() {
	__qiWeiOnce.Do(func() {
		__qiWeiCH = make(chan struct{}, runtime.NumCPU())
	})
	__qiWeiCH <- struct{}{}
}

func qiWeiRelease() {
	<-__qiWeiCH
}
