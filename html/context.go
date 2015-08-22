package html

import "github.com/gopherjs/gopherjs/js"

type Style interface {
	isStyleMarker()
}

type Context struct {
	*js.Object

	FillStyle   Style `js:"fillStyle"`
	StrokeStyle Style `js:"strokeStyle"`

	LineWidth      int      `js:"lineWidth"`
	LineCap        LineCap  `js:"lineCap"`
	LineJoin       LineJoin `js:"lineJoin"`
	MiterLimit     int      `js:"miterLimit"`
	LineDashOffset int      `js:"lineDashOffset"`

	Font         string        `js:"font"`
	TextAlign    TextAlign     `js:"textAlign"`
	TextBaseline TextBaseline  `js:"textBaseline"`
	Direction    TextDirection `js:"direction"`

	ShadowBlur    int   `js:"shadowBlur"`
	ShadowColor   Color `js:"shadowColor"`
	ShadowOffsetX int   `js:"shadowOffsetX"`
	ShadowOffsetY int   `js:"shadowOffsetY"`

	GlobalAlpha              float64            `js:"globalAlpha"`
	GlobalCompositeOperation CompositeOperation `js:"globalCompositeOperation"`
}

// Really floats work here but half pixels are ugly

func (c *Context) FillRect(x, y, width, height int) {
	c.Call("fillRect", x, y, width, height)
}

func (c *Context) ClearRect(x, y, width, height int) {
	c.Call("clearRect", x, y, width, height)
}

func (c *Context) StrokeRect(x, y, width, height int) {
	c.Call("strokeRect", x, y, width, height)
}

func (c *Context) FillText(text string, x, y int) {
	c.Call("fillText", text, x, y)
}

func (c *Context) StrokeText(text string, x, y int) {
	c.Call("strokeText", text, x, y)
}

func (c *Context) MeasureText(text string) *TextMetrics {
	return &TextMetrics{
		Object: c.Call("measureText", text),
	}
}

func (c *Context) LineDash() []int {
	return c.Call("getLineDash").Interface().([]int)
}

func (c *Context) SetLineDash(vals []int) {
	c.Call("setLineDash", vals)
}

func (c *Context) CreateLinearGradient(x0, y0, x1, y1 int) *CanvasGradient {
	return &CanvasGradient{
		Object: c.Call("createLinearGradient", x0, y0, x1, y1),
	}
}

func (c *Context) CreateRadialGradient(x0, y0, r0, x1, y1, r1 int) *CanvasGradient {
	return &CanvasGradient{
		Object: c.Call("createRadialGradient", x0, y0, r0, x1, y1, r1),
	}
}

func (c *Context) CreatePattern(image ImageSource, repetition Repetition) *CanvasPattern {
	return &CanvasPattern{
		Object: c.Call("createPattern", image, repetition),
	}
}

func (c *Context) BeginPath() {
	c.Call("beginPath")
}

func (c *Context) ClosePath() {
	c.Call("closePath")
}

func (c *Context) MoveTo(x, y int) {
	c.Call("moveTo", x, y)
}

func (c *Context) LineTo(x, y int) {
	c.Call("lineTo", x, y)
}

func (c *Context) BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y int) {
	c.Call("bezierCurveTo", cp1x, cp1y, cp2x, cp2y, x, y)
}

func (c *Context) QuadraticCurveTo(cpx, cpy, x, y int) {
	c.Call("quadraticCurveTo", cpx, cpy, x, y)
}

func (c *Context) Arc(x, y, radius int, startAngle, endAngle float64, antiClockwise bool) {
	c.Call("arc", x, y, radius, startAngle, endAngle, antiClockwise)
}

func (c *Context) ArcTo(x1, y1, x2, y2, radius int) {
	c.Call("arcTo", x1, y1, x2, y2, radius)
}

func (c *Context) Rect(x, y, width, height int) {
	c.Call("rect", x, y, width, height)
}

func (c *Context) Fill() {
	c.Call("fill")
}

func (c *Context) Stroke() {
	c.Call("stroke")
}

func (c *Context) Clip() {
	c.Call("clip")
}

func (c *Context) IsPointInPath(x, y int) bool {
	return c.Call("isPointInPath", x, y).Bool()
}

func (c *Context) IsPointInStroke(x, y int) bool {
	return c.Call("isPointInStroke", x, y).Bool()
}

func (c *Context) Rotate(angle float64) {
	c.Call("rotate", angle)
}

func (c *Context) Scale(x, y float64) {
	c.Call("scale", x, y)
}

func (c *Context) Translate(x, y int) {
	c.Call("translate", x, y)
}

