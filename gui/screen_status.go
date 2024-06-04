package gui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/KiraCore/kensho/helper/httph"
	"github.com/KiraCore/kensho/types"
)

func makeStatusScreen(_ fyne.Window, g *Gui) fyne.CanvasObject {
	const STATUS_Unavailable = "Unavailable"
	const STATUS_Running = "Running"

	// shidaiStatusBinding := binding.NewUntyped()

	// getShidaiStatus := func() shidai.Status {
	// 	status, _ := shidaiStatusBinding.Get()
	// 	return status.(shidai.Status)
	// }

	// setShidaiStatus := func(status shidai.Status) {
	// 	shidaiStatusBinding.Set(status)
	// }

	interxStatusCheck := binding.NewBool()
	interxInfraCheck := binding.NewBool()
	interxStatusInfo := widget.NewLabel("")
	interxInfoBox := container.NewHBox(
		widget.NewLabel("Interx:"),
		interxStatusInfo,
	)

	shidaiStatusCheck := binding.NewBool()
	shidaiInfraCheck := binding.NewBool()
	shidaiStatusInfo := widget.NewLabel("")
	shidaiInfoBox := container.NewHBox(
		widget.NewLabel("Shidai:"),
		shidaiStatusInfo,
	)

	sekaiStatusCheck := binding.NewBool()
	sekaiInfraCheck := binding.NewBool()
	sekaiStatusInfo := widget.NewLabel("")
	sekaiInfoBox := container.NewHBox(
		widget.NewLabel("Sekai:"),
		sekaiStatusInfo,
	)

	var dataListenerForSuccesses binding.DataListener
	deployButton := widget.NewButton("Deploy", func() {
		showDeployDialog(g, dataListenerForSuccesses, shidaiInfraCheck)
	})
	deployButton.Disable()

	checkInterxStatus := func() {
		_, err := httph.GetInterxStatus(g.Host.IP)
		if err != nil {
			log.Printf("ERROR getting interx status: %v", err)
			interxStatusInfo.SetText(STATUS_Unavailable)
			interxStatusCheck.Set(false)

		} else {
			interxStatusCheck.Set(true)
			interxStatusInfo.SetText(STATUS_Running)
		}

	}

	checkShidaiStatus := func() {
		shidaiStatus, err := httph.GetShidaiStatus(g.sshClient, types.DEFAULT_SHIDAI_PORT)
		if err != nil {
			log.Printf("ERROR: %v", err)
			shidaiStatusInfo.SetText(STATUS_Unavailable)
			shidaiStatusCheck.Set(false)

		} else {
			log.Println("switching  shidai state")
			shidaiStatusInfo.SetText(STATUS_Running)
			shidaiInfraCheck.Set(shidaiStatus.Shidai.Infra)
			sekaiInfraCheck.Set(shidaiStatus.Sekai.Infra)
			interxInfraCheck.Set(shidaiStatus.Interx.Infra)
			shidaiStatusCheck.Set(true)
		}
	}

	checkSekaiStatus := func() {
		_, err := httph.GetSekaiStatus(g.Host.IP, "26657")
		if err != nil {
			log.Printf("ERROR: %v", err)
			sekaiStatusInfo.SetText(STATUS_Unavailable)
			sekaiStatusCheck.Set(false)

		} else {
			sekaiStatusInfo.SetText(STATUS_Running)
			sekaiStatusCheck.Set(true)
		}
	}
	refresh := func() {
		g.WaitDialog.ShowWaitDialog()
		checkInterxStatus()
		checkShidaiStatus()
		checkSekaiStatus()
		shidaiCheck, _ := shidaiStatusCheck.Get()
		sekaiCheck, _ := sekaiStatusCheck.Get()
		interxCheck, _ := interxStatusCheck.Get()

		shidaiInfra, _ := shidaiInfraCheck.Get()
		sekaiInfra, _ := sekaiInfraCheck.Get()
		interxInfra, _ := interxInfraCheck.Get()
		log.Printf("CHECKS: shidaiCheck:%v sekaiCheck:%v interxCheck:%v shidaiInfra:%v sekaiInfra:%v interxInfra:%v",
			shidaiCheck, sekaiCheck, interxCheck, shidaiInfra, sekaiInfra, interxInfra)

		// TODO: first maybe we should try to restart first if shidai is not running
		deployButtonCheck := ((!shidaiCheck && !interxCheck && !sekaiCheck) || (shidaiInfra && (!sekaiInfra && !interxInfra)))
		log.Println("enable state: ", deployButtonCheck)
		if deployButtonCheck {
			deployButton.Enable()
		}

		defer g.WaitDialog.HideWaitDialog()
	}

	refreshButton := widget.NewButton("Refresh", func() {
		refresh()
	})

	dataListenerForSuccesses = binding.NewDataListener(func() {
		log.Println("triggering dataListenerForSuccesses")

		deployButton.Disable()
		refresh()
	})
	defer refresh()
	return container.NewBorder(nil, refreshButton, nil, nil,
		container.NewVBox(
			deployButton,
			interxInfoBox,
			sekaiInfoBox,
			shidaiInfoBox,
		))

}
