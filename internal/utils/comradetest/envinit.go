package comradetest

import (
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func InitEnv() {
	_, dir, _, _ := runtime.Caller(0)

	env_dir := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(dir))))
	env_file := filepath.Join(env_dir, ".env-test")

	err := godotenv.Load(env_file)

	if err != nil {
		panic(err.Error())
	}

}
