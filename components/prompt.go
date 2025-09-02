package components

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

// 定义状态
type mode int

const (
	CmdMode mode = iota
	BatMode
	FileMode
	ArgListMode
)

// 定义模型
type Promptmodel struct {
	Input textinput.Model

	Suggetions *CmdTreeChoose
	BatCmds    *BatChoose
	FileSelect *FileChoose
	Argchoose  *ArgListChoose

	debug string

	Mode mode
}

// 模型初始化
func InitialModel() Promptmodel {
	km := textinput.DefaultKeyMap //空格后pos也会变，如果不加，行末空格pos不会马上更新
	//km.CharacterForward.SetKeys(" ")
	ti := textinput.New()
	ti.KeyMap = km

	ti.Prompt = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#980612ff")).Render(">>")
	ti.Placeholder = "Char Asm......"
	ti.Focus()
	ti.Width = 150

	v := viper.New()
	v.SetConfigType("toml")
	if err := v.ReadConfig(strings.NewReader(configStr)); err != nil {
		fmt.Println("读取 TOML 失败: %w", err)
		panic(err)
	}

	sg := CreateCmdTreeChoose(v)
	bc := CreateBatChoose(v)
	fs := CreateFileChoose("./")
	ac := CreateArgListMode()

	p := Promptmodel{
		Input:      ti,
		Suggetions: sg,
		BatCmds:    bc,
		FileSelect: fs,
		Argchoose:  ac,

		debug: "",
		Mode:  CmdMode,
	}

	return p
}

// tea初始化函数
func (p Promptmodel) Init() tea.Cmd {
	return textinput.Blink
}

// 更新函数
func (p Promptmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	//p.vp, cmd = p.vp.Update(msg)
	p.Input, cmd = p.Input.Update(msg)
	currentword, startpos, endpos := wordAt(p.Input.Value(), p.Input.Position())

	switch p.Mode {
	//cmd模式
	case CmdMode:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc": //"ctrl+c",

				return p, tea.Quit
			case "up":
				p.Suggetions.Up()
				if p.Suggetions.GiveWord() != "" {
					p.replaceCursorWord(currentword, startpos, endpos, p.Suggetions.GiveWord())
				}
			case "down":
				p.Suggetions.Down()
				if p.Suggetions.GiveWord() != "" {
					p.replaceCursorWord(currentword, startpos, endpos, p.Suggetions.GiveWord())
				}

			case "tab":

				if p.Suggetions.GiveWord() != "" {
					p.replaceCursorWord(currentword, startpos, endpos, p.Suggetions.GiveWord())
				}
				p.Suggetions.Next()
			case " ": //空格表示新的一个word要输入，因此重新匹配
				p.Suggetions.ResetMatched()
				p.Suggetions.Matching(p.Input.Value(), currentword, p.Input.Position())
				return p, cmd

			case "enter":
				//p.Suggetions.ResetMatched()
				//p.vlog.SetValue(p.Input.Value())
				return p, tea.Quit

			case "ctrl+p": //切换到bat模式
				p.Mode = BatMode
				p.BatCmds.ResetMatched()
				//p.Input.Blur()
				//p.vlog.SetValue(p.Input.Value())
				return p, cmd

			case "ctrl+o": //切换到file模式
				p.Mode = FileMode
				p.FileSelect.ResetMatched()
				//p.Input.Blur()
				//p.vlog.SetValue(p.Input.Value())
				return p, cmd

			case "ctrl+n": //切换到arglist模式
				p.Mode = ArgListMode
				p.Argchoose.SetInput()
				//p.Input.Blur()
				//p.vlog.SetValue(p.Input.Value())
				return p, cmd

			default:
				p.Suggetions.Matching(p.Input.Value(), currentword, p.Input.Position())
				//p.vlog.SetValue(p.Input.Value())
			}

		}
		return p, cmd

	//file模式
	case FileMode:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up":
				p.FileSelect.Up()

			case "down":
				p.FileSelect.Down()

			case "shift+left":

				p.FileSelect.Left()
				p.FileSelect.Match(currentword)

				//p.Input, cmd = p.Input.Update(msg)
				return p, cmd

			case "shift+right":

				p.FileSelect.Right()
				p.FileSelect.Match("")

				//p.Input, cmd = p.Input.Update(msg)
				return p, cmd

			/* case "ctrl+c": //选择
			p.FileSelect.Choose() */
			//return p, cmd

			case " ":
				p.FileSelect.ResetMatched()

			case "tab": //选择
				p.FileSelect.Choose()
				//p.FileSelect.ResetMatched()

			case "enter": //替换并退出到cmdmode
				if p.FileSelect.GiveWord() != "" {
					p.replaceCursorWord(currentword, startpos, endpos, p.FileSelect.GiveWord())
				}
				p.Mode = CmdMode
				//p.FileSelect.ResetMatched()
				//p.Input.Focus()
				return p, cmd

			default:
				p.FileSelect.Match(currentword)
				//p.vlog.SetValue(p.Input.Value())

			}
			//p.Input.Update(msg)
		}

	//arglist模式
	case ArgListMode:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up":
				p.Argchoose.Up()

			case "down":
				p.Argchoose.Down()

			case " ":
				p.Argchoose.ResetMatched()

			case "tab": //选择
				p.Argchoose.Choose()
				//p.FileSelect.ResetMatched()

			case "enter": //替换并退出到cmdmode
				if p.Argchoose.GiveWord() != "" {
					p.replaceCursorWord(currentword, startpos, endpos, p.Argchoose.GiveWord())
				}
				p.Mode = CmdMode
				return p, cmd

			default:
				p.Argchoose.Match(currentword)
				//p.vlog.SetValue(p.Input.Value())

			}
			//p.Input.Update(msg)
		}

	//bat模式
	case BatMode:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up":
				p.BatCmds.Up()
				//p.BatCmds.Match(currentword)
			case "down":
				p.BatCmds.Down()
				//p.BatCmds.Match(currentword)
			case "shift+left":
				p.BatCmds.Left()
				p.BatCmds.Match("")

			case "shift+right":
				p.BatCmds.Right()
				p.BatCmds.Match("")

			case " ": //空格表示新的一个word要输入，因此重新匹配
				p.BatCmds.ResetMatched()
				return p, cmd
			case "tab": //替换
				if p.BatCmds.GiveWord() != "" {
					p.replaceCursorWord(currentword, startpos, endpos, p.BatCmds.GiveWord())
				}

			case "enter": //退出到cmdmode
				p.Mode = CmdMode
				//p.Input.Focus()
				return p, cmd

			default:
				p.BatCmds.Match(currentword)
				//p.vlog.SetValue(p.Input.Value())

			}

		}
		//return p, cmd

	}
	//p.vlog.SetValue(p.Input.Value())
	return p, cmd
}

