package goo_cos

var (
	// 文档: https://cloud.tencent.com/document/product/436/31923#.E6.9F.A5.E8.AF.A2.E5.AF.B9.E8.B1.A1.E5.88.97.E8.A1.A8
	DefaultAction = []string{
		"name/cos:PutObject",
		"name/cos:PutObjectCopy",
		"name/cos:PostObject",
		"name/cos:HeadObject",
		"name/cos:GetObject",
		"name/cos:DeleteObject",
	}

	// fmt.sprintf(ResourceTemplate, region, appid, bucket)
	ResourceTemplate = "qcs::cos:%s:uid/%s:%s/*"

	// 临时密钥有效时长，单位是秒，默认 8 小时
	DurationSeconds int64 = 8 * 3600
)
