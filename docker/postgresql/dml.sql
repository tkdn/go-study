SET search_path = public;

INSERT INTO users(name, age) VALUES('Mike', 10);
INSERT INTO users(name, age) VALUES('Jane', 6);
INSERT INTO users(name, age) VALUES('George', 13);

INSERT INTO posts(user_id, text) VALUES(1, 'foo bar post...チョコレート1');
INSERT INTO posts(user_id, text) VALUES(2, 'foo bar post...チョコレート2');
INSERT INTO posts(user_id, text) VALUES(1, 'foo bar post...チョコレート3');
