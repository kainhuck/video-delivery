package clamber

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Ifeng struct {
	path string
}

func NewIfeng(path string) Clamber {
	return &Ifeng{path: path}
}

func (i Ifeng) Crawl(uri string) (articleFile string, imageFile string, err error) {
	resp, err := http.Get(uri)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", err
	}

	articleContent := make([]string, 0)

	doc.Find("#articleBox p").Each(func(i int, selection *goquery.Selection) {
		articleContent = append(articleContent, selection.Text())
	})

	img := ""
	doc.Find("div#articleBox>img").Each(func(i int, selection *goquery.Selection) {
		img, _ = selection.Attr("src")
	})

	title := ""
	doc.Find("h2[class^=index_title_]").Each(func(index int, selection *goquery.Selection) {
		title = selection.Text()
	})

	now := time.Now().Format("2006-01-02")
	title = now + "|" + title

	articleFile = filepath.Join(i.path, "article", title) + ".txt"

	// 保存文本
	article, err := os.Create(articleFile)
	if err != nil {
		return "", "", err
	}
	if _, err := article.WriteString(strings.Join(articleContent, "\n")); err != nil {
		return "", "", err
	}

	// 保存图片
	resp, err = http.Get(img)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	imageFile = filepath.Join(i.path, "image", title) + ".jpg"

	image, err := os.Create(imageFile)
	if err != nil {
		return "", "", err
	}
	_, err = io.Copy(image, resp.Body)
	if err != nil {
		return "", "", err
	}

	return articleFile, imageFile, nil
}
