package downloader

var count int32



type Downloader  interface{
	Download(chan interface{})
}

