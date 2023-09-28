package deliverer

/*
	deliverer
	视频投放
*/

type Deliverer interface {
	// Delivery 视频投放
	Delivery(videoFile string, cover string, title string, desc string, custom ...interface{}) error
}
