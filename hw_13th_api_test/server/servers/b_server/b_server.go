package b_server

import (
	"log"
	"net/http"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

type msg struct {
	Larry string `json:"larry" binding:"required"`
}

func B_server(port string) error {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(auth)
	r.GET("/v1/person", vi_person)
	r.POST("/v1/book", vi_book)
	err := endless.ListenAndServe(":"+port, r)
	if err != nil {
		return err
	}
	return nil
}

type req_msg struct {
	Number string `json:"number"`
}

type res_msg[T result_all] struct {
	Status string `json:"status"`
	Text   string `json:"text"`
	Result T      `json:"result"`
}

type result_all interface {
	person | book | map[string]string
}

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Job  string `json:"job"`
}

type book struct {
	Name     string `json:"name"`
	Position string `json:"position"`
}

// id錯誤
var id_fail = res_msg[map[string]string]{
	Status: "N",
	Text:   "error",
	Result: nil,
}

// 驗證失敗資訊
var token_fail = res_msg[map[string]string]{
	Status: "N",
	Text:   "error",
	Result: nil,
}

// 驗證
func auth(c *gin.Context) {
	token := c.Request.Header.Get("api-key")
	if token == "321" {
		c.Next()
	} else {
		c.JSON(http.StatusOK, token_fail)
		c.Abort()
		return
	}
}

func vi_person(c *gin.Context) {
	// 成功資訊
	var success = res_msg[person]{
		Status: "Y",
		Text:   "ok",
		Result: person{Name: "Larry", Age: 99, Job: "凱瑞"},
	}

	// 辨識身份
	id := c.Query("no")
	if id == "1" {
		c.JSON(http.StatusOK, success)
	} else {
		c.JSON(http.StatusOK, id_fail)
	}
}

func vi_book(c *gin.Context) {
	// 成功資訊
	success := res_msg[book]{
		Status: "Y",
		Text:   "ok",
		Result: book{Name: "Larry", Position: "1"},
	}

	// 辨識身份
	var r req_msg
	err := c.ShouldBindJSON(&r)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, id_fail)
	}
	if r.Number == "1" {
		c.JSON(http.StatusOK, success)
	} else {
		c.JSON(http.StatusOK, id_fail)
	}
}
