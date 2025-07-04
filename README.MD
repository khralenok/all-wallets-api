# AllWallet API

With AllWallets API users might create and share wallets and track their finances, which make it fit great for family or small team finance tracking

## API Endpoints

### User management
- <code>POST /signin</code> - Create user
- <code>POST /login</code> - Return JWT token
- <code>GET /profile</code> - Return such user data as ID, Username, Base Currency, Balance and all wallets summary user participate in
- <code>PUT /delete-user</code> - Delete user

### Wallet management
- <code>POST /new-wallet</code> - Create wallet
- <code>GET /wallet/:id</code> - Get wallet data
- <code>/delete-wallet/:id</code> - Delete wallet

### Wallet Users management
- <code>POST /share-wallet</code> - Add user to the list of wallet users
- <code>DELETE /remove-wallet-user/wallet/:wallet_id/username/:username/</code> - Remove user from list of wallet users

### Transaction managemet
- <code>POST /add_income</code> - Add new income transaction
- <code>POST /add_expense</code> - Add new expense transaction

## Worker features
To perform admin level actions All Wallets have worker functions that runs in CLI. 

### Updating exchange rates
This function essential and should be run daily or more often to keep exchange rates updated.

### Snapshot
This function needs to reduce amount of calculations if there are many transactions in database. Recommended to run daily for each user. In current implementation runs manually

## In Next Releases
- User roles managent
- Implementing user permissions logic
- Chanhgin user base currency logic
- ...


## Join my Telegram and Threads for Updates

- [Telegram Channel](https://t.me/khralenok_com)
- [Threads Channel](https://www.threads.com/@khralenok.studio)