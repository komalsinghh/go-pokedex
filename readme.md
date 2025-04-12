# Go Pokedex

Go Pokedex is a command-line application that allows users to explore Pokémon locations, catch Pokémon, and manage their own Pokedex. It interacts with the [PokéAPI](https://pokeapi.co/) to fetch data about Pokémon and their habitats.

## Features

- **Explore Locations**: View the next or previous 20 Pokémon locations.
- **Explore Pokémon in a Location**: List all Pokémon available in a specific location.
- **Catch Pokémon**: Attempt to catch a Pokémon and add it to your Pokedex.
- **View Pokedex**: Display all the Pokémon you have caught.
- **Inspect Pokémon**: View detailed stats and information about a specific Pokémon.
- **Clear Screen**: Clear the terminal screen.
- **Help**: Display a list of available commands.
- **Exit**: Exit the application.

## Commands

| Command       | Description                                   |
|---------------|-----------------------------------------------|
| `map`         | Displays the next 20 locations.              |
| `mapb`        | Displays the previous 20 locations.          |
| `explore`     | Displays a list of all Pokémon in a location. |
| `catch`       | Try to catch a Pokémon.                      |
| `pokedex`     | Displays your caught Pokémon.                |
| `inspect`     | Inspect a specific Pokémon.                  |
| `clear`       | Clears the terminal screen.                  |
| `help`        | Displays a help message.                     |
| `exit`        | Exits the application.                       |

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/komalsinghh/go-pokedex.git
    cd go-pokedex
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

3. Run the application:
    ```bash
    go run main.go
    ```

## Usage

Once the application is running, you can type any of the commands listed above to interact with the Pokedex. For example:

- To explore the next 20 locations:
  ```
  map
  ```

- To catch a Pokémon:
  ```
  catch pikachu
  ```

- To view your Pokedex:
  ```
  pokedex
  ```

- To inspect a Pokémon:
  ```
  inspect pikachu
  ```
## Technologies Used

- **Go**: Backend programming language.
- **PokéAPI**: Public API for Pokémon data.

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch:
    ```bash
    git checkout -b feature-name
    ```
3. Commit your changes:
    ```bash
    git commit -m "Add feature-name"
    ```
4. Push to the branch:
    ```bash
    git push origin feature-name
    ```
5. Open a pull request.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgments

- [PokéAPI](https://pokeapi.co/) for providing Pokémon data.
- The Go community for their excellent resources and support.
