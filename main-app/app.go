package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type jsonRequest struct {
	Url string
}

type jsonFormat struct {
	Url        string
	ShortenUrl string
}

func middlewareApi(normal_ctx context.Context, redis_cli *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("ctx", normal_ctx)
		ctx.Set("redis", redis_cli)
		ctx.Next()
	}

}

func redirect(c *gin.Context) {
	redis_cli := c.MustGet("redis").(*redis.Client)
	ctx := c.MustGet("ctx").(context.Context)

	url_name := c.Params.ByName("url")

	full_shorten_url := c.Request.Host + "/" + url_name

	redis_data, err := redis_cli.HGetAll(ctx, full_shorten_url).Result()

	if err != nil {
		c.JSON(405, gin.H{})
		return
	}

	c.Redirect(
		303,
		redis_data["url"],
	)
	c.JSON(200, gin.H{
		"shorten url": full_shorten_url,
	})

}

func shortenUrl(c *gin.Context) {

	json_data_reader, json_err := io.ReadAll(c.Request.Body)

	err_handler(json_err)

	var json_data_request jsonRequest
	err_unmarshal := json.Unmarshal(json_data_reader, &json_data_request)
	err_handler(err_unmarshal)

	_, err := url.ParseRequestURI(json_data_request.Url)
	fmt.Println(json_data_request.Url)
	err_handler(err)

	res := c.Request.Host + "/" + randSeq(8)
	json_data := jsonFormat{
		Url:        json_data_request.Url,
		ShortenUrl: res,
	}

	redis_cli := c.MustGet("redis").(*redis.Client)
	ctx := c.MustGet("ctx").(context.Context)
	redis_err := redis_cli.HSet(ctx, res, map[string]string{"url": json_data.Url, "shortenUrl": json_data.ShortenUrl}).Err()

	err_handler(redis_err)
	// re_val, err := redis_cli.HGetAll(ctx, res).Result()

	// err_handler(err)

	// fmt.Println(res, re_val["url"])

	c.Bind(json_data)

	c.JSON(200, gin.H{
		"url":         json_data.Url,
		"shorten_url": json_data.ShortenUrl,
	})

}

func StartApp() {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Username: "default",
		Password: "",
		DB:       1,
	})

	ctx := context.Background()

	r := gin.New()

	r.Use(middlewareApi(ctx, redisClient))

	// r.GET("/ping", returnJson)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r.POST("/short", shortenUrl)
	r.GET("/:url", redirect)

	r.Run(":8080")

}
