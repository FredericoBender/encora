package spider

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"encora/service"

	"github.com/falqondev/selenium"
	"github.com/falqondev/selenium/chrome"
	"github.com/pkg/errors"
)

// EncoraExtractor is responsible per extract all jobs details from the encora website
type EncoraExtractor struct {
	Config     EncoraConfig
	TimeLimit  int
	Logger     *log.Logger
	DebugLevel int8
	Wd         selenium.WebDriver
}

type EncoraConfig struct {
	ChromeDriverPath string
	Port             int32
}

func (s *EncoraExtractor) prepareSelenium() (selenium.WebDriver, *selenium.Service, error) {
	ops := []selenium.ServiceOption{}
	fmt.Println(s.Config.ChromeDriverPath)
	service, err := selenium.NewChromeDriverService(s.Config.ChromeDriverPath, int(s.Config.Port), ops...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error starting the ChromeDriver server: %v")
	}
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	chromeCaps := chrome.Capabilities{
		Path: "",
		Args: []string{
			`--headless`,
			`--window-size=1500,1200`,
			`--whitelisted-ips=""`,
			`--log-level=3`,
		},
	}

	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:"+strconv.Itoa(int(s.Config.Port))+"/wd/hub")
	if err != nil {
		return nil, nil, errors.Wrap(err, "selenium call browser")
	}
	return wd, service, nil
}

func (s *EncoraExtractor) checkDependencies() error {
	if s.Config.ChromeDriverPath == "" {
		return errors.New("missing Config.ChromeDriverPath")
	}
	if s.Config.Port == 0 {
		return errors.New("missing Config.Port")
	}
	if s.TimeLimit == 0 {
		return errors.New("missing TimeLimit bigger than 0")
	}
	return nil
}

func (s *EncoraExtractor) Run() (service.EncoraJobs, error) {
	jobsTitleXPATH := "//h3"
	jobAreasXPATH := "//div[contains(@id,'jobs-cards-container')]//p[contains(@class,'job-category weight-bold m-0 mr-half')]"
	countriesXPATH := "//div[contains(@id,'jobs-cards-container')]//p[contains(@class,'job-location weight-bold m-0')]"

	err := s.checkDependencies()
	if err != nil {
		return service.EncoraJobs{}, errors.Wrap(err, "dependencies:")
	}

	wd, seleniumService, err := s.prepareSelenium()
	if err != nil {
		return service.EncoraJobs{}, errors.Wrap(err, "prepare Selenium")
	}
	s.Wd = wd
	defer seleniumService.Stop()
	defer wd.Quit()

	if err := s.Wd.Get("https://careers.encora.com/search?category=All"); err != nil {
		return service.EncoraJobs{}, errors.Wrap(err, "get careers.encora.com")
	}

	time.Sleep(1 * time.Second)
	jobsTitle, err := s.GetListOfElements(jobsTitleXPATH)
	if err != nil {
		return service.EncoraJobs{}, errors.Wrap(err, "get jobs title")
	}
	jobsAreas, err := s.GetListOfElements(jobAreasXPATH)
	if err != nil {
		return service.EncoraJobs{}, errors.Wrap(err, "get jobs areas")
	}
	jobsCountries, err := s.GetListOfElements(countriesXPATH)
	if err != nil {
		return service.EncoraJobs{}, errors.Wrap(err, "get jobs countries")
	}
	jobsDetailsURLs, err := s.GetDetailsURL()
	if err != nil {
		return service.EncoraJobs{}, errors.Wrap(err, "get URL from jobs details")
	}
	fmt.Println("enter in get details from jobs")
	jobsDescriptions, err := s.GetDetailsFromJobURL(jobsDetailsURLs)
	if err != nil {
		return service.EncoraJobs{}, errors.Wrap(err, "get details from jobs")
	}

	encoraJobs := service.EncoraJobs{
		JobsTitle:       jobsTitle,
		JobAreas:        jobsAreas,
		JobsCountries:   jobsCountries,
		JobsDetailsURLs: jobsDetailsURLs,
		Description:     jobsDescriptions,
	}

	return encoraJobs, nil
}

func (s *EncoraExtractor) GetListOfElements(path string) ([]string, error) {
	jobs, err := s.Wd.FindElements(selenium.ByXPATH, path)
	if err != nil {
		return nil, errors.Wrap(err, "find elements")
	}
	listOfData := []string{}
	for _, currentValue := range jobs {
		textOfElement, err := currentValue.Text()
		if err != nil {
			return nil, errors.Wrap(err, " current value gettText()")
		}
		listOfData = append(listOfData, textOfElement)
	}
	return listOfData, nil
}

func (s *EncoraExtractor) GetDetailsURL() ([]string, error) {
	jobDetailsUrls, err := s.Wd.FindElements(selenium.ByXPATH, "//div[contains(@id,'jobs-cards-container')]//a")
	if err != nil {
		return nil, errors.Wrap(err, "find elements")
	}
	listOfData := []string{}
	for _, currentValue := range jobDetailsUrls {
		linkOfElement, err := currentValue.GetAttribute("href")
		if err != nil {
			return nil, errors.Wrap(err, " current value gettText()")
		}
		listOfData = append(listOfData, linkOfElement)
	}
	return listOfData, nil

}

func (s *EncoraExtractor) GetDetailsFromJobURL(urls []string) ([]string, error) {
	listOfData := []string{}
	for index, currentUrl := range urls {
		fmt.Println(index, "of", len(urls))
		if err := s.Wd.Get(currentUrl); err != nil {
			return nil, errors.Wrap(err, "get careers.encora.com")
		}
		for {
			jobDetails, err := s.Wd.FindElement(selenium.ByXPATH, "//p[contains(@id,'description-job')]")
			if err != nil {
				return nil, errors.Wrap(err, "find elements")
			}
			textOfElement, err := jobDetails.Text()
			if err != nil {
				return nil, errors.Wrap(err, " current value gettText()")
			}
			if len(textOfElement) > 1 {
				listOfData = append(listOfData, textOfElement)
				break
			}
			time.Sleep(1 * time.Second)
		}
	}
	return listOfData, nil

}
