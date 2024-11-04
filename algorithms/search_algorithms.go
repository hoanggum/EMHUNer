package algorithms

import (
	"EMHUNer/models"
	"EMHUNer/utility"
	"fmt"
)

type SearchAlgorithms struct {
	UtilityArray        *models.UtilityArray
	Beta                map[int]bool
	ItemList            []int
	FilteredPrimary     []int
	FilteredSecondary   []int
	HighUtilityItemsets []*models.HighUtilityItemset
}

func NewSearchAlgorithms(utilityArray *models.UtilityArray) *SearchAlgorithms {
	return &SearchAlgorithms{
		UtilityArray:        utilityArray,
		Beta:                make(map[int]bool),
		HighUtilityItemsets: []*models.HighUtilityItemset{},
	}
}
func (s *SearchAlgorithms) Search(eta []int, X map[int]bool, itemTransactionMap map[int][]*models.Transaction, primary []int, secondary []int, minU float64) {
	if len(primary) == 0 {
		return
	}

	for _, item := range primary {
		// Sao chép X vào Beta và thêm item
		s.Beta = copyMap(X)
		s.Beta[item] = true
		s.ItemList = mapKeys(s.Beta)

		// Tạo projectedItemTransactionMap và tính utility của Beta trong cùng một bước
		projectedItemTransactionMap, utilityBeta := s.createProjectedItemTransactionMapAndCalculateUtility(itemTransactionMap, s.ItemList)
		s.printProjectedDatabase(projectedItemTransactionMap)
		if utilityBeta >= minU {
			fmt.Printf("U(%d) = %.2f >= %.2f HUI Found: %v\n", item, utilityBeta, minU, s.Beta)
			s.HighUtilityItemsets = append(s.HighUtilityItemsets, models.NewHighUtilityItemset(s.ItemList, utilityBeta))
		} else {
			fmt.Printf("%.2f < %.2f so %v is not a HUI.\n", utilityBeta, minU, s.Beta)
		}

		// Đệ quy với projectedItemTransactionMap thay vì itemTransactionMap đầy đủ
		if utilityBeta > minU {
			s.SearchN(eta, s.Beta, projectedItemTransactionMap, minU)
		}

		// Tạo FilteredPrimary và FilteredSecondary dựa trên RSU và RLU
		s.FilteredPrimary = []int{}
		s.FilteredSecondary = []int{}
		utility.CalculateRSUForAllItem(projectedItemTransactionMap, s.ItemList, secondary, s.UtilityArray)
		utility.CalculateRLUForAllItem(projectedItemTransactionMap, s.ItemList, secondary, s.UtilityArray)
		// In ra RSU và RLU để kiểm tra
		for i, secItem := range secondary {
			if secItem == item {
				continue
			}
			if i > indexOf(secondary, item) {
				rsu := s.UtilityArray.GetRSU(secItem)
				rlu := s.UtilityArray.GetRLU(secItem)
				fmt.Printf("RSU(%d): %.2f, RLU(%d): %.2f\n", secItem, rsu, secItem, rlu) // In giá trị RSU và RLU

				if rsu >= minU {
					s.FilteredPrimary = append(s.FilteredPrimary, secItem)
				}
				if rlu >= minU {
					s.FilteredSecondary = append(s.FilteredSecondary, secItem)
				}
			}
		}

		// fmt.Printf("Beta= %v\n", s.Beta)
		// fmt.Printf("Primary%v = %v\n", s.ItemList, s.FilteredPrimary)
		// fmt.Printf("Secondary%v = %v\n", s.ItemList, s.FilteredSecondary)

		// Đệ quy gọi lại Search với `projectedItemTransactionMap` đã thu hẹp
		s.Search(eta, s.Beta, projectedItemTransactionMap, s.FilteredPrimary, s.FilteredSecondary, minU)
	}
}

