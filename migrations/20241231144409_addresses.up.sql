CREATE TABLE IF NOT EXISTS addresses(
                                        id SERIAL PRIMARY KEY UNIQUE ,
                                        private_key VARCHAR(200) UNIQUE,
                                        address VARCHAR(42) UNIQUE,
                                        currency VARCHAR(20),
                                        nonce VARCHAR(100) UNIQUE
);