package gui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	dialogWizard "github.com/KiraCore/kensho/gui/dialogs"
)

func showSekaiExecuteDialog(g *Gui) {
	var wizard *dialogWizard.Wizard

	cmdEntry := widget.NewEntry()
	cmdEntry.MultiLine = true
	cmdEntry.Wrapping = fyne.TextWrapWord

	doneAction := binding.NewDataListener(func() {
		g.WaitDialog.ShowWaitDialog()
		log.Printf("Trying to execute: %v", cmdEntry.Text)

		// execute submit
		g.WaitDialog.HideWaitDialog()
		wizard.Hide()
	})

	submitButton := widget.NewButton("Submit", func() {
		log.Printf("Submitting sekai cmd: %v", cmdEntry.Text)
		warningMessage := fmt.Sprintf("Are you sure you want to execute this?\n\nCommand: <%v>\n\nYou cannot revert changes", cmdEntry.Text)
		showWarningMessageWithConfirmation(g, warningMessage, doneAction)

	})
	closeButton := widget.NewButton("Cancel", func() {
		wizard.Hide()
	})
	submitButton.Importance = widget.HighImportance
	content := container.NewBorder(
		widget.NewLabel("Enter sekai command below"),
		container.NewVBox(submitButton, closeButton),
		nil, nil,
		cmdEntry,
	)

	wizard = dialogWizard.NewWizard("Sekai cmd executor", content)
	wizard.Show(g.Window)
	wizard.Resize(fyne.NewSize(900, 200))
}
