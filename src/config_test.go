package main

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	const testKey = "TEST_GET_ENV_VAR"

	// Убеждаемся, что переменная не установлена
	os.Unsetenv(testKey)

	// Тест дефолтного значения (переменная отсутствует)
	if val := getEnv(testKey, "default_val"); val != "default_val" {
		t.Errorf("getEnv() = %s, ожидалось %s", val, "default_val")
	}

	// Тест дефолтного значения (переменная пуста)
	os.Setenv(testKey, "")
	if val := getEnv(testKey, "default_val"); val != "default_val" {
		t.Errorf("getEnv() = %s, ожидалось %s", val, "default_val")
	}

	// Тест установленной переменной
	os.Setenv(testKey, "custom_val")
	defer os.Unsetenv(testKey)

	if val := getEnv(testKey, "default_val"); val != "custom_val" {
		t.Errorf("getEnv() = %s, ожидалось %s", val, "custom_val")
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	// Очищаем переменные окружения перед тестом
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("UPLOAD_DIR")

	cfg := LoadConfig()

	expectedDB := "postgres://petya_db:arelun06_db@localhost:5666/Data_cloud"
	if cfg.DatabaseURL != expectedDB {
		t.Errorf("LoadConfig().DatabaseURL = %s, ожидалось %s", cfg.DatabaseURL, expectedDB)
	}

	expectedPort := "8080"
	if cfg.ServerPort != expectedPort {
		t.Errorf("LoadConfig().ServerPort = %s, ожидалось %s", cfg.ServerPort, expectedPort)
	}

	expectedUploadDir := "./uploading"
	if cfg.UploadDir != expectedUploadDir {
		t.Errorf("LoadConfig().UploadDir = %s, ожидалось %s", cfg.UploadDir, expectedUploadDir)
	}
}

func TestLoadConfig_Custom(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://custom:pass@customhost:5432/customdb")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("UPLOAD_DIR", "/custom/upload")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("UPLOAD_DIR")
	}()

	cfg := LoadConfig()

	if cfg.DatabaseURL != "postgres://custom:pass@customhost:5432/customdb" {
		t.Errorf("LoadConfig().DatabaseURL = %s, ожидалось postgres://custom:pass@customhost:5432/customdb", cfg.DatabaseURL)
	}
	if cfg.ServerPort != "9090" {
		t.Errorf("LoadConfig().ServerPort = %s, ожидалось 9090", cfg.ServerPort)
	}
	if cfg.UploadDir != "/custom/upload" {
		t.Errorf("LoadConfig().UploadDir = %s, ожидалось /custom/upload", cfg.UploadDir)
	}
}
