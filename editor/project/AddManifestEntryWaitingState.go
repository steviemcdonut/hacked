package project

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/inkyblackness/hacked/ss1/resource"
	"github.com/inkyblackness/hacked/ss1/resource/lgres"
	"github.com/inkyblackness/hacked/ss1/world"
	"github.com/inkyblackness/imgui-go"
)

type addManifestEntryWaitingState struct {
	view *View
}

func (state addManifestEntryWaitingState) Render() {
	if imgui.BeginPopupModalV("Add static world data", nil,
		imgui.WindowFlagsNoResize|imgui.WindowFlagsNoMove|imgui.WindowFlagsNoSavedSettings) {

		imgui.TextUnformatted(`Waiting for folders/files.

From your file browser drag'n'drop the folder (or files)
of the static data you want to reference into the editor window.
Typically, you would use the main "data" directory of the game
(where all the .res files are).
`)
		imgui.Separator()
		if imgui.Button("Cancel") {
			state.view.fileState = &idlePopupState{}
			imgui.CloseCurrentPopup()
		}
		imgui.EndPopup()
	} else {
		state.view.fileState = &idlePopupState{}
	}
}

type fileStaging struct {
	failedFiles int
	resources   map[string]resource.Provider
}

func (staging *fileStaging) stage(name string, enterDir bool) {
	fileInfo, err := os.Stat(name)
	if err != nil {
		staging.failedFiles++
		return
	}
	file, err := os.Open(name)
	if err != nil {
		return
	}
	defer file.Close()

	if fileInfo.IsDir() {
		if enterDir {
			subNames, _ := file.Readdirnames(0)
			for _, subName := range subNames {
				staging.stage(filepath.Join(name, subName), false)
			}
		}
	} else {
		fileData, err := ioutil.ReadAll(file)
		if err != nil {
			staging.failedFiles++
		}

		reader, err := lgres.ReaderFrom(bytes.NewReader(fileData))
		if err == nil {
			staging.resources[name] = reader
		}

		if err != nil {
			staging.failedFiles++
		}
	}
}

func (state addManifestEntryWaitingState) HandleFiles(names []string) {
	staging := fileStaging{
		resources: make(map[string]resource.Provider),
	}

	for _, name := range names {
		staging.stage(name, true)
	}
	if len(staging.resources) > 0 {
		entry := &world.ManifestEntry{
			ID: names[0],
		}

		for filename, provider := range staging.resources {
			localized := resource.LocalizeResourcesByFilename(provider, filename)
			entry.Resources = append(entry.Resources, localized)
		}

		state.view.requestAddManifestEntry(entry)
		state.view.fileState = &idlePopupState{}
	} else {
		// TODO: add failed state, notifying failure
		fmt.Printf("Failed...%d files\n", staging.failedFiles)
		state.view.fileState = &idlePopupState{}
	}
}
