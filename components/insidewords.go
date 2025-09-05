package components

const configStr = `
[cmd]
root = ["go","goctl","brew","hashcmp"]

[hashcmp]
dump = "--byname,--byhash,--output"
sml = "--simdist,--output"
cpuhash = "--sha256,--md5,--fastfinger,--output,--force,--ff"
mgtvs = "--desfile"

# Go 命令全集 TOML 主命令固定为 go
[go]
build = "-o,-a,-i,-n,-p,-race,-msan,-asan,-v,-work,-x,-asmflags,-buildmode,-buildvcs,-compiler,-gccgoflags,-gcflags,-installsuffix,-ldflags,-linkshared,-mod,-modcacherw,-modfile,-pkgdir,-tags,-trimpath,-h"
run = "-a,-n,-p,-race,-msan,-asan,-v,-work,-x,-asmflags,-buildmode,-buildvcs,-compiler,-gccgoflags,-gcflags,-installsuffix,-ldflags,-linkshared,-mod,-modcacherw,-modfile,-pkgdir,-tags,-trimpath,-h"
test = "-c,-i,-json,-o,-bench,-benchmem,-benchtime,-count,-cover,-covermode,-coverpkg,-cpu,-failfast,-fuzz,-fuzzminimizetime,-fuzztime,-list,-parallel,-run,-short,-shuffle,-skip,-timeout,-v,-vet,-work,-x,-asmflags,-buildmode,-buildvcs,-compiler,-gccgoflags,-gcflags,-installsuffix,-ldflags,-linkshared,-mod,-modcacherw,-modfile,-pkgdir,-tags,-trimpath,-h"
clean = "-cache,-modcache,-fuzzcache,-testcache,-n,-r,-v,-x,-h"
fmt = "-n,-x,-mod,-modfile,-h"
vet = "-n,-x,-v,-mod,-modfile,-h"
install = "-a,-i,-n,-p,-race,-msan,-asan,-v,-work,-x,-asmflags,-buildmode,-buildvcs,-compiler,-gccgoflags,-gcflags,-installsuffix,-ldflags,-linkshared,-mod,-modcacherw,-modfile,-pkgdir,-tags,-trimpath,-h"
get = "-d,-f,-t,-u,-v,-x,-insecure,-h"
list = "-deps,-f,-json,-m,-mod,-modfile,-u,-v,-versions,-h"
doc = "-all,-c,-cmd,-u,-h"
version = "-m,-v,-h"
env = "-json,-u,-w,-h"
bug = "-h"
fix = "-fix,-force,-h"
generate = "-run,-n,-v,-x,-h"

# === mod 子命令 ===
[go.mod]
download = "-x,-json,-h"
edit = "-fmt,-go,-module,-print,-require,-replace,-exclude,-droprequire,-dropreplace,-dropexclude,-json,-h"
graph = "-h"
init = "-h"
tidy = "-e,-v,-x,-diff,-go=<version>,-compat=<version>,-h"
vendor = "-e,-v,-o,-h"
verify = "-e,-v,-h"

# === work 子命令 ===
[go.work]
edit = "-fmt,-go,-module,-print,-replace,-use,-dropreplace,-dropuse,-json,-h"
init = "-h"
use = "-r,-h"
sync = "-h"

# === tool 子命令 ===
[go.tool]
addr2line = "-h"
asm = "-h"
buildid = "-h"
cgo = "-h"
compile = "-h"
cover = "-h"
dist = "-h"
doc = "-h"
fix = "-h"
link = "-h"
nm = "-h"
objdump = "-h"
pack = "-h"
pprof = "-h"
test2json = "-h"
trace = "-h"
vet = "-h"

# goctl 命令全集 TOML 主命令固定为 goctl
[goctl]
docker = "-h,--go,--version,--home,--remote,--branch,--style"
kube = "-h,-n,--namespace,--image,--secret,--replicas,--port,--node-port,--home,--remote,--branch,--style"
upgrade = "-h"
env = "-h"
bug = "-h"
completion = "fish,bash,zsh,powershell"

[goctl.api]
new = "-h,--dir,--api,--force,--style,--home,--remote,--branch"
format = "-h,--dir,--api,--style,--home,--remote,--branch"
validate = "-h,--dir,--api,--style,--home,--remote,--branch"
doc = "-h,--dir,--api,--style,--home,--remote,--branch,-o"
go = "-h,--dir,--api,--force,--style,--home,--remote,--branch"
java = "-h,--dir,--api,--force,--style,--home,--remote,--branch"
kt = "-h,--dir,--api,--force,--style,--home,--remote,--branch"
dart = "-h,--dir,--api,--force,--style,--home,--remote,--branch"
ts = "-h,--dir,--api,--force,--style,--home,--remote,--branch"
plugin = "-h,--dir,--api,--force,--style,--home,--remote,--branch"

# ============== RPC ==============
[goctl.rpc]
new = "-h,--o,--home,--remote,--branch,--style,-v"
protoc = "-h,--o,--home,--remote,--branch,--style,-v"
template = "-h,--home,--remote,--branch,--style,-v"

# ============== MODEL ==============
[goctl.model]
mysql = "-h,--src,--dir,--table,--cache,--database,--style,--home,--remote,--branch,--strict"
pg = "-h,--src,--dir,--table,--cache,--database,--style,--home,--remote,--branch,--strict"
mongo = "-h,-c,--dir,-d,-e,--home,--remote,--branch,--style,-t"
ddl = "-h,--src,--dir,--table,--cache,--database,--style,--home,--remote,--branch,--strict"
datasource = "-h,--url,--table,--dir,--cache,--database,--style,--home,--remote,--branch,--strict"

# ============== 其他 ==============
[goctl.template]
init = "-h,--home"
clean = "-h,--home"
revert = "-h,--category,--name,--home"
update = "-h,--home"

[brew]
install = "-h,--formula,--cask,--env,--ignore-dependencies,--only-dependencies,--build-from-source,--force-bottle,--include-test,--HEAD,--fetch-HEAD,--keep-tmp,--build-bottle,--force,--verbose,--debug,--dry-run,--quiet"
uninstall = "-h,--formula,--cask,--force,--ignore-dependencies,--zap,--quiet"
reinstall = "-h,--formula,--cask,--build-from-source,--force-bottle,--include-test,--HEAD,--fetch-HEAD,--keep-tmp,--force,--verbose,--debug,--dry-run,--quiet"
update = "-h,--merge,--force,--quiet"
upgrade = "-h,--formula,--cask,--fetch-HEAD,--dry-run,--ignore-pinned,--greedy,--force,--verbose,--debug,--quiet"
list = "-h,--formula,--cask,--full-name,--versions,--multiple,--pinned,--1,-l,-r,-t"
search = "-h,--formula,--cask,--desc,--macports,--fink,--opensuse,--fedora,--debian,--archlinux,--ubuntu"
info = "-h,--formula,--cask,--analytics,--github,--json,--installed,--all"
config = "-h"
doctor = "-h,--list-checks"
cleanup = "-h,-n,--prune,--dry-run,-s"
edit = "-h,--formula,--cask"
cat = "-h,--formula,--cask"
deps = "-h,--formula,--cask,--tree,--all,--installed,--include-build,--include-test,--include-optional,--skip-recommended"
uses = "-h,--formula,--cask,--installed,--recursive,--include-build,--include-test,--include-optional,--skip-recommended"
home = "-h,--formula,--cask"
link = "-h,--formula,--overwrite,--dry-run"
unlink = "-h,--formula,--dry-run"
pin = "-h,--formula"
unpin = "-h,--formula"
tap = "-h,--full,--force-auto-update,--custom-remote,--no-git"
untap = "-h,--force"
create = "-h,--autotools,--cmake,--crystal,--go,--meson,--node,--perl,--python,--ruby,--rust"
extract = "-h,--force,--version"
log = "-h,--formula,--cask,--oneline,--graph,--max-count"
leaves = "-h,--installed-on-request,--installed-as-dependency"
outdated = "-h,--formula,--cask,--fetch-HEAD,--greedy,--json,--quiet,--verbose"
missing = "-h"


[brew.services]
list = "-h,--json"
start = "-h,--all,--file,--verbose"
stop = "-h,--all,--file,--verbose"
restart = "-h,--all,--file,--verbose"
run = "-h,--all,--file,--verbose"
cleanup = "-h,--verbose"

[shell]
CommonAlias = [
  "type",
  "clc",
  "copy",
  "cpi",
  "-replace",
  "move",
  "del",
  "erase",
  "ren",
  "rni",
  "pwd",
  "popd",
  "pushd",
  "cls",
  "clear",
  "kill",
  "spps",
  "sleep",
  "sort",
  "select",
  "where",
  "foreach",
  "group",
  "measure",
  "compare",
  "tee",
  "write",
  "write-host",
  "write-error",
  "write-warning",
  "write-verbose",
  "read-host",
  "start-job",
  "receive-job",
  "remove-job",
  "stop-job",
  "wait-job",
  "get-job",
  "invoke-command",
  "test-connection",
  "get-date",
  "set-date",
  "get-random",
  "get-alias",
  "set-alias",
  "export-alias",
  "import-alias",
  "get-variable",
  "set-variable",
  "clear-variable",
  "remove-variable",
  "new-variable",
  "get-psdrive",
  "new-psdrive",
  "remove-psdrive",
  "get-psprovider",
  "test-path",
  "resolve-path",
  "convert-path",
  "join-path",
  "split-path",
  "start-process",
  "stop-process",
  "start-service",
  "stop-service",
  "restart-service",
  "get-service",
  "get-wmiobject",
  "get-ciminstance",
  "invoke-wmimethod",
  "export-csv",
  "import-csv",
  "convertTo-csv",
  "convertFrom-csv",
  "export-clixml",
  "import-clixml",
  "format-table",
  "format-list",
  "format-wide",
  "format-custom",
  "find",
  "head",
  "tail",
  "less",
  "more",
  "uniq",
  "cut",
  "diff",
  "cmp",
  "file",
  "stat",
  "touch",
  "mkdir",
  "rmdir",
  "chmod",
  "chown",
  "tar",
  "gzip",
  "gunzip",
  "zip",
  "unzip",
  "curl",
  "wget",
  "scp",
  "rsync",
]

[batcommands]
windows = [
 'dir /s /b#递归列出所有文件及完整路径',
 'dir /od#按修改时间排序显示目录',
 'Remove-Item -Force#直接删除文件,不受权限影响',
 'tree /f#以树形结构显示文件和文件夹',	
 'echo %PATH%#显示当前 PATH 环境变量',
 'set#列出所有环境变量',
 'setx JAVA_HOME "C:\jdk" /m#永久设置系统级 JAVA_HOME',
 'ipconfig#查看本机 IP 配置',
 'ipconfig /flushdns#清空 DNS 缓存',
 'ipconfig /release#释放当前 DHCP IP',
 'ipconfig /renew#重新获取 DHCP IP',
 'hostname#显示计算机名',
 'ping 127.0.0.1#测试本地回环网络',
 'tracert google.com#跟踪到谷歌的路由路径',
 'netstat -ano#列出所有端口及占用 PID',
 'netstat -b#列出端口及对应可执行文件',
 'arp -a#查看 ARP 缓存表',
 'nslookup google.com#域名解析测试',
 'whoami#显示当前登录用户',
 'whoami /priv#显示当前用户权限',
 'systeminfo#查看系统详细配置',
 'tasklist#列出所有运行中的进程',
 'taskkill /pid 1234 /f#强制结束指定 PID 进程',
 'sfc /scannow#扫描并修复系统文件',
 'chkdsk C: /f /r#检测并修复磁盘错误',
 'diskpart#进入磁盘分区管理交互界面',
 'driverquery#列出已安装驱动',
 'powercfg /energy#生成能耗报告',
 'powercfg /batteryreport#生成电池健康报告',
 'regedit#打开注册表编辑器',
 'reg query "HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall"#查询已安装程序',
 'reg add HKCU\Environment /v PATH /t REG_EXPAND_SZ /d "%PATH%;C:\mybin" /f#追加目录到用户 PAT',
 'sc query#列出所有服务',
 'net start#列出正在运行的服务',
 'net stop#可配合服务名停止服务',
 'net use#查看/映射网络驱动器',
 'net user#列出本地用户',
 'net localgroup administrators#显示管理员组成员',
 'net share#列出共享资源',
 'netsh interface show interface#列出网络接口',
 'netsh advfirewall show allprofiles#查看防火墙配置',
 'netsh wlan show profiles#列出保存的 Wi-Fi 配置',
 'netsh wlan show profile name=SSID key=clear#查看指定 Wi-Fi 明文密码',
 'netsh winsock reset#重置网络套接字',
 'wmic os get Caption,CSDVersion /value#简短系统版本信息',
 'wmic cpu get Name,NumberOfCores,NumberOfLogicalProcessors#CPU 核心信息',
 'wmic memorychip get Capacity,Speed#内存条容量与频率',
 'wmic logicaldisk get Caption,FreeSpace,Size#磁盘剩余/总空间',
 'wmic process where name="chrome.exe" get ProcessId#获取 Chrome 进程 PID',
 'wmic product get name,version#已安装软件列表',
 'wmic startup get caption,command#开机启动项',
 'wmic bios get serialnumber#主板序列号',
 'wmic csproduct get uuid#主板 UUID',
 'wmic nic get Name,MACAddress#网卡名称与 MAC',
 'wmic computersystem get TotalPhysicalMemory#总物理内存',
 'dxdiag /t dx.txt#导出 DirectX 诊断报告',
 'compact /c /s /i#压缩当前目录及子目录',
 'robocopy C:\source C:\dest /mir#镜像同步目录',
 'xcopy /s /e /h /i /y#带属性复制目录',
 'copy /y file1 file2#强制覆盖复制文件',
 'move /y file1 file2#强制覆盖移动文件',
 'del /q /s *.tmp#静默删除所有 tmp 文件',
 'rmdir /s /q folder#强制删除目录',
 'mkdir /p path\to\dir#递归创建目录',
 'attrib +h file.txt#设置隐藏属性',
 'attrib -h file.txt#取消隐藏属性',
 'findstr /s /i /n "ERROR" *.log#在所有日志中搜索 ERROR 并显示行号',
 'more file.txt#分页查看文本',
 'type file.txt#一次性输出文本',
 'fc file1.txt file2.txt#比较两个文件差异',
 'where git#查找可执行文件路径',
 'assoc .txt#查看 .txt 关联程序',
 'ftype txtfile#查看 .txt 打开方式',
 'winget list#列出已安装应用',
 'winget upgrade#检查可升级应用',
 'msinfo32 /report msinfo.txt#导出系统信息报告',
 'cmd /c start .#在当前目录打开资源管理器',
 'powershell -NoProfile -ExecutionPolicy Bypass -Command "Get-Process"#无策略限制执行脚本', 
]

zsh = [
 'ls -la -lh -GF -R#长格式列出所有文件,人类可读大小',
 'ipconfig getifaddr en0#查看本机ip',
 'diff -u <(cut -f1 file1.tsv | sort -u) <(cut -f1 file2.tsv | sort -u)#比较两个文件第一列cut -f1:取第一列TAB 分隔',
 'cd -#回到上一次目录',
 'system_profiler SPPowerDataType | grep Max#查看电池信息最大容量',
 '''tr ' ' '\n' | cat#空格转换行''',
 'top -o -cpu#实时进程',
 'du -sh *#目录大小',
 'mkdir -p a/b/c#递归创建目录',
 'rmdir empty_dir#删除空目录',
 'rm -rf dir#强制删除目录',
 'cp file1 file2#复制文件',
 'cp -r dir1 dir2#复制目录',
 'mv src dst#移动/重命名',
 'touch file#创建空文件',
 'cat file#查看文件内容',
 'tail -f log#实时追踪日志',
 'grep pattern file#搜索文本',
 'grep -i pattern#忽略大小写',
 'grep -r pattern dir#递归搜索',
 'find . -name "*.log"#按名称查找',
 'find . -type f -mtime -1#找今天修改的文件',
 'which cmd#查找可执行路径',
 'type cmd#显示命令类型',
 'history#查看历史命令',
 '!123#执行历史第 123 条',
 '!!#重复上一条命令',
 '!$#用上条最后一个参数',
 '''date +%Y%m%d#显示当前日期时间,可以变更格式''',
 'cal#显示月历',
 'whoami#当前用户',
 'who#登录用户列表',
 'ps aux#查看进程',
 'kill 1234#终止进程',
 'killall chrome#按名终止',
 'tar -czf file.tgz dir#打包压缩',
 'tar -xzf file.tgz#解压',
 'gzip file#单文件压缩',
 'gunzip file.gz#解压',
 'zip -r file.zip dir#zip 压缩',
 'unzip file.zip#解压 zip',
 'chmod 755 script#改权限',
 'chown user:group file#改属主',
 'df -h#磁盘使用',
 'uptime#系统负载',
 'uname -a#内核信息',
 'lsb_release -a#发行版信息',
 'env#环境变量',
 'export VAR=value#导出变量',
 'source ~/.zshrc#重载配置',
 'echo $SHELL#查看默认 shell',
 'printenv PATH#打印某变量',
 'wc -l file#统计行数',
 'sort file#排序',
 'sort -u file#去重排序',
 'uniq file#去重',
 'uniq -c file#计数去重',
 'cut -d: -f1 /etc/passwd#切列',
 'paste file1 file2#合并列',
 '''tr 'a-z' 'A-Z'#大小写转换''',
 '''sed 's/old/new/g' file#替换文本''',
 '''awk '{print $1}' file#取第一列''',
 'man cmd#查看手册',
 'whatis cmd#一句话描述',
 'zsh --version#查看 zsh 版本',
 'sudo /System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister -kill -r -domain local -domain system -domain user#执行终端命令刷新登录项缓存',

]

general = [
 'gox -osarch="windows/amd64" -ldflags="-s -w" -output="dist/{{.OS}}_{{.Arch}}/myapp" ',
 'go install -ldflags="-s -w" -trimpath',
 '''chasm <<< "$(fzf -m)"#与fzf的正确用法,先fzf,后chasm''', 
 '''fzf -m <<< "$(chasm | tr ' ' '\n' | cat )" #与fzf的正确用法,先chasm,后fzf''', 
]

`
