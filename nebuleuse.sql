SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;


CREATE TABLE `neb_achievements` (
  `id` int(10) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `max` int(10) UNSIGNED NOT NULL,
  `fullName` varchar(255) NOT NULL,
  `fullDesc` varchar(255) NOT NULL,
  `icon` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

INSERT INTO `neb_achievements` (`id`, `name`, `max`, `fullName`, `fullDesc`, `icon`) VALUES
(1, 'test', 24, 'test', 'test', 'http://i.imgur.com/oyrwt3a.png');

CREATE TABLE `neb_config` (
  `name` varchar(255) NOT NULL,
  `value` varchar(255) NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

INSERT INTO `neb_config` (`name`, `value`) VALUES
('gameName', ''),
('gameVersion', '1'),
('updaterVersion', '1'),
('sessionTimeout', '1800'),
('autoRegister', 'true'),
('defaultAvatar', 'http://i.imgur.com/oyrwt3a.png'),
('productionBranch', 'master'),
('gitRepositoryPath', './'),
('updateSystem', 'GitPatch');

CREATE TABLE `neb_mirrors` (
  `id` smallint(5) UNSIGNED NOT NULL,
  `Address` varchar(255) NOT NULL,
  `Name` varchar(255) NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

CREATE TABLE `neb_sessions` (
  `userid` int(10) UNSIGNED NOT NULL,
  `lastAlive` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `sessionId` varchar(36) NOT NULL,
  `sessionStart` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00'
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

CREATE TABLE `neb_stats_tables` (
  `tableName` varchar(255) NOT NULL,
  `fields` text NOT NULL,
  `autoCount` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

INSERT INTO `neb_stats_tables` (`tableName`, `fields`, `autoCount`) VALUES
('users', '[{"Name":"timeplayed","Type":"int","Size":11}]', 0),
('kills', '[{"Name":"userid","Type":"int","Size":11},{"Name":"x","Type":"int","Size":11},{"Name":"y","Type":"int","Size":11},{"Name":"z","Type":"int","Size":11},{"Name":"map","Type":"string","Size":255},{"Name":"weapon","Type":"string","Size":255}]', 1);

CREATE TABLE `neb_updates` (
  `build` int(11) NOT NULL,
  `branch` varchar(255) NOT NULL,
  `size` int(11) NOT NULL,
  `rollback` tinyint(1) NOT NULL,
  `semver` varchar(255) NOT NULL,
  `log` text NOT NULL,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

INSERT INTO `neb_updates` (`build`, `branch`, `size`, `rollback`, `semver`, `log`, `date`) VALUES
(1, 'public', 2, 1, '', '', '2016-03-22 10:37:52'),
(2, 'public', 25, 0, '', '', '2016-03-22 10:37:52'),
(3, 'public', 5, 0, '', '', '2016-03-22 10:37:52');

CREATE TABLE `neb_updates_branches` (
  `name` varchar(255) NOT NULL,
  `rank` int(11) NOT NULL,
  `activeBuild` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

INSERT INTO `neb_updates_branches` (`name`, `rank`, `activeBuild`) VALUES
('public', 0, 1);

CREATE TABLE `neb_updates_builds` (
  `id` int(11) NOT NULL,
  `commit` varchar(255) NOT NULL,
  `log` text NOT NULL,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `changelist` text NOT NULL,
  `obselete` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

INSERT INTO `neb_updates_builds` (`id`, `commit`, `log`, `date`, `changelist`, `obselete`) VALUES
(1, '', 'Changed things', '2016-03-15 22:03:12', '', 0),
(2, '', 'again', '2016-03-15 22:11:24', '', 1),
(3, '', 'a', '2016-03-15 22:08:42', '', 0);

CREATE TABLE `neb_users` (
  `id` int(11) NOT NULL,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `rank` tinyint(3) NOT NULL,
  `avatars` varchar(255) NOT NULL,
  `hash` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

INSERT INTO `neb_users` (`id`, `username`, `password`, `rank`, `avatars`, `hash`) VALUES
(1, 'test', 'q4F_1BnvOQERMAtuwNHoocjO6DiHvt15ol2krqZ60v-NW-tb0_IooASPuZq6iv1tjjT60JIIhA1MZvjTcGhDqA==', 8, '', 'O-z_gcTHzvgoM3ndhFnKVbM-tUcnGZDz_o6mhkWFiL0VTnvCvHFVOBYnvBp23pbz1ZIafCoH_JO51gXlVkmf8w=='),
(2, 'test2', 'gx-5jrB9JIe2bc-Ou03lLhy4QYus0vjyuGdcFsnGItsvup4cxOFiBWO9h5S6uZFjiGQySWZKd4JSxHWjKW6EZQ==', 1, '', 'fJ_EfI7QylpTJTnrq_9hxzTGcUOHawhBUgDSoUkBVkXUubzAyIeefSTjXlK5DBru3G37wxczGVd_ILkkvkDvqA==');

CREATE TABLE `neb_users_achievements` (
  `userid` int(10) UNSIGNED NOT NULL,
  `achievementid` int(10) UNSIGNED NOT NULL,
  `progress` int(10) UNSIGNED NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

INSERT INTO `neb_users_achievements` (`userid`, `achievementid`, `progress`) VALUES
(2, 1, 24);

CREATE TABLE `neb_users_stats` (
  `userid` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `value` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

INSERT INTO `neb_users_stats` (`userid`, `name`, `value`) VALUES
(2, 'kills', 13);

CREATE TABLE `neb_users_stats_kills` (
  `userid` int(11) DEFAULT NULL,
  `x` int(11) DEFAULT NULL,
  `y` int(11) DEFAULT NULL,
  `z` int(11) DEFAULT NULL,
  `map` varchar(255) DEFAULT NULL,
  `weapon` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

INSERT INTO `neb_users_stats_kills` (`userid`, `x`, `y`, `z`, `map`, `weapon`) VALUES
(2, 5, 5, 5, 'test', 'Flower'),
(2, 5, 5, 5, 'test', 'Flower');


ALTER TABLE `neb_achievements`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `neb_mirrors`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `neb_sessions`
  ADD PRIMARY KEY (`userid`);

ALTER TABLE `neb_updates`
  ADD PRIMARY KEY (`build`),
  ADD KEY `branch` (`branch`),
  ADD KEY `build` (`build`),
  ADD KEY `branch_2` (`branch`);

ALTER TABLE `neb_updates_branches`
  ADD UNIQUE KEY `name` (`name`);

ALTER TABLE `neb_updates_builds`
  ADD PRIMARY KEY (`id`),
  ADD KEY `id` (`id`),
  ADD KEY `id_2` (`id`);

ALTER TABLE `neb_users`
  ADD PRIMARY KEY (`id`);


ALTER TABLE `neb_achievements`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;
ALTER TABLE `neb_updates_builds`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;
ALTER TABLE `neb_users`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
