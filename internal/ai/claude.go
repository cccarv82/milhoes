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
	
	// Implementar retry logic com exponential backoff
	var resp *http.Response
	maxRetries := 3
	baseDelay := 2 * time.Second
	
	for attempt := 0; attempt < maxRetries; attempt++ {
		req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, fmt.Errorf("erro ao criar requisição: %w", err)
		}
		
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", c.apiKey)
		req.Header.Set("anthropic-version", "2023-06-01")
		
		resp, err = c.httpClient.Do(req)
		if err != nil {
			if attempt < maxRetries-1 {
				delay := baseDelay * time.Duration(1<<attempt) // Exponential backoff
				if config.IsVerbose() {
					fmt.Printf("⚠️  Tentativa %d falhou, tentando novamente em %v...\n", attempt+1, delay)
				}
				time.Sleep(delay)
				continue
			}
			return nil, fmt.Errorf("erro na requisição após %d tentativas: %w", maxRetries, err)
		}
		break
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
	// Serializar apenas os dados essenciais
	drawsJSON, _ := json.MarshalIndent(request.Draws, "", "  ")
	prefsJSON, _ := json.MarshalIndent(request.Preferences, "", "  ")
	rulesJSON, _ := json.MarshalIndent(request.Rules, "", "  ")
	
	prompt := fmt.Sprintf(`Você é um especialista em análise estatística de loterias brasileiras.

MISSÃO: Analisar dados históricos e criar estratégia otimizada para maximizar chances de ganhar.

DADOS HISTÓRICOS (%d sorteios):
%s

PREFERÊNCIAS:
%s

REGRAS:
%s

INSTRUÇÕES:
1. 🎯 ANÁLISE: Calcule frequências, identifique padrões, números quentes/frios
2. 🧠 ESTRATÉGIA: Use probabilidade, evite padrões óbvios, distribua números
3. 💰 ORÇAMENTO: Otimize distribuição, respeite limite exato
4. 📊 RESPOSTA EM JSON:
{
  "strategy": {
    "games": [{"type": "megasena", "numbers": [1,2,3,4,5,6], "cost": 5.0, "probability": 0.000002}],
    "totalCost": 50.0,
    "budget": 50.0,
    "reasoning": "Explicação da estratégia...",
    "statistics": {"totalDraws": %d, "hotNumbers": [], "coldNumbers": []},
    "createdAt": "2025-01-27T10:00:00Z"
  },
  "confidence": 0.85
}

REQUIREMENTS: JSON válido, jogos únicos, orçamento exato, números válidos, justificativa estatística.`, 
		len(request.Draws), drawsJSON, prefsJSON, rulesJSON, len(request.Draws))

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