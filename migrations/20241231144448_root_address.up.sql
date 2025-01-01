CREATE TABLE IF NOT EXISTS root_address(
                                           id SERIAL PRIMARY KEY UNIQUE ,
                                           private_key VARCHAR(200) UNIQUE,
                                           address VARCHAR(42) UNIQUE,
                                           currency VARCHAR(20) UNIQUE,
                                           nonce VARCHAR(100) UNIQUE
);