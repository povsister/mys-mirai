package bot

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/povsister/mys-mirai/pkg/log"
	"github.com/povsister/mys-mirai/pkg/util/fs"
	"github.com/povsister/mys-mirai/resources"
)

// 优先读取目录下的 device.json .
// 然后读取 embed device.json 或者随机生成一个
func (b *Bot) loadDeviceJSON() {
	if fs.FileExists("device.json") {
		log.Info().Msg("使用当前目录下的 device.json")
		read, err := fs.ReadFile("device.json")
		if err != nil {
			log.Fatal().Err(err).Msg("读取当前目录下的 device.json 失败")
		}
		if err = client.SystemDeviceInfo.ReadJson(read); err != nil {
			log.Fatal().Err(err).Msg("解析当前目录下的 device.json 失败")
		}
		return
	}

	if embd, err := resources.FS.ReadFile("device.json"); err == nil {
		log.Info().Msg("使用编译时提供的 embed device.json")
		if err = client.SystemDeviceInfo.ReadJson(embd); err == nil {
			return
		}
		log.Fatal().Err(err).Msg("程序错误 无效的 embed device.json")
	}

	log.Info().Msg("编译时未提供 embed device.json 将随机生成设备信息")
	client.GenRandomDevice()
	device := client.SystemDeviceInfo.ToJson()
	if err := fs.WriteFile("device.json", device, 0644); err != nil {
		log.Warn().Err(err).Msg("无法向当前目录下写入生成的 device.json")
		log.Warn().Msg("下次启动时仍将随机生成设备信息")
		return
	}
	log.Info().Msg("已将生成的设备信息写入 device.json")
	return
}
