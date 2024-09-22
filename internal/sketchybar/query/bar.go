package query

type Bar struct {
	Position      string   `json:"position"`
	Topmost       string   `json:"topmost"`
	Sticky        string   `json:"sticky"`
	Hidden        string   `json:"hidden"`
	Shadow        string   `json:"shadow"`
	FontSmoothing string   `json:"font_smoothing"`
	BlurRadius    int      `json:"blur_radius"`
	Margin        int      `json:"margin"`
	Drawing       string   `json:"drawing"`
	Color         string   `json:"color"`
	BorderColor   string   `json:"border_color"`
	BorderWidth   int      `json:"border_width"`
	Height        int      `json:"height"`
	CornerRadius  int      `json:"corner_radius"`
	PaddingLeft   int      `json:"padding_left"`
	PaddingRight  int      `json:"padding_right"`
	YOffset       int      `json:"y_offset"`
	Clip          float64  `json:"clip"`
	ImageValue    string   `json:"image.value"`
	ImageDrawing  string   `json:"image.drawing"`
	ImageScale    float64  `json:"image.scale"`
	Items         []string `json:"items"`
}
