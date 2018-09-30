package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mrl-athomelab/website/application/secure"
)

func (s *Server) render(ctx *gin.Context, code int, name string, data gin.H) {
	ctx.HTML(code, name, map[string]interface{}{
		"page":         name,
		"data":         data,
		"unixtime":     time.Now().Unix(),
		"current_path": ctx.Request.RequestURI,
		"secure_token": secure.MD5Hash(ctx.Request.RemoteAddr + s.config.SecretKey),
	})
}

func (s *Server) staticRender(page string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s.render(ctx, http.StatusOK, page, nil)
	}
}
