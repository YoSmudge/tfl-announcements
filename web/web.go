package web

import(
  "fmt"
  log "github.com/Sirupsen/logrus"
  "github.com/gin-gonic/gin"
  "time"
  "mime"
  "github.com/samarudge/tfl-announcements/fetcher"
  "github.com/satori/go.uuid"
)

var currentStatus *fetcher.StatusUpdate

func logger(c *gin.Context) {
  t := time.Now()
  c.Next()
  latency := time.Since(t)
  status := c.Writer.Status()
  clientIP := c.ClientIP()
  method := c.Request.Method

  log.WithFields(log.Fields{
    "Duration": latency,
    "Status": status,
    "ClientIP": clientIP,
    "Method": method,
    "Path": c.Request.URL.String(),
  }).Info("request")
}

func StartServer(bindAddress string){
  gin.SetMode(gin.ReleaseMode)
  r := gin.New()
  r.Use(logger, gin.Recovery())

  r.GET("/", redirectHome)
  r.GET("/Asset", returnAsset)
  r.GET("/UpdateAudio", getAudio)

  r.GET("/Feed", func(c *gin.Context){
    feedHandler(c.Writer, c.Request)
  })

  log.WithFields(log.Fields{
    "bind": bindAddress,
  }).Info("Starting web interface")
  r.Run(bindAddress)
}

func redirectHome(c *gin.Context){
  c.Redirect(302, "/Asset?File=index.html")
}

func returnAsset(c *gin.Context){
  asset := c.Query("File")
  filePath := fmt.Sprintf("web/assets/%s", asset)

  mimeType := mime.TypeByExtension(filePath)
  data, err := Asset(filePath)
  if err != nil{
    c.String(404, fmt.Sprintf("Asset error: %s", err))
  }

  c.Data(200, mimeType, data)
}

func getAudio(c *gin.Context){
  updateIdStr := c.Query("update_id")
  updateId, err := uuid.FromString(updateIdStr)
  if err != nil || updateId != currentStatus.Id{
    c.String(404, "Invalid update ID")
    return
  }
  c.Header("Content-Length", fmt.Sprintf("%d", len(currentStatus.Audio)))
  c.Data(200, "audio/mpeg", currentStatus.Audio)
}
