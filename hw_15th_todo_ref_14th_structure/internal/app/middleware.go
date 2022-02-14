package app

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (s *Server) Temp() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload_map := make(map[string]string)
		f, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Error().Caller().Str("func", "io.ReadAll").Err(err).Msg("Server")

		}
		err = json.Unmarshal(f, &payload_map)
		if err != nil {
			log.Error().Caller().Str("func", "Unmarshal").Err(err).Msg("Server")

		}
		payload_map["user_id"] = "1"
		jf, err := json.Marshal(payload_map)
		c.Request.Body = io.NopCloser(bytes.NewReader(jf))
	}
}
