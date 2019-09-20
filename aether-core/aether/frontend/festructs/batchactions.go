package festructs

import (
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"fmt"
)

// GetAllBoards gets all boards within our observable universe. Observable universe is all entities that it could possibly be changed by the entities we have in this current bucket, at this single instant in time.
func GetAllBoards(targetPtr *[]BoardCarrier, observableUniverse map[string]bool) error {
	var bc []BoardCarrier
	for fp, _ := range observableUniverse {
		b := BoardCarrier{}
		err := globals.KvInstance.One("Fingerprint", fp, &b)
		if err != nil {
			errTxt := fmt.Sprintf("Reading Board on GetAllBoards failed. Error: %#v", err)
			logging.Logf(1, errTxt)
			// return fmt.Errorf(errTxt)
			continue
		}
		bc = append(bc, b)
	}
	*targetPtr = bc
	return nil
}

func GetAllUserHeaderCarriers(targetPtr *[]UserHeaderCarrier, observableUniverse map[string]bool) error {
	var uhcs []UserHeaderCarrier
	for fp, _ := range observableUniverse {
		uhc := UserHeaderCarrier{}
		err := globals.KvInstance.One("Fingerprint", fp, &uhc)
		if err != nil {
			errTxt := fmt.Sprintf("Reading UserHeaderCarrier on GetAllUserHeaderCarriers failed. Error: %#v", err)
			logging.Logf(1, errTxt)
			return fmt.Errorf(errTxt)
		}
		uhcs = append(uhcs, uhc)
	}
	*targetPtr = uhcs
	return nil
}
