package pattern

// import (
// 	"fmt"
// 	"sync"
// )

// type singletonMap map[string]Singleton

// type Singleton interface
// {

// }

// func init() {

// 	singletonMap = make(map[string]Singleton)
// 	fmt.Println("api init all")
// 	singletonMap["logger"] = logger{}
// 	singletonMap["signUp"] = signUpHandler{}
// 	singletonMap["addCF"] = NewAddCFApi()
// 	singletonMap["addTaskPlan"] = NewAddTaskPlanApi()
// 	singletonMap["addTask"] = NewAddTaskApi()
// }

// type foo struct {
//     Val int
// }
// var instance *foo
// var once sync.Once

// func GetSingleton(name string) *Singleton {

//     once.Do(func () {
//         instance = &foo{1}
//     })
//     return instance
// }