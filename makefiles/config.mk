REGISTRY := 192.168.56.10:5000
TAG      := dev

GENERATOR_DIR := services/generator
INGESTION_DIR := services/ingestion
INDICATOR_DIR := services/indicator-calculator
SIGNAL_DIR    := services/signal-decision
BFF_DIR       := services/bff
FRONTEND_DIR  := services/frontend

GENERATOR_IMAGE := $(REGISTRY)/trading-generator:$(TAG)
INGESTION_IMAGE := $(REGISTRY)/trading-ingestion:$(TAG)
INDICATOR_IMAGE := $(REGISTRY)/trading-indicator-calculator:$(TAG)
SIGNAL_IMAGE    := $(REGISTRY)/trading-signal-decision:$(TAG)
BFF_IMAGE       := $(REGISTRY)/trading-bff:$(TAG)

PLATFORM_NAMESPACE := platform
PLATFORM_RELEASE   := platform
PLATFORM_CHART     := charts/platform

APP_NAMESPACE := default
APP_RELEASE   := app
APP_CHART     := charts/app

GO_MODULES := services/ingestion services/signal-decision services/bff gen/go
JAVA_MODULE := services/indicator-calculator
JS_DIR := $(FRONTEND_DIR)
JS_NPM := npm

VAGRANT_DIR := infra
VAGRANT     := cd $(VAGRANT_DIR) && vagrant