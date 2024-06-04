package server

import (
	"github.com/xmonv/ssh-tunnel-manager/public"
	"github.com/gin-gonic/gin"
	"net/http"
	"io/fs"
	"fmt"
	"io"
)

var static fs.FS

func initStatic() {
	dist, _ := fs.Sub(public.Public, "dist")
	static = dist
}

func Static(app *gin.Engine) {
	initStatic()
	indexFile, _ := static.Open("index.html")
	defer func() {
		indexFile.Close()
	}()
	index, _ := io.ReadAll(indexFile)
	indexHtml := string(index)

	folders := []string{"assets"}
	for i, folder := range folders {
		sub, err := fs.Sub(static, folder)
		if err == nil {
			app.StaticFS(fmt.Sprintf("/%s/", folders[i]), http.FS(sub))
		}
	}

	app.NoRoute(func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.Status(200)
		_, _ = c.Writer.WriteString(indexHtml)
		c.Writer.Flush()
		c.Writer.WriteHeaderNow()
	})
}