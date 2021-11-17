package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 資料庫連接資料
const (
	Username string = "root"
	Password string = "root"
	Addr string = "localhost:3306"
	Database string = "test"
	Max_lifetime int = 10
	Max_openconns int = 10
	Max_idleconns int = 10
)

// 建立table
type Todolist struct {
	ID			uint         	`gorm:"type:bigint(20) NOT NULL auto_increment;primary_key;"`
	Subject		string       	`gorm:"type:varchar(30) NOT NULL;"`
	Status		int       		`gorm:"type:int NOT NULL;"`
	CreatedAt	time.Time    	`gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP"`
}

// 接收create Json
type Msg_new struct{
	Type, Subject string
}

// 接收update Json
type Msg_update struct{
	Type, Id string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// 資料庫連接
var dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", Username, Password, Addr, Database)
var db_conn, db_conn_err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
var db, db_err = db_conn.DB()

// 幾筆資料一頁
var slice_target = 3

func main() {
	if db_conn_err != nil{
		println("gorm.Open failed: ", db_conn_err)
		return
	}
	if db_err != nil{
		println("conn.DB failed: ", db_err)
		return
	}

	// 確保資料庫安全關閉
	db.SetConnMaxLifetime(time.Duration(Max_lifetime) * time.Second)
	// 閒置連接數（官方建議跟SetMaxOpenConns一致）
	db.SetMaxIdleConns(Max_idleconns)
	// 限制資料庫連接數
	db.SetMaxOpenConns(Max_openconns)
	
	// 更新資料庫資料
	db_conn.Debug().AutoMigrate(&Todolist{})

	// 確認table存在
	migrator := db_conn.Migrator()
	has := migrator.HasTable(&Todolist{})
	if !has {
		fmt.Println("table not exist")
		return
	}
	// 開啟gin
	set_router()
}

func set_router(){
	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		c.File("./template")
	})
	r.GET("/ws", ws_get)  // 首頁
	r.GET("/ws_page", ws_page)  // 換頁
	r.GET("/ws_create", ws_create)  // 建立
	r.GET("/ws_update", ws_update)  // 更新
	r.GET("/ws_delete", ws_delete)	// 刪除
	err := r.Run(":3000")
	if err != nil{
		println("line 99, can't run route: ", err)
	}
}

func ws_get(c *gin.Context) {
	if websocket.IsWebSocketUpgrade(c.Request) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("line 107, websocket upgrader.Upgrade failed: ", err)
			return
		}
		defer conn.Close()
		
		var todolists []Todolist
		
		// 查詢所有非0（軟刪除的）data
		result := db_conn.Debug().Where("Status <> ?", 0).Find(&todolists)
		if result.Error != nil{
			if errors.Is(result.Error, gorm.ErrRecordNotFound){
				fmt.Println("line 117, can't find this subject", result.Error)
				return
			}else {
				fmt.Println("line 117, result.Error", result.Error)
				return
			}
		}

		// reverse 查詢後的資料並依照slice_target數呈現資料
		reverse_todolists := reverse(todolists)
		slice_data := reverse_todolists[:slice_target]
		err = conn.WriteJSON(slice_data)
		if err != nil {
			fmt.Println("line 131, conn.WriteJson failed:", err)
			return
		}
	} else {
		fmt.Println("line 106, can't connect websocket")
		return
	}
}

