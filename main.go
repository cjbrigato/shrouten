package main

import (
	"flag"
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
)

func main() {
	addressPtr := flag.String("address", "127.0.0.1", "the shroutening happens on this address")
	portPtr := flag.String("port", "8000", "bonding the shroutening to this tcp prot")
	dbPtr := flag.String("db", "deso.le.db", "the shroutening persiztenza medium")
	flag.Parse()
	db := NewDB(*dbPtr)
	api(db, *addressPtr, *portPtr)
}

func atLeastOne(n int) bool {
	return n > 0
}

func api(db *DB, address string, port string) {

	uuidFactory := NewFactory(DefaultGenerator, db)
	namedFactory := NewFactory(CustomNamed, db)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"atLeastOne": atLeastOne,
	})
	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"TOTAL_SHROUTENAGE": int(db.Len()),
		})
	})

	r.GET("/u/:shrot/", func(c *gin.Context) {
		u := c.Param("shrot")
		full := db.Get(u)
		c.Redirect(302, full)
	})
	r.GET("/shroutened", func(c *gin.Context) {
		c.IndentedJSON(200, db.All())
	})

	r.POST("/shrouten", func(c *gin.Context) {
		url := c.PostForm("url")
		key, err := uuidFactory.Gen(url, "")
		if err != nil {
			c.HTML(400, "index.html", gin.H{
				"FORM_RESULT": "Niet Valid Shroutenable"})
		} else {
			if err = db.Set(key, url); err != nil {
				c.HTML(500, "index.html", gin.H{"TOTAL_SHROUTENAGE": int(db.Len()), "FORM_RESULT": "Intrenol errur Shroutening:" + err.Error()})
			} else {
				c.HTML(200, "index.html", gin.H{"TOTAL_SHROUTENAGE": int(db.Len()), "FORM_RESULT": template.HTML("Shrouten suczfullness ! <pre><a target='_new' href='/u/" + key + "'>" + key + " </a></pre>")})
			}
		}
	})
	r.POST("/shrouten/:named", func(c *gin.Context) {
		fmt.Println("named.")
		key := c.Param("named")
		url := c.PostForm("url")
		key, err := namedFactory.Gen(url, key)
		if err != nil {
			c.HTML(400, "index.html", gin.H{
				"FORM_RESULT": "Niet Valid Shroutenable"})
		} else {
			if err = db.Set(key, url); err != nil {
				c.HTML(500, "index.html", gin.H{"TOTAL_SHROUTENAGE": int(db.Len()), "FORM_RESULT": "Intrenol errur Shroutening:" + err.Error()})
			} else {
				c.HTML(200, "index.html", gin.H{"TOTAL_SHROUTENAGE": int(db.Len()), "FORM_RESULT": template.HTML("Shrouten suczfullness ! <pre><a target='_new' href='/u/" + key + "'>" + key + " </a></pre>")})
			}
		}
	})
	r.POST("/pruge", func(c *gin.Context) {
		db.Clear()
		c.Redirect(302, "/")
	})
	routerconfig := address + ":" + port
	r.Run(routerconfig)
}
