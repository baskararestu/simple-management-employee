# ************************************************************
# Antares - SQL Client
# Version 0.7.34
# 
# https://antares-sql.app/
# https://github.com/antares-sql/antares
# 
# Host: 127.0.0.1 (Ubuntu 24.04 10.11.11)
# Database: db_simple_management
# Generation time: 2025-05-14T17:42:14+07:00
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table annual_leaves
# ------------------------------------------------------------

DROP TABLE IF EXISTS `annual_leaves`;

CREATE TABLE `annual_leaves` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `start_date` datetime(3) NOT NULL,
  `end_date` datetime(3) NOT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `status` varchar(20) NOT NULL DEFAULT 'pending',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_annual_leaves_user` (`user_id`),
  CONSTRAINT `fk_annual_leaves_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

LOCK TABLES `annual_leaves` WRITE;
/*!40000 ALTER TABLE `annual_leaves` DISABLE KEYS */;

INSERT INTO `annual_leaves` (`id`, `user_id`, `start_date`, `end_date`, `reason`, `status`, `created_at`, `updated_at`) VALUES
	("243e19c5-51a8-44bc-9718-3fae2e886395", "4cacbb80-da98-4b58-9508-c36008d25b5d", "2025-04-01 00:00:00.000", "2025-04-02 00:00:00.000", "test", "pending", "2025-05-14 10:01:03.569", "2025-05-14 10:01:03.569"),
	("c23258ea-d97d-4ba8-9656-2ac51798a469", "c3db7cd6-e95e-4774-b66d-093dae9adb22", "0000-00-00 00:00:00.000", "0000-00-00 00:00:00.000", "test", "approved", "0000-00-00 00:00:00.000", "2025-05-14 10:27:14.961"),
	("c249228b-9c5b-4534-afa2-ea2102950316", "4cacbb80-da98-4b58-9508-c36008d25b5d", "2025-04-01 00:00:00.000", "2025-04-02 00:00:00.000", "test", "rejected", "2025-05-14 09:56:51.752", "2025-05-14 09:56:51.752");

/*!40000 ALTER TABLE `annual_leaves` ENABLE KEYS */;
UNLOCK TABLES;



# Dump of table roles
# ------------------------------------------------------------

DROP TABLE IF EXISTS `roles`;

CREATE TABLE `roles` (
  `id` varchar(36) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

LOCK TABLES `roles` WRITE;
/*!40000 ALTER TABLE `roles` DISABLE KEYS */;

INSERT INTO `roles` (`id`, `name`, `created_at`, `updated_at`) VALUES
	("23d504f7-abc9-4356-a5ba-92f491545265", "admin", "2025-05-14 08:51:55.662", "2025-05-14 08:51:55.662"),
	("5288949c-4271-4617-aa55-1bd070d5aebd", "employee", "2025-05-14 08:51:55.669", "2025-05-14 08:51:55.669");

/*!40000 ALTER TABLE `roles` ENABLE KEYS */;
UNLOCK TABLES;



# Dump of table users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` varchar(36) NOT NULL,
  `first_name` varchar(255) DEFAULT NULL,
  `last_name` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `phone_number` varchar(20) DEFAULT NULL,
  `gender` enum('MALE','FEMALE') DEFAULT NULL,
  `role_id` varchar(36) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_users_email` (`email`),
  KEY `fk_users_role` (`role_id`),
  CONSTRAINT `fk_users_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;

INSERT INTO `users` (`id`, `first_name`, `last_name`, `email`, `password`, `address`, `phone_number`, `gender`, `role_id`, `created_at`, `updated_at`) VALUES
	("222b3e8d-ea01-4ce0-a2fb-df3ec9298a67", "admin", "ke4", "adminke2@example.com", "$2a$12$95O9OlMCrpiiDx9uk1EmzuxsFplkzaGxAEm3qmDlC6xh/7da0HhjG", NULL, NULL, NULL, "23d504f7-abc9-4356-a5ba-92f491545265", "0000-00-00 00:00:00.000", "2025-05-14 09:32:47.532"),
	("39d9e303-409c-43b1-bb1d-75545080613f", "John", "Doe", "johndoe@example.com", "$2a$12$9MoYz2JevYItYkGHrZNwu.XzQg4QHS27dWaNDEK4OhuVouSn4/xqy", NULL, NULL, NULL, "23d504f7-abc9-4356-a5ba-92f491545265", "2025-05-14 08:51:56.087", "2025-05-14 08:51:56.087"),
	("4cacbb80-da98-4b58-9508-c36008d25b5d", "pekerja", "ke222222222", "pekerjake22222222222@example.com", "$2a$12$d5TZLpC.GjFDpZvgOvv2LOMobjuMySzdDMQexpxxqS0RfbMYmTpY2", "Jalan ghaib", "000999888", "MALE", "5288949c-4271-4617-aa55-1bd070d5aebd", "2025-05-14 09:54:57.062", "2025-05-14 09:54:57.062"),
	("999d581b-9534-46bc-811b-3402b809446b", "admin", "ke555", "adminke7070@example.com", "$2a$12$dzI11XNAvsfAXfms5bXbau48jkdnk8tTVmkVy.mIgNEHcIIgY5L.y", NULL, NULL, NULL, "23d504f7-abc9-4356-a5ba-92f491545265", "0000-00-00 00:00:00.000", "2025-05-14 09:38:50.197"),
	("c3db7cd6-e95e-4774-b66d-093dae9adb22", "pekerja", "ke1", "pekerjaZZZ@example.com", "$2a$12$a6d0QDr6B7ZtOFTfxqGrf.LieLO1LnAlP/FbN/KDF9a5kRbE/uyd2", "Jalan sudirman", "111111", NULL, "5288949c-4271-4617-aa55-1bd070d5aebd", "0000-00-00 00:00:00.000", "2025-05-14 09:44:59.274");

/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;



# Dump of views
# ------------------------------------------------------------

# Creating temporary tables to overcome VIEW dependency errors


/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;

# Dump completed on 2025-05-14T17:42:15+07:00
