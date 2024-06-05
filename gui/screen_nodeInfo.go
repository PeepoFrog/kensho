package gui

import (
	"encoding/json"
	"fmt"
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
	validatorControlButton.Disable()
	validatorControlButton.Hide()

	claimSeatButton := widget.NewButton("Claim validator seat", func() {})
	claimSeatButton.Hide()

	latestBlockData := binding.NewString()
	latestBlockLabel := widget.NewLabelWithData(latestBlockData)
	latestBlockBox := container.NewHBox(
		widget.NewLabel("Block:"), latestBlockLabel,
	)

	// validator address box
	validatorAddressData := binding.NewString()
	validatorAddressLabel := widget.NewLabelWithData(validatorAddressData)
	validatorAddressBox := container.NewHBox(
		widget.NewLabel("Val.Addr: "), validatorAddressLabel,
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
		widget.NewLabel("Val.Status:\t\t"), validatorStatusLabel,
		validatorControlButton,
	)
	// nodeID
	nodeIDData := binding.NewString()
	nodeIDLabel := widget.NewLabelWithData(nodeIDData)
	nodeIDBox := container.NewHBox(
		widget.NewLabel("Val.NodeID:\t\t"), nodeIDLabel,
	)

	topData := binding.NewString()
	topLabel := widget.NewLabelWithData(topData)
	topBox := container.NewHBox(
		widget.NewLabel("Rank:\t\t"), topLabel,
	)

	// miss chance box
	missChanceData := binding.NewString()
	missChanceLabel := widget.NewLabelWithData(missChanceData)
	missChanceBox := container.NewHBox(
		widget.NewLabel("Miss:\t\t"), missChanceLabel,
	)

	lastProducedBlockData := binding.NewString()
	lastProducedLabel := widget.NewLabelWithData(lastProducedBlockData)
	lastProducedBox := container.NewHBox(
		widget.NewLabel("Latest Block:\t\t"), lastProducedLabel,
	)

	// miss Chance Confidence box
	missChanceConfidenceData := binding.NewString()
	missChanceConfidenceLabel := widget.NewLabelWithData(missChanceConfidenceData)
	missChanceConfidenceBox := container.NewHBox(
		widget.NewLabel("Miss conf.\t\t"), missChanceConfidenceLabel,
	)

	startHeightData := binding.NewString()
	startHeightLabel := widget.NewLabelWithData(startHeightData)
	startHeightBox := container.NewHBox(
		widget.NewLabel("Start Height:\t\t"), startHeightLabel,
	)

	producedBlocksData := binding.NewString()
	producedBlocksLabel := widget.NewLabelWithData(producedBlocksData)
	producedBlocksBox := container.NewHBox(
		widget.NewLabel("Produced:\t\t"), producedBlocksLabel,
	)

	monikerData := binding.NewString()
	monikerLabel := widget.NewLabelWithData(monikerData)
	monikerBox := container.NewHBox(
		widget.NewLabel("Moniker:\t\t"), monikerLabel,
	)

	genesisChecksumData := binding.NewString()
	genesisChecksumLabel := widget.NewLabelWithData(genesisChecksumData)
	genesisChecksumBox := container.NewHBox(
		widget.NewLabel("Gen. SHA256::\t\t"), genesisChecksumLabel,
	)

	// Numerical validators status
	activeValidatorsData := binding.NewInt()
	activeValidatorsLabel := widget.NewLabelWithData(binding.IntToString(activeValidatorsData))
	activeValidatorsBox := container.NewHBox(
		widget.NewLabel("Active:"), activeValidatorsLabel,
	)

	pausedValidatorsData := binding.NewInt()
	pausedValidatorsLabel := widget.NewLabelWithData(binding.IntToString(pausedValidatorsData))
	pausedValidatorsBox := container.NewHBox(
		widget.NewLabel("Paused:"), pausedValidatorsLabel,
	)

	inactiveValidatorsData := binding.NewInt()
	inactiveValidatorsLabel := widget.NewLabelWithData(binding.IntToString(inactiveValidatorsData))
	inactiveValidatorsBox := container.NewHBox(
		widget.NewLabel("Inactive:"), inactiveValidatorsLabel,
	)

	jailedValidatorsData := binding.NewInt()
	jailedValidatorsLabel := widget.NewLabelWithData(binding.IntToString(jailedValidatorsData))
	jailedValidatorsBox := container.NewHBox(
		widget.NewLabel("Jailed:"), jailedValidatorsLabel,
	)

	waitingValidatorsData := binding.NewInt()
	waitingValidatorsLabel := widget.NewLabelWithData(binding.IntToString(waitingValidatorsData))
	waitingValidatorsBox := container.NewHBox(
		widget.NewLabel("Waiting:"), waitingValidatorsLabel,
	)

	chainIDData := binding.NewString()
	chainIDLabel := widget.NewLabelWithData(chainIDData)
	chainIDBox := container.NewHBox(
		widget.NewLabel("ChainID:"), chainIDLabel,
	)

	// dateData := binding.NewString()
	// dateLabel := widget.NewLabelWithData(dateData)
	// dateBox := container.NewHBox(
	// 	widget.NewLabel("Date:\t\t"), dateLabel,
	// )

	streakData := binding.NewString()
	streakLabel := widget.NewLabelWithData(streakData)
	streakBox := container.NewHBox(
		widget.NewLabel("Streak:\t\t"), streakLabel,
	)

	//
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
			errBinding.Set(fmt.Errorf("ERROR: getting dashboard info: %w", err))
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

		validatorStatusData.Set(dashboardData.ValidatorStatus)

		latestBlockData.Set(dashboardData.Blocks)
		lastProducedBlockData.Set(dashboardData.LastProducedBlock)
		missChanceConfidenceData.Set(dashboardData.MischanceConfidence)
		producedBlocksData.Set(dashboardData.ProducedBlocks)

		startHeightData.Set(dashboardData.StartHeight)
		monikerData.Set(dashboardData.Moniker)
		activeValidatorsData.Set(dashboardData.ActiveValidators)
		pausedValidatorsData.Set(dashboardData.PausedValidators)
		inactiveValidatorsData.Set(dashboardData.InactiveValidators)
		jailedValidatorsData.Set(dashboardData.JailedValidators)
		waitingValidatorsData.Set(dashboardData.WaitingValidators)
		genesisChecksumData.Set(dashboardData.GenesisChecksum)
		chainIDData.Set(dashboardData.ChainID)
		// dateData.Set(dashboardData.Date)
		streakData.Set(dashboardData.Streak)
		// dashboardData.CatchingUp,
		// widget.NewLabel()

	}
	refreshBinding = binding.NewDataListener(func() {
		refreshScreen()
		err, _ := errBinding.Get()
		if err != nil {
			if e, ok := err.(error); ok {
				g.showErrorDialog(e, binding.NewDataListener(func() {}))
			}
		}
	})
	refreshButton := widget.NewButton("Refresh", refreshBinding.DataChanged)
	// sendSekaiCommandButton := widget.NewButton("Execute sekai command", func() { showSekaiExecuteDialog(g) })
	mainInfo := container.NewVScroll(
		container.NewVBox(
			container.NewVBox(
				widget.NewLabelWithStyle("CHAIN INFO", fyne.TextAlignCenter,
					fyne.TextStyle{Monospace: true, Bold: true, Symbol: true}),
			),
			container.NewHBox(activeValidatorsBox, pausedValidatorsBox),
			container.NewHBox(inactiveValidatorsBox, waitingValidatorsBox),
			jailedValidatorsBox,
			chainIDBox,
			widget.NewLabel(""),
			monikerBox,

			latestBlockBox,
			lastProducedBox,
			producedBlocksBox,
			streakBox,
			topBox,
			startHeightBox,
			missChanceBox,
			missChanceConfidenceBox,
			validatorAddressBox,
			validatorStatusBox,
			nodeIDBox,
			genesisChecksumBox,

			// validatorAddressBox,
			// validatorStatusBox,
			// nodeIDBox,
			// topBox,
			// latestBlockBox,
			// lastProducedBox,
			// missChanceBox,
			// missChanceConfidenceBox,
			// startHeightBox,
			// producedBlocksBox,
			// monikerBox,
			// genesisChecksumBox,
			// chainIDBox,
			// dateBox,
			// streakBox,

			// container.N
		),
	)
	// mainInfo := container.NewVScroll(
	// 	container.NewVBox(
	// 		container.NewVBox(
	// 			widget.NewLabelWithStyle("Validator Info", fyne.TextAlignCenter,
	// 				fyne.TextStyle{Monospace: true, Bold: true, Symbol: true}),
	// 		),
	// 		validatorAddressBox,
	// 		validatorStatusBox,
	// 		nodeIDBox,
	// 		topBox,
	// 		latestBlockBox,
	// 		lastProducedBox,
	// 		missChanceBox,
	// 		missChanceConfidenceBox,
	// 		startHeightBox,
	// 		producedBlocksBox,
	// 		monikerBox,
	// 		genesisChecksumBox,
	// 		chainIDBox,
	// 		dateBox,
	// 		streakBox,

	// 		// container.N
	// 	),
	// )
	// validatorsRightPart := container.NewVBox(widget.NewLabel("Vali\tdato\nrs"), activeValidatorsBox, pausedValidatorsBox, inactiveValidatorsBox, jailedValidatorsBox, waitingValidatorsBox)

	// return container.NewBorder(nil, refreshButton, nil, validatorsRightPart, mainInfo)
	return container.NewBorder(nil, refreshButton, nil, nil, mainInfo)
}
