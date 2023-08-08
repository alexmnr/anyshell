SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";

CREATE TABLE `connections` (
  `ID` int(11) NOT NULL,
  `HostID` int(11) NOT NULL,
  `ServerPort` int(6) NOT NULL,
  `LastUsed` datetime NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `hosts` (
  `ID` int(4) NOT NULL,
  `Name` varchar(20) NOT NULL,
  `User` varchar(20) NOT NULL,
  `Port` int(6) NOT NULL,
  `PublicIP` varchar(15) DEFAULT NULL,
  `LocalIP` varchar(15) DEFAULT NULL,
  `Online` int(1) NOT NULL,
  `LastOnline` datetime NOT NULL DEFAULT current_timestamp(),
  `Version` int(5) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `requests` (
  `ID` int(3) NOT NULL,
  `HostID` int(11) NOT NULL,
  `LastUsed` datetime NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


ALTER TABLE `connections`
  ADD PRIMARY KEY (`ID`),
  ADD KEY `HostID` (`HostID`) USING BTREE;

ALTER TABLE `hosts`
  ADD PRIMARY KEY (`ID`);

ALTER TABLE `requests`
  ADD PRIMARY KEY (`ID`),
  ADD KEY `HostID` (`HostID`);

CREATE DEFINER=`root`@`%` EVENT `ConnectionsCleanup` ON SCHEDULE EVERY 1 SECOND STARTS '2022-04-04 15:13:06' ON COMPLETION NOT PRESERVE ENABLE DO DELETE FROM connections WHERE `LastUsed` <  (NOW() - INTERVAL 10 SECOND);

CREATE DEFINER=`root`@`%` EVENT `HostsStatus` ON SCHEDULE EVERY 1 SECOND STARTS '2022-04-04 15:13:47' ON COMPLETION NOT PRESERVE ENABLE DO UPDATE hosts SET `Online`='0' WHERE `LastOnline` <  (NOW() - INTERVAL 10 SECOND);

CREATE DEFINER=`root`@`%` EVENT `RequestsCleanup` ON SCHEDULE EVERY 1 SECOND STARTS '2022-04-04 15:12:51' ON COMPLETION NOT PRESERVE ENABLE DO DELETE FROM requests WHERE `LastUsed` <  (NOW() - INTERVAL 10 SECOND);

ALTER TABLE `connections`
  ADD CONSTRAINT `connections_ibfk_1` FOREIGN KEY (`HostID`) REFERENCES `hosts` (`ID`);

ALTER TABLE `requests`
  ADD CONSTRAINT `HostID` FOREIGN KEY (`HostID`) REFERENCES `hosts` (`ID`) ON DELETE CASCADE;
COMMIT;

