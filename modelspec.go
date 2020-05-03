package runway

type ModelInfo struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Inputs      []DataType `json:"inputs"`
	Outputs     []DataType `json:"outputs"`
}

type DataType struct {
	// Common
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Default     interface{} `json:"default,omitempty"`

	// Array and Text
	ItemType  string `json:"itemType,omitempty"`
	MinLength int    `json:"minLength,omitempty"`
	MaxLength int    `json:"maxLength,omitempty"`

	// Image
	Channels            int    `json:"channels,omitempty"`
	MinWidth            int    `json:"minWidth,omitempty"`
	MaxWidth            int    `json:"maxWidth,omitempty"`
	MinHeight           int    `json:"minHeight,omitempty"`
	MaxHeight           int    `json:"maxHeight,omitempty"`
	Width               int    `json:"width,omitempty"`
	Height              int    `json:"height,omitempty"`
	DefaultOutputFormat string `json:"defaultOutputFormat,omitempty"`

	// Vector
	Length       int     `json:"length,omitempty"`
	SamplingMean float64 `json:"samplingMean,omitempty"`
	SamplingStd  float64 `json:"samplingStd,omitempty"`

	// Category
	OneOf []string `json:"oneOf,omitempty"`

	// Number
	Max  float64 `json:"max,omitempty"`
	Min  float64 `json:"min,omitempty"`
	Step float64 `json:"step,omitempty"`

	// File
	IsDirectory bool   `json:"isDirectory,omitempty"`
	Extension   string `json:"extension,omitempty"`

	// Segmentation
	Labels       []string          `json:"labels,omitempty"`
	DefaultLabel string            `json:"defaultLabel,omitempty"`
	LabelToId    map[string]int    `json:"labelToId,omitempty"`
	LabelToColor map[string]string `json:"labelToColor,omitempty"`

	// Image Landmarks
	Connections [][]string `json:"connections,omitempty"`
}
