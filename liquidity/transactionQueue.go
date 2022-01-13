package liquidity

import "container/list"

type TransactionQueue struct {
	v *list.List
}

func CreateTransactionQueue() *TransactionQueue {
	return &TransactionQueue{list.New()}
}

func (tq *TransactionQueue) Push(val interface{}) {
	tq.v.PushBack(val)
}

func (tq *TransactionQueue) Pop() interface{} {
	front := tq.v.Front()
	if front != nil {
		return tq.v.Remove(front)
	}
	return nil
}

func (tq *TransactionQueue) Len() interface{} {
	return tq.v.Len()
}
