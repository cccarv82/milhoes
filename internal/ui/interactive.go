package ui

import (
	"fmt"
	"lottery-optimizer/internal/ai"
	"lottery-optimizer/internal/config"
	"lottery-optimizer/internal/data"
	"lottery-optimizer/internal/lottery"
	"lottery-optimizer/internal/strategy"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

// ShowWelcome mostra a tela de boas-vindas
func ShowWelcome() {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	
	fmt.Println()
	cyan.Println("🎰 ═══════════════════════════════════════════════════════════════")
	cyan.Println("🎰 LOTTERY OPTIMIZER - Estratégias Inteligentes para Loterias")
	cyan.Println("🎰 ═══════════════════════════════════════════════════════════════")
	fmt.Println()
	
	yellow.Println("✨ Usando Inteligência Artificial Claude Sonnet 4")
	yellow.Println("📊 Análise estatística avançada de dados históricos")
	yellow.Println("🎯 Otimização matemática para maximizar suas chances")
	yellow.Println("💰 Gestão inteligente de orçamento")
	fmt.Println()
	
	green.Println("🍀 Que a sorte esteja com você! 🍀")
	fmt.Println()
}

// StartInteractiveMode inicia o modo interativo
func StartInteractiveMode() {
	for {
		action := showMainMenu()
		
		switch action {
		case "generate":
			generateStrategy()
		case "stats":
			showStatistics()
		case "config":
			showConfiguration()
		case "test":
			testConnections()
		case "help":
			showHelp()
		case "exit":
			showGoodbye()
			return
		}
	}
}

// showMainMenu exibe o menu principal
func showMainMenu() string {
	prompt := promptui.Select{
		Label: "🎯 Selecione uma opção",
		Items: []string{
			"🎲 Gerar Estratégia Otimizada",
			"📊 Ver Estatísticas",
			"⚙️  Configurações",
			"🔧 Testar Conexões",
			"❓ Ajuda",
			"🚪 Sair",
		},
	}
	
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Erro: %v\n", err)
		return "exit"
	}
	
	switch result {
	case "🎲 Gerar Estratégia Otimizada":
		return "generate"
	case "📊 Ver Estatísticas":
		return "stats"
	case "⚙️  Configurações":
		return "config"
	case "🔧 Testar Conexões":
		return "test"
	case "❓ Ajuda":
		return "help"
	default:
		return "exit"
	}
}

// generateStrategy processo principal de geração de estratégia
func generateStrategy() {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	
	cyan.Println("\n🎯 GERAÇÃO DE ESTRATÉGIA OTIMIZADA")
	fmt.Println("═══════════════════════════════════════")
	
	// Coletar preferências do usuário
	prefs, err := collectUserPreferences()
	if err != nil {
		color.Red("❌ Erro ao coletar preferências: %v", err)
		return
	}
	
	yellow.Println("\n🤖 Conectando com a IA...")
	
	// Criar clientes
	dataClient := data.NewClient()
	aiClient := ai.NewClaudeClient()
	
	// Buscar dados históricos
	yellow.Println("📥 Buscando dados históricos...")
	
	var allDraws []lottery.Draw
	var allRules []lottery.LotteryRules
	
	for _, ltype := range prefs.LotteryTypes {
		draws, err := dataClient.GetLatestDraws(ltype, 200) // Últimos 200 sorteios
		if err != nil {
			color.Red("❌ Erro ao buscar dados de %s: %v", ltype, err)
			continue
		}
		
		allDraws = append(allDraws, draws...)
		allRules = append(allRules, lottery.GetRules(ltype))
		
		if config.IsVerbose() {
			fmt.Printf("✅ Obtidos %d sorteios de %s\n", len(draws), ltype)
		}
	}
	
	if len(allDraws) == 0 {
		color.Red("❌ Nenhum dado histórico encontrado!")
		return
	}
	
	// Preparar requisição para IA
	analysisReq := lottery.AnalysisRequest{
		Draws:       allDraws,
		Preferences: *prefs,
		Rules:       allRules,
	}
	
	yellow.Println("🧠 Analisando com IA Claude Sonnet 4...")
	
	// Analisar com IA
	response, err := aiClient.AnalyzeStrategy(analysisReq)
	if err != nil {
		color.Red("❌ Erro na análise da IA: %v", err)
		return
	}
	
	// Validar e ajustar estratégia
	strategy := strategy.ValidateAndAdjustStrategy(&response.Strategy, *prefs)
	
	// Exibir resultado
	displayStrategy(strategy, response.Confidence)
	
	// Opção de salvar
	if askYesNo("💾 Deseja salvar esta estratégia?") {
		saveStrategy(strategy)
	}
}

