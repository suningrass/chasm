package components

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

type Tree map[string][]string

type CmdTreeChoose struct {
	//string:root,cmdwords:every cmd
	Suggestions Tree //提示单词库
	SingleWords []string

	//匹配相关
	CurrentIndex int      // 当前匹配的提示单词索引
	Matched      []string // 匹配的提示单词列表
	IsMatched    bool     // 是否有匹配的提示单词

	//内部变量
	keywords  map[string][]string //string=root,[]string-subcmd
	rootwords []string            //主命令root，suggestion.pre =""，永远都在推荐里面
	//help
	debug string
}

func CreateCmdTreeChoose(v *viper.Viper) *CmdTreeChoose {

	//suggestion
	cmdtree, err := ParseTree(v)
	if err != nil {
		panic(err)
	}

	//rootword
	rw := []string{}
	//keyword
	kw := make(map[string][]string, 0)

	for cmdname := range cmdtree {
		parts := strings.Split(cmdname, ".")
		if !slices.Contains(rw, parts[0]) {
			rw = append(rw, parts[0])
		}

		for _, word := range parts[1:] {
			if !slices.Contains(kw[parts[0]], word) {
				kw[parts[0]] = append(kw[parts[0]], word)
			}
		}

	}

	//singlewords
	sw := v.GetStringSlice("shell.CommonAlias")

	ctc := &CmdTreeChoose{
		Suggestions: cmdtree,
		SingleWords: sw,
		rootwords:   rw,
		keywords:    kw,

		CurrentIndex: -1,
		Matched:      []string{},
		IsMatched:    false,

		debug: "",
	}

	return ctc
}

// ParseTree 一次性解析 TOML 字符串，返回统一 Tree
// 规则：
//  1. 空值（"" 或仅空白）直接丢弃；
//  2. 段名本身作为 key，value 为该段直接子键的有序切片；
//  3. 叶子键 value 为按英文逗号拆分后的非空片段切片，**不含整串原文**。
func ParseTree(v *viper.Viper) (map[string][]string, error) {

	// 2. 取得需要解析的根段列表
	rootList := v.GetStringSlice("cmd.root")
	if len(rootList) == 0 {
		return nil, fmt.Errorf("缺少 [cmd] root 配置")
	}

	// 3. 建立白名单，方便 O(1) 判断
	white := make(map[string]struct{}, len(rootList))
	for _, r := range rootList {
		white[strings.TrimSpace(r)] = struct{}{}
	}

	// 4. 复用原算法数据结构
	sectionSet := make(map[string]map[string]struct{})
	leafSet := make(map[string][]string)

	for _, fullKey := range v.AllKeys() {
		// 只保留以白名单段开头的 key
		if !inWhite(fullKey, white) {
			continue
		}

		raw := strings.TrimSpace(v.GetString(fullKey))
		parts := strings.Split(fullKey, ".")

		// 4.1 记录段-子键关系（与原算法一致）
		for i := 1; i < len(parts); i++ {
			sec := strings.Join(parts[:i], ".")
			child := parts[i]
			if sectionSet[sec] == nil {
				sectionSet[sec] = make(map[string]struct{})
			}
			sectionSet[sec][child] = struct{}{}
		}

		// 4.2 叶子节点值处理
		if raw != "" {
			subParts := strings.FieldsFunc(raw, func(r rune) bool { return r == ',' })
			var clean []string
			for _, p := range subParts {
				if t := strings.TrimSpace(p); t != "" {
					clean = append(clean, t)
				}
			}
			if len(clean) == 0 {
				clean = []string{""}
			}
			leafSet[fullKey] = clean
		}
	}

	// 5. 合并到统一 Tree（与原算法一致）
	tree := make(Tree)
	for sec, set := range sectionSet {
		list := make([]string, 0, len(set))
		for k := range set {
			list = append(list, k)
		}
		sort.Strings(list)
		tree[sec] = list
	}
	for k, v := range leafSet {
		tree[k] = v
	}
	return tree, nil
}

// inWhite 判断 key 是否以白名单段开头
func inWhite(key string, white map[string]struct{}) bool {
	// 例如 key = "git.notes.add"，需要判断 "git" 是否在 white
	top := strings.SplitN(key, ".", 2)[0]
	_, ok := white[top]
	return ok
}

// SplitBeforePos 返回 rune 光标 pos 之前的内容按空白切分后的 []string
// 光标之后的内容 **绝不出现**
func splitBeforePos(line string, pos int) []string {
	// 计算光标对应的字节位置
	bytePos := 0
	count := 0
	for i := 0; count < pos && i < len(line); {
		_, size := utf8.DecodeRuneInString(line[i:])
		bytePos = i + size
		count++
		i += size
	}
	// 只取光标之前字节
	sub := line[:bytePos]

	// 在 sub 内按空白切分
	var out []string
	start := 0
	for i := 0; i < len(sub); {
		r, size := utf8.DecodeRuneInString(sub[i:])
		if unicode.IsSpace(r) {
			if start < i {
				out = append(out, sub[start:i])
			}
			start = i + size
		}
		i += size
	}
	if start < len(sub) {
		out = append(out, sub[start:])
	}
	return out
}

