package cntdialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/containers"
	"github.com/containers/podman-tui/pdcs/images"
	"github.com/containers/podman-tui/pdcs/networks"
	"github.com/containers/podman-tui/pdcs/pods"
	"github.com/containers/podman-tui/pdcs/volumes"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	containerCreateDialogMaxWidth = 100
	containerCreateDialogHeight   = 18
)

const (
	createContainerFormFocus = 0 + iota
	createCategoriesFocus
	createCategoryPagesFocus
	createContainerNameFieldFocus
	createContainerImageFieldFocus
	createcontainerPodFieldFocis
	createContainerLabelsFieldFocus
	createContainerRemoveFieldFocus
	createContainerSelinuxLabelFieldFocus
	createContainerApprarmorFieldFocus
	createContainerSeccompFeildFocus
	createContainerMaskFieldFocus
	createContainerUnmaskFieldFocus
	createContainerNoNewPrivFieldFocus
	createContainerPortExposeFieldFocus
	createContainerPortPublishFieldFocus
	createContainerPortPublishAllFieldFocus
	createContainerHostnameFieldFocus
	createContainerIPAddrFieldFocus
	createContainerMacAddrFieldFocus
	createContainerNetworkFieldFocus
	createContainerDNSServersFieldFocus
	createContainerDNSOptionsFieldFocus
	createContainerDNSSearchFieldFocus
	createContainerImageVolumeFieldFocus
	createContainerVolumeFieldFocus
	containerHealthCmdFieldFocus
	containerHealthStartupCmdFieldFocus
	containerHealthOnFailureFieldFocus
	containerHealthIntervalFieldFocus
	containerHealthStartupIntervalFieldFocus
	containerHealthTimeoutFieldFocus
	containerHealthStartupTimeoutFieldFocus
	containerHealthRetriesFieldFocus
	containerHealthStartupRetriesFieldFocus
	containerHealthStartPeriodFieldFocus
	containerHealthStartupSuccessFieldFocus
)

const (
	basicInfoPageIndex = 0 + iota
	dnsPageIndex
	healthPageIndex
	networkingPageIndex
	portPageIndex
	securityOptsPageIndex
	volumePageIndex
)

// ContainerCreateDialog implements container create dialog
type ContainerCreateDialog struct {
	*tview.Box
	layout                              *tview.Flex
	categoryLabels                      []string
	categories                          *tview.TextView
	categoryPages                       *tview.Pages
	basicInfoPage                       *tview.Flex
	securityOptsPage                    *tview.Flex
	portPage                            *tview.Flex
	networkingPage                      *tview.Flex
	dnsPage                             *tview.Flex
	volumePage                          *tview.Flex
	healthPage                          *tview.Flex
	form                                *tview.Form
	display                             bool
	activePageIndex                     int
	focusElement                        int
	imageList                           []images.ImageListReporter
	podList                             []*entities.ListPodsReport
	containerNameField                  *tview.InputField
	containerImageField                 *tview.DropDown
	containerPodField                   *tview.DropDown
	containerLabelsField                *tview.InputField
	containerRemoveField                *tview.Checkbox
	containerSelinuxLabelField          *tview.InputField
	containerApparmorField              *tview.InputField
	containerSeccompField               *tview.InputField
	containerMaskField                  *tview.InputField
	containerUnmaskField                *tview.InputField
	containerNoNewPrivField             *tview.Checkbox
	containerPortExposeField            *tview.InputField
	containerPortPublishField           *tview.InputField
	ContainerPortPublishAllField        *tview.Checkbox
	containerHostnameField              *tview.InputField
	containerIPAddrField                *tview.InputField
	containerMacAddrField               *tview.InputField
	containerNetworkField               *tview.DropDown
	containerDNSServersField            *tview.InputField
	containerDNSOptionsField            *tview.InputField
	containerDNSSearchField             *tview.InputField
	containerHealthCmdField             *tview.InputField
	containerHealthIntervalField        *tview.InputField
	containerHealthOnFailureField       *tview.DropDown
	containerHealthRetriesField         *tview.InputField
	containerHealthStartPeriodField     *tview.InputField
	containerHealthTimeoutField         *tview.InputField
	containerHealthStartupCmdField      *tview.InputField
	containerHealthStartupIntervalField *tview.InputField
	containerHealthStartupRetriesField  *tview.InputField
	containerHealthStartupSuccessField  *tview.InputField
	containerHealthStartupTimeoutField  *tview.InputField
	containerVolumeField                *tview.DropDown
	containerImageVolumeField           *tview.DropDown
	cancelHandler                       func()
	createHandler                       func()
}

// NewContainerCreateDialog returns new container create dialog primitive ContainerCreateDialog.
func NewContainerCreateDialog() *ContainerCreateDialog {
	containerDialog := ContainerCreateDialog{
		Box:              tview.NewBox(),
		layout:           tview.NewFlex().SetDirection(tview.FlexRow),
		categories:       tview.NewTextView(),
		categoryPages:    tview.NewPages(),
		basicInfoPage:    tview.NewFlex(),
		securityOptsPage: tview.NewFlex(),
		networkingPage:   tview.NewFlex(),
		dnsPage:          tview.NewFlex(),
		portPage:         tview.NewFlex(),
		volumePage:       tview.NewFlex(),
		healthPage:       tview.NewFlex(),
		form:             tview.NewForm(),
		categoryLabels: []string{
			"Basic Information",
			"DNS Settings",
			"Health check",
			"Network Settings",
			"Ports Settings",
			"Security Options",
			"Volumes Settings"},
		activePageIndex:                     0,
		display:                             false,
		containerNameField:                  tview.NewInputField(),
		containerImageField:                 tview.NewDropDown(),
		containerPodField:                   tview.NewDropDown(),
		containerLabelsField:                tview.NewInputField(),
		containerRemoveField:                tview.NewCheckbox(),
		containerSelinuxLabelField:          tview.NewInputField(),
		containerApparmorField:              tview.NewInputField(),
		containerSeccompField:               tview.NewInputField(),
		containerMaskField:                  tview.NewInputField(),
		containerUnmaskField:                tview.NewInputField(),
		containerNoNewPrivField:             tview.NewCheckbox(),
		containerPortExposeField:            tview.NewInputField(),
		containerPortPublishField:           tview.NewInputField(),
		ContainerPortPublishAllField:        tview.NewCheckbox(),
		containerHostnameField:              tview.NewInputField(),
		containerIPAddrField:                tview.NewInputField(),
		containerMacAddrField:               tview.NewInputField(),
		containerNetworkField:               tview.NewDropDown(),
		containerDNSServersField:            tview.NewInputField(),
		containerDNSOptionsField:            tview.NewInputField(),
		containerDNSSearchField:             tview.NewInputField(),
		containerVolumeField:                tview.NewDropDown(),
		containerImageVolumeField:           tview.NewDropDown(),
		containerHealthCmdField:             tview.NewInputField(),
		containerHealthIntervalField:        tview.NewInputField(),
		containerHealthOnFailureField:       tview.NewDropDown(),
		containerHealthRetriesField:         tview.NewInputField(),
		containerHealthStartPeriodField:     tview.NewInputField(),
		containerHealthTimeoutField:         tview.NewInputField(),
		containerHealthStartupCmdField:      tview.NewInputField(),
		containerHealthStartupIntervalField: tview.NewInputField(),
		containerHealthStartupRetriesField:  tview.NewInputField(),
		containerHealthStartupSuccessField:  tview.NewInputField(),
		containerHealthStartupTimeoutField:  tview.NewInputField(),
	}

	containerDialog.setupLayout()
	containerDialog.setActiveCategory(0)
	containerDialog.initCustomInputHanlers()

	return &containerDialog
}

