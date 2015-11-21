package base

import (
// "encoding/json"
// "fmt"
)

type user struct {
	Name string `json:"name"`
}

type Config struct {
	User user `json:"user"`
}

// func main() {
// 	fmt.Println(Config{User: user{Name: "Nishanth"}})
// 	var b []byte
// 	b, _ = json.Marshal(Config{User: user{Name: "Nishanth"}})
// 	fmt.Println(string(b))
// 	var c Config
// 	json.Unmarshal(b, &c)
// 	fmt.Println(c)
// }
