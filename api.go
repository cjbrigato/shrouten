package main

import (
	"fmt"
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

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, ashtml, gin.H{
			"TOTAL_SHROUTENAGE": int(store.Len()),
		})
	})

	r.GET("/u/:shrot/", func(c *gin.Context) {
		u := c.Param("shrot")
		full := store.Get(u)
		c.Redirect(http.StatusFound, full)
	})
	r.GET("/shroutened", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, store.All())
	})

	r.POST("/shrouten", func(c *gin.Context) {
		url := c.PostForm("url")
		key, err := uuidFactory.Gen(url, "")
		if err != nil {
			c.HTML(400, ashtml, gin.H{
				"FORM_RESULT": "Niet Valid Shroutenable"})
		} else {
			if err = store.Set(key, url); err != nil {
				c.HTML(500, ashtml, gin.H{"TOTAL_SHROUTENAGE": int(store.Len()), "FORM_RESULT": "Intrenol errur Shroutening:" + err.Error()})
			} else {
				c.HTML(200, ashtml, gin.H{"TOTAL_SHROUTENAGE": int(store.Len()), "FORM_RESULT": template.HTML("Shrouten suczfullness ! <pre><a target='_new' href='/u/" + key + "'>" + key + "</a></pre>")})
			}
		}
	})
	r.POST("/shrouten/:named", func(c *gin.Context) {
		fmt.Println("named.")
		key := c.Param("named")
		url := c.PostForm("url")
		key, err := namedFactory.Gen(url, key)
		if err != nil {
			c.HTML(400, ashtml, gin.H{
				"FORM_RESULT": "Niet Valid Shroutenable"})
		} else {
			if err = store.Set(key, url); err != nil {
				c.HTML(500, ashtml, gin.H{"TOTAL_SHROUTENAGE": int(store.Len()), "FORM_RESULT": "Intrenol errur Shroutening:" + err.Error()})
			} else {
				c.HTML(200, ashtml, gin.H{"TOTAL_SHROUTENAGE": int(store.Len()), "FORM_RESULT": template.HTML("Shrouten suczfullness ! <pre><a target='_new' href='/u/" + key + "'>" + key + "</a></pre>")})
			}
		}
	})
	r.POST("/pruge", func(c *gin.Context) {
		store.Clear()
		c.Redirect(http.StatusFound, "/")
	})

	routerconfig := address + ":" + port
	r.Run(routerconfig)
}
