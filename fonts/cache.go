package fonts

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"

	"github.com/rockwell-uk/go-text/fonts/ttf"
)

type MyFontCache map[string]*truetype.Font

func (fc MyFontCache) Store(fd draw2d.FontData, font *truetype.Font) {
	fc[fd.Name] = font
}

func (fc MyFontCache) Load(fd draw2d.FontData) (*truetype.Font, error) {
	font, stored := fc[fd.Name]

	if !stored {
		var c string
		pc, _, ln, ok := runtime.Caller(1)
		if ok {
			details := runtime.FuncForPC(pc)
			c, ln = details.FileLine(pc)
			path, _ := os.Getwd()
			c = strings.TrimPrefix(c, path+"/")
		} else {
			c = "unknown"
		}

		path := fmt.Sprintf("[%s:%v]", c, ln)

		return nil, fmt.Errorf("font %s is not stored in font cache. %v", fd.Name, path)
	}

	return font, nil
}

func init() {
	fontCache := MyFontCache{}

	TTFs := map[string]([]byte){
		"regular": ttf.Univers,
		"bold":    ttf.UniversBold,
	}

	for fontName, TTF := range TTFs {
		font, err := truetype.Parse(TTF)
		if err != nil {
			panic(err)
		}
		fontCache.Store(draw2d.FontData{Name: fontName}, font)
	}

	draw2d.SetFontCache(fontCache)
}
