package swagger

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

type Config struct {
	// "api前缀，例如/api/v1"，默认为空
	BasePath string

	// swagger文档标题
	Title string

	// swagger文档描述
	Description string

	// 文档版本
	DocVersion string

	// swagger ui的地址
	SwaggerUiUrl string

	// 文档Url地址，例如开发服务器http://baidu.com
	// 如果本值是doc，那么http://baidu.com/doc就可以重定向到@SwaggerUiUrl
	SwaggerUrlPrefix string

	// swagger文档的地址，用于调试，release直接打包到二进制里面。默认为空
	DocFilePath string

	// 用于支持swagger ui认证头的参数
	Headers []SecurityDefinition

	// 是否调试模式
	Debug bool
}

func (this *Config) initDefault() {
	if len(this.Title) == 0 {
		this.Title = "Swagger文档"
	}
	if len(this.Description) == 0 {
		this.Title = "Swagger文档描述"
	}
	if len(this.DocVersion) == 0 {
		this.Title = "0.0.1"
	}
	if len(this.SwaggerUiUrl) == 0 {
		// http://swagger.qiujinwu.com
		this.SwaggerUiUrl = "http://petstore.swagger.io/"
	}
	if len(this.SwaggerUrlPrefix) == 0 {
		this.SwaggerUrlPrefix = "apidoc"
	}
}

func InitializeApiRoutes(grouter *gin.Engine, config *Config, docLoader DocLoader) {
	if gOption != nil {
		panic("swagger inited yet")
		return
	}

	if config == nil || docLoader == nil {
		panic("invalid swagger parameter")
	}
	config.initDefault()
	gOption = newOptions(config)
	gOption.docLoader = docLoader

	grouter.GET("/"+config.SwaggerUrlPrefix+"/spec", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		swaggerData1 := gOption.swaggerData

		headersDef := make(map[string]SecurityDefinition)
		if len(config.Headers) > 0 {
			for _, v := range config.Headers {
				key := v.Type
				v.In = "header"
				v.Type = "apiKey"
				headersDef[key] = v
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"basePath": config.BasePath,
			"swagger":  swaggerVersion,
			"info": struct {
				Description string `json:"description"`
				Title       string `json:"title"`
				Version     string `json:"version"`
			}{
				Description: config.Description,
				Title:       config.Title,
				Version:     config.DocVersion,
			},
			"definition":          struct{}{},
			"paths":               swaggerData1,
			"securityDefinitions": headersDef,
		})

	})

	grouter.GET("/"+config.SwaggerUrlPrefix, func(c *gin.Context) {
		scheme := "http://"
		if c.Request.TLS != nil {
			scheme = "https://"
		}
		host := scheme + c.Request.Host + "/" + config.SwaggerUrlPrefix + "/spec"
		host = config.SwaggerUiUrl + "?url=" + url.QueryEscape(host)
		c.Redirect(http.StatusFound, host)
	})
}
