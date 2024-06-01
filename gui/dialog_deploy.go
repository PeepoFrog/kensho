package gui

import (
	"encoding/json"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	dialogWizard "github.com/KiraCore/kensho/gui/dialogs"
	"github.com/KiraCore/kensho/helper/httph"
	"github.com/KiraCore/kensho/types"
)

func showDeployDialog(g *Gui, doneListener binding.DataListener) {
	var wizard *dialogWizard.Wizard

	ipToJoinEntry := widget.NewEntry()

	interxPortToJoinEntry := widget.NewEntry()
	interxPortToJoinEntry.SetPlaceHolder("11000")

	sekaiRPCPortToJoinEntry := widget.NewEntry()
	sekaiRPCPortToJoinEntry.SetPlaceHolder("26657")

	sekaiP2PPortEntry := widget.NewEntry()
	sekaiP2PPortEntry.SetPlaceHolder("26656")

	localCheckBinding := binding.NewBool()
	localCheck := widget.NewCheckWithData("local", localCheckBinding)

	sudoPasswordBinding := binding.NewString()
	mnemonicBinding := binding.NewString()
	sudoCheck := binding.NewBool()
	mnemonicCheck := binding.NewBool()

	sudoPasswordEntryButton := widget.NewButtonWithIcon("Password (sudo)", theme.CancelIcon(), func() {
		showSudoEnteringDialog(g, sudoPasswordBinding, sudoCheck)
	})

	doneMnemonicDataListener := binding.NewDataListener(func() {
		mnemonicCheck.Set(true)
		confirmedMnemonic, _ := mnemonicBinding.Get()
		log.Println("Confirmed mnemonic:", confirmedMnemonic)
	})
	mnemonicManagerDialogButton := widget.NewButtonWithIcon("Mnemonic", theme.CancelIcon(), func() {
		showMnemonicManagerDialog(g, mnemonicBinding, doneMnemonicDataListener)
	})

	constructJoinCmd := func() (*types.RequestDeployPayload, error) {
		rpcPort := sekaiRPCPortToJoinEntry.Text
		if rpcPort == "" {
			rpcPort = sekaiRPCPortToJoinEntry.PlaceHolder
		} else {
			validate := httph.ValidatePortRange(rpcPort)
			if !validate {
				sekaiP2PPortEntry.SetValidationError(fmt.Errorf("invalid port"))
				return nil, fmt.Errorf("RPC port is not valid")
			}
		}
		p2pPort := sekaiP2PPortEntry.Text
		if p2pPort == "" {
			p2pPort = sekaiP2PPortEntry.PlaceHolder
		} else {
			validate := httph.ValidatePortRange(p2pPort)
			if !validate {
				return nil, fmt.Errorf("P2P port is not valid")
			}
		}
		interxPort := interxPortToJoinEntry.Text
		if interxPort == "" {
			interxPort = interxPortToJoinEntry.PlaceHolder
		} else {
			validate := httph.ValidatePortRange(rpcPort)
			if !validate {
				return nil, fmt.Errorf("interx port is not valid")
			}
		}

		ip := ipToJoinEntry.Text
		validate := httph.ValidateIP(ip)
		if !validate {
			return nil, fmt.Errorf(`ip <%v> is not valid`, ip)
		}

		mnemonic, _ := mnemonicBinding.Get()

		lCheck, _ := localCheckBinding.Get()

		payload := &types.RequestDeployPayload{
			Command: "join",
			Args: types.Args{
				IP:         ip,
				InterxPort: interxPort,
				RPCPort:    rpcPort,
				P2PPort:    p2pPort,
				Mnemonic:   mnemonic,
				Local:      lCheck,
			},
		}

		return payload, nil
	}

	deployErrorBinding := binding.NewBool()
	errorMessageBinding := binding.NewString()

	deployButton := widget.NewButton("Deploy", func() {
		payload, err := constructJoinCmd()
		if err != nil {
			g.showErrorDialog(err, binding.NewDataListener(func() {}))
			return
		}

		sP, _ := sudoPasswordBinding.Get()

		cmdForDeploy := fmt.Sprintf(`echo '%v' | sudo -S sh -c "$(curl -s --show-error --fail %v 2>&1)"`, sP, types.BOOTSTRAP_SCRIPT)
		// cmdForDeploy = fmt.Sprintf(`echo %v`, sP) // for quick testing
		showCmdExecDialogAndRunCmdV4(g, "Deploying", cmdForDeploy, true, deployErrorBinding, errorMessageBinding)

		errB, _ := deployErrorBinding.Get()
		if errB {
			errMsg, _ := errorMessageBinding.Get()
			g.showErrorDialog(fmt.Errorf("error while checking the sudo password: %v ", errMsg), binding.NewDataListener(func() {}))
			return
		}

		g.WaitDialog.ShowWaitDialog()
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Failed to marshal JSON payload: %v", err)
			return
		}
		out, err := httph.ExecHttpRequestBySSHTunnel(g.sshClient, types.SEKIN_EXECUTE_ENDPOINT, "POST", jsonPayload)
		log.Printf("Executing http payload for join: %+v", payload)
		log.Printf("log: %v\nerr: %v", string(out), err)
		g.WaitDialog.HideWaitDialog()

		if err != nil {
			g.showErrorDialog(err, binding.NewDataListener(func() {}))
			return
		}

		errB, _ = deployErrorBinding.Get()
		if errB {
			errMsg, _ := errorMessageBinding.Get()
			g.showErrorDialog(fmt.Errorf("error when executing join command: %v ", errMsg), binding.NewDataListener(func() {}))
			return
		}

		doneListener.DataChanged()
		wizard.Hide()

	})

	deployButton.Disable()

	deployActivatorDataListener := binding.NewDataListener(func() {
		sCheck, _ := sudoCheck.Get()
		if sCheck {
			sudoPasswordEntryButton.Icon = theme.ConfirmIcon()
			sudoPasswordEntryButton.Refresh()
		} else {
			sudoPasswordEntryButton.Icon = theme.CancelIcon()
			sudoPasswordEntryButton.Refresh()
		}
		mCheck, _ := mnemonicCheck.Get()

		if mCheck {
			mnemonicManagerDialogButton.Icon = theme.ConfirmIcon()
			mnemonicManagerDialogButton.Refresh()
		} else {
			mnemonicManagerDialogButton.Icon = theme.CancelIcon()
			mnemonicManagerDialogButton.Refresh()
		}

		if sCheck && mCheck {
			deployButton.Enable()
		} else {
			if !deployButton.Disabled() {
				deployButton.Disable()
			}
		}
	})

	closeButton := widget.NewButton("Close", func() {
		wizard.Hide()
	})

	mnemonicCheck.AddListener(deployActivatorDataListener)
	sudoCheck.AddListener(deployActivatorDataListener)

	content := container.NewVBox(
		widget.NewLabel("Trusted IP address"),
		ipToJoinEntry,
		localCheck,
		widget.NewLabel("RPC Port"),
		sekaiRPCPortToJoinEntry,
		widget.NewLabel("P2P Port"),
		sekaiP2PPortEntry,
		widget.NewLabel("Interx Port"),
		interxPortToJoinEntry,
		sudoPasswordEntryButton,
		mnemonicManagerDialogButton,
		deployButton,
		closeButton,
	)

	wizard = dialogWizard.NewWizard("Connect", content)
	wizard.Show(g.Window)
	wizard.Resize(fyne.NewSize(400, 500))

}

