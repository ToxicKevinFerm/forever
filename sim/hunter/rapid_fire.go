package hunter

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
)

var rapidFireRank = spellData.RapidFire.Highest()

func (hunter *Hunter) registerRapidFire() {
	aura := hunter.RegisterAura(spelldata.AuraConfig(rapidFireRank))

	// The melee haste parses. A_MOD_RANGED_HASTE has no row in the parse table, so the ranged haste
	// the row states beside it is attached by hand.
	spelldata.ParseEffects(&hunter.Character, aura, rapidFireRank)
	aura.AttachMultiplyRangedHaste(1 + rapidFireRank.Effect(dbcenums.A_MOD_RANGED_HASTE, 0).Percent())

	config := spelldata.SpellConfig(&hunter.Unit, rapidFireRank, spelldata.Flags(core.SpellFlagAPL))

	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
		aura.Activate(sim)
	}

	config.RelatedSelfBuff = aura

	hunter.RapidFire = hunter.RegisterSpell(config)

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: hunter.RapidFire,
		Type:  core.CooldownTypeDPS,
	})
}
