SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";

CREATE TABLE `connections` (
  `ID` int(11) NOT NULL,
  `Host-ID` int(11) NOT NULL,
  `Server-Port` int(6) NOT NULL,
  `last-used` datetime NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `hosts` (
  `ID` int(4) NOT NULL,
  `Name` varchar(20) NOT NULL,
  `User` varchar(20) NOT NULL,
  `Port` int(6) NOT NULL,
  `publicIP` varchar(15) DEFAULT NULL,
  `localIP` varchar(15) DEFAULT NULL,
  `online` int(1) NOT NULL,
  `last-online` datetime NOT NULL DEFAULT current_timestamp(),
  `version` int(5) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `requests` (
  `ID` int(3) NOT NULL,
  `Host-ID` int(11) NOT NULL,
  `last-used` datetime NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


ALTER TABLE `connections`
  ADD PRIMARY KEY (`ID`),
  ADD KEY `Host-ID` (`Host-ID`) USING BTREE;

ALTER TABLE `hosts`
  ADD PRIMARY KEY (`ID`);

ALTER TABLE `requests`
  ADD PRIMARY KEY (`ID`),
  ADD KEY `host_ID` (`Host-ID`);

CREATE DEFINER=`root`@`%` EVENT `connections_cleanup` ON SCHEDULE EVERY 1 SECOND STARTS '2022-04-04 15:13:06' ON COMPLETION NOT PRESERVE ENABLE DO DELETE FROM connections WHERE `last-used` <  (NOW() - INTERVAL 10 SECOND);

CREATE DEFINER=`root`@`%` EVENT `hosts_status` ON SCHEDULE EVERY 1 SECOND STARTS '2022-04-04 15:13:47' ON COMPLETION NOT PRESERVE ENABLE DO UPDATE hosts SET `online`='0' WHERE `last-online` <  (NOW() - INTERVAL 10 SECOND);

CREATE DEFINER=`root`@`%` EVENT `requests_cleanup` ON SCHEDULE EVERY 1 SECOND STARTS '2022-04-04 15:12:51' ON COMPLETION NOT PRESERVE ENABLE DO DELETE FROM requests WHERE `last-used` <  (NOW() - INTERVAL 10 SECOND);

ALTER TABLE `connections`
  ADD CONSTRAINT `connections_ibfk_1` FOREIGN KEY (`Host-ID`) REFERENCES `hosts` (`ID`);

ALTER TABLE `requests`
  ADD CONSTRAINT `host_ID` FOREIGN KEY (`Host-ID`) REFERENCES `hosts` (`ID`) ON DELETE CASCADE;
COMMIT;