func (d *ContainerCreateDialog) setupLayout() {
	bgColor := style.DialogBgColor

	d.categories.SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)
	d.categories.SetBackgroundColor(bgColor)
	d.categories.SetBorder(true)
	d.categories.SetBorderColor(style.DialogSubBoxBorderColor)

	// category pages
	d.categoryPages.SetBackgroundColor(bgColor)
	d.categoryPages.SetBorder(true)
	d.categoryPages.SetBorderColor(style.DialogSubBoxBorderColor)

	d.setupBasicInfoPageUI()
	d.setupDnsPageUI()
	d.setupHealthPageUI()
	d.setupNetworkPageUI()
	d.setupPortsPageUI()
	d.setupSecurityPageUI()
	d.setupVolumePageUI()

	// form
	d.form.SetBackgroundColor(bgColor)
	d.form.AddButton("Cancel", nil)
	d.form.AddButton("Create", nil)
	d.form.SetButtonsAlign(tview.AlignRight)
	d.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// adding category pages
	d.categoryPages.AddPage(d.categoryLabels[basicInfoPageIndex], d.basicInfoPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[dnsPageIndex], d.dnsPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[healthPageIndex], d.healthPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[networkingPageIndex], d.networkingPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[portPageIndex], d.portPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[securityOptsPageIndex], d.securityOptsPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[volumePageIndex], d.volumePage, true, true)

	// add it to layout.
	d.layout.SetBackgroundColor(bgColor)
	d.layout.SetBorder(true)
	d.layout.SetBorderColor(style.DialogBorderColor)
	d.layout.SetTitle("PODMAN CONTAINER CREATE")

	_, layoutWidth := utils.AlignStringListWidth(d.categoryLabels)
	layout := tview.NewFlex().SetDirection(tview.FlexColumn)

	layout.AddItem(d.categories, layoutWidth+6, 0, true)
	layout.AddItem(d.categoryPages, 0, 1, true)
	layout.SetBackgroundColor(bgColor)
	d.layout.AddItem(layout, 0, 1, true)

	d.layout.AddItem(d.form, dialogs.DialogFormHeight, 0, true)
}

func (d *ContainerCreateDialog) setupBasicInfoPageUI() {
	bgColor := style.DialogBgColor
	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected
	inputFieldBgColor := style.InputFieldBgColor
	basicInfoPageLabelWidth := 14

	// name field
	d.containerNameField.SetLabel("name:")
	d.containerNameField.SetLabelWidth(basicInfoPageLabelWidth)
	d.containerNameField.SetBackgroundColor(bgColor)
	d.containerNameField.SetLabelColor(style.DialogFgColor)
	d.containerNameField.SetFieldBackgroundColor(inputFieldBgColor)

	// image field
	d.containerImageField.SetLabel("select image:")
	d.containerImageField.SetLabelWidth(basicInfoPageLabelWidth)
	d.containerImageField.SetBackgroundColor(bgColor)
	d.containerImageField.SetLabelColor(style.DialogFgColor)
	d.containerImageField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	d.containerImageField.SetFieldBackgroundColor(inputFieldBgColor)

	// pod field
	d.containerPodField.SetLabel("select pod:")
	d.containerPodField.SetLabelWidth(basicInfoPageLabelWidth)
	d.containerPodField.SetBackgroundColor(bgColor)
	d.containerPodField.SetLabelColor(style.DialogFgColor)
	d.containerPodField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	d.containerPodField.SetFieldBackgroundColor(inputFieldBgColor)

	// labels field
	d.containerLabelsField.SetLabel("labels:")
	d.containerLabelsField.SetLabelWidth(basicInfoPageLabelWidth)
	d.containerLabelsField.SetBackgroundColor(bgColor)
	d.containerLabelsField.SetLabelColor(style.DialogFgColor)
	d.containerLabelsField.SetFieldBackgroundColor(inputFieldBgColor)

	// remove field
	d.containerRemoveField.SetLabel("remove container after exit ")
	d.containerRemoveField.SetBackgroundColor(bgColor)
	d.containerRemoveField.SetLabelColor(style.DialogFgColor)
	d.containerRemoveField.SetChecked(true)
	d.containerRemoveField.SetFieldBackgroundColor(inputFieldBgColor)

	d.basicInfoPage.SetDirection(tview.FlexRow)
	d.basicInfoPage.AddItem(d.containerNameField, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.containerImageField, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.containerPodField, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.containerLabelsField, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.containerRemoveField, 1, 0, true)
	d.basicInfoPage.SetBackgroundColor(bgColor)
}

func (d *ContainerCreateDialog) setupDnsPageUI() {
	bgColor := style.DialogBgColor
	inputFieldBgColor := style.InputFieldBgColor
	dnsPageLabelWidth := 13

	// hostname field
	d.containerDNSServersField.SetLabel("dns servers:")
	d.containerDNSServersField.SetLabelWidth(dnsPageLabelWidth)
	d.containerDNSServersField.SetBackgroundColor(bgColor)
	d.containerDNSServersField.SetLabelColor(style.DialogFgColor)
	d.containerDNSServersField.SetFieldBackgroundColor(inputFieldBgColor)

	// IP field
	d.containerDNSOptionsField.SetLabel("dns options:")
	d.containerDNSOptionsField.SetLabelWidth(dnsPageLabelWidth)
	d.containerDNSOptionsField.SetBackgroundColor(bgColor)
	d.containerDNSOptionsField.SetLabelColor(style.DialogFgColor)
	d.containerDNSOptionsField.SetFieldBackgroundColor(inputFieldBgColor)

	// mac field
	d.containerDNSSearchField.SetLabel("dns search:")
	d.containerDNSSearchField.SetLabelWidth(dnsPageLabelWidth)
	d.containerDNSSearchField.SetBackgroundColor(bgColor)
	d.containerDNSSearchField.SetLabelColor(style.DialogFgColor)
	d.containerDNSSearchField.SetFieldBackgroundColor(inputFieldBgColor)

	d.dnsPage.SetDirection(tview.FlexRow)
	d.dnsPage.AddItem(d.containerDNSServersField, 1, 0, true)
	d.dnsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.dnsPage.AddItem(d.containerDNSOptionsField, 1, 0, true)
	d.dnsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.dnsPage.AddItem(d.containerDNSSearchField, 1, 0, true)
	d.dnsPage.SetBackgroundColor(bgColor)
}

