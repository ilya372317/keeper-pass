package ui

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ilya372317/pass-keeper/internal/client/domain"
	"github.com/rivo/tview"
)

const (
	mainPageName     = "index"
	errModalPageName = "error_modal"

	shortWidth = 5
)

type passKeeperService interface {
	All(context.Context) ([]domain.ShortData, error)
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
	app := tview.NewApplication().SetRoot(pages, true).EnableMouse(true).EnablePaste(true)

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
		cc, err := ui.passKeeperService.ShowCreditCard(ctx, int(d.ID))
		if err != nil {
			return nil, fmt.Errorf("failed show credit card page: %w", err)
		}
		return ui.buildShowCreditCardForm(cc), nil
	case domain.KindText:
	case domain.KindBinary:
	}
	return nil, fmt.Errorf("ui can`t show data of type: %s", d.StringKind())
}

func (ui *UserInterface) buildShowCreditCardForm(cc *domain.CreditCard) *tview.Form {
	const expItemCount = 2
	const mothIndex = 0
	const yearIndex = 1
	expItems := strings.Split(cc.Exp, "/")
	var expMonth, expYear string
	if len(expItems) >= expItemCount {
		expMonth, expYear = expItems[mothIndex], expItems[yearIndex]
	}
	form := tview.NewForm()
	form.SetBorder(true).
		SetTitle(fmt.Sprintf("credit card data #%d", cc.ID)).
		SetTitleAlign(tview.AlignCenter)
	form.AddInputField("number: ", cc.CardNumber, 0, nil, nil)
	form.AddInputField("mm: ", expMonth, shortWidth, nil, nil)
	form.AddInputField("yy: ", expYear, shortWidth, nil, nil)
	form.AddInputField("CVV: ", strconv.Itoa(cc.Code), 0, nil, nil)
	form.AddInputField("bank name: ", cc.BankName, 0, nil, nil)
	form.AddButton("update", func() {
	})
	form.AddButton("back", func() {
		ui.pages.SwitchToPage(mainPageName)
		ui.pages.RemovePage(getItemPageName(int64(cc.ID)))
	})

	return form
}

func (ui *UserInterface) buildShowLoginPassForm(lp *domain.LoginPass) *tview.Form {
	form := tview.NewForm()
	form.AddInputField("login: ", lp.Login, 0, nil, nil)
	form.AddInputField("password: ", lp.Password, 0, nil, nil)
	form.AddInputField("URL: ", lp.URL, 0, nil, nil)
	form.AddButton("update", func() {
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