func (s *SearchAlgorithms) SearchN(eta []int, beta map[int]bool, itemTransactionMap map[int][]*models.Transaction, minU float64) {
	if len(eta) == 0 {
		return
	}

	for _, item := range eta {
		// Sao chép beta vào betaNew và thêm item
		betaNew := copyMap(beta)
		betaNew[item] = true
		itemList := mapKeys(betaNew)

		// Tạo projectedItemTransactionMap và tính utility của betaNew trong cùng một bước
		projectedDBNew, utilityBetaNew := s.createProjectedItemTransactionMapAndCalculateUtility(itemTransactionMap, itemList)
		s.printProjectedDatabase(projectedDBNew)

		if utilityBetaNew >= minU {
			fmt.Printf("U(%d) = %.2f >= %.2f HUI Found: %v\n", item, utilityBetaNew, minU, betaNew)
			s.HighUtilityItemsets = append(s.HighUtilityItemsets, models.NewHighUtilityItemset(itemList, utilityBetaNew))
		} else {
			fmt.Printf("%.2f < %.2f so %v is not a HUI.\n", utilityBetaNew, minU, betaNew)
		}

		// Tạo FilteredPrimary dựa trên RSU
		itemIndex := indexOf(eta, item)
		filteredPrimary := []int{}
		utility.CalculateRSUForAllItem(projectedDBNew, itemList, eta, s.UtilityArray)

		for _, secItem := range eta {
			if secItem == item {
				continue
			}
			if indexOf(eta, secItem) > itemIndex {
				rsu := s.UtilityArray.GetRSU(secItem)
				if rsu >= minU {
					filteredPrimary = append(filteredPrimary, secItem)
				}
			}
		}
		fmt.Printf("Primary = %v\n", filteredPrimary)

		// Đệ quy gọi lại SearchN với projectedItemTransactionMap đã thu hẹp
		s.SearchN(filteredPrimary, betaNew, projectedDBNew, minU)
	}
}

func (s *SearchAlgorithms) createProjectedItemTransactionMapAndCalculateUtility(itemTransactionMap map[int][]*models.Transaction, items []int) (map[int][]*models.Transaction, float64) {
	projectedItemTransactionMap := make(map[int][]*models.Transaction)
	totalUtility := 0.0

	// Duyệt qua từng item trong items
	for _, item := range items {
		transactions, found := itemTransactionMap[item]
		if !found {
			continue
		}

		// Duyệt qua từng giao dịch chứa item
		for _, transaction := range transactions {
			// Kiểm tra nếu giao dịch chứa tất cả các items trong danh sách `items`
			if containsAllItems(transaction.Items, items) {
				var projectedItems []int
				var projectedUtilities []float64
				transactionUtility := 0.0

				// Lấy danh sách items và utilities từ transaction
				for i := 0; i < len(transaction.Items); i++ {
					projectedItems = append(projectedItems, transaction.Items[i])
					projectedUtilities = append(projectedUtilities, transaction.Utilities[i])
				}

				// Tính tổng utility cho tất cả các items trong danh sách `items`
				for _, item := range items {
					itemIndex := indexOf(transaction.Items, item)
					if itemIndex != -1 {
						transactionUtility += transaction.Utilities[itemIndex]
					}
				}

				// Cộng thêm utility của giao dịch này vào `totalUtility`
				totalUtility += transactionUtility

				// Tạo transaction đã projected và thêm vào map
				projectedTransaction := models.NewTransaction(projectedItems, projectedUtilities, transactionUtility)
				projectedItemTransactionMap[item] = append(projectedItemTransactionMap[item], projectedTransaction)
			}
		}
	}

	return projectedItemTransactionMap, totalUtility
}

// func (s *SearchAlgorithms) Search(eta []int, X map[int]bool, transactions []*models.Transaction, primary []int, secondary []int, minU float64) {
// 	if len(primary) == 0 {
// 		return
// 	}

// 	for _, item := range primary {

// 		s.Beta = copyMap(X)
// 		s.Beta[item] = true
// 		s.ItemList = mapKeys(s.Beta)

// 		projectedDB, utilityBeta := s.projectDatabase(transactions, s.ItemList)

