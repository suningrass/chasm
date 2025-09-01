package components

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type ArgListChoose struct {
	ArgList      []string //所有列表
	CurrentIndex int      //当前匹配index
	Matched      []string // 已经匹配列表
	IfMatching   bool     // 是否有匹配的提示单词
	//选择列表
	chooselist map[string]bool
}

func CreateArgListMode() *ArgListChoose {
	alc := &ArgListChoose{
		ArgList:      []string{},
		CurrentIndex: -1,
		Matched:      []string{},
		IfMatching:   false,
		chooselist:   make(map[string]bool),
	}

	return alc
}

func (a *ArgListChoose) SetInput() {
	//去除空白行
	a.ArgList = slices.DeleteFunc(a.ArgList, func(s string) bool {
		return strings.TrimSpace(s) == ""
	})

	for _, arg := range a.ArgList {
		a.chooselist[arg] = false
	}

	a.Match("")

}

func (a *ArgListChoose) Up() {
	if a.CurrentIndex > 0 {
		a.CurrentIndex -= 1
	}
}

func (a *ArgListChoose) Down() {
	if a.CurrentIndex < len(a.Matched)-1 {
		a.CurrentIndex += 1
	}
}

func (a *ArgListChoose) Choose() {
	if a.IfMatching && len(a.Matched) > 0 {
		arg := a.Matched[a.CurrentIndex]
		if a.chooselist[arg] {
			a.chooselist[arg] = false
		} else {
			a.chooselist[arg] = true
		}
	}
}

func (a *ArgListChoose) GiveWord() string {
	output := ""
	if a.IfMatching && len(a.Matched) > 0 {
		for name, ischoose := range a.chooselist {
			if ischoose {
				output += name + " "
			}
		}
		return output
	}
	return ""
}

func (a *ArgListChoose) ResetMatched() {
	a.CurrentIndex = -1
	a.Matched = []string{}
	for name := range a.chooselist {
		a.chooselist[name] = false
	}
}

func (a *ArgListChoose) Match(currentword string) {
	a.Matched = []string{}
	if currentword != "" {
		for _, arg := range a.ArgList {
			if strings.Contains(strings.ToLower(arg), strings.ToLower(currentword)) {
				a.Matched = append(a.Matched, arg)
			}
		}
	} else {
		a.Matched = a.ArgList
	}

	if len(a.Matched) > 0 {
		a.IfMatching = true
		a.CurrentIndex = 0 // 初始化为第一个匹配项
	} else {
		a.IfMatching = false
		a.CurrentIndex = -1
	}

}

func (a *ArgListChoose) Render() string {
	//分页设置
	height := 10
	var top, end int
	if a.IfMatching /* && len(f.Matched) > 0 */ {
		top = (a.CurrentIndex / height) * height // 每页起始
		/* if top > len(f.Matched)-height {
			top = len(f.Matched) - height
		} */
		if top < 0 {
			top = 0
		}
		end = top + height
		if end > len(a.Matched) {
			end = len(a.Matched)
		}
	}

	// 高亮显示选中的提示单词
	highlight := lipgloss.NewStyle().Background(lipgloss.Color("#ed858dff")).Bold(true)
	modehighlight := lipgloss.NewStyle().Foreground(lipgloss.Color("#087c16ff")).Bold(true)
	choosehighlight := lipgloss.NewStyle().
		Background(lipgloss.Color("#b29eefff")). // 16/256/TrueColor 均可
		Foreground(lipgloss.Color("#5ce50cd9"))  // 前景色随意
	//output := lipgloss.NewStyle().Underline(true).Bold(true).Render(f.CurrentPath) + "\n"
	output := ""
	if a.IfMatching {
		for i := top; i < end; i++ {
			suggestion := a.Matched[i]

			if i == a.CurrentIndex {
				if a.chooselist[suggestion] {
					output += choosehighlight.Render("✓") + highlight.Render(suggestion) + "\n"
					continue
				} else {
					output += highlight.Render(suggestion) + "\n"
					continue
				}

			}
			if a.chooselist[suggestion] && i != a.CurrentIndex {
				output += choosehighlight.Render("✓") + suggestion + "\n"
				continue
			}
			output += suggestion + "\n"
		}
	}
	output += fmt.Sprintf("total: %d, currentNO: %d\n", len(a.Matched), a.CurrentIndex)

	return lipgloss.JoinHorizontal(lipgloss.Left, modehighlight.Render("ARGMODE: "),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#087c16ff")).Bold(true).Render(output),
	)

}
