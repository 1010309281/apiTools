>> API文档

api返回数据结构: 
```text
{
"code": 0,    // 请求返回状态码(0 请求成功, 1 请求失败)
"data": {}    // 返回的实际数据
}
```

#### 1. whois域名查询

- 请求地址：`/api/whoisquery`

- 请求方式: `get/post`

- 请求数据: 
```text
{
"domain": "baidu.com",  // whois查询的域名 -> 必传字段
"type": "json",         // 返回的数据类型(text 文本串, json json格式数据), 默认json  -> 可选字段    
}
```
- 返回数据:
```text
{
// 0 获取到域名whois信息
// 1 域名解析失败
// 2 域名未注册
// 3 暂不支持此域名后缀查询
// 4 域名查询失败
// 5 请求数据错误
"status": 0,      // 域名查询状态
"data": "",       // whois数据
"msg": "",        // 消息
}
```

#### 1. 短链接转换

*1) 转换为短链接* 
 
- 请求地址：`/api/toshorturl`

- 请求方式: `get/post`

- 请求数据: 
```text
{
"url": "xxx",             // 长链接url -> 必传字段
// 配置domain需要将自己的域名配置CNAME记录解析到服务提供网站的域名
"domain": "json",         // 短链接域名绑定自己的域名 (默认为系统当前域名) -> 可选字段  
"expireTime": 30,         // 设置过期时间, (以分钟为单位, -1代表用不过期)  -> 可选字段  
}
```
- 返回数据:
```text
!!! 如果使用了自己的`domain`没有配置CNAME记录则短链接生成后无法解析
{
"code": 0             // 转换成功状态码(0 成功, 非零 失败) 
"domain": "xxxx",     // 短地址配置的域名
"shortUrl": "xxx",    // 短链接地址
"msg": "",            // 消息
}
```

*2) 短链接解析回长链接*

- 请求地址：`/api/parseshorturl`

- 请求方式: `get/post`

- 请求数据: 
```text
{
"shortUrl": "xxxxx",    // 短链接地址  -> 必传字段
}
```
- 返回数据:
```text
!!! 如果使用了自己的`domain`没有配置CNAME记录则短链接生成后无法解析
{
"code": 0             // 转换成功状态码(0 成功, 非零 失败) 
"domain": "xxxx",     // 短地址配置的域名
"longUrl": "xxxx",    // 原长连接地址
"msg": "",            // 消息
}
```