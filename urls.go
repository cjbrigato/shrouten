package main

type ShortenedUrls []ShortenedUrl

type ShortenedUrl struct {
	Short string `json:"short"`
	Uri   string `json:"uri"`
}

func (urls *ShortenedUrls) add(url ShortenedUrl) {
	*urls = append(*urls, url)
}

func (urls *ShortenedUrls) addInPlace(short string, uri string) {
	urls.add(ShortenedUrl{short, uri})
}
