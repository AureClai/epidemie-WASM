package main

import (
	"image/color"
	"syscall/js"
)

const SIZE float64 = 600

// Const definition before config file

func main() {

	// settings := new(SimulationSettings)
	settings := DefaultSettings()
	// Open our jsonFile
	// jsonFile, err := os.Open("settings.json")
	// // if we os.Open returns an error then handle it
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// defer jsonFile.Close()
	// read our opened jsonFile as a byte array.
	// byteValue, _ := ioutil.ReadAll(jsonFile)

	// json.Unmarshal(byteValue, settings)

	// walls := NewEmptyWallList()
	// walls = append(walls, NewWall(10, 10, 0, 14, 0.2))
	// walls = append(walls, NewWall(10, 10, 16, 30, 0.2))

	// startAgParams := make([](*StartAgentsParameters), 0)
	// startAgParams = append(startAgParams, &StartAgentsParameters{
	// 	Position: Vect2{2, 2},
	// 	Speed:    Vect2{0, 3},
	// 	State:    1,
	// 	Movable:  true,
	// })

	// settings := SimulationSettings{
	// 	Walls:               walls,
	// 	WindowSizeX:         30.0,
	// 	WindowSizeY:         30.0,
	// 	Duration:            30.0,
	// 	Dt:                  1 / 60,
	// 	TimeToRecover:       5.0,
	// 	FracRandomUnmovable: 0.25,
	// 	NbRandomAgents:      200,
	// 	NbRandomSicks:       1,
	// 	StartAgParam:        startAgParams,
	// 	PDeath:              0.01,
	// 	AgentRadius:         0.2,
	// 	AgentStartSpeed:     3,
	// }
	// // TO DELETE
	// jsonFile, _ := json.MarshalIndent(settings, "", " ")
	// _ = ioutil.WriteFile("settings.json", jsonFile, 0644)

	simulation := NewSimulation(settings)
	//simulation.Run()

	cv, err := NewCanvas(simulation.Settings.WindowSizeX, SIZE)
	if err != nil {
		println("Impossible to create canvas")
	}

	//cv.ClearGC(color.RGBA{0xff, 0xff, 0xff, 0xff})

	var loop js.Func
	loop = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		simulation.NextStep()
		cv.ctx.Call("clearRect", 0, 0, cv.width, cv.height)

		// Make solid white background
		cv.ClearGC(color.RGBA{0xff, 0xff, 0xff, 0xff})

		// Walls
		for _, wall := range simulation.Walls {
			wall.Draw(cv)
		}

		//Agents
		for _, agent := range simulation.Agents {
			agent.Draw(cv)
		}

		cv.Render()
		cv.window.Call("requestAnimationFrame", loop)
		return nil
	})

	cv.window.Call("requestAnimationFrame", loop)

	<-make(chan bool)

	// simulation.SaveResults()
}