func (c *Context) DrawImage(image ImageSource, dx, dy int) {
	c.Call("drawImage", image, dx, dy)
}

func (c *Context) DrawImageScaled(image ImageSource, dx, dy, dw, dh int) {
	c.Call("drawImage", image, dx, dy, dw, dh)
}

func (c *Context) DrawImageSection(image ImageSource, sx, sy, sw, sh, dx, dy, dw, dh int) {
	c.Call("drawImage", image, sx, sy, sw, sh, dx, dy, dw, dh)
}

func (c *Context) SetImageSmoothingEnabled(b bool) {
	c.Set("mozImageSmoothingEnabled", b)
	c.Set("msImageSmoothingEnabled", b)
	c.Set("webkitImageSmoothingEnabled", b)
	c.Set("oImageSmoothingEnabled", b)
	c.Set("imageSmoothingEnabled", b)
}

func (c *Context) CreateImageData(width, height int) *ImageData {
	return &ImageData{
		Object: c.Call("createImageData", width, height),
	}
}

func (c *Context) GetImageData(sx, sy, sw, sh int) *ImageData {
	return &ImageData{
		Object: c.Call("getImageData", sx, sy, sw, sh),
	}
}

func (c *Context) PutImageData(imageData ImageData, dx, dy int) {
	c.Call("putImageData", dx, dy)
}

func (c *Context) Save() {
	c.Call("save")
}

func (c *Context) Restore() {
	c.Call("restore")
}

type ImageData struct {
	*js.Object
	Width  int    `js:"width"`
	Height int    `js:"height"`
	Data   []byte `js:"data"`
}

type CanvasPattern struct {
	*js.Object
}

func (*CanvasPattern) isStyleMarker() {}

type CanvasGradient struct {
	*js.Object
}

func (c *CanvasGradient) AddColorStop(offset float64, color Color) {
	c.Call("addColorStop", offset, color)
}
func (*CanvasGradient) isStyleMarker() {}

type TextMetrics struct {
	*js.Object
	Width int `js:"width"`
}

type LineCap string

const (
	Butt   LineCap = "butt"
	Round  LineCap = "round"
	Square LineCap = "square"
)

type LineJoin string

const (
	Bevel     LineJoin = "bevel"
	RoundJoin LineJoin = "round"
	Miter     LineJoin = "miter"
)

type TextAlign string

const (
	Left   TextAlign = "left"
	Right  TextAlign = "right"
	Center TextAlign = "center"
	Start  TextAlign = "start"
	End    TextAlign = "end"
)

type TextBaseline string

const (
	Top         TextBaseline = "top"
	Hanging     TextBaseline = "hanging"
	Middle      TextBaseline = "middle"
	Alphabetic  TextBaseline = "alphabetic"
	Ideographic TextBaseline = "ideographic"
	Bottom      TextBaseline = "bottom"
)

type TextDirection string

const (
	LTR     TextDirection = "ltr"
	RTL     TextDirection = "rtl"
	Inherit TextDirection = "inherit"
)

type Repetition string

const (
	Repeat   Repetition = "repeat"
	RepeatX  Repetition = "repeat-x"
	RepeatY  Repetition = "repeat-y"
	NoRepeat Repetition = "no-repeat"
)

type CompositeOperation string

const (
	SourceOver      CompositeOperation = "source-over"
	SourceIn        CompositeOperation = "source-in"
	SourceOut       CompositeOperation = "source-out"
	SourceAtop      CompositeOperation = "source-atop"
	DestinationOver CompositeOperation = "destination-over"
	DestinationIn   CompositeOperation = "destination-in"
	DestinationOut  CompositeOperation = "destination-out"
	DestinationAtop CompositeOperation = "destination-atop"
	Lighter         CompositeOperation = "lighter"
	Copy            CompositeOperation = "copy"
	Xor             CompositeOperation = "xor"
	Multiply        CompositeOperation = "multiply"
	Screen          CompositeOperation = "screen"
	Overlay         CompositeOperation = "overlay"
	Darken          CompositeOperation = "darken"
	Lighten         CompositeOperation = "lighten"
	ColorDodge      CompositeOperation = "color-dodge"
	ColorBurn       CompositeOperation = "color-burn"
	HardLight       CompositeOperation = "hard-light"
	SoftLight       CompositeOperation = "soft-light"
	Difference      CompositeOperation = "difference"
	Exclusion       CompositeOperation = "exclusion"
	Hue             CompositeOperation = "hue"
	Saturation      CompositeOperation = "saturation"
	ColorOperation  CompositeOperation = "color"
	Luminosity      CompositeOperation = "luminosity"
)
