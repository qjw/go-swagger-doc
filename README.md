## go-swagger-doc是一个用于生成swagger-ui依赖的json文件的工具库

默认包含了[gin](https://github.com/gin-gonic/gin)的handle，亦可单独适配其他框架

输出格式类似于<http://petstore.swagger.io/v2/swagger.json>

支持
1. yml文件
2. 对象自动推导

# YML文件定义

yml文件是用yml格式定义一个接口描述的文本文件，虽然繁琐但是某些情况下是必需，比如输入是文件上传，或者输出是图片的情形，例如
``` yaml
upload_tmp_media:
  summary: 上传临时素材
  tags:
    - 素材管理
  parameters:
    - name: file
      in: formData
      description: 待上传的文件
      required: true
      type: file
    - name: type
      in: formData
      description: 类型,image/voice/video/file
      type: string
      required: true
  consumes:
    - multipart/form-data
  responses:
    200:
      schema:
        properties:
          result:
            description: 错误码，默认0，成功
            type: integer
          message:
            description: 错误内容，默认OK
            type: string
```

或者
``` yaml
get_tmp_media:
  summary: 获取临时素材
  tags:
    - 素材管理
  parameters:
    - in: query
      name: media_id
      type: string
  produces:
    - image/png
    - image/jpg
    - image/jpeg
    - image/gif
  responses:
    200:
      type: file
```

# struct自动推导
在结构中增加json的tag，例如
``` go
type SendTextMsgParam struct {
	Sender   corp.KfMsgUserObj `json:"sender"`
	Receiver corp.KfMsgUserObj `json:"receiver"`
	Content  string            `json:"content"`
}
```

然后调用下面代码即可
``` go
swagger.Swagger2(urouter,"send_text_msg","post",&swagger.StructParam{
    JsonData:&SendTextMsgParam{},
    ResponseData:&swagger.SuccessResp{},
    Tags:[]string{"企业号-客服"},
    Summary:"发送文本消息",
})
router.POST("send_text_msg",func(c *gin.Context) {

}
```

对于yaml文件的定义，通过下面的代码关联
``` go
swagger.Swagger(urouter,"upload_tmp_media","post","corp.yaml:upload_tmp_media")
router.POST("upload_tmp_media",func(c *gin.Context) {
}
```

## StructParam
StructParam结构定义了函数的参数、返回值等依赖定义，FormData/JsonData不能共存
``` go
type StructParam struct {
	FormData     interface{} // form参数
	JsonData     interface{} // json参数
	QueryData    interface{} // query参数
	PathData     interface{} // path参数
	ResponseData interface{} // 返回值
	Description  string
	Summary      string
	Tags         []string
}
```

# 授权
在初始化传入的Config中，传入SecurityDefinition即可
``` go
[]swagger.SecurityDefinition{
    swagger.SecurityDefinition{
        Type:"apiKey",
        Name:"X-APP-ID",
        Description:"app id,用于区分是那个应用申请授权",
    },
    swagger.SecurityDefinition{
        Type:"apiHash",
        Name:"X-APP-HASH",
        Description:"app hash,用自己的私钥对uri和参数签名的hash",
    },
},
```