package controllers

import (
	"bytes"
	"fmt"
	"github.com/mingzhehao/goutils/filetool"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type UploadController struct {
	BaseController
}

type Sizer interface {
	Size() int64
}

const (
	LOCAL_FILE_DIR    = "static/uploads/file"
	MIN_FILE_SIZE     = 1       // bytes
	MAX_FILE_SIZE     = 5000000 // bytes
	IMAGE_TYPES       = "(jpg|gif|p?jpeg|(x-)?png)"
	ACCEPT_FILE_TYPES = IMAGE_TYPES
	EXPIRATION_TIME   = 300 // seconds
	THUMBNAIL_PARAM   = "=s80"
)

var (
	imageTypes      = regexp.MustCompile(IMAGE_TYPES)
	acceptFileTypes = regexp.MustCompile(ACCEPT_FILE_TYPES)
)

type FileInfo struct {
	Url          string `json:"url,omitempty"`
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Size         int64  `json:"size"`
	Error        string `json:"error,omitempty"`
	DeleteUrl    string `json:"deleteUrl,omitempty"`
	DeleteType   string `json:"deleteType,omitempty"`
}

func (fi *FileInfo) ValidateType() (valid bool) {
	if acceptFileTypes.MatchString(fi.Type) {
		return true
	}
	fi.Error = "Filetype not allowed"
	return false
}

func (fi *FileInfo) ValidateSize() (valid bool) {
	if fi.Size < MIN_FILE_SIZE {
		fi.Error = "File is too small"
	} else if fi.Size > MAX_FILE_SIZE {
		fi.Error = "File is too big"
	} else {
		return true
	}
	return false
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func escape(s string) string {
	return strings.Replace(url.QueryEscape(s), "+", "%20", -1)
}

func getFormValue(p *multipart.Part) string {
	var b bytes.Buffer
	io.CopyN(&b, p, int64(1<<20)) // Copy max: 1 MiB
	return b.String()
}

func (this *UploadController) Handle() {
	f, h, err := this.GetFile("file")
	t := time.Now()
	path := LOCAL_FILE_DIR+t.Format("2006-01-02")

	defer f.Close()
	if err != nil {
		fmt.Println("getfile err ", err)
		this.Data["json"] = "no file"
		this.ServeJSON()
		return
	} else {
		var Url string
		ext := filetool.Ext(h.Filename)
		fi := &FileInfo{
			Name: h.Filename,
			Type: ext,
		}
		if !fi.ValidateType() {
			this.Data["json"] = "invalid file type"
			this.ServeJSON()
			return
		}
		var fileSize int64
		if sizeInterface, ok := f.(Sizer); ok {
			fileSize = sizeInterface.Size()
			fmt.Println(fileSize)
		}
		fileExt := strings.TrimLeft(ext, ".")
		fileSaveName := fmt.Sprintf("%s_%d%s", fileExt, time.Now().Unix(), ext)
		imgPath := fmt.Sprintf("%s/%s", path, fileSaveName)

		filetool.InsureDir(path)

		this.SaveToFile("file", imgPath) // 保存位置在 static/upload,没有文件夹要先创建
		if err == nil {
			Url = "/" + imgPath
		}
		this.Data["json"] = Url
		this.ServeJSON()
		return
	}
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}