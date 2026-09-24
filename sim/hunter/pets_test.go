package hunter

import (
	"testing"

	"github.com/wowsims/forever/sim/core/proto"
)

// The generated family table and the options enum name the same families: a family the client
// adds needs its enum value, and an enum value the client no longer describes has to go.
func TestPetFamiliesMatchTheOptions(t *testing.T) {
	for name := range petFamilies {
		if _, ok := proto.HunterOptions_PetType_value[name]; !ok {
			t.Errorf("pet family %q has no HunterOptions.PetType value", name)
		}
	}
	for name, value := range proto.HunterOptions_PetType_value {
		if value == 0 {
			continue
		}
		if _, ok := petFamilies[name]; !ok {
			t.Errorf("HunterOptions.PetType %s names no generated pet family", name)
		}
	}
	for name, value := range proto.HunterOptions_PetAttackSpeed_value {
		if value == 0 {
			continue
		}
		if _, ok := petAttackSpeeds[proto.HunterOptions_PetAttackSpeed(value)]; !ok {
			t.Errorf("HunterOptions.PetAttackSpeed %s names no passive", name)
		}
	}
}
