package components

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"syscall"

	"github.com/charmbracelet/lipgloss"
)

type FileChoose struct {
	filelist    map[string]os.DirEntry
	CurrentPath string

	CurrentIndex int      //当前匹配index
	Matched      []string // 已经匹配列表
	IfMatching   bool     // 是否有匹配的提示单词
	//选择列表
	chooselist map[string]bool
}

func CreateFileChoose(path string) *FileChoose {

	fc := &FileChoose{
		filelist:     loaddir(path),
		CurrentIndex: -1,
		IfMatching:   false,
		Matched:      []string{},
		CurrentPath:  path,
	}

	fc.chooselist = setchooselist(fc.filelist)

	return fc
}

func loaddir(path string) map[string]os.DirEntry {
	// 1. 长度 0 或含 NUL 直接失败
	if path == "" || strings.IndexByte(path, 0) >= 0 {
		return nil
	}
	// 2. 让 filepath.Clean 做一次规范化，若结果为空则非法
	cleanpath := filepath.Clean(path)

	fl := make(map[string]os.DirEntry)
	// 1. 把当前目录转成绝对路径
	absDir, err := filepath.Abs(cleanpath)
	if err != nil {
		return nil
	}

	//检查是不是文件夹
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}
	if !info.IsDir() {
		return nil // 不是目录，返回空列表即可
	}

	entries, err := os.ReadDir(absDir)
	if err != nil {
		// Windows 受保护目录会抛出 ERROR_ACCESS_DENIED
		// 5 is ERROR_ACCESS_DENIED, 3 is ERROR_PATH_NOT_FOUND on Windows
		if errno, ok := err.(syscall.Errno); ok {
			if errno == 5 || errno == 3 {
				return nil // 跳过
			}
		}
		return nil

	}

	for _, e := range entries {
		// 2. 拼接绝对路径
		fullPath := filepath.Join(absDir, e.Name())
		fl[fullPath] = e
	}
	return fl
}

func setchooselist(entries map[string]os.DirEntry) map[string]bool {
	cl := make(map[string]bool)
	if len(entries) > 0 {
		for n := range entries {
			cl[n] = false
		}
	}
	return cl
}

func (f *FileChoose) Up() {
	if f.CurrentIndex > 0 {
		f.CurrentIndex -= 1
	}
}

func (f *FileChoose) Down() {
	if f.CurrentIndex < len(f.Matched)-1 {
		f.CurrentIndex += 1
	}

}

// listDrives 返回所有可用盘符（Windows）或根目录（Unix）的 os.DirEntry 映射。
// key 为路径，value 为对应的 os.DirEntry。
func listWindowsDrives() map[string]os.DirEntry {
	out := make(map[string]os.DirEntry)
	// 跨平台根路径
	//rootPath := "/"

	// Windows：枚举 A-Z
	for c := 'A'; c <= 'Z'; c++ {
		p := string(c) + ":\\"
		if fi, err := os.Stat(p); err == nil {
			out[p] = dirEntry{name: p, fi: fi}
		}
	}
	return out

}

// dirEntry 轻量实现 os.DirEntry
type dirEntry struct {
	name string
	fi   os.FileInfo
}

func (d dirEntry) Name() string               { return d.name }
func (d dirEntry) IsDir() bool                { return d.fi.IsDir() }
func (d dirEntry) Type() os.FileMode          { return d.fi.Mode() & os.ModeType }
func (d dirEntry) Info() (os.FileInfo, error) { return d.fi, nil }

func (f *FileChoose) Left() {
	// 获取当前目录
	dir, _ := filepath.Abs(f.CurrentPath)
	// 得到上级目录
	parent := filepath.Dir(dir)

	if parent == dir { // 已到根
		if runtime.GOOS == "windows" {
			// 获取全部盘符并当作目录刷新
			f.filelist = listWindowsDrives()
			f.chooselist = setchooselist(f.filelist)
			f.CurrentPath = "\\" // 约定：盘符列表

			return
		}

	}

	f.filelist = loaddir(parent)
	f.chooselist = setchooselist(f.filelist)
	//更新currentpath
	f.CurrentPath = parent

}

