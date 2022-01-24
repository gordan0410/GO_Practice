package todolist

import (
	"hw_ninth/models"
	"hw_ninth/tools"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Delete_msg struct {
	User_id string `json:"user_id"`
}

func delete(c *gin.Context) {
	var dm Delete_msg
	// 接收前端訊息
	if err := c.ShouldBindJSON(&dm); err != nil {
		log.Error().Caller().Str("func", "ShouldBindJSON(&subject)").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "subject type error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	user_id, err := tools.Parse_int(dm.User_id)
	if err != nil {
		log.Error().Caller().Str("func", "tools.Parse_int(user_id_raw)").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "user_id convert error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	// 開DB
	db_conn, err := models.Db_conn_web(c)
	defer models.Db_close(db_conn)
	if err != nil {
		log.Error().Caller().Str("func", "tools.Db_conn(c)").Err(err).Msg("Db")
		if err := tools.Msg_send(c, "error", "db error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	var todolist []models.Todolist
	count := 0

	// 若資料量>100, 則分批刪除
	for {
		err := db_conn.Where("Status = ? AND account_id = ? ", 2, user_id).Limit(100).Select("ID", "Status").Find(&todolist).Error
		if err != nil {
			log.Error().Caller().Str("func", "db_conn.Where(\"Status = ?\", 2)").Err(err).Msg("Web")
			if err := tools.Msg_send(c, "error", "db query error", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			}
			break
		}
		// 判斷是否有資料，無資料則回傳"no object selected"，以提示無資料經選取
		if len(todolist) == 0 && count == 0 {
			if err := tools.Msg_send(c, "error", "no object selected", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			}
			break
			// 無剩餘可刪除，刪除結束
		} else if len(todolist) == 0 && count > 0 {
			if err := tools.Msg_send(c, "success", "all deleted", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			}
			break
			// 有資料則更改status為0（軟刪除）並count++
		} else {
			err = db_conn.Model(&todolist).Update("Status", "0").Error
			if err != nil {
				log.Error().Caller().Str("func", " db_conn.Model(&todolist)").Err(err).Msg("Web")
				if err := tools.Msg_send(c, "error", "db deleted error", nil); err != nil {
					log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
				}
				break
			}
			count++
		}
	}
}
