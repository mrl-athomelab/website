package server

import (
	"bytes"
	"net/http"
	"time"

	"github.com/mrl-athomelab/website/application/logger"

	"github.com/gin-gonic/gin"

	"github.com/mrl-athomelab/website/application/secure"
)

func (s *Server) render(ctx *gin.Context, code int, name string, data gin.H) {
	buf := new(bytes.Buffer)

	err := s.template.ExecuteTemplate(buf, name, map[string]interface{}{
		"page":         name,
		"data":         data,
		"unixtime":     time.Now().Unix(),
		"current_path": ctx.Request.RequestURI,
		"secure_token": secure.MD5Hash(ctx.Request.RemoteAddr + s.config.SecretKey),
		"csrf_token":   secure.GenerateToken(ctx.GetString("raytoken"), s.config.SecretKey),
	})
	if err != nil {
		logger.Warn("Error on executing template %s ...", name)
		logger.Error("%v", err)
		ctx.String(http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	buf.WriteTo(ctx.Writer)
}

func (s *Server) staticRender(page string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s.render(ctx, http.StatusOK, page, nil)
	}
}
