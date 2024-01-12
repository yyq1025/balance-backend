INSERT INTO `networks` (`chain_id`, `name`, `url`, `symbol`, `explorer`) VALUES
('0x1', 'Ethereum', 'https://eth.public-rpc.com', 'ETH', 'https://etherscan.io'),
('0x38', 'BSC', 'https://bsc-dataseed.binance.org/', 'BNB', 'https://bscscan.com');

INSERT INTO `wallets` (`user_id`, `address`, `network_name`, `token`) VALUES
('1', decode('0000000000000000000000000000000000000000', 'hex'), 'Ethereum', decode('0000000000000000000000000000000000000000', 'hex'));
