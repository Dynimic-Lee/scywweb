package models

import (
	"SCYWWeb/pool"
	"SCYWWeb/utils"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	fmtdate "github.com/metakeule/fmtdate"
	"gopkg.in/mgo.v2/bson"
)

type WebModel struct {
}

type ImageData struct {
	FileName   string
	FileSize   string
	ImgSize    string
	UploadDate string
}

const ImagePath = "C:\\GOAPP\\src\\SCYWWeb\\data\\image"

func NewWebModel() *WebModel {
	return &WebModel{}
}

func (mm *WebModel) GetImageDatas() ([]ImageData, error) {
	files, err := ioutil.ReadDir(ImagePath)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	imgDatas := make([]ImageData, len(files))
	for i, f := range files {
		if true == f.IsDir() {
			continue
		}

		ig := ImageData{
			FileName:   f.Name(),
			FileSize:   fmt.Sprintf("%v KB", (f.Size() / 1024)),
			UploadDate: fmtdate.Format("YYYY-MM-DD hh:mm:ss", f.ModTime()),
		}

		ImgSize, _ := utils.GetImageSize(filepath.Join(ImagePath, f.Name()))
		ig.ImgSize = ImgSize

		// log.Printf("%s||%s||%s||%s\n", ig.FileName, ig.FileSize, ig.SaveDate, ig.ImgSize)

		imgDatas[i] = ig
	}

	return imgDatas, nil
}

func (mm *WebModel) GetImageData(Filename string) (*ImageData, error) {
	f, err := os.Stat(ImagePath + "\\" + Filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ig := &ImageData{
		FileName:   f.Name(),
		FileSize:   fmt.Sprintf("%v KB", (f.Size() / 1024)),
		UploadDate: fmtdate.Format("YYYY-MM-DD hh:mm:ss", f.ModTime()),
	}

	ImgSize, _ := utils.GetImageSize(filepath.Join(ImagePath, f.Name()))
	ig.ImgSize = ImgSize

	// log.Printf("%s||%s||%s||%s\n", ig.FileName, ig.FileSize, ig.SaveDate, ig.ImgSize)

	return ig, nil
}

func (mm *WebModel) LoadShowHaed() (string, error) {
	strUpload, err := ioutil.ReadFile("views/html/show/showHeadTemplate.html")
	if err != nil {
		return "", err
	}

	return string(strUpload), nil
}

// image data들을 가져와서 show 페이지 body에 들어갈 html 템플릿을 스트링으로 전달.
func (mm *WebModel) LoadShowBody() (string, error) {
	ImageDatas, err := mm.GetImageDatas()
	if err != nil {
		return "", err
	}

	var htmlUL bytes.Buffer
	tempThum, err := template.ParseFiles("views/html/show/showTemplate.html")
	if err != nil {
		return "", err
	}

	err = tempThum.Execute(&htmlUL, ImageDatas)
	if err != nil {
		return "", err
	}

	return htmlUL.String(), nil
}

func (mm *WebModel) LoadShowScript() (string, error) {
	showJs, err := ioutil.ReadFile("views/html/show/showjs.html")
	if err != nil {
		return "", err
	}

	return string(showJs), nil
}

func (mm *WebModel) LoadUploadHeader() (string, error) {
	uploadHeader, err := ioutil.ReadFile("views/html/upload/uploadHeadTemplate.html")
	if err != nil {
		return "", err
	}

	return string(uploadHeader), nil
}

func (mm *WebModel) LoadUploadBody() (string, error) {
	uploadBody, err := ioutil.ReadFile("views/html/upload/uploadTemp.html")
	if err != nil {
		return "", err
	}

	return string(uploadBody), nil
}

func (mm *WebModel) SaveImageFile(imageData multipart.File, imageHeader *multipart.FileHeader) error {
	return utils.SaveFile(imageData, imageHeader)
}

// DB에서 FileName을 갖는 Document가 있는지 Find() 한다
func (mm *WebModel) FindImageData(FileName string) ([]ImageData, int, error) {
	s := pool.GetSession()
	defer s.Close()

	var results []ImageData
	err := s.DB("test").C("imageDatas").Find(bson.M{"FileName": FileName}).Sort("-timestamp").All(&results)
	if err != nil {
		return nil, 0, err
	}

	return results, len(results), nil
}

// DB에 데이터를 insert() 한다.
func (mm *WebModel) InsertDataToMDB(FileName string, FileSize string, ImgSize string) error {
	s := pool.GetSession()
	defer s.Close()

	NowTime := fmtdate.Format("YYYY-MM-DD hh:mm:ss", time.Now())
	err := s.DB("test").C("imageDatas").Insert(
		&ImageData{
			FileName:   FileName,
			FileSize:   FileSize,
			ImgSize:    ImgSize,
			UploadDate: NowTime,
		},
	)

	return err
}

func (mm *WebModel) DeleteImage(FileName string) error {

	err := os.Remove("./data/image/" + FileName)
	if err != nil {
		return err
	}

	s := pool.GetSession()
	defer s.Close()

	selecter := map[string]string{
		"filename": FileName,
	}

	err = s.DB("test").C("imageDatas").Remove(selecter)

	return err
}