func ws_page(c *gin.Context){
	if websocket.IsWebSocketUpgrade(c.Request) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("line 144, websocket upgrader.Upgrade failed: ", err)
			return
		}
		defer conn.Close()
		
		var todolists []Todolist

		// 查詢所有非0（軟刪除的）data
		result := db_conn.Debug().Where("Status <> ?", 0).Find(&todolists)
		if result.Error != nil{
			if errors.Is(result.Error, gorm.ErrRecordNotFound){
				fmt.Println("line 154, can't find this subject", result.Error)
				return
			}else {
				fmt.Println("line 154, result.Error", result.Error)
				return
			}
		}
		// reverse querying data 並根據前端頁面回傳資料切片
		reverse_todolists := reverse(todolists)
		_, page, err:=conn.ReadMessage()
		if err != nil {
			fmt.Println("line 166, conn.ReadMessage failed: ", err)
			return
		}
		page_int, err := strconv.Atoi(string(page))
		if err != nil {
			fmt.Println("line 171, strconv.Atoi failed: ", err)
			return
		}
		total_object := len(reverse_todolists)
		add := 0
		if total_object % slice_target > 0 {
			add = 1
		}
		total_page := total_object / slice_target + add
		if slice_target * page_int <= total_object && slice_target * page_int > 0 {
			slice_data := reverse_todolists[slice_target * (page_int - 1) : slice_target * page_int]
			err = conn.WriteJSON(slice_data)
			if err != nil {
				fmt.Println("line 184, conn.WriteJson failed:", err)
				return
			}
		} else if slice_target * page_int > total_object && page_int == total_page {
			slice_data := reverse_todolists[slice_target * (page_int - 1) : total_object]
			err = conn.WriteJSON(slice_data)
			if err != nil {
				fmt.Println("line 191, conn.WriteJson failed:", err)
				return
			}
		} else {
			fmt.Println("no other page")
			err = conn.WriteMessage(websocket.TextMessage, []byte("no other page"))
			if err != nil {
				fmt.Println("line 198, conn.Writemessage failed:", err)
				return
			}
		}
	} else {
		fmt.Println("can't connect websocket")
		return
	}
}

// reverse querying Todolist struct for ws_get and we_page
func reverse(array []Todolist) []Todolist {
	for i, j := 0, len(array)-1; i < j; i, j = i+1, j-1 {
		array[i], array[j] = array[j], array[i]
	}
	return array
}

func ws_create(c *gin.Context) {
	if websocket.IsWebSocketUpgrade(c.Request) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("line 220, websocket upgrader.Upgrade failed: ", err)
			return
		}
		defer conn.Close()

		// 接收前端訊息
		var msg Msg_new
		err = conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("line 229, conn.ReadJson failed:", err)
			return
		}
		// 判斷訊息類型是否為"create"，若是則新增
		if msg.Type == "create" {
			create := Todolist{Subject: msg.Subject, Status: 1}
			result := db_conn.Debug().Create(&create)
			if result.Error != nil {
				fmt.Println("line 237, Create failt")
				return
			}
			if result.RowsAffected != 1 {
				fmt.Println("line 237, RowsAffected Number failt")
				return
			}
			conn.WriteMessage(websocket.TextMessage, []byte("create success"))
			if err != nil {
				fmt.Println("line 246, conn.WriteMessage falied: ", err)
				return
			}
		// 若無則回傳錯誤
		} else {
			fmt.Println("message type wrong")
			conn.WriteMessage(websocket.TextMessage, []byte("message type wrong"))
			if err != nil {
				fmt.Println("line 254, conn.WriteMessage falied: ", err)
				return
			}
		}
	} else {
		fmt.Println("can't connect websocket")
		return
	}
}

