package impl

import (
	"net/http"
	framework "tumor_server/framework"
	"tumor_server/message"
	"tumor_server/model"
	"tumor_server/service"

	"github.com/sirupsen/logrus"
)

func AllRouteMiddleware(c *framework.Context) {
	logrus.Print("AllRouteMiddleware")
	c.Next()
}

func V1RouteMiddleware(c *framework.Context) {
	logrus.Print("V1RouteMiddleware")
	c.Next()
}

func V1AuthMiddleware(c *framework.Context) {
	logrus.Print("V1AuthMiddleware")
	token := c.GetAuthorization()
	tokenInfo, err := service.TokenValidate(token, c.Req.Host)
	if err != nil {
		reponse := model.HttpResponse{
			Code: http.StatusUnauthorized,
			Msg:  message.AuthorityError,
		}
		c.AbortWithJSON(http.StatusUnauthorized, reponse)
		return
	}
	c.SetExtra(tokenInfo)
	c.Next()
}

func V1AdminMiddleware(c *framework.Context) {
	logrus.Print("V1AdminMiddleware")
	token := c.GetAuthorization()
	tokenInfo, err := service.TokenValidate(token, c.Req.Host)
	if err != nil {
		reponse := model.HttpResponse{
			Code: http.StatusUnauthorized,
			Msg:  message.AuthorityError,
		}
		c.AbortWithJSON(http.StatusUnauthorized, reponse)
		return
	}
	c.SetExtra(tokenInfo)
	c.Next()
}

func V1SuperAdminMiddleware(c *framework.Context) {
	logrus.Print("V1SuperAdminMiddleware")
	token := c.GetAuthorization()
	tokenInfo, err := service.TokenValidate(token, c.Req.Host)
	if err != nil {
		reponse := model.HttpResponse{
			Code: http.StatusUnauthorized,
			Msg:  message.AuthorityError,
		}
		c.AbortWithJSON(http.StatusUnauthorized, reponse)
		return
	}
	c.SetExtra(tokenInfo)
	c.Next()
}

func V1DipperMiddleware(c *framework.Context) {
	logrus.Print("V1AuthMiddleware")
	token := c.GetAuthorization()
	tokenInfo, err := service.TokenValidate(token, c.Req.Host)
	if err != nil {
		reponse := model.HttpResponse{
			Code: http.StatusUnauthorized,
			Msg:  message.AuthorityError,
		}
		c.AbortWithJSON(http.StatusUnauthorized, reponse)
		return
	}
	c.SetExtra(tokenInfo)
	c.Next()
}
