package clamber

/*
	clamber
	从指定网址爬取文章和封面存储到data/article和data/image
*/

type Clamber interface {
	Crawl(uri string) (articleFile string, imageFile string, err error)
}
