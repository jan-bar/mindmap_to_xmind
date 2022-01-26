package main

import (
	"archive/zip"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage: %s xmind mindmap [sheet name]\n", os.Args[0])
		return
	}

	xm, err := ReadXMind(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	var sheetName string
	if len(os.Args) > 3 {
		sheetName = os.Args[3]
	}
	err = SaveMindMap(xm, os.Args[2], sheetName)
	if err != nil {
		fmt.Println(err)
	}
}

type (
	MindMap struct {
		Nodes        []Nodes           `json:"nodes"`
		ReadOnly     bool              `json:"readOnly"` // false
		Toolbar      Toolbar           `json:"toolbar"`
		Remarks      map[string]string `json:"remarks,omitempty"` // id:remarks
		Version      string            `json:"version"`           // 2.0
		ResourceList string            `json:"resourceList"`      // ""
	}
	CustomStyle struct {
		FontFamily     []string `json:"fontFamily"`               // 字体
		FontSize       string   `json:"fontSize,omitempty"`       // 字体大小14px
		FontWeight     string   `json:"fontWeight,omitempty"`     // 加粗bold
		FontStyle      string   `json:"fontStyle,omitempty"`      // 斜体italic
		TextDecoration string   `json:"textDecoration,omitempty"` // 删除线 underline line-through
		Color          string   `json:"color,omitempty"`          // 前景色 #206153
		BorderColor    string   `json:"borderColor,omitempty"`    // 背景色 #4D94FF
	}
	Nodes struct {
		ID          string      `json:"id"`               // 唯一id
		IsRoot      bool        `json:"isroot,omitempty"` // true表示根节点
		Topic       string      `json:"topic"`            // 内容
		CustomStyle CustomStyle `json:"customStyle,omitempty"`
		Link        interface{} `json:"link,omitempty"`   // null,不知道干嘛的都写这个吧
		Expanded    bool        `json:"expanded"`         // true表示展开
		ParentId    interface{} `json:"parentid"`         // 父节点id
		Style       Style       `json:"style"`            // 应该是会员才能切换的画布:经典/简洁,我们这里保持{}
		Remark      string      `json:"remark,omitempty"` // 评价,与外层Remarks保持一致
	}
	Toolbar struct {
		LineType    string `json:"lineType"`    // default
		Strategy    string `json:"strategy"`    // 布局样式: logic_right,logic_left,logic_left_right
		Zoom        int    `json:"zoom"`        // 1
		Loading     bool   `json:"loading"`     // true
		BorderColor string `json:"borderColor"` // #4D94FF
	}
)

type Style struct{}

func (s Style) MarshalJSON() ([]byte, error) {
	return []byte("{}"), nil // 全部返回{}
}

func SaveMindMap(xm *XMindContent, dst, sheetName string) error {
	root := xm.Sheet[0].Topic // 默认只转换第一个标签
	if sheetName != "" {
		// 循环查找指定名称的标签
		for _, s := range xm.Sheet {
			if s.Title == sheetName {
				root = s.Topic
				break
			}
		}
	}

	mind := MindMap{
		Version:      "2.0",
		ResourceList: "",
		Remarks:      map[string]string{},
		ReadOnly:     false,
		Toolbar: Toolbar{
			LineType:    "default",
			Strategy:    "logic_right", // 逻辑图向右
			Zoom:        1,
			Loading:     true,
			BorderColor: "#1BD5D2",
		},
		Nodes: []Nodes{
			{ // 组装默认根节点
				ID:     "root",
				IsRoot: true,
				Topic:  root.Title,
				CustomStyle: CustomStyle{
					FontFamily: []string{"Microsoft YaHei", "STXihei"},
				},
				Expanded: true,
				Link:     json.RawMessage("null"), // root节点为null,子节点中没有
				ParentId: json.RawMessage("null"), // root节点为null,子节点为父节点id字符串
			},
		},
	}

	// 递归将xmind节点转换为有道云笔记的节点
	AppendNode(root.Children, "root", &mind)

	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer fw.Close()
	// 保存为json文件
	return json.NewEncoder(fw).Encode(mind)
}

func AppendNode(ch Children, parentId string, mind *MindMap) {
	for _, n := range ch.Topics.Topic {
		mind.Nodes = append(mind.Nodes, Nodes{
			ID:    n.ID, // 复用xmind中的唯一id
			Topic: n.Title,
			CustomStyle: CustomStyle{
				FontFamily: []string{"Microsoft YaHei", "STXihei"},
			},
			Expanded: true,
			ParentId: parentId,
			Remark:   n.Notes.Plain,
		})
		if n.Notes.Plain != "" {
			// xmind的批注保存成有道云笔记的备注
			mind.Remarks[n.ID] = n.Notes.Plain
		}
		AppendNode(n.Children, n.ID, mind)
	}
}

type (
	XMindContent struct {
		XMLName xml.Name `xml:"xmap-content"`
		// Text       string   `xml:",chardata"`
		Xmlns      string  `xml:"xmlns,attr"`
		Fo         string  `xml:"fo,attr"`
		SVG        string  `xml:"svg,attr"`
		Xhtml      string  `xml:"xhtml,attr"`
		Xlink      string  `xml:"xlink,attr"`
		ModifiedBy string  `xml:"modified-by,attr"`
		Timestamp  string  `xml:"timestamp,attr"`
		Version    string  `xml:"version,attr"`
		Sheet      []Sheet `xml:"sheet"`
	}
	Sheet struct {
		// Text       string    `xml:",chardata"`
		ID         string    `xml:"id,attr"`
		ModifiedBy string    `xml:"modified-by,attr"`
		Theme      string    `xml:"theme,attr"`
		Timestamp  string    `xml:"timestamp,attr"`
		Topic      RootTopic `xml:"topic"`
		Title      string    `xml:"title"`
	}
	RootTopic struct {
		// Text           string   `xml:",chardata"`
		ID             string   `xml:"id,attr"`
		ModifiedBy     string   `xml:"modified-by,attr"`
		StructureClass string   `xml:"structure-class,attr"`
		Timestamp      string   `xml:"timestamp,attr"`
		Title          string   `xml:"title"`
		Children       Children `xml:"children"`
	}
	Children struct {
		// Text   string `xml:",chardata"`
		Topics Topics `xml:"topics"`
	}
	Topics struct {
		// Text  string  `xml:",chardata"`
		Type  string  `xml:"type,attr"`
		Topic []Topic `xml:"topic"`
	}
	Topic struct {
		// Text       string `xml:",chardata"`
		ID         string `xml:"id,attr"`
		ModifiedBy string `xml:"modified-by,attr"`
		Timestamp  string `xml:"timestamp,attr"`
		Title      string `xml:"title"`
		Children   struct {
			// Text   string `xml:",chardata"`
			Topics Topics `xml:"topics"`
		} `xml:"children"`
		Notes struct {
			Plain string `xml:"plain"`
		} `xml:"notes"`
	}
)

func ReadXMind(p string) (*XMindContent, error) {
	zr, err := zip.OpenReader(p)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer zr.Close()

	for _, f := range zr.File {
		if f.Name == "content.xml" {
			fr, err := f.Open()
			if err != nil {
				return nil, err
			}

			var xm XMindContent
			err = xml.NewDecoder(fr).Decode(&xm)
			_ = fr.Close()
			if err != nil {
				return nil, err
			}
			return &xm, nil
		}
	}
	return nil, errors.New("can not find content.xml")
}
