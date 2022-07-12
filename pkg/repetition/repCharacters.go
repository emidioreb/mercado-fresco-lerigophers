package repetition

func RepetitionCharacters(rep int) string {
	cont := ""
	for i := 0; i < rep; i++ {
		cont = cont + "a"
	}
	return cont
}
