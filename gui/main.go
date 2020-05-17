package gui

import (
	"fmt"
	"github.com/Asutorufa/yuhaiin/net/common"
	"github.com/Asutorufa/yuhaiin/net/delay"
	"github.com/Asutorufa/yuhaiin/process"
	"github.com/Asutorufa/yuhaiin/subscr"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"time"
)

func (sGui *SGui) createMainWindow() {
	sGui.MainWindow = widgets.NewQMainWindow(nil, 0)
	sGui.MainWindow.SetFixedSize2(600, 400)
	sGui.MainWindow.SetWindowTitle("yuhaiin")
	img := gui.NewQPixmap()
	img.LoadFromData2(core.QByteArray_FromBase64(core.NewQByteArray2(icon, len(icon))), "svg", core.Qt__AutoColor)
	icon2 := gui.NewQIcon2(img)
	sGui.MainWindow.SetWindowIcon(icon2)

	trayIcon := widgets.NewQSystemTrayIcon(sGui.MainWindow)
	trayIcon.SetIcon(icon2)
	menu := widgets.NewQMenu(nil)
	ssrMicroClientTrayIconMenu := widgets.NewQAction2("yuhaiin", sGui.MainWindow)
	ssrMicroClientTrayIconMenu.ConnectTriggered(func(bool2 bool) {
		sGui.openMainWindow()
	})
	subscriptionTrayIconMenu := widgets.NewQAction2("subscription", sGui.MainWindow)
	subscriptionTrayIconMenu.ConnectTriggered(func(bool2 bool) {
		sGui.openSubscriptionWindow()
	})

	settingTrayIconMenu := widgets.NewQAction2("setting", sGui.MainWindow)
	settingTrayIconMenu.ConnectTriggered(func(bool2 bool) {
		sGui.openSettingWindow()
	})

	exit := widgets.NewQAction2("exit", sGui.MainWindow)
	exit.ConnectTriggered(func(bool2 bool) {
		sGui.App.Quit()
	})

	actions := []*widgets.QAction{ssrMicroClientTrayIconMenu,
		subscriptionTrayIconMenu, settingTrayIconMenu, exit}
	menu.AddActions(actions)
	trayIcon.SetContextMenu(menu)
	trayIcon.SetToolTip("")
	trayIcon.Show()

	//statusLabel := widgets.NewQLabel2("status", sGui.MainWindow, core.Qt__WindowType(0x00000000))
	//statusLabel.SetGeometry(core.NewQRect2(core.NewQPoint2(40, 10), core.NewQPoint2(130, 40)))
	statusLabel2 := widgets.NewQLabel2("", sGui.MainWindow, core.Qt__WindowType(0x00000000))
	statusLabel2.SetGeometry(core.NewQRect2(core.NewQPoint2(40, 10), core.NewQPoint2(560, 40)))

	nowNodeLabel := widgets.NewQLabel2("Now Use", sGui.MainWindow, core.Qt__WindowType(0x00000000))
	nowNodeLabel.SetGeometry(core.NewQRect2(core.NewQPoint2(40, 60), core.NewQPoint2(130, 90)))
	//nowNodeName, nowNodeGroup := subscr.GetNowNodeGroupAndName()
	nowNodeLabel2 := widgets.NewQLabel2("", sGui.MainWindow, core.Qt__WindowType(0x00000000))
	nowNodeLabel2.SetGeometry(core.NewQRect2(core.NewQPoint2(130, 60), core.NewQPoint2(560, 90)))

	groupLabel := widgets.NewQLabel2("Group", sGui.MainWindow, core.Qt__WindowType(0x00000000))
	groupLabel.SetGeometry(core.NewQRect2(core.NewQPoint2(40, 110), core.NewQPoint2(130, 140)))
	groupCombobox := widgets.NewQComboBox(sGui.MainWindow)
	//group, err := subscr.GetGroup()
	//if err != nil {
	//	sGui.MessageBox(err.Error())
	//	return
	//}
	//groupCombobox.AddItems(group)
	//groupCombobox.SetCurrentTextDefault(nowNodeGroup)
	groupCombobox.SetGeometry(core.NewQRect2(core.NewQPoint2(130, 110), core.NewQPoint2(450, 140)))
	//refreshButton := widgets.NewQPushButton2("refresh", sGui.MainWindow)
	//refreshButton.SetGeometry(core.NewQRect2(core.NewQPoint2(460, 110), core.NewQPoint2(560, 140)))

	nodeLabel := widgets.NewQLabel2("Node", sGui.MainWindow, core.Qt__WindowType(0x00000000))
	nodeLabel.SetGeometry(core.NewQRect2(core.NewQPoint2(40, 160), core.NewQPoint2(130, 190)))
	nodeCombobox := widgets.NewQComboBox(sGui.MainWindow)
	//node, err := subscr.GetNode(groupCombobox.CurrentText())
	//if err != nil {
	//	sGui.MessageBox(err.Error())
	//	return
	//}
	//nodeCombobox.AddItems(node)
	//nodeCombobox.SetCurrentTextDefault(nowNodeName)
	nodeCombobox.SetGeometry(core.NewQRect2(core.NewQPoint2(130, 160), core.NewQPoint2(450, 190)))

	startButton := widgets.NewQPushButton2("Use", sGui.MainWindow)
	startButton.ConnectClicked(func(bool2 bool) {
		group := groupCombobox.CurrentText()
		remarks := nodeCombobox.CurrentText()
		if err := subscr.ChangeNowNode(group, remarks); err != nil {
			sGui.MessageBox(err.Error())
			return
		}
		if err := process.ChangeNode(); err != nil {
			sGui.MessageBox(err.Error())
		}
		nowNodeLabel2.SetText(remarks + " - " + group)
	})

	startButton.SetGeometry(core.NewQRect2(core.NewQPoint2(460, 160),
		core.NewQPoint2(560, 190)))

	delayLabel := widgets.NewQLabel2("Latency", sGui.MainWindow,
		core.Qt__WindowType(0x00000000))
	delayLabel.SetGeometry(core.NewQRect2(core.NewQPoint2(40, 210),
		core.NewQPoint2(130, 240)))
	delayLabel2 := widgets.NewQLabel2("", sGui.MainWindow,
		core.Qt__WindowType(0x00000000))
	delayLabel2.SetGeometry(core.NewQRect2(core.NewQPoint2(130, 210),
		core.NewQPoint2(450, 240)))
	delayButton := widgets.NewQPushButton2("Get", sGui.MainWindow)
	delayButton.ConnectClicked(func(bool2 bool) {
		go func() {
			group := groupCombobox.CurrentText()
			remarks := nodeCombobox.CurrentText()
			server, port := subscr.GetOneNodeAddress(group, remarks)
			delayTmp, isSuccess, err := delay.TCPDelay(server, port)
			var delayString string
			if err != nil {
				delayString = err.Error()
			} else {
				delayString = delayTmp.String()
			}
			if isSuccess == false {
				delayString = "delay > 3s or server can not connect"
			}
			delayLabel2.SetText(delayString)
		}()
	})
	delayButton.SetGeometry(core.NewQRect2(core.NewQPoint2(460, 210), core.NewQPoint2(560, 240)))

	groupCombobox.ConnectCurrentTextChanged(func(string2 string) {
		node, err := subscr.GetNode(groupCombobox.CurrentText())
		if err != nil {
			sGui.MessageBox(err.Error())
			return
		}
		nodeCombobox.Clear()
		nodeCombobox.AddItems(node)
	})

	subButton := widgets.NewQPushButton2("Subscription Setting", sGui.MainWindow)
	subButton.SetGeometry(core.NewQRect2(core.NewQPoint2(40, 260), core.NewQPoint2(290, 290)))
	subButton.ConnectClicked(func(bool2 bool) {
		sGui.openSubscriptionWindow()
	})

	update := func() {
		group, err := subscr.GetGroup()
		if err != nil {
			sGui.MessageBox(err.Error())
			return
		}
		groupCombobox.Clear()
		groupCombobox.AddItems(group)
		node, err := subscr.GetNode(groupCombobox.CurrentText())
		if err != nil {
			sGui.MessageBox(err.Error())
			return
		}
		nodeCombobox.Clear()
		nodeCombobox.AddItems(node)

		nowNodeName, nowNodeGroup := subscr.GetNowNodeGroupAndName()
		groupCombobox.SetCurrentText(nowNodeGroup)
		nodeCombobox.SetCurrentText(nowNodeName)
		nowNodeLabel2.SetText(nowNodeName)
	}

	subUpdateButton := widgets.NewQPushButton2("Subscription Update", sGui.MainWindow)
	subUpdateButton.SetGeometry(core.NewQRect2(core.NewQPoint2(300, 260), core.NewQPoint2(560, 290)))
	subUpdateButton.ConnectClicked(func(bool2 bool) {
		message := widgets.NewQMessageBox(sGui.MainWindow)
		message.SetText("Updating!")
		message.Show()
		if err := subscr.GetLinkFromInt(); err != nil {
			sGui.MessageBox(err.Error())
		}
		message.SetText("Updated!")
		update()
	})

	settingButton := widgets.NewQPushButton2("Setting", sGui.MainWindow)
	settingButton.SetGeometry(core.NewQRect2(core.NewQPoint2(40, 300), core.NewQPoint2(290, 330)))
	settingButton.ConnectClicked(func(bool2 bool) {
		sGui.openSettingWindow()
	})

	statusRefreshIsRun := false
	sGui.MainWindow.ConnectShowEvent(func(event *gui.QShowEvent) {
		go func() {
			if statusRefreshIsRun {
				return
			}
			statusRefreshIsRun = true
			downloadTmp, downRate := 0, 0
			uploadTmp, uploadRate := 0, 0
			for {
				if sGui.MainWindow.IsHidden() {
					break
				}
				downRate = common.DownloadTotal - downloadTmp
				downloadTmp = common.DownloadTotal
				uploadRate = common.UploadTotal - uploadTmp
				uploadTmp = common.UploadTotal
				statusLabel2.SetText(fmt.Sprintf("Download<sub><i>(%s)</i></sub>: %s/S , Upload<sub><i>(%s)</i></sub>: %s/S",
					common.ReducedUnit2(float64(downloadTmp)),
					common.ReducedUnit2(float64(downRate)),
					common.ReducedUnit2(float64(uploadTmp)),
					common.ReducedUnit2(float64(uploadRate))))
				time.Sleep(time.Second)
			}
			statusRefreshIsRun = false
		}()
		update()
	})
}