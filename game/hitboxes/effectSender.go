package hitboxes

//Передача эффектов

func ABSender(a HitBoxer, b HitBoxer) {
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

	if !sender.CanSendEffects() {
		//log.Println("sender.CanSendEffects() == false")
		return
	}

	effects := sender.GetEffectsForTransfer()
	if len(effects) == 0 {
		//log.Println("effects пустые")
		return
	}

	receiver.OnCollision(effects)
	sender.OnEffectsSent()
}
