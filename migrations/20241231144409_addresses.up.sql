CREATE TABLE IF NOT EXISTS addresses(
                                        id SERIAL PRIMARY KEY UNIQUE ,
                                        private_key VARCHAR(100) UNIQUE,
                                        address VARCHAR(42) UNIQUE,
                                        Currency VARCHAR(20)
);