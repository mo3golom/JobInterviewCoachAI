package model

var (
	InterviewConfig = InterviewAvailableValues{
		Nodes: map[string]Node{
			"developer": {
				Children: map[string]SubNode{
					"golang": {
						Position: "golang developer",
					},
					"php": {
						Position: "php developer",
					},
					"python": {
						Position: "python developer",
					},
					"rust": {
						Position: "rust developer",
					},
					"javascript": {
						Position: "javascript developer",
					},
					"swift": {
						Position: "swift developer",
					},
					"java": {
						Position: "java developer",
					},
					"c_plus_plus": {
						Position: "c plus plus developer",
					},
					"c_sharp": {
						Position: "c sharp developer",
					},
				},
			},
			"qa": {
				Position: "qa engineer",
			},
			"project_manager": {
				Position: "project manager",
			},
			"product_manager": {
				Position: "product manager",
			},
			"product_designer": {
				Position: "product designer",
			},
			string(BehavioralPosition): {
				Position: "behavioral interview",
			},
		},
	}
)

type (
	Node struct {
		Position Position
		Children map[string]SubNode
	}

	SubNode struct {
		Position Position
	}

	InterviewAvailableValues struct {
		Nodes map[string]Node
	}
)
