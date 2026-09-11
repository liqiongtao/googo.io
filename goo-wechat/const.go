package goowechat

const (
	cgi_token_key = "wx:cgi:token:%s"
	cgi_token_url = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"

	cgi_ticket_key = "wx:cgi:ticket:%s"
	cgi_ticket_url = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"

	jsapi_ticket_qs = "jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s"

	menu_create_url = "https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s"
	menu_get_url    = "https://api.weixin.qq.com/cgi-bin/menu/get?access_token=%s"
	menu_del_url    = "https://api.weixin.qq.com/cgi-bin/menu/delete?access_token=%s"

	oauth2_authorize_url       = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=%s#wechat_redirect"
	sns_oauth2_accessToken_url = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	sns_userinfo_url           = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s"
	sns_jsscode2sess_url       = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

	message_tpl_send_url = "https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%s"

	getwxacodeunlimit_url = "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s"
)
