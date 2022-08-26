#

## SS58AuthData

format:

发起请求时，在header中增加`SS58AuthData`字段，格式如下：
```http request
SS58地址:原始数据:签名后的数据的十六进制表示
```

example:

```http request
SS58AuthData: 5HpGQhD72vZGgAFMMiCDY61mHYtANs6B4kZXrpptGm276KnT:hello:b2146a773345dce02a4c7c7416a9b215d19157842f427f6ad991e3f40e24271add31cecc28c6d8a610a0c1cb74e24b6218c7139345ee57b5b7fbd1ba96fb6688
```