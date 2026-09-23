package hunter

// The pet is not modelled, and every Beast Mastery talent but two is the pet's: Endurance
// Training, Focused Fire, Improved Aspect of the Monkey, Pathfinding, Improved Revive Pet, Bestial
// Swiftness, Unleashed Fury, Improved Mend Pet, Ferocity, Spirit Bond, Intimidation, Frenzy and
// Bestial Wrath, and Summon Hawk (1293527), which summons one.
func (hunter *Hunter) registerBeastMasteryTalents() {
	// Tier 1
	// Deadly Aspects: aspects.go

	// Tier 5
	hunter.registerBestialDiscipline()
}

// Effect 1 is the pet's focus regeneration. Effect 2 is the hunter's own mana regeneration while
// casting, A_MOD_MANA_REGEN_INTERRUPT, which the parse table has no row for.
func (hunter *Hunter) registerBestialDiscipline() {
	if hunter.Talents.BestialDiscipline == 0 {
		return
	}

	hunter.PseudoStats.SpiritRegenRateCasting += spellData.BestialDiscipline.EffectAt(2).FractionAt(hunter.Talents.BestialDiscipline)
}
