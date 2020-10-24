package main

import "os"

func setEnvVars() {
	os.Setenv("DB_NAME", "test")
	os.Setenv("DB_USERNAME", "19xlr95")
	os.Setenv("DB_PASSWORD", "19xlr95")
	os.Setenv("DB_LOCATION", "Europe%2FIstanbul")
}
