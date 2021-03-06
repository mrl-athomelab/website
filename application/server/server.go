package server

import (
	htmlTemplate "html/template"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	"github.com/mrl-athomelab/website/application/config"
	"github.com/mrl-athomelab/website/application/database"
	"github.com/mrl-athomelab/website/application/secure"
	"github.com/mrl-athomelab/website/application/template"
)

type Server struct {
	config   *config.Configuration
	router   *gin.Engine
	template *htmlTemplate.Template
}

func autoReloadTemplate(s *Server, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		s.template = template.New(path)
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

	err = database.Open(s.config.Database.Provider, s.config.Database.ConnString)
	if err != nil {
		return
	}
	administrator := &database.Administrator{Username: "admin", Password: secure.MD5Hash("admin", s.config.SecretKey)}
	exists := administrator.Get(database.ByUsernamePassword)
	if !exists {
		err = administrator.Save()
		if err != nil {
			return
		}
	}

	engine := gin.New()
	engine.Use(gin.Recovery(), gin.Logger())
	engine.Use(rayTokenGenerator())

	s.template = template.New(s.config.TemplatePath)
	if s.config.Development {
		engine.Use(autoReloadTemplate(s, s.config.TemplatePath))
	}

	engine.Static("/static", s.config.StaticPath)
	engine.GET("/", s.staticRender("index"))

	pages := engine.Group("/pages")
	{
		pages.GET("/blog")
		pages.GET("/gallery")
		pages.GET("/aboutus", s.staticRender("aboutus"))
		pages.GET("/join")
	}

	admin := engine.Group("/admin")
	{
		admin.Any("/login", s.adminLoginAnyHandler)
		admin.GET("/panel", s.adminPanelGetHandler, s.adminAuthMiddleware())

		administrators := admin.Group("/administrators", s.adminAuthMiddleware())
		{
			administrators.GET("/", s.adminAdministratorsGetHandler)
			rest := administrators.Group("/rest")
			{
				rest.POST("/", s.adminAdministratorsRestPostHandler)
				rest.DELETE("/:id", s.adminAdministratorsRestDeleteHandler)
			}
		}

		members := admin.Group("/members", s.adminAuthMiddleware())
		{
			members.GET("/", s.adminMembersGetHandler)
			members.POST("/image", s.adminMembersImagePostHandler)

			members.GET("/edit/:id", s.adminMembersEditGetHandler)

			rest := members.Group("/rest")
			{
				rest.POST("/", s.adminMembersRestPostHandler)
				rest.DELETE("/:id", s.adminMembersRestDeleteHandler)
			}
		}

		news := admin.Group("/news", s.adminAuthMiddleware())
		{
			news.GET("/", s.adminNewsGetHandler)
			news.GET("/edit/:id", s.adminNewsEditGetHandler)

			rest := news.Group("/rest")
			{
				rest.POST("/", s.adminNewsRestPostHandler)
				rest.PUT("/:id", s.adminNewsRestPutHandler)
				rest.DELETE("/:id", s.adminNewsRestDeleteHandler)
			}
		}
	}

	s.router = engine
	return
}

func (s *Server) Run() error {
	return s.router.Run(s.config.ListenAddr)
}

func rayTokenGenerator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		guid := xid.New()
		ctx.Set("raytoken", guid.String())
		ctx.Next()
	}
}
