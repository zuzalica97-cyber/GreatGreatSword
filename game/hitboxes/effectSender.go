package hitboxes

//Передача эффектов

func ABSender(a, b any) {
	// Проверяем, может ли a отправить эффект
	sender, ok := a.(LetterSender)
	if !ok {
		//log.Println("a не является LetterSender")
		return
	}

	receiver, ok := b.(LetterReceiver)
	if !ok {
		//log.Println("b не является LetterReceiver")
		return
	}

	if !sender.CanSendEffects(b) {
		//log.Println("sender.CanSendEffects() == false")
		return
	}

	effects := sender.GetEffectsForTransfer(b)
	if len(effects) == 0 {
		return
	}

	receiver.OnCollision(effects)
}
