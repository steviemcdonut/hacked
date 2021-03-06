package levels

import "github.com/inkyblackness/hacked/ss1/world"

type controlViewModel struct {
	selectedLevel                   int
	selectedAtlasIndex              int
	selectedSurveillanceObjectIndex int
	selectedTextureAnimationIndex   int

	restoreFocus bool
	windowOpen   bool
}

func freshControlViewModel() controlViewModel {
	return controlViewModel{
		selectedLevel:                 world.StartingLevel,
		selectedTextureAnimationIndex: 1,
	}
}
