package definedfmt

import (
	"fmt"
	"github.com/HouzuoGuo/LinuxManagementConsole/txtedit"
	"testing"
)

var samples = []struct {
	config txtedit.AnalyserConfig
	input  string
}{

	{Sysconfig, `## Path:        Hardware/Console
## Description: Text console settings (see also Hardware/Keyboard)
## Type:        string
## Default:     ""
## ServiceRestart: kbd
#
# Console settings.
# Note: The KBD_TTY setting from Hardware/Keyboard (sysconfig/keyboard)
# also applies for the settings here.
#
# Load this console font on bootup:
# (/usr/share/kbd/consolefonts/)
#
CONSOLE_FONT="lat9w-16.psfu"

## Type:        string
## Default:     ""
#
# Some fonts come without a unicode map.
# (.psfu fonts supposedly have it, others often not.)
# You can then specify the unicode mapping of your font
# explicitly. (/usr/share/kbd/unimaps/)
# Normally not needed.
#
CONSOLE_UNICODEMAP=""
`},

	{Sysctl, `# some "comment"
vm.swappiness=10

# more comment
vm.zone_reclaim_mode=1
kernel.sysrq = 0
kernel.randomize_va_space=2
kernel.panic = 10
`},

	{Systemd, `#  This file is part of systemd.
#
# Entries in this file show the compile time defaults.
# You can change settings by editing this file.
# Defaults can be restored by simply deleting this file.
#
# See systemd-system.conf(5) for details.

[Manager]
#LogTarget=journal-or-kmsg
#LogColor=yes
#CrashShell=no
#CrashChVT=1
#CPUAffinity=1 2
#JoinControllers=cpu,cpuacct net_cls,net_prio
#RuntimeWatchdogSec=0
#ShutdownWatchdogSec=10min
#SystemCallArchitectures=
#DefaultTimerAccuracySec=1min
#DefaultStandardOutput=journal
#DefaultTimeoutStopSec=90s
#DefaultRestartSec=100ms
#DefaultStartLimitInterval=10s
#DefaultStartLimitBurst=5
#DefaultEnvironment=
#  This file is part of systemd.
#
#  systemd is free software; you can redistribute it and/or modify it
#  under the terms of the GNU Lesser General Public License as published by
#  the Free Software Foundation; either version 2.1 of the License, or
#  (at your option) any later version.
#
# Entries in this file show the compile time defaults.
# You can change settings by editing this file.
# Defaults can be restored by simply deleting this file.
#
# See logind.conf(5) for details.

[Login]
#ReserveVT=6
#KillUserProcesses=no
#HibernateKeyIgnoreInhibited=no
#LidSwitchIgnoreInhibited=yes
#HoldoffTimeoutSec=30s
#IdleAction=ignore
#IdleActionSec=30min
#RuntimeDirectorySize=10%
#RemoveIPC=yes
`}}

func TestTextBreakdown(t *testing.T) {
	for _, sample := range samples {
		analyser := txtedit.NewAnalyser(sample.input, &sample.config, &txtedit.PrintDebugger{})
		rootNode := analyser.Run()
		reproducedText := rootNode.TextString()
		fmt.Println(txtedit.DebugNode(rootNode, 0))
		fmt.Println("Reproduced:")
		fmt.Println(reproducedText)
		if reproducedText != sample.input {
			t.Fatal("mismatch")
		}
	}
}
