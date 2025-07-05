# TUNI Sportuni Booking

A tool to automate booking your favorite sports activities at TUNI using *Playwright* (Go).

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

## Configuration

The application requires a YAML configuration file. The config file can be located in one of the following places:

- The current working directory (where you run the program):
  - `./config.yaml`
- The user home directory:
  - On Unix/MacOS: `~/.sportuni/config.yaml`
  - On Windows: `C:\Users\<YourUsername>\.sportuni\config.yaml`

You can also specify a custom config file location using the CLI flag `-f` or `--file`:

```sh
./sb -f /path/to/your/config.yaml
```

### Example `config.yaml`

```yaml
base_url: "https://www.tuni.fi/sportuni/omasivu/?newPage=selection&lang=en"
email: "your.email@tuni.fi"
password: "your_password"
state_file_name: "ms_user.json"

# activity values: Badminton, Billiards
# course_area values: hervanta, kauppi, citycentre
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
