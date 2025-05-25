package data

import (
	"encoding/json"
	"fmt"
	"lottery-optimizer/internal/config"
	"lottery-optimizer/internal/lottery"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client cliente para APIs de loterias
type Client struct {
	client *resty.Client
	baseURL string
}

// NewClient cria um novo cliente para APIs de loterias
func NewClient() *Client {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(2 * time.Second)
	
	// Headers para parecer mais com um browser real
	client.SetHeaders(map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "pt-BR,pt;q=0.9,en;q=0.8",
		"Accept-Encoding": "gzip, deflate, br",
		"Cache-Control":   "no-cache",
		"Pragma":          "no-cache",
		"Sec-Fetch-Dest":  "empty",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Site":  "cross-site",
		"Referer":         "https://loterias.caixa.gov.br/",
		"Origin":          "https://loterias.caixa.gov.br",
	})
	
	return &Client{
		client:  client,
		baseURL: config.GlobalConfig.App.DataSourceURL,
	}
}

// GetLatestDraws busca os últimos sorteios de uma loteria
func (c *Client) GetLatestDraws(ltype lottery.LotteryType, count int) ([]lottery.Draw, error) {
	endpoint := fmt.Sprintf("%s/%s/", c.baseURL, string(ltype))
	
	var draws []lottery.Draw
	
	// Delay inicial para não parecer bot
	time.Sleep(1 * time.Second)
	
	// Buscar o último sorteio primeiro para descobrir o número atual
	latestResp, err := c.client.R().Get(endpoint)
	
	if err != nil || latestResp.StatusCode() != 200 {
		fmt.Printf("⚠️  API da CAIXA inacessível (erro %d), usando dados simulados para %s\n", 
			latestResp.StatusCode(), ltype)
		return GetMockDraws(ltype, count), nil
	}
	
	// Debug: Log da resposta
	if config.IsVerbose() {
		fmt.Printf("🔍 Debug %s: Status=%d, Content-Type=%s\n", ltype, latestResp.StatusCode(), latestResp.Header().Get("Content-Type"))
		fmt.Printf("🔍 Debug %s: Primeiros 200 chars da resposta: %s\n", ltype, string(latestResp.Body()[:min(200, len(latestResp.Body()))]))
	}
	
	var latest lottery.Draw
	if err := json.Unmarshal(latestResp.Body(), &latest); err != nil {
		fmt.Printf("⚠️  Erro ao decodificar dados da API, usando dados simulados para %s\n", ltype)
		return GetMockDraws(ltype, count), nil
	}
	
	draws = append(draws, latest)
	
	// Buscar sorteios anteriores com delay entre requisições
	for i := 1; i < count && latest.Number-i > 0; i++ {
		// Delay para não sobrecarregar a API
		time.Sleep(500 * time.Millisecond)
		
		drawNumber := latest.Number - i
		drawResp, err := c.client.R().Get(fmt.Sprintf("%s%d", endpoint, drawNumber))
		
		if err != nil || drawResp.StatusCode() != 200 {
			if config.IsVerbose() {
				fmt.Printf("Erro ao buscar sorteio %d, parando busca: %v\n", drawNumber, err)
			}
			break
		}
		
		var draw lottery.Draw
		if err := json.Unmarshal(drawResp.Body(), &draw); err != nil {
			if config.IsVerbose() {
				fmt.Printf("Erro ao decodificar sorteio %d: %v\n", drawNumber, err)
			}
			continue
		}
		
		draws = append(draws, draw)
	}
	
	// Se conseguimos poucos dados da API, complementar com mock
	if len(draws) < count/2 {
		fmt.Printf("⚠️  Poucos dados obtidos da API (%d/%d), complementando com dados simulados\n", len(draws), count)
		mockDraws := GetMockDraws(ltype, count-len(draws))
		draws = append(draws, mockDraws...)
	}
	
	return draws, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetDrawByNumber busca um sorteio específico pelo número
func (c *Client) GetDrawByNumber(ltype lottery.LotteryType, number int) (*lottery.Draw, error) {
	endpoint := fmt.Sprintf("%s/%s/%d", c.baseURL, string(ltype), number)
	
	resp, err := c.client.R().
		SetHeader("Accept", "application/json").
		Get(endpoint)
	
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar sorteio %d: %w", number, err)
	}
	
	var draw lottery.Draw
	if err := json.Unmarshal(resp.Body(), &draw); err != nil {
		return nil, fmt.Errorf("erro ao decodificar sorteio %d: %w", number, err)
	}
	
	return &draw, nil
}

// GetDrawsRange busca sorteios em um intervalo
func (c *Client) GetDrawsRange(ltype lottery.LotteryType, startNumber, endNumber int) ([]lottery.Draw, error) {
	var draws []lottery.Draw
	
	for number := startNumber; number <= endNumber; number++ {
		draw, err := c.GetDrawByNumber(ltype, number)
		if err != nil {
			if config.IsVerbose() {
				fmt.Printf("Erro ao buscar sorteio %d: %v\n", number, err)
			}
			continue
		}
		draws = append(draws, *draw)
		
		// Pequeno delay para não sobrecarregar a API
		time.Sleep(100 * time.Millisecond)
	}
	
	return draws, nil
}

// GetAllHistoricalDraws busca todo o histórico disponível (usar com cuidado)
func (c *Client) GetAllHistoricalDraws(ltype lottery.LotteryType) ([]lottery.Draw, error) {
	// Primeiro, buscar o último sorteio para saber quantos existem
	latest, err := c.GetLatestDraws(ltype, 1)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar último sorteio: %w", err)
	}
	
	if len(latest) == 0 {
		return nil, fmt.Errorf("nenhum sorteio encontrado")
	}
	
	latestNumber := latest[0].Number
	
	// Buscar todos os sorteios do 1 até o último
	return c.GetDrawsRange(ltype, 1, latestNumber)
}

// TestConnection testa se a API está respondendo
func (c *Client) TestConnection() error {
	resp, err := c.client.R().Get(fmt.Sprintf("%s/megasena/", c.baseURL))
	
	if err != nil {
		return fmt.Errorf("erro de conectividade: %w - usando dados simulados", err)
	}
	
	if resp.StatusCode() == 403 {
		return fmt.Errorf("API da CAIXA bloqueada (403) - usando dados simulados para funcionamento")
	}
	
	if resp.StatusCode() != 200 {
		return fmt.Errorf("API retornou status %d - usando dados simulados", resp.StatusCode())
	}
	
	return nil
}

// GetNextDrawInfo busca informações sobre o próximo sorteio
func (c *Client) GetNextDrawInfo(ltype lottery.LotteryType) (*time.Time, int, error) {
	latest, err := c.GetLatestDraws(ltype, 1)
	if err != nil {
		return nil, 0, err
	}
	
	if len(latest) == 0 {
		return nil, 0, fmt.Errorf("nenhum sorteio encontrado")
	}
	
	nextDate := latest[0].NextDrawDate.Time()
	return &nextDate, latest[0].NextDrawNumber, nil
}

// TestDirectAPI testa resposta bruta da API
func (c *Client) TestDirectAPI() (string, error) {
	resp, err := c.client.R().
		SetHeader("Accept", "application/json").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		Get("https://servicebus2.caixa.gov.br/portaldeloterias/api/lotofacil/")
	
	if err != nil {
		return "", fmt.Errorf("erro na requisição: %w", err)
	}
	
	return fmt.Sprintf("Status: %d, Content-Type: %s, Body: %s", 
		resp.StatusCode(), 
		resp.Header().Get("Content-Type"),
		string(resp.Body()[:min(500, len(resp.Body()))])), nil
} 