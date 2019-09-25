package main

import (
	// "github.com/gin-contrib/gzip"
	// "github.com/gin-contrib/sessions"

	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

const VERSION = "0.1"
const ISOSTRING = "2006-01-02T15:04:05.99Z"

var store sessions.CookieStore

var Router *gin.Engine

var year int

func init() {
	InitConfig()
	store = sessions.NewCookieStore([]byte("MFDQmJQ4TF"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 6,
		HttpOnly: true,
	})

}

func main() {
	gin.SetMode(gin.ReleaseMode)
	Router = gin.New()
	Router.Use(gzip.Gzip(gzip.DefaultCompression))

	Router.StaticFile("/favicon.ico", "public/favicon.ico")
	Router.Static("/public", "public")
	Router.Static("/vendor", "vendor")

	// templates
	Router.HTMLRender = initTemplates()

	Router.Use(gin.Logger())
	Router.Use(checkRecover)

	Router.Use(sessions.Sessions("session", store))

	//refresh year every minute
	go func() {
		for {
			year, _, _ = time.Now().Date()
			time.Sleep(time.Minute)
		}
	}()

	Router.SetFuncMap(template.FuncMap{
		"errForUsernameCode": noescape,
	})

	Router.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		oauthMessage, exist := session.Get("oauthMessage").(string)
		session.Delete("oauthMessage")
		session.Save()

		// body := &bytes.Buffer{}
		// io.Copy(body, c.Request.Body)
		// body, _ := c.Request.GetBody()
		// b, _ := ioutil.ReadAll(body)

		// t, err := txt.New("code").Parse(string(b))
		// if err != nil {
		// 	log.Fatalf("error while parsing file: %v\n", err)
		// }
		// t.ExecuteTemplate(c.Writer, "errForUsernameCode", "<%- errfor.username ? \"has-error\" : \"\" %>")
		// t.Execute(c.Writer, "<%- errfor.username ? \"has-error\" : \"\" %>")

		c.Set("oauthMessage", oauthMessage)
		c.Set("oauthMessageExist", exist)
		c.Set("ProjectName", config.ProjectName)
		c.Set("CopyrightYear", year)
		c.Set("CopyrightName", config.CompanyName)
		c.Set("CacheBreaker", "br34k-01")

		c.Set("eachCode", template.HTML("<% _.each(errors, function(err) { %>"))
		c.Set("errorCode", template.HTML(`<%- err %>`))

		c.Set("ifSuccessCode", template.HTML("<% if (success) { %>"))
		c.Set("ifIDUndefinedCode", template.HTML("<% if (id == undefined) { %>"))
		c.Set("ifNotSuccessCode", template.HTML("<% if (!success) { %>"))
		c.Set("ifNotSuccessANDIDNotUndefinedCode", template.HTML("<% if (!success && id != undefined) { %>"))

		c.Set("endCode", template.HTML("<% }); %>"))
		c.Set("endCode2", template.HTML("<% } %>"))

		c.Set("firstCode", template.HTML("<%- errfor.first %>"))
		c.Set("middleCode", template.HTML("<%- errfor.middle %>"))
		c.Set("mapLastCode", template.HTML("<%- errfor['last'] %>"))

		c.Set("newPasswordCode", template.HTML("<%- errfor.newPassword %>"))
		c.Set("zipCode", template.HTML("<%- errfor.zip %>"))
		c.Set("phoneCode", template.HTML("<%- errfor.phone %>"))
		c.Set("companyCode", template.HTML("<%- errfor.company %>"))
		c.Set("emailCode", template.HTML("<%- errfor.email %>"))
		c.Set("confirmCode", template.HTML("<%- errfor.confirm %>"))
		c.Set("usernameCode", template.HTML("<%- errfor.username %>"))
		c.Set("passwordCode", template.HTML("<%- errfor.password %>"))
		c.Set("errForUsernameCode", template.HTML(`<%- errfor.username ? "has-error" : "" %>`))
		c.Set("errForPasswordCode", template.HTML(`<%- errfor.password ? "has-error" : "" %>`))
		c.Next()
	})
	Router.Use(IsAuthenticated)
	bindRoutes(Router) // --> cmd/go-getting-started/routers.go

	Router.Run(":" + config.Port)

	// https
	// put path to cert instead of CONF.TLS.CERT
	// put path to key instead of CONF.TLS.KEY
	/*
		go func() {
				http.ListenAndServe(":80", http.HandlerFunc(redirectToHTTPS))
			}()
			errorHTTPS := router.RunTLS(":443", CONF.TLS.CERT, CONF.TLS.KEY)
			if errorHTTPS != nil {
				log.Fatal("HTTPS doesn't work:", errorHTTPS.Error())
			}
	*/
}

// func errForUsernameCode() strings.Title {

// }

func noescape(str string) template.HTML {
	return template.HTML(str)
}

// force redirect to https from http
// necessary only if you use https directly
// put your domain name instead of CONF.ORIGIN
func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	//http.Redirect(w, req, "https://" + CONF.ORIGIN + req.RequestURI, http.StatusMovedPermanently)
}