// 显示函数
func (p Promptmodel) View() string {
	s := ""
	//s += p.vlog + "\n"
	s += p.Input.View() + "\n"
	//debug
	/* currentword, startpos, endpos := wordAt(p.Input.Value(), p.Input.Position())
	s += fmt.Sprintf("currentword: %s, sp: %d, ep: %d, curentpos: %d\n", currentword, startpos, endpos, p.Input.Position())
	*/
	//end debug
	if p.Mode == CmdMode {
		s += p.Suggetions.Render()
	}
	if p.Mode == BatMode {
		s += p.BatCmds.Render()
	}
	if p.Mode == FileMode {
		s += p.FileSelect.Render()
	}
	if p.Mode == ArgListMode {
		s += p.Argchoose.Render()
	}
	s += "\n"
	//s += "debug: " + p.debug
	s += lipgloss.NewStyle().Foreground(lipgloss.Color("#077f15ff")).Bold(true).Render("ctrl+p: bat; ctrl+o: file; ctrl+n: arg; tab: choose") + "\n"
	return s
}

func (p *Promptmodel) replaceCursorWord(currentWord string, startpos, endpos int, suggestword string) {
	line := p.Input.Value()
	if currentWord != "" {
		newline, _, ne := replaceRange(line, startpos, endpos, suggestword)
		p.Input.SetValue(newline)
		p.Input.SetCursor(ne)
	} else { //表示空格处或者首尾
		newline, _, ne := replaceRange(line, p.Input.Position(), p.Input.Position(), suggestword)
		p.Input.SetValue(newline)
		p.Input.SetCursor(ne)
	}
}

func replaceRange(s string, startRune, endRune int, newSub string) (string, int, int) {
	runes := []rune(s)
	n := len(runes)
	if startRune < 0 || startRune > endRune || endRune > n {
		return s, -1, -1 // 越界直接原样返回
	}

	newRunes := []rune(newSub)
	newLen := len(newRunes)

	// 替换/插入
	replaced := make([]rune, 0, n-newRuneLen(endRune-startRune)+newLen)
	replaced = append(replaced, runes[:startRune]...)
	replaced = append(replaced, newRunes...)
	replaced = append(replaced, runes[endRune:]...)

	// 新子串在结果中的起止
	newStart := startRune
	newEnd := startRune + newLen
	return string(replaced), newStart, newEnd

}

// 辅助：计算被删除的 rune 长度
func newRuneLen(oldLen int) int {
	if oldLen < 0 {
		return 0
	}
	return oldLen
}

// 返回光标所在单词及区间 [start,end)。
// 以空格为分隔符，start/end 都是 rune 索引。
func wordAt(s string, pos int) (word string, start, end int) {
	rs := []rune(s)
	n := len(rs)
	if pos < 0 || pos > n {
		return "", -1, -1
	}

	// 左边界：从 pos 向左找第一个空格
	left := pos
	for left > 0 && !unicode.IsSpace(rs[left-1]) {
		left--
	}
	// 右边界：从 pos 向右找第一个空格
	right := pos
	for right < n && !unicode.IsSpace(rs[right]) {
		right++
	}
	if left == right {
		return "", -1, -1 // 光标在空格或两端
	}
	return string(rs[left:right]), left, right
}
