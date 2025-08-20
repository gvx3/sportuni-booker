# TUNI Sportuni Booking

A tool to automate booking your favorite sports activities at TUNI using *Playwright* (Go).

> [!NOTE]  
> TUNI Sportuni uses Microsoft 2FA so it is impossible to achieve full automation, instead a state file will be generated and live for a short period (a few days) to keep the login session.

## Prerequisites

- Go version [1.24.1+](https://go.dev/dl/) or newer must be installed.

## Getting Started

1. **Clone the repository:**

   ```sh
   git clone https://github.com/gvx3/sportuni-book.git
   cd sportuni-booker
   ```

2. **Install dependencies and Playwright browsers:**

   ```sh
   make setup
   ```

3. **At the root directory of the project, build the binary:**

    ```sh
    go build -o sb
    ```

4. **Run application**

    ```sh
    ./sb
    # or specify a config file
    ./sb -f /path/to/config.yaml
    ```

## Demo

![Demo of SportUni booking automation](/asset/demo.gif)

## Configuration

The application requires a YAML configuration file. The config file can be located in one of the following places:

- The current working directory (where you run the program):
  - `./config.yaml`
- The user home directory:
  - On Unix/MacOS: `~/.sportuni/config.yaml`
  - On Windows: `C:\Users\<YourUsername>\.sportuni\config.yaml`

You can also specify a custom config file location using the CLI flag `-f` or `--file` if the path could not be found for some reason:

```sh
./sb -f /path/to/your/config.yaml
```

## Configuration values

Available configuration can be found in the following table:

- base_url: "https://www.tuni.fi/sportuni/omasivu/?newPage=selection&lang=en" #required
- email: "your.email@tuni.fi" #required
- password: "your_password" #required
- state_file_name: "ms_user.json" #optional

| Key    | Values | Required? |
| -------- | ------- | ------- |
| base_url  | https://www.tuni.fi/sportuni/omasivu/?newPage=selection&lang=en   | Yes |
| email | your.email@tuni.fi     | Yes |
| password    | your_password    | Yes |
| state_file_name    | ms_user.json    | No |
| day    | Mon, Tue, Wed, Thu, Fri, Sat, Sun    | Yes |
|   date  |  Format: day.month. (Ex: 11.6.) | No |
|   activity  |  Badminton, Billiards | Yes |
|   course_area  |  hervanta, kauppi, citycentre | Yes |

### Example `config.yaml`

```yaml
base_url: "https://www.tuni.fi/sportuni/omasivu/?newPage=selection&lang=en"
email: "your.email@tuni.fi"
password: "your_password"
state_file_name: "ms_user.json"

activity_slots:
  - day: "Tue"
    date: ""
    hour: "08:00"
    activity: "Badminton"
    course_area: "hervanta"
  - day: "Thu"
    date: ""
    hour: "13:00"
    activity: "Billiards"
    course_area: "citycentre"
```
