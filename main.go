package main

import (
	"log"
	"os"
	"time"

	br "github.com/gvx3/sportuni-book/pkg/browser"
	"github.com/gvx3/sportuni-book/pkg/config"
	"github.com/playwright-community/playwright-go"
)

func main() {
	baseConfig, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	browser, err := br.NewBrowser()
	if err != nil {
		log.Fatalf("Could not initialize browser: %v", err)
	}
	defer browser.Close()

	freshContext, err := browser.NewContext()
	if err != nil {
		log.Fatalf("could not create fresh context: %v", err)
	}
	defer freshContext.Close()

	var activePage playwright.Page
	slots := baseConfig.ActivitySlots

	if !fileExist(baseConfig.StateFileName) {
		// No state file, use fresh login
		activePage, err = freshContext.NewPage()
		if err != nil {
			log.Fatalf("could not create page: %v", err)
		}

		err = br.FreshLogin(activePage, baseConfig.BaseURL, baseConfig.Email, baseConfig.Password)
		if err != nil {
			log.Fatalf("Fresh login failed: %v", err)
		}
	} else {

		existedContext, err := browser.NewContext(playwright.BrowserNewContextOptions{
			StorageStatePath: playwright.String(baseConfig.StateFileName),
		})
		if err != nil {
			log.Fatalf("could not create context with state file: %v", err)
		}
		defer existedContext.Close()

		activePage, err = existedContext.NewPage()
		if err != nil {
			log.Fatalf("could not create page: %v", err)
		}

		err = br.StateFileExpireLogin(activePage, baseConfig.BaseURL)
		if err != nil {
			log.Printf("State file expired: %v", err)

			if err := activePage.Close(); err != nil {
				log.Printf("Warning: failed to close existing page: %v", err)
			}
			// State file expired, use fresh login
			activePage, err = freshContext.NewPage()
			if err != nil {
				log.Fatalf("could not create page: %v", err)
			}

			err = br.FreshLogin(activePage, baseConfig.BaseURL, baseConfig.Email, baseConfig.Password)
			if err != nil {
				log.Fatalf("Fresh login failed: %v", err)
			}
		} else {
			// State file is valid
			err = br.StateFileSucceedLogin(activePage, baseConfig.BaseURL)
			if err != nil {
				log.Fatalf("State file login failed: %v", err)
			}
		}
	}

	err = br.BookCourse(activePage, slots)
	if err != nil {
		log.Fatalf("Failed to book course: %v", err)
	}

	time.Sleep(20 * time.Second)
}

func fileExist(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
