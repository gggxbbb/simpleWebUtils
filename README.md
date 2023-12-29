# simpleWebUtils

一系列与 Web/HTTP 相关的简易工具

[Github](https://github.com/gggxbbb/simpleWebUtils)

## code

-> [/code/404](https://t.gxb.icu/code/404)

根据提交的状态值返回状态值，例如：
```
/code/404
/code/502
```

### detail

-> [/code/detail/404](https://t.gxb.icu/code/detail/404)

根据提交的状态值返回状态值的详细信息，例如：
```
/code/detail/404
/code/detail/502
```

## ua

返回 User-Agent

-> [/ua](https://t.gxb.icu/ua)

返回 User-Agent

### ua/analyze

-> [ua/analyze](https://t.gxb.icu/ua/analyze)

分析 User-Agent

## ip

-> [/ip](https://t.gxb.icu/ip)

返回访问者 IP

## minecraft

### bedrock
获取 Minecraft Bedrock Motd 信息

```
GET /minecraft/bedrock/<address>
GET /minecraft/bedrock/<address>/<port>
```

或通过 POST 提交

```
POST /minecraft/bedrock
{
    "server": "<address>",
    "port": <port>
}
```
