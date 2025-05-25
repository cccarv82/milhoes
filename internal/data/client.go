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
	
	return &Client{
		client:  client,
		baseURL: config.GlobalConfig.App.DataSourceURL,
	}
}

// GetLatestDraws busca os últimos sorteios de uma loteria
func (c *Client) GetLatestDraws(ltype lottery.LotteryType, count int) ([]lottery.Draw, error) {
	endpoint := fmt.Sprintf("%s/%s/", c.baseURL, string(ltype))
	
	var draws []lottery.Draw
	
	// Buscar o último sorteio primeiro para descobrir o número atual
	latestResp, err := c.client.R().
		SetHeader("Accept", "application/json").
		Get(endpoint)
	
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar último sorteio: %w", err)
	}
	
	var latest lottery.Draw
	if err := json.Unmarshal(latestResp.Body(), &latest); err != nil {
		return nil, fmt.Errorf("erro ao decodificar último sorteio: %w", err)
	}
	
	draws = append(draws, latest)
	
	// Buscar sorteios anteriores
	for i := 1; i < count && latest.Number-i > 0; i++ {
		drawNumber := latest.Number - i
		drawResp, err := c.client.R().
			SetHeader("Accept", "application/json").
			Get(fmt.Sprintf("%s%d", endpoint, drawNumber))
		
		if err != nil {
			if config.IsVerbose() {
				fmt.Printf("Erro ao buscar sorteio %d: %v\n", drawNumber, err)
			}
			continue
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
	
	return draws, nil
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
	resp, err := c.client.R().
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("%s/megasena/", c.baseURL))
	
	if err != nil {
		return fmt.Errorf("erro de conectividade: %w", err)
	}
	
	if resp.StatusCode() != 200 {
		return fmt.Errorf("API retornou status %d", resp.StatusCode())
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
	
	return &latest[0].NextDrawDate, latest[0].NextDrawNumber, nil
} 