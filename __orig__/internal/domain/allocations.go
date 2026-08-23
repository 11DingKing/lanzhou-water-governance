package domain

import "sort"

type Allocation struct {
	Region   string
	Capacity int64
	Used     int64
	Reserved int64
}

func (a Allocation) Available() int64             { return a.Capacity - a.Used - a.Reserved }
func (a Allocation) CanReserve(amount int64) bool { return amount > 0 && a.Available() >= amount }
func (a *Allocation) Reserve(amount int64) error {
	if !a.CanReserve(amount) {
		return ErrConflict
	}
	a.Reserved += amount
	return nil
}
func (a *Allocation) Commit(amount int64) error {
	if amount <= 0 || a.Reserved < amount {
		return ErrConflict
	}
	a.Reserved -= amount
	a.Used += amount
	return nil
}
func (a *Allocation) Release(amount int64) error {
	if amount <= 0 || a.Reserved < amount {
		return ErrConflict
	}
	a.Reserved -= amount
	return nil
}
func SortAllocations(items []Allocation) []Allocation {
	copyItems := append([]Allocation(nil), items...)
	sort.Slice(copyItems, func(i, j int) bool { return copyItems[i].Available() > copyItems[j].Available() })
	return copyItems
}
func TotalAvailable(items []Allocation) int64 {
	total := int64(0)
	for _, item := range items {
		total += item.Available()
	}
	return total
}
func Balance(items []Allocation) map[string]int64 {
	result := make(map[string]int64, len(items))
	for _, item := range items {
		result[item.Region] = item.Available()
	}
	return result
}
func Rebalance(items []Allocation, amount int64) []Allocation {
	sorted := SortAllocations(items)
	for i := range sorted {
		if amount <= 0 {
			break
		}
		take := sorted[i].Available()
		if take > amount {
			take = amount
		}
		sorted[i].Reserved += take
		amount -= take
	}
	return sorted
}
