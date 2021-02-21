package main

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServeShroutenerForever(store Store, address string, port string) {

	ashtml := "index.tmpl"

	uuidFactory := NewFactory(DefaultGenerator, store)
	namedFactory := NewFactory(IdentityGenerator, store)

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.LoadHTMLFiles("templates/index.tmpl")

	r.GET("/", indexHandler(store, ashtml))
	r.GET("/u/:shrot/", redirectHandler(store))
	r.GET("/shroutened", listAllHandler(store))

	r.POST("/shrouten", shroutenHandler(store, ashtml, uuidFactory))                  //uuid Generator needs nothing but url to shorten
	r.POST("/shrouten/:named", shroutenHandler(store, ashtml, namedFactory, "named")) // identity Generator needs a string, here provided by the Uri end named. See shroutenHandler details for more

	r.POST("/pruge", prugeHandler(store))

	routerconfig := address + ":" + port
	r.Run(routerconfig)
}

// indexhandler need the store and index template
func indexHandler(s Store, tmpl string) func(*gin.Context) {
	return (func(c *gin.Context) {
		c.HTML(http.StatusOK, tmpl, gin.H{
			"TOTAL_SHROUTENAGE": int(s.Len()),
		})
	})
}

// indexhandler needs only store as it will just 302
func redirectHandler(s Store) func(*gin.Context) {
	return (func(c *gin.Context) {
		u := c.Param("shrot")
		full := s.Get(u)
		c.Redirect(http.StatusFound, full)
	})
}

// indexhandler needs only store , used ad debug listing of all shroutened urls
func listAllHandler(s Store) func(*gin.Context) {
	return (func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, s.All())
	})
}

// this one makes Generic Shroutening, delegating to chosen factory the key generation
// Just as Factory Generators may need optional parameters, the opt variadic permits early way of getting such needed parameters
// we will get as much as told QueryString params then send them as is to factory Generator
// This is far from perfect as it limits to strings where Generators and factory are made more flexible (see factory.go),
// and enforcing query string params is very limiting. Yet, it permitted quite much more agnosticism and cleaning than we first had.
// Doing better is left as an exercise for the reader
func shroutenHandler(s Store, tmpl string, f *Factory, opt ...string) func(*gin.Context) {
	return (func(c *gin.Context) {
		url := c.PostForm("url")
		var optionals []string
		if opt != nil {
			for _, value := range opt {
				option := c.Param(value)
				optionals = append(optionals, option)
			}
		}
		fopts := make([]interface{}, len(optionals), len(optionals))
		for i := range optionals {
			fopts[i] = optionals[i]
		}
		key, err := f.Gen(url, fopts...)
		if err != nil {
			c.HTML(http.StatusBadRequest, tmpl, gin.H{
				"FORM_RESULT": "Niet Valid Shroutenable"})
		} else {
			if err = s.Set(key, url); err != nil {
				c.HTML(http.StatusInternalServerError, tmpl, gin.H{"TOTAL_SHROUTENAGE": int(s.Len()), "FORM_RESULT": "Intrenol errur Shroutening:" + err.Error()})
			} else {
				c.HTML(http.StatusAccepted, tmpl, gin.H{"TOTAL_SHROUTENAGE": int(s.Len()), "FORM_RESULT": template.HTML("Shrouten suczfullness ! <pre><a target='_new' href='/u/" + key + "'>" + key + "</a></pre>")})
			}
		}
	})
}

func prugeHandler(s Store) func(*gin.Context) {
	return (func(c *gin.Context) {
		s.Clear()
		c.Redirect(http.StatusFound, "/")
	})
}
