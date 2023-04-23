CREATE TABLE `question` (
  `id` int NOT NULL AUTO_INCREMENT,
  `question_text` varchar(255) NOT NULL,
  `choice_a` varchar(255) NOT NULL,
  `choice_b` varchar(255) NOT NULL,
  `choice_c` varchar(255) NOT NULL,
  `choice_d` varchar(255) NOT NULL,
  `choice_e` varchar(255) NOT NULL,
  `choice_f` varchar(255) NOT NULL,
  `correct_answer` int NOT NULL,
  `status` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'ready',
  `recent_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)

CREATE TABLE `participant` (
  `code` varchar(255) NOT NULL,
  `name` varchar(255),
  `phone_number` varchar(255),
  `question_id` int,
  `answer` int,
  `registered_time` timestamp,
  PRIMARY KEY (`code`),
  KEY `question_id` (`question_id`),
  CONSTRAINT `participant_ibfk_1` FOREIGN KEY (`question_id`) REFERENCES `question` (`id`)
)

CREATE TABLE game_quota (
	quota int
);

INSERT INTO game_quota (quota) VALUES (10);