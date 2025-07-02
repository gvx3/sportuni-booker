package browser

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gvx3/sportuni-book/pkg/config"
	"github.com/playwright-community/playwright-go"
)

var ErrNoTimeSlots = errors.New("no time slot elements appeared on the page")

func headlessBrowserOptions(testing bool) playwright.BrowserTypeLaunchOptions {
	if !testing {
		return playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(false),
			SlowMo:   playwright.Float(500),
		}
	}
	return playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	}
}

func performLogin(page playwright.Page, email string, pwd string) error {
	if err := page.GetByPlaceholder("Email, phone, or Skype").Fill(email); err != nil {
		return fmt.Errorf("could not fill username: %w", err)
	}

	if err := page.Locator("input.win-button[value='Next']").Click(); err != nil {
		return fmt.Errorf("could not click next button: %w", err)
	}

	if err := page.GetByPlaceholder("Password").Fill(pwd); err != nil {
		return fmt.Errorf("could not fill password: %w", err)
	}

	if err := page.Locator("input.win-button[value='Sign in']").Click(); err != nil {
		return fmt.Errorf("could not click sign in button: %w", err)
	}
	return nil
}

func NewBrowser() (playwright.Browser, *playwright.Playwright, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to run Playwright: %w", err)
	}

	browser, err := pw.Firefox.Launch(headlessBrowserOptions(false))
	if err != nil {
		pw.Stop()
		return nil, nil, fmt.Errorf("could not run Firefox: %w", err)
	}

	return browser, pw, nil
}

func NavigateToLogin(page playwright.Page, baseUrl string) error {
	if _, err := page.Goto(baseUrl); err != nil {
		return fmt.Errorf("could not go to %s: %w ", baseUrl, err)
	}

	if err := page.Locator("a.ups-settings-link:has-text('Settings')").Click(); err != nil {
		return fmt.Errorf("could not click the settings button: %w", err)
	}

	if err := page.Locator("a.ui-btn:has-text('Sign in')").Click(); err != nil {
		return fmt.Errorf("could not click the sign in button: %w", err)
	}

	return nil
}

func StateFileExpireLogin(page playwright.Page, baseURL string) error {
	err := NavigateToLogin(page, baseURL)
	if err != nil {
		return fmt.Errorf("could not navigate to state file login: %w", err)
	}

	pickAccount := page.Locator("div[role='heading']").Filter(
		playwright.LocatorFilterOptions{HasText: "Pick an account"})

	err = pickAccount.WaitFor(playwright.LocatorWaitForOptions{
		Timeout: playwright.Float(10000),
	})

	if err == nil {
		return fmt.Errorf("===state file expired, closing context===")
	}
	return nil
}

func StateFileSucceedLogin(page playwright.Page, baseURL string) error {
	err := NavigateToLogin(page, baseURL)
	if err != nil {
		return fmt.Errorf("could not navigate to login: %w", err)
	}

	popUpSignIn := page.Locator("div[role='heading']:has-text('Stay signed in?')")

	err = popUpSignIn.WaitFor(playwright.LocatorWaitForOptions{
		Timeout: playwright.Float(3000),
	})

	if err == nil {
		log.Printf("Succeeded using state file. Proceed...")
		if err := page.Locator("input.win-button[value='No']").Click(); err != nil {
			return fmt.Errorf("cannot click no: %w", err)
		}
	}
	return nil
}

func FreshLogin(page playwright.Page, baseUrl string, email string, pwd string) error {
	log.Printf("===START FRESH LOGIN===")

	err := NavigateToLogin(page, baseUrl)
	if err != nil {
		return err
	}

	err = performLogin(page, email, pwd)
	if err != nil {
		return err
	}

	faAuth := page.Locator("#idDiv_SAOTCAS_Title")
	if err := faAuth.WaitFor(playwright.LocatorWaitForOptions{
		Timeout: playwright.Float(3000),
	}); err != nil {
		return fmt.Errorf("2FA Authentication does not show up: %w", err)
	}

	log.Println("===2FA Authentication REQUIRED===")
	log.Println("Waiting for entering 2FA code OR timeout (55s)...")

	locator := page.Locator("div#idRichContext_DisplaySign.displaySign")
	code, err := locator.TextContent()
	if err != nil {
		log.Fatalf("Failed to get 2FA code text: %v", err)
	}
	fmt.Println("2FA code is:", code)

	faAuthNext := page.GetByText("Stay signed in?")
	err = faAuthNext.WaitFor(playwright.LocatorWaitForOptions{
		Timeout: playwright.Float(55000),
	})

	if err != nil {
		return fmt.Errorf("2FA is not fulfilled: %w", err)
	}

	log.Println("FRESH LOGIN: Stay signed in appeared. Proceed..")
	if err := page.Locator("input.win-button[value='No']").Click(); err != nil {
		return fmt.Errorf("cannot click no: %w", err)
	}

	if _, err := page.Context().StorageState("ms_user.json"); err != nil {
		log.Printf("cannot save state file: %v", err)
	}

	return nil
}

