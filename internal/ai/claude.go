package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"lottery-optimizer/internal/config"
	"lottery-optimizer/internal/lottery"
	"net/http"
	"time"
)

// ClaudeClient cliente para API do Claude
type ClaudeClient struct {
	apiKey     string
	baseURL    string
	model      string
	maxTokens  int
	httpClient *http.Client
}

// ClaudeRequest estrutura da requisição para Claude
type ClaudeRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
}

// Message mensagem para Claude
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeResponse resposta da API do Claude
type ClaudeResponse struct {
	Content []Content `json:"content"`
	Usage   Usage     `json:"usage"`
}

// Content conteúdo da resposta
type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Usage uso de tokens
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// NewClaudeClient cria um novo cliente para Claude
func NewClaudeClient() *ClaudeClient {
	return &ClaudeClient{
		apiKey:    config.GetClaudeAPIKey(),
		baseURL:   "https://api.anthropic.com/v1/messages",
		model:     config.GetClaudeModel(),
		maxTokens: config.GetMaxTokens(),
		httpClient: &http.Client{
			Timeout: time.Duration(config.GlobalConfig.Claude.TimeoutSec) * time.Second,
		},
	}
}

// AnalyzeStrategy usa Claude para analisar dados e gerar estratégia
func (c *ClaudeClient) AnalyzeStrategy(request lottery.AnalysisRequest) (*lottery.AnalysisResponse, error) {
	prompt := c.buildAnalysisPrompt(request)
	
	claudeReq := ClaudeRequest{
		Model:     c.model,
		MaxTokens: c.maxTokens,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}
	
	reqBody, err := json.Marshal(claudeReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar requisição: %w", err)
	}
	
	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API retornou status %d", resp.StatusCode)
	}
	
	var claudeResp ClaudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&claudeResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}
	
	if len(claudeResp.Content) == 0 {
		return nil, fmt.Errorf("resposta vazia do Claude")
	}
	
	// Parsear a resposta JSON do Claude
	var analysisResp lottery.AnalysisResponse
	if err := json.Unmarshal([]byte(claudeResp.Content[0].Text), &analysisResp); err != nil {
		// Se não conseguir parsear como JSON, criar resposta básica
		analysisResp = lottery.AnalysisResponse{
			Strategy: lottery.Strategy{
				Reasoning: claudeResp.Content[0].Text,
				CreatedAt: time.Now(),
			},
			Confidence: 0.7,
		}
	}
	
	if config.IsVerbose() {
		fmt.Printf("Tokens usados: %d input + %d output = %d total\n", 
			claudeResp.Usage.InputTokens, claudeResp.Usage.OutputTokens, 
			claudeResp.Usage.InputTokens+claudeResp.Usage.OutputTokens)
	}
	
	return &analysisResp, nil
}

// buildAnalysisPrompt constrói o prompt para análise
func (c *ClaudeClient) buildAnalysisPrompt(request lottery.AnalysisRequest) string {
	// Serializar dados históricos
	drawsJSON, _ := json.MarshalIndent(request.Draws, "", "  ")
	prefsJSON, _ := json.MarshalIndent(request.Preferences, "", "  ")
	rulesJSON, _ := json.MarshalIndent(request.Rules, "", "  ")
	
	prompt := fmt.Sprintf(`Você é um especialista em análise estatística de loterias brasileiras com PhD em Matemática e anos de experiência em análise de dados.

MISSÃO: Analisar os dados históricos fornecidos e criar a estratégia MAIS OTIMIZADA possível para maximizar as chances de ganhar na Mega Sena e/ou Lotofácil.

DADOS HISTÓRICOS:
%s

PREFERÊNCIAS DO USUÁRIO:
%s

REGRAS DAS LOTERIAS:
%s

INSTRUÇÕES ESPECÍFICAS:
1. 🎯 ANÁLISE ESTATÍSTICA PROFUNDA:
   - Calcule frequências de cada número nos últimos sorteios
   - Identifique padrões temporais e tendências
   - Analise distribuição de somas e padrões de espaçamento
   - Detecte números "quentes" (mais sorteados) e "frios" (menos sorteados)
   - Analise correlações entre números

2. 🧠 ESTRATÉGIAS AVANÇADAS:
   - Use teoria de probabilidades para otimizar seleção
   - Evite padrões óbvios (sequências, múltiplos, etc.)
   - Distribua números de forma equilibrada pelo range
   - Considere estratégias de cobertura máxima
   - Otimize para o orçamento disponível

3. 💰 OTIMIZAÇÃO DE ORÇAMENTO:
   - Distribua o orçamento de forma inteligente entre os jogos
   - Priorize jogos com melhor custo-benefício
   - Considere jogos com mais números para aumentar chances
   - Balance entre quantidade de jogos e cobertura

4. 📊 RESPOSTA OBRIGATÓRIA EM JSON:
{
  "strategy": {
    "games": [
      {
        "type": "megasena",
        "numbers": [1, 2, 3, 4, 5, 6],
        "cost": 5.0,
        "expectedReturn": 0.0001,
        "probability": 0.000002
      }
    ],
    "totalCost": 50.0,
    "budget": 50.0,
    "expectedReturn": 0.001,
    "reasoning": "Explicação detalhada da estratégia...",
    "statistics": {
      "totalDraws": 2000,
      "analyzedDraws": 100,
      "numberFrequency": {},
      "hotNumbers": [7, 10, 23],
      "coldNumbers": [13, 32, 55],
      "patterns": {}
    },
    "createdAt": "2025-01-27T10:00:00Z"
  },
  "confidence": 0.85,
  "alternatives": [],
  "warnings": []
}

5. ⚡ REQUIREMENTS CRÍTICOS:
   - SEMPRE retorne JSON válido
   - Gere jogos ÚNICOS (sem repetição)
   - Respeite o orçamento EXATAMENTE
   - Números válidos para cada loteria
   - Justifique CADA decisão estatisticamente
   - Seja AGRESSIVO na otimização - o objetivo é GANHAR!

🎲 AGORA ANALISE E CRIE A ESTRATÉGIA MAIS FODA POSSÍVEL!`, drawsJSON, prefsJSON, rulesJSON)

	return prompt
}

// TestConnection testa conectividade com Claude
func (c *ClaudeClient) TestConnection() error {
	testReq := ClaudeRequest{
		Model:     c.model,
		MaxTokens: 10,
		Messages: []Message{
			{
				Role:    "user",
				Content: "Teste de conectividade. Responda apenas: OK",
			},
		},
	}
	
	reqBody, err := json.Marshal(testReq)
	if err != nil {
		return fmt.Errorf("erro ao serializar teste: %w", err)
	}
	
	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição teste: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("erro na requisição teste: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Claude API retornou status %d", resp.StatusCode)
	}
	
	return nil
} 