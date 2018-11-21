/**
*FileName: objetcs
*Create on 2018/11/20 上午1:07
*Create by mok
*/

package objetcs

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"dataServer/model/file"
)

//接收上传的文件,并且保存在model.o中
func PostFile(c *gin.Context){
	c.Request.ParseMultipartForm(5000000)
	username := c.PostForm("username")
	form,err := c.MultipartForm()
	if err != nil{
		c.JSON(500,gin.H{
			"message":"读取上传文件失败",
			"error":err.Error(),
		})
	}
	for path,fheaders := range form.File{
		for _,fheader := range fheaders{
			fileInfo,err := file.GetFileFromFeader(username,path,fheader)
			if err != nil{
				c.JSON(200,gin.H{
					"message":fmt.Sprintf("上传文件%s失败",fheader.Filename),
					"error":err.Error(),
				})
				continue
			}
			err = fileInfo.AddFile()
			if err != nil{
				c.JSON(200,gin.H{
					"message":fmt.Sprintf("上传文件%s失败",fheader.Filename),
					"error":err.Error(),
				})
			}else {
				c.JSON(200,gin.H{
					"message":fmt.Sprintf("上传文件%s成功",fheader.Filename),
				})
			}
		}
	}
}
