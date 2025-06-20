# go-pokedex-cli
[![en](https://img.shields.io/badge/lang-en-green?style=flat)](https://github.com/englandrecoil/go-pokedex-cli/blob/main/README.md)
[![ru](https://img.shields.io/badge/lang-ru-blue?style=flat)](https://github.com/englandrecoil/go-pokedex-cli/blob/main/README.ru.md#go-pokedex-cli)

**Pokedex CLI** lets you explore the Pokémon world, discover new locations, and catch and track Pokémon in your Pokedex. It uses caching, PokeAPI to fetch data, and generates ASCII images of the Pokémon you catch - all in terminal!

<img src="https://iili.io/2g8pbwu.jpg" width="550">

## :minidisc: Installation
To use, follow these steps:
1. Paste `git clone https://github.com/englandrecoil/go-pokedex-cli.git` in your terminal
2. Compile it yourself with `go build`(you need [go](https://go.dev)) or just use the compiled one with `go-pokedex-cli` in your terminal.

## :keyboard: Usage
Try writing the 'help' command to find out what commands are available. 
The list of available commands can also be found additionally below:
| Command  | Description |
| ------------- | ------------- |
| `pokedex`  | Displays all caught Pokémon |
| `map`  | Displays the names of the next 20 location areas |
| `mapb` | Displays the names of the next 20 location areas |
| `explore {location_area}` | Displays all the Pokémon in a given area |
| `inspect {pokemon_name}` | Inspect the caught pokemon |
| `catch {pokemon_name}` | Catch Pokemon with a certain chance |
| `battle {pokemon_name1} {pokemon_name2}` | Simulate battles between two captured Pokémon |
| `help` | Displays a help message |
| `exit` | Exit the Pokedex |
| `clear` | Clear the terminal screen |
| `cache {integer_number}` | Set the caching interval(in hours) after which cleaning will occur |
| `color {on/off}` | Configures the display of color output* |

\* To comply with the [standard](https://no-color.org) and not confuse users, it only works if the environment variable 'NO_COLOR' is empty. By default, it is set to the value NO_COLORS. If you haven't touched this variable, you're all set.

## :spiral_notepad: Future improvements and enhancements
- [X] Simulate battles between captured Pokémon
- [X] Save progress between sessions by saving the user's Pokédex to disk
- [X] Improve ASCII art generation
- [X] Add 'clear' command support for Windows
