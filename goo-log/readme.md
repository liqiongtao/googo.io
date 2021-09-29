# 日志对象

- 默认

```
init() 初始化了默认日志对象 __log
```

- 新建（基于控制台日志）

```
l := goo_log.New(goo_log.NewConsoleAdapter())
l.Debug("this is debug")
```

- 新建（基于文件日志）

```
l := goo_log.New(goo_log.NewFileAdapter(
    goo_log.FileTagOption("sql"),
    goo_log.FileDirnameOption("logs-logs/"),
    goo_log.FileMaxSizeOption(100*1024*1024),
))
l.Debug("this is debug")
```

# 设置适配器

## 文件适配器 (默认)

- 默认参数

```
goo_log.SetAdapter(goo_log.NewFileAdapter())
```

- 自定义参数

```
goo_log.SetAdapter(goo_log.NewFileAdapter(
    goo_log.FileTagOption("sql"),
    goo_log.FileDirnameOption("logs-logs/"),
    goo_log.FileMaxSizeOption(100*1024*1024),
))
```

## 设置文件适配器参数

```
goo_log.SetFileAdapterOptions(
    goo_log.FileTagOption("sql"),
    goo_log.FileDirnameOption("logs-logs/"),
    goo_log.FileMaxSizeOption(100*1024*1024),
)
```

## 命令行适配器

```
goo_log.SetAdapter(goo_log.NewConsoleAdapter())
```

# 设置过滤路径

```
goo_log.SetTrimPath("/a/b/c")
```

# 添加钩子函数

```
goo_log.AddHook(func(msg goo_log.Message) {
    fmt.Println(msg.Level, msg.String()
})
```

# 添加日志字段

```
goo_log.WithField("name", "hnatao").Debug()
```

# 日志级别输出

```
goo_log.Debug("this is debug")
goo_log.Info("this is info")
goo_log.Warn("this is warn")
goo_log.Error("this is error")
goo_log.Panic("this is panic")
goo_log.Fatal("this is fatal")
```