# 上传到本地

```
s := goo.NewServer()
s.GET("/upload/local", func(c *gin.Context) {
    filename, _ := goo_upload.Local("/upload").Upload(c)
    c.JSON(200, gin.H{"filename": filename})
})
s.Run(":18080")
```

# 上传到OSS

```
var (
	conf = goo_upload.OSSConfig{
		AccessKeyId:     "",
		AccessKeySecret: "",
		Endpoint:        "",
		Bucket:          "",
		Domain:          "",
	}
)

func main() {
	body, _ := ioutil.ReadFile("1.txt")
	url, _ := goo_upload.OSS.Upload("1.txt", body)
	fmt.Println(url)
}
```