// 		if utilityBeta >= minU {
// 			fmt.Printf("U(%d) = %.2f >= %.2f HUI Found: %v\n", item, utilityBeta, minU, s.Beta)
// 			s.HighUtilityItemsets = append(s.HighUtilityItemsets, models.NewHighUtilityItemset(s.ItemList, utilityBeta))
// 		} else {
// 			fmt.Printf("%.2f < %.2f so %v is not a HUI.\n", utilityBeta, minU, s.Beta)
// 		}

// 		if utilityBeta > minU {
// 			s.SearchN(eta, s.Beta, projectedDB, minU)
// 		}

// 		s.FilteredPrimary = []int{}
// 		s.FilteredSecondary = []int{}
// 		utility.CalculateRSUForAllItem(projectedDB, s.ItemList, secondary, s.UtilityArray)
// 		utility.CalculateRLUForAllItem(projectedDB, s.ItemList, secondary, s.UtilityArray)
// 		for i, secItem := range secondary {

// 			if secItem == item {
// 				continue
// 			}

// 			if i > indexOf(secondary, item) {
// 				rsu := s.UtilityArray.GetRSU(secItem)
// 				rlu := s.UtilityArray.GetRLU(secItem)

// 				if rsu >= minU {
// 					s.FilteredPrimary = append(s.FilteredPrimary, secItem)
// 				}
// 				if rlu >= minU {
// 					s.FilteredSecondary = append(s.FilteredSecondary, secItem)
// 				}
// 			}
// 		}

// 		fmt.Printf("Beta= %v\n", s.Beta)
// 		fmt.Printf("Primary%v = %v\n", s.ItemList, s.FilteredPrimary)
// 		fmt.Printf("Secondary%v = %v\n", s.ItemList, s.FilteredSecondary)

// 		s.Search(eta, s.Beta, projectedDB, s.FilteredPrimary, s.FilteredSecondary, minU)
// 	}
// }

// func (s *SearchAlgorithms) SearchN(eta []int, beta map[int]bool, transactions []*models.Transaction, minU float64) {
// 	if len(eta) == 0 {
// 		return
// 	}

// 	for _, item := range eta {
// 		betaNew := copyMap(beta)
// 		betaNew[item] = true

// 		itemList := mapKeys(betaNew)

// 		projectedDBNew, utilityBetaNew := s.projectDatabase(transactions, itemList)

// 		if utilityBetaNew >= minU {
// 			fmt.Printf("U(%d) = %.2f >= %.2f HUI Found: %v\n", item, utilityBetaNew, minU, betaNew)
// 			s.HighUtilityItemsets = append(s.HighUtilityItemsets, models.NewHighUtilityItemset(mapKeys(betaNew), utilityBetaNew))
// 		} else {
// 			fmt.Printf("%.2f < %.2f so %v is not a HUI.\n", utilityBetaNew, minU, betaNew)
// 		}

// 		itemIndex := indexOf(eta, item)
// 		filteredPrimary := []int{}
// 		utility.CalculateRSUForAllItem(projectedDBNew, itemList, eta, s.UtilityArray)
// 		for _, secItem := range eta {
// 			if secItem == item {
// 				continue
// 			}
// 			if indexOf(eta, secItem) > itemIndex {
// 				rsu := s.UtilityArray.GetRSU(secItem)
// 				if rsu >= minU {
// 					filteredPrimary = append(filteredPrimary, secItem)
// 				}
// 			}
// 		}
// 		fmt.Printf("Primary = %v\n", filteredPrimary)
// 		s.SearchN(filteredPrimary, betaNew, projectedDBNew, minU)
// 	}
// }

// func (s *SearchAlgorithms) projectDatabase(transactions []*models.Transaction, items []int) []*models.Transaction {
// 	var projectedDB []*models.Transaction

// 	for _, transaction := range transactions {
// 		if containsAllItems(transaction.Items, items) {
// 			var projectedItems []int
// 			var projectedUtilities []float64
// 			lastItemIndex := -1

// 			for _, item := range items {
// 				itemIndex := indexOf(transaction.Items, item)
// 				if itemIndex > lastItemIndex {
// 					lastItemIndex = itemIndex
// 				}
// 			}

