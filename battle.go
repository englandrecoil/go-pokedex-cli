package main

import (
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/englandrecoil/go-pokedex-cli/internal/pokeapi"
	"github.com/fatih/color"
)

/* TODO:
- добавить сюда жеребьевку того, кто первый ходит(можно определять это по
опыту покемона)
- добавить вывод логов битвы
*/

func startBattle(firstContestant, secondContestant pokeapi.Battler) error {
	color.Unset()
	defer color.Unset()

	const treshold = 30

	for firstContestant.Health > 0 && secondContestant.Health > 0 {
		damageFirst := rand.IntN(int(math.Round(float64(firstContestant.Attack*secondContestant.Defense) / 100)))
		damageSecond := rand.IntN(int(math.Round(float64(secondContestant.Attack*firstContestant.Defense) / 100)))

		chanceToAttackFirst := rand.IntN(firstContestant.Experience) + treshold
		chanceToAttackSecond := rand.IntN(secondContestant.Experience) + treshold

		if chanceToAttackFirst > secondContestant.Parry {
			if secondContestant.Health -= damageFirst; secondContestant.Health <= 0 {
				fmt.Printf("%s attacked! %s's health is 0\n", firstContestant.Name, secondContestant.Name)
				color.Set(color.FgBlue)
				fmt.Printf("%s is the WINNER!\n", firstContestant.Name)
				return nil
			}
			fmt.Printf("%s attacked! %s's health is %d\n", firstContestant.Name, secondContestant.Name, secondContestant.Health)
		} else {
			fmt.Printf("%s missed\n", firstContestant.Name)
		}

		if chanceToAttackSecond > firstContestant.Parry {
			if firstContestant.Health -= damageSecond; firstContestant.Health <= 0 {
				fmt.Printf("%s attacked! %s's health is 0\n", secondContestant.Name, firstContestant.Name)
				color.Set(color.FgBlue)
				fmt.Printf("%s is the WINNER!\n", secondContestant.Name)
				return nil
			}
			fmt.Printf("%s attacked! %s's health is %d\n", secondContestant.Name, firstContestant.Name, firstContestant.Health)
		} else {
			fmt.Printf("%s missed\n", secondContestant.Name)

		}

	}

	return nil
}
