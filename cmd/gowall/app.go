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

		// Default keys
		c.Set("oauthMessage", oauthMessage)
		c.Set("oauthMessageExist", exist)
		c.Set("ProjectName", config.ProjectName)
		c.Set("CopyrightYear", year)
		c.Set("CopyrightName", config.CompanyName)
		c.Set("CacheBreaker", "br34k-01")

		// Underscore each statements
		c.Set("eachCode", template.HTML("<% _.each(errors, function(err) { %>"))
		c.Set("eachPermissionCode", template.HTML("<% _.each(permissions, function(permission) { %>"))
		c.Set("eachGroupCode", template.HTML("<% _.each(groups, function(group) { %>"))

		// Variables
		c.Set("errorCode", template.HTML(`<%- err %>`))
		c.Set("nameFullCode", template.HTML("<%- name.full %>"))
		c.Set("nameCode", template.HTML("<%- name %>"))
		c.Set("usernameCode", template.HTML("<%= username %>"))
		c.Set("dataCode", template.HTML("<%- data %>"))
		c.Set("_idCode", template.HTML("<%= _id %>"))
		c.Set("phonenumberCode", template.HTML("<%- phone %>"))
		c.Set("pivotCode", template.HTML("<%- pivot %>"))
		c.Set("isActiveCode", template.HTML("<%= isActive %>"))
		c.Set("messageCode", template.HTML("<%= message %>"))
		c.Set("statusNameCode", template.HTML("<%- status.name %>"))
		c.Set("statusCreatedTimeCode", template.HTML("<%= status.userCreated.time %>"))
		c.Set("pagesCurrentCode", template.HTML("<%= pages.current %>"))
		c.Set("pagesTotalCode", template.HTML("<%= pages.total %>"))
		c.Set("itemsBeginCode", template.HTML("<%= items.begin %>"))
		c.Set("itemsEndCode", template.HTML("<%= items.end %>"))
		c.Set("itemsTotalCode", template.HTML("<%= items.total %>"))
		c.Set("createdNameCode,", template.HTML("<%= userCreated.name %>"))
		c.Set("createdTimeCode", template.HTML("<%= userCreated.time %>"))

		// If else statements
		c.Set("ifGroupsLengthZeroCode", template.HTML("<% if (groups.length == 0) { %>"))
		c.Set("ifPermissionPermittedCode", template.HTML("<% if (permission.permit) { %>"))
		c.Set("ifPermissionsLengthZeroCode", template.HTML("<% if (permissions.length == 0) { %>"))
		c.Set("ifNameCode", template.HTML("<% if (name) { %>"))
		c.Set("ifSuccessCode", template.HTML("<% if (success) { %>"))
		c.Set("ifIDUndefinedCode", template.HTML("<% if (id == undefined) { %>"))
		c.Set("ifNotSuccessCode", template.HTML("<% if (!success) { %>"))
		c.Set("ifNotSuccessANDIDNotUndefinedCode", template.HTML("<% if (!success && id != undefined) { %>"))
		c.Set("ifRolesANDRolesAccount", template.HTML("<% if (roles && roles.account) { %>"))
		c.Set("elseCode", template.HTML("<% } else { %>"))

		// Ends of functions and closures
		c.Set("endCode", template.HTML("<% }); %>"))
		c.Set("endCode2", template.HTML("<% } %>"))

		// Error variables
		c.Set("errorFirstCode", template.HTML("<%- errfor.first %>"))
		c.Set("errorMiddleCode", template.HTML("<%- errfor.middle %>"))
		c.Set("errorLastCode", template.HTML("<%- errfor.last %>"))
		c.Set("errorMapMiddleCode", template.HTML("<%- errfor['middle'] %>"))
		c.Set("errorMapLastCode", template.HTML("<%- errfor['last'] %>"))
		c.Set("errorNewAccountID", template.HTML("<%- errfor.newAccountId %>"))
		c.Set("errorNewAdminID", template.HTML("<%- errfor.newAdminId %>"))
		c.Set("errorNameCode", template.HTML("<%- errfor.name %>"))
		c.Set("errorPivotCode", template.HTML("<%- errfor.pivot %>"))
		c.Set("errorUsernameCode", template.HTML("<%- errfor.username %>"))
		c.Set("errorEmailCode", template.HTML("<%- errfor.email %>"))
		c.Set("errorMessageCode", template.HTML("<%- errfor.message %>"))
		c.Set("errorMembershipsCode", template.HTML("<%- errfor.memberships %>"))
		c.Set("errorSettingsCode", template.HTML("<%- errfor.settings %>"))
		c.Set("errorNewUsernameCode", template.HTML("<%- errfor.newUsername %>"))
		c.Set("errorNewPasswordCode", template.HTML("<%- errfor.newPassword %>"))
		c.Set("errorZipCode", template.HTML("<%- errfor.zip %>"))
		c.Set("errorPhoneCode", template.HTML("<%- errfor.phone %>"))
		c.Set("errorCompanyCode", template.HTML("<%- errfor.company %>"))
		c.Set("errorConfirmCode", template.HTML("<%- errfor.confirm %>"))
		c.Set("errorPasswordCode", template.HTML("<%- errfor.password %>"))

		// Error conditional expressions
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
