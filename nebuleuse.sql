SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;


CREATE TABLE IF NOT EXISTS `neb_achievements` (
`id` int(10) unsigned NOT NULL,
  `name` varchar(255) NOT NULL,
  `max` int(10) unsigned NOT NULL,
  `fullName` varchar(255) NOT NULL,
  `fullDesc` varchar(255) NOT NULL,
  `icon` varchar(255) NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;

INSERT INTO `neb_achievements` (`id`, `name`, `max`, `fullName`, `fullDesc`, `icon`) VALUES
(1, 'test', 24, 'test', 'test', 'http://i.imgur.com/oyrwt3a.png');

CREATE TABLE IF NOT EXISTS `neb_config` (
  `name` varchar(255) NOT NULL,
  `value` varchar(255) NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

INSERT INTO `neb_config` (`name`, `value`) VALUES
('maintenance', '0'),
('maintenanceMessage', 'Main server is undergoing maintenance.'),
('gameVersion', '1'),
('sessionTimeout', '120'),
('gameName', ''),
('updaterVersion', '1'),
('autoRegister', 'true'),
('defaultAvatar', 'http://i.imgur.com/oyrwt3a.png');

CREATE TABLE IF NOT EXISTS `neb_mirrors` (
`id` smallint(5) unsigned NOT NULL,
  `Address` varchar(255) NOT NULL,
  `Name` varchar(255) NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS `neb_sessions` (
  `userid` int(10) unsigned NOT NULL,
  `lastAlive` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `sessionId` varchar(36) NOT NULL,
  `sessionStart` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00'
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS `neb_stats_tables` (
  `tableName` varchar(255) NOT NULL,
  `fields` text NOT NULL,
  `autoCount` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

INSERT INTO `neb_stats_tables` (`tableName`, `fields`, `autoCount`) VALUES
('users', '[{"Name":"timeplayed","Type":"int","Size":11}]', 0),
('kills', '[{"Name":"userid","Type":"int","Size":11},{"Name":"x","Type":"int","Size":11},{"Name":"y","Type":"int","Size":11},{"Name":"z","Type":"int","Size":11},{"Name":"map","Type":"string","Size":255},{"Name":"weapon","Type":"string","Size":255}]', 1);

CREATE TABLE IF NOT EXISTS `neb_updates` (
  `version` int(11) NOT NULL,
  `log` text NOT NULL,
  `size` int(11) NOT NULL,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `commit` varchar(255) NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS `neb_users` (
`id` int(11) NOT NULL,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `rank` tinyint(3) NOT NULL,
  `avatars` varchar(255) NOT NULL,
  `hash` varchar(255) NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=latin1;

INSERT INTO `neb_users` (`id`, `username`, `password`, `rank`, `avatars`, `hash`) VALUES
(1, 'test', 'q4F_1BnvOQERMAtuwNHoocjO6DiHvt15ol2krqZ60v-NW-tb0_IooASPuZq6iv1tjjT60JIIhA1MZvjTcGhDqA==', 2, '', 'O-z_gcTHzvgoM3ndhFnKVbM-tUcnGZDz_o6mhkWFiL0VTnvCvHFVOBYnvBp23pbz1ZIafCoH_JO51gXlVkmf8w==');

CREATE TABLE IF NOT EXISTS `neb_users_achievements` (
  `userid` int(10) unsigned NOT NULL,
  `achievementid` int(10) unsigned NOT NULL,
  `progress` int(10) unsigned NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

INSERT INTO `neb_users_achievements` (`userid`, `achievementid`, `progress`) VALUES
(2, 1, 24);

CREATE TABLE IF NOT EXISTS `neb_users_stats` (
  `userid` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `value` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

INSERT INTO `neb_users_stats` (`userid`, `name`, `value`) VALUES
(2, 'kills', 11);

CREATE TABLE IF NOT EXISTS `neb_users_stats_kills` (
  `userid` int(11) DEFAULT NULL,
  `x` int(11) DEFAULT NULL,
  `y` int(11) DEFAULT NULL,
  `z` int(11) DEFAULT NULL,
  `map` varchar(255) DEFAULT NULL,
  `weapon` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


ALTER TABLE `neb_achievements`
 ADD PRIMARY KEY (`id`);

ALTER TABLE `neb_mirrors`
 ADD PRIMARY KEY (`id`);

ALTER TABLE `neb_sessions`
 ADD PRIMARY KEY (`userid`);

ALTER TABLE `neb_updates`
 ADD PRIMARY KEY (`version`);

ALTER TABLE `neb_users`
 ADD PRIMARY KEY (`id`);


ALTER TABLE `neb_achievements`
MODIFY `id` int(10) unsigned NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=0;
ALTER TABLE `neb_users`
MODIFY `id` int(11) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=0;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
