package common

type DBConfig struct {
    Engine          string `json:"engine"`
    DSN             string `json:"DSN"`
    MaxOpenConns    int    `json:"maxOpenConns"`
    MaxIdleConns    int    `json:"maxIdleConns"`
    ConnMaxLifetime int64  `json:"connMaxLifetime"`
}