func showSudoEnteringDialog(g *Gui, bindString binding.String, bindCheck binding.Bool) {
	var wizard *dialogWizard.Wizard

	sudoPasswordEntry := widget.NewEntryWithData(bindString)
	errorMessageBinding := binding.NewString()
	checkSudoPassword := func(p string) error {
		cmd := fmt.Sprintf("echo '%v' | sudo -S uname", p)
		errB := binding.NewBool()
		showCmdExecDialogAndRunCmdV4(g, "checking sudo password", cmd, true, errB, errorMessageBinding)
		errExec, _ := errB.Get()
		if errExec {
			errMsg, err := errorMessageBinding.Get()
			if err != nil {
				return err
			}
			return fmt.Errorf("error while checking the sudo password: %v ", errMsg)
		}
		return nil
	}

	okButton := widget.NewButton("Ok", func() {
		err := checkSudoPassword(sudoPasswordEntry.Text)
		if err == nil {
			bindCheck.Set(true)
			wizard.Hide()
		} else {
			bindCheck.Set(false)
			sudoPasswordEntry.SetValidationError(fmt.Errorf("sudo password is wrong: %w", err))
			showInfoDialog(g, "ERROR", fmt.Sprintf("error when checking sudo password: %v", err.Error()))
		}

	})
	cancelButton := widget.NewButton("Cancel", func() { wizard.Hide() })
	content := container.NewVBox(
		sudoPasswordEntry,
		container.NewHBox(
			okButton, container.NewCenter(), cancelButton,
		),
	)

	wizard = dialogWizard.NewWizard("Enter your sudo password", content)
	wizard.Show(g.Window)

}
