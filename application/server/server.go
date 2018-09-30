package server

import (
	"github.com/gin-gonic/gin"

	"github.com/mrl-athomelab/website/application/config"
	"github.com/mrl-athomelab/website/application/template"
)

type Server struct {
	config *config.Configuration
	router *gin.Engine
}

func autoReloadTemplate(engine *gin.Engine, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		engine.SetHTMLTemplate(template.New(path))
		c.Next()
	}
}

func Prepare(configFile string) (s *Server, err error) {
	s = &Server{}
	s.config, err = config.Read(configFile)
	if err != nil {
		return
	}

	if !s.config.Development {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery(), gin.Logger())

	if s.config.Development {
		engine.Use(autoReloadTemplate(engine, s.config.TemplatePath))
	}

	engine.SetHTMLTemplate(template.New(s.config.TemplatePath))

	engine.Static("/static", s.config.StaticPath)
	engine.GET("/", s.staticRender("index"))

	s.router = engine
	return
}

func (s *Server) Run() error {
	return s.router.Run(s.config.ListenAddr)
}
