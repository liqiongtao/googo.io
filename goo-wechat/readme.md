# 微信

微信SDK，实现微信菜单管理、H5授权、JSAPI、用户信息、小程序登录解析、小程序模板消息等

## 全局AccessToken

accessToken 缓存到redis, key: ```wx:cgi:token:%s```, 有效期7200秒

```
accessToken, _ := gooWeixin.GetCGIAccessToken(appid, secret)
```

设置自动刷新

```
gooWeixin.AutoRefreshCGIAccessToken(appid, secret)
```

## 全局Ticket

```
ticket, _ := gooWeixin.GetCGITicket(appid, secret)
```

## 设置菜单

```
func MenuCreate(appid, secret, content string) error
```

## 获取菜单

```
func MenuGet(appid, secret string) (string, error)
```

## 删除菜单

```
func MenuDelete(appid, secret string) error
```

## JSAPI

```
func JsApi(appid, secret, urlStr string) goo.Params
```

## 获取H5授权链接

```
func Oauth2AuthorizeUrl(appid, authorizeUrl, originUrl string) string
```

- authorizeUrl: 授权地址
- originUrl: 授权成功后的回跳地址

## 获取AccessToken

```
func Oauth2AccessToken(appid, secret, code string) (*Oauth2AccessTokenResponse, error)
```

## 获取用户信息 - sns

```
func SnsUserInfo(accessToken, openid string) (*SnsUserInfoResponse, error)
```

## 小程序登录

```
func JsCode2Session(appid, secret, code string) (*JsCode2SessionResponse, error)
```

## 小程序用户信息

```
func MinipUserInfo(sessionKey, encryptedData, iv string) (*MinipUserInfoResponse, error)
```

## 小程序发送模板消息

```
func SendTemplateMessage(appid, secret, openid, templateId, page, formId string, data interface{}) error
```
