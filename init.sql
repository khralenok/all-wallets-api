CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  user_pwd TEXT NOT NULL,
  base_currency TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_deleted BOOLEAN DEFAULT FALSE,
  deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE wallets (
  id SERIAL PRIMARY KEY,
  wallet_name TEXT NOT NULL,
  currency TEXT NOT NULL,
  balance INT NOT NULL DEFAULT 0,
  last_snapshot TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE wallet_users (
  wallet_id INT REFERENCES wallets(id),
  user_id INT REFERENCES users(id),
  user_role TEXT CHECK (user_role IN('spectator','user','admin')),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,
  amount INT NOT NULL,
  is_deposit BOOLEAN NOT NULL,
  category TEXT NOT NULL, 
  wallet_id INT REFERENCES wallets(id),
  creator_id INT REFERENCES users(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE recurrent_payments (
  id SERIAL PRIMARY KEY,
  amount INT NOT NULL,
  is_deposit BOOLEAN NOT NULL,
  category TEXT NOT NULL, 
  wallet_id INT REFERENCES wallets(id),
  frequency TEXT CHECK (frequency IN ('daily', 'weekly', 'monthly','yearly')),
  scheduled_day INT, -- for monthly/yearly
  scheduled_weekday INT, -- for weekly
  scheduled_month INT, -- for yearly
  next_run TIMESTAMP,
  end_at TIMESTAMP, -- null if infinite
  creator_id INT REFERENCES users(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_next_recurrent_payments ON recurrent_payments(next_run);