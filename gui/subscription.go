package gui

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type subscription struct {
	subWindow    *widgets.QMainWindow
	subLabel     *widgets.QLabel
	subCombobox  *widgets.QComboBox
	deleteButton *widgets.QPushButton
	lineText     *widgets.QLineEdit
	addButton    *widgets.QPushButton
}

func NewSubscription(parent *widgets.QMainWindow) *widgets.QMainWindow {
	s := &subscription{}
	s.subWindow = widgets.NewQMainWindow(nil, core.Qt__Window)
	s.subWindow.SetWindowTitle("subscription")
	s.subWindow.ConnectCloseEvent(func(event *gui.QCloseEvent) {
		event.Ignore()
		s.subWindow.Hide()
	})
	s.subWindow.ConnectShowEvent(func(event *gui.QShowEvent) {
		s.subCombobox.Clear()
		//link, err := subscr.GetLink()
		//if err != nil {
		//	MessageBox(err.Error())
		//}
		links, err := apiC.GetSubLinks(apiCtx(), &empty.Empty{})
		if err != nil {
			MessageBox(err.Error())
			return
		}
		s.subCombobox.AddItems(links.Value)
	})

	s.subInit()
	s.setLayout()
	//s.setGeometry()
	s.setListener()

	return s.subWindow
}

func (s *subscription) subInit() {
	s.subLabel = widgets.NewQLabel2("subscription", s.subWindow, core.Qt__WindowType(0x00000000))
	s.subCombobox = widgets.NewQComboBox(s.subWindow)
	s.deleteButton = widgets.NewQPushButton2("delete", s.subWindow)
	s.lineText = widgets.NewQLineEdit(s.subWindow)
	s.addButton = widgets.NewQPushButton2("add", s.subWindow)
}

func (s *subscription) setLayout() {
	windowLayout := widgets.NewQGridLayout2()
	windowLayout.AddWidget2(s.subLabel, 0, 0, 0)
	windowLayout.AddWidget2(s.subCombobox, 0, 1, 0)
	windowLayout.AddWidget2(s.deleteButton, 0, 2, 0)
	windowLayout.AddWidget3(s.lineText, 1, 0, 1, 2, 0)
	windowLayout.AddWidget2(s.addButton, 1, 2, 0)

	centralWidget := widgets.NewQWidget(s.subWindow, 0)
	centralWidget.SetLayout(windowLayout)
	s.subWindow.SetCentralWidget(centralWidget)
}

func (s *subscription) setGeometry() {
	s.subLabel.SetGeometry(core.NewQRect2(core.NewQPoint2(10, 10), core.NewQPoint2(130, 40)))
	s.subCombobox.SetGeometry(core.NewQRect2(core.NewQPoint2(115, 10), core.NewQPoint2(600, 40)))
	s.deleteButton.SetGeometry(core.NewQRect2(core.NewQPoint2(610, 10), core.NewQPoint2(690, 40)))
	s.lineText.SetGeometry(core.NewQRect2(core.NewQPoint2(115, 50), core.NewQPoint2(600, 80)))
	s.addButton.SetGeometry(core.NewQRect2(core.NewQPoint2(610, 50), core.NewQPoint2(690, 80)))
}

func (s *subscription) setListener() {
	s.deleteButton.ConnectClicked(func(bool2 bool) {
		//linkToDelete := s.subCombobox.CurrentText()
		//if err := subscr.RemoveLinkJSON(linkToDelete); err != nil {
		//	MessageBox(err.Error())
		//	return
		//}
		links, err := apiC.DeleteSubLink(apiCtx(), &wrappers.StringValue{Value: s.subCombobox.CurrentText()})
		if err != nil {
			MessageBox(err.Error())
			return
		}
		s.subCombobox.Clear()
		//link, err := subscr.GetLink()
		//if err != nil {
		//	MessageBox(err.Error())
		//	return
		//}
		s.subCombobox.AddItems(links.Value)
	})

	s.addButton.ConnectClicked(func(bool2 bool) {
		//link, err := subscr.GetLink()
		//if err != nil {
		//	MessageBox(err.Error())
		//	return
		//}
		links, err := apiC.GetSubLinks(apiCtx(), &empty.Empty{})
		if err != nil {
			MessageBox(err.Error())
			return
		}
		linkToAdd := s.lineText.Text()
		if linkToAdd == "" {
			return
		}
		for index := range links.Value {
			if links.Value[index] == linkToAdd {
				return
			}
		}

		links, err = apiC.AddSubLink(apiCtx(), &wrappers.StringValue{Value: linkToAdd})
		if err != nil {
			MessageBox(err.Error())
			return
		}
		s.subCombobox.Clear()
		s.subCombobox.AddItems(links.Value)
		s.lineText.Clear()
	})
}
