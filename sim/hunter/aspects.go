package hunter

import (
	"github.com/wowsims/forever/sim/core"
	"github.com/wowsims/forever/sim/core/dbcenums"
	"github.com/wowsims/forever/sim/core/spelldata"
)

var aspectOfTheHawkRank = spellData.AspectOfTheHawk.Highest()
var aspectOfTheBeastRank = spellData.AspectOfTheBeast.Highest()
var aspectOfTheViperRank = spellData.AspectOfTheViper.Highest()

// The hastes the aspects proc: Quick Shots off Hawk, Quick Strikes off Beast.
var quickShots = spellData.AspectOfTheHawkTriggered.Highest()
var quickStrikes = spellData.AspectOfTheBeastTriggered.Highest()

func (hunter *Hunter) registerAspects() {
	hunter.AspectOfTheHawkAura = hunter.registerAspect(aspectOfTheHawkRank)
	hunter.AspectOfTheBeastAura = hunter.registerAspect(aspectOfTheBeastRank)
	hunter.AspectOfTheViperAura = hunter.registerAspect(aspectOfTheViperRank)

	// A_OBS_MOD_POWER has no row in the parse table: a share of maximum mana every period, by hand.
	regen := aspectOfTheViperRank.Effect(dbcenums.A_OBS_MOD_POWER, 0)
	manaMetrics := hunter.NewManaMetrics(core.ActionID{SpellID: aspectOfTheViperRank.ID})
	var regenTick *core.PendingAction
	hunter.AspectOfTheViperAura.ApplyOnGain(func(_ *core.Aura, sim *core.Simulation) {
		regenTick = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
			Period: regen.Period(),
			OnAction: func(sim *core.Simulation) {
				hunter.AddMana(sim, hunter.MaxMana()*regen.Percent(), manaMetrics)
			},
		})
	}).ApplyOnExpire(func(_ *core.Aura, sim *core.Simulation) {
		regenTick.Cancel(sim)
	})

	hunter.registerDeadlyAspects()
}

// One aspect at a time: a permanent self-buff the row's effects parse onto, and the spell that swaps
// to it.
func (hunter *Hunter) registerAspect(row *spelldata.Spell) *core.Aura {
	aura := hunter.RegisterAura(spelldata.AuraConfig(row))
	spelldata.ParseEffects(&hunter.Character, aura, row)
	aura.NewExclusiveEffect("Aspect", true, core.ExclusiveEffect{})

	config := spelldata.SpellConfig(&hunter.Unit, row, spelldata.Flags(core.SpellFlagAPL))

	config.ApplyEffects = func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
		aura.Activate(sim)
	}

	config.RelatedSelfBuff = aura

	hunter.RegisterSpell(config)

	return aura
}

// Deadly Aspects is the chance the aspects' own procs fire with: the Hawk and Beast rows state the
// trigger and no rate, and the talent adds SPELLMOD_CHANCE_OF_SUCCESS to each - 2% a rank - which
// the parse has no row for. The tooltip names Auto Shot and melee auto attacks where the rows hear
// every ranged and melee hit, so each mask is narrowed to the white hits by hand.
func (hunter *Hunter) registerDeadlyAspects() {
	rank := hunter.Talents.DeadlyAspects
	if rank == 0 {
		return
	}

	// A_MOD_RANGED_HASTE has no row in the parse table, so Quick Shots' haste is attached by hand;
	// Quick Strikes' melee haste parses.
	quickShotsAura := hunter.RegisterAura(spelldata.AuraConfig(quickShots)).
		AttachMultiplyRangedHaste(1 + quickShots.Effect(dbcenums.A_MOD_RANGED_HASTE, 0).Percent())
	quickStrikesAura := hunter.RegisterAura(spelldata.AuraConfig(quickStrikes))
	spelldata.ParseEffects(&hunter.Character, quickStrikesAura, quickStrikes)

	hawk := spelldata.ProcTrigger(&hunter.Character, aspectOfTheHawkRank,
		func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			quickShotsAura.Activate(sim)
		},
		spelldata.Chance(spellData.DeadlyAspects.EffectAt(1).FractionAt(rank)))
	hawk.Name = "Deadly Aspects - Quick Shots"
	hawk.ProcMask = core.ProcMaskRangedAuto
	hunter.AspectOfTheHawkAura.AttachProcTrigger(hawk)

	beast := spelldata.ProcTrigger(&hunter.Character, aspectOfTheBeastRank,
		func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			quickStrikesAura.Activate(sim)
		},
		spelldata.Chance(spellData.DeadlyAspects.EffectAt(2).FractionAt(rank)))
	beast.Name = "Deadly Aspects - Quick Strikes"
	beast.ProcMask = core.ProcMaskMeleeWhiteHit
	hunter.AspectOfTheBeastAura.AttachProcTrigger(beast)
}
