package runner

// import (
// 	// "context"
// 	// "fmt"
// 	// "os/exec"
// 	"path/filepath"
// 	// "time"

// 	"github.com/barrydevp/codeatest-runner-core/puller"
// 	"github.com/barrydevp/codeatest-runner-core/model"
// )

// type Golang struct {
//     BaseRunner
// }

// func (g *Golang) LoadData(data *puller.Data) error {
//     if data == nil {
//         return NewError(NilData)
//     }

//     g.Data = data

//     return nil
// }

// func NewGolang() Runner {

// 	abs, err := filepath.Abs("./tests/hello.go")

//     if err != nil {
//         return nil
//     }

//     runner := &Golang{
//         BaseRunner{
//             "golang",
//             "created",
//             &puller.Data{
//                 model.Job{},
//                 model.Quiz{},
//                 model.TestCase{},
//                 abs,
//             },
//             nil,
//             nil,
//             nil,
//             []string{"go", "run"},
//         },
//     }

//     return runner
// }
