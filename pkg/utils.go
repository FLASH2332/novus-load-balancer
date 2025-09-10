package pkg

type Node struct {
	Key   string
	Value []byte
	Prev  *Node
	Next  *Node
}

type DoublyLinkedList struct {
	head *Node
	tail *Node
}

func NewDoublyLinkedList() *DoublyLinkedList {
	return &DoublyLinkedList{}
}

// Add node to front (most recently used)
func (dll *DoublyLinkedList) PushFront(node *Node) {
	node.Prev = nil
	node.Next = dll.head

	if dll.head != nil {
		dll.head.Prev = node
	}
	dll.head = node

	if dll.tail == nil {
		dll.tail = node
	}
}

// Remove node from anywhere
func (dll *DoublyLinkedList) Remove(node *Node) {
	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		dll.head = node.Next
	}

	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		dll.tail = node.Prev
	}
}

// Move node to front
func (dll *DoublyLinkedList) MoveToFront(node *Node) {
	dll.Remove(node)
	dll.PushFront(node)
}

// Remove node from back (least recently used)
func (dll *DoublyLinkedList) PopBack() *Node {
	if dll.tail == nil {
		return nil
	}
	node := dll.tail
	dll.Remove(node)
	return node
}
