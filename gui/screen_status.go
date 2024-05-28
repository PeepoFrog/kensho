package gui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/KiraCore/kensho/helper/httph"
)

func makeStatusScreen(_ fyne.Window, g *Gui) fyne.CanvasObject {
	const STATUS_Unavailable = "Unavailable"
	const STATUS_Running = "Running"

	var dataListenerForSuccesses binding.DataListener
	deployButton := widget.NewButton("Deploy", func() {
		showDeployDialog(g, dataListenerForSuccesses)
	})
	deployButton.Disable()

	interxStatusBinding := binding.NewBool()
	interxStatusInfo := widget.NewLabel("")
	interxInfoBox := container.NewHBox(
		widget.NewLabel("Interx:"),
		interxStatusInfo,
	)

	shidaiStatusBinding := binding.NewBool()
	shidaiStatusInfo := widget.NewLabel("")
	shidaiInfoBox := container.NewHBox(
		widget.NewLabel("Shidai:"),
		shidaiStatusInfo,
	)

	sekaiStatusInfo := widget.NewLabel("")
	sekaiStatusBinding := binding.NewBool()
	sekaiStatusBinding.AddListener(binding.NewDataListener(func() {
		check, _ := sekaiStatusBinding.Get()
		if check {
			sekaiStatusInfo.SetText(STATUS_Running)
		} else {
			sekaiStatusInfo.SetText(STATUS_Unavailable)
		}
	}))
	sekaiInfoBox := container.NewHBox(
		widget.NewLabel("Sekai:"),
		sekaiStatusInfo,
	)

	checkInterxStatus := func() {
		_, err := httph.GetInterxStatus(g.Host.IP)
		if err != nil {
			interxStatusInfo.SetText(STATUS_Unavailable)
			log.Printf("ERROR getting interx status: %v", err)
			err = interxStatusBinding.Set(false)
			if err != nil {
				log.Printf("ERROR setting binding: %v", err)
				return
			}
		} else {
			err = interxStatusBinding.Set(true)
			interxStatusInfo.SetText(STATUS_Running)
			if err != nil {
				log.Printf("%v", err)
				return
			}
		}

	}

	checkShidaiStatus := func() {
		_, err := httph.MakeHttpRequest(fmt.Sprintf("http://%v:%v/status", g.Host.IP, 8282), "GET")
		if err != nil {
			log.Printf("ERROR: %v", err)
			shidaiStatusInfo.SetText(STATUS_Unavailable)
			err = shidaiStatusBinding.Set(false)
			if err != nil {
				log.Printf("ERROR: %v", err)
				return
			}
		} else {
			shidaiStatusInfo.SetText(STATUS_Running)
			err = shidaiStatusBinding.Set(true)
			if err != nil {
				return
			}
		}
	}
	checkSekaiStatus := func() {
		_, err := httph.GetSekaiStatus(g.Host.IP, "26657")
		if err != nil {
			sekaiStatusBinding.Set(false)
		} else {
			sekaiStatusBinding.Set(true)
		}
	}

	refresh := func() {
		g.WaitDialog.ShowWaitDialog()
		checkInterxStatus()
		checkShidaiStatus()
		checkSekaiStatus()
		shidaiCheck, _ := shidaiStatusBinding.Get()
		sekaiCheck, _ := sekaiStatusBinding.Get()
		interxCheck, _ := interxStatusBinding.Get()

		if !shidaiCheck && !interxCheck && !sekaiCheck {
			deployButton.Enable()
		}

		g.WaitDialog.HideWaitDialog()
	}

	refreshButton := widget.NewButton("Refresh", func() {
		refresh()
	})

	dataListenerForSuccesses = binding.NewDataListener(func() {
		deployButton.Disable()
		refresh()
	})

	// refresh()
	return container.NewBorder(nil, refreshButton, nil, nil,
		container.NewVBox(
			deployButton,
			interxInfoBox,
			sekaiInfoBox,
			shidaiInfoBox,
		))

}
