package services

func CalculateCarType(answers map[string]string) string {
	score := map[string]int{
		"SUV":       0,
		"Sedan":     0,
		"Hatchback": 0,
		"Electric":  0,
	}

	switch answers["1"] {
	case "Family":
		score["SUV"] += 3
		score["Sedan"] += 2
	case "Friends":
		score["Hatchback"] += 3
	case "Alone":
		score["Sedan"] += 2
	}

	switch answers["2"] {
	case "City":
		score["Hatchback"] += 3
	case "Highway":
		score["Sedan"] += 3
	case "Off-road":
		score["SUV"] += 4
	}

	switch answers["3"] {
	case "Electric":
		score["Electric"] += 5
	case "Hybrid":
		score["Sedan"] += 2
	}

	switch answers["7"] {
	case "Large":
		score["SUV"] += 3
	case "Small":
		score["Hatchback"] += 3
	}

	switch answers["10"] {
	case "Safety":
		score["SUV"] += 2
	case "Performance":
		score["Sedan"] += 3
	case "Technology":
		score["Electric"] += 3
	}

	bestType := "Sedan"
	max := 0
	for k, v := range score {
		if v > max {
			max = v
			bestType = k
		}
	}

	return bestType
}