func (d *ContainerCreateDialog) setupHealthPageUI() {
	bgColor := style.DialogBgColor
	inputFieldBgColor := style.InputFieldBgColor
	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected
	healthPageLabelWidth := 13
	healthPageSecColLabelWidth := 18
	healthPageMultiRowFieldWidth := 7

	// health cmd
	d.containerHealthCmdField.SetLabel("Command:")
	d.containerHealthCmdField.SetLabelWidth(healthPageLabelWidth)
	d.containerHealthCmdField.SetBackgroundColor(bgColor)
	d.containerHealthCmdField.SetLabelColor(style.DialogFgColor)
	d.containerHealthCmdField.SetFieldBackgroundColor(inputFieldBgColor)

	// startup cmd
	d.containerHealthStartupCmdField.SetLabel("Startup cmd:")
	d.containerHealthStartupCmdField.SetLabelWidth(healthPageLabelWidth)
	d.containerHealthStartupCmdField.SetBackgroundColor(bgColor)
	d.containerHealthStartupCmdField.SetLabelColor(style.DialogFgColor)
	d.containerHealthStartupCmdField.SetFieldBackgroundColor(inputFieldBgColor)

	// multi primitive row01
	// startup success
	d.containerHealthStartupSuccessField.SetLabel("Startup success:")
	d.containerHealthStartupSuccessField.SetLabelWidth(healthPageSecColLabelWidth)
	d.containerHealthStartupSuccessField.SetBackgroundColor(bgColor)
	d.containerHealthStartupSuccessField.SetLabelColor(style.DialogFgColor)
	d.containerHealthStartupSuccessField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerHealthStartupSuccessField.SetFieldWidth(healthPageMultiRowFieldWidth)

	// on-failure
	d.containerHealthOnFailureField.SetOptions([]string{"none", "kill", "restart", "stop"}, nil)
	d.containerHealthOnFailureField.SetLabel("On failure:")
	d.containerHealthOnFailureField.SetLabelWidth(healthPageLabelWidth)
	d.containerHealthOnFailureField.SetBackgroundColor(bgColor)
	d.containerHealthOnFailureField.SetLabelColor(style.DialogFgColor)
	d.containerHealthOnFailureField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	d.containerHealthOnFailureField.SetFieldBackgroundColor(inputFieldBgColor)

	// start period
	startPeroidLabel := fmt.Sprintf("%15s: ", "Start period")
	d.containerHealthStartPeriodField.SetLabel(startPeroidLabel)
	d.containerHealthStartPeriodField.SetBackgroundColor(bgColor)
	d.containerHealthStartPeriodField.SetLabelColor(style.DialogFgColor)
	d.containerHealthStartPeriodField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerHealthStartPeriodField.SetFieldWidth(healthPageMultiRowFieldWidth)

	multiItemRow01 := tview.NewFlex().SetDirection(tview.FlexColumn)
	multiItemRow01.AddItem(d.containerHealthOnFailureField, 0, 1, true)
	multiItemRow01.AddItem(d.containerHealthStartupSuccessField, 0, 1, true)
	multiItemRow01.AddItem(d.containerHealthStartPeriodField, 0, 1, true)
	multiItemRow01.SetBackgroundColor(bgColor)

	// multi primitive row02
	// startup interval
	d.containerHealthStartupIntervalField.SetLabel("Startup interval:")
	d.containerHealthStartupIntervalField.SetLabelWidth(healthPageSecColLabelWidth)
	d.containerHealthStartupIntervalField.SetBackgroundColor(bgColor)
	d.containerHealthStartupIntervalField.SetLabelColor(style.DialogFgColor)
	d.containerHealthStartupIntervalField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerHealthStartupIntervalField.SetFieldWidth(healthPageMultiRowFieldWidth)

	// interval
	d.containerHealthIntervalField.SetLabel("Interval:")
	d.containerHealthIntervalField.SetLabelWidth(healthPageLabelWidth)
	d.containerHealthIntervalField.SetBackgroundColor(bgColor)
	d.containerHealthIntervalField.SetLabelColor(style.DialogFgColor)
	d.containerHealthIntervalField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerHealthIntervalField.SetFieldWidth(healthPageMultiRowFieldWidth)

	multiItemRow02 := tview.NewFlex().SetDirection(tview.FlexColumn)
	multiItemRow02.AddItem(d.containerHealthIntervalField, 0, 1, true)
	multiItemRow02.AddItem(d.containerHealthStartupIntervalField, 0, 1, true)
	multiItemRow02.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, true)
	multiItemRow02.SetBackgroundColor(bgColor)

	// multi primitive row03
	// startup retries
	d.containerHealthStartupRetriesField.SetLabel("Startup retries:")
	d.containerHealthStartupRetriesField.SetLabelWidth(healthPageSecColLabelWidth)
	d.containerHealthStartupRetriesField.SetBackgroundColor(bgColor)
	d.containerHealthStartupRetriesField.SetLabelColor(style.DialogFgColor)
	d.containerHealthStartupRetriesField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerHealthStartupRetriesField.SetFieldWidth(healthPageMultiRowFieldWidth)

	// retires
	d.containerHealthRetriesField.SetLabel("Retries:")
	d.containerHealthRetriesField.SetLabelWidth(healthPageLabelWidth)
	d.containerHealthRetriesField.SetBackgroundColor(bgColor)
	d.containerHealthRetriesField.SetLabelColor(style.DialogFgColor)
	d.containerHealthRetriesField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerHealthRetriesField.SetFieldWidth(healthPageMultiRowFieldWidth)

	multiItemRow03 := tview.NewFlex().SetDirection(tview.FlexColumn)
	multiItemRow03.AddItem(d.containerHealthRetriesField, 0, 1, true)
	multiItemRow03.AddItem(d.containerHealthStartupRetriesField, 0, 1, true)
	multiItemRow03.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, true)
	multiItemRow03.SetBackgroundColor(bgColor)

	// multi primitive row04
	// startup timeout
	d.containerHealthStartupTimeoutField.SetLabel("Startup timeout:")
	d.containerHealthStartupTimeoutField.SetLabelWidth(healthPageSecColLabelWidth)
	d.containerHealthStartupTimeoutField.SetBackgroundColor(bgColor)
	d.containerHealthStartupTimeoutField.SetLabelColor(style.DialogFgColor)
	d.containerHealthStartupTimeoutField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerHealthStartupTimeoutField.SetFieldWidth(healthPageMultiRowFieldWidth)

	// timeout
	d.containerHealthTimeoutField.SetLabel("Timeout:")
	d.containerHealthTimeoutField.SetLabelWidth(healthPageLabelWidth)
	d.containerHealthTimeoutField.SetBackgroundColor(bgColor)
	d.containerHealthTimeoutField.SetLabelColor(style.DialogFgColor)
	d.containerHealthTimeoutField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerHealthTimeoutField.SetFieldWidth(healthPageMultiRowFieldWidth)

	multiItemRow04 := tview.NewFlex().SetDirection(tview.FlexColumn)
	multiItemRow04.AddItem(d.containerHealthTimeoutField, 0, 1, true)
	multiItemRow04.AddItem(d.containerHealthStartupTimeoutField, 0, 1, true)
	multiItemRow04.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, true)
	multiItemRow04.SetBackgroundColor(bgColor)

	// health page
	d.healthPage.SetDirection(tview.FlexRow)
	d.healthPage.AddItem(d.containerHealthCmdField, 1, 0, true)
	d.healthPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.healthPage.AddItem(d.containerHealthStartupCmdField, 1, 0, true)
	d.healthPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.healthPage.AddItem(multiItemRow01, 1, 0, true)
	d.healthPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.healthPage.AddItem(multiItemRow02, 1, 0, true)
	d.healthPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.healthPage.AddItem(multiItemRow03, 1, 0, true)
	d.healthPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.healthPage.AddItem(multiItemRow04, 1, 0, true)
	d.healthPage.SetBackgroundColor(bgColor)
}

