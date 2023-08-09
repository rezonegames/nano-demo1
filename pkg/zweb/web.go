package zweb

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"net/http"
)

// Response setting gin.JSON
func Response(c *gin.Context, data proto.Message) {
	c.Header("Content-Type", "application/x-protobuf")
	c.ProtoBuf(http.StatusOK, data)
}
