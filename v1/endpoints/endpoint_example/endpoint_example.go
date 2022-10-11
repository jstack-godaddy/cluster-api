package endpoint_example

import "github.com/gin-gonic/gin"

func Initialize(eg *gin.RouterGroup) {

	eg.GET("/helloworld", Helloworld)

}
