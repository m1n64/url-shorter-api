package services

import "github.com/mssola/user_agent"

type UserAgentInfo struct {
	Device     string
	OS         string
	OSVersion  string
	Browser    string
	BrowserVer string
}

type UserAgentParserService struct {
}

func NewUserAgentParserService() *UserAgentParserService {
	return &UserAgentParserService{}
}

func (s *UserAgentParserService) ParseUserAgent(uaString string) *UserAgentInfo {
	ua := user_agent.New(uaString)

	device := "PC"
	if ua.Mobile() {
		device = "Mobile"
	} else if ua.Bot() {
		device = "Bot"
	}

	os := ua.OSInfo()

	browser, browserVersion := ua.Browser()

	return &UserAgentInfo{
		Device:     device,
		OS:         os.Name,
		OSVersion:  os.Version,
		Browser:    browser,
		BrowserVer: browserVersion,
	}
}
