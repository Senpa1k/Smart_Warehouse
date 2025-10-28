package config

func Get(key string) (string, error) {
	return "postgresql://warehouse_user:secure_password@localhost:5432/warehouse_db?sslmode=disable", nil
	// if val := os.Getenv(key); val != "" {
	// 	return val, nil
	// }
	// return "", fmt.Errorf("have not acces to env ")

}
