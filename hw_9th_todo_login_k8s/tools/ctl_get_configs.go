package tools

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func Get_configs(c *gin.Context) (*Config_data, error) {
	configs_raw, b := c.Get("configs")
	if !b {
		err := errors.New("can't find configs")
		return nil, err
	}
	configs := configs_raw.(*Config_data)
	return configs, nil
}
