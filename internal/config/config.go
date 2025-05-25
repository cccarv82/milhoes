package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config estrutura de configuração da aplicação
type Config struct {
	Claude ClaudeConfig `yaml:"claude"`
	App    AppConfig    `yaml:"app"`
}

// ClaudeConfig configurações da API do Claude
type ClaudeConfig struct {
	APIKey     string `yaml:"api_key"`
	Model      string `yaml:"model"`
	MaxTokens  int    `yaml:"max_tokens"`
	TimeoutSec int    `yaml:"timeout_sec"`
}

// AppConfig configurações da aplicação
type AppConfig struct {
	CacheEnabled   bool   `yaml:"cache_enabled"`
	CacheDuration  int    `yaml:"cache_duration_hours"`
	DefaultBudget  int    `yaml:"default_budget"`
	LogLevel       string `yaml:"log_level"`
	DataSourceURL  string `yaml:"data_source_url"`
}

var GlobalConfig *Config

// Init inicializa a configuração global
func Init() {
	GlobalConfig = &Config{
		Claude: ClaudeConfig{
			APIKey:     getClaudeAPIKey(),
			Model:      viper.GetString("claude.model"),
			MaxTokens:  viper.GetInt("claude.max_tokens"),
			TimeoutSec: viper.GetInt("claude.timeout_sec"),
		},
		App: AppConfig{
			CacheEnabled:   viper.GetBool("app.cache_enabled"),
			CacheDuration:  viper.GetInt("app.cache_duration_hours"),
			DefaultBudget:  viper.GetInt("app.default_budget"),
			LogLevel:       viper.GetString("app.log_level"),
			DataSourceURL:  viper.GetString("app.data_source_url"),
		},
	}

	// Configurações padrão
	setDefaults()
}

// getClaudeAPIKey obtém a chave da API do Claude
func getClaudeAPIKey() string {
	// Prioridade: flag -> env var -> config file -> default
	if key := viper.GetString("api-key"); key != "" {
		return key
	}
	
	if key := os.Getenv("CLAUDE_API_KEY"); key != "" {
		return key
	}
	
	if key := viper.GetString("claude.api_key"); key != "" {
		return key
	}
	
	// Se não encontrar em nenhum lugar, usar a key padrão fornecida pelo usuário
	return "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"
}

// setDefaults define valores padrão para configurações
func setDefaults() {
	if GlobalConfig.Claude.Model == "" {
		GlobalConfig.Claude.Model = "claude-3-5-sonnet-20241022"
	}
	
	if GlobalConfig.Claude.MaxTokens == 0 {
		GlobalConfig.Claude.MaxTokens = 4000
	}
	
	if GlobalConfig.Claude.TimeoutSec == 0 {
		GlobalConfig.Claude.TimeoutSec = 30
	}
	
	if GlobalConfig.App.CacheDuration == 0 {
		GlobalConfig.App.CacheDuration = 24 // 24 horas
	}
	
	if GlobalConfig.App.DefaultBudget == 0 {
		GlobalConfig.App.DefaultBudget = 50 // R$ 50
	}
	
	if GlobalConfig.App.LogLevel == "" {
		GlobalConfig.App.LogLevel = "info"
	}
	
	if GlobalConfig.App.DataSourceURL == "" {
		GlobalConfig.App.DataSourceURL = "https://servicebus2.caixa.gov.br/portaldeloterias/api"
	}
	
	GlobalConfig.App.CacheEnabled = true
}

// ValidateConfig valida se a configuração está correta
func ValidateConfig() error {
	if GlobalConfig.Claude.APIKey == "" {
		return fmt.Errorf("chave da API do Claude não configurada. Use --api-key ou defina CLAUDE_API_KEY")
	}
	
	if GlobalConfig.App.DefaultBudget <= 0 {
		return fmt.Errorf("orçamento padrão deve ser maior que zero")
	}
	
	return nil
}

// GetClaudeAPIKey retorna a chave da API do Claude
func GetClaudeAPIKey() string {
	return GlobalConfig.Claude.APIKey
}

// GetClaudeModel retorna o modelo do Claude a ser usado
func GetClaudeModel() string {
	return GlobalConfig.Claude.Model
}

// GetMaxTokens retorna o número máximo de tokens para o Claude
func GetMaxTokens() int {
	return GlobalConfig.Claude.MaxTokens
}

// IsVerbose retorna se o modo verbose está ativo
func IsVerbose() bool {
	return viper.GetBool("verbose")
} 