package display

import (
	"get-service-version/entity"
	model "get-service-version/model"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var (
	titleEnvStyle = lipgloss.NewStyle().
			MarginLeft(2).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#5438a6")).
			Width(122).
			Align(lipgloss.Center).
			Bold(true)

	titleStyle = lipgloss.NewStyle().
			MarginLeft(2).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Width(60).
			Align(lipgloss.Center).
			Bold(true)

	keyStyle    = lipgloss.NewStyle().Width(15).MarginLeft(2).Bold(true)
	keyInaStyle = lipgloss.NewStyle().Width(15).Bold(true)
	valueStyle  = lipgloss.NewStyle().Width(45).MarginLeft(2)
)

func GetInfo(hostblue, hostgreen, hostlive string) (error, *entity.HealthCheck, *entity.HealthCheck, *entity.HealthCheck) {
	var (
		blue, green, live *entity.HealthCheck
		wg                sync.WaitGroup
		err               error
	)

	wg.Add(3)

	blue = &entity.HealthCheck{}
	green = &entity.HealthCheck{}
	live = &entity.HealthCheck{}

	go func() {
		defer wg.Done()
		_, er := NewClient().R().
			SetResult(blue).
			Get(hostblue)
		if er != nil {
			err = er
		}
	}()

	go func() {
		defer wg.Done()
		_, er := NewClient().R().
			SetResult(green).
			Get(hostgreen)
		if er != nil {
			err = er
		}
	}()

	go func() {
		defer wg.Done()
		_, er := NewClient().R().
			SetResult(live).
			Get(hostlive)
		if er != nil {
			err = er
		}
	}()

	wg.Wait()
	if err != nil {
		return err, nil, nil, nil
	}
	return nil, blue, green, live

}

func NewClient() *resty.Client {
	client := resty.New().
		SetDebug(false).
		SetDisableWarn(true).
		SetRetryCount(2).
		SetRetryWaitTime(500 * time.Millisecond).
		SetRetryMaxWaitTime(20 * time.Second).
		OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
			if resp.IsError() {
				return errors.Errorf("%s %s", resp.Status(), resp.String())
			}
			return nil
		}).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.IsError()
			},
		)
	return client
}

func Display(m *model.Model, title string) {
	var (
		err                                    error
		wg                                     sync.WaitGroup
		blue, green, live                      *entity.HealthCheck
		preprodblue, preprodgreen, preprodlive *entity.HealthCheck
	)
	r := ""
	wg.Add(2)
	go func() {
		defer wg.Done()
		err, blue, green, live = GetInfo(Services[title].ProdBlue, Services[title].ProdGreen, Services[title].Live)
	}()

	go func() {
		defer wg.Done()
		err, preprodblue, preprodgreen, preprodlive = GetInfo(Services[title].PreprodBlue, Services[title].PreprodGreen,
			Services[title].PreprodGreen)
	}()
	wg.Wait()

	if err != nil {
		m.Ready(err.Error())
		return
	}

	r += titleEnvStyle.Render("Production")
	if live.ClusterName == blue.ClusterName {
		r += Screen(blue, green)
	} else {
		r += Screen(green, blue)
	}

	r += "\n\n"
	r += titleEnvStyle.Render("PreProd")
	if preprodlive.ClusterName == preprodblue.ClusterName {
		r += Screen(preprodblue, preprodgreen)
	} else {
		r += Screen(preprodgreen, preprodblue)
	}

	r += "\n\n\n"
	r += "Press (esc) to go back"

	m.Ready(r)
}

func Screen(active, inactive *entity.HealthCheck) string {
	r := "\n"

	r += titleStyle.Render(active.ClusterName + " (ACTIVE)")
	r += titleStyle.Render(inactive.ClusterName)
	r += "\n"

	r += keyStyle.Render("appName")
	r += valueStyle.Render(active.AppName)
	r += keyInaStyle.Render("appName")
	r += valueStyle.Render(inactive.AppName)
	r += "\n"

	r += keyStyle.Render("appVersion")
	r += valueStyle.Render(active.AppVersion)
	r += keyInaStyle.Render("appVersion")
	r += valueStyle.Render(inactive.AppVersion)
	r += "\n"

	if !active.DeployedAt.IsZero() {
		r += keyStyle.Render("deployedAt")
		r += valueStyle.Render(active.DeployedAt.String())
		r += keyInaStyle.Render("deployedAt")
		r += valueStyle.Render(inactive.DeployedAt.String())
		r += "\n"
	}

	r += keyStyle.Render("git.hash")
	r += valueStyle.Render(active.Git.Hash)
	r += keyInaStyle.Render("git.hash")
	r += valueStyle.Render(inactive.Git.Hash)
	r += "\n"

	if active.Git.Branch != "" {
		r += keyStyle.Render("git.branch")
		r += valueStyle.Render(active.Git.Branch)
		r += keyInaStyle.Render("git.branch")
		r += valueStyle.Render(inactive.Git.Branch)
		r += "\n"
	}

	return r
}
