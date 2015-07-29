package main

import (
  "fmt"
  "github.com/gin-gonic/gin"
  "flag"
  "os"
  "strings"
  "net/http"
  "os/exec"
  "strconv"
)

func show_usage() {
	fmt.Fprintf(os.Stderr,
		"Usage: %s [-port=<port>] \n",
		os.Args[0])
	fmt.Fprintf(os.Stderr,
		"Flags:\n")
	flag.PrintDefaults()

}

var (
  g_http_port = flag.String("port", "7525", "run http server on specific port")
)

func main() {

  flag.Usage = show_usage
  flag.Parse()

  fmt.Printf("start http client on %s\n",*g_http_port)
  r := gin.Default()
  r.Use(CORSMiddleware())
  r.GET("/build", buildCommand)

  r.Run(":" + *g_http_port)
}


func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(200)
            return
        }
        c.Next()
    }
}





func buildCommand(c *gin.Context) {
  	//fmt.Fprintf(w, "start build, %s\n", o s.Getenv("GOPATH"))
  	var (
		cmdOut []byte
		err    error
	)
	var buildResult struct {
			Code int `json:"code"`
    		File string `json:"file"`
			Line int `json:"line"`
			Type string `json:"type"`
			Detail string `json:"detail"`
		}
  	cmdName := "go"
	cmdArgs := []string{"build"}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).CombinedOutput(); err != nil {
		out := string(cmdOut)
		outTab := strings.Split(out,"\n")
		errorLine := outTab[1]
		errorTab := strings.SplitN(errorLine, ":",4)
		
		fmt.Printf("There was an error running go command: %s\n details: %s\n", err, out)
		buildResult.Code = 1
		buildResult.File = errorTab[0]
		buildResult.Line,_ = strconv.Atoi(errorTab[1])
		buildResult.Type = errorTab[2]
		buildResult.Detail = errorTab[3]
   		c.JSON(http.StatusInternalServerError,buildResult)
	} else {
		buildResult.Code = 0
  		c.JSON(http.StatusOK, buildResult)
	}
	
}
