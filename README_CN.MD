# chasm — Character Assembly CLI  
> 本人用的交互式命令行字符串/命令组合器

用纯 Go 编写，单文件即可跨平台运行（macOS / Windows 均已实测）。  
核心目标：**用最少的击键次数，拼出最复杂的命令**。

---

## ✨ 主要功能

1. **智能补全**  
   补全逻辑全部在 `components/insidewords.go` 中维护，也可外置为 `insidewords.toml`，**无需重编译**即可热更新。

2. **文件导航**  
   内置迷你文件浏览器，支持模糊搜索、多选、Tab 快速补全路径。

3. **宏命令仓库**  
   把常用或冗长的命令存成宏，随时检索、插入、执行。

4. **管道友好**  
   `chasm ... | sh` 一行到位，也支持与任何 CLI 组合。

5. **剪贴板自动复制**  
   生成结果 **实时写入系统剪贴板**，Ctrl+V 即可粘贴执行。

6. **零依赖跨平台**  
   纯 Go 标准库，单文件分发，macOS / Windows 开箱即用。

---

## 🚀 安装

```bash
# 1. 直接运行
go run main.go

# 2. 本地编译
go build -o chasm

# 3. 全局安装
go install

🎮 交互模式速查
| 快捷键          | 模式           | 说明           |
| ------------ | ------------ | ------------ |
| `Tab` / `↑↓` | **CMDMOD**   | 默认命令输入+补全    |
| `Ctrl+P`     | **BATMODE**  | 调用宏命令，可搜索    |
| `Ctrl+O`     | **FILEMODE** | 选文件/文件夹，可多选  |
| `Ctrl+N`     | **ARGMODE**  | 处理上游管道列表，可多选 |

通用操作

Shift+←/→：在 BATMODE / FILEMODE 中切换类别
Enter：确认并返回 CMDMOD
顶部即时搜索框：实时过滤候选
📌 典型工作流

场景：执行 seed download <filepath>
启动 chasm → 进入 CMDMOD。
键入 seed → Tab 补全。
光标到 <filepath> 时按 Ctrl+O → FILEMODE 选文件。
Tab 选中 → Enter 回填路径 → Enter 退出。
字符串已复制到剪贴板；可直接 Ctrl+V 粘贴，或直接运行：
chasm seed download /path/to/file | sh

⚙️ 自定义补全
| 方式     | 文件                          | 说明             |
| ------ | --------------------------- | -------------- |
| **内嵌** | `components/insidewords.go` | 改源码后重新编译       |
| **外置** | `insidewords.toml`          | 同目录下放置即可，无需重编译 |

🛠️ 技术栈

语言：Go 1.20+
依赖：仅标准库 + 平台剪贴板库
终端：ANSI / Windows 控制台均可