func BookCourse(page playwright.Page, choices []config.ActivitySlot) error {
	var matchResult []config.ActivitySlot
	for _, c := range choices {

		if err := page.Locator("a.ui-btn.ui-btn-icon-right:has-text('Courses')").Click(); err != nil {
			return fmt.Errorf("could not find and click courses: %w ", err)
		}

		_, err := page.Locator("select#type").SelectOption(playwright.SelectOptionValues{
			Labels: playwright.StringSlice(c.DisplayCourseOption(c.Activity)),
		})
		if err != nil {
			return fmt.Errorf("cannot choose game selection %w", err)
		}

		_, err = page.Locator("select#area").SelectOption(playwright.SelectOptionValues{
			Labels: playwright.StringSlice(c.DisplayCourseArea(c.CourseArea)),
		})
		if err != nil {
			return fmt.Errorf("cannot choose game area: %w", err)
		}

		matchResult, err = matchSlot(page, c)
		if err != nil {
			if errors.Is(err, ErrNoTimeSlots) {
				log.Printf("[WARN] No time slots found for the current week (no <li> elements). Trying next week.")
				err := page.Locator("a.ui-btn:has-text('Next week')").Click()
				if err != nil {
					return fmt.Errorf("cannot click next week: %w", err)
				}
				matchResult, err = matchSlot(page, c)
				if err != nil {
					return fmt.Errorf("the book choices doesn't exist")
				}
			} else {
				return err
			}
		}

		if len(matchResult) == 0 {
			log.Printf("[WARN] Time slots found, but none match the intended criteria. Trying next week.")
			err := page.Locator("a.ui-btn:has-text('Next week')").Click()
			if err != nil {
				return fmt.Errorf("cannot click next week: %w", err)
			}
			matchResult, err = matchSlot(page, c)
			if err != nil {
				return fmt.Errorf("the book choices doesn't exist")
			}
		}
		log.Printf("Match result: %v\n", matchResult)

		for _, v := range matchResult {
			locator := fmt.Sprintf("li:has(a:text-is('%s %s')):has(span:text-is('%s %s'))", v.Hour, v.Activity, v.Day, v.Date)

			err = page.Locator(locator).Click()
			if err != nil {
				return fmt.Errorf("cannot click the exact schedule: %v", err)
			}

			if err := tryBookCourt(page, 6); err != nil {
				return fmt.Errorf("failed to book any court: %v", err)
			}

			err = page.Locator("div.ups-dialog-content:text-is('Thank you for your booking!')").WaitFor(playwright.LocatorWaitForOptions{
				Timeout: playwright.Float(2000),
			})
			if err == nil {
				if err = page.Locator("a.ui-btn:text-is('Ok')").Click(); err != nil {
					log.Printf("Close booking dialog")
					return nil
				}
			} else {
				return fmt.Errorf("cannot close booking dialog: %v", err)
			}
			locator = fmt.Sprintf("div:has(h1:text-is('%s')) a.ui-btn.ui-corner-all.ui-icon-delete[role='button']:has-text('Close')", c.DisplaySportDialogMap(c.Activity))
			err = page.Locator(locator).First().Click()
			if err != nil {
				return fmt.Errorf("failed closing booking court windows: %v", err)
			}
		}
	}
	return nil
}

func tryBookCourt(page playwright.Page, maxCourts int) error {
	for i := 1; i <= maxCourts; i++ {
		locator := fmt.Sprintf("a.ui-link:has-text('Book court %d')", i)
		err := page.Locator(locator).WaitFor(playwright.LocatorWaitForOptions{
			Timeout: playwright.Float(1000),
		})
		if err == nil {
			if err := page.Locator(locator).Click(); err == nil {
				log.Printf("Successfully booked court %d", i)
				return nil
			}
		}
	}
	return fmt.Errorf("no courts available (tried 1-%d)", maxCourts)
}

func tryReserveCourt(page playwright.Page) error {
	err := page.Locator("a.ui-link:has-text('Reserve')").WaitFor(playwright.LocatorWaitForOptions{
		Timeout: playwright.Float(1000),
	})
	if err == nil {
		if err := page.Locator("a.ui-link:has-text('Reserve')").Click(); err == nil {
			log.Printf("Successfully reserved court")
			return nil
		}
	}

	return fmt.Errorf("cannnot reserve court")
}

func matchSlot(page playwright.Page, choice config.ActivitySlot) ([]config.ActivitySlot, error) {
	// ==For Headless mode==
	err := page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	})
	if err != nil {
		log.Printf("Warning: Network not idle: %v", err)
	}
	// ================

	hourSlots, err := page.Locator("li:has(span)").All()
	if err != nil {
		return nil, fmt.Errorf("cannot find elements with date to book: %w", err)
	}

	var scrapedSlots []config.ActivitySlot
	slotMap := make(map[string]config.ActivitySlot)
	var matchResult []config.ActivitySlot

	for _, slot := range hourSlots {
		hourText, err := slot.TextContent()
		if err != nil {
			continue
		}
		//Format: "Wed 11.6. 20:00 Badminton"
		parts := strings.Fields(hourText)
		if len(parts) < 4 {
			continue
		}
		playslot := config.ActivitySlot{
			Day:      parts[0],
			Date:     parts[1],
			Hour:     parts[2],
			Activity: parts[3],
		}
		scrapedSlots = append(scrapedSlots, playslot)
	}

	for _, slot := range scrapedSlots {
		key := fmt.Sprintf("%v|%v|%v", slot.Day, slot.Hour, slot.Activity)
		slotMap[key] = slot
	}

	key := fmt.Sprintf("%v|%v|%v", choice.Day, choice.Hour, choice.Activity)
	if slot, exists := slotMap[key]; exists {
		matchResult = append(matchResult, slot)
	}

	return matchResult, nil
}
