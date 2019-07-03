package dealer

import "time"

func TransferMonthDay(dateStr *string) *string {
	date, _ := time.Parse("2006-01-02", *dateStr)
	formatDate := date.Format("1月2日")
	return &formatDate
}
