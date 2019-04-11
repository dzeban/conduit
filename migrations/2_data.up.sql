INSERT INTO articles (title, slug, description, BODY, created, updated)
VALUES ('How to train your dragon', 'how-to-train-your-dragon', 'Ever wonder how?', 'It takes a Jacobian', now(), now());
INSERT INTO articles (title, slug, description, BODY, created, updated)
VALUES ('My article 2', 'my-article-2', 'So toothless', 'Some really awesome content', now(), now());

-- password is test
INSERT INTO users (name, email, password)
VALUES ('test', 'test@example.com', '$argon2id$v=19$m=32768,t=5,p=1$TPhFoj0HehIQvKGTVLlD/Q$9p9FRCUDXTj5STagKuombK8QWlFOq9nJwO4lq9EdD0BfHzxZNyt+ih6FYKB42MJLzth5Xy4/elmKVD4YqVC3gA')
