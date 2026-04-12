package utils

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
)

// ExecuteSQLScript executes a SQL script file against the database
func ExecuteSQLScript(db *sql.DB, scriptPath string) error {
	script, err := os.Open(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to open script file: %v", err)
	}
	defer script.Close()

	scanner := bufio.NewScanner(script)
	var statement strings.Builder
	executedStatements := 0
	skippedStatements := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "--") || strings.HasPrefix(line, "#") {
			continue
		}

		statement.WriteString(line)
		statement.WriteString("\n")

		// Check if statement ends with semicolon
		if strings.HasSuffix(strings.TrimSpace(line), ";") {
			query := strings.TrimSpace(statement.String())
			if query != "" {
				if _, err := db.Exec(query); err != nil {
					// Continue execution even if some statements fail
					// This is important for INSERT IGNORE statements
					zap.L().Warn("skipped statement execution (may already exist)", zap.Error(err), zap.String("query", query))
					skippedStatements++
				} else {
					executedStatements++
				}
			}
			statement.Reset()
		}
	}

	// Execute any remaining statement that doesn't end with semicolon
	if statement.Len() > 0 {
		query := strings.TrimSpace(statement.String())
		if query != "" {
			if _, err := db.Exec(query); err != nil {
				zap.L().Warn("skipped final statement execution (may already exist)", zap.Error(err), zap.String("query", query))
				skippedStatements++
			} else {
				executedStatements++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading script file: %v", err)
	}

	zap.L().Info("executed SQL initialization script",
		zap.Int("executed_statements", executedStatements),
		zap.Int("skipped_statements", skippedStatements),
		zap.String("script_path", scriptPath))

	return nil
}

// ExecuteInitDataScript executes the default data initialization script
func ExecuteInitDataScript(db *sql.DB) error {
	scriptPath := "scripts/init_default_data.sql"

	// Check if the script file exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		zap.L().Warn("initialization script not found, skipping...", zap.String("path", scriptPath))
		return nil
	}

	zap.L().Info("executing default data initialization script", zap.String("path", scriptPath))

	if err := ExecuteSQLScript(db, scriptPath); err != nil {
		return fmt.Errorf("failed to execute initialization script: %v", err)
	}

	return nil
}