func (d *ContainerCreateDialog) setupNetworkPageUI() {
	bgColor := style.DialogBgColor
	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected
	inputFieldBgColor := style.InputFieldBgColor
	networkingPageLabelWidth := 13

	// hostname field
	d.containerHostnameField.SetLabel("hostname:")
	d.containerHostnameField.SetLabelWidth(networkingPageLabelWidth)
	d.containerHostnameField.SetBackgroundColor(bgColor)
	d.containerHostnameField.SetLabelColor(style.DialogFgColor)
	d.containerHostnameField.SetFieldBackgroundColor(inputFieldBgColor)

	// IP field
	d.containerIPAddrField.SetLabel("ip address:")
	d.containerIPAddrField.SetLabelWidth(networkingPageLabelWidth)
	d.containerIPAddrField.SetBackgroundColor(bgColor)
	d.containerIPAddrField.SetLabelColor(style.DialogFgColor)
	d.containerIPAddrField.SetFieldBackgroundColor(inputFieldBgColor)

	// mac field
	d.containerMacAddrField.SetLabel("mac address:")
	d.containerMacAddrField.SetLabelWidth(networkingPageLabelWidth)
	d.containerMacAddrField.SetBackgroundColor(bgColor)
	d.containerMacAddrField.SetLabelColor(style.DialogFgColor)
	d.containerMacAddrField.SetFieldBackgroundColor(inputFieldBgColor)

	// network field
	d.containerNetworkField.SetLabel("network:")
	d.containerNetworkField.SetLabelWidth(networkingPageLabelWidth)
	d.containerNetworkField.SetBackgroundColor(bgColor)
	d.containerNetworkField.SetLabelColor(style.DialogFgColor)
	d.containerNetworkField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	d.containerNetworkField.SetFieldBackgroundColor(inputFieldBgColor)

	// network settings page
	d.networkingPage.SetDirection(tview.FlexRow)
	d.networkingPage.AddItem(d.containerHostnameField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.networkingPage.AddItem(d.containerIPAddrField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.networkingPage.AddItem(d.containerMacAddrField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.networkingPage.AddItem(d.containerNetworkField, 1, 0, true)
	d.networkingPage.SetBackgroundColor(bgColor)
}

func (d *ContainerCreateDialog) setupPortsPageUI() {
	bgColor := style.DialogBgColor
	inputFieldBgColor := style.InputFieldBgColor
	portPageLabelWidth := 15

	// publish field
	d.containerPortPublishField.SetLabel("publish ports:")
	d.containerPortPublishField.SetLabelWidth(portPageLabelWidth)
	d.containerPortPublishField.SetBackgroundColor(bgColor)
	d.containerPortPublishField.SetLabelColor(style.DialogFgColor)
	d.containerPortPublishField.SetFieldBackgroundColor(inputFieldBgColor)

	// expose field
	d.containerPortExposeField.SetLabel("expose ports:")
	d.containerPortExposeField.SetLabelWidth(portPageLabelWidth)
	d.containerPortExposeField.SetBackgroundColor(bgColor)
	d.containerPortExposeField.SetLabelColor(style.DialogFgColor)
	d.containerPortExposeField.SetFieldBackgroundColor(inputFieldBgColor)

	// publish all field
	d.ContainerPortPublishAllField.SetLabel("publish all ")
	d.ContainerPortPublishAllField.SetLabelWidth(portPageLabelWidth)
	d.ContainerPortPublishAllField.SetBackgroundColor(bgColor)
	d.ContainerPortPublishAllField.SetLabelColor(style.DialogFgColor)
	d.ContainerPortPublishAllField.SetChecked(false)
	d.ContainerPortPublishAllField.SetFieldBackgroundColor(inputFieldBgColor)

	d.portPage.SetDirection(tview.FlexRow)
	d.portPage.AddItem(d.containerPortPublishField, 1, 0, true)
	d.portPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.portPage.AddItem(d.ContainerPortPublishAllField, 1, 0, true)
	d.portPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.portPage.AddItem(d.containerPortExposeField, 1, 0, true)
	d.portPage.SetBackgroundColor(bgColor)
}

func (d *ContainerCreateDialog) setupSecurityPageUI() {
	bgColor := style.DialogBgColor
	inputFieldBgColor := style.InputFieldBgColor
	securityOptsLabelWidth := 10

	// selinux label
	d.containerSelinuxLabelField.SetLabel("label:")
	d.containerSelinuxLabelField.SetLabelWidth(securityOptsLabelWidth)
	d.containerSelinuxLabelField.SetBackgroundColor(bgColor)
	d.containerSelinuxLabelField.SetLabelColor(style.DialogFgColor)
	d.containerSelinuxLabelField.SetFieldBackgroundColor(inputFieldBgColor)

	// apparmor
	d.containerApparmorField.SetLabel("apparmor:")
	d.containerApparmorField.SetLabelWidth(securityOptsLabelWidth)
	d.containerApparmorField.SetBackgroundColor(bgColor)
	d.containerApparmorField.SetLabelColor(style.DialogFgColor)
	d.containerApparmorField.SetFieldBackgroundColor(inputFieldBgColor)

	// seccomp
	d.containerSeccompField.SetLabel("seccomp:")
	d.containerSeccompField.SetLabelWidth(securityOptsLabelWidth)
	d.containerSeccompField.SetBackgroundColor(bgColor)
	d.containerSeccompField.SetLabelColor(style.DialogFgColor)
	d.containerSeccompField.SetFieldBackgroundColor(inputFieldBgColor)

	// mask
	d.containerMaskField.SetLabel("mask:")
	d.containerMaskField.SetLabelWidth(securityOptsLabelWidth)
	d.containerMaskField.SetBackgroundColor(bgColor)
	d.containerMaskField.SetLabelColor(style.DialogFgColor)
	d.containerMaskField.SetFieldBackgroundColor(inputFieldBgColor)

	// unmask
	d.containerUnmaskField.SetLabel("unmask:")
	d.containerUnmaskField.SetLabelWidth(securityOptsLabelWidth)
	d.containerUnmaskField.SetBackgroundColor(bgColor)
	d.containerUnmaskField.SetLabelColor(style.DialogFgColor)
	d.containerUnmaskField.SetFieldBackgroundColor(inputFieldBgColor)

	// no-new-privileges
	d.containerNoNewPrivField.SetLabel("no new privileges ")
	d.containerNoNewPrivField.SetBackgroundColor(bgColor)
	d.containerNoNewPrivField.SetLabelColor(style.DialogFgColor)
	d.containerNoNewPrivField.SetBackgroundColor(bgColor)
	d.containerNoNewPrivField.SetLabelColor(style.DialogFgColor)
	d.containerNoNewPrivField.SetChecked(false)
	d.containerNoNewPrivField.SetFieldBackgroundColor(inputFieldBgColor)

	// security options page
	d.securityOptsPage.SetDirection(tview.FlexRow)
	d.securityOptsPage.AddItem(d.containerSelinuxLabelField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.containerApparmorField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.containerSeccompField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.containerMaskField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.containerUnmaskField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.containerNoNewPrivField, 1, 0, true)
}

func (d *ContainerCreateDialog) setupVolumePageUI() {
	bgColor := style.DialogBgColor
	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected
	inputFieldBgColor := style.InputFieldBgColor
	volumePageLabelWidth := 14

	// volume
	d.containerVolumeField.SetLabel("volume:")
	d.containerVolumeField.SetLabelWidth(volumePageLabelWidth)
	d.containerVolumeField.SetBackgroundColor(bgColor)
	d.containerVolumeField.SetLabelColor(style.DialogFgColor)
	d.containerVolumeField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	d.containerVolumeField.SetFieldBackgroundColor(inputFieldBgColor)

	// image volume
	d.containerImageVolumeField.SetLabel("image volume:")
	d.containerImageVolumeField.SetLabelWidth(volumePageLabelWidth)
	d.containerImageVolumeField.SetBackgroundColor(bgColor)
	d.containerImageVolumeField.SetLabelColor(style.DialogFgColor)
	d.containerImageVolumeField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	d.containerImageVolumeField.SetFieldBackgroundColor(inputFieldBgColor)

	// volume settings page
	d.volumePage.SetDirection(tview.FlexRow)
	d.volumePage.AddItem(d.containerVolumeField, 1, 0, true)
	d.volumePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.volumePage.AddItem(d.containerImageVolumeField, 1, 0, true)
	d.volumePage.SetBackgroundColor(bgColor)
}

// Display displays this primitive.
func (d *ContainerCreateDialog) Display() {
	d.display = true
	d.initData()
	d.focusElement = createCategoryPagesFocus
}

// IsDisplay returns true if primitive is shown.
func (d *ContainerCreateDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *ContainerCreateDialog) Hide() {
	d.display = false
}

// HasFocus returns whether or not this primitive has focus.
func (d *ContainerCreateDialog) HasFocus() bool {
	if d.categories.HasFocus() || d.categoryPages.HasFocus() {
		return true
	}

	return d.Box.HasFocus() || d.form.HasFocus()
}

// dropdownHasFocus returns true if container create dialog dropdown primitives
// has focus.
func (d *ContainerCreateDialog) dropdownHasFocus() bool {
	if d.containerImageField.HasFocus() || d.containerPodField.HasFocus() {
		return true
	}
	if d.containerVolumeField.HasFocus() || d.containerImageVolumeField.HasFocus() {
		return true
	}
	if d.containerNetworkField.HasFocus() || d.containerHealthOnFailureField.HasFocus() {
		return true
	}
	return false
}

// Focus is called when this primitive receives focus.
func (d *ContainerCreateDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	// form has focus
	case createContainerFormFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = createCategoriesFocus // category text view
				d.Focus(delegate)
				d.form.SetFocus(0)
				return nil
			}
			if event.Key() == tcell.KeyEnter {
				//d.pullSelectHandler()
				return nil
			}
			return event
		})
		delegate(d.form)
	// category text view
	case createCategoriesFocus:
		d.categories.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = createCategoryPagesFocus // category page view
				d.Focus(delegate)
				return nil
			}
			// scroll between categories
			event = utils.ParseKeyEventKey(event)
			if event.Key() == tcell.KeyDown {
				d.nextCategory()
			}
			if event.Key() == tcell.KeyUp {
				d.previousCategory()
			}
			return event
		})
		delegate(d.categories)
	// basic info page
	case createContainerNameFieldFocus:
		delegate(d.containerNameField)
	case createContainerImageFieldFocus:
		delegate(d.containerImageField)
	case createcontainerPodFieldFocis:
		delegate(d.containerPodField)
	case createContainerLabelsFieldFocus:
		delegate(d.containerLabelsField)
	case createContainerRemoveFieldFocus:
		delegate(d.containerRemoveField)
	// security options page
	case createContainerSelinuxLabelFieldFocus:
		delegate(d.containerSelinuxLabelField)
	case createContainerApprarmorFieldFocus:
		delegate(d.containerApparmorField)
	case createContainerSeccompFeildFocus:
		delegate(d.containerSeccompField)
	case createContainerMaskFieldFocus:
		delegate(d.containerMaskField)
	case createContainerUnmaskFieldFocus:
		delegate(d.containerUnmaskField)
	case createContainerNoNewPrivFieldFocus:
		delegate(d.containerNoNewPrivField)
	// networking page
	case createContainerHostnameFieldFocus:
		delegate(d.containerHostnameField)
	case createContainerIPAddrFieldFocus:
		delegate(d.containerIPAddrField)
	case createContainerMacAddrFieldFocus:
		delegate(d.containerMacAddrField)
	case createContainerNetworkFieldFocus:
		delegate(d.containerNetworkField)
	// port page
	// networking page
	case createContainerPortPublishFieldFocus:
		delegate(d.containerPortPublishField)
	case createContainerPortPublishAllFieldFocus:
		delegate(d.ContainerPortPublishAllField)
	case createContainerPortExposeFieldFocus:
		delegate(d.containerPortExposeField)
	// dns page
	case createContainerDNSServersFieldFocus:
		delegate(d.containerDNSServersField)
	case createContainerDNSOptionsFieldFocus:
		delegate(d.containerDNSOptionsField)
	case createContainerDNSSearchFieldFocus:
		delegate(d.containerDNSSearchField)
	// volume page
	case createContainerVolumeFieldFocus:
		delegate(d.containerVolumeField)
	case createContainerImageVolumeFieldFocus:
		delegate(d.containerImageVolumeField)
	// health page
	case containerHealthCmdFieldFocus:
		delegate(d.containerHealthCmdField)
	case containerHealthStartupCmdFieldFocus:
		delegate(d.containerHealthStartupCmdField)
	case containerHealthOnFailureFieldFocus:
		delegate(d.containerHealthOnFailureField)
	case containerHealthIntervalFieldFocus:
		delegate(d.containerHealthIntervalField)
	case containerHealthStartupIntervalFieldFocus:
		delegate(d.containerHealthStartupIntervalField)
	case containerHealthTimeoutFieldFocus:
		delegate(d.containerHealthTimeoutField)
	case containerHealthStartupTimeoutFieldFocus:
		delegate(d.containerHealthStartupTimeoutField)
	case containerHealthRetriesFieldFocus:
		delegate(d.containerHealthRetriesField)
	case containerHealthStartupRetriesFieldFocus:
		delegate(d.containerHealthStartupRetriesField)
	case containerHealthStartPeriodFieldFocus:
		delegate(d.containerHealthStartPeriodField)
	case containerHealthStartupSuccessFieldFocus:
		delegate(d.containerHealthStartupSuccessField)
	// category page
	case createCategoryPagesFocus:
		delegate(d.categoryPages)
	}
}

