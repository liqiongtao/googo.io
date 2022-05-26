# `KEY` 命名规范

```
/命名空间/项目名/服务名/节点名
```

# 方法说明

- `Set` 设置 key-value
- `SetWithPrevKV` 设置 key-value 并且返回修改之前的 key-value
- `SetTTL` 设置有有效期的 key-value
- `SetTTLWithPrevKV` 设置有有效期的 key-value 并且返回修改之前的 key-value
- `Get` 返回 response 信息
- `GetString` 根据 key 前缀，返回 string-value
- `GetArray` 根据 key 前缀，返回 array-value
- `GetMap` 根据 key 前缀，返回 map-value
- `Del` 删除 key 并且返回删除之前的 key-value
- `DelWithPrefix` 根据 key 前缀删除， 并且返回删除之前的 key-value
- `RegisterService` 注册一个服务，并保持活跃
- `Watch` 根据 key 前缀观察，并返回 array-value