func ws_update(c *gin.Context) {
	if websocket.IsWebSocketUpgrade(c.Request) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("line 268, websocket upgrader.Upgrade failed: ", err)
			return
		}
		defer conn.Close()
		
		// 接收前端訊息
		var msg Msg_update
		err = conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("line 277, conn.ReadJson failed:", err)
			return
		}

		var todolist Todolist
		// 判斷訊息類型為"complete" or "active"
		// complete則更新指定ID的status為2（已完成的）
		if msg.Type == "complete" {
			// 轉前端id為int
			id_int, err := strconv.Atoi(msg.Id)
			if err != nil {
				fmt.Println("line 288, ID convert to int failed")
				return
			}
			// 查詢指定ID是否存在
			result := db_conn.Debug().Where("ID = ?", uint(id_int)).Take(&todolist)
			if result.Error != nil{
				if errors.Is(result.Error, gorm.ErrRecordNotFound){
					fmt.Println("line 294, can't find this subject", result.Error)
					return
				}else {
					fmt.Println("line 294, result.Error", result.Error)
					return
				}
			}
			// 更新
			result = db_conn.Debug().Model(&todolist).Update("Status", "2")
			if result.Error != nil{
				fmt.Println("line 305, result.Error", result.Error)
				return
			}
			conn.WriteMessage(websocket.TextMessage, []byte("data updated"))
			if err != nil {
				fmt.Println("line 305, conn.WriteMessage falied: ", err)
				return
			}
		// active則更新指定ID的status為1（未完成的）
		} else if msg.Type == "active" {
			// 轉前端id為int
			id_int, err := strconv.Atoi(msg.Id)
			if err != nil {
				fmt.Println("line 318, ID convert to int failed")
				return
			}
			// 查詢指定ID是否存在
			result := db_conn.Debug().Where("ID = ?", uint(id_int)).Take(&todolist)
			if result.Error != nil{
				if errors.Is(result.Error, gorm.ErrRecordNotFound){
					fmt.Println("line 324, can't find this subject", result.Error)
					return
				}else {
					fmt.Println("line 324, result.Error", result.Error)
					return
				}
			}
			// 更新
			result = db_conn.Debug().Model(&todolist).Update("Status", "1")
			if result.Error != nil{
				fmt.Println("line 335, result.Error", result.Error)
				return
			}
			conn.WriteMessage(websocket.TextMessage, []byte("data updated"))
			if err != nil {
				fmt.Println("line 340, conn.WriteMessage falied: ", err)
				return
			}
		// 非"complete" or "active"則錯誤
		} else {
			fmt.Println("unknow action type")
			conn.WriteMessage(websocket.TextMessage, []byte("updated in db failed"))
			if err != nil {
				fmt.Println("line 348, conn.WriteMessage falied: ", err)
				return
			}
		}
	} else {
		fmt.Println("can't connect websocket")
		return
	}
}

func ws_delete(c *gin.Context) {
	// 驗證websocket
	if websocket.IsWebSocketUpgrade(c.Request) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("line 363, websocket upgrader.Upgrade failed: ", err)
			return
		}
		defer conn.Close()
		
		// 讀取message是否為"clear"
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("line 371, conn.ReadMessage failed:", err)
			return
		}

		var todolist []Todolist
		if string(msg) == "clear" {
			// Query status為2（已完成的）的資料
			result := db_conn.Debug().Where("Status = ?", 2).Select("ID").Find(&todolist)
			if result.Error != nil{
				if errors.Is(result.Error, gorm.ErrRecordNotFound){
					fmt.Println("line 380, can't find this subject", result.Error)
					return
				}else {
					fmt.Println("line 380, result.Error", result.Error)
					return
				}
			}
			// 判斷是否有資料，無資料則回傳"no object selected"，以提示無資料經選取
			if len(todolist) == 0 {
				fmt.Println("no object selected")
				err := conn.WriteMessage(websocket.TextMessage, []byte("no object selected"))
				if err != nil {
					fmt.Println("line 393, conn.WriteMessage falied: ", err)
					return
				}
				return
			// 有資料則更改status為0（軟刪除）
			} else{
				result = db_conn.Debug().Model(&todolist).Update("Status", "0")
				if result.Error != nil{
					fmt.Println("line 401, result.Error", result.Error)
					return
				}
				err := conn.WriteMessage(websocket.TextMessage, []byte("data updated"))
				if err != nil {
					fmt.Println("line 406, conn.WriteMessage falied: ", err)
					return
				}
			}
		}
	} else {
		fmt.Println("can't connect websocket")
		return
	}
}