// 			for i := lastItemIndex + 1; i < len(transaction.Items); i++ {
// 				projectedItems = append(projectedItems, transaction.Items[i])
// 				projectedUtilities = append(projectedUtilities, transaction.Utilities[i])
// 			}

// 			if len(projectedItems) > 0 {
// 				projectedDB = append(projectedDB, models.NewTransaction(projectedItems, projectedUtilities, calculateTransactionUtility(projectedUtilities)))
// 			}
// 		}
// 	}

//		return projectedDB
//	}
func (s *SearchAlgorithms) projectDatabase(transactions []*models.Transaction, items []int) ([]*models.Transaction, float64) {
	var projectedDB []*models.Transaction
	totalUtility := 0.0

	for _, transaction := range transactions {
		if containsAllItems(transaction.Items, items) {
			var projectedItems []int
			var projectedUtilities []float64
			lastItemIndex := -1

			for _, item := range items {
				itemIndex := indexOf(transaction.Items, item)
				if itemIndex > lastItemIndex {
					lastItemIndex = itemIndex
				}
			}

			for i := 0; i < len(transaction.Items); i++ {
				projectedItems = append(projectedItems, transaction.Items[i])
				projectedUtilities = append(projectedUtilities, transaction.Utilities[i])
			}

			if len(projectedItems) > 0 {
				transactionUtility := calculateTransactionUtility(projectedUtilities)
				projectedDB = append(projectedDB, models.NewTransaction(projectedItems, projectedUtilities, transactionUtility))

				for _, item := range items {
					index := indexOf(transaction.Items, item)
					if index != -1 {
						totalUtility += transaction.Utilities[index]
					}
				}
			}
		}
	}
	return projectedDB, totalUtility
}

func (s *SearchAlgorithms) calculateUtility(transactions []*models.Transaction, itemset map[int]bool) float64 {
	totalUtility := 0.0

	for _, transaction := range transactions {
		if containsAllItemsMap(transaction.Items, itemset) {
			for item := range itemset {
				index := indexOf(transaction.Items, item)
				if index != -1 {
					itemUtility := transaction.Utilities[index]
					totalUtility += itemUtility
				}
			}
		}
	}

	return totalUtility
}

func copyMap(original map[int]bool) map[int]bool {
	copy := make(map[int]bool)
	for k, v := range original {
		copy[k] = v
	}
	return copy
}

func mapKeys(m map[int]bool) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func containsAllItems(items []int, itemset []int) bool {
	for _, item := range itemset {
		if indexOf(items, item) == -1 {
			return false
		}
	}
	return true
}

func indexOf(items []int, item int) int {
	for i, v := range items {
		if v == item {
			return i
		}
	}
	return -1
}

func containsAllItemsMap(items []int, itemset map[int]bool) bool {
	for item := range itemset {
		if indexOf(items, item) == -1 {
			return false
		}
	}
	return true
}

//	func (s *SearchAlgorithms) printProjectedDatabase(projectedDB []*models.Transaction, items []int) {
//		fmt.Printf("\nProjected Database after items %v:\n", items)
//		for _, transaction := range projectedDB {
//			fmt.Printf("Items: %v, Utilities: %v, Transaction Utility: %.2f\n",
//				transaction.Items, transaction.Utilities, calculateTransactionUtility(transaction.Utilities))
//		}
//		fmt.Println("----------------------------------")
//	}
func (s *SearchAlgorithms) printProjectedDatabase(projectedItemTransactionMap map[int][]*models.Transaction) {
	fmt.Println("\nProjected Item Transaction Map:")

	for item, transactions := range projectedItemTransactionMap {
		fmt.Printf("Item %d:\n", item)
		for _, transaction := range transactions {
			fmt.Printf("  Items: %v, Utilities: %v, Transaction Utility: %.2f\n",
				transaction.Items, transaction.Utilities, transaction.TransactionUtility)
		}
		fmt.Println("----------------------------------")
	}
}

func calculateTransactionUtility(utilities []float64) float64 { // Chuyển sang float64
	totalUtility := 0.0
	for _, utility := range utilities {
		totalUtility += utility
	}
	return totalUtility
}