func (d *ContainerCreateDialog) initCustomInputHanlers() {
	// pod name dropdown
	d.containerPodField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = utils.ParseKeyEventKey(event)
		return event
	})
	// container image volume
	d.containerImageField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = utils.ParseKeyEventKey(event)
		return event
	})
	// container network dropdown
	d.containerNetworkField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = utils.ParseKeyEventKey(event)
		return event
	})
	// container volume dropdown
	d.containerVolumeField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = utils.ParseKeyEventKey(event)
		return event
	})
	// container image volume dropdown
	d.containerImageVolumeField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = utils.ParseKeyEventKey(event)
		return event
	})
}

// InputHandler returns input handler function for this primitive.
func (d *ContainerCreateDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("container create dialog: event %v received", event)
		if event.Key() == tcell.KeyEsc && !d.dropdownHasFocus() {
			d.cancelHandler()
			return
		}
		if d.basicInfoPage.HasFocus() {
			if handler := d.basicInfoPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setBasicInfoPageNextFocus()
				}
				handler(event, setFocus)
				return
			}
		}
		if d.dnsPage.HasFocus() {
			if handler := d.dnsPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setDNSSettingsPageNextFocus()
				}
				handler(event, setFocus)
				return
			}
		}
		if d.networkingPage.HasFocus() {
			if handler := d.networkingPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setNetworkSettingsPageNextFocus()
				}
				handler(event, setFocus)
				return
			}
		}
		if d.portPage.HasFocus() {
			if handler := d.portPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setPortPageNextFocus()
				}
				handler(event, setFocus)
				return
			}
		}
		if d.securityOptsPage.HasFocus() {
			if handler := d.securityOptsPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setSecurityOptionsPageNextFocus()
				}
				handler(event, setFocus)
				return
			}
		}
		if d.volumePage.HasFocus() {
			if handler := d.volumePage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setVolumeSettingsPageNextFocus()
				}
				handler(event, setFocus)
				return
			}
		}

		if d.healthPage.HasFocus() {
			if handler := d.healthPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setHealthSettingsPageNextFocus()
				}
				handler(event, setFocus)
				return
			}
		}

		if d.categories.HasFocus() {
			if categroryHandler := d.categories.InputHandler(); categroryHandler != nil {
				categroryHandler(event, setFocus)
				return
			}
		}
		if d.form.HasFocus() {
			if formHandler := d.form.InputHandler(); formHandler != nil {
				if event.Key() == tcell.KeyEnter {
					enterButton := d.form.GetButton(d.form.GetButtonCount() - 1)
					if enterButton.HasFocus() {
						d.createHandler()
					}
				}
				formHandler(event, setFocus)
				return
			}
		}

	})
}

