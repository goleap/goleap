-- MySQL dump 10.13  Distrib 8.0.32, for Linux (aarch64)
--
-- Host: localhost    Database: acceptance
-- ------------------------------------------------------
-- Server version	8.0.32

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `comments`
--

DROP TABLE IF EXISTS `comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `comments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `post_id` int NOT NULL,
  `parent_id` int DEFAULT NULL,
  `content` text NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `post_id` (`post_id`),
  KEY `parent_id` (`parent_id`),
  CONSTRAINT `comments_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `comments_ibfk_2` FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`),
  CONSTRAINT `comments_ibfk_3` FOREIGN KEY (`parent_id`) REFERENCES `comments` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `comments`
--

LOCK TABLES `comments` WRITE;
/*!40000 ALTER TABLE `comments` DISABLE KEYS */;
INSERT INTO `comments` VALUES (1,2,1,NULL,'Je pratique la méditation tous les jours et cela a vraiment amélioré ma vie. Je recommande également cette pratique à tout le monde !','2023-03-28 02:38:19','2023-03-28 02:38:19'),(2,3,1,NULL,'Je suis entièrement d\'accord. La méditation est une excellente façon de se recentrer et de réduire le stress.','2023-03-29 02:38:19','2023-03-29 02:38:19'),(3,1,2,NULL,'Je suis jaloux, j\'aimerais voyager en Asie un jour ! Quel était ton endroit préféré que tu as visité ?','2023-03-30 02:38:19','2023-03-30 02:38:19'),(4,2,2,NULL,'C\'était difficile de choisir un endroit préféré, mais je pense que j\'ai adoré Tokyo au Japon. La nourriture était incroyable et il y avait tellement de choses à voir et à faire !','2023-03-28 02:38:19','2023-03-28 02:38:19'),(5,3,3,NULL,'Je suis en train d\'apprendre le français en ce moment, et je trouve que regarder des films et des émissions de télévision en français est vraiment utile pour pratiquer la langue.','2023-03-28 02:38:19','2023-03-28 02:38:19'),(6,1,4,NULL,'Je vais essayer cette recette ce week-end, merci pour le partage !','2023-03-28 02:38:19','2023-03-28 02:38:19'),(7,2,4,6,'J\'ai également essayé cette recette et elle est délicieuse. J\'ai ajouté des pommes de terre et des carottes pour compléter le plat.','2023-04-10 18:11:40','2023-04-10 18:11:40'),(8,3,4,7,'Je suis content de savoir que je ne suis pas le seul à aimer cette recette ! J\'adore aussi ajouter des légumes pour en faire un plat complet.','2023-04-10 18:11:49','2023-04-10 18:11:49');
/*!40000 ALTER TABLE `comments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `posts`
--

DROP TABLE IF EXISTS `posts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `posts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `title` varchar(255) NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `posts_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `posts`
--

LOCK TABLES `posts` WRITE;
/*!40000 ALTER TABLE `posts` DISABLE KEYS */;
INSERT INTO `posts` VALUES (1,1,'Les bienfaits de la méditation','La méditation est une pratique ancienne qui peut améliorer notre bien-être mental et physique. Des études ont montré qu\'elle peut réduire le stress, améliorer la concentration et même réduire la douleur chronique. Si vous cherchez à améliorer votre qualité de vie, essayez la méditation !','2023-03-28 02:38:05','2023-03-28 02:38:05'),(2,2,'Mon voyage en Asie','Je suis récemment allé en Asie et j\'ai eu l\'opportunité de découvrir de nouvelles cultures et de nouveaux paysages incroyables. J\'ai visité le Japon, la Chine et la Thaïlande, et chaque endroit avait sa propre beauté unique. Je recommande fortement à tout le monde de voyager en Asie au moins une fois dans leur vie.','2023-03-28 02:38:05','2023-03-28 02:38:05'),(3,3,'Comment apprendre une nouvelle langue rapidement','Apprendre une nouvelle langue peut sembler intimidant, mais avec la bonne approche, c\'est plus facile que vous ne le pensez. Voici quelques conseils pour apprendre une langue rapidement : pratiquez régulièrement, écoutez de la musique dans la langue que vous apprenez, regardez des films et des séries en version originale, et trouvez un tuteur ou un ami natif pour pratiquer avec vous.','2023-03-28 02:38:05','2023-03-28 02:38:05'),(4,1,'Ma recette préférée : poulet rôti aux herbes','Le poulet rôti aux herbes est l\'un de mes plats préférés. C\'est facile à faire et ça a tellement de saveur ! Voici ma recette : mélangez du romarin, du thym, de l\'ail, du sel et du poivre dans un bol. Ensuite, frottez le mélange sur un poulet entier et placez-le dans un plat allant au four. Faites cuire le poulet à 180°C pendant environ une heure et demie, ou jusqu\'à ce que le jus qui s\'écoule du poulet soit clair. Servez avec des légumes grillés pour un repas délicieux et sain.','2023-03-28 02:38:05','2023-03-28 02:38:05');
/*!40000 ALTER TABLE `posts` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `validated` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'john.doe@example.com','$2a$10$bgBxj3vO5WpU8B6U5m5U5Ov2QW40R8G6eH7UnhJXXvSPe63YAYMmi',1,'2023-03-24 02:31:25','2023-03-24 02:31:25'),(2,'jane.smith@example.com','$2a$10$t4Z/lLefwozbbIMFhJ8SOOwDeNlNKEO19T0GK/j9XspR0d0YDguhe',1,'2023-03-24 02:37:45','2023-03-24 02:37:45'),(3,'emily.davis@example.com','$2a$10$7.Pa/lqne5X7zZJK5Y2Q7e/MJW/kxTw7vR0eM5SPeh.sMfK8WxUyO',1,'2023-03-24 02:38:01','2023-03-24 02:38:01');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-04-10 22:40:22
