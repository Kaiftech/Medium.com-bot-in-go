package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/tebeka/selenium"
)

type BotConfig struct {
	Browser string
}

var pastInteractions []string

func main() {
	// Start the terminal interface
	terminalInterface()
}

func terminalInterface() {
	config := BotConfig{Browser: "Chrome"}
	launchBot(config)
}

func launchBot(config BotConfig) {
	fmt.Println("Launching bot with the following configuration:")
	fmt.Printf("Browser: %s\n", config.Browser)

	// Set up Selenium WebDriver
	seleniumPath := "" // Path to ChromeDriver
	port := 8081

	svc, err := selenium.NewChromeDriverService(seleniumPath, port)
	if err != nil {
		log.Fatalf("Error starting the ChromeDriver server: %v", err)
	}
	defer svc.Stop()

	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeArgs := []string{
		"--disable-blink-features=AutomationControlled",
		"--disable-infobars",
		"--start-maximized",
	}
	caps["goog:chromeOptions"] = map[string]interface{}{
		"args": chromeArgs,
	}

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		log.Fatalf("Error connecting to the WebDriver server: %v", err)
	}
	defer wd.Quit()

	// Open Medium.com and sign in
	err = signIn(wd)
	if err != nil {
		log.Fatal(err)
	}

	// Wait for user to confirm login
	for {
		fmt.Print("Please confirm if you have successfully logged in (yes/no): ")
		var input string
		fmt.Scanln(&input)
		if input == "yes" {
			break
		} else {
			fmt.Println("Waiting for you to complete the login...")
		}
	}

	// Interact with articles
	err = searchAndInteract(wd)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bot completed its tasks successfully!")
}

func signIn(wd selenium.WebDriver) error {
	fmt.Println("Signing in to Medium.com...")

	// Open Medium sign-in page
	signInURL := "https://medium.com/m/signin"
	if err := wd.Get(signInURL); err != nil {
		return err
	}

	// Wait for the "Sign in with Google" button to appear
	err := waitForElement(wd, selenium.ByCSSSelector, "button[jsname='Cuz2Ue']", 60*time.Second) // Update the selector to target the Google button reliably
	if err != nil {
		return fmt.Errorf("failed to find Google sign-in button: %v", err)
	}

	// Scroll to the Google sign-in button
	googleSignInButton, err := wd.FindElement(selenium.ByCSSSelector, "button[jsname='Cuz2Ue']")
	if err != nil {
		return fmt.Errorf("could not find Google sign-in button: %v", err)
	}

	_, err = wd.ExecuteScript("arguments[0].scrollIntoView(true);", []interface{}{googleSignInButton})
	if err != nil {
		return fmt.Errorf("could not scroll to Google sign-in button: %v", err)
	}

	time.Sleep(1 * time.Second) // Allow time for scroll action

	// Retry mechanism for clicking the Google sign-in button
	for attempts := 0; attempts < 5; attempts++ {
		err = googleSignInButton.Click()
		if err == nil {
			break
		}
		fmt.Printf("Attempt %d: could not click on Google sign-in button, retrying...\n", attempts+1)
		time.Sleep(2 * time.Second) // Wait before retrying
	}
	if err != nil {
		// If click still fails, use JavaScript to click
		_, jsErr := wd.ExecuteScript("arguments[0].click();", []interface{}{googleSignInButton})
		if jsErr != nil {
			return fmt.Errorf("could not click on Google sign-in button even with JavaScript: %v", jsErr)
		}
	}

	fmt.Println("Please manually enter your Google email and password in the browser window.")
	return nil
}