// SetRect set rects for this primitive.
func (d *ContainerCreateDialog) SetRect(x, y, width, height int) {

	if width > containerCreateDialogMaxWidth {
		emptySpace := (width - containerCreateDialogMaxWidth) / 2
		x = x + emptySpace
		width = containerCreateDialogMaxWidth
	}

	if height > containerCreateDialogHeight {
		emptySpace := (height - containerCreateDialogHeight) / 2
		y = y + emptySpace
		height = containerCreateDialogHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *ContainerCreateDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}
	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetCancelFunc sets form cancel button selected function.
func (d *ContainerCreateDialog) SetCancelFunc(handler func()) *ContainerCreateDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2)
	cancelButton.SetSelectedFunc(handler)
	return d
}

// SetCreateFunc sets form create button selected function.
func (d *ContainerCreateDialog) SetCreateFunc(handler func()) *ContainerCreateDialog {
	d.createHandler = handler
	enterButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	enterButton.SetSelectedFunc(handler)
	return d
}

func (d *ContainerCreateDialog) setActiveCategory(index int) {
	fgColor := style.DialogFgColor
	bgColor := style.ButtonBgColor
	ctgTextColor := style.GetColorHex(fgColor)
	ctgBgColor := style.GetColorHex(bgColor)

	d.activePageIndex = index
	d.categories.Clear()
	var ctgList []string
	alignedList, _ := utils.AlignStringListWidth(d.categoryLabels)
	for i := 0; i < len(alignedList); i++ {
		if i == index {
			ctgList = append(ctgList, fmt.Sprintf("[%s:%s:b]-> %s ", ctgTextColor, ctgBgColor, alignedList[i]))
			continue
		}
		ctgList = append(ctgList, fmt.Sprintf("[-:-:-]   %s ", alignedList[i]))
	}
	d.categories.SetText(strings.Join(ctgList, "\n"))

	// switch the page
	d.categoryPages.SwitchToPage(d.categoryLabels[index])
}

