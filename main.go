package main

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"io"
	"net/http"
	"net/url"
	"os"
)

var (
	HOST string
	PORT string
)

func init() {
	if HOST = os.Getenv("HOST"); HOST == "" {
		HOST = "127.0.0.1"
	}

	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "8080"
	}
}

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ProxyHeader:           "X-Forwarded-For",
	})

	// logger
	app.Use(logger.New())

	// proxy access
	app.All("/*", func(c *fiber.Ctx) error {
		destURL, err := url.Parse(c.GetReqHeaders()["Dest-Addr"])
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		fullURL := fmt.Sprintf("%s%s", destURL.String(), c.Path())
		params := c.Request().URI().QueryArgs().String()
		data := c.Request().Body()

		// 创建代理请求
		request, _ := http.NewRequest(c.Method(), fullURL, bytes.NewReader(data))
		request.URL.RawQuery = params

		// 复制请求头
		for k, v := range c.GetReqHeaders() {
			request.Header.Set(k, v)
		}
		request.Header.Set("Host", destURL.Host)
		request.Header.Set("Origin", destURL.Scheme+"://"+destURL.Host)
		request.Header.Set("Referer", fullURL)

		// 删除额外头 Dest-Addr
		request.Header.Del("Dest-Addr")

		// 发送代理请求
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		// 复制响应头
		for k, v := range response.Header {
			c.Set(k, v[0])
		}

		// 复制响应体
		_, err = io.Copy(c.Response().BodyWriter(), response.Body)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		// 返回响应
		return nil
	})

	_ = app.Listen(fmt.Sprintf("%s:%s", HOST, PORT))
}
