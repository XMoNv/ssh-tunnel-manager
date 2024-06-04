package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Static(app *gin.Engine) {
	app.StaticFile("/", "./public/")
	app.StaticFS("/assets/", http.Dir("./public/assets"))

	// folders := []string{"assets", "images", "streamer", "static"}
	// for i, folder := range folders {
	// 	sub, err := fs.Sub(static, folder)
	// 	if err != nil {
	// 		utils.Log.Fatalf("can't find folder: %s", folder)
	// 	}
	// 	r.StaticFS(fmt.Sprintf("/%s/", folders[i]), http.FS(sub))
	// }
}