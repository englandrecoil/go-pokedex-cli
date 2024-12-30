package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"time"

	"github.com/englandrecoil/go-pokedex-cli/internal/pokeapi"
	"github.com/fatih/color"
)

func startBattle(firstContestant, secondContestant pokeapi.Battler) error {
	color.Unset()
	defer color.Unset()

	const treshold = 30

	for firstContestant.Health > 0 && secondContestant.Health > 0 {
		damageFirst := rand.IntN(int(math.Round(float64(firstContestant.Attack*secondContestant.Defense) / 100)))
		damageSecond := rand.IntN(int(math.Round(float64(secondContestant.Attack*firstContestant.Defense) / 100)))

		chanceToAttackFirst := rand.IntN(firstContestant.Experience) + treshold
		chanceToAttackSecond := rand.IntN(secondContestant.Experience) + treshold

		time.Sleep(800 * time.Millisecond)
		if chanceToAttackFirst > secondContestant.Parry {
			if secondContestant.Health -= damageFirst; secondContestant.Health <= 0 {
				fmt.Printf("%s attacked! %s's health is 0\n", firstContestant.Name, secondContestant.Name)
				color.Set(color.FgGreen)
				fmt.Printf("%s is the WINNER!\n", firstContestant.Name)
				return nil
			}
			fmt.Printf("%s attacked! %s's health is %d\n", firstContestant.Name, secondContestant.Name, secondContestant.Health)
		} else {
			fmt.Printf("%s missed\n", firstContestant.Name)
		}

		time.Sleep(800 * time.Millisecond)
		if chanceToAttackSecond > firstContestant.Parry {
			if firstContestant.Health -= damageSecond; firstContestant.Health <= 0 {
				fmt.Printf("%s attacked! %s's health is 0\n", secondContestant.Name, firstContestant.Name)
				color.Set(color.FgHiGreen)
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