func (d *ContainerCreateDialog) nextCategory() {
	activePage := d.activePageIndex
	if d.activePageIndex < len(d.categoryLabels)-1 {
		activePage = activePage + 1
		d.setActiveCategory(activePage)
		return
	}
	d.setActiveCategory(0)
}

func (d *ContainerCreateDialog) previousCategory() {
	activePage := d.activePageIndex
	if d.activePageIndex > 0 {
		activePage = activePage - 1
		d.setActiveCategory(activePage)
		return
	}
	d.setActiveCategory(len(d.categoryLabels) - 1)
}

func (d *ContainerCreateDialog) initData() {
	// get available images
	imgList, _ := images.List()
	d.imageList = imgList
	imgOptions := []string{""}
	for i := 0; i < len(d.imageList); i++ {
		if d.imageList[i].ID == "<none>" {
			imgOptions = append(imgOptions, d.imageList[i].ID)
			continue
		}
		imgname := d.imageList[i].Repository + ":" + d.imageList[i].Tag
		imgOptions = append(imgOptions, imgname)
	}

	// get available pods
	podOptions := []string{""}
	podList, _ := pods.List()
	d.podList = podList
	for i := 0; i < len(podList); i++ {
		podOptions = append(podOptions, podList[i].Name)
	}

	// get available networks
	networkOptions := []string{""}
	networkList, _ := networks.List()
	for i := 0; i < len(networkList); i++ {
		networkOptions = append(networkOptions, networkList[i][1])
	}

	// get available volumes
	imageVolumeOptions := []string{"", "ignore", "tmpfs", "anonymous"}
	volumeOptions := []string{""}
	volList, _ := volumes.List()
	for i := 0; i < len(volList); i++ {
		volumeOptions = append(volumeOptions, volList[i].Name)
	}

	d.setActiveCategory(0)
	d.containerNameField.SetText("")
	d.containerImageField.SetOptions(imgOptions, nil)
	d.containerImageField.SetCurrentOption(0)
	d.containerPodField.SetOptions(podOptions, nil)
	d.containerPodField.SetCurrentOption(0)
	d.containerLabelsField.SetText("")
	d.containerRemoveField.SetChecked(false)
	d.containerSelinuxLabelField.SetText("")
	d.containerApparmorField.SetText("")
	d.containerSeccompField.SetText("")
	d.containerMaskField.SetText("")
	d.containerUnmaskField.SetText("")
	d.containerNoNewPrivField.SetChecked(false)
	d.containerPortPublishField.SetText("")
	d.ContainerPortPublishAllField.SetChecked(false)
	d.containerPortExposeField.SetText("")
	d.containerHostnameField.SetText("")
	d.containerIPAddrField.SetText("")
	d.containerMacAddrField.SetText("")
	d.containerNetworkField.SetOptions(networkOptions, nil)
	d.containerNetworkField.SetCurrentOption(0)
	d.containerDNSServersField.SetText("")
	d.containerDNSSearchField.SetText("")
	d.containerDNSOptionsField.SetText("")
	d.containerVolumeField.SetOptions(volumeOptions, nil)
	d.containerVolumeField.SetCurrentOption(0)
	d.containerImageVolumeField.SetOptions(imageVolumeOptions, nil)
	d.containerImageVolumeField.SetCurrentOption(0)

	// health options
	d.containerHealthCmdField.SetText("")
	d.containerHealthStartupCmdField.SetText("")
	d.containerHealthOnFailureField.SetCurrentOption(0)
	d.containerHealthIntervalField.SetText("")
	d.containerHealthStartupIntervalField.SetText("")
	d.containerHealthTimeoutField.SetText("")
	d.containerHealthStartupTimeoutField.SetText("")
	d.containerHealthRetriesField.SetText("")
	d.containerHealthStartupRetriesField.SetText("")
	d.containerHealthStartPeriodField.SetText("")
	d.containerHealthStartupSuccessField.SetText("")
}

func (d *ContainerCreateDialog) setPortPageNextFocus() {
	if d.containerPortPublishField.HasFocus() {
		d.focusElement = createContainerPortPublishAllFieldFocus
	} else if d.ContainerPortPublishAllField.HasFocus() {
		d.focusElement = createContainerPortExposeFieldFocus
	} else {
		d.focusElement = createContainerFormFocus
	}
}

func (d *ContainerCreateDialog) setBasicInfoPageNextFocus() {
	if d.containerNameField.HasFocus() {
		d.focusElement = createContainerImageFieldFocus
	} else if d.containerImageField.HasFocus() {
		d.focusElement = createcontainerPodFieldFocis
	} else if d.containerPodField.HasFocus() {
		d.focusElement = createContainerLabelsFieldFocus
	} else if d.containerLabelsField.HasFocus() {
		d.focusElement = createContainerRemoveFieldFocus
	} else {
		d.focusElement = createContainerFormFocus
	}
}

func (d *ContainerCreateDialog) setSecurityOptionsPageNextFocus() {
	if d.containerSelinuxLabelField.HasFocus() {
		d.focusElement = createContainerApprarmorFieldFocus
	} else if d.containerApparmorField.HasFocus() {
		d.focusElement = createContainerSeccompFeildFocus
	} else if d.containerSeccompField.HasFocus() {
		d.focusElement = createContainerMaskFieldFocus
	} else if d.containerMaskField.HasFocus() {
		d.focusElement = createContainerUnmaskFieldFocus
	} else if d.containerUnmaskField.HasFocus() {
		d.focusElement = createContainerNoNewPrivFieldFocus
	} else {
		d.focusElement = createContainerFormFocus
	}
}

func (d *ContainerCreateDialog) setNetworkSettingsPageNextFocus() {
	if d.containerHostnameField.HasFocus() {
		d.focusElement = createContainerIPAddrFieldFocus
	} else if d.containerIPAddrField.HasFocus() {
		d.focusElement = createContainerMacAddrFieldFocus
	} else if d.containerMacAddrField.HasFocus() {
		d.focusElement = createContainerNetworkFieldFocus
	} else {
		d.focusElement = createContainerFormFocus
	}
}

