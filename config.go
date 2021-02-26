package main

import (
	"encoding/json"
	"gioman/state"
	"os"
)

// [
//     {
//         name: 'qeweqqw',
//         url: 'weqweqweqwe',
//         method: 0
//     }
// ]
// TODO: Need to validate the config, if `method: 921`, it's not valid

func readConfig() (requests state.Requests, err error) {
	// TODO change the filename path, look for the config folder of the system?
	file, err := os.Open("config.json")
	if err != nil {
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&requests)
	if err != nil {
		return
	}
	return
}
