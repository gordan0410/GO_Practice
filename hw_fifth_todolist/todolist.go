package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 資料庫連接資料
type Db_message struct{
	Username string 
	Password string 
	Addr string 
	Database string 
	Max_lifetime int 
	Max_openconns int 
	Max_idleconns int 
}

// 建立table
type Todolist struct {
	ID			uint         	`gorm:"type:bigint(20) NOT NULL auto_increment;primary_key;"`
	Subject		string       	`gorm:"type:varchar(30) NOT NULL;"`
	Status		int       		`gorm:"type:int NOT NULL;"`
	CreatedAt	time.Time    	`gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP"`
}

// db連線
var db_conn *gorm.DB
var db_conn_err error

// 幾筆資料一頁
var slice_target = 4

func main() {
	// 開啟DB(user_info.json)
	load_database("./config.json")
	// 開啟gin
	set_router()
}

// 錯誤處理
func mistake_control(err error, line int ,func_name string){
	if err != nil{
		fmt.Printf("line %d, %s failed: %v", line, func_name, err)
		return
	} else {
		return
	}
}

func load_database(addr string){
	// 讀取config.json資料
	file, err := os.Open("./config.json")
	if err != nil {
		mistake_control(err, 62, "os.Open")
		return
	}
	var msg Db_message
	d_data := json.NewDecoder(file)
	err = d_data.Decode(&msg)
	if err != nil {
		mistake_control(err, 69, "json.Decode")
		return
	}
	
	// 資料庫連接
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", msg.Username, msg.Password, msg.Addr, msg.Database)
	db_conn, db_conn_err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if db_conn_err != nil{
		mistake_control(db_conn_err, 77, "gorm.Open")
		return
	}
	db, db_err := db_conn.DB()
	if db_err != nil{
		mistake_control(db_err, 82, "db_conn.DB()")
		return
	}

	// 確保資料庫安全關閉
	db.SetConnMaxLifetime(time.Duration(msg.Max_lifetime) * time.Second)
	// 閒置連接數（官方建議跟SetMaxOpenConns一致）
	db.SetMaxIdleConns(msg.Max_idleconns)
	// 限制資料庫連接數
	db.SetMaxOpenConns(msg.Max_openconns)
	
	// 更新資料庫資料
	db_conn.Debug().AutoMigrate(&Todolist{})

	// 確認table存在
	migrator := db_conn.Migrator()
	has := migrator.HasTable(&Todolist{})
	if !has {
		fmt.Println("table not exist")
		return
	}
}

func set_router(){
	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		c.File("./template")
	})
	r.GET("/api", get)  // 首頁
	r.POST("/api", create)  // 建立
	r.PATCH("/api", update)  // 更新
	r.DELETE("/api", delete)	// 刪除

	err := r.Run(":3000")
	if err != nil{
		mistake_control(err, 117, "router.Run")
		return
	}
}

// DB搜尋方法，參數為前端回傳group資料＋DB資料限制讀取數＋資料忽略數（之後可再把limit and offset給前端做調整）
func db_search_group(c *gin.Context, limit_record, offset_record int) func(db *gorm.DB) *gorm.DB{
	return func(db *gorm.DB) *gorm.DB{
		group := c.Query("group")
		switch {
		case group == "all":
			return db.Where("Status <> ?", 0).Offset(offset_record).Limit(limit_record)
		case group == "active":
			return db.Where("Status = ?", 1).Offset(offset_record).Limit(limit_record)
		case group == "complete":
			return db.Where("Status = ?", 2).Offset(offset_record).Limit(limit_record)
		default:
			return db.Where("Status <> ?", 0).Offset(offset_record).Limit(limit_record)
		}
	}
}

// 設定db頁面呈現方式，參數為前端回傳當前page資料＋欲顯示分頁數（之後可再把分頁數給前端做調整）
func db_search_page(c *gin.Context, slice_target int) func(db *gorm.DB) *gorm.DB{
	return func(db *gorm.DB) *gorm.DB{
		msg := c.Query("page")
		page, err := strconv.Atoi(msg)
		mistake_control(err ,145, "strconv.Atoi")
		offset := (page-1)*slice_target
		return db.Offset(offset).Limit(slice_target)
	}
}

func get(c *gin.Context) {

	var todolists []Todolist

	// 查詢所有非0（軟刪除的）data
	result := db_conn.Debug().Order("id desc").Scopes(db_search_group(c, 100, 0),db_search_page(c, slice_target)).Find(&todolists)	
	if result.Error != nil{
		mistake_control(result.Error, 157, "db query")
		c.String(http.StatusBadRequest, "db query error")
		return
	}

	if len(todolists)> 0{
		c.JSON(http.StatusOK, todolists)
		fmt.Println("data sent: ",todolists)
	} else {
		c.String(http.StatusBadRequest, "no other page")
		fmt.Println("no other page")
		return
	}
}


