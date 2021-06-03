package domain

type StreamIterator struct {
	index    int
	messages []Message
}

func (si StreamIterator) hasNext() bool {
	if si.index < len(si.messages) {
		return true
	}
	return false
}

func (si StreamIterator) getNext() Message {
	if !si.hasNext() {
		return nil
	}

	message := si.messages[si.index]
	si.index++
	return message
}

func NewStreamIterator(messages []Message) StreamIterator {
	return StreamIterator{0, messages}
}
