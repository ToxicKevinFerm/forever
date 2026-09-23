package hunter

import (
	"github.com/wowsims/forever/sim/common"
	_ "github.com/wowsims/forever/sim/common" // imported to get item effects included.
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/proto"

	"testing"
)

func init() {
	RegisterHunter()
	common.RegisterAllEffects()
}

func TestHunter(t *testing.T) {
	weaveRotation := core.GetAplRotation("../../ui/specs/hunter/dps/apls", "default")
	weaveRotation.Label = "weave"

	turretRotation := core.GetAplRotation("../../ui/specs/hunter/dps/apls", "default").Rotation
	turretRotation.ValueVariables[2] = &proto.APLValueVariable{
		Name: "Melee weave",
		Value: &proto.APLValue{
			Value: &proto.APLValue_Const{
				Const: &proto.APLValueConst{
					Val: "false",
				},
			},
		},
	}

	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassHunter,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceNightElf},
			GearSet:    core.GetGearSet("../../ui/specs/hunter/dps/gear_sets", "p1"),
			Talents:    DefaultMMTalents,
			OtherTalentSets: []core.TalentsCombo{
				{Label: "SV", Talents: DefaultSVTalents},
			},
			Consumables:      DefaultConsumables,
			SpecOptions:      core.SpecOptionsCombo{Label: "Default", SpecOptions: DefaultOptions},
			StartingDistance: 8,
			Profession1:      proto.Profession_Engineering,
			Profession2:      proto.Profession_Blacksmithing,

			Rotation: weaveRotation,
			OtherRotations: []core.RotationCombo{
				{Label: "Turret", Rotation: turretRotation},
			},

			ItemFilter: core.ItemFilter{
				ArmorType: proto.ArmorType_ArmorTypeMail,
				RangedWeaponTypes: []proto.RangedWeaponType{
					proto.RangedWeaponType_RangedWeaponTypeBow,
					proto.RangedWeaponType_RangedWeaponTypeCrossbow,
					proto.RangedWeaponType_RangedWeaponTypeGun,
				},
				WeaponTypes: []proto.WeaponType{
					proto.WeaponType_WeaponTypeAxe,
					proto.WeaponType_WeaponTypeDagger,
					proto.WeaponType_WeaponTypeFist,
					proto.WeaponType_WeaponTypePolearm,
					proto.WeaponType_WeaponTypeStaff,
					proto.WeaponType_WeaponTypeSword,
				},
				HandTypes: []proto.HandType{
					proto.HandType_HandTypeMainHand,
					proto.HandType_HandTypeOffHand,
					proto.HandType_HandTypeOneHand,
					proto.HandType_HandTypeTwoHand,
				},
			},
		},
	}))
}

var DefaultOptions = &proto.Player_Hunter{
	Hunter: &proto.Hunter{
		Options: &proto.Hunter_Options{
			ClassOptions: &proto.HunterOptions{
				PetType: proto.HunterOptions_PetNone,
			},
		},
	},
}

var DefaultMMTalents = "-30535525115023051-50000003"
var DefaultSVTalents = "-30535505100-500200030050020151"

var DefaultConsumables = &proto.ConsumesSpec{
	PotId:       22838,
	FlaskId:     22854,
	FoodId:      27658,
	ConjuredId:  22788,
	ExplosiveId: 30217,
	SuperSapper: true,
	ScrollAgi:   true,
	ScrollStr:   true,
}