func (d *ContainerCreateDialog) setDNSSettingsPageNextFocus() {
	if d.containerDNSServersField.HasFocus() {
		d.focusElement = createContainerDNSOptionsFieldFocus
	} else if d.containerDNSOptionsField.HasFocus() {
		d.focusElement = createContainerDNSSearchFieldFocus
	} else {
		d.focusElement = createContainerFormFocus
	}
}

func (d *ContainerCreateDialog) setHealthSettingsPageNextFocus() {
	if d.containerHealthCmdField.HasFocus() {
		d.focusElement = containerHealthStartupCmdFieldFocus
	} else if d.containerHealthStartupCmdField.HasFocus() {
		d.focusElement = containerHealthOnFailureFieldFocus
	} else if d.containerHealthOnFailureField.HasFocus() {
		d.focusElement = containerHealthStartupSuccessFieldFocus
	} else if d.containerHealthStartupSuccessField.HasFocus() {
		d.focusElement = containerHealthStartPeriodFieldFocus
	} else if d.containerHealthStartPeriodField.HasFocus() {
		d.focusElement = containerHealthIntervalFieldFocus
	} else if d.containerHealthIntervalField.HasFocus() {
		d.focusElement = containerHealthStartupIntervalFieldFocus
	} else if d.containerHealthStartupIntervalField.HasFocus() {
		d.focusElement = containerHealthRetriesFieldFocus
	} else if d.containerHealthRetriesField.HasFocus() {
		d.focusElement = containerHealthStartupRetriesFieldFocus
	} else if d.containerHealthStartupRetriesField.HasFocus() {
		d.focusElement = containerHealthTimeoutFieldFocus
	} else if d.containerHealthTimeoutField.HasFocus() {
		d.focusElement = containerHealthStartupTimeoutFieldFocus
	} else {
		d.focusElement = createContainerFormFocus
	}
}

func (d *ContainerCreateDialog) setVolumeSettingsPageNextFocus() {
	if d.containerVolumeField.HasFocus() {
		d.focusElement = createContainerImageVolumeFieldFocus
	} else {
		d.focusElement = createContainerFormFocus
	}
}

// ContainerCreateOptions returns new network options.
func (d *ContainerCreateDialog) ContainerCreateOptions() containers.CreateOptions {
	var (
		labels           []string
		imageID          string
		podID            string
		dnsServers       []string
		dnsOptions       []string
		dnsSearchDomains []string
		publish          []string
		expose           []string
		volume           string
		imageVolume      string
		selinuxOpts      []string
	)

	for _, label := range strings.Split(d.containerLabelsField.GetText(), " ") {
		if label != "" {
			labels = append(labels, label)
		}
	}

	selectedImageIndex, _ := d.containerImageField.GetCurrentOption()
	if len(d.imageList) > 0 && selectedImageIndex > 0 {
		imageID = d.imageList[selectedImageIndex-1].ID
	}
	selectedPodIndex, _ := d.containerPodField.GetCurrentOption()
	if len(d.podList) > 0 && selectedPodIndex > 0 {
		podID = d.podList[selectedPodIndex-1].Id
	}

	// ports
	for _, p := range strings.Split(d.containerPortPublishField.GetText(), " ") {
		if p != "" {
			publish = append(publish, p)
		}
	}
	for _, e := range strings.Split(d.containerPortExposeField.GetText(), " ") {
		if e != "" {
			expose = append(expose, e)
		}
	}
	// DNS setting
	for _, dns := range strings.Split(d.containerDNSServersField.GetText(), " ") {
		if dns != "" {
			dnsServers = append(dnsServers, dns)
		}
	}
	for _, do := range strings.Split(d.containerDNSOptionsField.GetText(), " ") {
		if do != "" {
			dnsOptions = append(dnsOptions, do)
		}
	}
	for _, ds := range strings.Split(d.containerDNSSearchField.GetText(), " ") {
		if ds != "" {
			dnsSearchDomains = append(dnsSearchDomains, ds)
		}
	}
	_, volume = d.containerVolumeField.GetCurrentOption()
	_, imageVolume = d.containerImageVolumeField.GetCurrentOption()

	// security options
	for _, selinuxLabel := range strings.Split(d.containerSelinuxLabelField.GetText(), " ") {
		if selinuxLabel != "" {
			selinuxOpts = append(selinuxOpts, selinuxLabel)
		}
	}

	// health check
	_, healthOnFailure := d.containerHealthOnFailureField.GetCurrentOption()

	_, network := d.containerNetworkField.GetCurrentOption()
	opts := containers.CreateOptions{
		Name:                  d.containerNameField.GetText(),
		Image:                 imageID,
		Pod:                   podID,
		Labels:                labels,
		Remove:                d.containerRemoveField.IsChecked(),
		Hostname:              d.containerHostnameField.GetText(),
		MacAddress:            d.containerMacAddrField.GetText(),
		IPAddress:             d.containerIPAddrField.GetText(),
		Network:               network,
		Publish:               publish,
		Expose:                expose,
		PublishAll:            d.ContainerPortPublishAllField.IsChecked(),
		DNSServer:             dnsServers,
		DNSOptions:            dnsOptions,
		DNSSearchDomain:       dnsSearchDomains,
		Volume:                volume,
		ImageVolume:           imageVolume,
		SelinuxOpts:           selinuxOpts,
		ApparmorProfile:       d.containerApparmorField.GetText(),
		Seccomp:               d.containerSeccompField.GetText(),
		NoNewPriv:             d.containerNoNewPrivField.IsChecked(),
		Mask:                  d.containerMaskField.GetText(),
		Unmask:                d.containerUnmaskField.GetText(),
		HealthCmd:             strings.TrimSpace(d.containerHealthCmdField.GetText()),
		HealthInterval:        strings.TrimSpace(d.containerHealthIntervalField.GetText()),
		HealthRetries:         strings.TrimSpace(d.containerHealthRetriesField.GetText()),
		HealthStartPeroid:     strings.TrimSpace(d.containerHealthStartPeriodField.GetText()),
		HealthTimeout:         strings.TrimSpace(d.containerHealthTimeoutField.GetText()),
		HealthOnFailure:       strings.TrimSpace(healthOnFailure),
		HealthStartupCmd:      strings.TrimSpace(d.containerHealthStartupCmdField.GetText()),
		HealthStartupInterval: strings.TrimSpace(d.containerHealthStartupIntervalField.GetText()),
		HealthStartupRetries:  strings.TrimSpace(d.containerHealthStartupRetriesField.GetText()),
		HealthStartupSuccess:  strings.TrimSpace(d.containerHealthStartupSuccessField.GetText()),
		HealthStartupTimeout:  strings.TrimSpace(d.containerHealthStartupTimeoutField.GetText()),
	}
	return opts
}
