package controllers

import (
	"SCYWWeb/models"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebController struct {
	mm *models.WebModel
}

func NewWebController() *WebController {
	return &WebController{
		mm: models.NewWebModel(),
	}
}

// Index 화면 출력
func (wc *WebController) Show(c *gin.Context) {
	showHead, err := wc.mm.LoadShowHaed()
	if err != nil {
		log.Printf("[WebController] Show - LoadShowHaed Fail [%s]", err.Error())
		c.HTML(http.StatusOK, "template.html", gin.H{
			"addhead":     "",
			"bodycontent": "",
			"BodyScript":  "",
			"uri":         c.Request.RequestURI,
		})
		return
	}

	showBody, err := wc.mm.LoadShowBody()
	if err != nil {
		log.Printf("[WebController] Show - LoadShowBody Fail [%s]", err.Error())
		c.HTML(http.StatusOK, "template.html", gin.H{
			"addhead":     template.HTML(showHead),
			"bodycontent": "",
			"BodyScript":  "",
			"uri":         c.Request.RequestURI,
		})
		return
	}

	showScript, err := wc.mm.LoadShowScript()
	if err != nil {
		log.Printf("[WebController] Show - LoadShowScript Fail [%s]", err.Error())
		c.HTML(http.StatusOK, "template.html", gin.H{
			"addhead":     template.HTML(showHead),
			"bodycontent": template.HTML(showBody),
			"BodyScript":  "",
			"uri":         c.Request.RequestURI,
		})
		return
	}

	c.HTML(http.StatusOK, "template.html", gin.H{
		"addhead":     template.HTML(showHead),
		"bodycontent": template.HTML(showBody),
		"BodyScript":  template.HTML(showScript),
		"uri":         c.Request.RequestURI,
	})
}

func (wc *WebController) RemoveImage(c *gin.Context) {
	Filename := c.PostForm("data")
	log.Println("RemoveImage data : ", Filename)

	err := wc.mm.DeleteImage(Filename)
	if err != nil {
		log.Printf("[WebController] RemoveImage Fail [%s]", err.Error())
		c.JSON(http.StatusForbidden, gin.H{
			"error":  err.Error(),
			"result": http.StatusForbidden,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

func (wc *WebController) Upload(c *gin.Context) {
	uploadHead, err := wc.mm.LoadUploadHeader()
	if err != nil {
		log.Printf("[WebController] Upload - LoadUploadHeader Fail [%s]", err.Error())
		c.HTML(http.StatusOK, "template.html", gin.H{
			"addhead":     "",
			"bodycontent": "",
			"uri":         c.Request.RequestURI,
		})
		return
	}

	uploadBody, err := wc.mm.LoadUploadBody()
	if err != nil {
		log.Printf("[WebController] Upload - LoadUploadBody Fail [%s]", err.Error())
		c.HTML(http.StatusOK, "template.html", gin.H{
			"addhead":     template.HTML(uploadHead),
			"bodycontent": "",
			"uri":         c.Request.RequestURI,
		})
		return
	}

	c.HTML(http.StatusOK, "template.html", gin.H{
		"addhead":     template.HTML(uploadHead),
		"bodycontent": template.HTML(uploadBody),
		"uri":         c.Request.RequestURI,
	})
}

func (wc *WebController) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("upload")
	if err != nil {
		log.Printf("UploadFile - FormFile [%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error": "file get fail",
		})
		return
	}

	// 동일한 파일명 체크
	_, hit, _ := wc.mm.FindImageData(header.Filename)
	if 0 < hit {
		log.Printf("UploadFile - Already File Names")
		c.JSON(http.StatusOK, gin.H{
			"error": "Already File Names",
		})
		return
	}

	// 파일 저장
	err = wc.mm.SaveImageFile(file, header)
	if err != nil {
		log.Printf("UploadFile - SaveFile [%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error": "file save fail",
		})
		return
	}

	// DB에 정보 insert 위해서 저장된 파일에서 파일정보 추출
	ig, err := wc.mm.GetImageData(header.Filename)
	if err != nil {
		log.Printf("UploadFile - GetImageData [%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error": "get file info fail",
		})
		return
	}

	// DB에 파일 정보 insert
	err = wc.mm.InsertDataToMDB(ig.FileName, ig.FileSize, ig.ImgSize)
	if err != nil {
		log.Printf("UploadFile - GetImageData [%s]", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"error": "get file info fail",
		})
		return
	}

	c.JSON(http.StatusOK, "")
}
