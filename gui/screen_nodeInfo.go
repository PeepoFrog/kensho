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
)

// type InfoScreen struct {
// 	NodeIP string
// }

func makeNodeInfoScreen(_ fyne.Window, g *Gui) fyne.CanvasObject {
	// TODO: only for testing, delete later
	// g.Host.IP = "148.251.69.56"
	//

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