func searchAndInteract(wd selenium.WebDriver) error {
	fmt.Println("Searching for articles on Medium.com...")

	// Go to Medium homepage
	homepageURL := "https://medium.com"
	if err := wd.Get(homepageURL); err != nil {
		return fmt.Errorf("could not navigate to Medium homepage: %v", err)
	}

	// Wait for the search input to appear
	err := waitForElement(wd, selenium.ByCSSSelector, "input[data-testid='headerSearchInput']", 30*time.Second)
	if err != nil {
		return fmt.Errorf("failed to find search input: %v", err)
	}

	// Enter a search query
	searchInput, err := wd.FindElement(selenium.ByCSSSelector, "input[data-testid='headerSearchInput']")
	if err != nil {
		return fmt.Errorf("could not find search input: %v", err)
	}
	err = searchInput.SendKeys("technology")
	if err != nil {
		return fmt.Errorf("could not enter search query: %v", err)
	}

	searchInput.SendKeys(selenium.EnterKey) // Press Enter to search
	time.Sleep(2 * time.Second)             // Wait for search results to load

	// Click on the first article link with a more generic selector
	articleSelectors := []string{
		"h2",
		"h3",
	}
	var articleLink selenium.WebElement
	found := false
	for _, selector := range articleSelectors {
		articleHeadings, err := wd.FindElements(selenium.ByCSSSelector, selector)
		if err == nil && len(articleHeadings) > 0 {
			for _, heading := range articleHeadings {
				link, err := heading.FindElement(selenium.ByXPATH, "ancestor::a")
				if err == nil {
					articleLink = link
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("could not find an article heading using any of the provided selectors")
	}

	_, err = wd.ExecuteScript("arguments[0].scrollIntoView(true);", []interface{}{articleLink})
	if err != nil {
		return fmt.Errorf("could not scroll to article link: %v", err)
	}
	err = articleLink.Click()
	if err != nil {
		return fmt.Errorf("could not click on article link: %v", err)
	}

	// Wait for the article page to load
	err = waitForElement(wd, selenium.ByCSSSelector, "button[data-testid='headerClapButton']", 30*time.Second)
	if err != nil {
		return fmt.Errorf("failed to load article page: %v", err)
	}

	// Interact with the article
	err = interactWithArticle(wd)
	if err != nil {
		return err
	}

	return nil
}

func interactWithArticle(wd selenium.WebDriver) error {
	fmt.Println("Interacting with the article...")

	// Clap for the article multiple times
	for i := 0; i < rand.Intn(10)+1; i++ {
		fmt.Println("Clapping for the article...")
		err := clapForArticle(wd)
		if err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond) // Short sleep between claps
	}

	/*

		// Randomly comment on the article
		if rand.Float32() < 0.4 {
			fmt.Println("Commenting on the article...")
			err := commentOnArticle(wd)
			if err != nil {
				return err
			}
			time.Sleep(time.Duration(rand.Intn(5)+2) * time.Second) // Random sleep after commenting
		}

	*/

	return nil
}

func clapForArticle(wd selenium.WebDriver) error {
	logInteraction("Clapped for an article.")
	clapButton, err := wd.FindElement(selenium.ByCSSSelector, "button[data-testid='headerClapButton']")
	if err != nil {
		return err
	}
	return clapButton.Click()
}

/*

func commentOnArticle(wd selenium.WebDriver) error {
	comments := []string{
		"Great read! Really enjoyed the insights.",
		"Very informative article, thanks for sharing!",
		"Loved the depth of information provided here.",
		"This was exactly what I was looking for. Great job!",
		"Interesting perspective! Thanks for writing this.",
		"Your article made me think. Thanks for sharing your thoughts!",
		"I appreciate the unique perspective presented here.",
	}
	comment := comments[rand.Intn(len(comments))]

	logInteraction(fmt.Sprintf("Commented on an article: '%s'", comment))

	respondButton, err := wd.FindElement(selenium.ByCSSSelector, "button[data-testid='ResponseRespondButton']")
	if err != nil {
		return err
	}
	err = respondButton.Click()
	if err != nil {
		return err
	}

	time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second) // Random sleep before typing

	commentBox, err := wd.FindElement(selenium.ByCSSSelector, "textarea")
	if err != nil {
		return err
	}
	err = commentBox.SendKeys(comment)
	if err != nil {
		return err
	}

	time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second) // Random sleep after typing

	publishButton, err := wd.FindElement(selenium.ByCSSSelector, "button[data-testid='ResponseRespondButton']")
	if err != nil {
		return err
	}
	return publishButton.Click()
}

*/

func waitForElement(wd selenium.WebDriver, by, value string, timeout time.Duration) error {
	for start := time.Now(); time.Since(start) < timeout; {
		_, err := wd.FindElement(by, value)
		if err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("element not found: %s %s", by, value)
}

func logInteraction(interaction string) {
	pastInteractions = append(pastInteractions, interaction)
}
