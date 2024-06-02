package gui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"

	"github.com/KiraCore/kensho/helper/httph"
	"github.com/KiraCore/kensho/types"
	"github.com/KiraCore/kensho/types/endpoint/shidai"
)

func makeNodeInfoScreen(_ fyne.Window, g *Gui) fyne.CanvasObject {

	nodeInfoTab := container.NewTabItem("Node Info", makeNodeInfoTab(g))
	validatorInfoTab := container.NewTabItem("Validator Info", makeValidatorInfoTab(g))
	return container.NewAppTabs(nodeInfoTab, validatorInfoTab)
}

func makeValidatorInfoTab(g *Gui) fyne.CanvasObject {
	validatorStatusCheck := binding.NewBool()
	claimSeatStatusCheck := binding.NewBool()
	validatorStateCheck := binding.NewString()

	validatorControlButton := widget.NewButton("", func() {

	})
	updateBinding := binding.NewDataListener(func() {})
	validatorControlButton.Disable()
	validatorControlButton.Hide()
	refreshFunc := func() {
		g.WaitDialog.ShowWaitDialog()
		status, err := httph.GetValidatorStatus(g.sshClient, types.DEFAULT_SHIDAI_PORT)
		if err != nil {
			g.showErrorDialog(err, binding.NewDataListener(func() {}))
			log.Println(err.Error())
			return
		}

		validatorStatusCheck.Set(status.Validator)
		claimSeatStatusCheck.Set(status.ClaimSeat)
		updateBinding.DataChanged()
		g.WaitDialog.HideWaitDialog()
	}

	refreshBinding := binding.NewDataListener(refreshFunc)

	pauseValidatorFunc := func() {
		// pause
		refreshBinding.DataChanged()
	}
	unpauseValidatorFunc := func() {
		// unpause tx
		refreshBinding.DataChanged()
	}
	activateValidatorFunc := func() {
		// activate
		refreshBinding.DataChanged()
	}

	updateBinding = binding.NewDataListener(func() {
		valCheck, _ := validatorStatusCheck.Get()
		claimCheck, _ := claimSeatStatusCheck.Get()

		if claimCheck {
			if validatorControlButton.Disabled() {
				validatorControlButton.Enable()
				validatorControlButton.Show()
				validatorControlButton.SetText("Claim validator seat")
			}
		}

		if valCheck {
			status, _ := validatorStateCheck.Get()
			switch status {
			case string(shidai.Active):
				if validatorControlButton.Disabled() {
					validatorControlButton.Enable()
					validatorControlButton.SetText("Pause")
					validatorControlButton.OnTapped = pauseValidatorFunc
				}
			case string(shidai.Paused):
				if validatorControlButton.Disabled() {
					validatorControlButton.Enable()
					validatorControlButton.Show()
					validatorControlButton.SetText("Activate")
					validatorControlButton.OnTapped = unpauseValidatorFunc
				}
			case string(shidai.Inactive):
				if validatorControlButton.Disabled() {
					validatorControlButton.Enable()
					validatorControlButton.Show()
					validatorControlButton.SetText("Activate")
					validatorControlButton.OnTapped = activateValidatorFunc
				}
			}
		}
	})

	refreshButton := widget.NewButton("Refresh", func() { refreshBinding.DataChanged() })
	container.NewHBox()

	return container.NewBorder(nil, refreshButton, nil, nil)
}

func makeNodeInfoTab(g *Gui) fyne.CanvasObject {
	// TODO: only for testing, delete later
	// g.Host.IP = "148.251.69.56"

	// latest block box
	latestBlockData := binding.NewString()
	latestBlockLabel := widget.NewLabelWithData(latestBlockData)
	latestBlockBox := container.NewHBox(
		widget.NewLabel("Latest Block"), latestBlockLabel,
	)

	// validator address box
	validatorAddressData := binding.NewString()
	validatorAddressLabel := widget.NewLabelWithData(validatorAddressData)
	validatorAddressBox := container.NewHBox(
		widget.NewLabel("Validator Address: "), validatorAddressLabel,
		widget.NewButtonWithIcon("Copy", theme.ContentCopyIcon(), func() {
			data, err := validatorAddressData.Get()
			if err != nil {
				log.Println(err)
				return
			}

			err = clipboard.WriteAll(data)
			if err != nil {
				return
			}
		}),
	)
	// public ip box
	publicIpData := binding.NewString()
	publicIpLabel := widget.NewLabelWithData(publicIpData)
	publicIpBox := container.NewHBox(
		widget.NewLabel("Public IP Address: "), publicIpLabel,
	)

	// miss chance box
	missChanceData := binding.NewString()
	missChanceLabel := widget.NewLabelWithData(missChanceData)
	missChanceBox := container.NewHBox(
		widget.NewLabel("Miss Chance: "), missChanceLabel,
	)

	refreshScreen := func() {
		g.WaitDialog.ShowWaitDialog()
		defer g.WaitDialog.HideWaitDialog()
		i, err := httph.GetInterxStatus(g.Host.IP)
		if err != nil {
			return
		}
		latestBlockData.Set(i.InterxInfo.LatestBlockHeight)

	}

	refreshButton := widget.NewButton("Refresh", refreshScreen)
	sendSekaiCommandButton := widget.NewButton("Execute sekai command", func() { showSekaiExecuteDialog(g) })
	mainInfo := container.NewVScroll(
		container.NewVBox(
			latestBlockBox,
			validatorAddressBox,
			publicIpBox,
			missChanceBox,
		),
	)
	return container.NewBorder(sendSekaiCommandButton, refreshButton, nil, nil, mainInfo)
}
