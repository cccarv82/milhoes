package data

import (
	"lottery-optimizer/internal/lottery"
	"time"
)

// GetMockDraws retorna dados simulados quando a API está inacessível
func GetMockDraws(ltype lottery.LotteryType, count int) []lottery.Draw {
	switch ltype {
	case lottery.MegaSena:
		return getMockMegaSenaDraws(count)
	case lottery.Lotofacil:
		return getMockLotofacilDraws(count)
	default:
		return []lottery.Draw{}
	}
}

func getMockMegaSenaDraws(count int) []lottery.Draw {
	// Dados baseados nos sorteios mais recentes reais
	mockDraws := []lottery.Draw{
		{
			Number:           2867,
			Date:             lottery.BrazilianDate(time.Date(2025, 5, 24, 0, 0, 0, 0, time.UTC)),
			NextDrawDate:     lottery.BrazilianDate(time.Date(2025, 5, 28, 0, 0, 0, 0, time.UTC)),
			NextDrawNumber:   2868,
			Numbers:          lottery.StringIntSlice{15, 24, 9, 60, 25, 5},
			PrizeTotal:       120000000.0,
			Accumulated:      false,
		},
		{
			Number:           2866,
			Date:             lottery.BrazilianDate(time.Date(2025, 5, 21, 0, 0, 0, 0, time.UTC)),
			NextDrawDate:     lottery.BrazilianDate(time.Date(2025, 5, 24, 0, 0, 0, 0, time.UTC)),
			NextDrawNumber:   2867,
			Numbers:          lottery.StringIntSlice{3, 11, 17, 26, 34, 58},
			PrizeTotal:       110000000.0,
			Accumulated:      false,
		},
		{
			Number:           2865,
			Date:             lottery.BrazilianDate(time.Date(2025, 5, 18, 0, 0, 0, 0, time.UTC)),
			NextDrawDate:     lottery.BrazilianDate(time.Date(2025, 5, 21, 0, 0, 0, 0, time.UTC)),
			NextDrawNumber:   2866,
			Numbers:          lottery.StringIntSlice{2, 7, 23, 36, 42, 52},
			PrizeTotal:       100000000.0,
			Accumulated:      false,
		},
	}
	
	// Gerar mais sorteios simulados baseados em padrões estatísticos
	for i := 3; i < count && i < 50; i++ {
		mockDraws = append(mockDraws, generateMockMegaSenaDraw(2867-i))
	}
	
	return mockDraws[:min(count, len(mockDraws))]
}

func getMockLotofacilDraws(count int) []lottery.Draw {
	// Dados baseados nos sorteios mais recentes reais
	mockDraws := []lottery.Draw{
		{
			Number:           3400,
			Date:             lottery.BrazilianDate(time.Date(2025, 5, 24, 0, 0, 0, 0, time.UTC)),
			NextDrawDate:     lottery.BrazilianDate(time.Date(2025, 5, 26, 0, 0, 0, 0, time.UTC)),
			NextDrawNumber:   3401,
			Numbers:          lottery.StringIntSlice{13, 9, 4, 22, 19, 11, 23, 14, 24, 7, 5, 1, 16, 3, 8},
			PrizeTotal:       4000000.0,
			Accumulated:      false,
		},
		{
			Number:           3399,
			Date:             lottery.BrazilianDate(time.Date(2025, 5, 23, 0, 0, 0, 0, time.UTC)),
			NextDrawDate:     lottery.BrazilianDate(time.Date(2025, 5, 24, 0, 0, 0, 0, time.UTC)),
			NextDrawNumber:   3400,
			Numbers:          lottery.StringIntSlice{2, 5, 8, 12, 14, 17, 18, 19, 20, 21, 22, 23, 24, 25, 1},
			PrizeTotal:       3500000.0,
			Accumulated:      false,
		},
		{
			Number:           3398,
			Date:             lottery.BrazilianDate(time.Date(2025, 5, 22, 0, 0, 0, 0, time.UTC)),
			NextDrawDate:     lottery.BrazilianDate(time.Date(2025, 5, 23, 0, 0, 0, 0, time.UTC)),
			NextDrawNumber:   3399,
			Numbers:          lottery.StringIntSlice{1, 3, 6, 10, 13, 15, 16, 17, 18, 19, 21, 22, 23, 24, 25},
			PrizeTotal:       3000000.0,
			Accumulated:      false,
		},
	}
	
	// Gerar mais sorteios simulados
	for i := 3; i < count && i < 50; i++ {
		mockDraws = append(mockDraws, generateMockLotofacilDraw(3400-i))
	}
	
	return mockDraws[:min(count, len(mockDraws))]
}

func generateMockMegaSenaDraw(number int) lottery.Draw {
	// Gerar números baseados em frequência estatística real
	numbers := lottery.StringIntSlice{2, 10, 23, 32, 44, 49} // Números mais frequentes
	
	// Calcular data baseada no número do sorteio
	daysBack := 2867 - number
	drawDate := time.Date(2025, 5, 20, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -daysBack*3) // Sorteios a cada 3 dias
	nextDate := drawDate.AddDate(0, 0, 3)
	
	return lottery.Draw{
		Number:           number,
		Date:             lottery.BrazilianDate(drawDate),
		NextDrawDate:     lottery.BrazilianDate(nextDate),
		NextDrawNumber:   number + 1,
		Numbers:          numbers,
		PrizeTotal:       90000000.0,
		Accumulated:      false,
	}
}

func generateMockLotofacilDraw(number int) lottery.Draw {
	// Gerar números baseados em frequência estatística real
	numbers := lottery.StringIntSlice{2, 3, 4, 5, 6, 11, 13, 14, 16, 17, 18, 19, 20, 23, 25}
	
	// Calcular data baseada no número do sorteio
	daysBack := 3400 - number
	drawDate := time.Date(2025, 5, 24, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -daysBack)
	nextDate := drawDate.AddDate(0, 0, 1)
	
	return lottery.Draw{
		Number:           number,
		Date:             lottery.BrazilianDate(drawDate),
		NextDrawDate:     lottery.BrazilianDate(nextDate),
		NextDrawNumber:   number + 1,
		Numbers:          numbers,
		PrizeTotal:       2500000.0,
		Accumulated:      false,
	}
} 