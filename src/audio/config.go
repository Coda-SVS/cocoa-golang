package audio

import (
	"github.com/Kor-SVS/cocoa/src/config"
	"github.com/gen2brain/malgo"
)

var (
	audioConfig        *config.Config
	audioContextConfig *config.Config
	audioDeviceConfig  *config.Config
)

func configInit() {
	audioConfig = config.RootConfig.GetSub("audio")
	audioContextConfig = audioConfig.GetSub("context")
	audioDeviceConfig = audioConfig.GetSub("device")

	audioContextConfig.SetDefault("ThreadPriority", malgo.ThreadPriorityHighest)
	audioDeviceConfig.SetDefault("DeviceName", "")
}
