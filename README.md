# Medium Automation Bot

## Overview

This project is an automation bot designed to interact with articles on Medium.com. The bot can clap for articles and is built using Go programming language along with the Selenium WebDriver for browser automation.

## Table of Contents

1. [Features](#features)
2. [Installation](#installation)
3. [Configuration](#configuration)
4. [Usage](#usage)
5. [Code Structure](#code-structure)
6. [License](#license)

## Features

- Automated interaction with Medium articles.
- Clap for articles multiple times.
- Simple terminal interface for configuration and interaction.

## Installation

### Prerequisites

Before you begin, ensure you have the following installed:

- Go (version 1.16 or later)
- Chrome web browser
- ChromeDriver (compatible with your Chrome version)

### Steps to Install

1. Clone the repository:

   ```bash
   git clone https://github.com/Kaiftech/medium-bot.git
   cd medium-bot
   ```

2. Install the required Go packages:

   ```bash
   go mod tidy
   ```

3. Ensure the path to the ChromeDriver is correctly set in the code.

## Configuration

Modify the `launchBot` function in the `main.go` file to set your ChromeDriver path:

```go
seleniumPath := "C:\\path\\to\\chromedriver.exe" // Update this path
```

## Usage

1. Open your terminal and navigate to the project directory.
2. Run the application:

   ```bash
   go run main.go
   ```

3. Follow the on-screen instructions to sign in to Medium.com. After logging in, type `yes` in the terminal to confirm successful login.

4. The bot will begin searching for articles and interacting with them.

## Code Structure

- `main.go`: The main entry point of the application.
- `BotConfig`: A struct that holds configuration details for the bot.
- `terminalInterface`: Starts the terminal interface and launches the bot.
- `launchBot`: Initializes the Selenium WebDriver and handles the sign-in process.
- `signIn`: Automates the sign-in process using the Google account.
- `searchAndInteract`: Searches for articles and interacts with them.
- `interactWithArticle`: Manages clapping for articles.
- `waitForElement`: Utility function to wait for specific elements to load.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Feel free to customize any sections as needed!
