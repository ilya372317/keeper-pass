package ui

import (
	"context"
	"fmt"

	"github.com/ilya372317/pass-keeper/internal/client/domain"
	"github.com/rivo/tview"
)

const (
	mainPageName     = "index"
	errModalPageName = "error_modal"
)

type passKeeperService interface {
	Login(context.Context, string, string) error
	All(context.Context) ([]domain.ShortData, error)
	Register(context.Context, string, string) error
	SaveLogin(context.Context, string, string, string) error
	SaveCard(context.Context, string, string, string, string) error
	SaveText(context.Context, string, string) error
	SaveBinary(context.Context, string, string) error
	Delete(context.Context, string) error
	ShowLoginPass(context.Context, int) (*domain.LoginPass, error)
	ShowCreditCard(context.Context, int) (*domain.CreditCard, error)
	ShowText(context.Context, int) (*domain.Text, error)
	ShowBinary(context.Context, int) (*domain.Binary, error)
}

type UserInterface struct {
	passKeeperService passKeeperService
	app               *tview.Application
	pages             *tview.Pages
}

func New(service passKeeperService) *UserInterface {
	pages := tview.NewPages()
	app := tview.NewApplication().SetRoot(pages, true).EnableMouse(true)

	return &UserInterface{
		passKeeperService: service,
		app:               app,
		pages:             pages,
	}
}

func (ui *UserInterface) Run(ctx context.Context) error {
	allData, err := ui.passKeeperService.All(ctx)
	if err != nil {
		return fmt.Errorf("failed get data for make list: %w", err)
	}
	ui.pages.AddPage(mainPageName, ui.buildMainPageWidget(ctx, allData), true, true)

	if err = ui.app.Run(); err != nil {
		return fmt.Errorf("failed on run ui app: %w", err)
	}

	return nil
}

func (ui *UserInterface) buildMainPageWidget(ctx context.Context, data []domain.ShortData) *tview.List {
	list := tview.NewList()
	for _, d := range data {
		list.InsertItem(
			int(d.ID),
			fmt.Sprintf("Info: %s", d.Info),
			fmt.Sprintf("Type: %s", d.StringKind()),
			'0',
			ui.handleShowClickFunc(ctx, d))
	}

	return list
}

func getItemPageName(itemID int64) string {
	return fmt.Sprintf("item-%d", itemID)
}

func (ui *UserInterface) handleShowClickFunc(ctx context.Context, d domain.ShortData) func() {
	return func() {
		form, err := ui.buildShowDataForm(ctx, d)
		if err != nil {
			ui.handleError(err, mainPageName)
			return
		}
		showPageName := getItemPageName(d.ID)
		ui.pages.AddPage(showPageName, form, true, true)
		ui.pages.SwitchToPage(getItemPageName(d.ID))
	}
}

func (ui *UserInterface) buildShowDataForm(ctx context.Context, d domain.ShortData) (*tview.Form, error) {
	switch d.Kind {
	case domain.KindLoginPass:
		lp, err := ui.passKeeperService.ShowLoginPass(ctx, int(d.ID))
		if err != nil {
			return nil, fmt.Errorf("failed show login pass page: %w", err)
		}
		return ui.buildShowLoginPassForm(lp), nil
	case domain.KindCreditCard:
	case domain.KindText:
	case domain.KindBinary:
	}
	return nil, fmt.Errorf("ui can`t show data of type: %s", d.StringKind())
}

func (ui *UserInterface) buildShowLoginPassForm(lp *domain.LoginPass) *tview.Form {
	form := tview.NewForm()
	form.AddInputField("login: ", lp.Login, 0, nil, nil)
	form.AddInputField("password: ", lp.Password, 0, nil, nil)
	form.AddInputField("URL: ", lp.URL, 0, nil, nil)
	form.AddButton("update", func() {
		ui.pages.SwitchToPage(mainPageName)
		ui.pages.RemovePage(getItemPageName(int64(lp.ID)))
	})
	form.AddButton("back", func() {
		ui.pages.SwitchToPage(mainPageName)
		ui.pages.RemovePage(getItemPageName(int64(lp.ID)))
	})
	form.
		SetBorder(true).
		SetTitle(fmt.Sprintf("login pass data #%d", lp.ID)).
		SetTitleAlign(tview.AlignCenter)
	return form
}

func (ui *UserInterface) handleError(err error, backPage string) {
	ui.
		pages.
		AddPage(
			errModalPageName,
			ui.buildErrorModelPage(err.Error(), backPage),
			true,
			true,
		)
	ui.pages.SwitchToPage(errModalPageName)
}

func (ui *UserInterface) buildErrorModelPage(errMsg, backPage string) *tview.Modal {
	return tview.NewModal().SetText(errMsg).AddButtons([]string{"back"}).SetDoneFunc(
		func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 {
				ui.pages.SwitchToPage(backPage)
				ui.pages.RemovePage(errModalPageName)
			}
		})
}
