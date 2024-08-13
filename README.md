# Pokedex CLI

Welcome to the Pokedex CLI project! This is a command-line application that allows you to explore the world of Pokemon, catch Pokemon, and manage your Pokedex. It is developed as part of the Boot.Dev Backend Development course.

## Features

- Explore different location areas in the Pokemon world
- Catch Pokemon and add them to your Pokedex
- View details of caught Pokemon, including their stats and types
- Simulate battles between caught Pokemon
- Manage a party of Pokemon and allow them to level up
- Evolve caught Pokemon after a set amount of time
- Persist your Pokedex progress between sessions
- Encounter wild Pokemon randomly while exploring
- Use different types of Poke Balls to catch Pokemon

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/pokedex-cli.git
   ```

2. Navigate to the project directory:
   ```
   cd pokedex-cli
   ```

3. Install the dependencies:
   ```
   go mod download
   ```

## Usage

To start the Pokedex CLI, run the following command:
```
go run main.go
```

Once the application is running, you can use the following commands:

- `help`: Displays a list of available commands and their descriptions.
- `map`: View a map of the Pokemon location areas.
- `explore <area>`: Explore a specific location area and encounter wild Pokemon.
- `catch <pokemon>`: Attempt to catch a Pokemon and add it to your Pokedex.
- `inspect <pokemon>`: View detailed information about a caught Pokemon.
- `pokedex`: Show a list of all caught Pokemon in your Pokedex.
- `party`: Manage your party of Pokemon.
- `battle <pokemon1> <pokemon2>`: Simulate a battle between two caught Pokemon.
- `save`: Manually save your Pokedex progress.
- `exit`: Exit the Pokedex CLI.

## Configuration

The Pokedex CLI uses the PokeAPI (https://pokeapi.co/) to fetch Pokemon data. Make sure you have an active internet connection to access the API.

You can configure the application by modifying the constants and variables in the `main.go` file, such as the cache duration, default location area, and more.

## Testing

The project includes unit tests to ensure the correctness and reliability of the code. To run the tests, use the following command:
```
go test ./...
```

## Contributing

Contributions to the Pokedex CLI project are welcome! If you find any bugs, have suggestions for improvements, or want to add new features, please open an issue or submit a pull request.

When contributing, please follow the existing code style and conventions, and make sure to update the tests and documentation accordingly.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgements

- [PokeAPI](https://pokeapi.co/) - The API used to fetch Pokemon data.
- [Boot.Dev](https://boot.dev/) - The platform that provided the Backend Development course and project inspiration.
