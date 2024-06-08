package gui

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

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

type nodeInfoScreen struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	executingStatus bool
}

func makeNodeInfoScreen(_ fyne.Window, g *Gui) fyne.CanvasObject {

	// nodeInfoTab := container.NewTabItem("Node Info", makeNodeInfoTab(g))
	// validatorInfoTab := container.NewTabItem("Validator Info", makeValidatorInfoTab(g))
	// return container.NewAppTabs(nodeInfoTab, validatorInfoTab)

	return makeNodeInfoTab(g)
}

func makeNodeInfoTab(g *Gui) fyne.CanvasObject {
	// TODO: only for testing, delete later
	// g.Host.IP = "148.251.69.56"
	ctx := context.Background()
	g.NodeInfo.ctx, g.NodeInfo.ctxCancel = context.WithCancel(ctx)

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
	// latestBlockBox := container.NewHBox(
	// 	widget.NewLabel("Block:"), latestBlockLabel,
	// )

	// validator address box
	validatorAddressData := binding.NewString()
	validatorAddressLabel := widget.NewLabelWithData(validatorAddressData)
	validatorAddressCopyButton := widget.NewButtonWithIcon("Copy", theme.ContentCopyIcon(), func() {
		data, err := validatorAddressData.Get()
		if err != nil {
			log.Println(err)
			return
		}

		err = clipboard.WriteAll(data)
		if err != nil {
			return
		}
	})

	//validator status (active, paused, etc...)
	validatorStatusData := binding.NewString()
	validatorStatusLabel := widget.NewLabelWithData(validatorStatusData)
	// validatorStatusBox := container.NewHBox(
	// 	widget.NewLabel("Val.Status:\t\t"), validatorStatusLabel,
	// 	validatorControlButton,
	// )
	// nodeID
	nodeIDData := binding.NewString()
	nodeIDLabel := widget.NewLabelWithData(nodeIDData)
	// nodeIDBox := container.NewHBox(
	// 	widget.NewLabel("Val.NodeID:\t\t"), nodeIDLabel,
	// )

	topData := binding.NewString()
	topLabel := widget.NewLabelWithData(topData)
	// topBox := container.NewHBox(
	// 	widget.NewLabel("Rank:\t\t"), topLabel,
	// )

	// miss chance box
	missChanceData := binding.NewString()
	missChanceLabel := widget.NewLabelWithData(missChanceData)
	// missChanceBox := container.NewHBox(
	// 	widget.NewLabel("Miss:\t\t"), missChanceLabel,
	// )

	lastProducedBlockData := binding.NewString()
	lastProducedLabel := widget.NewLabelWithData(lastProducedBlockData)
	// lastProducedBox := container.NewHBox(
	// 	widget.NewLabel("Latest Block:\t\t"), lastProducedLabel,
	// )

	// miss Chance Confidence box
	missChanceConfidenceData := binding.NewString()
	missChanceConfidenceLabel := widget.NewLabelWithData(missChanceConfidenceData)
	// missChanceConfidenceBox := container.NewHBox(
	// 	widget.NewLabel("Miss conf.\t\t"), missChanceConfidenceLabel,
	// )

	startHeightData := binding.NewString()
	startHeightLabel := widget.NewLabelWithData(startHeightData)
	// startHeightBox := container.NewHBox(
	// 	widget.NewLabel("Start Height:\t\t"), startHeightLabel,
	// )

	producedBlocksData := binding.NewString()
	producedBlocksLabel := widget.NewLabelWithData(producedBlocksData)
	// producedBlocksBox := container.NewHBox(
	// 	widget.NewLabel("Produced:\t\t"), producedBlocksLabel,
	// )

	monikerData := binding.NewString()
	monikerLabel := widget.NewLabelWithData(monikerData)
	// monikerBox := container.NewHBox(
	// 	widget.NewLabel("Moniker:\t\t"), monikerLabel,
	// )

	genesisChecksumData := binding.NewString()
	genesisChecksumLabel := widget.NewLabelWithData(genesisChecksumData)
	// genesisChecksumBox := container.NewHBox(
	// 	widget.NewLabel("Gen. SHA256::\t\t"), genesisChecksumLabel,
	// )

	nodeCatchingData := binding.NewString()
	nodeCatchingLabel := widget.NewLabelWithData(nodeCatchingData)

	// Numerical validators status
	activeValidatorsData := binding.NewInt()
	activeValidatorsLabel := widget.NewLabelWithData(binding.IntToString(activeValidatorsData))
	activeValidatorsLabel.Alignment = fyne.TextAlignCenter
	activeValidatorsInfoText := widget.NewLabel("Active:")
	activeValidatorsInfoText.TextStyle.Bold = true
	activeValidatorsBox := container.NewVBox(
		activeValidatorsInfoText, activeValidatorsLabel,
	)

	pausedValidatorsData := binding.NewInt()
	pausedValidatorsLabel := widget.NewLabelWithData(binding.IntToString(pausedValidatorsData))
	pausedValidatorsLabel.Alignment = fyne.TextAlignCenter
	pausedValidatorsInfoText := widget.NewLabel("Paused:")
	pausedValidatorsInfoText.TextStyle.Bold = true
	pausedValidatorsBox := container.NewVBox(
		pausedValidatorsInfoText, pausedValidatorsLabel,
	)

	inactiveValidatorsData := binding.NewInt()
	inactiveValidatorsLabel := widget.NewLabelWithData(binding.IntToString(inactiveValidatorsData))
	inactiveValidatorsLabel.Alignment = fyne.TextAlignCenter
	inactiveValidatorsInfoText := widget.NewLabel("Inactive:")
	inactiveValidatorsInfoText.TextStyle.Bold = true

	inactiveValidatorsBox := container.NewVBox(
		inactiveValidatorsInfoText, inactiveValidatorsLabel,
	)

	jailedValidatorsData := binding.NewInt()
	jailedValidatorsLabel := widget.NewLabelWithData(binding.IntToString(jailedValidatorsData))
	jailedValidatorsLabel.Alignment = fyne.TextAlignCenter

	jailedValidatorsInfoText := widget.NewLabel("Jailed:")
	jailedValidatorsInfoText.TextStyle.Bold = true

	jailedValidatorsBox := container.NewVBox(
		jailedValidatorsInfoText, jailedValidatorsLabel,
	)

	waitingValidatorsData := binding.NewInt()
	waitingValidatorsLabel := widget.NewLabelWithData(binding.IntToString(waitingValidatorsData))
	waitingValidatorsLabel.Alignment = fyne.TextAlignCenter

	waitingValidatorsInfoText := widget.NewLabel("Waiting:")
	waitingValidatorsInfoText.TextStyle.Bold = true

	waitingValidatorsBox := container.NewVBox(
		waitingValidatorsInfoText, waitingValidatorsLabel,
	)

	chainIDData := binding.NewString()
	chainIDLabel := widget.NewLabelWithData(chainIDData)

	// dateData := binding.NewString()
	// dateLabel := widget.NewLabelWithData(dateData)
	// dateBox := container.NewHBox(
	// 	widget.NewLabel("Date:\t\t"), dateLabel,
	// )

	streakData := binding.NewString()
	streakLabel := widget.NewLabelWithData(streakData)

	//
	valuesForm := widget.NewForm(
		widget.NewFormItem("ChainID:", chainIDLabel),
		widget.NewFormItem("Moniker:", monikerLabel),
		widget.NewFormItem("Val.Status:", validatorStatusLabel),
		widget.NewFormItem("Node Status:", nodeCatchingLabel),
		widget.NewFormItem("Block:", latestBlockLabel),
		widget.NewFormItem("Latest Block:", lastProducedLabel),
		widget.NewFormItem("Produced:", producedBlocksLabel),
		widget.NewFormItem("Streak:", streakLabel),
		widget.NewFormItem("Rank:", topLabel),
		widget.NewFormItem("Miss:", missChanceLabel),
		widget.NewFormItem("Miss conf.", missChanceConfidenceLabel),
		widget.NewFormItem("Start Height:", startHeightLabel),
		widget.NewFormItem("Val.Addr:", container.NewHBox(validatorAddressLabel, validatorAddressCopyButton)),
		widget.NewFormItem("Val.NodeID:", nodeIDLabel),
		widget.NewFormItem("Gen.SHA256:", genesisChecksumLabel),
	)

	loadingData := binding.NewFloat()

	loadMessage := binding.NewString()
	txExecLoadingWidget := container.NewHBox(
		widget.NewProgressBarWithData(loadingData), widget.NewLabelWithData(loadMessage),
	)

	startTXexec := func(command string) {
		txExecLoadingWidget.Show()
		g.NodeInfo.executingStatus = true

		loadMessage.Set(fmt.Sprintf("Executing <%v> command", command))
		validatorControlButton.Disable()
		maxRange := 40
		for i := range maxRange {
			time.Sleep(time.Second * time.Duration(i))
			percentage := (i * 100) / maxRange
			loadingData.Set(float64(percentage))

		}
		txExecLoadingWidget.Hide()
		g.NodeInfo.executingStatus = false
		refreshBinding.DataChanged()
	}

	execFunc := func(args types.ExecSekaiCmd) {
		go startTXexec(string(args.TX))
		g.WaitDialog.ShowWaitDialog()
		payload, err := json.Marshal(args)

		log.Printf("Executing: %v", args)
		if err != nil {
			g.showErrorDialog(err, binding.NewDataListener(func() {}))
		}
		httph.ExecHttpRequestBySSHTunnel(g.sshClient, types.SEKIN_EXECUTE_CMD_ENDPOINT, "POST", payload)
		g.WaitDialog.HideWaitDialog()

		refreshBinding.DataChanged()
	}

	monikerEntryData := binding.NewString()
	claimDataListener := binding.NewDataListener(func() {
		moniker, _ := monikerEntryData.Get()
		execFunc(types.ExecSekaiCmd{TX: types.ClaimValidatorSeat, Moniker: moniker})
	})

	claimSeatButton.OnTapped = func() {
		//claim seat
		showMonikerEntryDialog(g, monikerEntryData, claimDataListener)
	}
	pauseValidatorFunc := func() {
		// pause
		execFunc(types.ExecSekaiCmd{TX: types.Pause})
	}
	unpauseValidatorFunc := func() {
		// unpause tx
		execFunc(types.ExecSekaiCmd{TX: types.Unpause})
	}
	activateValidatorFunc := func() {
		// activate
		execFunc(types.ExecSekaiCmd{TX: types.Activate})
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

		if dashboardData.ValidatorStatus != "Unknown" && !g.NodeInfo.executingStatus {
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

		if dashboardData.CatchingUp {
			nodeCatchingData.Set("Syncing")
		} else {
			nodeCatchingData.Set("Running")
		}
		// dashboardData.CatchingUp,
		// widget.NewLabel()

	}
	refreshBinding = binding.NewDataListener(func() {
		refreshScreen()
		err, _ := errBinding.Get()
		if err != nil {
			if e, ok := err.(error); ok {
				// g.showErrorDialog(e, binding.NewDataListener(func() {}))
				log.Printf("Refreshing unsuccessful, reason: %v", e.Error())
			}
		}
	})
	go func() {
		refreshTime := 20 * time.Second
		log.Printf("Starting goroutine with refresh rate %v", refreshTime)
		timer := time.NewTimer(refreshTime)
		defer timer.Stop() // Clean up the timer when the goroutine ends

		refreshBinding.DataChanged()
		for {
			select {
			case <-g.NodeInfo.ctx.Done():
				log.Printf("Ending nodeInfo refresh goroutine")
				return
			case <-timer.C:
				refreshBinding.DataChanged()
				timer.Reset(refreshTime) // Reset the timer
			}
		}
	}()

	// refreshButton := widget.NewButton("Refresh", refreshBinding.DataChanged)
	// sendSekaiCommandButton := widget.NewButton("Execute sekai command", func() { showSekaiExecuteDialog(g) })
	mainInfo := container.NewVScroll(
		container.NewVBox(
			widget.NewSeparator(),
			valuesForm,
			validatorControlButton,
		),
	)

	validatorsTopPart := container.NewHBox(activeValidatorsBox, pausedValidatorsBox, inactiveValidatorsBox, jailedValidatorsBox, waitingValidatorsBox)

	// return container.NewBorder(nil, refreshButton, nil, validatorsRightPart, mainInfo)
	return container.NewBorder(container.NewCenter(validatorsTopPart), nil, nil, nil, mainInfo)
}
