package components

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

type BatChoose struct {
	BatList map[string][]string //所有列表

	category        []string
	currentCategory int

	CurrentIndex int      //当前匹配index
	Matched      []string // 已经匹配列表
	IfMatching   bool     // 是否有匹配的提示单词
}

func CreateBatChoose(v *viper.Viper) *BatChoose {

	batcmd := v.GetStringMapStringSlice("batcommands")

	cy := []string{}

	for name := range batcmd {
		cy = append(cy, name)
	}
	sort.Strings(cy)

	bc := &BatChoose{
		BatList:         batcmd,
		category:        cy,
		currentCategory: 0, //默认第一个

		CurrentIndex: -1,
		Matched:      []string{},
		IfMatching:   false,
	}

	return bc
}

func (b *BatChoose) Up() {
	if b.CurrentIndex > 0 {
		b.CurrentIndex -= 1
	}

}

func (b *BatChoose) Down() {
	if b.CurrentIndex < len(b.Matched)-1 {
		b.CurrentIndex += 1
	}

}

func (b *BatChoose) Left() {
	if b.currentCategory > 0 {
		b.currentCategory -= 1
	}

}

func (b *BatChoose) Right() {
	if b.currentCategory < len(b.category)-1 {
		b.currentCategory += 1
	}

}

func (b *BatChoose) Match(currentword string) {
	b.Matched = []string{}

	if currentword != "" {
		for _, bat := range b.BatList[b.category[b.currentCategory]] {
			if strings.Contains(strings.ToLower(bat), strings.ToLower(currentword)) {
				b.Matched = append(b.Matched, bat)
			}
		}
	} else {
		b.Matched = b.BatList[b.category[b.currentCategory]]
	}

	if len(b.Matched) > 0 {
		b.IfMatching = true
		b.CurrentIndex = 0 // 初始化为第一个匹配项
	} else {
		b.IfMatching = false
		b.CurrentIndex = -1
	}
	sort.Strings(b.Matched)

}

func (b *BatChoose) GiveWord() string {
	output := ""
	if b.IfMatching && len(b.Matched) > 0 {
		output = b.Matched[b.CurrentIndex]
		output = strings.Split(output, "#")[0]
	}
	return output
}

func (b *BatChoose) ResetMatched() {
	//b.currentCategory = 0
	b.Match("")
}

func (b *BatChoose) Render() string {
	//分页设置
	height := 10
	var top, end int
	if b.IfMatching /* && len(f.Matched) > 0 */ {
		top = (b.CurrentIndex / height) * height // 每页起始
		/* if top > len(f.Matched)-height {
			top = len(f.Matched) - height
		} */
		if top < 0 {
			top = 0
		}
		end = top + height
		if end > len(b.Matched) {
			end = len(b.Matched)
		}
	}
	// 高亮显示选中的提示单词
	highlight := lipgloss.NewStyle().Background(lipgloss.Color("#ed858dff")).Bold(true)
	modehighlight := lipgloss.NewStyle().Foreground(lipgloss.Color("#087c16ff")).Bold(true)
	output := ""
	for i, category := range b.category {
		if i == b.currentCategory {
			output += lipgloss.NewStyle().Underline(true).Foreground(lipgloss.Color("#e60a1cff")).Bold(true).Render(category) + " "
			continue
		}
		output += category + " "
	}

	output += "\n"

	if b.IfMatching {
		for i := top; i < end; i++ {
			suggestion := b.Matched[i]

			if i == b.CurrentIndex {
				output += highlight.Render(suggestion) + "\n"
				continue
			}
			output += suggestion + "\n"

		}

	}
	/* //debug
	if b.currentCategory >= 0 {
		output += fmt.Sprintf("\n%s,%d\n", b.category[b.currentCategory], b.currentCategory)
	}
	//debug */
	output += fmt.Sprintf("total: %d, currentNO: %d\n", len(b.Matched), b.CurrentIndex)

	return lipgloss.JoinHorizontal(lipgloss.Left, modehighlight.Render("BATMODE: "),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#087c16ff")).Bold(true).Render(output),
	)

}
