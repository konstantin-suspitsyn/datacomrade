package jsonlog

// Logger custom presaved responses are here
// Uses slog Default logger. It should be created separately

import (
	"context"
	"log/slog"
	"net/http"
)

func SimpleLog(logLevel slog.Level, message string, propertiesNoGroup map[string]string, propertiesGroup map[string]map[string]string) {

	attrs := []slog.Attr{}

	if propertiesNoGroup != nil {
		attrs = propertiesNoGropToAttrs(attrs, propertiesNoGroup)
	}

	if propertiesGroup != nil {
		attrs = propertiesGroupToAttrs(attrs, propertiesGroup)
	}

	slog.LogAttrs(
		context.Background(),
		logLevel,
		message,

		attrs...,
	)
}

// Helper function
func propertiesNoGropToAttrs(attrs []slog.Attr, propertiesNoGroup map[string]string) []slog.Attr {
	for attribute, property := range propertiesNoGroup {
		attrs = append(attrs, slog.String(attribute, property))
	}
	return attrs
}

// Helper function that converts map to slog.Attr slice as slog.Group for LogAttrs function
func propertiesGroupToAttrs(attrs []slog.Attr, propertiesGroup map[string]map[string]string) []slog.Attr {
	for groupName, mapOfAttributes := range propertiesGroup {
		groupAttributes := []any{}
		for key, val := range mapOfAttributes {
			groupAttributes = append(groupAttributes, key)
			groupAttributes = append(groupAttributes, val)
		}
		attrs = append(attrs, slog.Group(groupName, groupAttributes...))
	}
	return attrs
}

// Logs Error for HTTP Response without any cusom properties
func LogResponseErrorNoCustomProperties(r *http.Request, err error) {
	properties := map[string]string{
		"request_url":    r.URL.String(),
		"request_method": r.Method,
	}

	groupProrerties := map[string]map[string]string{
		"request": properties,
	}

	SimpleLog(slog.LevelError, err.Error(), nil, groupProrerties)
}

func PrintInfo(message string, propertiesNoGroup map[string]string, propertiesGroup map[string]map[string]string) {
	SimpleLog(slog.LevelInfo, message, propertiesNoGroup, propertiesGroup)
}

func PrintError(message string, propertiesNoGroup map[string]string, propertiesGroup map[string]map[string]string) {
	SimpleLog(slog.LevelError, message, propertiesNoGroup, propertiesGroup)
}