// 去重复
func unique[T comparable](s []T) []T {
	if len(s) == 0 {
		return nil // 避免空切片 → nil 语义更一致
	}

	seen := make(map[T]struct{}, len(s))
	out := make([]T, 0, len(s))

	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

func (c *CmdTreeChoose) getMathedToken(line string, pos int) []string {
	//空白行
	if len(line) == 0 {
		return []string{}
	}
	//划分words，可以处理中字符
	words := splitBeforePos(line, pos)

	//rootword 和对应的 index
	tokens := []string{}
	index := []int{}

	//查找所有rootword和他的index
	for i := len(words) - 1; i >= 0; i-- {
		if slices.Contains(c.rootwords, words[i]) {
			tokens = append(tokens, words[i])
			index = append(index, i)
		}
	}
	//查找每个rootword开始到后面的所有keyword
	outputs := []string{}

	for no, rootword := range tokens {
		keywords := []string{}
		for i := index[no]; i < len(words); i++ {
			if slices.Contains(c.keywords[rootword], words[i]) &&
				!slices.Contains(keywords, words[i]) {
				keywords = append(keywords, words[i])
			}
		}
		//
		if len(keywords) == 0 { //表示就一个rootword
			outputs = append(outputs, rootword)
		} else { //组合成 root.key.key的拉平结构
			output := rootword + "." + strings.Join(keywords, ".")
			outputs = append(outputs, output)
		}
	}

	//debug
	/* if len(outputs) != 0 {
		c.debug = fmt.Sprintf("match-token: %v", outputs)
	} */
	return outputs

	//return unique(outputs)
}

func (c *CmdTreeChoose) Matching(line string, currentword string, endpos int) {
	matchedtokens := c.getMathedToken(line, endpos)
	c.Matched = []string{}
	//先组成一个候选，有root的就要从候选里面match
	backlist := []string{}

	if len(matchedtokens) == 0 { //没有匹配到，就从rootword开始
		backlist = c.rootwords
	} else { //匹配到了直接查表,加入匹配项次
		for _, token := range matchedtokens {

			value, ok := c.Suggestions[token]
			if ok {
				backlist = append(backlist, value...)
			}

			// backlist = append(backlist, c.Suggestions[token]...)
			backlist = unique(backlist) //去重

		}
		//backlist = append(backlist, c.Suggestions[matchedtoken]...)
	}

	//表示currentword 存在就开始匹配-输入字符后开始匹配
	if currentword != "" {
		if len(matchedtokens) != 0 {
			backlist = append(backlist, c.rootwords...) //前面有rootword，还是随时可以新rootword
		}

		for _, suggestion := range backlist {
			if strings.HasPrefix(strings.ToLower(suggestion), strings.ToLower(currentword)) {
				c.Matched = append(c.Matched, suggestion)
			}
		}

		//必须要有输入才能匹配singleword，如然太多了，matched里面已经有的就不加入
		for _, suggestion := range c.SingleWords {
			if strings.HasPrefix(strings.ToLower(suggestion), strings.ToLower(currentword)) &&
				!slices.Contains(c.Matched, suggestion) {
				c.Matched = append(c.Matched, suggestion)
			}

		}

	} else { //在首尾或者空白-此处光标没有字符
		c.Matched = backlist
	}

	if len(c.Matched) > 0 {
		c.IsMatched = true
		c.CurrentIndex = 0 // 初始化为第一个匹配项
	} else {
		c.IsMatched = false
		c.CurrentIndex = -1
	}

	//sort.Strings(s.Matched)

}

func (c *CmdTreeChoose) Next() {
	suglen := len(c.Matched)
	// 匹配词选择和插入
	if c.IsMatched && suglen > 0 {
		c.CurrentIndex = (c.CurrentIndex + 1) % suglen
	}
}

func (c *CmdTreeChoose) Up() {
	if c.CurrentIndex > 0 {
		c.CurrentIndex -= 1
	}
}

func (c *CmdTreeChoose) Down() {
	if c.CurrentIndex < len(c.Matched)-1 {
		c.CurrentIndex += 1
	}
}

func (c *CmdTreeChoose) GiveWord() string {
	if c.IsMatched && len(c.Matched) > 0 {
		return c.Matched[c.CurrentIndex]
	}
	return ""
}

func (c *CmdTreeChoose) ResetMatched() {
	c.CurrentIndex = -1
	c.Matched = []string{}
}

func (c *CmdTreeChoose) Render() string {
	//分页设置
	height := 10
	var top, end int
	if c.IsMatched /* && len(f.Matched) > 0 */ {
		top = (c.CurrentIndex / height) * height // 每页起始
		/* if top > len(f.Matched)-height {
			top = len(f.Matched) - height
		} */
		if top < 0 {
			top = 0
		}
		end = top + height
		if end > len(c.Matched) {
			end = len(c.Matched)
		}
	}
	// 高亮显示选中的提示单词
	highlight := lipgloss.NewStyle().Background(lipgloss.Color("#ed858dff")).Bold(true)
	modehighlight := lipgloss.NewStyle().Foreground(lipgloss.Color("#0a7b17ff")).Bold(true)
	output := ""

	if c.IsMatched {
		for i := top; i < end; i++ {
			suggestion := c.Matched[i]

			if i == c.CurrentIndex {
				output += "\n" + highlight.Render(suggestion)
				continue
			}
			output += "\n" + suggestion
		}
	}

	output += fmt.Sprintf("\ntotal: %d, currentNO: %d\n", len(c.Matched), c.CurrentIndex)

	//debug
	//output += c.debug
	//end of debug

	return lipgloss.JoinHorizontal(lipgloss.Left, modehighlight.Render("CMDMODE: "),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#087815ff")).Bold(true).Render(output),
	)
}
