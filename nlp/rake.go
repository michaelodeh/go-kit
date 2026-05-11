package nlp

import (
	rake "github.com/afjoseph/RAKE.go"
)

func Rake(text string) string {

	// text := `The growing doubt of human autonomy and reason has created a state of moral confusion where man is left without the guidance of either revelation or reason. The result is the acceptance of a relativistic position which proposes that value judgements and ethical norms are exclusively matters of arbitrary preference and that no objectively valid statement can be made in this realm... But since man cannot live without values and norms, this relativism makes him an easy prey for irrational value systems.`

	candidates := rake.RunRake(text)

	return candidates[0].Key

	// fmt.Printf("\nsize: %d\n", len(candidates))
}
