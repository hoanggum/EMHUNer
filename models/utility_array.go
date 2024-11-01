package models

import "fmt"

type UtilityArray struct {
	RTWUs []float64
	RLUs  []float64
	RSUs  []float64
	Size  int
}

func NewUtilityArray(size int) *UtilityArray {
	size += 1
	return &UtilityArray{
		RTWUs: make([]float64, size),
		RLUs:  make([]float64, size),
		RSUs:  make([]float64, size),
		Size:  size,
	}
}

func (ua *UtilityArray) SetRTWU(index int, value float64) {
	if index >= 0 && index < ua.Size {
		ua.RTWUs[index] = value
	}
}

func (ua *UtilityArray) GetRTWU(index int) float64 {
	if index >= 0 && index < ua.Size {
		return ua.RTWUs[index]
	}
	return 0
}

func (ua *UtilityArray) SetRLU(index int, value float64) {
	if index >= 0 && index < ua.Size {
		ua.RLUs[index] = value
	}
}

func (ua *UtilityArray) GetRLU(index int) float64 {
	if index >= 0 && index < ua.Size {
		return ua.RLUs[index]
	}
	return 0
}

func (ua *UtilityArray) SetRSU(index int, value float64) {
	if index >= 0 && index < ua.Size {
		ua.RSUs[index] = value
	}
}

func (ua *UtilityArray) GetRSU(index int) float64 {
	if index >= 0 && index < ua.Size {
		return ua.RSUs[index]
	}
	return 0
}

func (ua *UtilityArray) PrintUtilityArray() {
	fmt.Println("RTWU Array:")
	fmt.Println(ua.RTWUs)

	fmt.Println("RLU Array:")
	fmt.Println(ua.RLUs)

	fmt.Println("RSU Array:")
	fmt.Println(ua.RSUs)
}
