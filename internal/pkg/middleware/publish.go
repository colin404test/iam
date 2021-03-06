// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/json"

	"github.com/marmotedu/iam/internal/authzserver/store/load"
	"github.com/marmotedu/iam/pkg/log"
	"github.com/marmotedu/iam/pkg/storage"
)

// Publish publish a redis event to specified redis channel when some action occurred.
func Publish() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		var resource string

		pathSplit := strings.Split(c.Request.URL.Path, "/")
		if len(pathSplit) > 2 {
			resource = pathSplit[2]
		}

		method := c.Request.Method

		switch resource {
		case "policies":
			notify(method, load.NoticePolicyChanged)
		case "secrets":
			notify(method, load.NoticeSecretChanged)
		default:
		}
	}
}

func notify(method string, command load.NotificationCommand) {
	switch method {
	case "POST", "PUT", "DELETE", "PATH":
		redisStore := &storage.RedisCluster{}
		message, _ := json.Marshal(load.Notification{Command: command})

		if err := redisStore.Publish(load.RedisPubSubChannel, string(message)); err != nil {
			log.Errorw("publish redis message failed", "error", err.Error())
		}
		log.Debugw("publish redis message", "method", method, "command", command)
	default:
	}
}
