package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

var err error

type Deb struct {
	MinVersion  string `json:"minVersion"`
	HeaderImage string `json:"headerImage"`
	TintColor   string `json:"tintColor"`
	Tabs        []*struct {
		Tabname string  `json:"tabname"`
		Views   []*View `json:"views,omitempty"`
		Class   string  `json:"class"`
	} `json:"tabs"`
	Class string `json:"class"`
}

type View struct {
	Class            string `json:"class"`
	UseBoldText      bool   `json:"useBoldText,omitempty"`
	UseBottomMargin  bool   `json:"useBottomMargin,omitempty"`
	UseRawFormat     bool   `json:"useRawFormat,omitempty"`
	OpenExternal     bool   `json:"openExternal,omitempty"`
	ItemCornerRadius bool   `json:"itemCornerRadius,omitempty"`
	ItemSize         string `json:"itemSize,omitempty"`
	TintColor        string `json:"tintColor,omitempty"`
	Text             string `json:"text,omitempty"`
	Action           string `json:"action,omitempty"`
	Title            string `json:"title,omitempty"`
	Markdown         string `json:"markdown,omitempty"`
	Orientation      string `json:"orientation,omitempty"`
	TextColor        string `json:"textColor,omitempty"`
	Margins          string `json:"margins,omitempty"`
	FontWeight       string `json:"fontWeight,omitempty"`
	FontSize         int    `json:"fontSize,omitempty"`
	Spacing          int    `json:"spacing,omitempty"`
	UsePadding       bool   `json:"usePadding,omitempty"`
	Screenshots      []*struct {
		AccessibilityText string `json:"accessibilityText,omitempty"`
		URL               string `json:"url,omitempty"`
	} `json:"screenshots,omitempty"`
	Views []*View `json:"views,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		panic("target url missing")
	}
	var req *http.Request
	if req, err = http.NewRequest("GET", os.Args[1], nil); err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "Sileo/2.0.0b9 CoreFoundation/1677.104 Darwin/19.6.0")
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var bs []byte
	if bs, err = ioutil.ReadAll(resp.Body); err != nil {
		panic(err)
	}
	// 吐槽：垃圾中文山寨源 连个json都写不对 还写个jb代码？
	var reg = regexp.MustCompile(`,(\s+)?\]`)
	bs = reg.ReplaceAll(bs, []byte("]"))
	var res Deb
	if err = json.Unmarshal(bs, &res); err != nil {
		panic(err)
	}
	var def Deb
	if err = json.Unmarshal([]byte(_default), &def); err != nil {
		panic(err)
	}
	for idx, tab := range res.Tabs {
		for _, lable := range []string{"插件介绍", "插件版本", "插件大小", "下载次数", "插件作者", "更新时间", "兼容系统", "浏览截图"} {
			if view := getViewByLabel(tab.Views, lable); view != nil {
				if def.Tabs[idx] == nil {
					panic("默认模板和拷贝对象的tab数量不一致，请检查数据格式。")
				}
				if set := getViewByLabel(def.Tabs[idx].Views, lable); set != nil {
					for _, s := range view.Screenshots {
						if s.URL == "https://apt.wxhbts.com/images/spng.png" ||
							s.URL == "https://apt.cydiaa.com/images/spng.png" {
							s.URL = "https://apt.otokaze.me/images/no_screenshots.png"
						}
					}
					set.Text = view.Text
					set.Markdown = view.Markdown
					set.UseRawFormat = view.UseRawFormat
					set.UseBoldText = view.UseBoldText
					set.OpenExternal = view.OpenExternal
					set.Spacing = view.Spacing
					set.TintColor = view.TintColor
					set.Screenshots = view.Screenshots
					set.ItemSize = view.ItemSize
					set.ItemCornerRadius = view.ItemCornerRadius
				}
			}
		}
	}
	if bs, err = json.Marshal(def); err != nil {
		panic(err)
	}
	fmt.Printf(string(bs))
	return
}

func getViewByLabel(views []*View, label string) *View {
	for idx, view := range views {
		if label == "插件介绍" &&
			(view.Title == "软件介绍" || view.Title == "插件介绍") &&
			views[idx+1] != nil {
			return views[idx+1]
		}
		if label == "浏览截图" &&
			view.Title == "浏览截图" && views[idx+1] != nil {
			return views[idx+1]
		}
		if view.Text == label && views[idx+1] != nil {
			return views[idx+1]
		}
		if view = getViewByLabel(view.Views, label); view != nil {
			return view
		}
	}
	return nil
}