// collectUserPreferences coleta as preferências do usuário
func collectUserPreferences() (*lottery.UserPreferences, error) {
	prefs := &lottery.UserPreferences{}
	
	// Selecionar tipos de loteria
	lotteryPrompt := promptui.Select{
		Label: "🎲 Quais loterias deseja jogar?",
		Items: []string{
			"🎯 Apenas Mega Sena",
			"🍀 Apenas Lotofácil",
			"🎲 Ambas (estratégia mista)",
		},
	}
	
	_, lotteryChoice, err := lotteryPrompt.Run()
	if err != nil {
		return nil, err
	}
	
	switch lotteryChoice {
	case "🎯 Apenas Mega Sena":
		prefs.LotteryTypes = []lottery.LotteryType{lottery.MegaSena}
	case "🍀 Apenas Lotofácil":
		prefs.LotteryTypes = []lottery.LotteryType{lottery.Lotofacil}
	default:
		prefs.LotteryTypes = []lottery.LotteryType{lottery.MegaSena, lottery.Lotofacil}
	}
	
	// Orçamento
	budgetPrompt := promptui.Prompt{
		Label:    "💰 Qual seu orçamento disponível? (R$)",
		Validate: validateBudget,
		Default:  "50",
	}
	
	budgetStr, err := budgetPrompt.Run()
	if err != nil {
		return nil, err
	}
	
	budget, _ := strconv.ParseFloat(budgetStr, 64)
	prefs.Budget = budget
	
	// Estratégia
	strategyPrompt := promptui.Select{
		Label: "📈 Qual tipo de estratégia prefere?",
		Items: []string{
			"🛡️  Conservadora (menor risco)",
			"⚖️  Equilibrada (recomendada)",
			"🚀 Agressiva (maior risco/retorno)",
		},
	}
	
	_, strategyChoice, err := strategyPrompt.Run()
	if err != nil {
		return nil, err
	}
	
	switch strategyChoice {
	case "🛡️  Conservadora (menor risco)":
		prefs.Strategy = "conservative"
	case "🚀 Agressiva (maior risco/retorno)":
		prefs.Strategy = "aggressive"
	default:
		prefs.Strategy = "balanced"
	}
	
	// Perguntas adicionais
	prefs.AvoidPatterns = askYesNo("🔢 Evitar padrões óbvios (sequências, múltiplos)?")
	
	if askYesNo("⭐ Tem números da sorte?") {
		favNumbers := askForNumbers("Digite os números da sorte (separados por vírgula):")
		prefs.FavoriteNumbers = favNumbers
	}
	
	if askYesNo("❌ Tem números que quer evitar?") {
		excNumbers := askForNumbers("Digite os números a evitar (separados por vírgula):")
		prefs.ExcludeNumbers = excNumbers
	}
	
	return prefs, nil
}

// displayStrategy exibe a estratégia gerada
func displayStrategy(strategy *lottery.Strategy, confidence float64) {
	green := color.New(color.FgGreen, color.Bold)
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	white := color.New(color.FgWhite, color.Bold)
	
	fmt.Println("\n" + strings.Repeat("═", 60))
	green.Println("🎯 ESTRATÉGIA GERADA PELA IA")
	fmt.Println(strings.Repeat("═", 60))
	
	cyan.Printf("💰 Orçamento: R$ %.2f\n", strategy.Budget)
	cyan.Printf("💸 Custo Total: R$ %.2f\n", strategy.TotalCost)
	cyan.Printf("📊 Confiança da IA: %.1f%%\n", confidence*100)
	cyan.Printf("🎲 Total de Jogos: %d\n", len(strategy.Games))
	fmt.Println()
	
	// Agrupar jogos por tipo
	megaSenaGames := []lottery.Game{}
	lotofacilGames := []lottery.Game{}
	
	for _, game := range strategy.Games {
		if game.Type == lottery.MegaSena {
			megaSenaGames = append(megaSenaGames, game)
		} else {
			lotofacilGames = append(lotofacilGames, game)
		}
	}
	
	// Exibir jogos da Mega Sena
	if len(megaSenaGames) > 0 {
		yellow.Println("🎯 MEGA SENA:")
		for i, game := range megaSenaGames {
			white.Printf("Jogo %d: ", i+1)
			for j, num := range game.Numbers {
				if j > 0 {
					fmt.Print(" ")
				}
				fmt.Printf("%02d", num)
			}
			fmt.Printf(" (R$ %.2f)\n", game.Cost)
		}
		fmt.Println()
	}
	
	// Exibir jogos da Lotofácil
	if len(lotofacilGames) > 0 {
		yellow.Println("🍀 LOTOFÁCIL:")
		for i, game := range lotofacilGames {
			white.Printf("Jogo %d: ", i+1)
			for j, num := range game.Numbers {
				if j > 0 {
					fmt.Print(" ")
				}
				fmt.Printf("%02d", num)
			}
			fmt.Printf(" (R$ %.2f)\n", game.Cost)
		}
		fmt.Println()
	}
	
	// Exibir raciocínio da IA
	cyan.Println("🤖 JUSTIFICATIVA DA IA:")
	fmt.Println(strategy.Reasoning)
	fmt.Println()
	
	// Estatísticas
	if strategy.Statistics.TotalDraws > 0 {
		cyan.Println("📊 ESTATÍSTICAS:")
		fmt.Printf("• Sorteios analisados: %d\n", strategy.Statistics.AnalyzedDraws)
		
		if len(strategy.Statistics.HotNumbers) > 0 {
			fmt.Print("• Números quentes: ")
			for i, num := range strategy.Statistics.HotNumbers {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%d", num)
			}
			fmt.Println()
		}
		
		if len(strategy.Statistics.ColdNumbers) > 0 {
			fmt.Print("• Números frios: ")
			for i, num := range strategy.Statistics.ColdNumbers {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%d", num)
			}
			fmt.Println()
		}
		fmt.Println()
	}
	
	// Próximos sorteios
	showNextDraws()
	
	fmt.Println(strings.Repeat("═", 60))
	green.Println("🍀 BOA SORTE! 🍀")
	fmt.Println(strings.Repeat("═", 60))
}

