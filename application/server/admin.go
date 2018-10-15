package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mrl-athomelab/website/application/cookie"
	"github.com/mrl-athomelab/website/application/database"
	"github.com/mrl-athomelab/website/application/jsonresponse"
	"github.com/mrl-athomelab/website/application/logger"
	"github.com/mrl-athomelab/website/application/secure"
)

func (s *Server) adminMembersGetHandler(ctx *gin.Context) {
	members := &database.Members{}
	s.render(ctx, http.StatusOK, "admin-members", gin.H{
		"members": members.All(),
	})
}

func (s *Server) adminMembersEditGetHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Redirect(http.StatusFound, "/admin/members/?err=bad_id")
		return
	}
	member := &database.Member{}
	member.ID = uint(id)
	exists := member.Get(database.ByID)
	if !exists {
		ctx.Redirect(http.StatusFound, "/admin/members/?err=not_found")
		return
	}
	s.render(ctx, http.StatusOK, "admin-member", gin.H{
		"member": member,
	})
}

func (s *Server) adminMembersRestDeleteHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		jsonresponse.Error(ctx, err)
		return
	}
	member := &database.Member{}
	member.ID = uint(id)
	exists := member.Get(database.ByID)
	if !exists {
		jsonresponse.Failed(ctx, gin.H{
			"message": "user not exists !",
		})
		return
	}
	member.Delete()
	jsonresponse.Success(ctx, gin.H{
		"message": "removed !",
	})
}

func (s *Server) adminMembersImagePostHandler(ctx *gin.Context) {

}

func (s *Server) adminMembersRestPostHandler(ctx *gin.Context) {
	var input struct {
		Firstname       string `json:"firstname"`
		Lastname        string `json:"lastname"`
		Biography       string `json:"biography"`
		Socialmedialink string `json:"socialmedialink"`
		Socialmediatype string `json:"socialmediatype"`
		Rule            string `json:"rule"`
	}
	err := ctx.BindJSON(&input)
	if err != nil {
		jsonresponse.Error(ctx, err)
		return
	}
	member := &database.Member{
		FirstName:       input.Firstname,
		LastName:        input.Lastname,
		ShortBiography:  input.Biography,
		SocialMediaLink: input.Socialmedialink,
		SocialMediaType: input.Socialmediatype,
		Rule:            input.Rule,
	}
	exists := member.Get(database.ByUsernamePassword)
	if exists {
		jsonresponse.Failed(ctx, gin.H{
			"user":    member,
			"message": "user already exists !",
		})
		return
	}
	err = member.Save()
	if err != nil {
		jsonresponse.Error(ctx, err)
		return
	}
	jsonresponse.Success(ctx, gin.H{
		"user":    member,
		"message": "saved !",
	})
}

func (s *Server) adminAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, err := cookie.Get(ctx, "administrator_cookie", s.config.SecretKey)
		if err != nil {
			logger.Warn("Admin auth error %v", err)
			ctx.Abort()
			ctx.Redirect(http.StatusFound, "/admin/login?redirect="+ctx.Request.RequestURI)
			return
		}
		if c.IsExpired() {
			ctx.Abort()
			ctx.Redirect(http.StatusFound, "/admin/login?redirect="+ctx.Request.RequestURI)
			return
		}
		ctx.Next()
	}
}

func (s *Server) adminPanelGetHandler(ctx *gin.Context) {
	s.render(ctx, http.StatusOK, "admin-panel", nil)
}

func (s *Server) adminAdministratorsGetHandler(ctx *gin.Context) {
	administrators := &database.Administrators{}
	s.render(ctx, http.StatusOK, "admin-administrators", gin.H{
		"administrators": administrators.All(),
	})
}

func (s *Server) adminAdministratorsRestDeleteHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		jsonresponse.Error(ctx, err)
		return
	}
	administrator := &database.Administrator{}
	administrator.ID = uint(id)
	exists := administrator.Get(database.ByID)
	if !exists {
		jsonresponse.Failed(ctx, gin.H{
			"message": "user not exists !",
		})
		return
	}
	administrator.Delete()
	jsonresponse.Success(ctx, gin.H{
		"message": "removed !",
	})
}

