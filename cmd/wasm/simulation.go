package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var simulationDuration float64 = 30 // seconds
var dt float64 = 1.0 / 60.0

type Simulation struct {
	Agents      AgentList
	AliveAgents AgentList
	Walls       WallList
	Infos       []string
	Settings    *SimulationSettings
	Time        float64
}

type SimulationSettings struct {
	Walls               WallList                   `json:"walls"`
	WindowSizeX         float64                    `json:"window_size_x"`
	WindowSizeY         float64                    `json:"window_size_y"`
	Duration            float64                    `json:"duration"`
	Dt                  float64                    `json:"dt"`
	TimeToRecover       float64                    `json:"time_to_recover"`
	FracRandomUnmovable float64                    `json:"frac_unmovable"`
	NbRandomAgents      uint                       `json:"nb_random_agents"`
	NbRandomSicks       uint                       `json:"nb_random_sick"`
	PDeath              float64                    `json:"death_proportion"`
	AgentStartSpeed     float64                    `json:"agents_start_speed"`
	AgentRadius         float64                    `json:"agents_radius"`
	StartAgParam        [](*StartAgentsParameters) `json="start_agents"`
}

func DefaultSettings() *SimulationSettings {
	settings := new(SimulationSettings)
	json_string := `
	{
		"walls": [
		  {
			"start": {
			  "x": 10,
			  "y": 0
			},
			"end": {
			  "x": 10,
			  "y": 14
			},
			"radius": 0.3
		  },
		  {
			"start": {
			  "x": 10,
			  "y": 16
			},
			"end": {
			  "x": 10,
			  "y": 30
			},
			"radius": 0.3
		  },
		  {
			"start": {
			  "x": 15,
			  "y": 10
			},
			"end": {
			  "x": 25,
			  "y": 20
			},
			"radius": 0.3
		  },
		  {
			"start": {
			  "x": 15,
			  "y": 20
			},
			"end": {
			  "x": 25,
			  "y": 10
			},
			"radius": 0.3
		  }
		],
		"window_size_x": 30,
		"window_size_y": 30,
		"duration": 30,
		"dt": 0,
		"time_to_recover": 5,
		"frac_unmovable": 0.25,
		"nb_random_agents": 200,
		"nb_random_sick": 1,
		"death_proportion": 0.01,
		"agents_start_speed": 3,
		"agents_radius": 0.2,
		"StartAgParam": [
		  {
			"position": {
			  "x": 2,
			  "y": 2
			},
			"speed": {
			  "x": 0,
			  "y": 3
			},
			"state": 1,
			"movable": true
		  }
		]
	  }	  
	`
	json.Unmarshal([]byte(json_string), settings)
	return settings
}

func NewSimulation(settings *SimulationSettings) *Simulation {
	walls := settings.Walls
	agents := instanciate_agents(walls, settings)
	return &Simulation{
		Agents:      agents,
		AliveAgents: CopyList(agents),
		Walls:       walls,
		Infos:       make([]string, 0),
		Settings:    settings,
		Time:        0,
	}
}

func (sim *Simulation) Run() {
	for sim.Time < sim.Settings.Duration {
		sim.NextStep()
	}
}

func (sim *Simulation) NextStep() {
	//fmt.Println(simu_time)
	sim.Time += sim.Settings.Dt
	// Collision agent/walls
	bouceWithWalls(sim.AliveAgents, sim.Walls, sim.Time)

	// Collision between agents
	bounce(sim.AliveAgents, sim.Time)
	hasDied := make(AgentList, 0)
	// Move alive Agents
	for _, agent := range sim.AliveAgents {
		isDead := agent.updatePos(sim.Time, sim.Settings)
		if isDead {
			hasDied = append(hasDied, agent)
		}
	}
	// Update dead agent
	for _, deadAgent := range hasDied {
		sim.AliveAgents.RemoveAgent(deadAgent)
		fmt.Printf("%v has been removed from alives \n", deadAgent.ID)
	}

	//Get  Info
	for _, agent := range sim.Agents {
		sim.Infos = append(sim.Infos, fmt.Sprintf("%v;%v", sim.Time, agent.GetInfo()))
	}
}

func (sim *Simulation) SaveResults() {
	dirName := "Results_" + time.Now().Format("20060201_150405")
	os.Mkdir("."+string(filepath.Separator)+dirName, 0777)
	// Positions
	file, err := os.Create("." + string(filepath.Separator) + dirName + string(filepath.Separator) + "positions.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"time;id;state;x;y"})
	if err != nil {
		fmt.Println(err)
	}

	for _, value := range sim.Infos {
		err = writer.Write([]string{value})
		if err != nil {
			fmt.Println(err)
		}
	}

	// Settings
	filepath := "." + string(filepath.Separator) + dirName + string(filepath.Separator) + "settings.json"
	jsonFile, _ := json.MarshalIndent(sim.Settings, "", " ")
	_ = ioutil.WriteFile(filepath, jsonFile, 0644)
}