func create(c *gin.Context) {
	// 接收前端訊息
	subject := c.PostForm("subject")
	if subject ==""{
		c.String(http.StatusBadRequest, "subject is empty")
		fmt.Println("subject is empty")
		return
	}

	// 建立資料並寫入
	create := Todolist{Subject: subject, Status: 1}
	result := db_conn.Debug().Create(&create)
		if result.Error != nil {
			mistake_control(result.Error, 186, "db_conn.Debug().Create(&create)")
			c.String(http.StatusBadRequest, "db create error")
			return
		}
		if result.RowsAffected != 1 {
			fmt.Println("add more than one object")
			c.String(http.StatusBadRequest, "bad request")
			return
		}
		c.String(http.StatusOK, "created")
		fmt.Println("new data created: ", create)
}

func update(c *gin.Context) {
	// 接收前端訊息
	id := c.PostForm("id")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		mistake_control(err, 204, "strconv.Atoi")
		c.String(http.StatusBadRequest, "id problem")
		return
	}

	status := c.PostForm("status")
	if status == "" {
		fmt.Println("status is empty")
		c.String(http.StatusBadRequest, "status is empty")
		return
	}
	
	var todolist Todolist

	// 先找出欲異動資料
	result := db_conn.Debug().Where("ID = ? AND Status <> ?", uint(id_int), 0).Take(&todolist)
	if result.Error != nil{
		if errors.Is(result.Error, gorm.ErrRecordNotFound){
			mistake_control(result.Error, 221, "subject can't find")
			c.String(http.StatusBadRequest, "can't find this subject")
			return
		}else {
			mistake_control(result.Error, 221, "db error")
			c.String(http.StatusBadRequest, "db error")
			return
		}
	}
		
	// 判斷訊息類型為"complete" or "active"
	// complete則更新指定ID的status為2（已完成的）
	if status == "complete" {
		// 更新
		result = db_conn.Debug().Model(&todolist).Update("Status", "2")
		if result.Error != nil{
			mistake_control(result.Error, 238, "db updated failed")
			c.String(http.StatusBadRequest, "db updated failed")
			return
		}
		c.String(http.StatusOK, "updated")
		fmt.Println("data updated: ", todolist)

		// active則更新指定ID的status為1（未完成的）
	} else if status == "active" {
		// 更新
		result = db_conn.Debug().Model(&todolist).Update("Status", "1")
		if result.Error != nil{
			mistake_control(result.Error, 250, "db updated failed")
			c.String(http.StatusBadRequest, "db updated failed")
			return
		}
		c.String(http.StatusOK, "updated")
		fmt.Println("data updated: ", todolist)

	// 判斷是否為更改標題
	} else if status == "subject_change" {
		subject := c.PostForm("subject")
		if subject == "" {
			fmt.Println("subject is empty")
			c.String(http.StatusBadRequest, "subject is empty")
			return
		}
		result = db_conn.Debug().Model(&todolist).Update("Subject", subject)
		if result.Error != nil{
			mistake_control(result.Error, 267, "db updated failed")
			c.String(http.StatusBadRequest, "db updated failed")
			return
		}
		c.String(http.StatusOK, "updated")
		fmt.Println("data updated: ", todolist)

	// 非以上則錯誤
	} else {
		fmt.Println("unknow action type")
		c.String(http.StatusBadRequest, "can't find this status")
		return
	}
}

func delete(c *gin.Context) {
	var todolist []Todolist
	count := 0
	// 若資料量>100, 則分批刪除
	for{
		result := db_conn.Debug().Where("Status = ?", 2).Limit(100).Select("ID", "Status").Find(&todolist)
		if result.Error != nil{
			if errors.Is(result.Error, gorm.ErrRecordNotFound){
				mistake_control(result.Error, 221, "subject can't find")
				c.String(http.StatusBadRequest, "can't find this subject")
				break
			}else {
				fmt.Println("line 380, result.Error", result.Error)
				break
			}
		}
		// 判斷是否有資料，無資料則回傳"no object selected"，以提示無資料經選取
		if len(todolist) == 0 && count == 0{
			fmt.Println("no object selected")
			c.String(http.StatusBadRequest, "no object selected")
			break
		// 無剩餘可刪除，刪除結束	
		} else if len(todolist) == 0 && count > 0{
			fmt.Println("deleted ok, no object remain")
			c.String(http.StatusOK, "all deleted")
			break
		// 有資料則更改status為0（軟刪除）
		} else{
			result = db_conn.Debug().Model(&todolist).Update("Status", "0")
			if result.Error != nil{
				mistake_control(result.Error, 312, "db error")
				c.String(http.StatusBadRequest, "db error")
				break
			}
			count++
		}
	}
}
