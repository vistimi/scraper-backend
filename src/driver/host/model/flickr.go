package host

// https://golangexample.com/pagser-a-simple-and-deserialize-html-page-to-struct-based-on-goquery-and-struct-tags-for-golang-crawler/
type SearchPhotPerPageData struct {
	Stat    string        `pagser:"rsp->attr(stat)"`
	Page    uint          `pagser:"photos->attr(page)"`
	Pages   uint          `pagser:"photos->attr(pages)"`
	PerPage uint          `pagser:"photos->attr(perpage)"`
	Total   uint          `pagser:"photos->attr(total)"`
	Photos  []PhotoFlickr `pagser:"photo"`
}
type PhotoFlickr struct {
	ID     string `pagser:"->attr(id)"`
	Secret string `pagser:"->attr(secret)"`
	Title  string `pagser:"->attr(title)"`
}

// https://golangexample.com/pagser-a-simple-and-deserialize-html-page-to-struct-based-on-goquery-and-struct-tags-for-golang-crawler/
type DownloadPhotoSingleData struct {
	Label  string `pagser:"->attr(label)"`
	Width  int    `pagser:"->attr(width)"`
	Height int    `pagser:"->attr(height)"`
	Source string `pagser:"->attr(source)"`
}

// https://golangexample.com/pagser-a-simple-and-deserialize-html-page-to-struct-based-on-goquery-and-struct-tags-for-golang-crawler/
type DownloadPhotoData struct {
	Stat   string                    `pagser:"rsp->attr(stat)"`
	Photos []DownloadPhotoSingleData `pagser:"size"`
}

// https://golangexample.com/pagser-a-simple-and-deserialize-html-page-to-struct-based-on-goquery-and-struct-tags-for-golang-crawler/
type InfoPhotoData struct {
	Stat           string `pagser:"rsp->attr(stat)"`
	ID             string `pagser:"photo->attr(id)"`
	Secret         string `pagser:"photo->attr(secret)"`
	OriginalSecret string `pagser:"photo->attr(originalsecret)"`
	OriginalFormat string `pagser:"photo->attr(originalformat)"`
	Title          string `pagser:"title"`
	Description    string `pagser:"description"`
	UserID         string `pagser:"owner->attr(nsid)"`
	UserName       string `pagser:"owner->attr(username)"`
	Tags           []Tag  `pagser:"tag"`
}

type Tag struct {
	Name string `pagser:"->text()"`
}