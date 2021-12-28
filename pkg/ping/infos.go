package ping

// infosVersion is version type of Infos.
type infosVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

// infosPlayersSample is sample type of Infos.Players.
type infosPlayersSample struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// infosPlayers is version type of Infos.
type infosPlayers struct {
	Max    int                  `json:"max"`
	Online int                  `json:"online"`
	Sample []infosPlayersSample `json:"sample"`
}

// Infos represents usual informations contained in ping response.
type Infos struct {
	Version     infosVersion `json:"version"`
	Players     infosPlayers `json:"players"`
	Description string       `json:"description"`
	Favicon     string       `json:"favicon"`
}

// Infos extracts informations from ping response properties (JSON), and put it into an Infos structure.
func (m *JSON) Infos() Infos {
	var infos Infos

	version, ok := (*m)["version"]
	if ok {
		versionJSON, ok := version.(map[string]interface{})
		if ok {
			name, ok := versionJSON["name"]
			if ok {
				nameString, ok := name.(string)
				if ok {
					infos.Version.Name = nameString
				}
			}

			protocol, ok := versionJSON["protocol"]
			if ok {
				protocolNumber, ok := protocol.(float64)
				if ok {
					infos.Version.Protocol = int(protocolNumber)
				}
			}
		}
	}

	players, ok := (*m)["players"]
	if ok {
		playersJSON, ok := players.(map[string]interface{})
		if ok {
			max, ok := playersJSON["max"]
			if ok {
				maxNumber, ok := max.(float64)
				if ok {
					infos.Players.Max = int(maxNumber)
				}
			}

			online, ok := playersJSON["online"]
			if ok {
				onlineNumber, ok := online.(float64)
				if ok {
					infos.Players.Online = int(onlineNumber)
				}
			}

			sample, ok := playersJSON["sample"]
			if ok {
				sampleArray, ok := sample.([]interface{})
				if ok {
					for _, player := range sampleArray {
						playerSampleJSON, ok := player.(map[string]interface{})
						if ok {
							player := infosPlayersSample{}

							name, ok := playerSampleJSON["name"]
							if ok {
								nameString, ok := name.(string)
								if ok {
									player.Name = nameString
								}
							}

							id, ok := playerSampleJSON["id"]
							if ok {
								idString, ok := id.(string)
								if ok {
									player.ID = idString
								}
							}

							infos.Players.Sample = append(infos.Players.Sample, player)
						}
					}
				}
			}
		}
	}

	description, ok := (*m)["description"]
	if ok {
		descriptionString, ok := description.(string)
		if ok {
			infos.Description = descriptionString
		} else {
			descriptionJSON, ok := description.(map[string]interface{})
			if ok {
				text, ok := descriptionJSON["text"]
				if ok {
					textString, ok := text.(string)
					if ok {
						infos.Description = textString
					}
				}
			}
		}
	}

	favicon, ok := (*m)["favicon"]
	if ok {
		faviconString, ok := favicon.(string)
		if ok {
			infos.Favicon = faviconString
		}
	}

	return infos
}