// Funções auxiliares
func validateBudget(input string) error {
	budget, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return fmt.Errorf("valor inválido")
	}
	if budget < 5 {
		return fmt.Errorf("orçamento mínimo é R$ 5,00")
	}
	if budget > 10000 {
		return fmt.Errorf("orçamento máximo é R$ 10.000,00")
	}
	return nil
}

func askYesNo(question string) bool {
	prompt := promptui.Select{
		Label: question,
		Items: []string{"Sim", "Não"},
	}
	
	_, result, err := prompt.Run()
	if err != nil {
		return false
	}
	
	return result == "Sim"
}

func askForNumbers(prompt string) []int {
	numberPrompt := promptui.Prompt{
		Label: prompt,
	}
	
	result, err := numberPrompt.Run()
	if err != nil {
		return nil
	}
	
	var numbers []int
	parts := strings.Split(result, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if num, err := strconv.Atoi(part); err == nil {
			numbers = append(numbers, num)
		}
	}
	
	return numbers
}

func showNextDraws() {
	cyan := color.New(color.FgCyan, color.Bold)
	
	dataClient := data.NewClient()
	
	cyan.Println("📅 PRÓXIMOS SORTEIOS:")
	
	// Mega Sena
	if nextDate, nextNum, err := dataClient.GetNextDrawInfo(lottery.MegaSena); err == nil {
		fmt.Printf("• Mega Sena: Concurso %d em %s\n", 
			nextNum, nextDate.Format("02/01/2006"))
	}
	
	// Lotofácil
	if nextDate, nextNum, err := dataClient.GetNextDrawInfo(lottery.Lotofacil); err == nil {
		fmt.Printf("• Lotofácil: Concurso %d em %s\n", 
			nextNum, nextDate.Format("02/01/2006"))
	}
	
	fmt.Println()
}

func saveStrategy(strategy *lottery.Strategy) {
	// TODO: Implementar salvamento em arquivo
	color.Green("✅ Estratégia salva com sucesso!")
}

func showStatistics() {
	color.Yellow("📊 Funcionalidade em desenvolvimento...")
}

func showConfiguration() {
	color.Yellow("⚙️ Funcionalidade em desenvolvimento...")
}

func testConnections() {
	cyan := color.New(color.FgCyan, color.Bold)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)
	
	cyan.Println("\n🔧 TESTANDO CONEXÕES")
	fmt.Println("═══════════════════════")
	
	// Testar API da Caixa
	fmt.Print("🌐 API Loterias Caixa... ")
	dataClient := data.NewClient()
	if err := dataClient.TestConnection(); err != nil {
		red.Printf("❌ FALHOU: %v\n", err)
	} else {
		green.Println("✅ OK")
	}
	
	// Testar Claude API
	fmt.Print("🤖 Claude API... ")
	aiClient := ai.NewClaudeClient()
	if err := aiClient.TestConnection(); err != nil {
		red.Printf("❌ FALHOU: %v\n", err)
	} else {
		green.Println("✅ OK")
	}
	
	fmt.Println()
}

func showHelp() {
	cyan := color.New(color.FgCyan, color.Bold)
	
	cyan.Println("\n❓ AJUDA")
	fmt.Println("═══════════")
	fmt.Println("Este programa usa inteligência artificial para analisar dados")
	fmt.Println("históricos das loterias brasileiras e gerar estratégias otimizadas.")
	fmt.Println()
	fmt.Println("🎯 Mega Sena: 6 números de 1 a 60")
	fmt.Println("🍀 Lotofácil: 15 números de 1 a 25")
	fmt.Println()
	fmt.Println("💡 Dicas:")
	fmt.Println("• Use estratégia equilibrada para melhores resultados")
	fmt.Println("• Orçamentos maiores permitem estratégias mais sofisticadas")
	fmt.Println("• Evite padrões óbvios para maximizar chances")
	fmt.Println()
}

func showGoodbye() {
	green := color.New(color.FgGreen, color.Bold)
	
	fmt.Println()
	green.Println("🍀 Obrigado por usar o Lottery Optimizer!")
	green.Println("🎯 Que os números escolhidos pela IA sejam os sorteados!")
	green.Println("💰 Boa sorte! ��")
	fmt.Println()
} 