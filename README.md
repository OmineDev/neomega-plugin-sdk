# NeOmega Go 插件 SDK

用 Go 为 Omega Frame（OMG）编写插件。插件可以注册游戏和终端命令、与玩家交互、监听数据包，以及读写组件数据。

[编写第一个插件](#编写第一个插件) · [SDK 接口](#sdk-接口) · [可运行的插件示例](https://github.com/OmineDev/neomega-dynamic-components)

## 编写第一个插件

这个插件给 OMG 终端添加一条 `hello` 命令。先按示例仓库根目录 README 的说明跑通复读机，再在同一个示例仓库根目录新建 `hello/`。

### 1. 注册终端命令

将下面的代码保存为 `hello/main.go`。插件在 `init` 中注册工厂，每份配置都会得到一个新的组件实例。`Inject` 在宿主中注册终端命令。

```go
package hello

import (
	"github.com/OmineDev/neomega-plugin-sdk/frame"
	"github.com/OmineDev/neomega-plugin-sdk/frame/objectentry"
)

type component struct{ frame.BasicComponent }

func (c *component) Inject(host frame.ExtendOmega) {
	c.BasicComponent.Inject(host)
	host.AddBackendMenuEntry(&frame.BackendMenuEntry{
		MenuEntry: frame.MenuEntry{Triggers: []string{"hello"}, Usage: "打印问候语"},
		OnTrigger: func(args []string) {
			host.Out().Info.Println("Hello from NeOmega!")
		},
	})
}

func init() {
	objectentry.Register(func() frame.Component { return &component{} })
}
```

### 2. 添加配置

在同一目录新建 `hello/hello.json`，内容如下。顶层的 `版本` 属于配置文件，组件自定义选项放在 `配置` 中。

```json
{
  "名称": "hello.a",
  "描述": "在终端打印问候语",
  "是否禁用": false,
  "版本": "0.0.1",
  "来源": "DynamicLoader",
  "配置": {}
}
```

### 3. 编译并运行

回到**示例仓库根目录**，按运行 OMG 的平台构建。下面以 Linux amd64 为目标。在 Windows 开发电脑上使用 `.\build.cmd -package ./hello linux`。

```sh
./build.sh -package ./hello linux
```

1. 把生成的 `buildenv/linux-amd64/hello.a` 放到 OMG 存储目录的 `lang/DynamicLoader/`。
2. 把 `hello/hello.json` 放到同一存储目录的 `config/`。
3. 在 OMG 终端输入 `reload`。看到 `成功加载组件: [hello.a]` 后，等终端菜单就绪，输入 `hello`。预期输出为 `Hello from NeOmega!`。

## SDK 接口

| 包 | 用途 |
| --- | --- |
| [frame](frame/) | 组件生命周期、游戏控制、玩家交互、菜单、存储和区块接口。 |
| [frame/objectentry](frame/objectentry/) | 在 `init` 中调用 `Register` 注册组件工厂。 |
| [qq](qq/) | 消息回调使用的群成员和频道数据类型。 |

构建时选择与目标 OMG 配套的 SDK 和运行平台。示例构建脚本会自动使用 SDK 指定的 Go 版本。数据包接口使用 SDK 对应版本的 gophertunnel 类型，`LevelChunk` 原始载荷与缓存 blob 的格式转换目前仅部分支持。
