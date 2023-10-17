package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/duanemay/advantage360/model"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
)

const (
	keymapJsonFileName           = "keymap.json"
	adv360KeymapTemplateFileName = "template.keymap"
	adv360KeymapOutputFileName   = "adv360.keymap"
)

func main() {
	app := &cli.App{
		Name:      "generate",
		Usage:     "Generate keyboard definitions",
		ArgsUsage: "file",
		Action: func(cCtx *cli.Context) error {
			if cCtx.Args().Len() != 1 {
				_ = cli.ShowAppHelp(cCtx)
				return fmt.Errorf("must specify a file")
			}
			orchestrator(cCtx.Args().Get(0))
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func orchestrator(file string) {
	content, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("Error when reading %s: %s", file, err)
	}

	keysFile := model.KeysFile{}
	err = json.Unmarshal(content, &keysFile)
	if err != nil {
		log.Fatal("Error in json file: ", err)
	}

	// remove comments from keyIds
	for i, keyId := range keysFile.KeyIds {
		if keyId.Comment != "" {
			keysFile.KeyIds = append(keysFile.KeyIds[:i], keysFile.KeyIds[i+1:]...)
		}
	}
	sort.Sort(keysFile.KeyIds)
	var keyIdList []string
	for _, keyId := range keysFile.KeyIds {
		keyIdList = append(keyIdList, keyId.KeyId)
	}

	keyMapFile := model.NewKeyMapFile()
	for _, layerName := range model.LayerOrder {
		layer, _ := keysFile.Layers.GetLayer(layerName)
		keyMap := map[string]model.Key{}
		for _, key := range layer.Keys {
			keyMap[key.Id] = key
		}

		var keyActions []string
		for _, keyId := range keyIdList {
			key, ok := keyMap[keyId]
			if ok {
				value := "&" + key.Action
				if key.Value != "" {
					value += " " + key.Value
				}
				keyActions = append(keyActions, value)
			} else {
				keyActions = append(keyActions, "&none")
			}
		}
		keyMapFile.Layers = append(keyMapFile.Layers, keyActions)
	}
	jsonOut, _ := json.Marshal(keyMapFile)
	err = os.WriteFile(keymapJsonFileName, jsonOut, 0644)
	if err != nil {
		log.Fatalf("Could not write %s: %s", keymapJsonFileName, err)
	}
	fmt.Printf("Wrote %s...\n", keymapJsonFileName)

	content, err = os.ReadFile(adv360KeymapTemplateFileName)
	if err != nil {
		log.Fatalf("Error when reading %s: %s", adv360KeymapTemplateFileName, err)
	}

	for i, layerName := range keyMapFile.LayerNames {
		layer := keyMapFile.Layers[i]
		bindingName := []byte("BINDINGS_" + string(layerName))
		binding := generateBindingString(layer, keysFile.KeyIds)
		content = bytes.ReplaceAll(content, bindingName, binding)
	}

	err = os.WriteFile(adv360KeymapOutputFileName, content, 0644)
	if err != nil {
		log.Fatalf("Could not write %s: %s", adv360KeymapOutputFileName, err)
	}
	fmt.Printf("Wrote %s...\n", adv360KeymapOutputFileName)

}

func generateBindingString(layer []string, ids model.KeyIdArray) []byte {
	idRow := 1
	var result []byte
	for idx, id := range ids {
		if id.Row > idRow {
			result = append(result, []byte("\n")...)
			idRow = id.Row
		}
		result = append(result, []byte("    ")...)
		result = append(result, []byte(layer[idx])...)
	}

	return result
}
