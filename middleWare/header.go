/**
*FileName: middleWare
*Create on 2018/11/20 下午2:49
*Create by mok
*/

package middleWare

import "github.com/gin-gonic/gin"

func CrossDomain(c * gin.Context){
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
	c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Header("Content-Type", "multipart/form-data")
}
