// Services > Meta Parser
// This service provides a way to read and write meta fields in a structured way.

package metaparse

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

/*----------  Meta payloads  ----------*/

type BoardMeta struct{}
type ThreadMeta struct{}
type PostMeta struct{}
type VoteMeta struct {
	/*----------  Follows guidelines  ----------*/
	FGReason string `json:"fg_reason,omitempty"`
	MAReason string `json:"ma_reason,omitempty"`
}
type KeyMeta struct{}
type TruststateMeta struct {
	CanonicalName string `json:"canonical_name,omitempty"`
}

func (e *BoardMeta) IsMeta()      {}
func (e *ThreadMeta) IsMeta()     {}
func (e *PostMeta) IsMeta()       {}
func (e *VoteMeta) IsMeta()       {}
func (e *KeyMeta) IsMeta()        {}
func (e *TruststateMeta) IsMeta() {}

type MetaStruct interface {
	IsMeta()
}

func ReadMeta(entityType, metaAsString string) (MetaStruct, error) {
	if len(metaAsString) == 0 {
		return nil, nil
	}
	switch entityType {
	case "Board":
		return nil, nil
	case "Thread":
		return nil, nil
	case "Post":
		return nil, nil
	case "Vote":
		em := VoteMeta{}
		err := json.Unmarshal([]byte(metaAsString), &em)
		if err != nil {
			return nil, err
		}
		return &em, nil
	case "Key":
		return nil, nil
	case "Truststate":
		em := TruststateMeta{}
		err := json.Unmarshal([]byte(metaAsString), &em)
		if err != nil {
			return nil, err
		}
		return &em, nil
	}
	return nil, fmt.Errorf("The entityType you gave to JSON Parser is unknown. You gave: %v", entityType)
}

func CreateMetaString(payloadStruct MetaStruct) (string, error) {
	jsonAsByte, err := json.Marshal(payloadStruct)
	return string(jsonAsByte), err
}