func (s *Server) adminAdministratorsRestPostHandler(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := ctx.BindJSON(&input)
	if err != nil {
		jsonresponse.Error(ctx, err)
		return
	}
	administrator := &database.Administrator{
		Username: input.Username,
		Password: input.Password,
	}
	exists := administrator.Get(database.ByUsernamePassword)
	if exists {
		jsonresponse.Failed(ctx, gin.H{
			"user":    administrator,
			"message": "user already exists !",
		})
		return
	}
	err = administrator.Save()
	if err != nil {
		jsonresponse.Error(ctx, err)
		return
	}
	jsonresponse.Success(ctx, gin.H{
		"user":    administrator,
		"message": "saved !",
	})
}

func (s *Server) adminLoginAnyHandler(ctx *gin.Context) {
	method := ctx.Request.Method

	if method == "GET" {
		s.render(ctx, http.StatusOK, "admin-login", gin.H{
			"message": "",
		})
		return
	}

	if method == "POST" {
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		csrfToken := ctx.PostForm("csrf_token")

		_, valid := secure.ExtractToken(csrfToken, s.config.SecretKey)
		if !valid {
			s.render(ctx, http.StatusFound, "admin-login", gin.H{
				"message": "Bad csrfToken",
			})
			return
		}

		administrator := &database.Administrator{}
		administrator.Username = username
		administrator.Password = secure.MD5Hash(password, s.config.SecretKey)
		exists := administrator.Get(database.ByUsernamePassword)
		if !exists {
			s.render(ctx, http.StatusFound, "admin-login", gin.H{
				"message": "Invalid username or password",
			})
			return
		}
		cookie := &cookie.Cookie{}
		cookie.TTL = time.Hour
		cookie.Payload = administrator.UUID
		cookie.CreatedAt = time.Now()
		err := cookie.Set(ctx, "administrator_cookie", s.config.SecretKey)
		if err != nil {
			s.render(ctx, http.StatusFound, "admin-login", gin.H{
				"message": "Couldn't create cookie !",
			})
			return
		}
		ctx.Redirect(http.StatusFound, "/admin/panel")
		return
	}

	s.render(ctx, http.StatusFound, "admin-login", gin.H{
		"message": "Bad method !",
	})
}

func (s *Server) adminNewsGetHandler(ctx *gin.Context) {
	news := &database.NewsList{}
	s.render(ctx, http.StatusOK, "admin-news-list", gin.H{
		"news": news.All(),
	})
}

func (s *Server) adminNewsRestPutHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		jsonresponse.Error(ctx, err)
		return
	}
	news := &database.News{}
	news.ID = uint(id)
	exists := news.Get(database.ByID)
	if !exists {
		jsonresponse.Error(ctx, err)
		return
	}
	var input struct {
		Content string `json:"content"`
		Title   string `json:"title"`
	}
	err = ctx.BindJSON(&input)
	if err != nil {
		jsonresponse.Error(ctx, err)
		return
	}
	news.Title = input.Title
	news.Content = input.Content
	err = news.Save()
	if err != nil {
		jsonresponse.Error(ctx, err)
		return
	}

	jsonresponse.Success(ctx, gin.H{
		"message": "saved !",
		"news":    news,
	})
}

func (s *Server) adminNewsEditGetHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Redirect(http.StatusFound, "/admin/members/?err=bad_id")
		return
	}
	news := &database.News{}
	news.ID = uint(id)
	exists := news.Get(database.ByID)
	if !exists {
		ctx.Redirect(http.StatusFound, "/admin/members/?err=not_exists")
		return
	}
	s.render(ctx, 200, "admin-news-edit", gin.H{
		"news": news,
	})
}

func (s *Server) adminNewsRestPostHandler(ctx *gin.Context) {
	var input struct {
		Content string `json:"content"`
		Title   string `json:"title"`
	}
	err := ctx.BindJSON(&input)
	if err != nil {
		jsonresponse.Error(ctx, err)
		return
	}
	news := &database.News{}
	news.Title = input.Title
	news.Content = input.Content
	err = news.Save()
	if err != nil {
		jsonresponse.Error(ctx, err)
		return
	}

	jsonresponse.Success(ctx, gin.H{
		"news":    news,
		"message": "saved !",
	})
}

func (s *Server) adminNewsRestDeleteHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		jsonresponse.Error(ctx, err)
		return
	}
	news := &database.News{}
	news.ID = uint(id)
	exists := news.Get(database.ByID)
	if !exists {
		jsonresponse.Failed(ctx, gin.H{
			"message": "user not exists !",
		})
		return
	}
	news.Delete()
	jsonresponse.Success(ctx, gin.H{
		"message": "removed !",
	})
}
