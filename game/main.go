package main

import (
	"strings"
)

type gamer struct {
	currRoom  string
	inventory []string
}
type jobParms struct {
	prefix  string
	postfix string
}
type room struct {
	jobs      map[string]jobParms
	exits     []string
	inventory map[string][]string
}

var user gamer
var door, bag bool
var kitchen, hall, livingRoom, street room
var actions []string
var rooms map[string]room

func printList(list []string) string {
	var tmp string
	var len = len(list) - 1
	for i, elem := range list {
		tmp += elem
		if i != len {
			tmp += ", "
		}
	}
	return tmp
}
func findObject(list []string, object string) (bool, int) {
	for i, elem := range list {
		if elem == object {
			return true, i
		}
	}
	return false, -1
}
func findObjectInMap(currMap map[string][]string, object string) (bool, string) {
	for key, value := range currMap {
		ok2, _ := findObject(value, object)
		if ok2 == true {
			return true, key
		}
	}
	return false, ""
}

func main() {
	initGame()
}

func initGame() {
	door = false
	bag = false
	user.inventory = []string{}
	actions = []string{
		"осмотреться", "идти", "надеть", "взять", "применить",
	}
	user.currRoom = "кухня"

	kitchen = room{
		jobs: map[string]jobParms{
			"осмотреться": jobParms{
				prefix:  "ты находишься на кухне, ",
				postfix: "можно пройти - ",
			},
			"идти": jobParms{
				prefix:  "",
				postfix: "кухня, ничего интересного. можно пройти - ",
			},
		},
		exits: []string{
			"коридор",
		},
		inventory: map[string][]string{
			"стол": {"чай"},
		},
	}

	hall = room{
		jobs: map[string]jobParms{
			"идти": jobParms{
				prefix:  "",
				postfix: "ничего интересного. можно пройти - ",
			},
			"применить": jobParms{
				prefix:  "",
				postfix: "",
			},
		},
		exits: []string{
			"кухня", "комната", "улица",
		},
		inventory: map[string][]string{
			"дверь": {},
		},
	}

	livingRoom = room{
		jobs: map[string]jobParms{
			"осмотреться": jobParms{
				prefix:  "",
				postfix: ". можно пройти - ",
			},
			"идти": jobParms{
				prefix:  "",
				postfix: "ты в своей комнате. можно пройти - ",
			},
			"надеть": jobParms{
				prefix:  "",
				postfix: "вы надели: ",
			},
			"взять": jobParms{
				prefix:  "",
				postfix: "предмет добавлен в инвентарь: ",
			},
		},
		exits: []string{
			"коридор",
		},
		inventory: map[string][]string{
			"стол": {"ключи", "конспекты"},
			"стул": {"рюкзак"},
		},
	}

	street = room{
		jobs: map[string]jobParms{
			"идти": jobParms{
				prefix:  "",
				postfix: "на улице весна. можно пройти - ",
			},
		},
		exits: []string{
			"домой",
		},
	}

	rooms = make(map[string]room)
	rooms["кухня"] = kitchen
	rooms["коридор"] = hall
	rooms["комната"] = livingRoom
	rooms["улица"] = street
}

func handleCommand(command string) string {
	var split = strings.Split(command, " ")
	var result string
	var parms = split[1:]

	currentRoom := rooms[user.currRoom]
	currJobParm := currentRoom.jobs[split[0]]

	ok, _ := findObject(actions, split[0])
	if ok == false {
		return "неизвестная команда"
	}

	result += currJobParm.prefix
	switch split[0] {
	case "осмотреться":
		{
			_, ok1 := currentRoom.inventory["стол"]
			if ok1 == true {
				result += "на столе: "
				result += printList(currentRoom.inventory["стол"])
			}
			_, ok2 := currentRoom.inventory["стул"]
			if ok2 == true {
				if ok1 == true {
					result += ", "
				}
				result = result + "на стуле: "
				result += printList(currentRoom.inventory["стул"])
			}
			if user.currRoom == "кухня" {
				if bag == false {
					result += ", надо собрать рюкзак и идти в универ. "
				} else {
					result += ", надо идти в универ. "
				}
			}
			if ok1 == false && ok2 == false {
				result += "пустая комната"
			}
		}
	case "идти":
		{
			ok, _ := findObject(currentRoom.exits, parms[0])
			if ok == false {
				return "нет пути в " + parms[0]
			}
			if parms[0] == "улица" && door == false {
				return "дверь закрыта"
			}
			user.currRoom = parms[0]
			currentRoom = rooms[user.currRoom]
			currJobParm = currentRoom.jobs[split[0]]
		}
	case "применить":
		{
			ok1, _ := findObject(user.inventory, parms[0])
			if ok1 == false {
				return "нет предмета в инвентаре - " + parms[0]
			}
			if parms[0] == "ключи" && parms[1] == "дверь" {
				door = true
				return "дверь открыта"
			}
			ok2, _ := findObjectInMap(currentRoom.inventory, parms[1])
			if ok2 == false {
				return "не к чему применить"
			}
		}
	case "взять":
		{
			if bag == false {
				return "некуда класть"
			}
			ok1, n := findObjectInMap(currentRoom.inventory, parms[0])
			if ok1 == false {
				return "нет такого"
			}
			_, i := findObject(currentRoom.inventory[n], parms[0])
			if len(currentRoom.inventory[n]) == 1 {
				delete(currentRoom.inventory, n)
			} else {
				currentRoom.inventory[n] = append(currentRoom.inventory[n][:i], currentRoom.inventory[n][i+1:]...)
			}
		}
	}
	result += currJobParm.postfix
	switch split[0] {
	case "осмотреться":
		{
			result += printList(currentRoom.exits)
		}
	case "идти":
		{
			result += printList(currentRoom.exits)
		}
	case "надеть":
		{
			if parms[0] == "рюкзак" {
				bag = true
				delete(currentRoom.inventory, "стул")
			}
			result += parms[0]
		}
	case "взять":
		{
			user.inventory = append(user.inventory, parms[0])
			result += parms[0]
		}
	}

	return result
}
