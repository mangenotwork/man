package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
	"html/template"
	"net/http"
)

func main() {
	router := gin.Default()

	router.GET("/profile", xssCase)

	router.Use(CSRFMiddle())

	router.GET("/form", func(c *gin.Context) {
		// 从请求中获取CSRF令牌（如果有的话）
		token := c.Request.FormValue("_csrf")
		c.HTML(http.StatusOK, "form.tmpl", gin.H{
			"csrfToken": token, // 将令牌传递给模板以便嵌入到表单中
		})
	})

	_ = router.Run(":18080")
}

/*
跨站脚本攻击（XSS）是一种攻击方式，攻击者通过在用户输入中注入恶意脚本，使之在用户浏览器中执行。
可以使用: html/template 包的 template.HTMLEscapeString 函数可以防止 XSS 攻击。
*/
func xssCase(c *gin.Context) {
	userInput := c.Query("input")

	// 防止 XSS 攻击
	safeHTML := template.HTMLEscapeString(userInput)

	c.HTML(http.StatusOK, "profile.tmpl", gin.H{
		"input": safeHTML,
	})
}

/*
跨站请求伪造（CSRF）是一种攻击方式，攻击者通过伪装成受信任用户的请求，以在用户不知情的情况下执行恶意操作。
在 Gin 框架中，可以使用 github.com/gorilla/csrf 中间件来防范 CSRF 攻击。
使用csrf中间件
*/
func CSRFMiddle() gin.HandlerFunc {
	csrfMiddleware := csrf.Protect([]byte("32-byte-long-auth-key"))
	// 这里使用adpater包将csrfMiddleware转换成gin的中间件返回值
	return adapter.Wrap(csrfMiddleware)
}

/*

扩展学习:  grpc的CSRF 验证  https://cloud.tencent.com/developer/article/1921883

*/
