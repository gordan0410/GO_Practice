package todolist

import (
	"fmt"
	"hw_ninth/models"
	"hw_ninth/tools"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func get(c *gin.Context) {
	// 拿取user_id
	user_id_raw, b := c.GetQuery("user_id")
	if !b {
		log.Error().Caller().Str("func", "c.GetQuery(\"user_id\")").Str("msg", "get user id failed").Msg("Web")
		if err := tools.Msg_send(c, "error", "get user id failed", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	// set user_id for other function
	c.Set("user_id", user_id_raw)
	user_id, err := tools.Parse_int(user_id_raw)
	if err != nil {
		log.Error().Caller().Str("func", "tools.Parse_int(user_id_raw)").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "user_id convert error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	// 取得分群
	group, b := c.GetQuery("group")
	if !b {
		log.Error().Caller().Str("func", "c.GetQuery(\"group\")").Str("msg", "get group failed").Msg("Web")
		if err := tools.Msg_send(c, "error", "get group failed", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	// 取得頁面顯示數
	slice_target_raw, b := c.GetQuery("slice_target")
	if !b {
		log.Error().Caller().Str("func", "c.GetQuery(\"slice_target\")").Str("msg", "get slice_target failed").Msg("Web")
		if err := tools.Msg_send(c, "error", "get slice_target failed", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}
	slice_target, err := tools.Parse_int(slice_target_raw)
	if err != nil {
		log.Error().Caller().Str("func", "tools.Parse_int(slice_target_raw)").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "slice_target convert error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	// 取得當前頁數
	page_raw, b := c.GetQuery("page")
	if !b {
		log.Error().Caller().Str("func", "c.GetQuery(\"page\")").Str("msg", "get page failed").Msg("Web")
		if err := tools.Msg_send(c, "error", "get page failed", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}
	page, err := tools.Parse_int(page_raw)
	if err != nil {
		log.Error().Caller().Str("func", "tools.Parse_int(page_raw)").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "page convert error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	// 無<1的頁面
	if page <= 0 {
		if err := tools.Msg_send(c, "error", "no other page", nil); err != nil {
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

	var todolists []models.Todolist
	// 查詢所有非0（軟刪除的）data
	err = db_conn.Order("id desc").Where("account_id = ? ", user_id).Scopes(tools.Db_query_group(group, 100, 0), tools.Db_query_page(page, slice_target)).Find(&todolists).Error
	if err != nil {
		log.Error().Caller().Str("func", "db_conn.Order(\"id desc\")").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "db query error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	// 根據資料量表示是否有後續分頁
	if len(todolists) > 0 {
		data := make(map[string]interface{})
		for i, v := range todolists {
			si := fmt.Sprint(i)
			data[si] = gin.H{"id": v.ID,
				"status":  v.Status,
				"subject": v.Subject}
		}
		if err := tools.Msg_send(c, "success", "list send", data); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			return
		}
		// 第一次啟動
	} else if len(todolists) == 0 && page == 1 {
		if err := tools.Msg_send(c, "success", "no data send", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			return
		}
		// 無分頁
	} else {
		if err := tools.Msg_send(c, "error", "no other page", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			return
		}
	}
}
