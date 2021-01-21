package main

import (
	"encoding/binary"
	"math"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type Telemetry struct {
	LapTime     uint32 `json:"lap_time"`
	LastLapTime uint32 `json:"last_lap_time"`
	Laps        uint32 `json:"laps"`

	Speed float32 `json:"speed"` // Km/h

	GFrontal    float32 `json:"g_frontal"`
	GHorizontal float32 `json:"g_horizontal"`
	GLateral    float32 `json:"g_lateral"`
}

func BytesToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func TransformEvent(event cloudevents.Event) (*cloudevents.Event, error) {
	rawData := event.Data()

	tel := Telemetry{
		LapTime:     binary.LittleEndian.Uint32(rawData[40:44]),
		LastLapTime: binary.LittleEndian.Uint32(rawData[44:48]),
		Laps:        binary.LittleEndian.Uint32(rawData[52:56]),
		Speed:       BytesToFloat32(rawData[8:12]),
		GFrontal:    BytesToFloat32(rawData[36:40]),
		GHorizontal: BytesToFloat32(rawData[32:36]),
		GLateral:    BytesToFloat32(rawData[28:32]),
	}

	outputEvent := cloudevents.NewEvent()
	outputEvent.SetID(event.ID())
	outputEvent.SetType("telemetry")
	outputEvent.SetSource("rawdataprocessor.myapp.com")
	return &outputEvent, outputEvent.SetData("application/json", tel)
}
