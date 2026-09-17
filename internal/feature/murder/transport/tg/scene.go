package murder_tg_transport

import (
	"strconv"
	"strings"
)

// Сцена собрана из рядов однотипных эмодзи: в Telegram у всех эмодзи одна
// ширина, поэтому «одна клетка = один эмодзи» даёт ровные колонки без
// моноширинного шрифта.
const (
	tokEmpty  = "⬜"  // пустая клетка
	tokPaw    = "🐾"  // лапа медведя
	tokPoke   = "👇"  // лапа продавливает улей
	tokTrail  = "💨"  // след прыжков лапы
	tokHive   = "🏠"  // целый улей
	tokBoom   = "💥"  // улей в момент тычка
	tokHoney  = "🍯"  // вскрытый улей с мёдом
	tokMiss   = "🏚️" // вскрытый пустой улей
	tokBee    = "🐝"
	tokSpark  = "✨"
	tokGround = "🟫"
)

// beesMode — что происходит в небе над пасекой.
type beesMode int

const (
	beesFly       beesMode = iota // пчёлы кружат
	beesSwarm                     // рой злых пчёл
	beesCelebrate                 // искры-салют
)

// robScene — один кадр анимации.
type robScene struct {
	caption   string   // реплика медведя
	paw       int      // над каким ульем лапа; 0 — лапа спрятана
	trail     bool     // рисовать след прыжков
	poke      bool     // лапа продавливает улей: 👇 сверху, 💥 на улье
	bees      beesMode // небо
	pawHive   int      // >0 — вскрытый пустой улей (промах)
	honeyHive int      // >0 — вскрытый улей с мёдом
}

// buildFrames собирает кадры анимации: лапа прыгает к выбранному улью,
// продавливает его и в финале вскрывает улей с мёдом — или показывает,
// что под лапой пусто.
func buildFrames(pawHive int, honeyHive int, win bool) []robScene {
	// Лапа скачет крупными шагами, чтобы анимация не длилась вечность.
	hop := (pawHive-1)/3 + 1

	frames := []robScene{
		{caption: "«где мёд?!»", bees: beesFly},
	}
	for p := 1; p < pawHive; p += hop {
		frames = append(frames, robScene{caption: "«щас вскрою…»", paw: p, trail: true, bees: beesFly})
	}
	frames = append(frames,
		robScene{caption: "«ну-ну…»", paw: pawHive, trail: true, bees: beesFly},
		robScene{caption: "«ХЛОП!»", paw: pawHive, poke: true, bees: beesFly},
		robScene{caption: "«тя-я-яну!»", paw: pawHive, bees: beesFly},
	)

	// Финал: на победе лапа держит горшок над вскрытым ульем,
	// на промахе мишка спрятал лапу, пчёлы злятся.
	final := robScene{bees: beesCelebrate, caption: "«МЁ-О-ОД!!!»"}
	if win {
		final.paw = pawHive
		final.honeyHive = honeyHive
	} else {
		final.caption = "«пчёлы злятся…»"
		final.bees = beesSwarm
		final.pawHive = pawHive
		final.honeyHive = honeyHive
	}
	return append(frames, final)
}

// sceneText собирает кадр: медведь с репликой и ряды пасеки.
func sceneText(scene robScene, hives int) string {
	var b strings.Builder
	b.WriteString("🐻 " + scene.caption + "\n")
	b.WriteString(beeRow(scene.bees, hives))
	b.WriteString(trailRow(scene, hives))
	b.WriteString(pokeRow(scene, hives))
	b.WriteString(hiveRow(scene, hives))
	b.WriteString(numberRow(hives))
	b.WriteString(groundRow(hives))
	return b.String()
}

// beeRow — небо: пчёлы кружат, рой злых пчёл или салют.
func beeRow(mode beesMode, hives int) string {
	var b strings.Builder
	for i := 1; i <= hives; i++ {
		switch mode {
		case beesSwarm:
			b.WriteString(tokBee)
		case beesCelebrate:
			b.WriteString(tokSpark)
		default:
			if i%2 == 1 {
				b.WriteString(tokBee)
			} else {
				b.WriteString(tokSpark)
			}
		}
	}
	b.WriteString("\n")
	return b.String()
}

// trailRow — строка лапы: 🐾 над текущим ульем и 💨 позади неё.
func trailRow(scene robScene, hives int) string {
	var b strings.Builder
	for i := 1; i <= hives; i++ {
		switch {
		case scene.paw == i:
			b.WriteString(tokPaw)
		case scene.trail && i < scene.paw:
			b.WriteString(tokTrail)
		default:
			b.WriteString(tokEmpty)
		}
	}
	b.WriteString("\n")
	return b.String()
}

// pokeRow — строка тычка: 👇 над ульем в момент удара.
func pokeRow(scene robScene, hives int) string {
	var b strings.Builder
	for i := 1; i <= hives; i++ {
		if scene.poke && scene.paw == i {
			b.WriteString(tokPoke)
		} else {
			b.WriteString(tokEmpty)
		}
	}
	b.WriteString("\n")
	return b.String()
}

// hiveRow — ряд ульев с вскрытыми клетками.
func hiveRow(scene robScene, hives int) string {
	var b strings.Builder
	for i := 1; i <= hives; i++ {
		switch {
		case scene.honeyHive == i:
			b.WriteString(tokHoney)
		case scene.pawHive == i:
			b.WriteString(tokMiss)
		case scene.poke && scene.paw == i:
			b.WriteString(tokBoom)
		default:
			b.WriteString(tokHive)
		}
	}
	b.WriteString("\n")
	return b.String()
}

// numberRow — номера ульев эмодзи-цифрами.
func numberRow(hives int) string {
	var b strings.Builder
	for i := 1; i <= hives; i++ {
		b.WriteString(keycap(i))
	}
	b.WriteString("\n")
	return b.String()
}

// keycap — эмодзи-цифра 1️⃣…9️⃣, для десятки — 🔟.
func keycap(i int) string {
	if i == 10 {
		return "🔟"
	}
	return strconv.Itoa(i) + "\uFE0F\u20E3"
}

// groundRow — земля под пасекой.
func groundRow(hives int) string {
	return strings.Repeat(tokGround, hives) + "\n"
}
