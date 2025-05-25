package lottery

import (
	"fmt"
	"time"
)

// LotteryType tipos de loteria suportados
type LotteryType string

const (
	MegaSena  LotteryType = "megasena"
	Lotofacil LotteryType = "lotofacil"
)

// LotteryRules regras de cada tipo de loteria
type LotteryRules struct {
	Name           string
	MinNumbers     int
	MaxNumbers     int
	NumberRange    int
	BasePrice      float64
	DrawDays       []time.Weekday
	ResultNumbers  int
}

// GetRules retorna as regras para cada tipo de loteria
func GetRules(ltype LotteryType) LotteryRules {
	switch ltype {
	case MegaSena:
		return LotteryRules{
			Name:           "Mega Sena",
			MinNumbers:     6,
			MaxNumbers:     15,
			NumberRange:    60,
			BasePrice:      5.00,
			DrawDays:       []time.Weekday{time.Wednesday, time.Saturday},
			ResultNumbers:  6,
		}
	case Lotofacil:
		return LotteryRules{
			Name:           "Lotofácil",
			MinNumbers:     15,
			MaxNumbers:     20,
			NumberRange:    25,
			BasePrice:      3.00,
			DrawDays:       []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
			ResultNumbers:  20,
		}
	default:
		return LotteryRules{}
	}
}

// Draw representa um sorteio individual
type Draw struct {
	Number         int       `json:"numero"`
	Date           time.Time `json:"dataApuracao"`
	Numbers        []int     `json:"dezenasSorteadasOrdemSorteio"`
	Winners        []Winner  `json:"listaRateioPremio"`
	PrizeTotal     float64   `json:"valorArrecadado"`
	NextDrawNumber int       `json:"numeroConcursoProximo"`
	NextDrawDate   time.Time `json:"dataProximoConcurso"`
	Accumulated    bool      `json:"acumulado"`
}

// Winner representa ganhadores por faixa de prêmio
type Winner struct {
	Description string  `json:"descricaoFaixa"`
	Winners     int     `json:"numeroDeGanhadores"`
	Prize       float64 `json:"valorPremio"`
}

// Game representa um jogo individual
type Game struct {
	Type           LotteryType `json:"type"`
	Numbers        []int       `json:"numbers"`
	Cost           float64     `json:"cost"`
	ExpectedReturn float64     `json:"expectedReturn"`
	Probability    float64     `json:"probability"`
}

// Strategy representa uma estratégia completa
type Strategy struct {
	Games          []Game    `json:"games"`
	TotalCost      float64   `json:"totalCost"`
	Budget         float64   `json:"budget"`
	ExpectedReturn float64   `json:"expectedReturn"`
	Reasoning      string    `json:"reasoning"`
	Statistics     Stats     `json:"statistics"`
	CreatedAt      time.Time `json:"createdAt"`
}

// Stats estatísticas de análise
type Stats struct {
	TotalDraws       int               `json:"totalDraws"`
	AnalyzedDraws    int               `json:"analyzedDraws"`
	NumberFrequency  map[int]int       `json:"numberFrequency"`
	SumDistribution  map[int]int       `json:"sumDistribution"`
	RecentTrends     []int             `json:"recentTrends"`
	HotNumbers       []int             `json:"hotNumbers"`
	ColdNumbers      []int             `json:"coldNumbers"`
	Patterns         map[string]string `json:"patterns"`
}

// UserPreferences preferências do usuário
type UserPreferences struct {
	LotteryTypes   []LotteryType `json:"lotteryTypes"`
	Budget         float64       `json:"budget"`
	Strategy       string        `json:"strategy"` // conservative, balanced, aggressive
	MaxGames       int           `json:"maxGames"`
	AvoidPatterns  bool          `json:"avoidPatterns"`
	FavoriteNumbers []int         `json:"favoriteNumbers"`
	ExcludeNumbers []int         `json:"excludeNumbers"`
}

// AnalysisRequest requisição para análise da IA
type AnalysisRequest struct {
	Draws       []Draw          `json:"draws"`
	Preferences UserPreferences `json:"preferences"`
	Rules       []LotteryRules  `json:"rules"`
}

// AnalysisResponse resposta da análise da IA
type AnalysisResponse struct {
	Strategy     Strategy `json:"strategy"`
	Confidence   float64  `json:"confidence"`
	Alternatives []Game   `json:"alternatives"`
	Warnings     []string `json:"warnings"`
}

// ValidateGame valida se um jogo está correto conforme as regras
func ValidateGame(game Game) error {
	rules := GetRules(game.Type)
	
	if len(game.Numbers) < rules.MinNumbers || len(game.Numbers) > rules.MaxNumbers {
		return fmt.Errorf("número de dezenas inválido para %s: deve estar entre %d e %d", 
			rules.Name, rules.MinNumbers, rules.MaxNumbers)
	}
	
	for _, num := range game.Numbers {
		if num < 1 || num > rules.NumberRange {
			return fmt.Errorf("número %d inválido para %s: deve estar entre 1 e %d", 
				num, rules.Name, rules.NumberRange)
		}
	}
	
	// Verificar se não há números repetidos
	seen := make(map[int]bool)
	for _, num := range game.Numbers {
		if seen[num] {
			return fmt.Errorf("número %d repetido no jogo", num)
		}
		seen[num] = true
	}
	
	return nil
}

// CalculateGameCost calcula o custo de um jogo baseado na quantidade de números
func CalculateGameCost(ltype LotteryType, numCount int) float64 {
	rules := GetRules(ltype)
	
	if numCount == rules.MinNumbers {
		return rules.BasePrice
	}
	
	// Cálculo combinatório para jogos com mais números
	// Custo = BasePrice * C(numCount, MinNumbers)
	combinations := calculateCombinations(numCount, rules.MinNumbers)
	return rules.BasePrice * float64(combinations)
}

// calculateCombinations calcula C(n, r) = n! / (r! * (n-r)!)
func calculateCombinations(n, r int) int {
	if r > n {
		return 0
	}
	if r == 0 || r == n {
		return 1
	}
	
	result := 1
	for i := 0; i < r; i++ {
		result = result * (n - i) / (i + 1)
	}
	return result
} 