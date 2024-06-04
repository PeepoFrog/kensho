package gui

import (
	"encoding/json"
	"log"
	"strings"

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

	// nodeInfoTab := container.NewTabItem("Node Info", makeNodeInfoTab(g))
	// validatorInfoTab := container.NewTabItem("Validator Info", makeValidatorInfoTab(g))
	// return container.NewAppTabs(nodeInfoTab, validatorInfoTab)

	return makeNodeInfoTab(g)
}

func makeNodeInfoTab(g *Gui) fyne.CanvasObject {
	// TODO: only for testing, delete later
	// g.Host.IP = "148.251.69.56"
	var claimSeat bool
	// latest block box
	var refreshBinding binding.DataListener
	validatorControlButton := widget.NewButton("", func() {})

	claimSeatButton := widget.NewButton("Claim validator seat", func() {})

	claimSeatButton.Hide()

	latestBlockData := binding.NewString()
	latestBlockLabel := widget.NewLabelWithData(latestBlockData)
	latestBlockBox := container.NewHBox(
		widget.NewLabel("Latest Block:"), latestBlockLabel,
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

	//validator status (active, paused, etc...)
	validatorStatusData := binding.NewString()
	validatorStatusLabel := widget.NewLabelWithData(validatorStatusData)
	validatorStatusBox := container.NewHBox(
		widget.NewLabel("Validator Status: "), validatorStatusLabel,
		validatorControlButton,
	)
	// nodeID
	nodeIDData := binding.NewString()
	nodeIDLabel := widget.NewLabelWithData(nodeIDData)
	nodeIDBox := container.NewHBox(
		widget.NewLabel("Node ID:"), nodeIDLabel,
	)

	topData := binding.NewString()
	topLabel := widget.NewLabelWithData(topData)
	topBox := container.NewHBox(
		widget.NewLabel("Top:"), topLabel,
	)

	// public ip box
	// publicIpData := binding.NewString()
	// publicIpLabel := widget.NewLabelWithData(publicIpData)
	// publicIpBox := container.NewHBox(
	// 	widget.NewLabel("Public IP Address: "), publicIpLabel,
	// )

	// miss chance box
	missChanceData := binding.NewString()
	missChanceLabel := widget.NewLabelWithData(missChanceData)
	missChanceBox := container.NewHBox(
		widget.NewLabel("Miss Chance: "), missChanceLabel,
	)

	lastProducedBlockData := binding.NewString()
	lastProducedLabel := widget.NewLabelWithData(lastProducedBlockData)
	lastProducedBox := container.NewHBox(
		widget.NewLabel("Last produced block: "), lastProducedLabel,
	)

	validatorControlButton.Disable()
	validatorControlButton.Hide()

	execFunc := func(arg types.Cmd) {
		g.WaitDialog.ShowWaitDialog()
		payload, err := json.Marshal(types.ExecSekaiCmd{TX: arg})
		log.Printf("Executing: %v", arg)
		if err != nil {
			g.showErrorDialog(err, binding.NewDataListener(func() {}))
		}
		httph.ExecHttpRequestBySSHTunnel(g.sshClient, types.SEKIN_EXECUTE_CMD_ENDPOINT, "POST", payload)
		g.WaitDialog.HideWaitDialog()

		refreshBinding.DataChanged()
	}

	claimSeatButton.OnTapped = func() {
		//claim seat
		execFunc(types.ClaimValidatorSeat)
	}
	pauseValidatorFunc := func() {
		// pause
		execFunc(types.Pause)
	}
	unpauseValidatorFunc := func() {
		// unpause tx
		execFunc(types.Unpause)
	}
	activateValidatorFunc := func() {
		// activate
		execFunc(types.Activate)

	}

	errBinding := binding.NewUntyped()
	refreshScreen := func() {
		g.WaitDialog.ShowWaitDialog()
		defer g.WaitDialog.HideWaitDialog()
		dashboardData, err := httph.GetDashboardInfo(g.sshClient, 8282)
		if err != nil {
			errBinding.Set(err)
			return
		}

		log.Printf("dashboard %+v", dashboardData)
		claimSeat = dashboardData.SeatClaimAvailable
		if claimSeat {
			claimSeatButton.Show()
		} else {
			claimSeatButton.Hide()
		}

		if dashboardData.ValidatorStatus != "Unknown" {
			status := strings.ToUpper(dashboardData.ValidatorStatus)
			switch status {
			case string(shidai.Active):
				if validatorControlButton.Disabled() {
					validatorControlButton.Enable()
					validatorControlButton.Show()
					validatorControlButton.SetText("Pause")
					validatorControlButton.OnTapped = pauseValidatorFunc
				}
			case string(shidai.Paused):
				if validatorControlButton.Disabled() {
					validatorControlButton.Enable()
					validatorControlButton.Show()
					validatorControlButton.SetText("Unpause")
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
		nodeIDData.Set(dashboardData.NodeID)
		topData.Set(dashboardData.Top)
		validatorAddressData.Set(dashboardData.ValidatorAddress)
		missChanceData.Set(dashboardData.Mischance)
		latestBlockData.Set(dashboardData.Blocks)
		lastProducedBlockData.Set(dashboardData.LastProducedBlock)
		validatorStatusData.Set(dashboardData.ValidatorStatus)
	}
	refreshBinding = binding.NewDataListener(func() {
		g.WaitDialog.ShowWaitDialog()
		refreshScreen()
		g.WaitDialog.HideWaitDialog()
		err, _ := errBinding.Get()
		if err != nil {
			if e, ok := err.(error); ok {
				g.showErrorDialog(e, binding.NewDataListener(func() {}))
			}
		}
	})

	refreshButton := widget.NewButton("Refresh", refreshBinding.DataChanged)
	sendSekaiCommandButton := widget.NewButton("Execute sekai command", func() { showSekaiExecuteDialog(g) })
	mainInfo := container.NewVScroll(
		container.NewVBox(
			// publicIpBox,
			validatorAddressBox,
			validatorStatusBox,
			nodeIDBox,
			topBox,
			latestBlockBox,
			lastProducedBox,
			missChanceBox,
		),
	)
	return container.NewBorder(sendSekaiCommandButton, refreshButton, nil, nil, mainInfo)
}