func (f *FileChoose) Right() {
	if f.IfMatching && len(f.Matched) > 0 /* && f.CurrentIndex >= 0 */ {
		subdir := f.Matched[f.CurrentIndex]
		/* if f.filelist[subdir].IsDir() {
			f.CurrentPath = subdir
			//path, _ := filepath.Abs(subdir.Name())
			//f.filelist = loaddir(subdir)
			//f.chooselist = setchooselist(f.filelist)
		} */
		info, err := os.Stat(subdir)
		if err != nil {
			return
		}
		if !info.IsDir() {
			return // 不是目录，返回空列表即可
		} else {
			f.CurrentPath = subdir
		}

	}

	f.filelist = loaddir(f.CurrentPath)
	f.chooselist = setchooselist(f.filelist)
	//f.Match("")

}

func (f *FileChoose) Choose() {
	if f.IfMatching && len(f.Matched) > 0 {
		filepath := f.Matched[f.CurrentIndex]
		if f.chooselist[filepath] {
			f.chooselist[filepath] = false
		} else {
			f.chooselist[filepath] = true
		}
	}
}

func (f *FileChoose) Match(currentword string) {
	f.Matched = []string{}
	if currentword != "" {
		//f.Matched = []string{}
		for choosepath := range f.chooselist {
			if strings.Contains(strings.ToLower(choosepath), strings.ToLower(currentword)) {
				f.Matched = append(f.Matched, choosepath)
			}
		}
	} else {
		for name := range f.chooselist {
			f.Matched = append(f.Matched, name)
		}

	}
	if len(f.Matched) > 0 {
		f.IfMatching = true
		f.CurrentIndex = 0 // 初始化为第一个匹配项
	} else {
		f.IfMatching = false
		f.CurrentIndex = -1
	}
	sort.Strings(f.Matched)
}

func (f *FileChoose) GiveWord() string {
	output := ""
	if f.IfMatching && len(f.Matched) > 0 {
		for name, ischoose := range f.chooselist {
			if ischoose {
				output += name + " "
			}
		}
		return output
	}
	return ""
}

func (f *FileChoose) ResetMatched() {

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	f.CurrentPath = pwd
	for name := range f.chooselist {
		f.chooselist[name] = false
	}
	// f.CurrentIndex = -1
	// f.Matched = []string{}
	f.Match("")
}

func (f *FileChoose) Render() string {
	//分页设置
	height := 10
	var top, end int
	if f.IfMatching /* && len(f.Matched) > 0 */ {
		top = (f.CurrentIndex / height) * height // 每页起始
		/* if top > len(f.Matched)-height {
			top = len(f.Matched) - height
		} */
		if top < 0 {
			top = 0
		}
		end = top + height
		if end > len(f.Matched) {
			end = len(f.Matched)
		}
	}

	// 高亮显示选中的提示单词
	highlight := lipgloss.NewStyle().Background(lipgloss.Color("#ed858dff")).Bold(true)
	modehighlight := lipgloss.NewStyle().Foreground(lipgloss.Color("#087c16ff")).Bold(true)
	dirhighlight := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#06a290ff"))
	choosehighlight := lipgloss.NewStyle().
		Background(lipgloss.Color("#b29eefff")). // 16/256/TrueColor 均可
		Foreground(lipgloss.Color("#5ce50cd9"))  // 前景色随意
	output := lipgloss.NewStyle().Underline(true).Foreground(lipgloss.Color("#e60a1cff")).Bold(true).Render(f.CurrentPath) + "\n"

	if f.IfMatching /* && top >= 0 && end > 0 */ {

		for i := top; i < end; i++ {
			suggestion := f.Matched[i]

			if i == f.CurrentIndex {
				if f.chooselist[suggestion] {
					output += choosehighlight.Render("✓") + highlight.Render(suggestion) + "\n"
					//continue
				} else {
					output += highlight.Render(suggestion) + "\n"
				}
				continue
			}
			if f.chooselist[suggestion] && i != f.CurrentIndex {
				output += choosehighlight.Render("✓") + suggestion + "\n"
				continue
			}
			if f.filelist[suggestion].IsDir() {
				output += dirhighlight.Render(suggestion) + "\n"
				continue
			}
			output += suggestion + "\n"
		}
	}
	//debug
	/* if f.IfMatching && len(f.Matched) > 0 && f.CurrentIndex >= 0 {
		subdir := f.Matched[f.CurrentIndex]
		output += "subdir: " + subdir + "\n"
	} */

	//output += fmt.Sprintf("debug: top %d, end %d, matchlen %d, currentindex %d\n", top, end, len(f.Matched), f.CurrentIndex)
	//end debug
	output += fmt.Sprintf("total: %d, currentNO: %d\n", len(f.Matched), f.CurrentIndex)

	return lipgloss.JoinHorizontal(lipgloss.Left, modehighlight.Render("FILEMODE: "), output)

}
